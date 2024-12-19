package vcd

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vmware/go-vcloud-director/v3/govcd"
	"github.com/vmware/go-vcloud-director/v3/types/v56"
)

const labelTmOrgNetworkingSettings = "Org Networking Settings"

func resourceVcdTmOrgNetworkingSettings() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVcdTmOrgNetworkingSettingsCreate,
		ReadContext:   resourceVcdTmOrgNetworkingSettingsRead,
		UpdateContext: resourceVcdTmOrgNetworkingSettingsUpdate,
		DeleteContext: resourceVcdTmOrgNetworkingSettingsDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceVcdTmOrgNetworkingSettingsImport,
		},

		Schema: map[string]*schema.Schema{
			"org_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: fmt.Sprintf("Org ID for %s", labelTmOrgNetworkingSettings),
			},
			"log_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Log name for Org",
			},
		},
	}
}

func resourceVcdTmOrgNetworkingSettingsCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vcdClient := meta.(*VCDClient)

	createQosConfigInEdgeCluster := func(config *types.TmOrgNetworkingSettings) (*govcd.TmOrgNetworkingSettings, error) {
		ec, err := vcdClient.GetTmOrgNetworkingSettingsByOrgId(d.Get("org_id").(string))
		if err != nil {
			return nil, fmt.Errorf("error looking up %s by ID: %s", labelTmOrgNetworkingSettings, err)
		}
		return ec.Update(config)
	}

	c := crudConfig[*govcd.TmOrgNetworkingSettings, types.TmOrgNetworkingSettings]{
		entityLabel:      labelTmOrgNetworkingSettings,
		getTypeFunc:      getTmOrgNetworkingSettingsType,
		stateStoreFunc:   setTmOrgNetworkingSettingsData,
		createFunc:       createQosConfigInEdgeCluster,
		resourceReadFunc: resourceVcdTmOrgNetworkingSettingsRead,
	}
	return createResource(ctx, d, meta, c)
}

func resourceVcdTmOrgNetworkingSettingsUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vcdClient := meta.(*VCDClient)
	c := crudConfig[*govcd.TmOrgNetworkingSettings, types.TmOrgNetworkingSettings]{
		entityLabel:      labelTmOrgNetworkingSettings,
		getTypeFunc:      getTmOrgNetworkingSettingsType,
		getEntityFunc:    vcdClient.GetTmOrgNetworkingSettingsByOrgId,
		resourceReadFunc: resourceVcdTmOrgNetworkingSettingsRead,
	}

	return updateResource(ctx, d, meta, c)
}

func resourceVcdTmOrgNetworkingSettingsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vcdClient := meta.(*VCDClient)
	c := crudConfig[*govcd.TmOrgNetworkingSettings, types.TmOrgNetworkingSettings]{
		entityLabel:    labelTmOrgNetworkingSettings,
		getEntityFunc:  vcdClient.GetTmOrgNetworkingSettingsByOrgId,
		stateStoreFunc: setTmOrgNetworkingSettingsData,
	}
	return readResource(ctx, d, meta, c)
}

func resourceVcdTmOrgNetworkingSettingsDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vcdClient := meta.(*VCDClient)

	c := crudConfig[*govcd.TmOrgNetworkingSettings, types.TmOrgNetworkingSettings]{
		entityLabel:   labelTmOrgNetworkingSettings,
		getEntityFunc: vcdClient.GetTmOrgNetworkingSettingsByOrgId,
	}

	return deleteResource(ctx, d, meta, c)
}

func resourceVcdTmOrgNetworkingSettingsImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	_ = meta.(*VCDClient)

	d.SetId("???")
	return []*schema.ResourceData{d}, nil
}

func getTmOrgNetworkingSettingsType(vcdClient *VCDClient, d *schema.ResourceData) (*types.TmOrgNetworkingSettings, error) {
	t := &types.TmOrgNetworkingSettings{
		OrgNameForLogs: d.Get("log_name").(string),
	}

	return t, nil
}

func setTmOrgNetworkingSettingsData(_ *VCDClient, d *schema.ResourceData, org *govcd.TmOrgNetworkingSettings) error {
	if org == nil || org.TmOrgNetworkingSettings == nil {
		return fmt.Errorf("nil value received for %s", labelTmOrgNetworkingSettings)
	}

	d.SetId(org.OrgId)
	dSet(d, "log_name", org.TmOrgNetworkingSettings.OrgNameForLogs)

	return nil
}
