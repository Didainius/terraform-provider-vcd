package vcd

import (
	"context"
	"fmt"

	"github.com/vmware/go-vcloud-director/v2/govcd"
	"github.com/vmware/go-vcloud-director/v2/types/v56"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceVcdNetworkRoutedV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVcdNetworkRoutedV2Create,
		ReadContext:   resourceVcdNetworkRoutedV2Read,
		UpdateContext: resourceVcdNetworkRoutedV2Update,
		DeleteContext: resourceVcdNetworkRoutedV2Delete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceVcdNetworkRoutedV2Import,
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
				Description: "Edge gateway name in which NAT Rule is located",
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"gateway": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"prefix_length": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},
			"dns1": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"dns2": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"dns_suffix": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"static_ip_pool": &schema.Schema{
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "IP ranges used for static pool allocation in the network",
				Elem:        networkV2IpRange,
			},
		},
	}
}

// resourceVcdNetworkRoutedV2Create
func resourceVcdNetworkRoutedV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vcdClient := meta.(*VCDClient)

	// vcdClient.lockParentEdgeGtw(d)
	// defer vcdClient.unLockParentEdgeGtw(d)

	_, vdc, err := vcdClient.GetOrgAndVdcFromResource(d)
	if err != nil {
		return diag.Errorf("error retrieving VDC: %s", err)
	}

	networkType, err := getOpenApiOrgVdcNetworkType(d, vdc)
	if err != nil {
		return diag.FromErr(err)
	}

	orgNetwork, err := vdc.CreateNsxtOrgVdcNetwork(networkType)
	if err != nil {
		return diag.Errorf("error creating Org Vdc routed network: %s", err)
	}

	d.SetId(orgNetwork.OrgVdcNetwork.ID)

	return resourceVcdNetworkRoutedV2Read(ctx, d, meta)
}

// resourceVcdNetworkRoutedV2Update
func resourceVcdNetworkRoutedV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vcdClient := meta.(*VCDClient)

	// vcdClient.lockParentEdgeGtw(d)
	// defer vcdClient.unLockParentEdgeGtw(d)

	_, vdc, err := vcdClient.GetOrgAndVdcFromResource(d)
	if err != nil {
		return diag.Errorf("error retrieving VDC: %s", err)
	}

	orgNetwork, err := vdc.GetNsxtOrgVdcNetworkById(d.Id())
	// If object is not found -
	if govcd.ContainsNotFound(err) {
		d.SetId("")
		return nil
	}
	if err != nil {
		return diag.Errorf("error getting Org Vdc network: %s", err)
	}

	networkType, err := getOpenApiOrgVdcNetworkType(d, vdc)
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = orgNetwork.Update(networkType)
	if err != nil {
		return diag.Errorf("error updating Org Vdc network: %s", err)
	}

	return resourceVcdNetworkRoutedV2Read(ctx, d, meta)
}

// resourceVcdNetworkRoutedV2Read
func resourceVcdNetworkRoutedV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vcdClient := meta.(*VCDClient)

	// vcdClient.lockParentEdgeGtw(d)
	// defer vcdClient.unLockParentEdgeGtw(d)

	_, vdc, err := vcdClient.GetOrgAndVdcFromResource(d)
	if err != nil {
		return diag.Errorf("error retrieving VDC: %s", err)
	}

	orgNetwork, err := vdc.GetNsxtOrgVdcNetworkById(d.Id())
	// If object is not found -
	if govcd.ContainsNotFound(err) {
		d.SetId("")
		return nil
	}
	if err != nil {
		return diag.Errorf("error getting Org Vdc network: %s", err)
	}

	err = setOpenApiOrgVdcNetworkData(d, orgNetwork.OrgVdcNetwork)
	if err != nil {
		return diag.Errorf("error setting Org Vdc network data: %s", err)
	}

	d.SetId(orgNetwork.OrgVdcNetwork.ID)

	return nil
}

// resourceVcdNetworkRoutedV2Delete
func resourceVcdNetworkRoutedV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vcdClient := meta.(*VCDClient)

	// vcdClient.lockParentEdgeGtw(d)
	// defer vcdClient.unLockParentEdgeGtw(d)

	_, vdc, err := vcdClient.GetOrgAndVdcFromResource(d)
	if err != nil {
		return diag.Errorf("error retrieving VDC: %s", err)
	}

	orgNetwork, err := vdc.GetNsxtOrgVdcNetworkById(d.Id())
	if err != nil {
		return diag.Errorf("error getting Org Vdc network: %s", err)
	}

	orgNetwork.Delete()
	if err != nil {
		return diag.Errorf("error deleting Org Vdc network: %s", err)
	}

	return nil
}

// resourceVcdNetworkRoutedV2Import
func resourceVcdNetworkRoutedV2Import(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	return []*schema.ResourceData{d}, nil
}

func setOpenApiOrgVdcNetworkData(d *schema.ResourceData, orgVdcNetwork *types.OpenApiOrgVdcNetwork) error {

	_ = d.Set("name", orgVdcNetwork.Name)
	_ = d.Set("description", orgVdcNetwork.Description)
	// Check if values are not empty
	_ = d.Set("edge_gateway_id", orgVdcNetwork.Connection.RouterRef.ID)

	_ = d.Set("gateway", orgVdcNetwork.Subnets.Values[0].Gateway)
	_ = d.Set("prefix_length", orgVdcNetwork.Subnets.Values[0].PrefixLength)
	_ = d.Set("dns1", orgVdcNetwork.Subnets.Values[0].DNSServer1)
	_ = d.Set("dns2", orgVdcNetwork.Subnets.Values[0].DNSServer2)
	_ = d.Set("dns_suffix", orgVdcNetwork.Subnets.Values[0].DNSSuffix)

	// If any IP sets are available
	if len(orgVdcNetwork.Subnets.Values[0].IPRanges.Values) > 0 {
		ipRangeSlice := make([]interface{}, len(orgVdcNetwork.Subnets.Values[0].IPRanges.Values))
		for ii, ipRange := range orgVdcNetwork.Subnets.Values[0].IPRanges.Values {
			ipRangeMap := make(map[string]interface{})
			ipRangeMap["start_address"] = ipRange.StartAddress
			ipRangeMap["end_address"] = ipRange.EndAddress

			ipRangeSlice[ii] = ipRangeMap
		}
		ipRangeSet := schema.NewSet(schema.HashResource(networkV2IpRange), ipRangeSlice)

		err := d.Set("static_ip_pool", ipRangeSet)
		if err != nil {
			return fmt.Errorf("error setting 'static_ip_pool': %s", err)
		}
	}

	return nil
}

func getOpenApiOrgVdcNetworkType(d *schema.ResourceData, vdc *govcd.Vdc) (*types.OpenApiOrgVdcNetwork, error) {
	orgVdcNetworkConfig := &types.OpenApiOrgVdcNetwork{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		OrgVdc:      &types.OpenApiReference{ID: vdc.Vdc.ID},

		NetworkType: types.OrgVdcNetworkTypeRouted,

		// Connection is used for "routed" network
		Connection: &types.Connection{
			RouterRef: types.OpenApiReference{
				ID: d.Get("edge_gateway_id").(string),
			},
			ConnectionType: "INTERNAL",
		},
		Subnets: types.OrgVdcNetworkSubnets{
			Values: []types.OrgVdcNetworkSubnetValues{
				{
					Gateway:      d.Get("gateway").(string),
					PrefixLength: d.Get("prefix_length").(int),
					DNSServer1:   d.Get("dns1").(string),
					DNSServer2:   d.Get("dns2").(string),
					DNSSuffix:    d.Get("dns_suffix").(string),
					IPRanges: types.OrgVdcNetworkSubnetIPRanges{
						Values: processIpRangesS(d.Get("static_ip_pool").(*schema.Set)),
					},
				},
			},
		},
	}

	return orgVdcNetworkConfig, nil
}

func processIpRangesS(staticIpPool *schema.Set) []types.ExternalNetworkV2IPRange {
	subnetRng := make([]types.ExternalNetworkV2IPRange, len(staticIpPool.List()))
	for rangeIndex, subnetRange := range staticIpPool.List() {
		subnetRangeStr := convertToStringMap(subnetRange.(map[string]interface{}))
		oneRange := types.ExternalNetworkV2IPRange{
			StartAddress: subnetRangeStr["start_address"],
			EndAddress:   subnetRangeStr["end_address"],
		}
		subnetRng[rangeIndex] = oneRange
	}
	return subnetRng
}
