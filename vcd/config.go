package vcd

import (
	"bytes"
	"compress/gzip"
	"crypto/sha1"
	"encoding/base64"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/mutexkv"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/vmware/go-vcloud-director/v2/govcd"
)

type Config struct {
	User            string
	Password        string
	Token           string // Token used instead of user and password
	SysOrg          string // Org used for authentication
	Org             string // Default Org used for API operations
	Vdc             string // Default (optional) VDC for API operations
	Href            string
	MaxRetryTimeout int
	InsecureFlag    bool
}

type VCDClient struct {
	*govcd.VCDClient
	SysOrg          string
	Org             string // name of default Org
	Vdc             string // name of default VDC
	MaxRetryTimeout int
	InsecureFlag    bool
}

// Type used to simplify reading resource definitions
type StringMap map[string]interface{}

const (
	// Most common error messages in the library

	// Used when a call to GetOrgAndVdc fails. The placeholder is for the error
	errorRetrievingOrgAndVdc = "error retrieving Org and VDC: %s"

	// Used when a call to GetOrgAndVdc fails. The placeholders are for vdc, org, and the error
	errorRetrievingVdcFromOrg = "error retrieving VDC %s from Org %s: %s"

	// Used when we can't get a valid edge gateway. The placeholder is for the error
	errorUnableToFindEdgeGateway = "unable to find edge gateway: %s"

	// Used when a task fails. The placeholder is for the error
	errorCompletingTask = "error completing tasks: %s"

	// Used when a call to GetAdminOrgFromResource fails. The placeholder is for the error
	errorRetrievingOrg = "error retrieving Org: %s"
)

// Cache values for vCD connection.
// When the Client() function is called with the same parameters, it will return
// a cached value instead of connecting again.
// This makes the Client() function both deterministic and fast.
type cachedConnection struct {
	initTime   time.Time
	connection *VCDClient
}

type cacheStorage struct {
	// conMap holds cached VDC authenticated connection
	conMap map[string]cachedConnection
	// cacheClientServedCount records how many times we have cached a connection
	cacheClientServedCount int
	sync.Mutex
}

var (
	// Enables the caching of authenticated connections
	enableConnectionCache bool = os.Getenv("VCD_CACHE") != ""

	// Cached VDC authenticated connection
	cachedVCDClients = &cacheStorage{conMap: make(map[string]cachedConnection)}

	// Invalidates the cache after a given time (connection tokens usually expire after 20 to 30 minutes)
	maxConnectionValidity time.Duration = 20 * time.Minute

	enableDebug bool = os.Getenv("GOVCD_DEBUG") != ""
	enableTrace bool = os.Getenv("GOVCD_TRACE") != ""

	// Separation string used for import operations
	// Can be changed usin either "import_separator" property in Provider
	// or environment variable "VCD_IMPORT_SEPARATOR"
	ImportSeparator = "."
)

// Displays conditional messages
func debugPrintf(format string, args ...interface{}) {
	// When GOVCD_TRACE is enabled, we also display the function that generated the message
	if enableTrace {
		format = fmt.Sprintf("[%s] %s", filepath.Base(callFuncName()), format)
	}
	// The formatted message passed to this function is displayed only when GOVCD_DEBUG is enabled.
	if enableDebug {
		fmt.Printf(format, args...)
	}
}

// This is a global MutexKV for all resources
var vcdMutexKV = mutexkv.NewMutexKV()

func (cli *VCDClient) lockVapp(d *schema.ResourceData) {
	vappName := d.Get("name").(string)
	if vappName == "" {
		panic("vApp name not found")
	}
	key := fmt.Sprintf("org:%s|vdc:%s|vapp:%s", cli.getOrgName(d), cli.getVdcName(d), vappName)
	vcdMutexKV.Lock(key)
}

func (cli *VCDClient) unLockVapp(d *schema.ResourceData) {
	vappName := d.Get("name").(string)
	if vappName == "" {
		panic("vApp name not found")
	}
	key := fmt.Sprintf("org:%s|vdc:%s|vapp:%s", cli.getOrgName(d), cli.getVdcName(d), vappName)
	vcdMutexKV.Unlock(key)
}

// locks an edge gateway resource
// Differs from lockParentEdgeGtw in the resource name. When EGW is the parent,
// it's named "edge_gateway". When it's the main resource, it's found at "name"
func (cli *VCDClient) lockEdgeGateway(d *schema.ResourceData) {
	edgeGatewayName := d.Get("name").(string)
	if edgeGatewayName == "" {
		panic("edge gateway name not found")
	}
	key := fmt.Sprintf("org:%s|vdc:%s|edge:%s", cli.getOrgName(d), cli.getVdcName(d), edgeGatewayName)
	vcdMutexKV.Lock(key)
}

// unlocks an edge gateway resource
// Differs from unlockParentEdgeGtw in the resource name. When EGW is the parent,
// it's named "edge_gateway". When it's the main resource, it's found at "name"
func (cli *VCDClient) unlockEdgeGateway(d *schema.ResourceData) {
	edgeGatewayName := d.Get("name").(string)
	if edgeGatewayName == "" {
		panic("edge gateway name not found")
	}
	key := fmt.Sprintf("org:%s|vdc:%s|edge:%s", cli.getOrgName(d), cli.getVdcName(d), edgeGatewayName)
	vcdMutexKV.Unlock(key)
}

// function lockParentVapp locks using vapp_name name existing in resource parameters.
// Parent means the resource belongs to the vApp being locked
func (cli *VCDClient) lockParentVapp(d *schema.ResourceData) {
	vappName := d.Get("vapp_name").(string)
	if vappName == "" {
		panic("vApp name not found")
	}
	key := fmt.Sprintf("org:%s|vdc:%s|vapp:%s", cli.getOrgName(d), cli.getVdcName(d), vappName)
	vcdMutexKV.Lock(key)
}

func (cli *VCDClient) unLockParentVapp(d *schema.ResourceData) {
	vappName := d.Get("vapp_name").(string)
	if vappName == "" {
		panic("vApp name not found")
	}
	key := fmt.Sprintf("org:%s|vdc:%s|vapp:%s", cli.getOrgName(d), cli.getVdcName(d), vappName)
	vcdMutexKV.Unlock(key)
}

// lockParentVm locks using vapp_name and vm_name names existing in resource parameters.
// Parent means the resource belongs to the VM being locked
func (cli *VCDClient) lockParentVm(d *schema.ResourceData) {
	vappName := d.Get("vapp_name").(string)
	if vappName == "" {
		panic("vApp name not found")
	}
	vmName := d.Get("vm_name").(string)
	if vmName == "" {
		panic("vmName name not found")
	}
	key := fmt.Sprintf("org:%s|vdc:%s|vapp:%s|vm:%s", cli.getOrgName(d), cli.getVdcName(d), vappName, vmName)
	vcdMutexKV.Lock(key)
}

func (cli *VCDClient) unLockParentVm(d *schema.ResourceData) {
	vappName := d.Get("vapp_name").(string)
	if vappName == "" {
		panic("vApp name not found")
	}
	vmName := d.Get("vm_name").(string)
	if vmName == "" {
		panic("vmName name not found")
	}
	key := fmt.Sprintf("org:%s|vdc:%s|vapp:%s|vm:%s", cli.getOrgName(d), cli.getVdcName(d), vappName, vmName)
	vcdMutexKV.Unlock(key)
}

// function lockParentEdgeGtw locks using edge_gateway name existing in resource parameters.
// Parent means the resource belongs to the edge gateway being locked
func (cli *VCDClient) lockParentEdgeGtw(d *schema.ResourceData) {
	edgeGtwName := d.Get("edge_gateway").(string)
	if edgeGtwName == "" {
		panic("edge gateway not found")
	}
	key := fmt.Sprintf("org:%s|vdc:%s|edge:%s", cli.getOrgName(d), cli.getVdcName(d), edgeGtwName)
	vcdMutexKV.Lock(key)
}

func (cli *VCDClient) unLockParentEdgeGtw(d *schema.ResourceData) {
	edgeGtwName := d.Get("edge_gateway").(string)
	if edgeGtwName == "" {
		panic("edge gateway not found")
	}
	key := fmt.Sprintf("org:%s|vdc:%s|edge:%s", cli.getOrgName(d), cli.getVdcName(d), edgeGtwName)
	vcdMutexKV.Unlock(key)
}

func (cli *VCDClient) getOrgName(d *schema.ResourceData) string {
	orgName := d.Get("org").(string)
	if orgName == "" {
		orgName = cli.Org
	}
	return orgName
}

func (cli *VCDClient) getVdcName(d *schema.ResourceData) string {
	orgName := d.Get("vdc").(string)
	if orgName == "" {
		orgName = cli.Vdc
	}
	return orgName
}

// GetOrgAndVdc finds a pair of org and vdc using the names provided
// in the args. If the names are empty, it will use the default
// org and vdc names from the provider.
func (cli *VCDClient) GetOrgAndVdc(orgName, vdcName string) (org *govcd.Org, vdc *govcd.Vdc, err error) {

	if orgName == "" {
		orgName = cli.Org
	}
	if orgName == "" {
		return nil, nil, fmt.Errorf("empty Org name provided")
	}
	if vdcName == "" {
		vdcName = cli.Vdc
	}
	if vdcName == "" {
		return nil, nil, fmt.Errorf("empty VDC name provided")
	}
	org, err = cli.VCDClient.GetOrgByName(orgName)
	if err != nil {
		return nil, nil, fmt.Errorf("error retrieving Org %s: %s", orgName, err)
	}
	if org.Org.Name == "" || org.Org.HREF == "" || org.Org.ID == "" {
		return nil, nil, fmt.Errorf("empty Org %s found ", orgName)
	}
	vdc, err = org.GetVDCByName(vdcName, false)
	if err != nil {
		return nil, nil, fmt.Errorf("error retrieving VDC %s: %s", vdcName, err)
	}
	if vdc == nil || vdc.Vdc.ID == "" || vdc.Vdc.HREF == "" || vdc.Vdc.Name == "" {
		return nil, nil, fmt.Errorf("error retrieving VDC %s: not found", vdcName)
	}
	return org, vdc, err
}

// GetAdminOrg finds org using the names provided in the args.
// If the name is empty, it will use the default
// org name from the provider.
func (cli *VCDClient) GetAdminOrg(orgName string) (org *govcd.AdminOrg, err error) {

	if orgName == "" {
		orgName = cli.Org
	}
	if orgName == "" {
		return nil, fmt.Errorf("empty Org name provided")
	}

	org, err = cli.VCDClient.GetAdminOrgByName(orgName)
	if err != nil {
		return nil, fmt.Errorf("error retrieving Org %s: %s", orgName, err)
	}
	if org.AdminOrg.Name == "" || org.AdminOrg.HREF == "" || org.AdminOrg.ID == "" {
		return nil, fmt.Errorf("empty org %s found", orgName)
	}
	return org, err
}

// Same as GetOrgAndVdc, but using data from the resource, if available.
func (cli *VCDClient) GetOrgAndVdcFromResource(d *schema.ResourceData) (org *govcd.Org, vdc *govcd.Vdc, err error) {
	orgName := d.Get("org").(string)
	vdcName := d.Get("vdc").(string)
	return cli.GetOrgAndVdc(orgName, vdcName)
}

// Same as GetOrgAndVdc, but using data from the resource, if available.
func (cli *VCDClient) GetAdminOrgFromResource(d *schema.ResourceData) (org *govcd.AdminOrg, err error) {
	orgName := d.Get("org").(string)
	return cli.GetAdminOrg(orgName)
}

// Gets an edge gateway when you don't need org or vdc for other purposes
func (cli *VCDClient) GetEdgeGateway(orgName, vdcName, edgeGwName string) (eg *govcd.EdgeGateway, err error) {

	if edgeGwName == "" {
		return nil, fmt.Errorf("empty Edge Gateway name provided")
	}
	_, vdc, err := cli.GetOrgAndVdc(orgName, vdcName)
	if err != nil {
		return nil, fmt.Errorf("error retrieving org and vdc: %s", err)
	}
	eg, err = vdc.GetEdgeGatewayByName(edgeGwName, true)

	if err != nil {
		if os.Getenv("GOVCD_DEBUG") != "" {
			return nil, fmt.Errorf(fmt.Sprintf("(%s) [%s] ", edgeGwName, callFuncName())+errorUnableToFindEdgeGateway, err)
		}
		return nil, fmt.Errorf(errorUnableToFindEdgeGateway, err)
	}
	return eg, nil
}

// Same as GetEdgeGateway, but using data from the resource, if available
// edgeGatewayFieldName is the name used in the resource. It is usually "edge_gateway"
// for all resources that *use* an edge gateway, and when the resource is vcd_edgegateway, it is "name"
func (cli *VCDClient) GetEdgeGatewayFromResource(d *schema.ResourceData, edgeGatewayFieldName string) (eg *govcd.EdgeGateway, err error) {
	orgName := d.Get("org").(string)
	vdcName := d.Get("vdc").(string)
	edgeGatewayName := d.Get(edgeGatewayFieldName).(string)
	egw, err := cli.GetEdgeGateway(orgName, vdcName, edgeGatewayName)
	if err != nil {
		if os.Getenv("GOVCD_DEBUG") != "" {
			return nil, fmt.Errorf("(%s) [%s] : %s", edgeGatewayName, callFuncName(), err)
		}
		return nil, err
	}
	return egw, nil
}

func ProviderAuthenticate(client *govcd.VCDClient, user, password, token, org string) error {
	var err error
	if token != "" {
		err = client.SetToken(org, govcd.AuthorizationHeader, token)
		if err != nil {
			err = fmt.Errorf("error during token-based authentication: %s", err)
		}
	} else {
		err = client.Authenticate(user, password, org)
	}
	return err
}

func ProviderAuthenticateSaml(client *govcd.VCDClient, domainUser, domainPassword, org string) error {

	url := client.Client.VCDHREF

	backupRedirectChecker := client.Client.Http.CheckRedirect
	defer func() {
		client.Client.Http.CheckRedirect = backupRedirectChecker
	}()

	client.Client.Http.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	loginUrl := url.Scheme + "://" + url.Host + "/tenant/" + org

	resp, err := client.Client.Http.Get(loginUrl)
	if err != nil {
		return err
	}

	// https://192.168.1.160/login?service=tenant:my-org&redirectTo=%2Ftenant%2Fmy-org%2F
	nextHop, _ := resp.Location()

	resp, err = client.Client.Http.Get(nextHop.String())
	if err != nil {
		return err
	}

	// https://192.168.1.160/login/?service=tenant:my-org
	nextHop, _ = resp.Location()

	resp, err = client.Client.Http.Get(nextHop.String())
	if err != nil {
		return err
	}

	// https://192.168.1.160/login/my-org/saml/login/alias/vcd?service=tenant:my-org
	nextHop, _ = resp.Location()
	// https://win-60g606n0afg.test-forest.net/adfs/ls/?SAMLRequest=fZFNT8MwDIbv%2FIoq9zYflDKitdMEQkICCcHgwM1K3S5Tm4w4HfDvyQYTcOESK9Lrx86T%2BeJ9HLIdBrLe1UwWgmXojG%2Bt62v2tLrOZ2zRnMwJxkFt9XKKa%2FeArxNSzJZEGGLqu%2FSOphHDI4adNXjjWnyvWSJdpZh1EA%2FsdYxb0py%2FWZdXoq9E5QR0fRFTKO98SKVwGDm0HfGBOMuufTB4GFmzDgZClt1c1QxKs5mhrJSFslRdrzYbgWDkWQVnQp6mEN0Dkd3hTxvRlPaiCC7WTAklclHmQq7EuVZKy4uinJ2%2BsOz5KELtRSQ1jvTX02s2Bac9kCXtYETS0ejH5d2tTlG9DT564wfWfJnSh4HhN%2BF%2FABxdsuboSV6oQlazQqZTcDP4qeU%2B9Hz8yPdlD%2BUjRmghAofBAvGdaef89wLN9%2FXvzzWf&RelayState=aHR0cHM6Ly8xOTIuMTY4LjEuMTYwL3RlbmFudC9teS1vcmc%3D&SigAlg=http%3A%2F%2Fwww.w3.org%2F2000%2F09%2Fxmldsig%23rsa-sha1&Signature=W3mUfGiecEJudJLqIV%2F2cFuJcbPiQxgayxVbJf6hOMp8ZQcqG01NR1Rm3qTqaote5dSkprw42dVOMIHdeiJL1g7%2FW9ON6%2BvzJHvL3rdy652%2BeSv6q0r9wDJ8eKC5DpwcmW0UUATHHt4ENMPa6w6MgE2Mwm1F6eYu1c5CcIC306lzQNiSWwNA08frX1wxl3RtSrrm9qo9K9UoQOAULkYjAgghI65Dr%2BEWjiu%2FYgVw1SuXMKRQcQ1Q8MQ2uhDqjlfXuO3Fnp582zLh1uMx1ZiFO1LPaqTT7K%2BBvlUISzkCe6YkSOlr%2Fz08t7A4fkMxmhnTd2gBu6WhzcUUWkqLMzSwgw%3D%3D

	resp, err = client.Client.Http.Get(nextHop.String())
	if err != nil {
		return err
	}

	nextHop, _ = resp.Location()

	log.Println("DAINIUS nexthop ", nextHop)

	authEndPoint := nextHop.Scheme + "://" + nextHop.Hostname() + "/adfs/services/trust/13/usernamemixed"

	body := strings.NewReader(requestBody)

	log.Println("DAINIUS posting to endpoint: ", authEndPoint)

	resp, err = client.Client.Http.Post(authEndPoint, "application/soap+xml", body)
	if err != nil {
		return err
	}

	log.Println("DAINIUS auth body: ", requestBody)

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

	// e, err := xml.Marshal(responseStruct.Body.RequestSecurityTokenResponseCollection.RequestSecurityTokenResponse.RequestedSecurityToken.EncryptedAss)
	// if err != nil {
	// 	return fmt.Errorf("unable marshal : %s", err)
	// }

	log.Printf("DAINIUS encrypted assertion: %s", tokenPart)

	// g, err := gzip.

	var b bytes.Buffer
	gz := gzip.NewWriter(&b)
	if _, err := gz.Write([]byte(tokenPart)); err != nil {
		log.Fatal(err)
	}
	if err := gz.Close(); err != nil {
		log.Fatal(err)
	}
	// fmt.Println(b.Bytes())
	// b.S
	// bbb, err := base64.StdEncoding.E
	sEnc := base64.StdEncoding.EncodeToString(b.Bytes())
	// sEnc := base64.StdEncoding.DecodeString(b)

	log.Printf("DAINIUS encoded test %s", sEnc)

	/// # Got data, try to authenticate against vCD

	req, err := http.NewRequest(http.MethodPost, "https://192.168.1.160/api/sessions", nil)
	if err != nil {
		return fmt.Errorf("error posting: %s", err)
	}
	req.Header.Add("Accept", "application/*+xml;version=29.0")
	req.Header.Add("Authorization", `SIGN token="`+sEnc+`",org="`+org+`"`)

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

var requestBody = `<s:Envelope xmlns:s="http://www.w3.org/2003/05/soap-envelope" xmlns:a="http://www.w3.org/2005/08/addressing" xmlns:u="http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-wssecurity-utility-1.0.xsd">
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
        <o:Username>test@test-forest.net</o:Username>
        <o:Password o:Type="http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-username-token-profile-1.0#PasswordText">R00tas!23</o:Password>
      </o:UsernameToken>
    </o:Security>
  </s:Header>
  <s:Body>
    <trust:RequestSecurityToken xmlns:trust="http://docs.oasis-open.org/ws-sx/ws-trust/200512">
      <wsp:AppliesTo xmlns:wsp="http://schemas.xmlsoap.org/ws/2004/09/policy">
        <a:EndpointReference>
          <a:Address>https://192.168.1.160/cloud/org/my-org/saml/metadata/alias/vcd</a:Address>
        </a:EndpointReference>
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
        <a:EndpointReference>
          <a:Address>{{.EntityId}}</a:Address>
        </a:EndpointReference>
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

func (c *Config) Client() (*VCDClient, error) {
	rawData := c.User + "#" +
		c.Password + "#" +
		c.Token + "#" +
		c.SysOrg + "#" +
		c.Href
	checksum := fmt.Sprintf("%x", sha1.Sum([]byte(rawData)))

	// The cached connection is served only if the variable VCD_CACHE is set
	cachedVCDClients.Lock()
	client, ok := cachedVCDClients.conMap[checksum]
	cachedVCDClients.Unlock()
	if ok && enableConnectionCache {
		cachedVCDClients.Lock()
		cachedVCDClients.cacheClientServedCount += 1
		cachedVCDClients.Unlock()
		// debugPrintf("[%s] cached connection served %d times (size:%d)\n",
		elapsed := time.Since(client.initTime)
		if elapsed > maxConnectionValidity {
			debugPrintf("cached connection invalidated after %2.0f minutes \n", maxConnectionValidity.Minutes())
			cachedVCDClients.Lock()
			delete(cachedVCDClients.conMap, checksum)
			cachedVCDClients.Unlock()
		} else {
			return client.connection, nil
		}
	}

	authUrl, err := url.ParseRequestURI(c.Href)
	if err != nil {
		return nil, fmt.Errorf("something went wrong while retrieving URL: %s", err)
	}

	vcdClient := &VCDClient{
		VCDClient: govcd.NewVCDClient(*authUrl, c.InsecureFlag,
			govcd.WithMaxRetryTimeout(c.MaxRetryTimeout)),
		SysOrg:          c.SysOrg,
		Org:             c.Org,
		Vdc:             c.Vdc,
		MaxRetryTimeout: c.MaxRetryTimeout,
		InsecureFlag:    c.InsecureFlag}

	// err = ProviderAuthenticate(vcdClient.VCDClient, c.User, c.Password, c.Token, c.SysOrg)
	err = ProviderAuthenticateSaml(vcdClient.VCDClient, c.User, c.Password, c.SysOrg)
	if err != nil {
		return nil, fmt.Errorf("something went wrong during authentication: %s", err)
	}
	cachedVCDClients.Lock()
	cachedVCDClients.conMap[checksum] = cachedConnection{initTime: time.Now(), connection: vcdClient}
	cachedVCDClients.Unlock()

	return vcdClient, nil
}

// Returns the name of the function that called the
// current function.
// It is used for tracing
func callFuncName() string {
	fpcs := make([]uintptr, 1)
	n := runtime.Callers(3, fpcs)
	if n > 0 {
		fun := runtime.FuncForPC(fpcs[0] - 1)
		if fun != nil {
			return fun.Name()
		}
	}
	return ""
}

func init() {
	separator := os.Getenv("VCD_IMPORT_SEPARATOR")
	if separator != "" {
		ImportSeparator = separator
	}
}
