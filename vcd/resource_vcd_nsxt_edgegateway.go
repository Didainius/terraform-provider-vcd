package vcd

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

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
			// "org": &schema.Schema{
			// 	Type:     schema.TypeString,
			// 	Optional: true,
			// 	ForceNew: true,
			// 	Description: "The name of organization to use, optional if defined at provider " +
			// 		"level. Useful when connected as sysadmin working across different organizations",
			// },
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
			"subnets": {
				Optional: true,
				Computed: true,
				ForceNew: true,
				Type:     schema.TypeSet,
				MinItems: 1,
				Elem: &schema.Resource{
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
				},
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

	// Making sure the parent entities are available
	orgName := vcdClient.getOrgName(d)
	vdcName := vcdClient.getVdcName(d)

	var missing []string
	if orgName == "" {
		missing = append(missing, "org")
	}
	if vdcName == "" {
		missing = append(missing, "vdc")
	}
	if len(missing) > 0 {
		return fmt.Errorf("missing properties. %v should be given either in the resource or at provider level", missing)
	}

	org, vdc, err := vcdClient.GetOrgAndVdc(orgName, vdcName)
	if err != nil {
		return err
	}
	if org == nil {
		return fmt.Errorf("no valid Organization named '%s' was found", orgName)
	}
	if vdc == nil || vdc.Vdc.HREF == "" || vdc.Vdc.ID == "" || vdc.Vdc.Name == "" {
		return fmt.Errorf("no valid VDC named '%s' was found", vdcName)
	}

	adminOrg, err := vcdClient.GetAdminOrgFromResource(d)
	if err != nil {
		return fmt.Errorf("error getting adminOrg: %s", err)
	}

	adminOrg.CreateNsxtEdgeGateway()

	return nil
}

// resourceVcdNsxtEdgeGatewayUpdate
func resourceVcdNsxtEdgeGatewayUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[TRACE] edge gateway update initiated")

	vcdClient := meta.(*VCDClient)
	vdcName := vcdClient.getVdcName(d)
	vdc, err := vcdClient.getVdc(vdcName)
	if err != nil {
		return fmt.Errorf("could not get VDC '%s': %s", vdcName, err)
	}

	return nil
}

// resourceVcdNsxtEdgeGatewayRead
func resourceVcdNsxtEdgeGatewayRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[TRACE] edge gateway read initiated")

	vcdClient := meta.(*VCDClient)
	vdcName := vcdClient.getVdcName(d)
	vdc, err := vcdClient.getVdc(vdcName)
	if err != nil {
		return fmt.Errorf("could not get VDC '%s': %s", vdcName, err)
	}

	return nil
}

// resourceVcdNsxtEdgeGatewayDelete
func resourceVcdNsxtEdgeGatewayDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[TRACE] edge gateway deletion initiated")

	vcdClient := meta.(*VCDClient)
	vdcName := vcdClient.getVdcName(d)
	vdc, err := vcdClient.getVdc(vdcName)
	if err != nil {
		return fmt.Errorf("could not get VDC '%s': %s", vdcName, err)
	}

	return nil
}

// resourceVcdNsxtEdgeGatewayImport
func resourceVcdNsxtEdgeGatewayImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[TRACE] edge gateway import initiated")

	vcdClient := meta.(*VCDClient)
	vdcName := vcdClient.getVdcName(d)
	vdc, err := vcdClient.getVdc(vdcName)
	if err != nil {
		return nil, fmt.Errorf("could not get VDC '%s': %s", vdcName, err)
	}

	return []*schema.ResourceData{d}, nil
}
