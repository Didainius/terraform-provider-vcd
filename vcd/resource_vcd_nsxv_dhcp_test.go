// +build gateway ALL functional

package vcd

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccVcdNsxvDhcpPool(t *testing.T) {

	// String map to fill the template
	var params = StringMap{
		"Org":         testConfig.VCD.Org,
		"Vdc":         testConfig.VCD.Vdc,
		"EdgeGateway": testConfig.Networking.EdgeGateway,
		"Tags":        "gateway",
	}

	configText1 := templateFill(testAccVcdNsxvDhcpPool, params)
	debugPrintf("#[DEBUG] CONFIGURATION for step 1: %s", configText1)
	//
	params["FuncName"] = t.Name() + "-step1"
	configText2 := templateFill(testAccVcdNsxvDhcpPool2, params)
	debugPrintf("#[DEBUG] CONFIGURATION for step 2: %s", configText2)

	params["FuncName"] = t.Name() + "-step2"
	configText3 := templateFill(testAccVcdNsxvDhcpPool3, params)
	debugPrintf("#[DEBUG] CONFIGURATION for step 3: %s", configText3)

	if vcdShortTest {
		t.Skip(acceptanceTestsSkipped)
		return
	}

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviders,
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy:      testAccCheckVcdDhcpRelaySettingsEmpty(),
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: configText1,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("vcd_nsxv_dhcp.pool_config", "id", regexp.MustCompile(`^urn:vcloud:gateway:`)),
					resource.TestCheckResourceAttr("vcd_nsxv_dhcp.pool_config", "dhcp_enabled", "true"),
					resource.TestCheckTypeSetElemNestedAttrs("vcd_nsxv_dhcp.pool_config", "dhcp_pool.*", map[string]string{
						"ip_range":    "192.168.2.252-192.168.2.253",
						"domain_name": "simple.hostname",
						"gateway":     "192.168.2.1",
						"subnet_mask": "255.255.255.0",
						"lease_time":  "86400",
					}),
				),
			},
			resource.TestStep{
				Config: configText2,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("vcd_nsxv_dhcp.pool_config", "id", regexp.MustCompile(`^urn:vcloud:gateway:`)),
					resource.TestCheckResourceAttr("vcd_nsxv_dhcp.pool_config", "dhcp_enabled", "false"),

					resource.TestCheckTypeSetElemNestedAttrs("vcd_nsxv_dhcp.pool_config", "dhcp_pool.*", map[string]string{
						"ip_range":    "192.168.2.252-192.168.2.253",
						"domain_name": "simple.hostname",
						"gateway":     "192.168.2.1",
						"subnet_mask": "255.255.255.0",
						"lease_time":  "86400",
					}),
				),
			},

			resource.TestStep{
				Config: configText3,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("vcd_nsxv_dhcp.pool_config", "id", regexp.MustCompile(`^urn:vcloud:gateway:`)),
					resource.TestCheckResourceAttr("vcd_nsxv_dhcp.pool_config", "dhcp_enabled", "true"),

					resource.TestCheckTypeSetElemNestedAttrs("vcd_nsxv_dhcp.pool_config", "dhcp_pool.*", map[string]string{
						"ip_range":    "192.168.2.252-192.168.2.253",
						"domain_name": "simple.hostname",
						"gateway":     "192.168.2.1",
						"subnet_mask": "255.255.255.0",
						"lease_time":  "infinite",
					}),
					resource.TestCheckTypeSetElemNestedAttrs("vcd_nsxv_dhcp.pool_config", "dhcp_pool.*", map[string]string{
						"ip_range":    "192.168.2.250-192.168.2.251",
						"domain_name": "tricky.hostname",
						"gateway":     "192.168.2.1",
						"subnet_mask": "255.255.255.0",
						"lease_time":  "86400",
					}),
				),
			},
			// resource.TestStep{
			// 	ResourceName:      "vcd_nsxv_dhcp.pool_config",
			// 	ImportState:       true,
			// 	ImportStateVerify: true,
			// 	ImportStateId:     testConfig.VCD.Org + "." + testConfig.VCD.Vdc + "." + testConfig.Networking.EdgeGateway,
			// },
		},
	})
}

const testAccVcdNsxvDhcpPoolEgw = `
data "vcd_edgegateway" "nsxv" {
  org          = "{{.Org}}"
  vdc          = "{{.Vdc}}"

  name = "{{.EdgeGateway}}"
}
`

const testAccVcdNsxvDhcpPool = testAccVcdNsxvDhcpPoolEgw + `
resource "vcd_nsxv_dhcp" "pool_config" {
  org          = "{{.Org}}"
  vdc          = "{{.Vdc}}"
  edge_gateway_id = data.vcd_edgegateway.nsxv.id
  dhcp_enabled = true
  dhcp_pool {
    ip_range    = "192.168.2.252-192.168.2.253"
    domain_name = "simple.hostname"
    gateway     = "192.168.2.1"
    subnet_mask = "255.255.255.0"
  }
}
`

const testAccVcdNsxvDhcpPool2 = testAccVcdNsxvDhcpPoolEgw + `
resource "vcd_nsxv_dhcp" "pool_config" {
  org          = "{{.Org}}"
  vdc          = "{{.Vdc}}"
  edge_gateway_id = data.vcd_edgegateway.nsxv.id

  dhcp_enabled = false
  dhcp_pool {
    ip_range    = "192.168.2.252-192.168.2.253"
    domain_name = "simple.hostname"
    gateway     = "192.168.2.1"
    subnet_mask = "255.255.255.0"
	lease_time  = 86400
  }
}
`

const testAccVcdNsxvDhcpPool3 = testAccVcdNsxvDhcpPoolEgw + `
resource "vcd_nsxv_dhcp" "pool_config" {
  org          = "{{.Org}}"
  vdc          = "{{.Vdc}}"
  edge_gateway_id = data.vcd_edgegateway.nsxv.id

  dhcp_enabled = true

  dhcp_pool {
    ip_range    = "192.168.2.252-192.168.2.253"
    domain_name = "simple.hostname"
    gateway     = "192.168.2.1"
    subnet_mask = "255.255.255.0"
	lease_time  = "infinite"
  }

dhcp_pool {
    ip_range    = "192.168.2.250-192.168.2.251"
    domain_name = "tricky.hostname"
    gateway     = "192.168.2.1"
    subnet_mask = "255.255.255.0"
	lease_time  = 86400
  }
}
`
