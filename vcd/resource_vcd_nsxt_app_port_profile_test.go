// +build network nsxt ALL functional

package vcd

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// TestAccVcdNsxtAppPortProfileTenant tests possible options for
func TestAccVcdNsxtAppPortProfileTenant(t *testing.T) {
	preTestChecks(t)
	if vcdShortTest {
		t.Skip(acceptanceTestsSkipped)
		return
	}
	vcdClient := createTemporaryVCDConnection()
	if vcdClient.Client.APIVCDMaxVersionIs("< 34.0") {
		t.Skip(t.Name() + " requires at least API v34.0 (vCD 10.1.1+)")
	}
	skipNoNsxtConfiguration(t)

	// String map to fill the template
	var params = StringMap{
		"Org":         testConfig.VCD.Org,
		"NsxtVdc":     testConfig.Nsxt.Vdc,
		"EdgeGw":      testConfig.Nsxt.EdgeGateway,
		"NetworkName": t.Name(),
		"Tags":        "network",
	}

	configText1 := templateFill(testAccVcdNsxtAppPortProfileTenantStep1, params)
	debugPrintf("#[DEBUG] CONFIGURATION for step 1: %s", configText1)

	params["FuncName"] = t.Name() + "-step2"
	configText2 := templateFill(testAccVcdNsxtAppPortProfileTenantStep2, params)
	debugPrintf("#[DEBUG] CONFIGURATION for step 2: %s", configText2)

	//params["FuncName"] = t.Name() + "-step2"
	//configText2 := templateFill(TestAccVcdNetworkRoutedV2NsxtStep2, params)
	//debugPrintf("#[DEBUG] CONFIGURATION for step 2: %s", configText2)
	//
	//params["FuncName"] = t.Name() + "-step3"
	//configText3 := templateFill(TestAccVcdNetworkRoutedV2NsxtStep3, params)
	//debugPrintf("#[DEBUG] CONFIGURATION for step 3: %s", configText3)

	// Ensure the resource is never recreated - ID stays the same
	//cachedId := &testCachedFieldValue{}

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviders,
		PreCheck:          func() { testAccPreCheck(t) },
		//CheckDestroy:      testAccCheckOpenApiVcdNetworkDestroy(testConfig.Nsxt.Vdc, t.Name()),
		Steps: []resource.TestStep{
			resource.TestStep{ // step 1
				Config: configText1,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("vcd_nsxt_app_port_profile.LDAP", "id"),
					resource.TestCheckResourceAttr("vcd_nsxt_app_port_profile.LDAP", "name", "ldap_app_prof"),
					resource.TestCheckResourceAttr("vcd_nsxt_app_port_profile.LDAP", "description", "Application port profile for LDAP"),
					resource.TestCheckResourceAttr("vcd_nsxt_app_port_profile.LDAP", "scope", "TENANT"),
					resource.TestCheckTypeSetElemNestedAttrs("vcd_nsxt_app_port_profile.LDAP", "app_port.*", map[string]string{
						"protocol": "ICMPv4",
					}),
				),
			},
			resource.TestStep{ // step 1
				Config: configText2,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("vcd_nsxt_app_port_profile.LDAP", "id"),
					resource.TestCheckResourceAttr("vcd_nsxt_app_port_profile.LDAP", "name", "ldap_app_prof-updated"),
					resource.TestCheckResourceAttr("vcd_nsxt_app_port_profile.LDAP", "description", "Application port profile for LDAP-updated"),
					resource.TestCheckResourceAttr("vcd_nsxt_app_port_profile.LDAP", "scope", "TENANT"),
					resource.TestCheckTypeSetElemNestedAttrs("vcd_nsxt_app_port_profile.LDAP", "app_port.*", map[string]string{
						"protocol": "ICMPv6",
					}),
					resource.TestCheckTypeSetElemNestedAttrs("vcd_nsxt_app_port_profile.LDAP", "app_port.*", map[string]string{
						"protocol": "TCP",
					}),
					resource.TestCheckTypeSetElemNestedAttrs("vcd_nsxt_app_port_profile.LDAP", "app_port.*", map[string]string{
						"protocol": "UDP",
					}),
					resource.TestCheckTypeSetElemAttr("vcd_nsxt_app_port_profile.LDAP", "app_port.*.port.*", "2000"),
					resource.TestCheckTypeSetElemAttr("vcd_nsxt_app_port_profile.LDAP", "app_port.*.port.*", "2010-2020"),
					resource.TestCheckTypeSetElemAttr("vcd_nsxt_app_port_profile.LDAP", "app_port.*.port.*", "12345"),
					resource.TestCheckTypeSetElemAttr("vcd_nsxt_app_port_profile.LDAP", "app_port.*.port.*", "65000"),
					resource.TestCheckTypeSetElemAttr("vcd_nsxt_app_port_profile.LDAP", "app_port.*.port.*", "40000-60000"),
				),
			},
		},
	})
	postTestChecks(t)
}

const testAccVcdNsxtAppPortProfileTenantStep1 = `
resource "vcd_nsxt_app_port_profile" "LDAP" {
  org  = "{{.Org}}"
  vdc  = "{{.NsxtVdc}}"
  name = "ldap_app_prof"

  description = "Application port profile for LDAP"
  scope = "TENANT"
  
  app_port {
    protocol = "ICMPv4"
  }
}
`

const testAccVcdNsxtAppPortProfileTenantStep2 = `
resource "vcd_nsxt_app_port_profile" "LDAP" {
  org  = "{{.Org}}"
  vdc  = "{{.NsxtVdc}}"
  name = "ldap_app_prof-updated"

  description = "Application port profile for LDAP-updated"
  scope = "TENANT"
  
  app_port {
    protocol = "ICMPv6"
  }

  app_port {
    protocol = "TCP"
	port = ["2000", "2010-2020", "12345", "65000"]
  }

  app_port {
    protocol = "UDP"
    port = ["40000-60000"]
  }
}
`
