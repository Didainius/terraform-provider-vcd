---
layout: "vcd"
page_title: "VMware Cloud Director: vcd_nsxt_nat_rule"
sidebar_current: "docs-vcd-resource-nsxt-nat-rule"
description: |-
  Provides a resource to manage NSX-T NAT rules. To change the source IP address from a private to a
  public IP address, you create a source NAT (SNAT) rule. To change the destination IP address from 
  a public to a private IP address, you create a destination NAT (DNAT) rule.
---

# vcd\_nsxt\_nat\_rule

Supported in provider *v3.3+* and VCD 10.1+ with NSX-T backed VDCs.

Provides a resource to manage NSX-T NAT rules. To change the source IP address from a private to a
public IP address, you create a source NAT (SNAT) rule. To change the destination IP address from 
a public to a private IP address, you create a destination NAT (DNAT) rule.

-> When you configure a SNAT or a DNAT rule on an edge gateway in the VMware Cloud Director
environment, you always configure the rule from the perspective of your organization VDC.

## Example Usage 1 (SNAT rule translates to primary Edge Gateway IP for traffic going from 11.11.11.2 to 8.8.8.8)

```hcl
resource "vcd_nsxt_nat_rule" "snat" {
  org  = "dainius"
  vdc  = "nsxt-vdc-dainius"

  edge_gateway_id = data.vcd_nsxt_edgegateway.existing.id

  name        = "SNAT rule"
  rule_type   = "SNAT"
  description = "description"

  # Using primary_ip from edge gateway
  external_addresses         = tolist(data.vcd_nsxt_edgegateway.existing.subnet)[0].primary_ip
  internal_addresses         = "11.11.11.2"
  snat_destination_addresses = "8.8.8.8"
  logging = true
}
```

## Example Usage 2 (No SNAT rule)
```hcl
resource "vcd_nsxt_nat_rule" "no-snat" {
  org  = "dainius"
  vdc  = "nsxt-vdc-dainius"

  edge_gateway_id = data.vcd_nsxt_edgegateway.existing.id

  name        = "test-no-snat-rule"
  rule_type   = "NO_SNAT"
  description = "description"

  # Using primary_ip from edge gateway
  internal_addresses         = "11.11.11.2"
}
```

## Example Usage 3 (DNAT rule)
```hcl
resource "vcd_nsxt_nat_rule" "dnat" {
  org = "my-org"
  vdc = "nsxt-vdc"

  edge_gateway_id = data.vcd_nsxt_edgegateway.existing.id

  name        = "test-dnat-rule"
  rule_type   = "DNAT"
  description = "description"

  # Using primary_ip from edge gateway
  external_addresses = tolist(data.vcd_nsxt_edgegateway.existing.subnet)[0].primary_ip
  internal_addresses = "11.11.11.2"
  logging            = true
}
```

## Example Usage 4 (No DNAT rule)
```hcl
resource "vcd_nsxt_nat_rule" "no-dnat" {
  org = "my-org"
  vdc = "nsxt-vdc"

  edge_gateway_id = data.vcd_nsxt_edgegateway.existing.id

  name      = "test-no-dnat-rule"
  rule_type = "NO_DNAT"


  # Using primary_ip from edge gateway
  external_addresses = tolist(data.vcd_nsxt_edgegateway.existing.subnet)[0].primary_ip
  dnat_external_port = 7777

  # app_port_profile_id =
}
```

## Argument Reference

The following arguments are supported:

* `org` - (Optional) The name of organization to use, optional if defined at provider level. Useful
  when connected as sysadmin working across different organisations.
* `vdc` - (Optional) The name of VDC to use, optional if defined at provider level.
* `edge_gateway_id` - (Required) The ID of the edge gateway (NSX-T only). Can be looked up using
  `vcd_nsxt_edgegateway` datasource
* `name` - (Required) A name for NAT rule
* `description` - (Optional) An optional description of the NAT rule
* `rule_type` - (Required) One of `DNAT`, `NO_DNAT`, `SNAT`, `NO_SNAT`

  * `DNAT` rule translates the IP address and, optionally, the port of packets received by an
    organization VDC network that are coming from an external network or from another organization
    VDC network.
  * `NO_DNAT` rule prevents the translation of the external IP address of packets received by an
    organization VDC from an external network or from another organization VDC network.
  * `SNAT` rule translates the source IP address of packets sent from an organization VDC network
    out to an external network or to another organization VDC network.
  * `NO_SNAT` rule prevents the translation of the internal IP address of packets sent from an
    organization VDC out to an external network or to another organization VDC network.

* `external_addresses` (Optional) value depends on `rule_type`
  * `SNAT` - the public IP address of the edge gateway for which you are configuring the SNAT rule
  * `NO_SNAT` - leave empty
  * `DNAT` - the public IP address of the edge gateway for which you are configuring the DNAT rule.
    The IP addresses that you enter must belong to the suballocated IP range of the edge gateway.

* `internal_addresses` (Optional) Enter the IP address or a range of IP addresses of the virtual
  machines for which you are configuring SNAT, so that they can send traffic to the external
  network.
  
* `app_port_profile_id` (Optional) - Select a specific application port profile to which to apply
  the rule. The application port profile includes a port and a protocol that the incoming traffic
  uses on the edge gateway to connect to the internal network.
* `dnat_external_port` (Optional) - Enter a port into which the DNAT rule is translating for the
  packets inbound to the virtual machines.
* `snat_destination_addresses` (Optional) For `SNAT` only. If you want the rule to apply only for
  traffic to a specific domain, enter an IP address for this domain or an IP address range in CIDR
  format. If you leave this text box blank, the SNAT rule applies to all destinations outside of the
  local subnet.
* `logging` (Optional) - to have the address translation performed by this rule logged, toggle on
  the Logging option
* `enabled` (Optional) - allows to enable or disable NAT rule (default `true`)
* `app_port_profile_id` (Optional) Reference to Application Port Profile. `vcd_app

* `firewall_match` (Optional, VCD 10.2.2+) - You can set a firewall match rule to determine how
  firewall is applied during NAT. One of `MATCH_INTERNAL_ADDRESS`, `MATCH_EXTERNAL_ADDRESS`,
  `BYPASS`

  * `MATCH_INTERNAL_ADDRESS` - applies firewall rules to the internal address of a NAT rule
  * `MATCH_EXTERNAL_ADDRESS` - applies firewall rules to the external address of a NAT rule
  * `BYPASS` - skip applying firewall rules to NAT rule


* `priority` (Optional, VCD 10.2.2+) - if an address has multiple NAT rules, you can assign these
  rules different priorities to determine the order in which they are applied. A lower value means a
  higher priority for this rule. 

## Importing

~> The current implementation of Terraform import can only import resources into the state.
It does not generate configuration. [More information.](https://www.terraform.io/docs/import/)

An existing NAT Rule configuration can be [imported][docs-import] into this resource
via supplying the full dot separated path for your NAT Rule name or ID. An example is
below:

[docs-import]: https://www.terraform.io/docs/import/

Supplying Name
```
terraform import vcd_nsxt_nat_rule.imported my-org.my-org-vdc.my-nsxt-edge-gateway.my-nat-rule-name
```



-> When there are multiple NAT rules with the same name they will all be listed so that one can pick
it by ID

```
$ terraform import vcd_nsxt_nat_rule.dnat my-org.nsxt-vdc.nsxt-gw.dnat1

vcd_nsxt_nat_rule.dnat: Importing from ID "my-org.nsxt-vdc.nsxt-gw.dnat1"...
# The following NAT rules with Name 'dnat1' are available
# Please use ID instead of Name in import path to pick exact rule
ID                                   Name  Rule Type Internal Addresses External Addresses
04fde766-2cbd-4986-93bb-7f57e59c6b19 dnat1 DNAT      1.1.1.1            10.1.2.139
f40e3d68-cfa6-42ea-83ed-5571659b3e7b dnat1 DNAT      2.2.2.2            10.1.2.139

$ terraform import vcd_nsxt_nat_rule.imported my-org.my-org-vdc.my-nsxt-edge-gateway.0214a26b-fc30-4202-88e5-7ed551aa6c19
```

The above would import the `my-nat-rule-name` NAT Rule config settings that are defined
on NSX-T Edge Gateway `my-nsxt-edge-gateway` which is configured in organization named `my-org` and
VDC named `my-org-vdc`.
