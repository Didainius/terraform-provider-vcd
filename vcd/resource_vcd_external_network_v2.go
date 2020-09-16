package vcd

import (
	"fmt"
	"log"
	"strconv"

	"github.com/vmware/go-vcloud-director/v2/types/v56"

	"github.com/vmware/go-vcloud-director/v2/govcd"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceVcdExternalNetworkV2() *schema.Resource {
	return &schema.Resource{
		Create: resourceVcdExternalNetworkV2Create,
		// Update: resourceVcdExternalNetworkV2Update,
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
				Type:        schema.TypeSet,
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
							Type:        schema.TypeString,
							Required:    true,
							ForceNew:    true,
							Description: "Network mask",
							// ValidateFunc: validation.IsIPAddress,
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
						// "enabled": &schema.Schema{
						// 	Type:        schema.TypeString,
						// 	Optional:    true,
						// 	ForceNew:    true,
						// 	Description: "If subnet is enabled",
						// },
						"static_ip_pool": &schema.Schema{
							Type:        schema.TypeSet,
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
				Type:     schema.TypeSet,
				Optional: true,
				// ExactlyOneOf: []string{"vsphere_network", "nsxt_network"},
				ForceNew:    true,
				Description: "A list of port groups that back this network. Each referenced DV_PORTGROUP or NETWORK must exist on a vCenter server registered with the system.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"vcenter_id": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							ForceNew:    true,
							Description: "The vCenter server name",
						},
						"portgroup_id": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							ForceNew:    true,
							Description: "The name of the port group",
						},
						"portgroup_type": &schema.Schema{
							Type:         schema.TypeString,
							Required:     true,
							ForceNew:     true,
							Description:  "The vSphere port group type. One of: DV_PORTGROUP (distributed virtual port group), NETWORK",
							ValidateFunc: validation.StringInSlice([]string{"DV_PORTGROUP", "NETWORK"}, false),
						},
					},
				},
			},
			"nsxt_network": &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				// ExactlyOneOf: []string{"vsphere_network", "nsxt_network"},
				ForceNew:    true,
				Description: "A list of port groups that back this network. Each referenced DV_PORTGROUP or NETWORK must exist on a vCenter server registered with the system.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"nsxt_manager_id": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							ForceNew:    true,
							Description: "ID of NSX-T manager",
						},
						"nsxt_tier0_router_id": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							ForceNew:    true,
							Description: "The vSphere port group type. One of: DV_PORTGROUP (distributed virtual port group), NETWORK",
							// ValidateFunc: validation.StringInSlice([]string{"DV_PORTGROUP", "NETWORK"}, false),
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
	vcdClient := meta.(*VCDClient)
	log.Printf("[TRACE] external network V2 creation initiated")

	netType, err := getExternalNetworkV2Type(d)
	if err != nil {
		return fmt.Errorf("could not get network data: %s", err)
	}

	extNet, err := govcd.CreateExternalNetworkV2(vcdClient.VCDClient, netType)
	if err != nil {
		return fmt.Errorf("error applying data: %s", err)
	}

	// Only store ID and leave all the rest to "READ"
	d.SetId(extNet.ExternalNetwork.ID)

	return resourceVcdExternalNetworkV2Read(d, meta)
}

func resourceVcdExternalNetworkV2Update(d *schema.ResourceData, meta interface{}) error {
	vcdClient := meta.(*VCDClient)
	log.Printf("[TRACE] update network V2 creation initiated")

	extNet, err := govcd.GetExternalNetworkById(vcdClient.VCDClient, d.Id())
	if err != nil {
		return fmt.Errorf("could not find external network by ID '%s': %s", d.Id(), err)
	}

	netType, err := getExternalNetworkV2Type(d)
	if err != nil {
		return fmt.Errorf("could not get network data: %s", err)
	}

	netType.ID = extNet.ExternalNetwork.ID
	extNet.ExternalNetwork = netType

	_, err = extNet.Update()

	return err
}

func resourceVcdExternalNetworkV2Read(d *schema.ResourceData, meta interface{}) error {
	vcdClient := meta.(*VCDClient)
	log.Printf("[TRACE] external network V2 creation initiated")

	extNet, err := govcd.GetExternalNetworkById(vcdClient.VCDClient, d.Id())
	if err != nil {
		return fmt.Errorf("could not find external network by ID '%s': %s", d.Id(), err)
	}

	return setExternalNetworkV2Data(d, extNet.ExternalNetwork)
}

func resourceVcdExternalNetworkV2Delete(d *schema.ResourceData, meta interface{}) error {
	vcdClient := meta.(*VCDClient)
	log.Printf("[TRACE] external network V2 creation initiated")

	extNet, err := govcd.GetExternalNetworkById(vcdClient.VCDClient, d.Id())
	if err != nil {
		return fmt.Errorf("could not find external network by ID '%s': %s", d.Id(), err)
	}

	return extNet.Delete()
}

func getExternalNetworkV2Type(d *schema.ResourceData) (*types.ExternalNetworkV2, error) {

	// Subnets
	subnets := d.Get("ip_scope").(*schema.Set)
	listSubnets := make([]types.Subnet, len(subnets.List()))
	for subnetIndex, subnet := range subnets.List() { // Loop over ip_scopes

		subnetMap := subnet.(map[string]interface{})

		prefixInt, _ := strconv.Atoi(subnetMap["netmask"].(string))
		subnet := types.Subnet{
			Gateway:      subnetMap["gateway"].(string),
			DNSSuffix:    subnetMap["dns_suffix"].(string),
			DNSServer1:   subnetMap["dns1"].(string),
			DNSServer2:   subnetMap["dns2"].(string),
			PrefixLength: prefixInt,
			Enabled:      true,
			// IPRanges:     types.IPRanges2{},
			// UsedIPCount:  0,
			// TotalIPCount: 0,
		}

		// Loop over IP ranges (static IP pools)
		rrr := subnetMap["static_ip_pool"].(*schema.Set)
		subnetRng := make([]types.IPRange2, len(rrr.List()))
		for rangeIndex, subnetRange := range rrr.List() {
			subnetRangeStr := convertToStringMap(subnetRange.(map[string]interface{}))
			oneRange := types.IPRange2{
				StartAddress: subnetRangeStr["start_address"],
				EndAddress:   subnetRangeStr["end_address"],
			}
			// subnetRng = append(subnetRng, oneRange)
			subnetRng[rangeIndex] = oneRange
		}
		// Add all ranges
		subnet.IPRanges = types.IPRanges2{Values: subnetRng}

		// listSubnets = append(listSubnets, subnet)
		listSubnets[subnetIndex] = subnet
	}

	//
	// relayAgentsSlice := relayAgentsSet.List()
	// relayAgentsStruct := make([]types.EdgeDhcpRelayAgent, len(relayAgentsSlice))
	// for index, relayAgent := range relayAgentsSlice {
	// 	relayAgentMap := convertToStringMap(relayAgent.(map[string]interface{}))

	var backing types.NetworkBacking
	// Network backings - NSX-T
	nsxtNetwork := d.Get("nsxt_network").(*schema.Set)
	nsxtNetworkSlice := nsxtNetwork.List()
	if len(nsxtNetworkSlice) > 0 {
		nsxtNetworkStrings := convertToStringMap(nsxtNetworkSlice[0].(map[string]interface{}))
		backing = types.NetworkBacking{
			BackingID:   nsxtNetworkStrings["nsxt_tier0_router_id"], // Tier 0- router
			Name:        "",
			BackingType: types.ExternalNetworkBackingTypeNsxtTier0Router,
			NetworkProvider: types.NetworkProvider{
				// Name: "",
				ID: nsxtNetworkStrings["nsxt_manager_id"], // NSX-T manager
			},
		}
	}

	// Network backings - NSX-V
	nsxvNetwork := d.Get("vsphere_network").(*schema.Set)
	nsxvNetworkSlice := nsxvNetwork.List()
	if len(nsxvNetworkSlice) > 0 {
		nsxvNetworkStrings := convertToStringMap(nsxvNetworkSlice[0].(map[string]interface{}))
		backing = types.NetworkBacking{
			BackingID:   nsxvNetworkStrings["portgroup_id"],
			BackingType: nsxvNetworkStrings["portgroup_type"],
			NetworkProvider: types.NetworkProvider{
				ID: nsxvNetworkStrings["vcenter_id"],
			},
		}
	}

	newExtNet := &types.ExternalNetworkV2{
		Name:            d.Get("name").(string),
		Description:     d.Get("description").(string),
		Subnets:         types.Subnets{Values: listSubnets},
		NetworkBackings: types.NetworkBackings{[]types.NetworkBacking{backing}},
	}

	return newExtNet, nil
}

func setExternalNetworkV2Data(d *schema.ResourceData, net *types.ExternalNetworkV2) error {
	return nil
}
