package vcd

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/vmware/go-vcloud-director/v2/govcd"
)

func GetSamlMetadata(client *govcd.VCDClient, org string) (*EntityDescriptor, error) {
	// https://192.168.1.160/cloud/org/my-org/saml/metadata/alias/vcd
	// urlString := https://
	url := client.Client.VCDHREF
	SamlMetadataUrl := url.Scheme + "://" + url.Host + "/cloud/org/" + org + "/saml/metadata/alias/vcd"
	resp, err := client.Client.Http.Get(SamlMetadataUrl)
	if err != nil {
		return nil, fmt.Errorf("could not fetch SAML metadata from URL %s: %s", SamlMetadataUrl, err)
	}

	allBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("unable to read body for SAML metadata: %s", err)
	}

	metadata := EntityDescriptor{}
	err = xml.Unmarshal(allBody, &metadata)
	if err != nil {
		return nil, fmt.Errorf("unable to unmarshal SAML metadata: %s", err)
	}
	return &metadata, nil
}

func GetSamlEntityId(client *govcd.VCDClient, org string) (string, error) {
	metadata, err := GetSamlMetadata(client, org)
	if err != nil {
		return "", fmt.Errorf("unable to get SAML entity ID: %s", err)
	}

	return metadata.EntityID, nil
}

func ProviderAuthenticateSaml(client *govcd.VCDClient, domainUser, domainPassword, org string) error {
	url := client.Client.VCDHREF

	backupRedirectChecker := client.Client.Http.CheckRedirect
	// Restore client at the end of retrieval
	defer func() {
		client.Client.Http.CheckRedirect = backupRedirectChecker
	}()

	// Patch http client to avoid following redirects
	client.Client.Http.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	loginUrl := url.Scheme + "://" + url.Host + "/login/my-org/saml/login/alias/vcd?service=tenant:" + org
	log.Printf("[DEBUG] SAML looking up IdP host redirect in: %s", loginUrl)
	redirect, err := client.Client.Http.Get(loginUrl)
	if err != nil {
		return fmt.Errorf("SAML unable to get IdP address via %s: %s", loginUrl, err)
	}
	adfsEndpoint, err := redirect.Location()
	if err != nil {
		return fmt.Errorf("SAML no redirect location for %s found: %s", loginUrl, err)
	}

	authEndPoint := adfsEndpoint.Scheme + "://" + adfsEndpoint.Hostname() + "/adfs/services/trust/13/usernamemixed"
	log.Printf("[DEBUG] SAML got IdP login endpoint: %s", authEndPoint)

	samlEntityId, err := GetSamlEntityId(client, org)
	if err != nil {
		return fmt.Errorf("unable to get SAML entity ID: %s", err)
	}
	log.Printf("[DEBUG] SAML got entity ID: %s", samlEntityId)

	requestBody := getSamlRequestBody(domainUser, domainPassword, samlEntityId)

	body := strings.NewReader(requestBody)
	log.Printf("[DEBUG] SAML posting login data to IdP: %s", authEndPoint)

	resp, err := client.Client.Http.Post(authEndPoint, "application/soap+xml", body)
	if err != nil {
		return fmt.Errorf("SAML got error after posting login data to IdP %s: %s", authEndPoint, err)
	}
	if resp.StatusCode != http.StatusOK {
		errorBody, _ := ioutil.ReadAll(resp.Body)
		log.Printf("[ERROR] SAML authentication error against IdP '%s' for entity with ID '%s': %s",
			authEndPoint, samlEntityId, errorBody)

		errParse := ErrorEnvelope{}
		_ = xml.Unmarshal(errorBody, &errParse)
		return fmt.Errorf("could not SAML authenticate against IdP endpoint '%s' for entity with ID '%s'. Got status code %d and response %s",
			authEndPoint, samlEntityId, resp.StatusCode, errParse.Body.Fault.Reason.Text)
	}

	authResponseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	responseStruct := AuthResponseEnvelope{}
	err = xml.Unmarshal(authResponseBody, &responseStruct)
	if err != nil {
		return fmt.Errorf("unable to unmarshal response body: %s", err)
	}

	tokenString := responseStruct.Body.RequestSecurityTokenResponseCollection.RequestSecurityTokenResponse.RequestedSecurityTokenTxt.Text
	base64GzippedToken, err := gzipAndBase64Encode(tokenString)
	if err != nil {
		return fmt.Errorf("SAML error encoding SIGN token: %s", err)
	}

	log.Printf("[DEBUG] SAML got SIGN token from IdP '%s' for entity with ID '%s'",
		authEndPoint, samlEntityId)

	req, err := http.NewRequest(http.MethodPost, url.Scheme+"://"+url.Host+"/api/sessions", nil)
	if err != nil {
		return fmt.Errorf("error making new request with SAML SIGN token to %s: %s", req.URL.String(), err)
	}
	req.Header.Add("Accept", "application/*+xml;version="+client.Client.APIVersion)
	req.Header.Add("Authorization", `SIGN token="`+base64GzippedToken+`",org="`+org+`"`)

	resp, err = client.Client.Http.Do(req)
	if err != nil {
		return fmt.Errorf("error submitting SIGN token for authentication to %s: %s",
			req.URL.String(), err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("got unexpected response status code %d", resp.StatusCode)
	}

	accessToken := resp.Header.Get("X-Vcloud-Authorization")
	log.Printf("[DEBUG] SAML settings access token for further requests")
	err = client.SetToken(org, govcd.AuthorizationHeader, accessToken)
	if err != nil {
		return fmt.Errorf("error during token-based authentication: %s", err)
	}

	return nil
}

func gzipAndBase64Encode(token string) (string, error) {
	var gzipBuffer bytes.Buffer
	gz := gzip.NewWriter(&gzipBuffer)
	if _, err := gz.Write([]byte(token)); err != nil {
		return "", fmt.Errorf("error writing to gzip buffer: %s", err)
	}
	if err := gz.Close(); err != nil {
		return "", fmt.Errorf("error closing gzip buffer: %s", err)
	}
	base64GzippedToken := base64.StdEncoding.EncodeToString(gzipBuffer.Bytes())

	return base64GzippedToken, nil
}

func getSamlRequestBody(user, pass, samlEntityIderence string) string {
	return `<s:Envelope xmlns:s="http://www.w3.org/2003/05/soap-envelope" xmlns:a="http://www.w3.org/2005/08/addressing" xmlns:u="http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-wssecurity-utility-1.0.xsd">
	<s:Header>
	  <a:Action s:mustUnderstand="1">http://docs.oasis-open.org/ws-sx/ws-trust/200512/RST/Issue</a:Action>
	  <a:ReplyTo>
		<a:Address>http://www.w3.org/2005/08/addressing/anonymous</a:Address>
	  </a:ReplyTo>
	  <a:To s:mustUnderstand="1">https://win-60g606n0afg.test-forest.net/adfs/services/trust/13/usernamemixed</a:To>
	  <o:Security s:mustUnderstand="1" xmlns:o="http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-wssecurity-secext-1.0.xsd">
		<u:Timestamp u:Id="_0">
		  <u:Created>` + time.Now().Format(time.RFC3339) + `</u:Created>
		  <u:Expires>` + time.Now().Add(1*time.Minute).Format(time.RFC3339) + `</u:Expires>
		</u:Timestamp>
		<o:UsernameToken>
		  <o:Username>` + user + `</o:Username>
		  <o:Password o:Type="http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-username-token-profile-1.0#PasswordText">` + pass + `</o:Password>
		</o:UsernameToken>
	  </o:Security>
	</s:Header>
	<s:Body>
	  <trust:RequestSecurityToken xmlns:trust="http://docs.oasis-open.org/ws-sx/ws-trust/200512">
		<wsp:AppliesTo xmlns:wsp="http://schemas.xmlsoap.org/ws/2004/09/policy">
		  <a:samlEntityIderence>
			<a:Address>` + samlEntityIderence + `</a:Address>
		  </a:samlEntityIderence>
		</wsp:AppliesTo>
		<trust:KeySize>0</trust:KeySize>
		<trust:KeyType>http://docs.oasis-open.org/ws-sx/ws-trust/200512/Bearer</trust:KeyType>
		<i:RequestDisplayToken xml:lang="en" xmlns:i="http://schemas.xmlsoap.org/ws/2005/05/identity" />
		<trust:RequestType>http://docs.oasis-open.org/ws-sx/ws-trust/200512/Issue</trust:RequestType>
		<trust:TokenType>http://docs.oasis-open.org/wss/oasis-wss-saml-token-profile-1.1#SAMLV2.0</trust:TokenType>
	  </trust:RequestSecurityToken>
	</s:Body>
  </s:Envelope>
  `
}
