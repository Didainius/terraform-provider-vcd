package vcd

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vmware/go-vcloud-director/v2/govcd"
	"github.com/vmware/go-vcloud-director/v2/types/v56"
	"os"
)

func datasourceVcdNsxtIpSet() *schema.Resource {
	return &schema.Resource{
		ReadContext: datasourceVcdNsxtIpSetRead,

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
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "IP set name",
			},
			"edge_gateway_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Edge Gateway ID in which IP Set is located",
			},
			"description": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "IP set description",
			},
			"ip_addresses": {
				Computed:    true,
				Type:        schema.TypeSet,
				Description: "A set of IP address, CIDR, IP range objects",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func datasourceVcdNsxtIpSetRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vcdClient := meta.(*VCDClient)

	nsxtEdgeGateway, err := vcdClient.GetNsxtEdgeGatewayFromResourceById(d, "edge_gateway_id")
	if err != nil {
		return diag.Errorf(errorUnableToFindEdgeGateway, err)
	}

	// Name uniqueness is enforced by VCD for types.FirewallGroupTypeIpSet
	ipSet, err := nsxtEdgeGateway.GetNsxtFirewallGroupByName(d.Get("name").(string), types.FirewallGroupTypeIpSet)
	if err != nil {
		return diag.Errorf("error getting NSX-T IP Set with Name '%s': %s", d.Get("name").(string), err)
	}
	err = setNsxtIpSetData(d, ipSet.NsxtFirewallGroup)
	if err != nil {
		return diag.Errorf("error setting NSX-T IP Set: %s", err)
	}

	d.SetId(ipSet.NsxtFirewallGroup.ID)

	return nil
}

// ===================== To be removed before merge (will come from Security Group PR) ===========

// GetNsxtEdgeGatewayFromResource helps to retrieve NSX-T Edge Gateway when Org org VDC are not
// needed. It performs a query By ID.
func (cli *VCDClient) GetNsxtEdgeGatewayFromResourceById(d *schema.ResourceData, edgeGatewayFieldName string) (eg *govcd.NsxtEdgeGateway, err error) {
	orgName := d.Get("org").(string)
	vdcName := d.Get("vdc").(string)
	edgeGatewayId := d.Get(edgeGatewayFieldName).(string)
	egw, err := cli.GetNsxtEdgeGatewayById(orgName, vdcName, edgeGatewayId)
	if err != nil {
		if os.Getenv("GOVCD_DEBUG") != "" {
			return nil, fmt.Errorf("(%s) [%s] : %s", edgeGatewayId, callFuncName(), err)
		}
		return nil, err
	}
	return egw, nil
}

// GetNsxtEdgeGatewayById gets an NSX-T Edge Gateway when you don't need Org or VDC for other purposes
func (cli *VCDClient) GetNsxtEdgeGatewayById(orgName, vdcName, edgeGwId string) (eg *govcd.NsxtEdgeGateway, err error) {
	if edgeGwId == "" {
		return nil, fmt.Errorf("empty NSX-T Edge Gateway ID provided")
	}
	_, vdc, err := cli.GetOrgAndVdc(orgName, vdcName)
	if err != nil {
		return nil, fmt.Errorf("error retrieving Org and VDC: %s", err)
	}
	eg, err = vdc.GetNsxtEdgeGatewayById(edgeGwId)

	if err != nil {
		if os.Getenv("GOVCD_DEBUG") != "" {
			return nil, fmt.Errorf(fmt.Sprintf("(%s) [%s] ", edgeGwId, callFuncName())+errorUnableToFindEdgeGateway, err)
		}
		return nil, fmt.Errorf(errorUnableToFindEdgeGateway, err)
	}
	return eg, nil
}

// GetNsxtEdgeGateway gets an NSX-T Edge Gateway when you don't need Org or VDC for other purposes
func (cli *VCDClient) GetNsxtEdgeGateway(orgName, vdcName, edgeGwName string) (eg *govcd.NsxtEdgeGateway, err error) {

	if edgeGwName == "" {
		return nil, fmt.Errorf("empty NSX-T Edge Gateway name provided")
	}
	_, vdc, err := cli.GetOrgAndVdc(orgName, vdcName)
	if err != nil {
		return nil, fmt.Errorf("error retrieving Org and VDC: %s", err)
	}
	eg, err = vdc.GetNsxtEdgeGatewayByName(edgeGwName)

	if err != nil {
		if os.Getenv("GOVCD_DEBUG") != "" {
			return nil, fmt.Errorf(fmt.Sprintf("(%s) [%s] ", edgeGwName, callFuncName())+errorUnableToFindEdgeGateway, err)
		}
		return nil, fmt.Errorf(errorUnableToFindEdgeGateway, err)
	}
	return eg, nil
}

// EOF ===================== To be removed before merge (will come from Security Group PR) ===========
