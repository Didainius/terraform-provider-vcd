package vcd

import (
	"context"
	"fmt"
	"strings"

	"github.com/vmware/go-vcloud-director/v2/govcd"
	"github.com/vmware/go-vcloud-director/v2/types/v56"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var appPortDefinition = &schema.Resource{
	Schema: map[string]*schema.Schema{
		"protocol": {
			Required:         true,
			Type:             schema.TypeString,
			ValidateFunc:     validation.StringInSlice([]string{"ICMPv4", "ICMPv6", "TCP", "UDP"}, true),
			DiffSuppressFunc: suppressCase,
		},
		"port": {
			Optional:    true,
			Type:        schema.TypeSet,
			Description: "Set of ports or ranges",
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
	},
}

func resourceVcdNsxtAppPortProfile() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVcdNsxtAppPortProfileCreate,
		ReadContext:   resourceVcdNsxtAppPortProfileRead,
		UpdateContext: resourceVcdNsxtAppPortProfileUpdate,
		DeleteContext: resourceVcdNsxtAppPortProfileDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceVcdNsxtAppPortProfileImport,
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
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Application Port Profile name",
			},
			"description": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Application Port Profile description",
			},
			"scope": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				Description:  "Scope - 'PROVIDER' or 'TENANT'",
				ValidateFunc: validation.StringInSlice([]string{"PROVIDER", "TENANT"}, false),
			},
			"nsxt_manager_id": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "ID of NSX-T manager. Only required for 'PROVIDER' scope",
			},
			"app_port": {
				Required: true,
				MinItems: 1,
				Type:     schema.TypeSet,
				Elem:     appPortDefinition,
			},
		},
	}
}

func resourceVcdNsxtAppPortProfileCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vcdClient := meta.(*VCDClient)
	// Leaving this commented to remember and check if locks are required at NSX-T manager or Org VDC
	//vcdClient.lockParentEdgeGtw(d)
	//defer vcdClient.unLockParentEdgeGtw(d)

	org, vdc, err := vcdClient.GetOrgAndVdcFromResource(d)
	if err != nil {
		return diag.Errorf(errorRetrievingOrgAndVdc, err)
	}

	appPortProfile := getNsxtAppPortProfileType(d, org, vdc)

	createdAppPortProfile, err := org.CreateNsxtAppPortProfile(appPortProfile)
	if err != nil {
		return diag.Errorf("error creating NSX-T Application Port Profile '%s': %s", appPortProfile.Name, err)
	}

	d.SetId(createdAppPortProfile.NsxtAppPortProfile.ID)

	return resourceVcdNsxtAppPortProfileRead(ctx, d, meta)
}

func resourceVcdNsxtAppPortProfileUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vcdClient := meta.(*VCDClient)

	org, vdc, err := vcdClient.GetOrgAndVdcFromResource(d)
	if err != nil {
		return diag.Errorf(errorRetrievingOrgAndVdc, err)
	}

	appPortProfile, err := org.GetNsxtAppPortProfileById(d.Id())
	if err != nil {
		return diag.Errorf("error getting NSX-T Security Group: %s", err)
	}

	updateappPortProfile := getNsxtAppPortProfileType(d, org, vdc)
	// Inject existing ID for update
	updateappPortProfile.ID = d.Id()

	_, err = appPortProfile.Update(updateappPortProfile)
	if err != nil {
		return diag.Errorf("error updating NSX-T Application Port Profile '%s': %s", updateappPortProfile.Name, err)
	}

	return resourceVcdNsxtAppPortProfileRead(ctx, d, meta)
}

func resourceVcdNsxtAppPortProfileRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vcdClient := meta.(*VCDClient)

	org, _, err := vcdClient.GetOrgAndVdcFromResource(d)
	if err != nil {
		return diag.Errorf(errorRetrievingOrgAndVdc, err)
	}

	appPortProfile, err := org.GetNsxtAppPortProfileById(d.Id())
	if err != nil {
		if govcd.ContainsNotFound(err) {
			d.SetId("")
			return nil
		}
		return diag.Errorf("error getting NSX-T Application Port Profile with ID '%s': %s", d.Id(), err)
	}

	err = setNsxtAppPortProfileData(d, appPortProfile.NsxtAppPortProfile)
	if err != nil {
		return diag.Errorf("error reading NSX-T Application Port Profile: %s", err)
	}

	return nil
}

func resourceVcdNsxtAppPortProfileDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vcdClient := meta.(*VCDClient)

	org, _, err := vcdClient.GetOrgAndVdcFromResource(d)
	if err != nil {
		return diag.Errorf(errorRetrievingOrgAndVdc, err)
	}

	appPortProfile, err := org.GetNsxtAppPortProfileById(d.Id())
	if err != nil {
		return diag.Errorf("error getting NSX-T Application Port Profile: %s", err)
	}

	err = appPortProfile.Delete()
	if err != nil {
		return diag.Errorf("error deleting NSX-T Application Port Profile: %s", err)
	}

	d.SetId("")

	return nil
}

func resourceVcdNsxtAppPortProfileImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {

	return []*schema.ResourceData{d}, nil
}

func getNsxtAppPortProfileType(d *schema.ResourceData, org *govcd.Org, vdc *govcd.Vdc) *types.NsxtAppPortProfile {
	appPortProfileConfig := &types.NsxtAppPortProfile{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Scope:       d.Get("scope").(string),
	}

	scope := d.Get("scope").(string)
	switch strings.ToUpper(scope) {
	case types.ApplicationPortProfileScopeProvider:
		appPortProfileConfig.Scope = scope
		nsxtManagerUrn := d.Get("nsxt_manager_id").(string)
		appPortProfileConfig.ContextEntityId = nsxtManagerUrn
	case types.ApplicationPortProfileScopeTenant:
		appPortProfileConfig.Scope = scope
		appPortProfileConfig.OrgRef = &types.OpenApiReference{ID: org.Org.ID}
		appPortProfileConfig.ContextEntityId = vdc.Vdc.ID
	}

	appPortSet := d.Get("app_port").(*schema.Set)
	if appPortSet != nil {
		appPortSlice := appPortSet.List()
		applicationPorts := make([]types.NsxtAppPortProfilePort, len(appPortSlice))
		for index, singlePort := range appPortSlice {
			appPortMap := singlePort.(map[string]interface{})
			onePortDef := types.NsxtAppPortProfilePort{
				Protocol:         appPortMap["protocol"].(string),
				DestinationPorts: convertSchemaSetToSliceOfStrings(appPortMap["port"].(*schema.Set)),
			}
			applicationPorts[index] = onePortDef
		}
		appPortProfileConfig.ApplicationPorts = applicationPorts
	}

	return appPortProfileConfig
}

func setNsxtAppPortProfileData(d *schema.ResourceData, appPortProfile *types.NsxtAppPortProfile) error {
	_ = d.Set("name", appPortProfile.Name)
	_ = d.Set("description", appPortProfile.Description)
	_ = d.Set("scope", appPortProfile.Scope)

	if appPortProfile.ApplicationPorts != nil && len(appPortProfile.ApplicationPorts) > 0 {

		resultSet := make([]interface{}, len(appPortProfile.ApplicationPorts))

		for index, value := range appPortProfile.ApplicationPorts {
			appPortMap := make(map[string]interface{})
			appPortMap["protocol"] = value.Protocol

			destinationPortInterface := convertToTypeSet(value.DestinationPorts)
			desinationPortSet := schema.NewSet(schema.HashSchema(&schema.Schema{Type: schema.TypeString}), destinationPortInterface)
			appPortMap["port"] = desinationPortSet

			resultSet[index] = appPortMap

		}

		appPortSet := schema.NewSet(schema.HashResource(appPortDefinition), resultSet)
		err := d.Set("app_port", appPortSet)
		if err != nil {
			return fmt.Errorf("error setting Application Port Profile: %s", err)
		}
	}

	return nil
}
