// +build network ALL functional

package vcd

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccVcdNetworkRoutedV2(t *testing.T) {
	// String map to fill the template
	var params = StringMap{
		"Org": testConfig.VCD.Org,
		"Vdc": testConfig.VCD.Vdc,
		// "Tags": "lb lbVirtualServer",
	}

	configText := templateFill(TestAccVcdNetworkRoutedV2Step0, params)
	debugPrintf("#[DEBUG] CONFIGURATION for step 0: %s", configText)

	// params["FuncName"] = t.Name() + "-step2"
	// params["VirtualServerName"] = t.Name() + "-step2"
	// configText2 := templateFill(testAccVcdLbVirtualServer_step2, params)
	// debugPrintf("#[DEBUG] CONFIGURATION for step 2: %s", configText2)

	if vcdShortTest {
		t.Skip(acceptanceTestsSkipped)
		return
	}

	resource.ParallelTest(t, resource.TestCase{
		ProviderFactories: testAccProviders,
		PreCheck:          func() { testAccPreCheck(t) },
		// CheckDestroy:      testAccCheckVcdLbVirtualServerDestroy(params["VirtualServerName"].(string)),
		Steps: []resource.TestStep{
			resource.TestStep{ // step 0
				Config: configText,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("vcd_network_routed_v2.net1", "id"),
				),
			},

			// Check that import works
			// resource.TestStep{ // step 1
			// 	ResourceName:      "vcd_lb_virtual_server.http",
			// 	ImportState:       true,
			// 	ImportStateVerify: true,
			// 	ImportStateIdFunc: importStateIdEdgeGatewayObject(testConfig, testConfig.Networking.EdgeGateway, t.Name()),
			// },
		},
	})
}

const TestAccVcdNetworkRoutedV2Step0 = `
data "vcd_nsxt_edgegateway" "existing" {
  org  = "dainius"
  vdc  = "nsxt-vdc-dainius"
  name = "nsxt-edge"
}

resource "vcd_network_routed_v2" "net1" {
  org  = "dainius"
  vdc  = "nsxt-vdc-dainius"
  name = "TestRoutedNsxt"
  description = "NSX-T routed network test OpenAPI"

  edge_gateway_id = data.vcd_nsxt_edgegateway.existing.id
  
  gateway = "1.1.1.1"
  prefix_length = 24
  
}
`
