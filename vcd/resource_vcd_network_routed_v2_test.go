// +build network ALL functional

package vcd

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccVcdNetworkRoutedV2Nsxv(t *testing.T) {
	// String map to fill the template
	var params = StringMap{
		"Org":    testConfig.VCD.Org,
		"Vdc":    testConfig.VCD.Vdc,
		"EdgeGw": testConfig.Networking.EdgeGateway,
		// "Tags": "lb lbVirtualServer",
	}

	configText := templateFill(TestAccVcdNetworkRoutedV2NsxvStep0, params)
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

const TestAccVcdNetworkRoutedV2NsxvStep0 = `
data "vcd_edgegateway" "existing" {
  org  = "{{.Org}}"
  vdc  = "{{.Vdc}}"
  name = "{{.EdgeGw}}"
}

resource "vcd_network_routed_v2" "net1" {
  org  = "{{.Org}}"
  vdc  = "{{.Vdc}}"
  name = "TestRoutedNsxv"
  description = "NSX-V routed network test OpenAPI"

  edge_gateway_id = data.vcd_edgegateway.existing.id
  
  gateway = "1.1.1.1"
  prefix_length = 24


  static_ip_pool {
	start_address = "1.1.1.10"
    end_address = "1.1.1.20"
  }
  
}
`

func TestAccVcdNetworkRoutedV2Nsxt(t *testing.T) {
	// String map to fill the template
	var params = StringMap{
		"Org":    testConfig.VCD.Org,
		"Vdc":    testConfig.Nsxt.Vdc,
		"EdgeGw": "nsxt-edge",
		// "Tags": "lb lbVirtualServer",
	}

	configText := templateFill(TestAccVcdNetworkRoutedV2NsxtStep0, params)
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

const TestAccVcdNetworkRoutedV2NsxtStep0 = `
data "vcd_nsxt_edgegateway" "existing" {
  org  = "{{.Org}}"
  vdc  = "{{.Vdc}}"
  name = "{{.EdgeGw}}"
}

resource "vcd_network_routed_v2" "net1" {
  org  = "{{.Org}}"
  vdc  = "{{.Vdc}}"
  name = "TestRoutedNsxt"
  description = "NSX-T routed network test OpenAPI"

  edge_gateway_id = data.vcd_nsxt_edgegateway.existing.id

  gateway = "1.1.1.1"
  prefix_length = 24

  static_ip_pool {
	start_address = "1.1.1.10"
    end_address = "1.1.1.20"
  }
}
`
