package vcd

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func datasourceVcdNsxtTier0Router() *schema.Resource {
	return &schema.Resource{
		Read: datasourceNsxtTier0RouterRead,
		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of NSX-T Tier-0 router.",
			},
			"nsxt_manager_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of NSX-T manager.",
			},
		},
	}
}

func datasourceNsxtTier0RouterRead(d *schema.ResourceData, meta interface{}) error {
	vcdClient := meta.(*VCDClient)
	nsxtManagerId := d.Get("nsxt_manager_id").(string)
	tier0RouterName := d.Get("name").(string)

	tier0Router, err := vcdClient.GetNsxtTier0RouterByName(tier0RouterName, nsxtManagerId)
	if err != nil {
		return fmt.Errorf("could not find NSX-T Tier-0 router by name '%s' in NSX-T manager %s: %s", tier0RouterName, nsxtManagerId, err)
	}

	d.SetId(tier0Router.NsxtTier0Router.ID)
	return nil
}
