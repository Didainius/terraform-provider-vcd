// +build network nsxt ALL functional

package vcd

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

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

	configText1 := templateFill(testAccNsxtNatDnat, params)
	debugPrintf("#[DEBUG] CONFIGURATION for step 1: %s", configText1)

	params["FuncName"] = t.Name() + "-step2"
	configText2 := templateFill(testAccNsxtNatDnatStep2, params)
	debugPrintf("#[DEBUG] CONFIGURATION for step 2: %s", configText2)

	if vcdShortTest {
		t.Skip(acceptanceTestsSkipped)
		return
	}

	natRuleId := &testCachedFieldValue{}

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviders,
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy: resource.ComposeAggregateTestCheckFunc(
			testAccCheckNsxtNatRuleDestroy("test-dnat-rule"),
			testAccCheckNsxtNatRuleDestroy("test-dnat-rule-updated"),
		),
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: configText1,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("vcd_nsxt_nat_rule.dnat", "id"),
					resource.TestMatchResourceAttr("vcd_nsxt_nat_rule.dnat", "edge_gateway_id", regexp.MustCompile(`^urn:vcloud:gateway:`)),
					resource.TestCheckResourceAttr("vcd_nsxt_nat_rule.dnat", "name", "test-dnat-rule"),
					resource.TestCheckResourceAttr("vcd_nsxt_nat_rule.dnat", "rule_type", "DNAT"),
					resource.TestCheckResourceAttr("vcd_nsxt_nat_rule.dnat", "description", "description"),
					resource.TestCheckResourceAttrSet("vcd_nsxt_nat_rule.dnat", "external_addresses"),
					resource.TestCheckResourceAttr("vcd_nsxt_nat_rule.dnat", "internal_addresses", "11.11.11.2"),
					resource.TestCheckResourceAttr("vcd_nsxt_nat_rule.dnat", "logging", "true"),
					resource.TestCheckNoResourceAttr("vcd_nsxt_nat_rule.dnat", "app_port_profile_id"),
					resource.TestCheckResourceAttr("vcd_nsxt_nat_rule.dnat", "snat_destination_addresses", ""),
					resource.TestCheckResourceAttr("vcd_nsxt_nat_rule.dnat", "enabled", "true"),
					natRuleId.cacheTestResourceFieldValue("vcd_nsxt_nat_rule.dnat", "id"),
				),
			},
			resource.TestStep{
				Config: configText2,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("vcd_nsxt_nat_rule.dnat", "id"),
					resource.TestMatchResourceAttr("vcd_nsxt_nat_rule.dnat", "edge_gateway_id", regexp.MustCompile(`^urn:vcloud:gateway:`)),
					resource.TestCheckResourceAttr("vcd_nsxt_nat_rule.dnat", "name", "test-dnat-rule-updated"),
					resource.TestCheckResourceAttr("vcd_nsxt_nat_rule.dnat", "rule_type", "DNAT"),
					resource.TestCheckResourceAttr("vcd_nsxt_nat_rule.dnat", "description", "updated-description"),
					resource.TestCheckResourceAttrSet("vcd_nsxt_nat_rule.dnat", "external_addresses"),
					resource.TestCheckResourceAttr("vcd_nsxt_nat_rule.dnat", "internal_addresses", "11.11.11.0/32"),
					resource.TestCheckResourceAttr("vcd_nsxt_nat_rule.dnat", "logging", "false"),
					resource.TestCheckNoResourceAttr("vcd_nsxt_nat_rule.dnat", "app_port_profile_id"),
					resource.TestCheckResourceAttr("vcd_nsxt_nat_rule.dnat", "dnat_external_port", "8888"),
					resource.TestCheckResourceAttr("vcd_nsxt_nat_rule.dnat", "snat_destination_addresses", ""),
					resource.TestCheckResourceAttr("vcd_nsxt_nat_rule.dnat", "enabled", "false"),
				),
			},
			// Try to import by Name
			resource.TestStep{
				ResourceName:      "vcd_nsxt_nat_rule.dnat",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: importStateIdNsxtEdgeGatewayObject(testConfig, testConfig.Nsxt.EdgeGateway, "test-dnat-rule-updated"),
			},
			// Try to import by rule UUID
			resource.TestStep{
				ResourceName:      "vcd_nsxt_nat_rule.dnat",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: importStateIdNsxtEdgeGatewayObjectLazyName(testConfig, testConfig.Nsxt.EdgeGateway,
					func() string { return natRuleId.fieldValue }),
			},
		},
	})
	postTestChecks(t)
}

const testAccNsxtNatDnat = testAccNsxtSecurityGroupPrereqsEmpty + `
resource "vcd_nsxt_nat_rule" "dnat" {
  org = "{{.Org}}"
  vdc = "{{.NsxtVdc}}"

  edge_gateway_id = data.vcd_nsxt_edgegateway.existing.id

  name        = "test-dnat-rule"
  rule_type   = "DNAT"
  description = "description"

  # Using primary_ip from edge gateway
  external_addresses = tolist(data.vcd_nsxt_edgegateway.existing.subnet)[0].primary_ip
  internal_addresses = "11.11.11.2"
  logging            = true
}
`

const testAccNsxtNatDnatStep2 = testAccNsxtSecurityGroupPrereqsEmpty + `
resource "vcd_nsxt_nat_rule" "dnat" {
  org  = "{{.Org}}"
  vdc  = "{{.NsxtVdc}}"

  edge_gateway_id = data.vcd_nsxt_edgegateway.existing.id

  name        = "test-dnat-rule-updated"
  rule_type  = "DNAT"
  description = "updated-description"
  
  # Using primary_ip from edge gateway
  external_addresses = tolist(data.vcd_nsxt_edgegateway.existing.subnet)[0].primary_ip
  internal_addresses = "11.11.11.0/32"
  dnat_external_port = 8888
  
  logging = false
  enabled = false
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

	if vcdShortTest {
		t.Skip(acceptanceTestsSkipped)
		return
	}

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviders,
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy:      testAccCheckNsxtNatRuleDestroy("test-no-dnat-rule"),
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: configText,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("vcd_nsxt_nat_rule.no-dnat", "id"),
					resource.TestCheckResourceAttr("vcd_nsxt_nat_rule.no-dnat", "name", "test-no-dnat-rule"),
					resource.TestCheckResourceAttr("vcd_nsxt_nat_rule.no-dnat", "rule_type", "NO_DNAT"),
					resource.TestCheckResourceAttr("vcd_nsxt_nat_rule.no-dnat", "description", ""),
					resource.TestCheckResourceAttrSet("vcd_nsxt_nat_rule.no-dnat", "external_addresses"),
					resource.TestCheckResourceAttr("vcd_nsxt_nat_rule.no-dnat", "internal_addresses", ""),
					resource.TestCheckResourceAttr("vcd_nsxt_nat_rule.no-dnat", "logging", "false"),
					resource.TestCheckResourceAttr("vcd_nsxt_nat_rule.no-dnat", "dnat_external_port", "7777"),
				),
			},
			resource.TestStep{
				ResourceName:      "vcd_nsxt_nat_rule.no-dnat",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: importStateIdNsxtEdgeGatewayObject(testConfig, testConfig.Nsxt.EdgeGateway, "test-no-dnat-rule"),
			},
		},
	})
	postTestChecks(t)
}

const testAccNsxtNatNoDnat = testAccNsxtSecurityGroupPrereqsEmpty + `
resource "vcd_nsxt_nat_rule" "no-dnat" {
  org  = "{{.Org}}"
  vdc  = "{{.NsxtVdc}}"

  edge_gateway_id = data.vcd_nsxt_edgegateway.existing.id

  name      = "test-no-dnat-rule"
  rule_type = "NO_DNAT"

  
  # Using primary_ip from edge gateway
  external_addresses = tolist(data.vcd_nsxt_edgegateway.existing.subnet)[0].primary_ip
  dnat_external_port = 7777

  # app_port_profile_id = 
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
  snat_destination_addresses = "8.8.8.8"
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

// TestAccVcdNsxtNatRuleFirewallMatchPriority explicitly tests support for two new fields introduced in API 35.2 (VCD 10.2.2)
// firewall_match and priority. For 10.2.2 versions this should work, while for lower versions it should return an error.
func TestAccVcdNsxtNatRuleFirewallMatchPriority(t *testing.T) {
	preTestChecks(t)
	skipNoNsxtConfiguration(t)

	// expectError must stay nil for versions > 10.2.2, because we expect it to work. For lower versions - it must have
	// match the runtime validation error
	var expectError *regexp.Regexp
	client := createTemporaryVCDConnection()
	if client.Client.APIVCDMaxVersionIs("< 35.2") {
		expectError = regexp.MustCompile(`firewall_match and priority fields can only be set for VCD 10.2.2+`)
	}

	// String map to fill the template
	var params = StringMap{
		"Org":           testConfig.VCD.Org,
		"NsxtVdc":       testConfig.Nsxt.Vdc,
		"EdgeGw":        testConfig.Nsxt.EdgeGateway,
		"NetworkName":   t.Name(),
		"Tags":          "network nsxt",
		"FirewallMatch": "MATCH_INTERNAL_ADDRESS",
		"Priority":      "10",
	}

	configText1 := templateFill(testAccNsxtNatFirewallMatchPriority, params)
	debugPrintf("#[DEBUG] CONFIGURATION for step 1: %s", configText1)

	params["FuncName"] = t.Name() + "-step2"
	params["FirewallMatch"] = "MATCH_EXTERNAL_ADDRESS"
	params["Priority"] = "30"
	configText2 := templateFill(testAccNsxtNatFirewallMatchPriority, params)
	debugPrintf("#[DEBUG] CONFIGURATION for step 2: %s", configText2)

	params["FuncName"] = t.Name() + "-step3"
	params["FirewallMatch"] = "BYPASS"
	params["Priority"] = "0"
	configText3 := templateFill(testAccNsxtNatFirewallMatchPriority, params)
	debugPrintf("#[DEBUG] CONFIGURATION for step 3: %s", configText3)

	if vcdShortTest {
		t.Skip(acceptanceTestsSkipped)
		return
	}

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviders,
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy:      testAccCheckNsxtNatRuleDestroy("test-dnat-rule-match-and-priority"),
		Steps: []resource.TestStep{
			resource.TestStep{
				Config:      configText1,
				ExpectError: expectError,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("vcd_nsxt_nat_rule.dnat-match", "id"),
					resource.TestCheckResourceAttr("vcd_nsxt_nat_rule.dnat-match", "name", "test-dnat-rule-match-and-priority"),
					resource.TestCheckResourceAttr("vcd_nsxt_nat_rule.dnat-match", "firewall_match", "MATCH_INTERNAL_ADDRESS"),
					resource.TestCheckResourceAttr("vcd_nsxt_nat_rule.dnat-match", "priority", "10"),
				),
			},
			resource.TestStep{
				Config:      configText2,
				ExpectError: expectError,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("vcd_nsxt_nat_rule.dnat-match", "id"),
					resource.TestCheckResourceAttr("vcd_nsxt_nat_rule.dnat-match", "name", "test-dnat-rule-match-and-priority"),
					resource.TestCheckResourceAttr("vcd_nsxt_nat_rule.dnat-match", "firewall_match", "MATCH_EXTERNAL_ADDRESS"),
					resource.TestCheckResourceAttr("vcd_nsxt_nat_rule.dnat-match", "priority", "30"),
				),
			},
			resource.TestStep{
				Config:      configText3,
				ExpectError: expectError,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("vcd_nsxt_nat_rule.dnat-match", "id"),
					resource.TestCheckResourceAttr("vcd_nsxt_nat_rule.dnat-match", "name", "test-dnat-rule-match-and-priority"),
					resource.TestCheckResourceAttr("vcd_nsxt_nat_rule.dnat-match", "firewall_match", "BYPASS"),
					resource.TestCheckResourceAttr("vcd_nsxt_nat_rule.dnat-match", "priority", "0"),
				),
			},
			resource.TestStep{
				ResourceName:      "vcd_nsxt_nat_rule.dnat-match",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: importStateIdNsxtEdgeGatewayObject(testConfig, testConfig.Nsxt.EdgeGateway, "test-dnat-rule-match-and-priority"),
			},
		},
	})
	postTestChecks(t)
}

const testAccNsxtNatFirewallMatchPriority = testAccNsxtSecurityGroupPrereqsEmpty + `
resource "vcd_nsxt_nat_rule" "dnat-match" {
  org = "{{.Org}}"
  vdc = "{{.NsxtVdc}}"

  edge_gateway_id = data.vcd_nsxt_edgegateway.existing.id

  name        = "test-dnat-rule-match-and-priority"
  rule_type   = "DNAT"
  description = "description"

  # Using primary_ip from edge gateway
  external_addresses = tolist(data.vcd_nsxt_edgegateway.existing.subnet)[0].primary_ip
  internal_addresses = "11.11.11.2"
  logging            = true

  firewall_match = "{{.FirewallMatch}}"
  priority       = "{{.Priority}}"
}
`

func testAccCheckNsxtNatRuleDestroy(natRuleIdentifier string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := testAccProvider.Meta().(*VCDClient)
		egw, err := conn.GetNsxtEdgeGateway(testConfig.VCD.Org, testConfig.Nsxt.Vdc, testConfig.Nsxt.EdgeGateway)
		if err != nil {
			return fmt.Errorf(errorUnableToFindEdgeGateway, testConfig.Nsxt.EdgeGateway)
		}

		_, errByName := egw.GetNatRuleByName(natRuleIdentifier)
		_, errById := egw.GetNatRuleById(natRuleIdentifier)

		if errByName == nil {
			return fmt.Errorf("got no errors for NSX-T NAT rule lookup by Name")
		}

		if errById == nil {
			return fmt.Errorf("got no errors for NSX-T NAT rule lookup by ID")
		}

		return nil
	}
}
