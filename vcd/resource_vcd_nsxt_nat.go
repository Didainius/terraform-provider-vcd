package vcd

import (
	"context"
	"fmt"

	"github.com/vmware/go-vcloud-director/v2/types/v56"

	"github.com/vmware/go-vcloud-director/v2/govcd"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceVcdNsxtNat() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVcdNsxtNatCreate,
		ReadContext:   resourceVcdNsxtNatRead,
		UpdateContext: resourceVcdNsxtNatUpdate,
		DeleteContext: resourceVcdNsxtNatDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceVcdNsxtNatImport,
		},

		Schema: map[string]*schema.Schema{
			"org": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Description: "The name of organization to use, optional if defined at provider " +
					"level. Useful when connected as sysadmin working across different organizations",
			},
			"vdc": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "The name of VDC to use, optional if defined at provider level",
			},
			"edge_gateway_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Edge gateway name in which NAT Rule is located",
			},
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "",
			},
			"rule_type": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "",
			},
			"description": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "",
			},
			"external_addresses": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "",
			},
			"internal_addresses": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "",
			},
			"app_port_profile_id": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "",
			},
			"dnat_external_port": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "",
			},
			"snat_destination_addresses": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "",
			},
			"logging": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "",
			},
			"enabled": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "",
			},
			"firewall_match": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "VCD 10.2.2+ Determines how the firewall matches the address during NATing if firewall stage is not skipped. Below are valid values.",
			},
			"priority": &schema.Schema{
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "VCD 10.2.2+ If an address has multiple NAT rules, the rule with the highest priority is applied. A lower value means a higher precedence for this rule.",
			},
		},
	}
}

// resourceVcdNsxtNatCreate
func resourceVcdNsxtNatCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vcdClient := meta.(*VCDClient)
	vcdClient.lockParentEdgeGtw(d)
	defer vcdClient.unLockParentEdgeGtw(d)

	orgName := d.Get("org").(string)
	vdcName := d.Get("vdc").(string)
	edgeGatewayId := d.Get("edge_gateway_id").(string)

	nsxtEdge, err := vcdClient.GetNsxtEdgeGatewayById(orgName, vdcName, edgeGatewayId)
	if err != nil {
		return diag.Errorf("error retrieving Edge Gateway: %s", err)
	}

	nsxtNatRule, err := getNsxtNatType(d, vcdClient)
	if err != nil {
		return diag.Errorf("error getting NSX-T NAT rule type: %s", err)
	}

	rule, err := nsxtEdge.CreateNatRule(nsxtNatRule)
	if err != nil {

		return diag.Errorf("error creating NSX-T NAT rule: %s", err)
	}

	d.SetId(rule.NsxtNatRule.ID)

	return resourceVcdNsxtNatRead(ctx, d, meta)
}

// resourceVcdNsxtNatUpdate
func resourceVcdNsxtNatUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vcdClient := meta.(*VCDClient)
	vcdClient.lockParentEdgeGtw(d)
	defer vcdClient.unLockParentEdgeGtw(d)

	orgName := d.Get("org").(string)
	vdcName := d.Get("vdc").(string)
	edgeGatewayId := d.Get("edge_gateway_id").(string)

	nsxtEdge, err := vcdClient.GetNsxtEdgeGatewayById(orgName, vdcName, edgeGatewayId)
	if err != nil {
		return diag.Errorf("error retrieving Edge Gateway: %s", err)
	}

	nsxtNatRule, err := getNsxtNatType(d, vcdClient)
	if err != nil {
		return diag.Errorf("error getting NSX-T NAT rule type: %s", err)
	}

	existingRule, err := nsxtEdge.GetNatRuleById(d.Id())
	if err != nil {
		return diag.Errorf("unable to find NSX-T NAT rule: %s", err)
	}

	// Inject ID for update
	nsxtNatRule.ID = existingRule.NsxtNatRule.ID
	_, err = existingRule.Update(nsxtNatRule)
	if err != nil {
		return diag.Errorf("error updating NSX-T NAT rule: %s", err)
	}

	return resourceVcdNsxtNatRead(ctx, d, meta)
}

// resourceVcdNsxtNatRead
func resourceVcdNsxtNatRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vcdClient := meta.(*VCDClient)

	orgName := d.Get("org").(string)
	vdcName := d.Get("vdc").(string)
	edgeGatewayId := d.Get("edge_gateway_id").(string)

	nsxtEdge, err := vcdClient.GetNsxtEdgeGatewayById(orgName, vdcName, edgeGatewayId)
	if err != nil {
		return diag.Errorf("error retrieving Edge Gateway: %s", err)
	}

	existingRule, err := nsxtEdge.GetNatRuleById(d.Id())
	if err != nil {
		if govcd.ContainsNotFound(err) {
			d.SetId("")
		}
		return diag.Errorf("unable to find NSX-T NAT rule: %s", err)
	}

	err = setNsxtNatRuleData(existingRule.NsxtNatRule, d, vcdClient)
	if err != nil {
		return diag.Errorf("error storing NSX-T NAT rule in statefile: %s", err)
	}

	return nil
}

// resourceVcdNsxtNatDelete
func resourceVcdNsxtNatDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	vcdClient := meta.(*VCDClient)
	vcdClient.lockParentEdgeGtw(d)
	defer vcdClient.unLockParentEdgeGtw(d)

	orgName := d.Get("org").(string)
	vdcName := d.Get("vdc").(string)
	edgeGatewayId := d.Get("edge_gateway_id").(string)

	nsxtEdge, err := vcdClient.GetNsxtEdgeGatewayById(orgName, vdcName, edgeGatewayId)
	if err != nil {
		return diag.Errorf("error retrieving Edge Gateway: %s", err)
	}

	rule, err := nsxtEdge.GetNatRuleById(d.Id())
	if err != nil {
		return diag.Errorf("error finding NSX-T NAT Rule: %s", err)
	}

	err = rule.Delete()
	if err != nil {
		return diag.Errorf("error deleting NSX-T NAT rule: %s", err)
	}

	d.SetId("")

	return nil
}

// resourceVcdNsxtNatImport
func resourceVcdNsxtNatImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	return []*schema.ResourceData{d}, nil
}

func getNsxtNatType(d *schema.ResourceData, client *VCDClient) (*types.NsxtNatRule, error) {

	firewallMatch, firewallMatchOk := d.GetOk("firewall_match")
	priority, priorityOk := d.GetOk("priority")

	// Only supported in VCD 10.2.2+ (API V35.2) and throw immediate error if used with older versions as API error is
	// opaque
	if (firewallMatchOk || priorityOk) && client.Client.APIVCDMaxVersionIs("< 35.2") {
		return nil, fmt.Errorf("firewall_match and priority fields can only be set for VCD 10.2.2+")
	}

	nsxtNatRule := &types.NsxtNatRule{
		Name:                     d.Get("name").(string),
		Description:              d.Get("description").(string),
		Enabled:                  d.Get("enabled").(bool),
		RuleType:                 d.Get("rule_type").(string),
		ExternalAddresses:        d.Get("external_addresses").(string),
		InternalAddresses:        d.Get("internal_addresses").(string),
		DnatExternalPort:         d.Get("dnat_external_port").(string),
		SnatDestinationAddresses: d.Get("snat_destination_addresses").(string),
		Logging:                  d.Get("logging").(bool),
	}

	if appPortProf, ok := d.GetOk("app_port_profile_id"); ok {
		nsxtNatRule.ApplicationPortProfile = &types.OpenApiReference{ID: appPortProf.(string)}
	}

	if firewallMatchOk {
		nsxtNatRule.FirewallMatch = firewallMatch.(string)
	}

	if priorityOk {
		nsxtNatRule.Priority = takeIntPointer(priority.(int))
	}
	return nsxtNatRule, nil
}

func setNsxtNatRuleData(rule *types.NsxtNatRule, d *schema.ResourceData, client *VCDClient) error {
	_ = d.Set("name", rule.Name)
	_ = d.Set("description", rule.Description)
	_ = d.Set("external_addresses", rule.ExternalAddresses)
	_ = d.Set("internal_addresses", rule.InternalAddresses)
	_ = d.Set("dnat_external_port", rule.DnatExternalPort)
	_ = d.Set("snat_destination_addresses", rule.SnatDestinationAddresses)
	_ = d.Set("logging", rule.Logging)
	_ = d.Set("enabled", rule.Enabled)
	_ = d.Set("rule_type", rule.RuleType)

	if client.Client.APIVCDMaxVersionIs("< 35.2") {
		_ = d.Set("firewall_match", rule.FirewallMatch)
		_ = d.Set("priority", rule.Priority)
	}

	return nil
}
