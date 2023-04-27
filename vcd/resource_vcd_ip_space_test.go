//go:build network || nsxt || ALL || functional

package vcd

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccVcdIpSpacePublic(t *testing.T) {
	preTestChecks(t)
	skipIfNotSysAdmin(t)

	// String map to fill the template
	var params = StringMap{
		"TestName": t.Name(),

		"Tags": "network nsxt",
	}
	testParamsNotEmpty(t, params)

	params["FuncName"] = t.Name() + "step1"
	configText1 := templateFill(testAccVcdIpSpaceStep1, params)
	debugPrintf("#[DEBUG] CONFIGURATION for step 1: %s", configText1)

	// params["FuncName"] = t.Name() + "step2"
	// configText2 := templateFill(testAccNsxtEdgeGatewayInVdcDS, params)
	// debugPrintf("#[DEBUG] CONFIGURATION for step 2: %s", configText2)

	if vcdShortTest {
		t.Skip(acceptanceTestsSkipped)
		return
	}

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviders,
		// CheckDestroy:      testAccCheckVcdNsxtEdgeGatewayDestroy(params["NsxtEdgeGatewayVcd"].(string)),
		Steps: []resource.TestStep{
			{
				Config: configText1,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("vcd_ip_space.space1", "id"),
					// resource.TestCheckResourceAttr("vcd_nsxt_edgegateway.nsxt-edge", "name", params["NsxtEdgeGatewayVcd"].(string)),
					// resource.TestMatchResourceAttr("vcd_nsxt_edgegateway.nsxt-edge", "owner_id", regexp.MustCompile(`^urn:vcloud:vdc:`)),
					// resource.TestCheckResourceAttr("vcd_nsxt_edgegateway.nsxt-edge", "vdc", testConfig.Nsxt.Vdc),
					// resource.TestCheckResourceAttrPair("vcd_nsxt_edgegateway.nsxt-edge", "owner_id", "data.vcd_org_vdc.test", "id"),
				),
			},
			// {
			// 	Config: configText2,
			// 	Check: resource.ComposeTestCheckFunc(
			// 		resource.TestCheckResourceAttr("vcd_nsxt_edgegateway.nsxt-edge", "name", params["NsxtEdgeGatewayVcd"].(string)),
			// 		resource.TestMatchResourceAttr("vcd_nsxt_edgegateway.nsxt-edge", "owner_id", regexp.MustCompile(`^urn:vcloud:vdc:`)),
			// 		resource.TestCheckResourceAttr("vcd_nsxt_edgegateway.nsxt-edge", "vdc", testConfig.Nsxt.Vdc),
			// 		resource.TestCheckResourceAttrPair("vcd_nsxt_edgegateway.nsxt-edge", "owner_id", "data.vcd_org_vdc.test", "id"),
			// 		// Comparing data source and resource fields. Ignoring total field count '%' because data source does not have `starting_vdc_id`
			// 		resourceFieldsEqual("data.vcd_nsxt_edgegateway.nsxt-edge", "vcd_nsxt_edgegateway.nsxt-edge", []string{"%"}),
			// 	),
			// },
			// {
			// 	ResourceName:      "vcd_nsxt_edgegateway.nsxt-edge",
			// 	ImportState:       true,
			// 	ImportStateVerify: true,
			// 	ImportStateIdFunc: importStateIdOrgNsxtVdcObject(params["NsxtEdgeGatewayVcd"].(string)),
			// },
		},
	})
	postTestChecks(t)
}

const testAccVcdIpSpaceStep1 = `
resource "vcd_ip_space" "space1" {
  name = "{{.TestName}}"
  type = "PUBLIC"

  internal_scope = ["192.168.1.0/24"]

  route_advertisement_enabled = false
}
`
