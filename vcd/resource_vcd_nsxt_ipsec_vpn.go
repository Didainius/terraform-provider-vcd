package vcd

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

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
				Description: "Security type which is use for IPsec VPN Tunnel. It will be 'DEFAULT' if nothing is " +
					"customized and 'CUSTOM' if some changes are applied",
			},
			"status": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Overall IPsec VPN Tunnel Status",
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

			"security_profile_customization": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "Security profile customization",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{

						"ike_version": &schema.Schema{
							Type:         schema.TypeString,
							Required:     true,
							Description:  "",
							ValidateFunc: validation.StringInSlice([]string{"IKE_V1", "IKE_V2", "IKE_FLEX"}, false),
						},
						"ike_encryption_algorithms": &schema.Schema{
							Type:        schema.TypeSet,
							Required:    true,
							Description: "",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"ike_digest_algorithms": &schema.Schema{
							Type:        schema.TypeSet,
							Optional:    true,
							Description: "",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"ike_dh_groups": &schema.Schema{
							Type:        schema.TypeSet,
							Required:    true,
							Description: "",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"ike_sa_lifetime": &schema.Schema{
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Security Association life time (in seconds)",
						},

						"tunnel_pfs_enabled": &schema.Schema{
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     true,
							Description: "Perfect Forward Secrecy",
						},

						"tunnel_df_policy": &schema.Schema{
							Type:         schema.TypeString,
							Optional:     true,
							Default:      "COPY",
							Description:  "Perfect Forward Secrecy",
							ValidateFunc: validation.StringInSlice([]string{"COPY", "CLEAR"}, false),
						},

						"tunnel_encryption_algorithms": &schema.Schema{
							Type:        schema.TypeSet,
							Required:    true,
							Description: "",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"tunnel_digest_algorithms": &schema.Schema{
							Type:        schema.TypeSet,
							Optional:    true,
							Description: "",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"tunnel_dh_groups": &schema.Schema{
							Type:        schema.TypeSet,
							Required:    true,
							Description: "",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"tunnel_sa_lifetime": &schema.Schema{
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Security Association life time (in seconds)",
						},
						"dpd_probe_internal": &schema.Schema{
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Dead Peer Detection probe interval (in seconds)",
						},
					},
				},
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
		return diag.Errorf("error getting NSX-T IPsec VPN Tunnel configuration type: %s", err)
	}

	createdIpSecVpnConfig, err := nsxtEdge.CreateIpSecVpn(ipSecVpnConfig)
	if err != nil {
		return diag.Errorf("error creating NSX-T IPsec VPN Tunnel configuration: %s", err)
	}
	// IPSec VPN Tunnel is already create - store the ID
	d.SetId(createdIpSecVpnConfig.NsxtIpSecVpn.ID)

	// Tunnel Profile
	if _, isSet := d.GetOk("security_profile_customization"); isSet {
		tunnelProfileConfig, err := getNsxtIpSecVpnProfileTunnelConfigurationType(d)
		if err != nil {
			return diag.Errorf("error getting NSX-T IPsec VPN Tunnel Profile: %s", err)
		}

		_, err = createdIpSecVpnConfig.UpdateTunnelConnectionProperties(tunnelProfileConfig)
		if err != nil {
			return diag.Errorf("error setting VPN Tunnel Profile: %s", err)
		}

	}

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

	// Retrieve configuration
	existingIpSecVpnConfiguration, err := nsxtEdge.GetIpSecVpnById(d.Id())
	if err != nil {
		diag.Errorf("error retrieving existing NSX-T IPsec VPN Tunnel configuration: %s", err)
	}

	ipSecVpnConfig, err := getNsxtIpSecVpnTunnelType(d)
	if err != nil {
		return diag.Errorf("error getting NSX-T IPsec VPN Tunnel configuration type: %s", err)
	}
	// Inject ID for update
	ipSecVpnConfig.ID = d.Id()

	//
	securityProfileHasChange := d.HasChange("security_profile_customization")
	_, newSecurityProfile := d.GetChange("security_profile_customization")

	// Security Profile Customization settings work on two different endpoints:
	// * To set a custom security profile - there is a separate endpoint where all security profile settings can be
	// set. After setting them, parent IPsec VPN Tunnel `SecurityType` becomes "CUSTOM".
	// * To remove customization and switch back to NSX-T Default parameters the parent IPsec VPN Tunnel must be updated
	// and its field 'SecurityType' must be set to 'DEFAULT'
	if securityProfileHasChange && len(newSecurityProfile.([]interface{})) == 0 {
		ipSecVpnConfig.SecurityType = "DEFAULT"
	}

	// At first update IPsec VPN tunnel configuration
	// It will reset Security Profile to DEFAULT at the same shot if no customization exists in 'security_profile_customization'
	updatedIpSecVpnConfiguration, err := existingIpSecVpnConfiguration.Update(ipSecVpnConfig)
	if err != nil {
		return diag.Errorf("error updating NSX-T IPsec VPN Tunnel configuration '%s': %s", ipSecVpnConfig.Name, err)
	}

	// If Security Profile has change - it must be explicitly set using UpdateTunnelConnectionProperties which will
	// change value of parent NsxtIpSecVpnTunnel.SecurityType to 'CUSTOM' automatically.
	if securityProfileHasChange && newSecurityProfile != nil {
		ipSecTunnelProfileConfig, err := getNsxtIpSecVpnProfileTunnelConfigurationType(d)
		if err != nil {
			return diag.Errorf("error getting NSX-T IPsec VPN Tunnel Profile: %s", err)
		}

		// To set IPsec VPN Tunnel Connection Profile - it must be updated (HTTP PUT) with all the options configured
		if ipSecTunnelProfileConfig != nil {
			_, err = updatedIpSecVpnConfiguration.UpdateTunnelConnectionProperties(ipSecTunnelProfileConfig)
			if err != nil {
				return diag.Errorf("error setting NSX-T IPsec VPN Tunnel Profile: %s", err)
			}

		}
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
		return diag.Errorf("error retrieving NSX-T IPsec VPN Tunnel configuration for deletion: %s", err)
	}

	// Set general schema for configuration
	err = setNsxtIpSecVpnTunnelData(d, ipSecVpnConfig.NsxtIpSecVpn)
	if err != nil {
		return diag.Errorf("error storing NSX-T IPsec VPN Tunnel configuration to schema: %s", err)
	}

	// ipSe
	tunnelConnectionProperties, err := ipSecVpnConfig.GetTunnelConnectionProperties()
	if err != nil {
		return diag.Errorf("error reading NSX-T IPsec VPN Tunnel Security Customization: %s", err)
	}

	err = setNsxtIpSecVpnProfileTunnelConfigurationData(d, tunnelConnectionProperties)
	if err != nil {
		return diag.Errorf("error storing NSX-T IPsec VPN Tunnel Security Customization to schema: %s", err)
	}

	// Read tunnel status data from separate endpoint
	tunnelStatus, err := ipSecVpnConfig.GetStatus()
	if err != nil {
		return diag.Errorf("error reading NSX-T IPsec VPN Tunnel status: %s", err)
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
		return diag.Errorf("error retrieving NSX-T IPsec VPN Tunnel configuration for deletion: %s", err)
	}

	err = ipSecVpnConfig.Delete()
	if err != nil {
		return diag.Errorf("error deleting NSX-T IPsec VPN Tunnel configuration: %s", err)
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

func getNsxtIpSecVpnProfileTunnelConfigurationType(d *schema.ResourceData) (*types.NsxtIpSecVpnTunnelSecurityProfile, error) {
	tunnel, isSet := d.GetOk("security_profile_customization")

	if !isSet {
		return nil, nil
	}
	tunnelSlice := tunnel.([]interface{})
	tunnelMap := tunnelSlice[0].(map[string]interface{})

	nsxtIpSecVpnTunnelProfile := &types.NsxtIpSecVpnTunnelSecurityProfile{
		SecurityType: "CUSTOM", // Security Type must become CUSTOM, because we are configuring profile
		IkeConfiguration: types.NsxtIpSecVpnTunnelProfileIkeConfiguration{
			IkeVersion:           tunnelMap["ike_version"].(string),
			EncryptionAlgorithms: convertSchemaSetToSliceOfStrings(tunnelMap["ike_encryption_algorithms"].(*schema.Set)),
			DigestAlgorithms:     convertSchemaSetToSliceOfStrings(tunnelMap["ike_digest_algorithms"].(*schema.Set)),
			DhGroups:             convertSchemaSetToSliceOfStrings(tunnelMap["ike_dh_groups"].(*schema.Set)),
			SaLifeTime:           tunnelMap["ike_sa_lifetime"].(int),
		},
		TunnelConfiguration: types.NsxtIpSecVpnTunnelProfileTunnelConfiguration{
			PerfectForwardSecrecyEnabled: tunnelMap["tunnel_pfs_enabled"].(bool),
			DfPolicy:                     tunnelMap["tunnel_df_policy"].(string),
			EncryptionAlgorithms:         convertSchemaSetToSliceOfStrings(tunnelMap["tunnel_encryption_algorithms"].(*schema.Set)),
			DigestAlgorithms:             convertSchemaSetToSliceOfStrings(tunnelMap["tunnel_digest_algorithms"].(*schema.Set)),
			DhGroups:                     convertSchemaSetToSliceOfStrings(tunnelMap["tunnel_dh_groups"].(*schema.Set)),
			SaLifeTime:                   tunnelMap["tunnel_sa_lifetime"].(int),
		},
		DpdConfiguration: types.NsxtIpSecVpnTunnelProfileDpdConfiguration{
			ProbeInterval: tunnelMap["dpd_probe_internal"].(int),
		},
	}

	return nsxtIpSecVpnTunnelProfile, nil
}

func converListTotTypeSet(slice []string) *schema.Set {
	sliceOfInterfaces := convertToTypeSet(slice)
	set := schema.NewSet(schema.HashSchema(&schema.Schema{Type: schema.TypeString}), sliceOfInterfaces)

	return set
}

func setNsxtIpSecVpnProfileTunnelConfigurationData(d *schema.ResourceData, tunneConfig *types.NsxtIpSecVpnTunnelSecurityProfile) error {

	if tunneConfig.SecurityType == "DEFAULT" {
		err := d.Set("security_profile_customization", nil)
		if err != nil {
			return fmt.Errorf("error resetting 'security_profile_customization' to empty: %s", err)
		}
		// Return early because there is nothing to store
		return nil
	}

	object := make([]interface{}, 1)
	objectMap := make(map[string]interface{})

	objectMap["ike_version"] = tunneConfig.IkeConfiguration.IkeVersion

	objectMap["ike_encryption_algorithms"] = converListTotTypeSet(tunneConfig.IkeConfiguration.EncryptionAlgorithms)
	objectMap["ike_digest_algorithms"] = converListTotTypeSet(tunneConfig.IkeConfiguration.DigestAlgorithms)
	objectMap["ike_dh_groups"] = converListTotTypeSet(tunneConfig.IkeConfiguration.DhGroups)
	objectMap["ike_sa_lifetime"] = tunneConfig.IkeConfiguration.SaLifeTime

	objectMap["tunnel_pfs_enabled"] = tunneConfig.TunnelConfiguration.PerfectForwardSecrecyEnabled
	objectMap["tunnel_df_policy"] = tunneConfig.TunnelConfiguration.DfPolicy
	objectMap["tunnel_encryption_algorithms"] = converListTotTypeSet(tunneConfig.TunnelConfiguration.EncryptionAlgorithms)
	objectMap["tunnel_digest_algorithms"] = converListTotTypeSet(tunneConfig.TunnelConfiguration.DigestAlgorithms)
	objectMap["tunnel_dh_groups"] = converListTotTypeSet(tunneConfig.TunnelConfiguration.DhGroups)
	objectMap["tunnel_sa_lifetime"] = tunneConfig.TunnelConfiguration.SaLifeTime

	objectMap["dpd_probe_internal"] = tunneConfig.DpdConfiguration.ProbeInterval

	object[0] = objectMap

	err := d.Set("security_profile_customization", object)

	return err
}
