package vcd

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceVcdExternalNetworkV2() *schema.Resource {
	return &schema.Resource{
		Create: resourceVcdExternalNetworkV2Create,
		Update: resourceVcdExternalNetworkV2Update,
		Delete: resourceVcdExternalNetworkV2Delete,
		Read:   resourceVcdExternalNetworkV2Read,
		Importer: &schema.ResourceImporter{
			State: resourceVcdExternalNetworkImport,
		},
		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"ip_scope": &schema.Schema{
				Type:        schema.TypeList,
				Required:    true,
				ForceNew:    true,
				Description: "A list of IP scopes for the network",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"gateway": &schema.Schema{
							Type:         schema.TypeString,
							Required:     true,
							ForceNew:     true,
							Description:  "Gateway of the network",
							ValidateFunc: validation.IsIPAddress,
						},
						"netmask": &schema.Schema{
							Type:         schema.TypeString,
							Required:     true,
							ForceNew:     true,
							Description:  "Network mask",
							ValidateFunc: validation.IsIPAddress,
						},
						"dns1": &schema.Schema{
							Type:         schema.TypeString,
							Optional:     true,
							ForceNew:     true,
							Description:  "Primary DNS server",
							ValidateFunc: validation.IsIPAddress,
						},
						"dns2": &schema.Schema{
							Type:         schema.TypeString,
							Optional:     true,
							ForceNew:     true,
							Description:  "Secondary DNS server",
							ValidateFunc: validation.IsIPAddress,
						},
						"dns_suffix": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Description: "DNS suffix",
						},
						"static_ip_pool": &schema.Schema{
							Type:        schema.TypeList,
							Optional:    true,
							ForceNew:    true,
							Description: "IP ranges used for static pool allocation in the network",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"start_address": &schema.Schema{
										Type:         schema.TypeString,
										Required:     true,
										ForceNew:     true,
										Description:  "Start address of the IP range",
										ValidateFunc: validation.IsIPAddress,
									},
									"end_address": &schema.Schema{
										Type:         schema.TypeString,
										Required:     true,
										ForceNew:     true,
										Description:  "End address of the IP range",
										ValidateFunc: validation.IsIPAddress,
									},
								},
							},
						},
					},
				},
			},
			"vsphere_network": &schema.Schema{
				Type:        schema.TypeList,
				Required:    true,
				ForceNew:    true,
				Description: "A list of port groups that back this network. Each referenced DV_PORTGROUP or NETWORK must exist on a vCenter server registered with the system.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"vcenter": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							ForceNew:    true,
							Description: "The vCenter server name",
						},
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							ForceNew:    true,
							Description: "The name of the port group",
						},
						"type": &schema.Schema{
							Type:         schema.TypeString,
							Required:     true,
							ForceNew:     true,
							Description:  "The vSphere port group type. One of: DV_PORTGROUP (distributed virtual port group), NETWORK",
							ValidateFunc: validation.StringInSlice([]string{"DV_PORTGROUP", "NETWORK"}, false),
						},
					},
				},
			},
			"retain_net_info_across_deployments": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Default:     false,
				Description: "Specifies whether the network resources such as IP/MAC of router will be retained across deployments. Default is false.",
			},
		},
	}
}

func resourceVcdExternalNetworkV2Create(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceVcdExternalNetworkV2Update(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceVcdExternalNetworkV2Read(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceVcdExternalNetworkV2Delete(d *schema.ResourceData, meta interface{}) error {
	return nil
}
