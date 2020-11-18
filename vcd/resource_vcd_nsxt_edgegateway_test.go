// +build gateway ALL functional

package vcd

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccVcdNsxtEdgeGateway(t *testing.T) {
	// var (
	// 	edgeGatewayVcdName    string = "test_edge_gateway_basic"
	// 	newExternalNetwork    string = "TestExternalNetwork"
	// 	newExternalNetworkVcd string = "test_external_network"
	// )

	// String map to fill the template
	var params = StringMap{
		"Org":     testConfig.VCD.Org,
		"NsxtVdc": testConfig.Nsxt.Vdc,
		// "EdgeCluster": testConfig.Nsxt.,
		"EdgeGateway":        edgeGatewayNameBasic,
		"NsxtEdgeGatewayVcd": "nsxt-edge",
		"ExternalNetwork":    testConfig.Networking.ExternalNetwork,
		"Advanced":           "true",
		// "NewExternalNetwork":    newExternalNetwork,
		// "NewExternalNetworkVcd": newExternalNetworkVcd,
		"Type":      testConfig.Networking.ExternalNetworkPortGroupType,
		"PortGroup": testConfig.Networking.ExternalNetworkPortGroup,
		"Vcenter":   testConfig.Networking.Vcenter,
		"Tags":      "gateway",
	}
	configText := templateFill(testAccNsxtEdgeGateway, params)
	if vcdShortTest {
		t.Skip(acceptanceTestsSkipped)
		return
	}

	params["FuncName"] = t.Name() + "step1"
	configText1 := templateFill(testAccNsxtEdgeGatewayUpdate, params)
	if vcdShortTest {
		t.Skip(acceptanceTestsSkipped)
		return
	}

	if !usingSysAdmin() {
		t.Skip("Edge Gateway tests require system admin privileges")
		return
	}
	debugPrintf("#[DEBUG] CONFIGURATION: %s", configText)
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckVcdEdgeGatewayDestroy(edgeGatewayNameBasic),
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: configText,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("vcd_nsxt_edgegateway.nsxt-edge", "name", "nsxt-edge"),
					resource.TestCheckTypeSetElemNestedAttrs("vcd_nsxt_edgegateway.nsxt-edge", "subnet.*", map[string]string{
						"enabled":       "true",
						"gateway":       "10.150.191.253",
						"prefix_length": "19",
						"primary_ip":    "10.150.160.137",
					}),
					// resource.TestMatchResourceAttr("vcd_edgegateway.nsxt-edge", "name", ipV4Regex),
				),
			},
			resource.TestStep{
				Config: configText1,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("vcd_nsxt_edgegateway.nsxt-edge", "name", "nsxt-edge"),
					resource.TestCheckResourceAttr("vcd_nsxt_edgegateway.nsxt-edge", "description", "Updated-Description"),
					resource.TestCheckTypeSetElemNestedAttrs("vcd_nsxt_edgegateway.nsxt-edge", "subnet.*", map[string]string{
						"enabled":       "true",
						"gateway":       "10.150.191.253",
						"prefix_length": "19",
						"primary_ip":    "10.150.160.137",
					}),
					// resource.TestMatchResourceAttr("vcd_edgegateway.nsxt-edge", "name", ipV4Regex),
				),
			},
			// resource.TestStep{
			// 	ResourceName:            "vcd_edgegateway." + edgeGatewayNameBasic + "-import",
			// 	ImportState:             true,
			// 	ImportStateVerify:       true,
			// 	ImportStateIdFunc:       importStateIdOrgVdcObject(testConfig, edgeGatewayVcdName),
			// 	ImportStateVerifyIgnore: []string{"external_network", "external_networks"},
			// },
		},
	})
}

const testAccNsxtEdgeGatewayDataSources = `
#data "vcd_nsxt_edge_cluster" "ec" {
#	vdc  = "{{.NsxtVdc}}"
#	name = "{{.ExistingEdgeCluster}}"
#}

data "vcd_external_network_v2" "existing-extnet" {
	name = "nsxt-extnet-dainius"
}

data "vcd_nsxt_manager" "main" {
  name = "nsxManager1"
}
`

const testAccNsxtEdgeGateway = testAccNsxtEdgeGatewayDataSources + `
resource "vcd_nsxt_edgegateway" "nsxt-edge" {
  org                     = "{{.Org}}"
  vdc                     = "{{.NsxtVdc}}"
  name                    = "{{.NsxtEdgeGatewayVcd}}"
  description             = "Description"
#  edge_cluster_id         = data.vcd_nsxt_edge_cluster.ec.id

  #nsxt_manager_id     = data.vcd_nsxt_manager.main.id
  external_network_id = data.vcd_external_network_v2.existing-extnet.id

  subnet {
     gateway               = "10.150.191.253"
     prefix_length         = "19"
     primary_ip            = "10.150.160.137"
     allocated_ips {
       start_address = "10.150.160.137"
       end_address   = "10.150.160.137"
     }
  }
}
`

const testAccNsxtEdgeGatewayUpdate = testAccNsxtEdgeGatewayDataSources + `
resource "vcd_nsxt_edgegateway" "nsxt-edge" {
  org                     = "{{.Org}}"
  vdc                     = "{{.NsxtVdc}}"
  name                    = "{{.NsxtEdgeGatewayVcd}}"
  description             = "Updated-Description"
#  edge_cluster_id         = data.vcd_nsxt_edge_cluster.ec.id

  #nsxt_manager_id     = data.vcd_nsxt_manager.main.id
  external_network_id = data.vcd_external_network_v2.existing-extnet.id

  subnet {
     gateway               = "10.150.191.253"
     prefix_length         = "19"
     primary_ip            = "10.150.160.137"
     allocated_ips {
       start_address = "10.150.160.137"
       end_address   = "10.150.160.137"
     }
  }
}
`
