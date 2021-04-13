package vcd

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vmware/go-vcloud-director/v2/types/v56"
)

func resourceVcdSecurityGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVcdSecurityGroupCreate,
		ReadContext:   resourceVcdSecurityGroupRead,
		UpdateContext: resourceVcdSecurityGroupUpdate,
		DeleteContext: resourceVcdSecurityGroupDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceVcdSecurityGroupImport,
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
				Description: "Edge Gateway ID in which security group is located",
			},
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Security Group name",
			},
			"description": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Security Group description",
			},
			"member_org_network_ids": {
				Optional:    true,
				Type:        schema.TypeSet,
				Description: "Set of Org VDC network IDs attached to this security group",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"member_vm_ids": {
				Type:        schema.TypeSet,
				Computed:    true,
				Description: "Set of VM IDs",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

// resourceVcdSecurityGroupCreate
func resourceVcdSecurityGroupCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vcdClient := meta.(*VCDClient)
	vcdClient.lockParentEdgeGtw(d)
	defer vcdClient.unLockParentEdgeGtw(d)

	_, vdc, err := vcdClient.GetOrgAndVdcFromResource(d)
	if err != nil {
		return diag.Errorf(errorRetrievingOrgAndVdc, err)
	}

	fwGroup := getNsxtFirewallGroupType(d)

	createdFwGroup, err := vdc.CreateNsxtFirewallGroup(fwGroup)
	if err != nil {
		return diag.Errorf("error creating NSX-T Security Group '%s': %s", fwGroup.Name, err)
	}

	d.SetId(createdFwGroup.NsxtFirewallGroup.ID)

	return nil
}

// resourceVcdSecurityGroupUpdate
func resourceVcdSecurityGroupUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vcdClient := meta.(*VCDClient)
	vcdClient.lockParentEdgeGtw(d)
	defer vcdClient.unLockParentEdgeGtw(d)

	_, vdc, err := vcdClient.GetOrgAndVdcFromResource(d)
	if err != nil {
		return diag.Errorf(errorRetrievingOrgAndVdc, err)
	}

	fwGroup, err := vdc.GetNsxtFirewallGroupById(d.Id())
	if err != nil {
		return diag.Errorf("error getting NSX-T Security Group: %s", err)
	}

	updateFwGroup := getNsxtFirewallGroupType(d)
	// Inject existing ID for update
	updateFwGroup.ID = d.Id()

	_, err = fwGroup.Update(updateFwGroup)
	if err != nil {
		return diag.Errorf("error updating NSX-T Security Group '%s': %s", fwGroup.NsxtFirewallGroup.Name, err)
	}

	return nil
}

// resourceVcdSecurityGroupRead
func resourceVcdSecurityGroupRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vcdClient := meta.(*VCDClient)
	vcdClient.lockParentEdgeGtw(d)
	defer vcdClient.unLockParentEdgeGtw(d)

	_, vdc, err := vcdClient.GetOrgAndVdcFromResource(d)
	if err != nil {
		return diag.Errorf(errorRetrievingOrgAndVdc, err)
	}

	fwGroup, err := vdc.GetNsxtFirewallGroupById(d.Id())
	if err != nil {
		return diag.Errorf("error getting NSX-T Security Group: %s", err)
	}

	err = setNsxtFirewallGroupData(d, fwGroup.NsxtFirewallGroup)
	if err != nil {
		return diag.Errorf("error reading NSX-T Security Group: %s", err)
	}

	return nil
}

// resourceVcdSecurityGroupDelete
func resourceVcdSecurityGroupDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vcdClient := meta.(*VCDClient)
	vcdClient.lockParentEdgeGtw(d)
	defer vcdClient.unLockParentEdgeGtw(d)

	_, vdc, err := vcdClient.GetOrgAndVdcFromResource(d)
	if err != nil {
		return diag.Errorf(errorRetrievingOrgAndVdc, err)
	}

	fwGroup, err := vdc.GetNsxtFirewallGroupById(d.Id())
	if err != nil {
		return diag.Errorf("error getting NSX-T Security Group: %s", err)
	}

	err = fwGroup.Delete()
	if err != nil {
		return diag.Errorf("error deleting NSX-T Security Group: %s", err)
	}

	d.SetId("")

	return nil
}

// resourceVcdSecurityGroupImport
func resourceVcdSecurityGroupImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	resourceURI := strings.Split(d.Id(), ImportSeparator)
	if len(resourceURI) != 4 {
		return nil, fmt.Errorf("resource name must be specified as org-name.vdc-name.edge_gateway_name.security_group_name")
	}
	orgName, vdcName, edgeGatewayName, securityGroupName := resourceURI[0], resourceURI[1], resourceURI[2], resourceURI[3]

	vcdClient := meta.(*VCDClient)
	org, err := vcdClient.GetAdminOrg(orgName)
	if err != nil {
		return nil, fmt.Errorf("unable to find Org %s: %s", orgName, err)
	}
	vdc, err := org.GetVDCByName(vdcName, false)
	if err != nil {
		return nil, fmt.Errorf("unable to find VDC %s: %s", vdcName, err)
	}

	if !vdc.IsNsxt() {
		return nil, errors.New("security groups are only supported by NSX-T VDCs")
	}

	edgeGateway, err := vdc.GetNsxtEdgeGatewayByName(edgeGatewayName)
	if err != nil {
		return nil, fmt.Errorf("unable to find Edge Gateway '%s': %s", edgeGatewayName, err)
	}

	firewallGroup, err := edgeGateway.GetNsxtFirewallGroupByName(securityGroupName)
	if err != nil {
		return nil, fmt.Errorf("unable to find Security Group '%s': %s", edgeGatewayName, err)
	}

	if !firewallGroup.IsSecurityGroup() {
		return nil, fmt.Errorf("Firewall Group '%s' is not a Security Group, but '%s'",
			firewallGroup.NsxtFirewallGroup.Name, firewallGroup.NsxtFirewallGroup.Type)
	}

	_ = d.Set("org", orgName)
	_ = d.Set("vdc", vdcName)
	_ = d.Set("edge_gateway_id", edgeGateway.EdgeGateway.ID)
	d.SetId(firewallGroup.NsxtFirewallGroup.ID)

	return []*schema.ResourceData{d}, nil
}

func setNsxtFirewallGroupData(d *schema.ResourceData, fw *types.NsxtFirewallGroup) error {

	_ = d.Set("name", fw.Name)
	_ = d.Set("description", fw.Description)

	netIds := make([]string, len(fw.Members))
	for i := range fw.Members {
		netIds[i] = fw.Members[i].ID
	}

	// Convert `member_org_network_ids` to set
	memberNetIds := convertToTypeSet(netIds)
	memberNetSet := schema.NewSet(schema.HashSchema(&schema.Schema{Type: schema.TypeString}), memberNetIds)

	err := d.Set("member_org_network_ids", memberNetSet)
	if err != nil {
		return fmt.Errorf("error setting 'member_org_network_ids': %s", err)
	}

	return nil

}

func getNsxtFirewallGroupType(d *schema.ResourceData) *types.NsxtFirewallGroup {
	fwGroup := &types.NsxtFirewallGroup{
		// ID:          "",
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		// IpAddresses: []string{},
		// Members: []types.OpenApi Reference{},
		OwnerRef: &types.OpenApiReference{
			ID: d.Get("edge_gateway_id").(string),
		},
		// EdgeGatewayRef: &types.OpenApiReference{
		// 	ID: d.Get("edge_gateway_id").(string),
		// },
		Type: types.FirewallGroupTypeSecurityGroup,
	}

	// Expand member networks
	orgNetworkIds := convertSchemaSetToSliceOfStrings(d.Get("member_org_network_ids").(*schema.Set))
	memberReferences := make([]types.OpenApiReference, len(orgNetworkIds))
	for i := range orgNetworkIds {
		memberReferences[i].ID = orgNetworkIds[i]
	}

	if len(memberReferences) > 0 {
		fwGroup.Members = memberReferences
	}

	return fwGroup
}
