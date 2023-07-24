package vcd

//lint:file-ignore SA1019 ignore deprecated functions
import (
	"net/netip"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// suppressWordToEmptyString is a DiffSuppressFunc which ignore the change from word to empty string "".
// This is useful when API returns some default value but it is not set (and not sent via API) in config.
func suppressWordToEmptyString(word string) schema.SchemaDiffSuppressFunc {
	return func(k string, old string, new string, d *schema.ResourceData) bool {
		if old == word && new == "" {
			return true
		}
		return false
	}

}

// suppressNetworkUpgradedInterface is used to silence the changes in
// property "interface_type" in routed networks.
// In the old the version, the "internal" interface was implicit,
// while in the new one it is one of several.
// This function only considers the "internal" value, as the other interface types
// were not possible in the previous version
func suppressNetworkUpgradedInterface() schema.SchemaDiffSuppressFunc {
	return func(k string, old string, new string, d *schema.ResourceData) bool {
		if old == "" && new == "internal" {
			return true
		}
		return false
	}
}

// falseBoolSuppress suppresses change if value is set to false or is empty
func falseBoolSuppress() schema.SchemaDiffSuppressFunc {
	return func(k string, old string, new string, d *schema.ResourceData) bool {
		_, isTrue := d.GetOkExists(k)
		return !isTrue
	}
}

// suppressNewFalse always suppresses when new value is false
func suppressFalse() schema.SchemaDiffSuppressFunc {
	return func(k string, old string, new string, d *schema.ResourceData) bool {
		return new == "false"
	}
}

// suppressCase is a schema.SchemaDiffSuppressFunc which ignore case changes
func suppressCase(k, old, new string, d *schema.ResourceData) bool {
	return strings.EqualFold(old, new)
}

// suppressEqualIp is mainly useful for IPv6 where multiple formats might exist - long
// (2a00:a555:3000:1:0:0:0:32) and short (2a00:a555:3000:1::32), but also works with IPv4. VCD
// always returns long format IPv6 address, but a user might supply short format (as well as
// Terraform's native `cidrhost` function, therefore we want to avoid diff when a user uses his own
// format)
func suppressEqualIp(k, old, new string, d *schema.ResourceData) bool {
	oldIp, err := netip.ParseAddr(old)
	if err != nil {
		return false
	}

	newIp, err := netip.ParseAddr(new)
	if err != nil {
		return false
	}

	return oldIp.Compare(newIp) == 0
}
