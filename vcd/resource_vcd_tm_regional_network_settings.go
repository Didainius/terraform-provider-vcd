package vcd

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vmware/go-vcloud-director/v3/govcd"
	"github.com/vmware/go-vcloud-director/v3/types/v56"
)

const labelTmRegionalNetworkingSettings = "Tm Regional Networking Settings"

func resourceVcdTmRegionalNetworkSettings() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVcdTmRegionalNetworkSettingsCreate,
		ReadContext:   resourceVcdTmRegionalNetworkSettingsRead,
		UpdateContext: resourceVcdTmRegionalNetworkSettingsUpdate,
		DeleteContext: resourceVcdTmRegionalNetworkSettingsDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceVcdTmRegionalNetworkSettingsImport,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: fmt.Sprintf("Org ID for %s", labelTmOrgNetworkingSettings),
			},
			"org_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: fmt.Sprintf("Org ID for %s", labelTmOrgNetworkingSettings),
			},
			"region_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: fmt.Sprintf("Org ID for %s", labelTmOrgNetworkingSettings),
			},
			"provider_gateway_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: fmt.Sprintf("Org ID for %s", labelTmOrgNetworkingSettings),
			},
			"edge_cluster_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: fmt.Sprintf("Org ID for %s", labelTmOrgNetworkingSettings),
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: fmt.Sprintf("Status of %s", labelTmRegionalNetworkingSettings),
			},
		},
	}
}

func resourceVcdTmRegionalNetworkSettingsCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vcdClient := meta.(*VCDClient)
	c := crudConfig[*govcd.TmRegionalNetworkingSettings, types.TmRegionalNetworkingSettings]{
		entityLabel:      labelTmRegionalNetworkingSettings,
		getTypeFunc:      getTmRegionalNetworkingSettingsType,
		stateStoreFunc:   setTmRegionalNetworkingSettingsData,
		createFunc:       vcdClient.CreateTmRegionalNetworkingSettings,
		resourceReadFunc: resourceVcdTmRegionalNetworkSettingsRead,
	}
	return createResource(ctx, d, meta, c)
}

func resourceVcdTmRegionalNetworkSettingsUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vcdClient := meta.(*VCDClient)
	c := crudConfig[*govcd.TmRegionalNetworkingSettings, types.TmRegionalNetworkingSettings]{
		entityLabel:      labelTmRegionalNetworkingSettings,
		getTypeFunc:      getTmRegionalNetworkingSettingsType,
		getEntityFunc:    vcdClient.GetTmRegionalNetworkingSettingsById,
		resourceReadFunc: resourceVcdTmRegionalNetworkSettingsRead,
		// preUpdateHooks: []outerEntityHookInnerEntityType[*govcd.TmRegionalNetworkingSettings, *types.TmRegionalNetworkingSettings]{},
	}

	return updateResource(ctx, d, meta, c)
}

func resourceVcdTmRegionalNetworkSettingsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vcdClient := meta.(*VCDClient)
	c := crudConfig[*govcd.TmRegionalNetworkingSettings, types.TmRegionalNetworkingSettings]{
		entityLabel:    labelTmRegionalNetworkingSettings,
		getEntityFunc:  vcdClient.GetTmRegionalNetworkingSettingsById,
		stateStoreFunc: setTmRegionalNetworkingSettingsData,
	}
	return readResource(ctx, d, meta, c)
}

func resourceVcdTmRegionalNetworkSettingsDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vcdClient := meta.(*VCDClient)

	c := crudConfig[*govcd.TmRegionalNetworkingSettings, types.TmRegionalNetworkingSettings]{
		entityLabel:   labelTmRegionalNetworkingSettings,
		getEntityFunc: vcdClient.GetTmRegionalNetworkingSettingsById,
		// preDeleteHooks: []outerEntityHook[*govcd.TmRegionalNetworkingSettings]{},
	}

	return deleteResource(ctx, d, meta, c)
}

func resourceVcdTmRegionalNetworkSettingsImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	_ = meta.(*VCDClient)

	d.SetId("???")
	return []*schema.ResourceData{d}, nil
}

func getTmRegionalNetworkingSettingsType(vcdClient *VCDClient, d *schema.ResourceData) (*types.TmRegionalNetworkingSettings, error) {
	t := &types.TmRegionalNetworkingSettings{
		Name:               d.Get("name").(string),
		OrgRef:             types.OpenApiReference{ID: d.Get("org_id").(string)},
		RegionRef:          types.OpenApiReference{ID: d.Get("region_id").(string)},
		ProviderGatewayRef: types.OpenApiReference{ID: d.Get("provider_gateway_id").(string)},
		// ServiceEdgeClusterRef: ,
	}

	return t, nil
}

func setTmRegionalNetworkingSettingsData(_ *VCDClient, d *schema.ResourceData, org *govcd.TmRegionalNetworkingSettings) error {

	d.SetId(org.TmRegionalNetworkingSettings.ID)
	dSet(d, "name", org.TmRegionalNetworkingSettings.Name)
	dSet(d, "org_id", org.TmRegionalNetworkingSettings.OrgRef.ID)
	dSet(d, "region_id", org.TmRegionalNetworkingSettings.RegionRef.ID)
	dSet(d, "provider_gateway_id", org.TmRegionalNetworkingSettings.ProviderGatewayRef.ID)
	dSet(d, "status", org.TmRegionalNetworkingSettings.Status)
	// IMPLEMENT
	return nil
}
