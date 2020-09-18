// +build functional network extnetwork nsxt ALL

package vcd

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccVcdExternalNetworkV2Nsxt(t *testing.T) {

	if !usingSysAdmin() {
		t.Skip(t.Name() + " requires system admin privileges")
		return
	}

	startAddress := "192.168.30.51"
	endAddress := "192.168.30.62"
	description := "Test External Network"
	gateway := "192.168.30.49"
	netmask := "24"
	dns1 := "192.168.0.164"
	dns2 := "192.168.0.196"
	var params = StringMap{
		"NsxtManager":     testConfig.Nsxt.Manager,
		"NsxtTier0Router": testConfig.Nsxt.Tier0router,

		"ExternalNetworkName": TestAccVcdExternalNetwork,
		"Type":                testConfig.Networking.ExternalNetworkPortGroupType,
		"PortGroup":           testConfig.Networking.ExternalNetworkPortGroup,
		"Vcenter":             testConfig.Networking.Vcenter,
		"StartAddress":        startAddress,
		"EndAddress":          endAddress,
		"Description":         description,
		"Gateway":             gateway,
		"Netmask":             netmask,
		"Dns1":                dns1,
		"Dns2":                dns2,
		"Tags":                "network extnetwork",
	}

	configText := templateFill(testAccCheckVcdExternalNetworkV2Nsxt, params)
	if vcdShortTest {
		t.Skip(acceptanceTestsSkipped)
		return
	}
	debugPrintf("#[DEBUG] CONFIGURATION: %s", configText)

	resourceName := "vcd_external_network_v2.ext-net-nsxt"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		// CheckDestroy: testAccCheckExternalNetworkDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: configText,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", TestAccVcdExternalNetwork),
					// resource.TestCheckResourceAttr(resourceName, "vsphere_network.0.vcenter", testConfig.Networking.Vcenter),
					// resource.TestCheckResourceAttr(resourceName, "vsphere_network.0.name", testConfig.Networking.ExternalNetworkPortGroup),
					// resource.TestCheckResourceAttr(resourceName, "vsphere_network.0.type", testConfig.Networking.ExternalNetworkPortGroupType),
					// resource.TestCheckResourceAttr(resourceName, "ip_scope.0.gateway", gateway),
					// resource.TestCheckResourceAttr(resourceName, "ip_scope.0.netmask", netmask),
					// resource.TestCheckResourceAttr(resourceName, "ip_scope.0.dns1", dns1),
					// resource.TestCheckResourceAttr(resourceName, "ip_scope.0.dns2", dns2),
					// resource.TestCheckResourceAttr(resourceName, "ip_scope.0.static_ip_pool.0.start_address", startAddress),
					// resource.TestCheckResourceAttr(resourceName, "ip_scope.0.static_ip_pool.0.end_address", endAddress),
					resource.TestCheckResourceAttr(resourceName, "description", description),
					// resource.TestCheckResourceAttr(resourceName, "retain_net_info_across_deployments", "false"),
				),
			},
		},
	})
}

const testAccCheckVcdExternalNetworkV2Nsxt = `
data "vcd_nsxt_manager" "main" {
  name = "{{.NsxtManager}}"
}

data "vcd_nsxt_tier0_router" "router" {
  name            = "{{.NsxtTier0Router}}"
  nsxt_manager_id = data.vcd_nsxt_manager.main.id
}

resource "vcd_external_network_v2" "ext-net-nsxt" {
  name        = "{{.ExternalNetworkName}}"
  description = "{{.Description}}"

  nsxt_network {
    nsxt_manager_id      = data.vcd_nsxt_manager.main.id
    nsxt_tier0_router_id = data.vcd_nsxt_tier0_router.router.id
  }

  ip_scope {
    enabled       = false
    gateway       = "{{.Gateway}}"
    prefix_length = "{{.Netmask}}"

    static_ip_pool {
      start_address = "{{.StartAddress}}"
      end_address   = "{{.EndAddress}}"
    }
  }

  ip_scope {
    gateway       = "14.14.14.1"
    prefix_length = "24"

    static_ip_pool {
      start_address = "14.14.14.10"
      end_address   = "14.14.14.15"
    }
    
    static_ip_pool {
      start_address = "14.14.14.20"
      end_address   = "14.14.14.25"
    }
  }
}
`

func TestAccVcdExternalNetworkV2Nsxv(t *testing.T) {

	if !usingSysAdmin() {
		t.Skip(t.Name() + " requires system admin privileges")
		return
	}

	startAddress := "192.168.30.51"
	endAddress := "192.168.30.62"
	description := "Test External Network"
	gateway := "192.168.30.49"
	netmask := "24"
	dns1 := "192.168.0.164"
	dns2 := "192.168.0.196"
	var params = StringMap{

		"ExternalNetworkName": TestAccVcdExternalNetwork,
		"Type":                testConfig.Networking.ExternalNetworkPortGroupType,
		"PortGroup":           testConfig.Networking.ExternalNetworkPortGroup,
		"Vcenter":             testConfig.Networking.Vcenter,
		"StartAddress":        startAddress,
		"EndAddress":          endAddress,
		"Description":         description,
		"Gateway":             gateway,
		"Netmask":             netmask,
		"Dns1":                dns1,
		"Dns2":                dns2,
		"Tags":                "network extnetwork",
	}

	configText := templateFill(testAccCheckVcdExternalNetworkV2Nsxv, params)
	if vcdShortTest {
		t.Skip(acceptanceTestsSkipped)
		return
	}
	debugPrintf("#[DEBUG] CONFIGURATION: %s", configText)

	resourceName := "vcd_external_network_v2.ext-net-nsxt"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		// CheckDestroy: testAccCheckExternalNetworkDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: configText,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", TestAccVcdExternalNetwork),
					resource.TestCheckResourceAttr(resourceName, "description", description),

					// resource.TestCheckResourceAttr(resourceName, "vsphere_network.0.vcenter", testConfig.Networking.Vcenter),
					// resource.TestCheckResourceAttr(resourceName, "vsphere_network.0.name", testConfig.Networking.ExternalNetworkPortGroup),
					// resource.TestCheckResourceAttr(resourceName, "vsphere_network.0.type", testConfig.Networking.ExternalNetworkPortGroupType),
					// resource.TestCheckResourceAttr(resourceName, "ip_scope.0.gateway", gateway),
					// resource.TestCheckResourceAttr(resourceName, "ip_scope.0.netmask", netmask),
					// resource.TestCheckResourceAttr(resourceName, "ip_scope.0.dns1", dns1),
					// resource.TestCheckResourceAttr(resourceName, "ip_scope.0.dns2", dns2),
					// resource.TestCheckResourceAttr(resourceName, "ip_scope.0.static_ip_pool.0.start_address", startAddress),
					// resource.TestCheckResourceAttr(resourceName, "ip_scope.0.static_ip_pool.0.end_address", endAddress),
					// resource.TestCheckResourceAttr(resourceName, "retain_net_info_across_deployments", "false"),
				),
			},
		},
	})
}

const testAccCheckVcdExternalNetworkV2Nsxv = `
data "vcd_vcenter" "vc" {
  name = "{{.Vcenter}}"
}

data "vcd_portgroup" "sw" {
  name = "{{.PortGroup}}"
  type = "{{.Type}}"
}

resource "vcd_external_network_v2" "ext-net-nsxt" {
  name        = "{{.ExternalNetworkName}}"
  description = "{{.Description}}"

  vsphere_network {
    vcenter_id     = data.vcd_vcenter.vc.id
    portgroup_id   = data.vcd_portgroup.sw.id
  }

  ip_scope {
    gateway       = "{{.Gateway}}"
    prefix_length = "{{.Netmask}}"
    dns1          = "{{.Dns1}}"
    dns2          = "{{.Dns2}}"
    dns_suffix    = "company.biz"

    static_ip_pool {
      start_address = "{{.StartAddress}}"
      end_address   = "{{.EndAddress}}"
    }
  }
}
`
