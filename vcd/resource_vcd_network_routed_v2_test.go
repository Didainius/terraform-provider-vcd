// +build network ALL functional

package vcd

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// TestAccVcdNetworkRoutedV2NsxvInterfaceTypes attempts to test all supported interface types for
// NSX-V Org VDC routed network
func TestAccVcdNetworkRoutedV2NsxvInterfaceTypes(t *testing.T) {
	// String map to fill the template
	var params = StringMap{
		"Org":           testConfig.VCD.Org,
		"Vdc":           testConfig.VCD.Vdc,
		"EdgeGw":        testConfig.Networking.EdgeGateway,
		"InterfaceType": "internal",
		"NetworkName":   t.Name(),
		// "Tags": "lb lbVirtualServer",
	}

	// interface_type = "INTERNAL"
	// interface_type = "SUBINTERFACE"
	// interface_type = "TRUNK"
	// interface_type = "UPLINK"
	// interface_type = "DISTRIBUTED"
	// INTERNAL UPLINK TRUNK SUBINTERFACE

	configText := templateFill(testAccVcdNetworkRoutedV2Nsxv, params)
	debugPrintf("#[DEBUG] CONFIGURATION for step 0: %s", configText)

	params["FuncName"] = t.Name() + "-step1"
	params["InterfaceType"] = "subinterface"
	configText1 := templateFill(testAccVcdNetworkRoutedV2Nsxv, params)
	debugPrintf("#[DEBUG] CONFIGURATION for step 1: %s", configText1)

	// params["FuncName"] = t.Name() + "-step2"
	// params["InterfaceType"] = "distributed"
	// configText2 := templateFill(testAccVcdNetworkRoutedV2Nsxv, params)
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
			resource.TestStep{ // step 0
				Config: configText1,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("vcd_network_routed_v2.net1", "id"),
				),
			},
			// resource.TestStep{ // step 0
			// 	Config: configText2,
			// 	Check: resource.ComposeAggregateTestCheckFunc(
			// 		resource.TestCheckResourceAttrSet("vcd_network_routed_v2.net1", "id"),
			// 	),
			// },

			// Check that import works
			resource.TestStep{ // step 1
				ResourceName:      "vcd_network_routed_v2.net1",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: importStateIdOrgVdcObject(testConfig, t.Name()),
			},
		},
	})
}

const testAccVcdNetworkRoutedV2Nsxv = `
data "vcd_edgegateway" "existing" {
  org  = "{{.Org}}"
  vdc  = "{{.Vdc}}"
  name = "{{.EdgeGw}}"
}

resource "vcd_network_routed_v2" "net1" {
  org  = "{{.Org}}"
  vdc  = "{{.Vdc}}"
  name = "{{.NetworkName}}"
  description = "NSX-V routed network test OpenAPI"

  interface_type = "{{.InterfaceType}}"

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
	vcdClient := createTemporaryVCDConnection()
	if vcdClient.Client.APIVCDMaxVersionIs("< 34.0") {
		t.Skip(t.Name() + " requires at least API v34.0 (vCD 10.1.1+)")
	}
	skipNoNsxtConfiguration(t)

	// String map to fill the template
	var params = StringMap{
		"Org":         testConfig.VCD.Org,
		"Vdc":         testConfig.Nsxt.Vdc,
		"EdgeGw":      testConfig.Nsxt.EdgeGateway,
		"NetworkName": t.Name(),
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
			resource.TestStep{ // step 1
				ResourceName:      "vcd_network_routed_v2.net1",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: importStateIdOrgNsxtVdcObject(testConfig, t.Name()),
			},
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
  name = "{{.NetworkName}}"
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
