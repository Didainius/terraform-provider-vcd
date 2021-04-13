// +build network nsxt ALL functional

package vcd

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccVcdNsxtFirewallGroup(t *testing.T) {
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

	configText := templateFill(testAccNsxtFirewallGroup, params)
	debugPrintf("#[DEBUG] CONFIGURATION for step 1: %s", configText)

	params["FuncName"] = t.Name() + "-step1"
	configText1 := templateFill(testAccNsxtFirewallGroup2, params)
	debugPrintf("#[DEBUG] CONFIGURATION for step 2: %s", configText1)

	if vcdShortTest {
		t.Skip(acceptanceTestsSkipped)
		return
	}

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviders,
		PreCheck:          func() { testAccPreCheck(t) },
		// CheckDestroy:      testAccCheckOpenApiVcdNetworkDestroy(testConfig.Nsxt.Vdc, "nsxt-routed-dhcp"),
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: configText,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("vcd_nsxt_security_group.group1", "id", regexp.MustCompile(`^urn:vcloud:firewallGroup:.*$`)),
					resource.TestCheckResourceAttr("vcd_nsxt_security_group.group1", "name", "test-security-group"),
					resource.TestCheckResourceAttr("vcd_nsxt_security_group.group1", "description", "test-security-group-description"),
					resource.TestCheckNoResourceAttr("vcd_nsxt_security_group.group1", "member_org_network_ids"),
					resource.TestCheckNoResourceAttr("vcd_nsxt_security_group.group1", "member_vm_ids"),
				),
			},
			resource.TestStep{
				Config: configText1,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("vcd_nsxt_security_group.group1", "id", regexp.MustCompile(`^urn:vcloud:firewallGroup:.*$`)),
					resource.TestCheckResourceAttr("vcd_nsxt_security_group.group1", "name", "test-security-group-changed"),
					resource.TestCheckResourceAttr("vcd_nsxt_security_group.group1", "description", ""),
					resource.TestCheckNoResourceAttr("vcd_nsxt_security_group.group1", "member_org_network_ids"),
					resource.TestCheckNoResourceAttr("vcd_nsxt_security_group.group1", "member_vm_ids"),
				),
			},
			// resource.TestStep{
			// 	Config: configText1,
			// 	Check: resource.ComposeAggregateTestCheckFunc(
			// 		resource.TestMatchResourceAttr("vcd_nsxt_network_dhcp.pools", "id", regexp.MustCompile(`^urn:vcloud:network:.*$`)),
			// 		resource.TestCheckTypeSetElemNestedAttrs("vcd_nsxt_network_dhcp.pools", "pool.*", map[string]string{
			// 			"start_address": "7.1.1.100",
			// 			"end_address":   "7.1.1.110",
			// 		}),
			// 		resource.TestCheckTypeSetElemNestedAttrs("vcd_nsxt_network_dhcp.pools", "pool.*", map[string]string{
			// 			"start_address": "7.1.1.130",
			// 			"end_address":   "7.1.1.140",
			// 		}),
			// 	),
			// },
			resource.TestStep{
				ResourceName:      "vcd_nsxt_security_group.group1",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: importStateIdNsxtEdgeGatewayObject(testConfig, testConfig.Nsxt.EdgeGateway, "test-security-group-changed"),
			},
		},
	})
	postTestChecks(t)
}

const testAccNsxtFirewallGroupPrereqs = `
data "vcd_nsxt_edgegateway" "existing" {
	org  = "{{.Org}}"
	vdc  = "{{.NsxtVdc}}"

	name = "{{.EdgeGw}}"
}
`

const testAccNsxtFirewallGroup = testAccNsxtFirewallGroupPrereqs + `
resource "vcd_nsxt_security_group" "group1" {
  org  = "{{.Org}}"
  vdc  = "{{.NsxtVdc}}"

  edge_gateway_id = data.vcd_nsxt_edgegateway.existing.id

  name = "test-security-group"
  description = "test-security-group-description"
}
`

const testAccNsxtFirewallGroup2 = testAccNsxtFirewallGroupPrereqs + `
resource "vcd_nsxt_security_group" "group1" {
  org  = "{{.Org}}"
  vdc  = "{{.NsxtVdc}}"

  edge_gateway_id = data.vcd_nsxt_edgegateway.existing.id

  name = "test-security-group-changed"
}
`
