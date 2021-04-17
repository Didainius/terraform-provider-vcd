// +build network nsxt ALL functional

package vcd

import (
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// TestAccVcdNsxtSecurityGroupEmpty tests out capabilities to setup Security Groups without
// attaching member networks
func TestAccVcdNsxtSecurityGroupEmpty(t *testing.T) {
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

	configText := templateFill(testAccNsxtSecurityGroupEmpty, params)
	debugPrintf("#[DEBUG] CONFIGURATION for step 1: %s", configText)

	params["FuncName"] = t.Name() + "-step1"
	configText1 := templateFill(testAccNsxtSecurityGroupEmpty2, params)
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

const testAccNsxtSecurityGroupPrereqsEmpty = `
data "vcd_nsxt_edgegateway" "existing" {
	org  = "{{.Org}}"
	vdc  = "{{.NsxtVdc}}"

	name = "{{.EdgeGw}}"
}
`

const testAccNsxtSecurityGroupEmpty = testAccNsxtSecurityGroupPrereqs + `
resource "vcd_nsxt_security_group" "group1" {
  org  = "{{.Org}}"
  vdc  = "{{.NsxtVdc}}"

  edge_gateway_id = data.vcd_nsxt_edgegateway.existing.id

  name = "test-security-group"
  description = "test-security-group-description"
}
`

const testAccNsxtSecurityGroupEmpty2 = testAccNsxtSecurityGroupPrereqs + `
resource "vcd_nsxt_security_group" "group1" {
  org  = "{{.Org}}"
  vdc  = "{{.NsxtVdc}}"

  edge_gateway_id = data.vcd_nsxt_edgegateway.existing.id

  name = "test-security-group-changed"
}
`

// TestAccVcdNsxtSecurityGroup is similar to TestAccVcdNsxtFirewallGroupEmpty, but it also creates
// Org VDC networks and attaches them to security group.

// Additionally it tests `vcd_nsxt_security_group` datasource to save testing time and avoid creating
// the same prerequisite resources.
func TestAccVcdNsxtSecurityGroup(t *testing.T) {
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

	configText := templateFill(testAccNsxtSecurityGroup, params)
	debugPrintf("#[DEBUG] CONFIGURATION for step 1: %s", configText)

	params["FuncName"] = t.Name() + "-step2"
	configText1 := templateFill(testAccNsxtSecurityGroupDatasource, params)
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
					resource.TestCheckTypeSetElemNestedAttrs("vcd_nsxt_security_group.group1", "member_vms.*", map[string]string{
						"vm_name":   "vapp-vm",
						"vapp_name": "web",
					}),
					resource.TestCheckTypeSetElemNestedAttrs("vcd_nsxt_security_group.group1", "member_vms.*", map[string]string{
						"vm_name":   "standalone-VM",
						"vapp_name": "",
					}),
				),
			},
			resource.TestStep{
				Config: configText1,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("vcd_nsxt_security_group.group1", "id", regexp.MustCompile(`^urn:vcloud:firewallGroup:.*$`)),
					resource.TestCheckResourceAttr("vcd_nsxt_security_group.group1", "name", "test-security-group"),
					resource.TestCheckResourceAttr("vcd_nsxt_security_group.group1", "description", "test-security-group-description"),
					resource.TestCheckTypeSetElemNestedAttrs("vcd_nsxt_security_group.group1", "member_vms.*", map[string]string{
						"vm_name":   "vapp-vm",
						"vapp_name": "web",
					}),
					resource.TestCheckTypeSetElemNestedAttrs("vcd_nsxt_security_group.group1", "member_vms.*", map[string]string{
						"vm_name":   "standalone-VM",
						"vapp_name": "",
					}),
					// Ensure datasource has all the fields
					resourceFieldsEqual("vcd_nsxt_security_group.group1", "data.vcd_nsxt_security_group.group1", []string{}),
				),
			},
			// resource.TestStep{
			// 	Config: configText1,
			// 	Check: resource.ComposeAggregateTestCheckFunc(
			// 		resource.TestMatchResourceAttr("vcd_nsxt_security_group.group1", "id", regexp.MustCompile(`^urn:vcloud:firewallGroup:.*$`)),
			// 		resource.TestCheckResourceAttr("vcd_nsxt_security_group.group1", "name", "test-security-group-changed"),
			// 		resource.TestCheckResourceAttr("vcd_nsxt_security_group.group1", "description", ""),
			// 		resource.TestCheckNoResourceAttr("vcd_nsxt_security_group.group1", "member_org_network_ids"),
			// 		resource.TestCheckNoResourceAttr("vcd_nsxt_security_group.group1", "member_vm_ids"),
			// 	),
			// },
			resource.TestStep{
				ResourceName:      "vcd_nsxt_security_group.group1",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: importStateIdNsxtEdgeGatewayObject(testConfig, testConfig.Nsxt.EdgeGateway, "test-security-group"),
			},
		},
	})
	postTestChecks(t)
}

func sleepTester() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		fmt.Println("sleeping")
		time.Sleep(2 * time.Minute)
		return nil
	}
}

const testAccNsxtSecurityGroupPrereqs = testAccNsxtSecurityGroupPrereqsEmpty + `
resource "vcd_network_routed_v2" "nsxt-backed" {
	# This value could be larger to test more members, but left 2 for the sake of testing speed
	count = 2

	org  = "{{.Org}}"
	vdc  = "{{.NsxtVdc}}"
	name        = "nsxt-routed-${count.index}"
	description = "My routed Org VDC network backed by NSX-T"
  
	edge_gateway_id = data.vcd_nsxt_edgegateway.existing.id
  
	gateway       = "212.1.${count.index}.1"
	prefix_length = 24
  
	static_ip_pool {
	  start_address = "212.1.${count.index}.10"
	  end_address   = "212.1.${count.index}.20"
	}
}

# Create stanadlone VM to check membership
resource "vcd_vm" "emptyVM" {
	org  = "{{.Org}}"
	vdc  = "{{.NsxtVdc}}"

	name          = "standalone-VM"
	computer_name = "emptyVM"
	memory        = 2048
	cpus          = 2
	cpu_cores     = 1
  
	os_type = "sles10_64Guest"
	hardware_version = "vmx-14"
	
	network {
		type               = "org"
		name               = vcd_network_routed_v2.nsxt-backed[0].name
		ip_allocation_mode = "POOL"
		is_primary         = true
	}

	depends_on = [vcd_network_routed_v2.nsxt-backed]
  }

# Create a vApp and VM
resource "vcd_vapp" "web" {
  org  = "{{.Org}}"
  vdc  = "{{.NsxtVdc}}"

  name = "web"
}

resource "vcd_vapp_org_network" "vappOrgNet" {
	org  = "{{.Org}}"
	vdc  = "{{.NsxtVdc}}"
  
	vapp_name         = vcd_vapp.web.name
  
   # Comment below line to create an isolated vApp network
	org_network_name  = vcd_network_routed_v2.nsxt-backed[1].name

	depends_on = [vcd_vapp.web,vcd_network_routed_v2.nsxt-backed]
  }

resource "vcd_vapp_vm" "emptyVM" {
  org  = "{{.Org}}"
  vdc  = "{{.NsxtVdc}}"

  vapp_name     = vcd_vapp.web.name
  name          = "vapp-vm"
  computer_name = "emptyVM"
  memory        = 2048
  cpus          = 2
  cpu_cores     = 1

  os_type = "sles10_64Guest"
  hardware_version = "vmx-14"

  network {
	type               = "org"
	name               = vcd_network_routed_v2.nsxt-backed[1].name
	ip_allocation_mode = "POOL"
	is_primary         = true
  }

  depends_on = [vcd_vapp_org_network.vappOrgNet]
}
`

const testAccNsxtSecurityGroup = testAccNsxtSecurityGroupPrereqs + `
resource "vcd_nsxt_security_group" "group1" {
  org  = "{{.Org}}"
  vdc  = "{{.NsxtVdc}}"

  edge_gateway_id = data.vcd_nsxt_edgegateway.existing.id

  name = "test-security-group"
  description = "test-security-group-description"

  member_org_network_ids = vcd_network_routed_v2.nsxt-backed.*.id

  depends_on = [vcd_vapp_vm.emptyVM, vcd_vm.emptyVM]
}
`

const testAccNsxtSecurityGroupDatasource = testAccNsxtSecurityGroup + `
data "vcd_nsxt_security_group" "group1" {
	org  = "{{.Org}}"
	vdc  = "{{.NsxtVdc}}"

	edge_gateway_id = data.vcd_nsxt_edgegateway.existing.id
	name            = "test-security-group"
}
`
