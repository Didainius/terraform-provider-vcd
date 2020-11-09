package vcd

import (
	"fmt"
	"log"
	"strings"

	"github.com/vmware/go-vcloud-director/v2/govcd"
	"github.com/vmware/go-vcloud-director/v2/types/v56"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

var externalNetworkResource2 = &schema.Resource{
	Schema: map[string]*schema.Schema{
		"external_network_id": {
			Required:    true,
			ForceNew:    true,
			Type:        schema.TypeString,
			Description: "External network name",
		},

		"subnet": {
			Optional: true,
			Computed: true,
			ForceNew: true,
			Type:     schema.TypeSet,
			MinItems: 1,
			Elem:     subnetResource2,
		},
	},
}

var subnetResource2 = &schema.Resource{
	Schema: map[string]*schema.Schema{
		"gateway": {
			Required:    true,
			ForceNew:    true,
			Description: "Gateway address for a subnet",
			Type:        schema.TypeString,
		},
		"prefix_length": {
			Required:    true,
			ForceNew:    true,
			Description: "Netmask address for a subnet",
			Type:        schema.TypeString,
		},
		"ip_address": {
			Optional:    true,
			Type:        schema.TypeString,
			ForceNew:    true,
			Description: "IP address on the edge gateway - will be auto-assigned if not defined",
		},
		"use_for_default_route": {
			Optional:    true,
			Default:     false,
			ForceNew:    true,
			Type:        schema.TypeBool,
			Description: "Defines if this subnet should be used as default gateway for edge",
		},
		"suballocate_pool": {
			Optional:    true,
			Type:        schema.TypeSet,
			ForceNew:    true,
			Description: "Define zero or more blocks to sub-allocate pools on the edge gateway",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"start_address": {
						Required: true,
						Type:     schema.TypeString,
						ForceNew: true,
					},
					"end_address": {
						Required: true,
						Type:     schema.TypeString,
						ForceNew: true,
					},
				},
			},
		},
	},
}

func resourceVcdNsxtEdgeGateway() *schema.Resource {
	return &schema.Resource{
		Create: resourceVcdNsxtEdgeGatewayCreate,
		Read:   resourceVcdNsxtEdgeGatewayRead,
		Update: resourceVcdNsxtEdgeGatewayUpdate,
		Delete: resourceVcdNsxtEdgeGatewayDelete,
		Importer: &schema.ResourceImporter{
			State: resourceVcdNsxtEdgeGatewayImport,
		},

		Schema: map[string]*schema.Schema{
			"org": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Description: "The name of organization to use, optional if defined at provider " +
					"level. Useful when connected as sysadmin working across different organizations",
			},
			"vdc": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "The name of VDC to use, optional if defined at provider level",
			},
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Edge Gateway name",
			},
			"description": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Edge Gateway description",
			},
			"dedicate_external_network": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Dedicating the External Network will enable Route Advertisement for this Edge Gateway.",
			},
			"external_network": {
				Description: "One or more blocks with external network information to be attached to this gateway's interface",
				ForceNew:    true,
				Required:    true,
				Type:        schema.TypeSet,
				Elem:        externalNetworkResource2,
			},
			"edge_cluster_id": &schema.Schema{
				// TODO datasource - `vcd_nsxt_edge_cluster`
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Select specific NSX-T Edge Cluster. Will be inherited from external network if not specified",
			},
		},
	}
}

// resourceVcdNsxtEdgeGatewayCreate
func resourceVcdNsxtEdgeGatewayCreate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[TRACE] edge gateway creation initiated")

	vcdClient := meta.(*VCDClient)
	_, vdc, err := vcdClient.GetOrgAndVdcFromResource(d)
	if err != nil {
		return fmt.Errorf("error retrieving VDC: %s", err)
	}

	adminOrg, err := vcdClient.GetAdminOrgFromResource(d)
	if err != nil {
		return fmt.Errorf("error getting adminOrg: %s", err)
	}

	t, err := getNsxtEdgeGatewayType(d, vdc)
	if err != nil {
		return fmt.Errorf("could not create edge gateway type: %s", err)
	}

	_, err = adminOrg.CreateNsxtEdgeGateway(t)
	if err != nil {
		return fmt.Errorf("error creating edge gateway: %s", err)
	}

	return nil
}

// resourceVcdNsxtEdgeGatewayUpdate
func resourceVcdNsxtEdgeGatewayUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[TRACE] edge gateway update initiated")

	vcdClient := meta.(*VCDClient)
	_, vdc, err := vcdClient.GetOrgAndVdcFromResource(d)
	if err != nil {
		return fmt.Errorf("error retrieving VDC: %s", err)
	}

	adminOrg, err := vcdClient.GetAdminOrgFromResource(d)
	if err != nil {
		return fmt.Errorf("error getting adminOrg: %s", err)
	}

	edge, err := adminOrg.GetNsxtEdgeGatewayById(d.Id())
	if err != nil {
		return fmt.Errorf("could not retrieve edge gateway: %s", err)
	}

	edge.EdgeGateway, err = getNsxtEdgeGatewayType(d, vdc)
	if err != nil {
		return fmt.Errorf("error creating edge gateway type: %s", err)
	}

	_, err = edge.Update(edge.EdgeGateway)
	if err != nil {
		return fmt.Errorf("error updating edge gateway with ID '%s': %s", d.Id(), err)
	}

	return nil
}

// resourceVcdNsxtEdgeGatewayRead
func resourceVcdNsxtEdgeGatewayRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[TRACE] edge gateway read initiated")

	vcdClient := meta.(*VCDClient)
	// _, vdc, err := vcdClient.GetOrgAndVdcFromResource(d)
	// if err != nil {
	// 	return fmt.Errorf("error retrieving VDC: %s", err)
	// }

	adminOrg, err := vcdClient.GetAdminOrgFromResource(d)
	if err != nil {
		return fmt.Errorf("error getting adminOrg: %s", err)
	}

	edge, err := adminOrg.GetNsxtEdgeGatewayById(d.Id())
	if err != nil {
		return fmt.Errorf("could not retrieve edge gateway: %s", err)
	}

	err = setNsxtEdgeGatewayData(edge.EdgeGateway, d)
	if err != nil {
		return fmt.Errorf("error reading edge gateway data: %s", err)
	}
	return nil
}

// resourceVcdNsxtEdgeGatewayDelete
func resourceVcdNsxtEdgeGatewayDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[TRACE] edge gateway deletion initiated")

	vcdClient := meta.(*VCDClient)
	// _, vdc, err := vcdClient.GetOrgAndVdcFromResource(d)
	// if err != nil {
	// 	return fmt.Errorf("error retrieving VDC: %s", err)
	// }

	adminOrg, err := vcdClient.GetAdminOrgFromResource(d)
	if err != nil {
		return fmt.Errorf("error getting adminOrg: %s", err)
	}

	edge, err := adminOrg.GetNsxtEdgeGatewayById(d.Id())
	if err != nil {
		return fmt.Errorf("could not retrieve edge gateway: %s", err)
	}

	err = edge.Delete()
	if err != nil {
		return fmt.Errorf("error deleting edge gateway: %s", err)
	}

	return nil
}

// resourceVcdNsxtEdgeGatewayImport
func resourceVcdNsxtEdgeGatewayImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[TRACE] edge gateway import initiated")

	resourceURI := strings.Split(d.Id(), ImportSeparator)
	if len(resourceURI) != 3 {
		return nil, fmt.Errorf("resource name must be specified as org-name.vdc-name.edge-gw-name (or edge-gw-ID)")
	}
	orgName, _, edgeName := resourceURI[0], resourceURI[1], resourceURI[2]

	vcdClient := meta.(*VCDClient)
	adminOrg, err := vcdClient.GetAdminOrg(orgName)
	if err != nil {
		return nil, fmt.Errorf("unable to find org %s: %s", orgName, err)
	}
	// vdc, err := adminOrg.GetVDCByName(vdcName, false)
	// if err != nil {
	// 	return nil, fmt.Errorf("unable to find VDC %s: %s", vdcName, err)
	// }

	edge, err := adminOrg.GetNsxtEdgeGatewayByName(edgeName)
	if err != nil {
		return nil, fmt.Errorf("could not retrieve edge gateway with ID '%s': %s", d.Id(), err)
	}

	d.SetId(edge.EdgeGateway.ID)

	return []*schema.ResourceData{d}, nil
}

// getNsxtEdgeGatewayType
func getNsxtEdgeGatewayType(d *schema.ResourceData, vdc *govcd.Vdc) (*types.NsxtEdgeGateway, error) {

	e := types.NsxtEdgeGateway{
		// Status:                    "",
		// ID:                        "",
		Name:                      d.Get("name").(string),
		Description:               d.Get("description").(string),
		DistributedRoutingEnabled: true, // ???NSX-T is always distributed???
		EdgeGatewayUplinks: []types.EdgeGatewayUplinks{types.EdgeGatewayUplinks{
			UplinkID:                 "",
			UplinkName:               "",
			Subnets:                  types.NsxtSubnets{},
			Connected:                false,
			QuickAddAllocatedIPCount: nil,
			Dedicated:                false,
		}},
		// OrgVdcNetworkCount:        0,
		GatewayBacking: types.GatewayBacking{},
		OrgVdc:         types.OrgVdc{},
		OrgRef:         types.OrgRef{},
		// ServiceNetworkDefinition:  "",
		// EdgeClusterConfig:         types.EdgeClusterConfig{},
	}

	// e := types.NsxtEdgeGateway{
	// 	// Status:      "",
	// 	// ID:          "",
	// 	Name:        d.Get("name").(string),
	// 	Description: d.Get("description").(string),
	// 	// OrgVdc: struct {
	// 	// 	ID string `json:"id"`
	// 	// }{ID: vdc.Vdc.ID},
	// 	// EdgeGatewayUplinks: nil,
	// }
	//
	// t, err := getNsxtEdgeGatewayType(d, vdc)
	// if err != nil {
	// 	return fmt.Errorf("could not create edge gateway type: %s", err)
	// }
	//
	// _, err = adminOrg.UpdateNsxtEdgeGateway(t)

	return &e, nil
}

// setNsxtEdgeGatewayData
func setNsxtEdgeGatewayData(e *types.NsxtEdgeGateway, d *schema.ResourceData) error {
	return nil
}
