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

	loginUrl := url.Scheme + "://" + url.Host + "/tenant/" + org
	log.Println("DAINIUS nexthop ", loginUrl)
	resp, err := client.Client.Http.Get(loginUrl)
	if err != nil {
		return err
	}
	// https://192.168.1.160/login?service=tenant:my-org&redirectTo=%2Ftenant%2Fmy-org%2F
	nextHop, _ := resp.Location()
	log.Println("DAINIUS nexthop ", nextHop)
	resp, err = client.Client.Http.Get(nextHop.String())
	if err != nil {
		return err
	}

	// https://192.168.1.160/login/?service=tenant:my-org
	nextHop, _ = resp.Location()
	log.Println("DAINIUS nexthop ", nextHop)
	resp, err = client.Client.Http.Get(nextHop.String())
	if err != nil {
		return err
	}

	// https://192.168.1.160/login/my-org/saml/login/alias/vcd?service=tenant:my-org
	nextHop, _ = resp.Location()
	log.Println("DAINIUS nexthop ", nextHop)
	// https://win-60g606n0afg.test-forest.net/adfs/ls/?SAMLRequest=fZFNT8MwDIbv%2FIoq9zYflDKitdMEQkICCcHgwM1K3S5Tm4w4HfDvyQYTcOESK9Lrx86T%2BeJ9HLIdBrLe1UwWgmXojG%2Bt62v2tLrOZ2zRnMwJxkFt9XKKa%2FeArxNSzJZEGGLqu%2FSOphHDI4adNXjjWnyvWSJdpZh1EA%2FsdYxb0py%2FWZdXoq9E5QR0fRFTKO98SKVwGDm0HfGBOMuufTB4GFmzDgZClt1c1QxKs5mhrJSFslRdrzYbgWDkWQVnQp6mEN0Dkd3hTxvRlPaiCC7WTAklclHmQq7EuVZKy4uinJ2%2BsOz5KELtRSQ1jvTX02s2Bac9kCXtYETS0ejH5d2tTlG9DT564wfWfJnSh4HhN%2BF%2FABxdsuboSV6oQlazQqZTcDP4qeU%2B9Hz8yPdlD%2BUjRmghAofBAvGdaef89wLN9%2FXvzzWf&RelayState=aHR0cHM6Ly8xOTIuMTY4LjEuMTYwL3RlbmFudC9teS1vcmc%3D&SigAlg=http%3A%2F%2Fwww.w3.org%2F2000%2F09%2Fxmldsig%23rsa-sha1&Signature=W3mUfGiecEJudJLqIV%2F2cFuJcbPiQxgayxVbJf6hOMp8ZQcqG01NR1Rm3qTqaote5dSkprw42dVOMIHdeiJL1g7%2FW9ON6%2BvzJHvL3rdy652%2BeSv6q0r9wDJ8eKC5DpwcmW0UUATHHt4ENMPa6w6MgE2Mwm1F6eYu1c5CcIC306lzQNiSWwNA08frX1wxl3RtSrrm9qo9K9UoQOAULkYjAgghI65Dr%2BEWjiu%2FYgVw1SuXMKRQcQ1Q8MQ2uhDqjlfXuO3Fnp582zLh1uMx1ZiFO1LPaqTT7K%2BBvlUISzkCe6YkSOlr%2Fz08t7A4fkMxmhnTd2gBu6WhzcUUWkqLMzSwgw%3D%3D

	resp, err = client.Client.Http.Get(nextHop.String())
	if err != nil {
		return err
	}

	nextHop, _ = resp.Location()

	log.Println("DAINIUS nexthop ", nextHop)

	authEndPoint := nextHop.Scheme + "://" + nextHop.Hostname() + "/adfs/services/trust/13/usernamemixed"
	// samlEntityId := url.Scheme + "://" + url.Host + "/cloud/org/my-org/saml/metadata/alias/vcd"
	// samlEntityId := "asd"

	samlEntityId, err := GetSamlEntityId(client, org)
	if err != nil {
		return fmt.Errorf("unable to get SAML entity ID: %s", err)
	}
	log.Printf("[DEBUG] got SAML entity ID: %s", samlEntityId)

	requestBody := getSamlRequestBody(domainUser, domainPassword, samlEntityId)
	body := strings.NewReader(requestBody)

	log.Printf("[DEBUG] Posting login data to IdP: %s", authEndPoint)

	resp, err = client.Client.Http.Post(authEndPoint, "application/soap+xml", body)
	if err != nil {
		return err
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

	log.Println("DAINIUS response status: ", resp.StatusCode)
	log.Printf("DAINIUS response %+#v", resp)

	r, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	log.Println("DAINIUS auth response: ", string(r))

	responseStruct := ResponseEnvelope{}

	err = xml.Unmarshal(r, &responseStruct)
	if err != nil {
		return fmt.Errorf("unable to unmarshal response body: %s", err)
	}

	tokenPart := responseStruct.Body.RequestSecurityTokenResponseCollection.RequestSecurityTokenResponse.RequestedSecurityTokenTxt.Text
	log.Printf("DAINIUS EncryptedAssertion %+#v", tokenPart)
	log.Printf("DAINIUS encrypted assertion: %s", tokenPart)

	// Gzip and base64 encode security token as per endpoint requirements
	var gzipBuffer bytes.Buffer
	gz := gzip.NewWriter(&gzipBuffer)
	if _, err := gz.Write([]byte(tokenPart)); err != nil {
		log.Fatal(err)
	}
	if err := gz.Close(); err != nil {
		log.Fatal(err)
	}
	base64GzippedToken := base64.StdEncoding.EncodeToString(gzipBuffer.Bytes())

	log.Printf("DAINIUS encoded test %s", base64GzippedToken)

	/// # Got data, try to authenticate against vCD

	req, err := http.NewRequest(http.MethodPost, url.Scheme+"://"+url.Host+"/api/sessions", nil)
	if err != nil {
		return fmt.Errorf("error posting: %s", err)
	}
	req.Header.Add("Accept", "application/*+xml;version=29.0")
	req.Header.Add("Authorization", `SIGN token="`+base64GzippedToken+`",org="`+org+`"`)

	log.Printf("DAINIUS headers: %+#v", req.Header)
	log.Printf("DAINIUS sign: %s", req.Header.Get("Authorization"))

	resp, err = client.Client.Http.Do(req)
	if err != nil {
		return err
	}

	iii, _ := ioutil.ReadAll(resp.Body)
	log.Printf("DAINIUS status code: %d", resp.StatusCode)
	log.Printf("DAINIUS body: %s", string(iii))
	log.Printf("DAINIUS %+#v", resp.Header)

	accessToken := resp.Header.Get("X-Vcloud-Authorization")
	log.Printf("DAINIUS %s", accessToken)

	err = client.SetToken(org, govcd.AuthorizationHeader, accessToken)
	if err != nil {
		return fmt.Errorf("error during token-based authentication: %s", err)
	}

	return nil
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

const requestBodyTemplate = `
<s:Envelope xmlns:s="http://www.w3.org/2003/05/soap-envelope" xmlns:a="http://www.w3.org/2005/08/addressing" xmlns:u="http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-wssecurity-utility-1.0.xsd">
  <s:Header>
    <a:Action s:mustUnderstand="1">http://docs.oasis-open.org/ws-sx/ws-trust/200512/RST/Issue</a:Action>
    <a:ReplyTo>
      <a:Address>http://www.w3.org/2005/08/addressing/anonymous</a:Address>
    </a:ReplyTo>
    <a:To s:mustUnderstand="1">{{.DestinationUrl}}</a:To>
    <o:Security s:mustUnderstand="1" xmlns:o="http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-wssecurity-secext-1.0.xsd">
      <u:Timestamp u:Id="_0">
        <u:Created>{{.TsUtcNow}}</u:Created>
        <u:Expires>{{.TsUtcThen}}</u:Expires>
      </u:Timestamp>
      <o:UsernameToken>
        <o:Username>{{.Username}}</o:Username>
        <o:Password o:Type="http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-username-token-profile-1.0#PasswordText">{{.Password}}</o:Password>
      </o:UsernameToken>
    </o:Security>
  </s:Header>
  <s:Body>
    <trust:RequestSecurityToken xmlns:trust="http://docs.oasis-open.org/ws-sx/ws-trust/200512">
      <wsp:AppliesTo xmlns:wsp="http://schemas.xmlsoap.org/ws/2004/09/policy">
        <a:samlEntityIderence>
          <a:Address>{{.EntityId}}</a:Address>
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
