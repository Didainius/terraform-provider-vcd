// +build network nsxt ALL functional

package vcd

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccVcdNsxtNatRuleDnat(t *testing.T) {
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

	configText := templateFill(testAccNsxtNatDnat, params)
	debugPrintf("#[DEBUG] CONFIGURATION for step 1: %s", configText)

	//params["FuncName"] = t.Name() + "-step1"
	//configText1 := templateFill(testAccNsxtSecurityGroupEmpty2, params)
	//debugPrintf("#[DEBUG] CONFIGURATION for step 2: %s", configText1)

	if vcdShortTest {
		t.Skip(acceptanceTestsSkipped)
		return
	}

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviders,
		PreCheck:          func() { testAccPreCheck(t) },
		//CheckDestroy: resource.ComposeAggregateTestCheckFunc(
		//	testAccCheckNsxtFirewallGroupDestroy(testConfig.Nsxt.Vdc, "test-security-group", types.FirewallGroupTypeSecurityGroup),
		//	testAccCheckNsxtFirewallGroupDestroy(testConfig.Nsxt.Vdc, "test-security-group-changed", types.FirewallGroupTypeSecurityGroup),
		//),
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: configText,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("vcd_nsxt_nat_rule.dnat", "id"),
					//resource.TestCheckResourceAttr("vcd_nsxt_security_group.group1", "name", "test-security-group"),
					//resource.TestCheckResourceAttr("vcd_nsxt_security_group.group1", "description", "test-security-group-description"),
					//resource.TestCheckNoResourceAttr("vcd_nsxt_security_group.group1", "member_org_network_ids"),
					//resource.TestCheckNoResourceAttr("vcd_nsxt_security_group.group1", "member_vm_ids"),
				),
			},
			//resource.TestStep{
			//	Config: configText1,
			//	Check: resource.ComposeAggregateTestCheckFunc(
			//		resource.TestMatchResourceAttr("vcd_nsxt_security_group.group1", "id", regexp.MustCompile(`^urn:vcloud:firewallGroup:.*$`)),
			//		resource.TestCheckResourceAttr("vcd_nsxt_security_group.group1", "name", "test-security-group-changed"),
			//		resource.TestCheckResourceAttr("vcd_nsxt_security_group.group1", "description", ""),
			//		resource.TestCheckNoResourceAttr("vcd_nsxt_security_group.group1", "member_org_network_ids"),
			//		resource.TestCheckNoResourceAttr("vcd_nsxt_security_group.group1", "member_vm_ids"),
			//	),
			//},
			//resource.TestStep{
			//	ResourceName:      "vcd_nsxt_security_group.group1",
			//	ImportState:       true,
			//	ImportStateVerify: true,
			//	ImportStateIdFunc: importStateIdNsxtEdgeGatewayObject(testConfig, testConfig.Nsxt.EdgeGateway, "test-security-group-changed"),
			//},
		},
	})
	postTestChecks(t)
}

const testAccNsxtNatDnat = testAccNsxtSecurityGroupPrereqsEmpty + `
resource "vcd_nsxt_nat_rule" "dnat" {
  org  = "{{.Org}}"
  vdc  = "{{.NsxtVdc}}"

  edge_gateway_id = data.vcd_nsxt_edgegateway.existing.id

  name = "test-dnat-rule"
  rule_type = "DNAT"
  description = "description"
  
  # Using primary_ip from edge gateway
  external_addresses = tolist(data.vcd_nsxt_edgegateway.existing.subnet)[0].primary_ip
  internal_addresses = "11.11.11.2"
  logging = true
}
`

func TestAccVcdNsxtNatRuleNoDnat(t *testing.T) {
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

	configText := templateFill(testAccNsxtNatNoDnat, params)
	debugPrintf("#[DEBUG] CONFIGURATION for step 1: %s", configText)

	//params["FuncName"] = t.Name() + "-step1"
	//configText1 := templateFill(testAccNsxtSecurityGroupEmpty2, params)
	//debugPrintf("#[DEBUG] CONFIGURATION for step 2: %s", configText1)

	if vcdShortTest {
		t.Skip(acceptanceTestsSkipped)
		return
	}

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviders,
		PreCheck:          func() { testAccPreCheck(t) },
		//CheckDestroy: resource.ComposeAggregateTestCheckFunc(
		//	testAccCheckNsxtFirewallGroupDestroy(testConfig.Nsxt.Vdc, "test-security-group", types.FirewallGroupTypeSecurityGroup),
		//	testAccCheckNsxtFirewallGroupDestroy(testConfig.Nsxt.Vdc, "test-security-group-changed", types.FirewallGroupTypeSecurityGroup),
		//),
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: configText,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("vcd_nsxt_nat_rule.no-dnat", "id"),
					//resource.TestCheckResourceAttr("vcd_nsxt_security_group.group1", "name", "test-security-group"),
					//resource.TestCheckResourceAttr("vcd_nsxt_security_group.group1", "description", "test-security-group-description"),
					//resource.TestCheckNoResourceAttr("vcd_nsxt_security_group.group1", "member_org_network_ids"),
					//resource.TestCheckNoResourceAttr("vcd_nsxt_security_group.group1", "member_vm_ids"),
				),
			},
			//resource.TestStep{
			//	Config: configText1,
			//	Check: resource.ComposeAggregateTestCheckFunc(
			//		resource.TestMatchResourceAttr("vcd_nsxt_security_group.group1", "id", regexp.MustCompile(`^urn:vcloud:firewallGroup:.*$`)),
			//		resource.TestCheckResourceAttr("vcd_nsxt_security_group.group1", "name", "test-security-group-changed"),
			//		resource.TestCheckResourceAttr("vcd_nsxt_security_group.group1", "description", ""),
			//		resource.TestCheckNoResourceAttr("vcd_nsxt_security_group.group1", "member_org_network_ids"),
			//		resource.TestCheckNoResourceAttr("vcd_nsxt_security_group.group1", "member_vm_ids"),
			//	),
			//},
			//resource.TestStep{
			//	ResourceName:      "vcd_nsxt_security_group.group1",
			//	ImportState:       true,
			//	ImportStateVerify: true,
			//	ImportStateIdFunc: importStateIdNsxtEdgeGatewayObject(testConfig, testConfig.Nsxt.EdgeGateway, "test-security-group-changed"),
			//},
		},
	})
	postTestChecks(t)
}

const testAccNsxtNatNoDnat = testAccNsxtSecurityGroupPrereqsEmpty + `
resource "vcd_nsxt_nat_rule" "no-dnat" {
  org  = "{{.Org}}"
  vdc  = "{{.NsxtVdc}}"

  edge_gateway_id = data.vcd_nsxt_edgegateway.existing.id

  name = "test-dnat-rule"
  rule_type = "NO_DNAT"
  description = "description"
  
  # Using primary_ip from edge gateway
  external_addresses = tolist(data.vcd_nsxt_edgegateway.existing.subnet)[0].primary_ip
  //internal_addresses = "11.11.11.2"
  logging = true
}
`

func TestAccVcdNsxtNatRuleSnat(t *testing.T) {
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

	configText := templateFill(testAccNsxtNatSnat, params)
	debugPrintf("#[DEBUG] CONFIGURATION for step 1: %s", configText)

	//params["FuncName"] = t.Name() + "-step1"
	//configText1 := templateFill(testAccNsxtSecurityGroupEmpty2, params)
	//debugPrintf("#[DEBUG] CONFIGURATION for step 2: %s", configText1)

	if vcdShortTest {
		t.Skip(acceptanceTestsSkipped)
		return
	}

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviders,
		PreCheck:          func() { testAccPreCheck(t) },
		//CheckDestroy: resource.ComposeAggregateTestCheckFunc(
		//	testAccCheckNsxtFirewallGroupDestroy(testConfig.Nsxt.Vdc, "test-security-group", types.FirewallGroupTypeSecurityGroup),
		//	testAccCheckNsxtFirewallGroupDestroy(testConfig.Nsxt.Vdc, "test-security-group-changed", types.FirewallGroupTypeSecurityGroup),
		//),
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: configText,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("vcd_nsxt_nat_rule.snat", "id"),
					//resource.TestCheckResourceAttr("vcd_nsxt_security_group.group1", "name", "test-security-group"),
					//resource.TestCheckResourceAttr("vcd_nsxt_security_group.group1", "description", "test-security-group-description"),
					//resource.TestCheckNoResourceAttr("vcd_nsxt_security_group.group1", "member_org_network_ids"),
					//resource.TestCheckNoResourceAttr("vcd_nsxt_security_group.group1", "member_vm_ids"),
				),
			},
			//resource.TestStep{
			//	Config: configText1,
			//	Check: resource.ComposeAggregateTestCheckFunc(
			//		resource.TestMatchResourceAttr("vcd_nsxt_security_group.group1", "id", regexp.MustCompile(`^urn:vcloud:firewallGroup:.*$`)),
			//		resource.TestCheckResourceAttr("vcd_nsxt_security_group.group1", "name", "test-security-group-changed"),
			//		resource.TestCheckResourceAttr("vcd_nsxt_security_group.group1", "description", ""),
			//		resource.TestCheckNoResourceAttr("vcd_nsxt_security_group.group1", "member_org_network_ids"),
			//		resource.TestCheckNoResourceAttr("vcd_nsxt_security_group.group1", "member_vm_ids"),
			//	),
			//},
			//resource.TestStep{
			//	ResourceName:      "vcd_nsxt_security_group.group1",
			//	ImportState:       true,
			//	ImportStateVerify: true,
			//	ImportStateIdFunc: importStateIdNsxtEdgeGatewayObject(testConfig, testConfig.Nsxt.EdgeGateway, "test-security-group-changed"),
			//},
		},
	})
	postTestChecks(t)
}

const testAccNsxtNatSnat = testAccNsxtSecurityGroupPrereqsEmpty + `
resource "vcd_nsxt_nat_rule" "snat" {
  org  = "{{.Org}}"
  vdc  = "{{.NsxtVdc}}"
	
  edge_gateway_id = data.vcd_nsxt_edgegateway.existing.id

  name        = "test-dnat-rule"
  rule_type   = "SNAT"
  description = "description"
  
  # Using primary_ip from edge gateway
  external_addresses         = tolist(data.vcd_nsxt_edgegateway.existing.subnet)[0].primary_ip
  internal_addresses         = "11.11.11.2"
  snat_destination_addresses = "11.11.11.4"
  logging = true
}
`

func TestAccVcdNsxtNatRuleNoSnat(t *testing.T) {
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

	configText := templateFill(testAccNsxtNatNoSnat, params)
	debugPrintf("#[DEBUG] CONFIGURATION for step 1: %s", configText)

	//params["FuncName"] = t.Name() + "-step1"
	//configText1 := templateFill(testAccNsxtSecurityGroupEmpty2, params)
	//debugPrintf("#[DEBUG] CONFIGURATION for step 2: %s", configText1)

	if vcdShortTest {
		t.Skip(acceptanceTestsSkipped)
		return
	}

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviders,
		PreCheck:          func() { testAccPreCheck(t) },
		//CheckDestroy: resource.ComposeAggregateTestCheckFunc(
		//	testAccCheckNsxtFirewallGroupDestroy(testConfig.Nsxt.Vdc, "test-security-group", types.FirewallGroupTypeSecurityGroup),
		//	testAccCheckNsxtFirewallGroupDestroy(testConfig.Nsxt.Vdc, "test-security-group-changed", types.FirewallGroupTypeSecurityGroup),
		//),
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: configText,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("vcd_nsxt_nat_rule.snat", "id"),
					//resource.TestCheckResourceAttr("vcd_nsxt_security_group.group1", "name", "test-security-group"),
					//resource.TestCheckResourceAttr("vcd_nsxt_security_group.group1", "description", "test-security-group-description"),
					//resource.TestCheckNoResourceAttr("vcd_nsxt_security_group.group1", "member_org_network_ids"),
					//resource.TestCheckNoResourceAttr("vcd_nsxt_security_group.group1", "member_vm_ids"),
				),
			},
			//resource.TestStep{
			//	Config: configText1,
			//	Check: resource.ComposeAggregateTestCheckFunc(
			//		resource.TestMatchResourceAttr("vcd_nsxt_security_group.group1", "id", regexp.MustCompile(`^urn:vcloud:firewallGroup:.*$`)),
			//		resource.TestCheckResourceAttr("vcd_nsxt_security_group.group1", "name", "test-security-group-changed"),
			//		resource.TestCheckResourceAttr("vcd_nsxt_security_group.group1", "description", ""),
			//		resource.TestCheckNoResourceAttr("vcd_nsxt_security_group.group1", "member_org_network_ids"),
			//		resource.TestCheckNoResourceAttr("vcd_nsxt_security_group.group1", "member_vm_ids"),
			//	),
			//},
			//resource.TestStep{
			//	ResourceName:      "vcd_nsxt_security_group.group1",
			//	ImportState:       true,
			//	ImportStateVerify: true,
			//	ImportStateIdFunc: importStateIdNsxtEdgeGatewayObject(testConfig, testConfig.Nsxt.EdgeGateway, "test-security-group-changed"),
			//},
		},
	})
	postTestChecks(t)
}

const testAccNsxtNatNoSnat = testAccNsxtSecurityGroupPrereqsEmpty + `
resource "vcd_nsxt_nat_rule" "snat" {
  org  = "{{.Org}}"
  vdc  = "{{.NsxtVdc}}"

  edge_gateway_id = data.vcd_nsxt_edgegateway.existing.id

  name        = "test-dnat-rule"
  rule_type   = "NO_SNAT"
  description = "description"
  
  # Using primary_ip from edge gateway
  //external_addresses         = tolist(data.vcd_nsxt_edgegateway.existing.subnet)[0].primary_ip
  internal_addresses         = "11.11.11.2"
  //snat_destination_addresses = "11.11.11.4"
  logging = true
}
`

// TestAccVcdNsxtNatRuleFirewallMatchPriority explicitly tests support for two new fields introduces in API 35.2 (VCD 10.2.2)
// firewall_match and priority. For 10.2.2 versions this should work, while for lower versions it should return an error.
func TestAccVcdNsxtNatRuleFirewallMatchPriority(t *testing.T) {
	preTestChecks(t)
	skipNoNsxtConfiguration(t)

	// expectError must stay nil for versions > 10.2.2, but should match anything for lower versions
	var expectError *regexp.Regexp

	client := createTemporaryVCDConnection()
	if client.Client.APIVCDMaxVersionIs("< 35.2") {
		expectError = regexp.MustCompile(`firewall_match and priority fields can only be set for VCD 10.2.2+`)
	}

	// String map to fill the template
	var params = StringMap{
		"Org":         testConfig.VCD.Org,
		"NsxtVdc":     testConfig.Nsxt.Vdc,
		"EdgeGw":      testConfig.Nsxt.EdgeGateway,
		"NetworkName": t.Name(),
		"Tags":        "network nsxt",
	}

	configText := templateFill(testAccNsxtNatFirewallMatchPriority, params)
	debugPrintf("#[DEBUG] CONFIGURATION for step 1: %s", configText)

	//params["FuncName"] = t.Name() + "-step1"
	//configText1 := templateFill(testAccNsxtSecurityGroupEmpty2, params)
	//debugPrintf("#[DEBUG] CONFIGURATION for step 2: %s", configText1)

	if vcdShortTest {
		t.Skip(acceptanceTestsSkipped)
		return
	}

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviders,
		PreCheck:          func() { testAccPreCheck(t) },
		//CheckDestroy: resource.ComposeAggregateTestCheckFunc(
		//	testAccCheckNsxtFirewallGroupDestroy(testConfig.Nsxt.Vdc, "test-security-group", types.FirewallGroupTypeSecurityGroup),
		//	testAccCheckNsxtFirewallGroupDestroy(testConfig.Nsxt.Vdc, "test-security-group-changed", types.FirewallGroupTypeSecurityGroup),
		//),
		Steps: []resource.TestStep{
			resource.TestStep{
				Config:      configText,
				ExpectError: expectError,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("vcd_nsxt_nat_rule.dnat", "id"),
					//resource.TestCheckResourceAttr("vcd_nsxt_security_group.group1", "name", "test-security-group"),
					//resource.TestCheckResourceAttr("vcd_nsxt_security_group.group1", "description", "test-security-group-description"),
					//resource.TestCheckNoResourceAttr("vcd_nsxt_security_group.group1", "member_org_network_ids"),
					//resource.TestCheckNoResourceAttr("vcd_nsxt_security_group.group1", "member_vm_ids"),
				),
			},
			//resource.TestStep{
			//	Config: configText1,
			//	Check: resource.ComposeAggregateTestCheckFunc(
			//		resource.TestMatchResourceAttr("vcd_nsxt_security_group.group1", "id", regexp.MustCompile(`^urn:vcloud:firewallGroup:.*$`)),
			//		resource.TestCheckResourceAttr("vcd_nsxt_security_group.group1", "name", "test-security-group-changed"),
			//		resource.TestCheckResourceAttr("vcd_nsxt_security_group.group1", "description", ""),
			//		resource.TestCheckNoResourceAttr("vcd_nsxt_security_group.group1", "member_org_network_ids"),
			//		resource.TestCheckNoResourceAttr("vcd_nsxt_security_group.group1", "member_vm_ids"),
			//	),
			//},
			//resource.TestStep{
			//	ResourceName:      "vcd_nsxt_security_group.group1",
			//	ImportState:       true,
			//	ImportStateVerify: true,
			//	ImportStateIdFunc: importStateIdNsxtEdgeGatewayObject(testConfig, testConfig.Nsxt.EdgeGateway, "test-security-group-changed"),
			//},
		},
	})
	postTestChecks(t)
}

const testAccNsxtNatFirewallMatchPriority = testAccNsxtSecurityGroupPrereqsEmpty + `
resource "vcd_nsxt_nat_rule" "dnat" {
  org  = "{{.Org}}"
  vdc  = "{{.NsxtVdc}}"

  edge_gateway_id = data.vcd_nsxt_edgegateway.existing.id

  name = "test-dnat-rule"
  rule_type = "DNAT"
  description = "description"
  
  # Using primary_ip from edge gateway
  external_addresses = tolist(data.vcd_nsxt_edgegateway.existing.subnet)[0].primary_ip
  internal_addresses = "11.11.11.2"
  logging = true

  firewall_match = "MATCH_INTERNAL_ADDRESS"
  priority       = 100
}
`

//// TestAccVcdNsxtNatRuleReflexive (VCD 10.2.2+)
//func TestAccVcdNsxtNatRuleReflexive(t *testing.T) {
//	preTestChecks(t)
//	skipNoNsxtConfiguration(t)
//
//	client := createTemporaryVCDConnection()
//	if client.Client.APIVCDMaxVersionIs("< 35.2") {
//		t.Skip(t.Name() + " requires at least API v35.2 (vCD 10.2.2)")
//	}
//
//	// String map to fill the template
//	var params = StringMap{
//		"Org":         testConfig.VCD.Org,
//		"NsxtVdc":     testConfig.Nsxt.Vdc,
//		"EdgeGw":      testConfig.Nsxt.EdgeGateway,
//		"NetworkName": t.Name(),
//		"Tags":        "network nsxt",
//	}
//
//	configText := templateFill(testAccNsxtNatReflexive, params)
//	debugPrintf("#[DEBUG] CONFIGURATION for step 1: %s", configText)
//
//	//params["FuncName"] = t.Name() + "-step1"
//	//configText1 := templateFill(testAccNsxtSecurityGroupEmpty2, params)
//	//debugPrintf("#[DEBUG] CONFIGURATION for step 2: %s", configText1)
//
//	if vcdShortTest {
//		t.Skip(acceptanceTestsSkipped)
//		return
//	}
//
//	resource.Test(t, resource.TestCase{
//		ProviderFactories: testAccProviders,
//		PreCheck:          func() { testAccPreCheck(t) },
//		//CheckDestroy: resource.ComposeAggregateTestCheckFunc(
//		//	testAccCheckNsxtFirewallGroupDestroy(testConfig.Nsxt.Vdc, "test-security-group", types.FirewallGroupTypeSecurityGroup),
//		//	testAccCheckNsxtFirewallGroupDestroy(testConfig.Nsxt.Vdc, "test-security-group-changed", types.FirewallGroupTypeSecurityGroup),
//		//),
//		Steps: []resource.TestStep{
//			resource.TestStep{
//				Config: configText,
//				Check: resource.ComposeAggregateTestCheckFunc(
//					resource.TestCheckResourceAttrSet("vcd_nsxt_nat_rule.reflexive", "id"),
//					//resource.TestCheckResourceAttr("vcd_nsxt_security_group.group1", "name", "test-security-group"),
//					//resource.TestCheckResourceAttr("vcd_nsxt_security_group.group1", "description", "test-security-group-description"),
//					//resource.TestCheckNoResourceAttr("vcd_nsxt_security_group.group1", "member_org_network_ids"),
//					//resource.TestCheckNoResourceAttr("vcd_nsxt_security_group.group1", "member_vm_ids"),
//				),
//			},
//			//resource.TestStep{
//			//	Config: configText1,
//			//	Check: resource.ComposeAggregateTestCheckFunc(
//			//		resource.TestMatchResourceAttr("vcd_nsxt_security_group.group1", "id", regexp.MustCompile(`^urn:vcloud:firewallGroup:.*$`)),
//			//		resource.TestCheckResourceAttr("vcd_nsxt_security_group.group1", "name", "test-security-group-changed"),
//			//		resource.TestCheckResourceAttr("vcd_nsxt_security_group.group1", "description", ""),
//			//		resource.TestCheckNoResourceAttr("vcd_nsxt_security_group.group1", "member_org_network_ids"),
//			//		resource.TestCheckNoResourceAttr("vcd_nsxt_security_group.group1", "member_vm_ids"),
//			//	),
//			//},
//			//resource.TestStep{
//			//	ResourceName:      "vcd_nsxt_security_group.group1",
//			//	ImportState:       true,
//			//	ImportStateVerify: true,
//			//	ImportStateIdFunc: importStateIdNsxtEdgeGatewayObject(testConfig, testConfig.Nsxt.EdgeGateway, "test-security-group-changed"),
//			//},
//		},
//	})
//	postTestChecks(t)
//}
//
//const testAccNsxtNatReflexive = testAccNsxtSecurityGroupPrereqsEmpty + `
//resource "vcd_nsxt_nat_rule" "reflexive" {
//  org  = "{{.Org}}"
//  vdc  = "{{.NsxtVdc}}"
//
//  edge_gateway_id = data.vcd_nsxt_edgegateway.existing.id
//
//  name        = "test-reflexive-rule"
//  rule_type   = "REFLEXIVE"
//  description = "description"
//
//  # Using primary_ip from edge gateway
//  external_addresses         = tolist(data.vcd_nsxt_edgegateway.existing.subnet)[0].primary_ip
//  internal_addresses         = "11.11.11.2"
//  //snat_destination_addresses = "11.11.11.4"
//  //logging = true
//}
//`
