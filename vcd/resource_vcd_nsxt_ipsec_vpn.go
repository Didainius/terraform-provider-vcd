package vcd

import (
	"context"
	"fmt"

	"github.com/vmware/go-vcloud-director/v2/govcd"

	"github.com/vmware/go-vcloud-director/v2/types/v56"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceVcdNsxtIpSecVpnTunnel() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVcdNsxtIpSecVpnTunnelCreate,
		ReadContext:   resourceVcdNsxtIpSecVpnTunnelRead,
		UpdateContext: resourceVcdNsxtIpSecVpnTunnelUpdate,
		DeleteContext: resourceVcdNsxtIpSecVpnTunnelDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceVcdNsxtIpSecVpnTunnelImport,
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
				Description: "Edge gateway name in which IP Sec VPN configuration is located",
			},
			"enabled": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Enables or disables this configuration (default true)",
			},
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of IP Sec VPN configuration",
			},
			"description": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Description of NAT rule",
			},
			"pre_shared_key": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				Description: "Pre-Shared Key (PSK)",
			},
			"local_ip_address": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "IPv4 Address for the endpoint. This has to be a suballocated IP on the Edge Gateway.",
			},
			"local_networks": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Set of local networks in CIDR format. Leaving it empty is interpreted as 0.0.0.0/0",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"remote_ip_address": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Public IPv4 Address of the remote device terminating the VPN connection",
			},
			"remote_networks": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Set of remote networks in CIDR format. Leaving it empty is interpreted as 0.0.0.0/0",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"logging": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Sets whether logging for the tunnel is enabled or not. (default - false)",
			},

			"security_profile": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				Description: "Security type which is use for IPSec VPN Tunnel. It will be 'DEFAULT' if nothing is " +
					"customized and 'CUSTOM' if some changes are applied",
			},
			"status": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Overall IPSec VPN Tunnel Status",
			},
			"ike_service_status": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Status for the actual IKE Session for the given tunnel",
			},
			"ike_fail_reason": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Provides more details of failure if the IKE service is not UP",
			},
		},
	}
}

func resourceVcdNsxtIpSecVpnTunnelCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vcdClient := meta.(*VCDClient)
	vcdClient.lockParentEdgeGtw(d)
	defer vcdClient.unLockParentEdgeGtw(d)

	orgName := d.Get("org").(string)
	vdcName := d.Get("vdc").(string)
	edgeGatewayId := d.Get("edge_gateway_id").(string)

	nsxtEdge, err := vcdClient.GetNsxtEdgeGatewayById(orgName, vdcName, edgeGatewayId)
	if err != nil {
		return diag.Errorf("error retrieving Edge Gateway: %s", err)
	}

	ipSecVpnConfig, err := getNsxtIpSecVpnTunnelType(d)
	if err != nil {
		return diag.Errorf("error getting NSX-T IPSec VPN Tunnel configuration type: %s", err)
	}

	createdIpSecVpnConfig, err := nsxtEdge.CreateIpSecVpn(ipSecVpnConfig)
	if err != nil {
		return diag.Errorf("error creating NSX-T IPSec VPN Tunnel configuration: %s", err)
	}

	d.SetId(createdIpSecVpnConfig.NsxtIpSecVpn.ID)

	return resourceVcdNsxtIpSecVpnTunnelRead(ctx, d, meta)
}

func resourceVcdNsxtIpSecVpnTunnelUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vcdClient := meta.(*VCDClient)
	vcdClient.lockParentEdgeGtw(d)
	defer vcdClient.unLockParentEdgeGtw(d)

	orgName := d.Get("org").(string)
	vdcName := d.Get("vdc").(string)
	edgeGatewayId := d.Get("edge_gateway_id").(string)

	nsxtEdge, err := vcdClient.GetNsxtEdgeGatewayById(orgName, vdcName, edgeGatewayId)
	if err != nil {
		return diag.Errorf("error retrieving Edge Gateway: %s", err)
	}

	ipSecVpnConfig, err := getNsxtIpSecVpnTunnelType(d)
	if err != nil {
		return diag.Errorf("error getting NSX-T IPSec VPN Tunnel configuration type: %s", err)
	}
	// Inject ID for update
	ipSecVpnConfig.ID = d.Id()

	_, err = nsxtEdge.CreateIpSecVpn(ipSecVpnConfig)
	if err != nil {
		return diag.Errorf("error updating NSX-T IPSec VPN Tunnel configuration '%s': %s", ipSecVpnConfig.Name, err)
	}

	return resourceVcdNsxtIpSecVpnTunnelRead(ctx, d, meta)
}

func resourceVcdNsxtIpSecVpnTunnelRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vcdClient := meta.(*VCDClient)

	orgName := d.Get("org").(string)
	vdcName := d.Get("vdc").(string)
	edgeGatewayId := d.Get("edge_gateway_id").(string)

	nsxtEdge, err := vcdClient.GetNsxtEdgeGatewayById(orgName, vdcName, edgeGatewayId)
	if err != nil {
		return diag.Errorf("error retrieving Edge Gateway: %s", err)
	}

	ipSecVpnConfig, err := nsxtEdge.GetIpSecVpnById(d.Id())
	if err != nil {
		if govcd.ContainsNotFound(err) {
			d.SetId("")
		}
		return diag.Errorf("error retrieving NSX-T IPSec VPN Tunnel configuration for deletion: %s", err)
	}

	// Set general schema for configuration
	err = setNsxtIpSecVpnTunnelData(d, ipSecVpnConfig.NsxtIpSecVpn)
	if err != nil {
		return diag.Errorf("error storing NSX-T IPSec VPN Tunnel configuration to schema: %s", err)
	}

	// Read tunnel status data from separate endpoint
	tunnelStatus, err := ipSecVpnConfig.GetStatus()
	if err != nil {
		return diag.Errorf("error reading NSX-T IPSec VPN Tunnel status: %s", err)
	}
	setNsxtIpSecVpnTunnelStatusData(d, tunnelStatus)

	return nil
}

func resourceVcdNsxtIpSecVpnTunnelDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vcdClient := meta.(*VCDClient)
	vcdClient.lockParentEdgeGtw(d)
	defer vcdClient.unLockParentEdgeGtw(d)

	orgName := d.Get("org").(string)
	vdcName := d.Get("vdc").(string)
	edgeGatewayId := d.Get("edge_gateway_id").(string)

	nsxtEdge, err := vcdClient.GetNsxtEdgeGatewayById(orgName, vdcName, edgeGatewayId)
	if err != nil {
		return diag.Errorf("error retrieving Edge Gateway: %s", err)
	}

	ipSecVpnConfig, err := nsxtEdge.GetIpSecVpnById(d.Id())
	if err != nil {
		return diag.Errorf("error retrieving NSX-T IPSec VPN Tunnel configuration for deletion: %s", err)
	}

	err = ipSecVpnConfig.Delete()
	if err != nil {
		return diag.Errorf("error deleting NSX-T IPSec VPN Tunnel configuration: %s", err)
	}

	d.SetId("")

	return nil
}

func resourceVcdNsxtIpSecVpnTunnelImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	return []*schema.ResourceData{d}, nil
}

func getNsxtIpSecVpnTunnelType(d *schema.ResourceData) (*types.NsxtIpSecVpnTunnel, error) {
	ipSecVpnConfig := &types.NsxtIpSecVpnTunnel{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Enabled:     d.Get("enabled").(bool),
		LocalEndpoint: types.NsxtIpSecVpnTunnelLocalEndpoint{
			LocalId:       d.Get("local_ip_address").(string),
			LocalAddress:  d.Get("local_ip_address").(string),
			LocalNetworks: convertSchemaSetToSliceOfStrings(d.Get("local_networks").(*schema.Set)),
		},
		RemoteEndpoint: types.NsxtIpSecVpnTunnelRemoteEndpoint{
			RemoteId:       d.Get("remote_ip_address").(string),
			RemoteAddress:  d.Get("remote_ip_address").(string),
			RemoteNetworks: convertSchemaSetToSliceOfStrings(d.Get("remote_networks").(*schema.Set)),
		},
		PreSharedKey: d.Get("pre_shared_key").(string),
		//SecurityType:            "",
		Logging: d.Get("logging").(bool),
		//AuthenticationMode:      "",
		//ConnectorInitiationMode: "",
		//Version:                 nil,
	}

	return ipSecVpnConfig, nil
}

func setNsxtIpSecVpnTunnelData(d *schema.ResourceData, ipSecVpnConfig *types.NsxtIpSecVpnTunnel) error {
	_ = d.Set("name", ipSecVpnConfig.Name)
	_ = d.Set("description", ipSecVpnConfig.Description)
	_ = d.Set("pre_shared_key", ipSecVpnConfig.PreSharedKey)
	_ = d.Set("enabled", ipSecVpnConfig.Enabled)
	_ = d.Set("local_ip_address", ipSecVpnConfig.LocalEndpoint.LocalAddress)
	_ = d.Set("enabled", ipSecVpnConfig.Enabled)
	_ = d.Set("logging", ipSecVpnConfig.Logging)
	_ = d.Set("security_profile", ipSecVpnConfig.SecurityType)

	localNetworks := convertToTypeSet(ipSecVpnConfig.LocalEndpoint.LocalNetworks)
	localNetworksSet := schema.NewSet(schema.HashSchema(&schema.Schema{Type: schema.TypeString}), localNetworks)
	err := d.Set("local_networks", localNetworksSet)
	if err != nil {
		return fmt.Errorf("error storing 'local_networks': %s", err)
	}

	_ = d.Set("remote_ip_address", ipSecVpnConfig.RemoteEndpoint.RemoteAddress)
	remoteNetworks := convertToTypeSet(ipSecVpnConfig.RemoteEndpoint.RemoteNetworks)
	remoteNetworksSet := schema.NewSet(schema.HashSchema(&schema.Schema{Type: schema.TypeString}), remoteNetworks)
	err = d.Set("remote_networks", remoteNetworksSet)
	if err != nil {
		return fmt.Errorf("error storing 'remote_networks': %s", err)
	}

	return nil
}

func setNsxtIpSecVpnTunnelStatusData(d *schema.ResourceData, ipSecVpnStatus *types.NsxtIpSecVpnTunnelStatus) {
	_ = d.Set("status", ipSecVpnStatus.TunnelStatus)
	_ = d.Set("ike_service_status", ipSecVpnStatus.IkeStatus.IkeServiceStatus)
	_ = d.Set("ike_fail_reason", ipSecVpnStatus.IkeStatus.FailReason)
	return
}
