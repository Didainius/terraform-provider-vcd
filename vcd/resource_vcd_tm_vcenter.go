package vcd

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vmware/go-vcloud-director/v2/govcd"
	"github.com/vmware/go-vcloud-director/v2/types/v56"
)

func resourceVcdTmVcenter() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVcdTmVcenterCreate,
		// ReadContext:   resourceVcdTmVcenterRead,
		// UpdateContext: resourceVcdTmVcenterUpdate,
		// DeleteContext: resourceVcdTmVcenterDelete,
		// Importer: &schema.ResourceImporter{
		// 	StateContext: resourceVcdTmVcenterImport,
		// },

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "",
			},
			"url": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "",
			},
			"username": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "",
			},
			"password": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				Description: "",
			},
			"is_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "",
			},

			"has_proxy": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "",
			},
			"is_connected": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "",
			},
			"mode": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "",
			},
			"listener_state": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "",
			},
			"cluster_health_status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "",
			},
			"version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "",
			},
			"uuid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "",
			},
		},
	}
}

const labelVirtualCenter = "Virtual Center"

func getTmVcenterType(d *schema.ResourceData) (*types.VSphereVirtualCenter, error) {
	t := &types.VSphereVirtualCenter{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Url:         d.Get("url").(string),
		Username:    d.Get("username").(string),
		Password:    d.Get("password").(string),
		IsEnabled:   d.Get("is_enabled").(bool),
	}

	return t, nil
}

func setTmVcenterData(d *schema.ResourceData, v *govcd.VCenter) error {
	if v == nil || v.VSphereVCenter == nil {
		return fmt.Errorf("nil object")
	}

	dSet(d, "name", v.VSphereVCenter.Name)
	dSet(d, "description", v.VSphereVCenter.Description)
	dSet(d, "url", v.VSphereVCenter.Url)
	dSet(d, "username", v.VSphereVCenter.Username)
	// dSet(d, "password", v.VSphereVCenter.Password) // password is never returned,
	dSet(d, "is_enabled", v.VSphereVCenter.IsEnabled)

	dSet(d, "has_proxy", v.VSphereVCenter.HasProxy)
	dSet(d, "is_connected", v.VSphereVCenter.IsConnected)
	dSet(d, "mode", v.VSphereVCenter.Mode)
	dSet(d, "listener_state", v.VSphereVCenter.ListenerState)
	dSet(d, "cluster_health_status", v.VSphereVCenter.ClusterHealthStatus)
	dSet(d, "version", v.VSphereVCenter.VcVersion)
	dSet(d, "uuid", v.VSphereVCenter.Uuid)

	d.SetId(v.VSphereVCenter.VcId)

	return nil
}

func resourceVcdTmVcenterCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vcdClient := meta.(*VCDClient)
	c := crudConfig[*govcd.VCenter, types.VSphereVirtualCenter]{
		entityLabel:    labelVirtualCenter,
		getTypeFunc:    getTmVcenterType,
		createFunc:     vcdClient.CreateVcenter,
		stateStoreFunc: setTmVcenterData,
		readFunc:       resourceVcdTmVcenterRead,
	}

	// return create(ctx, d, meta, labelVirtualCenter, getTmVcenterType, vcdClient.CreateVcenter, setTmVcenterData, resourceVcdTmVcenterRead)
	return create2(ctx, d, meta, c)
	// return nil
}

func resourceVcdTmVcenterUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vcdClient := meta.(*VCDClient)
	return update(ctx, d, meta, labelVirtualCenter, getTmVcenterType, vcdClient.GetVCenterById, resourceVcdTmVcenterRead)
}

func resourceVcdTmVcenterRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vcdClient := meta.(*VCDClient)
	return read(ctx, d, meta, labelVirtualCenter, vcdClient.GetVCenterById, setTmVcenterData)
}

func resourceVcdTmVcenterDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vcdClient := meta.(*VCDClient)
	return deleteRes(ctx, d, meta, labelVirtualCenter, vcdClient.GetVCenterById)
}

// func resourceVcdTmVcenterImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
// 	vcdClient := meta.(*VCDClient)

// 	v, err := vcdClient.GetVCenterByName(d.Id())
// 	if err != nil {
// 		return nil, fmt.Errorf("error retrieving vCenter by name: %s", err)
// 	}

// 	d.SetId(v.VSphereVCenter.VcId)
// 	return []*schema.ResourceData{d}, nil
// }
