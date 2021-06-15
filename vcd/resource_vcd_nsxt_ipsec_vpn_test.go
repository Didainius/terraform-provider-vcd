// +build network nsxt ALL functional

package vcd

import (
	"fmt"
	"testing"
	"time"

	"github.com/davecgh/go-spew/spew"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccVcdNsxtIpSecVpnTunnel(t *testing.T) {
	preTestChecks(t)
	skipNoNsxtConfiguration(t)

	// String map to fill the template
	var params = StringMap{
		"Org":         testConfig.VCD.Org,
		"NsxtVdc":     testConfig.Nsxt.Vdc,
		"EdgeGw":      testConfig.Nsxt.EdgeGateway,
		"NetworkName": t.Name(),
		"Tags":        "network nsxt",
	}

	configText := templateFill(testAccNsxtIpSecVpnTunnel, params)
	debugPrintf("#[DEBUG] CONFIGURATION for step 1: %s", configText)

	//params["FuncName"] = t.Name() + "-step1"
	//configText1 := templateFill(testAccNsxtIpSetEmpty2, params)
	//debugPrintf("#[DEBUG] CONFIGURATION for step 2: %s", configText1)
	//
	//params["FuncName"] = t.Name() + "-step11"
	//configText11 := templateFill(testAccNsxtIpSetEmpty2+testAccNsxtIpSetDS, params)
	//debugPrintf("#[DEBUG] CONFIGURATION for step 3: %s", configText11)
	//
	//params["FuncName"] = t.Name() + "-step2"
	//configText2 := templateFill(testAccNsxtIpSetIpRanges, params)
	//debugPrintf("#[DEBUG] CONFIGURATION for step 4: %s", configText2)

	if vcdShortTest {
		t.Skip(acceptanceTestsSkipped)
		return
	}

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviders,
		PreCheck:          func() { testAccPreCheck(t) },
		//CheckDestroy: resource.ComposeAggregateTestCheckFunc(
		//	testAccCheckNsxtFirewallGroupDestroy(testConfig.Nsxt.Vdc, "test-ip-set", types.FirewallGroupTypeIpSet),
		//	testAccCheckNsxtFirewallGroupDestroy(testConfig.Nsxt.Vdc, "test-ip-set-changed", types.FirewallGroupTypeIpSet),
		//),
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: configText,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("vcd_nsxt_ipsec_vpn_tunnel.tunnel1", "id"),
					//resource.TestMatchResourceAttr("vcd_nsxt_ip_set.set1", "id", regexp.MustCompile(`^urn:vcloud:firewallGroup:.*$`)),
					//resource.TestCheckResourceAttr("vcd_nsxt_ip_set.set1", "name", "test-ip-set"),
					//resource.TestCheckResourceAttr("vcd_nsxt_ip_set.set1", "description", "test-ip-set-description"),
					//resource.TestCheckResourceAttr("vcd_nsxt_ip_set.set1", "ip_addresses.#", "0"),
				),
			},
			resource.TestStep{
				Config: configText,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("vcd_nsxt_ipsec_vpn_tunnel.tunnel1", "id"),
					resource.TestCheckResourceAttrSet("vcd_nsxt_ipsec_vpn_tunnel.tunnel1", "status"),
					//resource.TestMatchResourceAttr("vcd_nsxt_ip_set.set1", "id", regexp.MustCompile(`^urn:vcloud:firewallGroup:.*$`)),
					//resource.TestCheckResourceAttr("vcd_nsxt_ip_set.set1", "name", "test-ip-set"),
					//resource.TestCheckResourceAttr("vcd_nsxt_ip_set.set1", "description", "test-ip-set-description"),
					//resource.TestCheckResourceAttr("vcd_nsxt_ip_set.set1", "ip_addresses.#", "0"),
				),
			},
			//resource.TestStep{
			//	Config: configText1,
			//	Check: resource.ComposeAggregateTestCheckFunc(
			//		resource.TestMatchResourceAttr("vcd_nsxt_ip_set.set1", "id", regexp.MustCompile(`^urn:vcloud:firewallGroup:.*$`)),
			//		resource.TestCheckResourceAttr("vcd_nsxt_ip_set.set1", "name", "test-ip-set-changed"),
			//		resource.TestCheckResourceAttr("vcd_nsxt_ip_set.set1", "description", ""),
			//		resource.TestCheckResourceAttr("vcd_nsxt_ip_set.set1", "ip_addresses.#", "0"),
			//	),
			//},
			//resource.TestStep{
			//	Config: configText11,
			//	Check: resource.ComposeAggregateTestCheckFunc(
			//		resource.TestMatchResourceAttr("vcd_nsxt_ip_set.set1", "id", regexp.MustCompile(`^urn:vcloud:firewallGroup:.*$`)),
			//		resource.TestCheckResourceAttr("vcd_nsxt_ip_set.set1", "name", "test-ip-set-changed"),
			//		resource.TestCheckResourceAttr("vcd_nsxt_ip_set.set1", "description", ""),
			//		resource.TestCheckResourceAttr("vcd_nsxt_ip_set.set1", "ip_addresses.#", "0"),
			//
			//		resourceFieldsEqual("vcd_nsxt_ip_set.set1", "data.vcd_nsxt_ip_set.ds", []string{}),
			//	),
			//},
			//// Test import with no IP addresses
			//resource.TestStep{
			//	ResourceName:      "vcd_nsxt_ip_set.set1",
			//	ImportState:       true,
			//	ImportStateVerify: true,
			//	ImportStateIdFunc: importStateIdNsxtEdgeGatewayObject(testConfig, testConfig.Nsxt.EdgeGateway, "test-ip-set-changed"),
			//},
			//resource.TestStep{
			//	Config: configText2,
			//	Check: resource.ComposeAggregateTestCheckFunc(
			//		resource.TestMatchResourceAttr("vcd_nsxt_ip_set.set1", "id", regexp.MustCompile(`^urn:vcloud:firewallGroup:.*$`)),
			//		resource.TestCheckResourceAttr("vcd_nsxt_ip_set.set1", "name", "test-ip-set-changed"),
			//		resource.TestCheckResourceAttr("vcd_nsxt_ip_set.set1", "description", ""),
			//		resource.TestCheckTypeSetElemAttr("vcd_nsxt_ip_set.set1", "ip_addresses.*", "12.12.12.1"),
			//		resource.TestCheckTypeSetElemAttr("vcd_nsxt_ip_set.set1", "ip_addresses.*", "10.10.10.0/24"),
			//		resource.TestCheckTypeSetElemAttr("vcd_nsxt_ip_set.set1", "ip_addresses.*", "11.11.11.1-11.11.11.2"),
			//		resource.TestCheckTypeSetElemAttr("vcd_nsxt_ip_set.set1", "ip_addresses.*", "2001:db8::/48"),
			//		resource.TestCheckTypeSetElemAttr("vcd_nsxt_ip_set.set1", "ip_addresses.*", "2001:db6:0:0:0:0:0:0-2001:db6:0:ffff:ffff:ffff:ffff:ffff"),
			//		resource.TestCheckResourceAttr("vcd_nsxt_ip_set.set1", "ip_addresses.#", "5"),
			//	),
			//},
			//// Test import with IP addresses
			//resource.TestStep{
			//	ResourceName:      "vcd_nsxt_ip_set.set1",
			//	ImportState:       true,
			//	ImportStateVerify: true,
			//	ImportStateIdFunc: importStateIdNsxtEdgeGatewayObject(testConfig, testConfig.Nsxt.EdgeGateway, "test-ip-set-changed"),
			//},
		},
	})
	postTestChecks(t)
}

func stateDumper() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		spew.Dump(s)
		return nil
	}
}

func sleepTester() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		fmt.Println("sleeping")
		time.Sleep(5 * time.Minute)
		return nil
	}
}

const testAccNsxtIpSecVpnTunnel = testAccNsxtIpSetPrereqs + `
resource "vcd_nsxt_ipsec_vpn_tunnel" "tunnel1" {
  org = "{{.Org}}"
  vdc = "{{.NsxtVdc}}"

  edge_gateway_id = data.vcd_nsxt_edgegateway.existing_gw.id

  name        = "test-tunnel-1"
  description = "test-tunnel-description"
  
  pre_shared_key    = "test-psk"
  # Primary IP address of Edge Gateway
  local_ip_address  = tolist(data.vcd_nsxt_edgegateway.existing_gw.subnet)[0].primary_ip
  local_networks    = ["10.10.10.0/24", "30.30.30.0/28", "40.40.40.1/32"]
  # That is a fake remote IP address
  remote_ip_address = "1.2.3.4"
  remote_networks   = ["192.168.1.0/24", "192.168.10.0/24", "192.168.20.0/28"]
}
`

func TestAccVcdNsxtIpSecVpnTunnelCustomProfile(t *testing.T) {
	preTestChecks(t)
	skipNoNsxtConfiguration(t)

	// String map to fill the template
	var params = StringMap{
		"Org":         testConfig.VCD.Org,
		"NsxtVdc":     testConfig.Nsxt.Vdc,
		"EdgeGw":      testConfig.Nsxt.EdgeGateway,
		"NetworkName": t.Name(),
		"Tags":        "network nsxt",
	}

	configText1 := templateFill(testAccNsxtIpSecVpnTunnelProfileStep1, params)
	debugPrintf("#[DEBUG] CONFIGURATION for step 1: %s", configText1)

	params["FuncName"] = t.Name() + "-step2"
	configText2 := templateFill(testAccNsxtIpSecVpnTunnelProfileStep2, params)
	debugPrintf("#[DEBUG] CONFIGURATION for step 2: %s", configText2)
	//
	//params["FuncName"] = t.Name() + "-step11"
	//configText11 := templateFill(testAccNsxtIpSetEmpty2+testAccNsxtIpSetDS, params)
	//debugPrintf("#[DEBUG] CONFIGURATION for step 3: %s", configText11)
	//
	//params["FuncName"] = t.Name() + "-step2"
	//configText2 := templateFill(testAccNsxtIpSetIpRanges, params)
	//debugPrintf("#[DEBUG] CONFIGURATION for step 4: %s", configText2)

	if vcdShortTest {
		t.Skip(acceptanceTestsSkipped)
		return
	}

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviders,
		PreCheck:          func() { testAccPreCheck(t) },
		//CheckDestroy: resource.ComposeAggregateTestCheckFunc(
		//	testAccCheckNsxtFirewallGroupDestroy(testConfig.Nsxt.Vdc, "test-ip-set", types.FirewallGroupTypeIpSet),
		//	testAccCheckNsxtFirewallGroupDestroy(testConfig.Nsxt.Vdc, "test-ip-set-changed", types.FirewallGroupTypeIpSet),
		//),
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: configText1,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("vcd_nsxt_ipsec_vpn_tunnel.tunnel1", "id"),
					resource.TestCheckResourceAttr("vcd_nsxt_ipsec_vpn_tunnel.tunnel1", "security_profile_customization.#", "1"),
					//stateDumper(),
					//sleepTester(),
					//resource.TestMatchResourceAttr("vcd_nsxt_ip_set.set1", "id", regexp.MustCompile(`^urn:vcloud:firewallGroup:.*$`)),
					//resource.TestCheckResourceAttr("vcd_nsxt_ip_set.set1", "name", "test-ip-set"),
					//resource.TestCheckResourceAttr("vcd_nsxt_ip_set.set1", "description", "test-ip-set-description"),
					//resource.TestCheckResourceAttr("vcd_nsxt_ip_set.set1", "ip_addresses.#", "0"),
				),
			},
			resource.TestStep{
				Config: configText2,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("vcd_nsxt_ipsec_vpn_tunnel.tunnel1", "id"),
					resource.TestCheckResourceAttr("vcd_nsxt_ipsec_vpn_tunnel.tunnel1", "security_profile_customization.#", "0"),

					//stateDumper(),
					//sleepTester(),
					//resource.TestMatchResourceAttr("vcd_nsxt_ip_set.set1", "id", regexp.MustCompile(`^urn:vcloud:firewallGroup:.*$`)),
					//resource.TestCheckResourceAttr("vcd_nsxt_ip_set.set1", "name", "test-ip-set"),
					//resource.TestCheckResourceAttr("vcd_nsxt_ip_set.set1", "description", "test-ip-set-description"),
					//resource.TestCheckResourceAttr("vcd_nsxt_ip_set.set1", "ip_addresses.#", "0"),
				),
			},
			//resource.TestStep{
			//	Config: configText1,
			//	Check: resource.ComposeAggregateTestCheckFunc(
			//		resource.TestCheckResourceAttrSet("vcd_nsxt_ipsec_vpn_tunnel.tunnel1", "id"),
			//		resource.TestCheckResourceAttrSet("vcd_nsxt_ipsec_vpn_tunnel.tunnel1", "status"),
			//		//resource.TestMatchResourceAttr("vcd_nsxt_ip_set.set1", "id", regexp.MustCompile(`^urn:vcloud:firewallGroup:.*$`)),
			//		//resource.TestCheckResourceAttr("vcd_nsxt_ip_set.set1", "name", "test-ip-set"),
			//		//resource.TestCheckResourceAttr("vcd_nsxt_ip_set.set1", "description", "test-ip-set-description"),
			//		//resource.TestCheckResourceAttr("vcd_nsxt_ip_set.set1", "ip_addresses.#", "0"),
			//	),
			//},
			//resource.TestStep{
			//	Config: configText2,
			//	Check: resource.ComposeAggregateTestCheckFunc(
			//		resource.TestMatchResourceAttr("vcd_nsxt_ip_set.set1", "id", regexp.MustCompile(`^urn:vcloud:firewallGroup:.*$`)),
			//		resource.TestCheckResourceAttr("vcd_nsxt_ip_set.set1", "name", "test-ip-set-changed"),
			//		resource.TestCheckResourceAttr("vcd_nsxt_ip_set.set1", "description", ""),
			//		resource.TestCheckResourceAttr("vcd_nsxt_ip_set.set1", "ip_addresses.#", "0"),
			//	),
			//},
			//resource.TestStep{
			//	Config: configText11,
			//	Check: resource.ComposeAggregateTestCheckFunc(
			//		resource.TestMatchResourceAttr("vcd_nsxt_ip_set.set1", "id", regexp.MustCompile(`^urn:vcloud:firewallGroup:.*$`)),
			//		resource.TestCheckResourceAttr("vcd_nsxt_ip_set.set1", "name", "test-ip-set-changed"),
			//		resource.TestCheckResourceAttr("vcd_nsxt_ip_set.set1", "description", ""),
			//		resource.TestCheckResourceAttr("vcd_nsxt_ip_set.set1", "ip_addresses.#", "0"),
			//
			//		resourceFieldsEqual("vcd_nsxt_ip_set.set1", "data.vcd_nsxt_ip_set.ds", []string{}),
			//	),
			//},
			//// Test import with no IP addresses
			//resource.TestStep{
			//	ResourceName:      "vcd_nsxt_ip_set.set1",
			//	ImportState:       true,
			//	ImportStateVerify: true,
			//	ImportStateIdFunc: importStateIdNsxtEdgeGatewayObject(testConfig, testConfig.Nsxt.EdgeGateway, "test-ip-set-changed"),
			//},
			//resource.TestStep{
			//	Config: configText2,
			//	Check: resource.ComposeAggregateTestCheckFunc(
			//		resource.TestMatchResourceAttr("vcd_nsxt_ip_set.set1", "id", regexp.MustCompile(`^urn:vcloud:firewallGroup:.*$`)),
			//		resource.TestCheckResourceAttr("vcd_nsxt_ip_set.set1", "name", "test-ip-set-changed"),
			//		resource.TestCheckResourceAttr("vcd_nsxt_ip_set.set1", "description", ""),
			//		resource.TestCheckTypeSetElemAttr("vcd_nsxt_ip_set.set1", "ip_addresses.*", "12.12.12.1"),
			//		resource.TestCheckTypeSetElemAttr("vcd_nsxt_ip_set.set1", "ip_addresses.*", "10.10.10.0/24"),
			//		resource.TestCheckTypeSetElemAttr("vcd_nsxt_ip_set.set1", "ip_addresses.*", "11.11.11.1-11.11.11.2"),
			//		resource.TestCheckTypeSetElemAttr("vcd_nsxt_ip_set.set1", "ip_addresses.*", "2001:db8::/48"),
			//		resource.TestCheckTypeSetElemAttr("vcd_nsxt_ip_set.set1", "ip_addresses.*", "2001:db6:0:0:0:0:0:0-2001:db6:0:ffff:ffff:ffff:ffff:ffff"),
			//		resource.TestCheckResourceAttr("vcd_nsxt_ip_set.set1", "ip_addresses.#", "5"),
			//	),
			//},
			//// Test import with IP addresses
			//resource.TestStep{
			//	ResourceName:      "vcd_nsxt_ip_set.set1",
			//	ImportState:       true,
			//	ImportStateVerify: true,
			//	ImportStateIdFunc: importStateIdNsxtEdgeGatewayObject(testConfig, testConfig.Nsxt.EdgeGateway, "test-ip-set-changed"),
			//},
		},
	})
	postTestChecks(t)
}

const testAccNsxtIpSecVpnTunnelProfileStep1 = testAccNsxtIpSetPrereqs + `
resource "vcd_nsxt_ipsec_vpn_tunnel" "tunnel1" {
  org = "{{.Org}}"
  vdc = "{{.NsxtVdc}}"

  edge_gateway_id = data.vcd_nsxt_edgegateway.existing_gw.id

  name        = "test-tunnel-1"
  description = "test-tunnel-description"
  
  pre_shared_key    = "test-psk"
  # Primary IP address of Edge Gateway
  local_ip_address  = tolist(data.vcd_nsxt_edgegateway.existing_gw.subnet)[0].primary_ip
  local_networks    = ["10.10.10.0/24", "30.30.30.0/28", "40.40.40.1/32"]
  # That is a fake remote IP address as there is nothing else to peer to
  remote_ip_address = "1.2.3.4"
  remote_networks   = ["192.168.1.0/24", "192.168.10.0/24", "192.168.20.0/28"]

  security_profile_customization {
    ike_version               = "IKE_V2"
    ike_encryption_algorithms = ["AES_128"]
    ike_digest_algorithms     = ["SHA2_256"]
    ike_dh_groups             = ["GROUP14"]
    ike_sa_lifetime           = 86400
    
	tunnel_pfs_enabled = true
	tunnel_df_policy = "COPY"
    tunnel_encryption_algorithms = ["AES_GCM_128"]
    # tunnel_digest_algorithms     = ["AES_GCM_128"]
    tunnel_dh_groups             = ["GROUP14"]
    tunnel_sa_lifetime           = 3600
    
    dpd_probe_internal = "30"
  }
}
`

const testAccNsxtIpSecVpnTunnelProfileStep2 = testAccNsxtIpSetPrereqs + `
resource "vcd_nsxt_ipsec_vpn_tunnel" "tunnel1" {
  org = "{{.Org}}"
  vdc = "{{.NsxtVdc}}"

  edge_gateway_id = data.vcd_nsxt_edgegateway.existing_gw.id

  name        = "test-tunnel-1"
  description = "test-tunnel-description"
  
  pre_shared_key    = "test-psk"
  # Primary IP address of Edge Gateway
  local_ip_address  = tolist(data.vcd_nsxt_edgegateway.existing_gw.subnet)[0].primary_ip
  local_networks    = ["10.10.10.0/24", "30.30.30.0/28", "40.40.40.1/32"]
  # That is a fake remote IP address as there is nothing else to peer to
  remote_ip_address = "1.2.3.4"
  remote_networks   = ["192.168.1.0/24", "192.168.10.0/24", "192.168.20.0/28"]
}
`
