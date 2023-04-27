package vcd

import (
	"context"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/vmware/go-vcloud-director/v2/govcd"
	"github.com/vmware/go-vcloud-director/v2/types/v56"
)

/// TODO
// IP space uplinks in Provider Gateways (separate API call)
// External network (provider gateway "ownership") (`dedicatedOrg` field in external network endpoint)

// var ipSpaceIpRange = &schema.Resource{
// 	Schema: map[string]*schema.Schema{
// 		"ip_ranges": {
// 			Type:        schema.TypeSet,
// 			Optional:    true,
// 			Description: "IP ranges (should match internal scope)",
// 			Elem:        ipSpaceIpRangeRange,
// 		},
// 		"default_quota": {
// 			Type:        schema.TypeString,
// 			Required:    true,
// 			Description: "Floating IP quota (-1 for unlimited, 0 - cannot be allocated)",
// 		},
// 	},
// }

var ipSpaceIpRangeRange = &schema.Resource{
	Schema: map[string]*schema.Schema{
		"start_address": {
			Type:         schema.TypeString,
			Required:     true,
			Description:  "Start address of the IP range",
			ValidateFunc: validation.IsIPAddress,
		},
		"end_address": {
			Type:         schema.TypeString,
			Required:     true,
			Description:  "End address of the IP range",
			ValidateFunc: validation.IsIPAddress,
		},
	},
}

var ipPrefixes = &schema.Resource{
	Schema: map[string]*schema.Schema{
		"prefix": {
			Type:        schema.TypeSet,
			Optional:    true,
			Description: "IP ranges (should match internal scope)",
			Elem:        ipSpacePrefix,
		},
		"default_quota": {
			Type:         schema.TypeString,
			Required:     true,
			Description:  "Floating IP quota",
			ValidateFunc: IsIntAndAtLeast(-1),
		},
	},
}

var ipSpacePrefix = &schema.Resource{
	Schema: map[string]*schema.Schema{
		"first_ip": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "First IP in CIDR format",
			// ValidateFunc: validation.IsIPAddress,
		},
		"prefix_length": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "First IP in CIDR format",
			// ValidateFunc: validation.IsIPAddress,
		},
		"prefix_count": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Prefix count",
			// ValidateFunc: validation.IsIPAddress,
		},
	},
}

// Quota is set for:
// * Floating IPs
// * Each of the defined prefixes
// var ipSpaceQuota = &schema.Resource{
// 	Schema: map[string]*schema.Schema{
// 		"floating_ips": {
// 			Type:        schema.TypeString,
// 			Required:    true,
// 			Description: "Floating IPs",
// 			// ValidateFunc: validation.IsIPAddress,
// 		},
// 		"32 prefix": {
// 			Type:        schema.TypeString,
// 			Required:    true,
// 			Description: "/32 prefix",
// 			// ValidateFunc: validation.IsIPAddress,
// 		},
// 	},
// }

func resourceVcdIpSpace() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVcdIpSpaceCreate,
		ReadContext:   resourceVcdIpSpaceRead,
		UpdateContext: resourceVcdIpSpaceUpdate,
		DeleteContext: resourceVcdIpSpaceDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceVcdIpSpaceImport,
		},

		Schema: map[string]*schema.Schema{
			"org_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "For 'SHARED' (Org bound) IP spaces - Org ID",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of IP space",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Description of IP space",
			},

			"type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Type of IP space",
				// PUBLIC, SHARED_SERVICES, PRIVATE
			},

			// metadata? Nothing in UI, but maybe we can leverage the new JSON metadata mechanism?

			// "org": {
			// 	Type:     schema.TypeString,
			// 	Optional: true,
			// 	ForceNew: true,
			// 	Description: "The name of organization to use, optional if defined at provider " +
			// 		"level. Useful when connected as sysadmin working across different organizations",
			// },
			"internal_scope": {
				Type:        schema.TypeSet,
				Required:    true,
				Description: "A set of up internal scope IPs in CIDR format",
				Elem: &schema.Schema{
					MinItems: 1,
					Type:     schema.TypeString,
				},
			},

			"ip_range_quota": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "IP ranges (should match internal scope)",
				ValidateFunc: IsIntAndAtLeast(-1),
			},
			"ip_ranges": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "IP ranges (should match internal scope)",
				Elem:        ipSpaceIpRangeRange,
			},

			"ip_prefixes": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "IP prefixes (should match internal scope)",
				Elem:        ipPrefixes,
			},
			"external_scope": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "External scope in CIDR format",
			},

			"route_advertisement_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Flag whether route advertisement should be enabled",
			},

			// "quota": {
			// 	Type:        schema.TypeList,
			// 	MaxItems:    1,
			// 	Optional:    true,
			// 	Description: "Quota",
			// 	Elem:        ipPrefixes,
			// },
		},
	}
}

func resourceVcdIpSpaceCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vcdClient := meta.(*VCDClient)
	log.Printf("[TRACE] IP Space creation initiated")

	ipSpaceConfig, err := getIpSpaceType(d)
	if err != nil {
		return diag.Errorf("could not get IP Space type: %s", err)
	}

	createdIpSpace, err := vcdClient.GenericCreateIpSpace(ipSpaceConfig)
	if err != nil {
		return diag.Errorf("error creating IP Space: %s", err)
	}

	d.SetId(createdIpSpace.IpSpace.ID)

	return resourceVcdIpSpaceRead(ctx, d, meta)
}

func resourceVcdIpSpaceUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return resourceVcdIpSpaceRead(ctx, d, meta)
}

func resourceVcdIpSpaceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vcdClient := meta.(*VCDClient)
	log.Printf("[TRACE] IP Space read initiated")

	ipSpace, err := vcdClient.GetIpSpaceById(d.Id())
	if err != nil {
		if govcd.ContainsNotFound(err) {
			d.SetId("")
			return nil
		}
		return diag.Errorf("error finding IP Space by ID '%s': %s", d.Id(), err)
	}

	err = setIpSpaceData(d, ipSpace.IpSpace)
	if err != nil {
		return diag.Errorf("error storing IP Space state: %s", err)
	}

	return nil
}

func resourceVcdIpSpaceDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vcdClient := meta.(*VCDClient)
	log.Printf("[TRACE] IP Space deletion initiated")

	ipSpace, err := vcdClient.GetIpSpaceById(d.Id())
	if err != nil {
		return diag.Errorf("error finding IP Space by ID '%s': %s", d.Id(), err)
	}

	err = ipSpace.Delete()
	if err != nil {
		return diag.Errorf("error deleting IP space by ID '%s': %s", d.Id(), err)
	}

	return nil
}

func resourceVcdIpSpaceImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	return []*schema.ResourceData{d}, nil
}

func getIpSpaceType(d *schema.ResourceData) (*types.IpSpace, error) {

	ipSpace := &types.IpSpace{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Type:        d.Get("type").(string),

		// Utilization: types.Utilization{
		// 	FloatingIPs: types.FloatingIPs{
		// 		TotalCount:          "",
		// 		AllocatedCount:      "",
		// 		UsedCount:           "",
		// 		UnusedCount:         "",
		// 		AllocatedPercentage: 0,
		// 		UsedPercentage:      0,
		// 	},
		// 	IPPrefixes: types.IPPrefixes{
		// 		TotalCount:               "",
		// 		AllocatedCount:           "",
		// 		UsedCount:                "",
		// 		UnusedCount:              "",
		// 		AllocatedPercentage:      0,
		// 		UsedPercentage:           0,
		// 		PrefixLengthUtilizations: []types.PrefixLengthUtilizations{},
		// 	},
		// },
		// IPSpaceRanges: types.IPSpaceRanges{
		// 	IPRanges: []types.IpSpaceRangeValues{
		// 		{
		// 			ID:             "",
		// 			StartIPAddress: "",
		// 			EndIPAddress:   "",
		// 			// TotalIPCount:          "",
		// 			// AllocatedIPCount:      "",
		// 			// AllocatedIPPercentage: 0,
		// 		},
		// 	},
		// 	DefaultFloatingIPQuota: 0,
		// },
		// IPSpacePrefixes: []types.IPSpacePrefixes{
		// 	{
		// 		IPPrefixSequence: []types.IPPrefixSequence{
		// 			{
		// 				// ID:                        "",
		// 				StartingPrefixIPAddress: "",
		// 				PrefixLength:            0,
		// 				TotalPrefixCount:        0,

		// 				// AllocatedPrefixCount:      0,
		// 				// AllocatedPrefixPercentage: 0,
		// 			},
		// 		},
		// 		DefaultQuotaForPrefixLength: 0,
		// 	},
		// },

		IPSpaceInternalScope:      convertSchemaSetToSliceOfStrings(d.Get("internal_scope").(*schema.Set)),
		IPSpaceExternalScope:      d.Get("external_scope").(string),
		RouteAdvertisementEnabled: d.Get("route_advertisement_enabled").(bool),
	}

	// IP Space Ranges
	ipRangeQuota := d.Get("ip_range_quota").(string)
	ipRangeQuotaInt, _ := strconv.Atoi(ipRangeQuota) // error is ignored because validation is enforced in schema

	ipSpace.IPSpaceRanges = types.IPSpaceRanges{
		DefaultFloatingIPQuota: ipRangeQuotaInt,
	}

	ipRanges := d.Get("ip_ranges").(*schema.Set)
	ipRangesSlice := ipRanges.List()
	if len(ipRangesSlice) > 0 {

		ipSpace.IPSpaceRanges.IPRanges = make([]types.IpSpaceRangeValues, len(ipRangesSlice))
		for ipRangeIndex := range ipRangesSlice {
			ipRangeStrings := convertToStringMap(ipRangesSlice[ipRangeIndex].(map[string]interface{}))

			ipSpace.IPSpaceRanges.IPRanges[ipRangeIndex].StartIPAddress = ipRangeStrings["start_address"]
			ipSpace.IPSpaceRanges.IPRanges[ipRangeIndex].EndIPAddress = ipRangeStrings["end_address"]

		}
	}

	// EOF // IP Space Ranges

	// IP Prefixes
	ipPrefixes := d.Get("ip_prefixes").(*schema.Set)
	ipPrefixesSlice := ipPrefixes.List()

	// Initialize structure
	if len(ipPrefixesSlice) > 0 {
		ipSpace.IPSpacePrefixes = []types.IPSpacePrefixes{}
	}

	for ipPrefixIndex := range ipPrefixesSlice {

		ipppppppPrefix := ipPrefixesSlice[ipPrefixIndex]
		ipPrefixMap := ipppppppPrefix.(map[string]interface{})
		ipPrefixQuota := ipPrefixMap["default_quota"].(string)
		ipPrefixQuotaInt, _ := strconv.Atoi(ipPrefixQuota) // ignoring error as validation is enforce in schema

		ipSpacePrexif := types.IPSpacePrefixes{
			DefaultQuotaForPrefixLength: ipPrefixQuotaInt,
		}

		// Extract IP prefixess

		// 'prefix'
		ipPrefixPrefix := ipPrefixMap["prefix"].(*schema.Set)
		ipPrefixPrefixSlice := ipPrefixPrefix.List()
		if len(ipPrefixPrefixSlice) > 0 {
			ipSpacePrexif.IPPrefixSequence = []types.IPPrefixSequence{}
		}

		for ipPrefixPrefixIndex := range ipPrefixPrefixSlice {

			ipPrefixMap := convertToStringMap(ipPrefixPrefixSlice[ipPrefixPrefixIndex].(map[string]interface{}))
			prefixLengthInt, _ := strconv.Atoi(ipPrefixMap["prefix_length"])
			prefixLengthCountInt, _ := strconv.Atoi(ipPrefixMap["prefix_count"])

			ipSpacePrexif.IPPrefixSequence = append(ipSpacePrexif.IPPrefixSequence, types.IPPrefixSequence{
				StartingPrefixIPAddress: ipPrefixMap["first_ip"],
				PrefixLength:            prefixLengthInt,
				TotalPrefixCount:        prefixLengthCountInt,
			})

			// "first_ip": {
			// 	Type:        schema.TypeString,
			// 	Required:    true,
			// 	Description: "First IP in CIDR format",
			// 	// ValidateFunc: validation.IsIPAddress,
			// },
			// "prefix_length": {
			// 	Type:        schema.TypeString,
			// 	Required:    true,
			// 	Description: "First IP in CIDR format",
			// 	// ValidateFunc: validation.IsIPAddress,
			// },
			// "prefix_count": {
			// 	Type:        schema.TypeString,
			// 	Required:    true,
			// 	Description: "Prefix count",
			// 	// ValidateFunc: validation.IsIPAddress,
			// },

		}

		// EOF // Extract IP prefixess

		// Add to the list
		ipSpace.IPSpacePrefixes = append(ipSpace.IPSpacePrefixes, ipSpacePrexif)

	}

	// EOF IP Prefixes

	// only with
	orgId := d.Get("org_id").(string)
	if orgId != "" {
		ipSpace.OrgRef = &types.OpenApiReference{ID: orgId}
	}

	return ipSpace, nil
}

func setIpSpaceData(d *schema.ResourceData, ipSpace *types.IpSpace) error {

	dSet(d, "name", ipSpace.Name)
	dSet(d, "description", ipSpace.Description)
	dSet(d, "type", ipSpace.Type)
	dSet(d, "route_advertisement_enabled", ipSpace.RouteAdvertisementEnabled)

	return nil
}
