package vcd

import (
	"context"

	"github.com/vmware/go-vcloud-director/v2/types/v56"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var nsxvDhcpPoolSetSchema = &schema.Resource{
	Schema: map[string]*schema.Schema{
		"ip_range": &schema.Schema{
			Type:        schema.TypeString,
			Required:    true,
			Description: "IP range for DHCP pools",
		},
		"domain_name": &schema.Schema{
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Domain name",
		},
		"autoconfigure_dns": &schema.Schema{
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "End address of the IP range",
		},
		"dns1": &schema.Schema{
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Primary DNS server",
		},
		"dns2": &schema.Schema{
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Secondary DNS server",
		},
		"gateway": &schema.Schema{
			Type:         schema.TypeString,
			Required:     true,
			Description:  "Default gateway",
			ValidateFunc: validation.IsIPAddress,
		},
		"subnet_mask": &schema.Schema{
			Type:        schema.TypeString,
			Required:    true,
			Description: "Subnet mask",
		},
		"lease_time": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
			// Computed: true,
			Default:     "86400",
			Description: "Lease time in seconds or 'infinite'",
		},
	},
}

func resourceVcdNsxvDhcpPools() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVcdNsxvDhcpCreate,
		ReadContext:   resourceVcdNsxvDhcpRead,
		UpdateContext: resourceVcdNsxvDhcpUpdate,
		DeleteContext: resourceVcdNsxvDhcpDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceVcdNsxvDhcpImport,
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
				Description: "Edge gateway ID for DHCP relay settings",
			},
			"dhcp_enabled": &schema.Schema{
				Type:        schema.TypeBool,
				Default:     true,
				Optional:    true,
				Description: "If DHCP service is enabled. Default 'true'",
			},
			"dhcp_pool": &schema.Schema{
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "IP ranges used for static pool allocation in the network",
				Elem:        nsxvDhcpPoolSetSchema,
			},
		},
	}
}

// resourceVcdNsxvDhcpCreate
func resourceVcdNsxvDhcpCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vcdClient := meta.(*VCDClient)
	vcdClient.lockParentEdgeGtw(d)
	defer vcdClient.unLockParentEdgeGtw(d)

	_, vdc, err := vcdClient.GetOrgAndVdcFromResource(d)
	if err != nil {
		return diag.Errorf("[NSX-V DHCP pool create] error retrieving VDC: %s", err)
	}
	egw, err := vdc.GetEdgeGatewayById(d.Get("edge_gateway_id").(string), false)
	if err != nil {
		return diag.Errorf("[NSX-V DHCP pool create] error looking up edge gateway: %s", err)
	}

	edgeDhcpType := getEdgeDhcpType(d)
	_, err = egw.UpdateDhcpPools(edgeDhcpType)

	if err != nil {
		return diag.Errorf("[NSX-V DHCP pool create] error setting DHCP pools: %s", err)
	}

	d.SetId(egw.EdgeGateway.ID)

	return resourceVcdNsxvDhcpRead(ctx, d, meta)
}

// resourceVcdNsxvDhcpUpdate
func resourceVcdNsxvDhcpUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vcdClient := meta.(*VCDClient)
	vcdClient.lockParentEdgeGtw(d)
	defer vcdClient.unLockParentEdgeGtw(d)

	_, vdc, err := vcdClient.GetOrgAndVdcFromResource(d)
	if err != nil {
		return diag.Errorf("[NSX-V DHCP pool update] error retrieving VDC: %s", err)
	}
	egw, err := vdc.GetEdgeGatewayById(d.Get("edge_gateway_id").(string), false)
	if err != nil {

		return diag.Errorf("[NSX-V DHCP pool update] error looking up edge gateway: %s", err)
	}

	edgeDhcpType := getEdgeDhcpType(d)
	_, err = egw.UpdateDhcpPools(edgeDhcpType)

	if err != nil {
		return diag.Errorf("[NSX-V DHCP pool update] error setting DHCP pools: %s", err)
	}

	return resourceVcdNsxvDhcpRead(ctx, d, meta)
}

// resourceVcdNsxvDhcpRead
func resourceVcdNsxvDhcpRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vcdClient := meta.(*VCDClient)

	_, vdc, err := vcdClient.GetOrgAndVdcFromResource(d)
	if err != nil {
		return diag.Errorf("[NSX-V DHCP pool read] error retrieving VDC: %s", err)
	}
	egw, err := vdc.GetEdgeGatewayById(d.Get("edge_gateway_id").(string), false)
	if err != nil {

		return diag.Errorf("[NSX-V DHCP pool read] error looking up edge gateway: %s", err)
	}

	dhcpConfig, err := egw.GetDhcpPoolsAndBindings()
	if err != nil {
		return diag.Errorf("[NSX-V DHCP pool read] error getting DHCP config: %s", err)
	}

	err = setNsxvDhcpData(dhcpConfig, d)
	if err != nil {
		return diag.Errorf("[NSX-V DHCP pool read] error setting DHCP config: %s", err)
	}

	return nil
}

// resourceVcdNsxvDhcpDelete
func resourceVcdNsxvDhcpDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vcdClient := meta.(*VCDClient)
	vcdClient.lockParentEdgeGtw(d)
	defer vcdClient.unLockParentEdgeGtw(d)

	_, vdc, err := vcdClient.GetOrgAndVdcFromResource(d)
	if err != nil {
		return diag.Errorf("[NSX-V DHCP pool delete] error retrieving VDC: %s", err)
	}
	egw, err := vdc.GetEdgeGatewayById(d.Get("edge_gateway_id").(string), false)
	if err != nil {
		return diag.Errorf("[NSX-V DHCP pool delete] error looking up edge gateway: %s", err)
	}

	err = egw.ResetDhcpPools()
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

// resourceVcdNsxvDhcpImport
func resourceVcdNsxvDhcpImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	return []*schema.ResourceData{d}, nil
}

func setNsxvDhcpData(dhcpConfig *types.EdgeDhcp, d *schema.ResourceData) error {
	err := d.Set("dhcp_enabled", dhcpConfig.Enabled)
	if err != nil {
		return err
	}

	if len(dhcpConfig.EdgeDhcpIpPools.EdgeDhcpIpPool) > 0 {
		poolInterfaceSlice := make([]interface{}, len(dhcpConfig.EdgeDhcpIpPools.EdgeDhcpIpPool))

		for index, pool := range dhcpConfig.EdgeDhcpIpPools.EdgeDhcpIpPool {
			onePool := make(map[string]interface{})
			onePool["autoconfigure_dns"] = pool.AutoConfigureDNS
			onePool["gateway"] = pool.DefaultGateway
			onePool["domain_name"] = pool.DomainName
			onePool["ip_range"] = pool.IpRange
			onePool["lease_time"] = pool.LeaseTime
			onePool["subnet_mask"] = pool.SubnetMask
			onePool["dns1"] = pool.PrimaryNameServer
			onePool["dns2"] = pool.SecondaryNameServer

			poolInterfaceSlice[index] = onePool
		}

		dhcpPoolSet := schema.NewSet(schema.HashResource(nsxvDhcpPoolSetSchema), poolInterfaceSlice)
		err = d.Set("dhcp_pool", dhcpPoolSet)
		if err != nil {
			return err
		}
	}

	return nil

}

func getEdgeDhcpType(d *schema.ResourceData) *types.EdgeDhcp {
	edgeDhcp := types.EdgeDhcp{
		Enabled: d.Get("dhcp_enabled").(bool),
		Logging: &types.LbLogging{Enable: false, LogLevel: "info"},
	}

	dhcpPool := d.Get("dhcp_pool")
	if dhcpPool == nil {
		return nil
	}

	dhcpPoolSet := dhcpPool.(*schema.Set)
	dhcpPoolList := dhcpPoolSet.List()

	if len(dhcpPoolList) > 0 {
		dhcpPools := make([]types.EdgeDhcpIpPool, len(dhcpPoolList))
		for index, pool := range dhcpPoolList {
			poolMap := pool.(map[string]interface{})
			onePool := types.EdgeDhcpIpPool{
				AutoConfigureDNS:    poolMap["autoconfigure_dns"].(bool),
				DefaultGateway:      poolMap["gateway"].(string),
				DomainName:          poolMap["domain_name"].(string),
				LeaseTime:           poolMap["lease_time"].(string),
				SubnetMask:          poolMap["subnet_mask"].(string),
				IpRange:             poolMap["ip_range"].(string),
				PrimaryNameServer:   poolMap["dns1"].(string),
				SecondaryNameServer: poolMap["dns2"].(string),
			}
			dhcpPools[index] = onePool
		}

		// Inject data into main structure
		edgeDhcp.EdgeDhcpIpPools = &types.EdgeDhcpIpPools{}
		edgeDhcp.EdgeDhcpIpPools.EdgeDhcpIpPool = dhcpPools
	}

	return &edgeDhcp

}
