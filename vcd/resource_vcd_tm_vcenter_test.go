//go:build tm || ALL || functional

package vcd

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccVcdTmVcenter(t *testing.T) {
	preTestChecks(t)

	skipIfNotSysAdmin(t)
	skipIfNotTm(t)

	if !testConfig.Tm.CreateVcenter {
		t.Skipf("Skipping vCenter creation")
	}

	var params = StringMap{
		"Org":             testConfig.Tm.Org,
		"VcenterUsername": testConfig.Tm.VcenterUsername,
		"VcenterPassword": testConfig.Tm.VcenterPassword,
		"VcenterUrl":      testConfig.Tm.VcenterUrl,

		"Testname": t.Name(),

		"Tags": "tm",
	}
	testParamsNotEmpty(t, params)

	configText1 := templateFill(testAccVcdTmVcenterStep1, params)
	// params["FuncName"] = t.Name() + "-step2"
	// configText2 := templateFill(testAccVcdTmVcenterStep2, params)

	debugPrintf("#[DEBUG] CONFIGURATION step2: %s\n", configText1)
	// debugPrintf("#[DEBUG] CONFIGURATION step1: %s\n", configText2)
	if vcdShortTest {
		t.Skip(acceptanceTestsSkipped)
		return
	}

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: configText1,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("vcd_tm_vcenter.test", "id"),
					resource.TestCheckResourceAttr("vcd_tm_vcenter.test", "name", t.Name()),
				),
			},
			// {
			// 	Config: configText2,
			// 	Check: resource.ComposeTestCheckFunc(
			// 		resourceFieldsEqual("vcd_tm_org.test", "data.vcd_tm_org.test", nil),
			// 	),
			// },
			// {
			// 	ResourceName:      "vcd_tm_org.test",
			// 	ImportState:       true,
			// 	ImportStateVerify: true,
			// 	ImportStateId:     params["Org"].(string),
			// },
		},
	})

	postTestChecks(t)
}

const testAccVcdTmVcenterStep1 = `
resource "vcd_tm_vcenter" "test" {
  name       = "{{.Testname}}"
  url        = "{{.VcenterUrl}}"
  username   = "{{.VcenterUsername}}"
  password   = "{{.VcenterPassword}}"
  is_enabled = false
}
`

// const testAccVcdTmVcenterStep2 = testAccVcdTmVcenterStep1 + `
// data "vcd_tm_org" "test" {
//   name = vcd_tm_org.test.name
// }
// `
