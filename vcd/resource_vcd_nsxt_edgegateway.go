package vcd

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vmware/go-vcloud-director/v2/govcd"
	"github.com/vmware/go-vcloud-director/v2/types/v56"
)

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
			Description: "Netmask address for a subnet (e.g. /24)",
			Type:        schema.TypeInt,
		},
		"enabled": {
			Optional: true,
			Default:  true,
			// ForceNew:    true,
			Description: "",
			Type:        schema.TypeBool,
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
		"allocated_ips": {
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
			"external_network_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "External network ID",
			},
			"nsxt_manager_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "NSX-T manager ID",
			},
			"subnet": {
				Description: "One or more blocks with external network information to be attached to this gateway's interface",
				ForceNew:    true,
				Required:    true,
				Type:        schema.TypeSet,
				Elem:        subnetResource2,
			},
			"edge_cluster_id": &schema.Schema{
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

	t, err := getNsxtEdgeGatewayType(d, adminOrg, vdc)
	if err != nil {
		return fmt.Errorf("could not create edge gateway type: %s", err)
	}

	createdEdgeGateway, err := adminOrg.CreateNsxtEdgeGateway(t)
	if err != nil {
		return fmt.Errorf("error creating edge gateway: %s", err)
	}

	d.SetId(createdEdgeGateway.EdgeGateway.ID)

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

	edge.EdgeGateway, err = getNsxtEdgeGatewayType(d, adminOrg, vdc)
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
func getNsxtEdgeGatewayType(d *schema.ResourceData, adminOrg *govcd.AdminOrg, vdc *govcd.Vdc) (*types.OpenAPIEdgeGateway, error) {

	e := types.OpenAPIEdgeGateway{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		// DistributedRoutingEnabled: true, // ???NSX-T is always distributed???
		EdgeGatewayUplinks: []types.EdgeGatewayUplinks{types.EdgeGatewayUplinks{
			UplinkID: d.Get("external_network_id").(string),
			// UplinkName:               "",
			Subnets: types.EdgeGatewaySubnets{getNsxtEdgeGatewayUplinks(d)},
			// Connected:                false,
			// QuickAddAllocatedIPCount: nil,
			Dedicated: d.Get("dedicate_external_network").(bool),
		}},
		// OrgVdcNetworkCount:        0,
		// GatewayBacking: types.GatewayBacking{},
		OrgVdc: &types.OpenApiReference{
			ID: vdc.Vdc.ID,
		},
		// Org: types.Org{
		// 	ID: adminOrg.AdminOrg.ID,
		// },
		// ServiceNetworkDefinition:  "",
		// EdgeClusterConfig:         types.EdgeClusterConfig{},
		// GatewayBacking: types.GatewayBacking{
		// 	GatewayType: "NSXT_BACKED",
		// 	// GatewayType: "NSXT_IMPORT",
		// 	NetworkProvider: types.NetworkProvider{
		// 		ID: d.Get("nsxt_manager_id").(string),
		// 	},
		// },
	}

	// Optional edge_cluster_id
	if clusterId, isSet := d.GetOk("edge_cluster_id"); isSet {
		e.EdgeClusterConfig.PrimaryEdgeCluster.BackingID = clusterId.(string)
	}

	return &e, nil
}

func getNsxtEdgeGatewayUplinks(d *schema.ResourceData) []types.EdgeGatewaySubnetValue {
	extNetworks := d.Get("subnet").(*schema.Set).List()
	subnetSlice := make([]types.EdgeGatewaySubnetValue, len(extNetworks))

	for index, singleSubnet := range extNetworks {
		subnetMap := singleSubnet.(map[string]interface{})
		subn := types.EdgeGatewaySubnetValue{
			Gateway:      subnetMap["gateway"].(string),
			PrefixLength: subnetMap["prefix_length"].(int),
			Enabled:      subnetMap["enabled"].(bool),
			// TotalIPCount:         0,
			// UsedIPCount:          nil,
			// PrimaryIP:            "",
			// AutoAllocateIPRanges: false,
		}
		// Only feed in ip range allocations if they are defined
		if ipRanges := getNsxtEdgeGatewayUplinkRanges(subnetMap); ipRanges != nil {
			subn.IPRanges = &types.OpenApiIPRanges{ipRanges}
		}

		subnetSlice[index] = subn
	}

	return subnetSlice
}

func getNsxtEdgeGatewayUplinkRanges(subnetMap map[string]interface{}) []types.OpenApiIPRangeValues {
	suballocatePoolSchema := subnetMap["allocated_ips"].(*schema.Set)
	subnetRanges := make([]types.OpenApiIPRangeValues, len(suballocatePoolSchema.List()))

	if len(subnetRanges) == 0 {
		return nil
	}

	for rangeIndex, subnetRange := range suballocatePoolSchema.List() {
		subnetRangeStr := convertToStringMap(subnetRange.(map[string]interface{}))
		oneRange := types.OpenApiIPRangeValues{
			StartAddress: subnetRangeStr["start_address"],
			EndAddress:   subnetRangeStr["end_address"],
		}
		subnetRanges[rangeIndex] = oneRange
	}
	return subnetRanges
}

// setNsxtEdgeGatewayData
func setNsxtEdgeGatewayData(e *types.OpenAPIEdgeGateway, d *schema.ResourceData) error {
	return nil
}
