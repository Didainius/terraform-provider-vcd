---
layout: "vcd"
page_title: "VMware Cloud Director: vcd_ip_space"
sidebar_current: "docs-vcd-resource-ip-space"
description: |-
  Provides a resource to manage IP Spaces for IP address management needs. IP Spaces provide structured approach to allocating public and private IP addresses by preventing the use of overlapping IP addresses across organizations and organization VDCs.
---

# vcd\_ip\_space

IP Spaces require VCD 10.4.1+ with NSX-T.

Provides a resource to manage IP Spaces for IP address management needs. IP Spaces provide
structured approach to allocating public and private IP addresses by preventing the use of
overlapping IP addresses across organizations and organization VDCs.


## Example Usage (Private)

```hcl
resource "vcd_ip_space" "space1" {
  name        = "TestAccVcdIpSpacePrivate"
  description = "added description"
  type        = "PRIVATE"
  org_id      = data.vcd_org.org1.id

  internal_scope = ["192.168.1.0/24","10.10.10.0/24", "11.11.11.0/24"]

  route_advertisement_enabled = false

  ip_prefixes {
	default_quota = -1 # no quota

	prefix {
		first_ip = "192.168.1.100"
		prefix_length = 30
		prefix_count = 4
	}
  }

  ip_prefixes {
	default_quota = -1 # no quota

	prefix {
		first_ip = "10.10.10.96"
		prefix_length = 29
		prefix_count = 4
	}
  }
}
```

## Example Usage (Public)

```hcl
resource "vcd_ip_space" "space1" {
  name        = "TestAccVcdIpSpacePublic"
  description = "added description"
  type        = "PUBLIC"

  internal_scope = ["192.168.1.0/24","10.10.10.0/24", "11.11.11.0/24"]

  route_advertisement_enabled = false

  ip_prefixes {
	default_quota = 2

	prefix {
		first_ip = "192.168.1.100"
		prefix_length = 30
		prefix_count = 4
	}
  }

  ip_prefixes {
	default_quota = -1

	prefix {
		first_ip = "10.10.10.96"
		prefix_length = 29
		prefix_count = 4
	}
  }
}
```

## Example Usage (Shared)

```hcl
resource "vcd_ip_space" "space1" {
  name        = "TestAccVcdIpSpaceShared"
  description = "added description"
  type        = "SHARED_SERVICES"

  internal_scope = ["192.168.1.0/24","10.10.10.0/24", "11.11.11.0/24"]

  route_advertisement_enabled = false

  ip_prefixes {
	 default_quota = 0 # no quota

	prefix {
		first_ip = "192.168.1.100"
		prefix_length = 30
		prefix_count = 4
	}

	prefix {
		first_ip = "192.168.1.200"
		prefix_length = 30
		prefix_count = 4
	}
  }

  ip_prefixes {
	default_quota = 0 # no quota

	prefix {
		first_ip = "10.10.10.96"
		prefix_length = 29
		prefix_count = 4
	}
  }

  ip_ranges {
	start_address = "11.11.11.100"
	end_address   = "11.11.11.110"
  }

  ip_ranges {
	start_address = "11.11.11.120"
	end_address   = "11.11.11.123"
  }
}
```

## Argument Reference

The following arguments are supported:

* `org_id` - (Optional) Required for `PRIVATE` type
* `name` - (Required) A name for IP Space
* `description` - (Optional) - Description of IP Space
* `type` - (Required) One of `PUBLIC`, `SHARED_SERVICES`, `PRIVATE`
  * `PUBLIC` - A public IP space is used by multiple organizations and is controlled by the service
    provider through a quota-based system. 
  * `SHARED_SERVICES` - An IP space for services and management networks that are required in the
    tenant space, but as a service provider, you don't want to expose it to organizations in your
    environment. 
  * `PRIVATE` - Private IP spaces are dedicated to a single tenant - a private IP space is used by
    only one organization that is specified during the space creation. For this organization, IP
    consumption is unlimited.

* `internal_scope` - (Required) The internal scope of an IP space is a list of CIDR notations that
  defines the exact span of IP addresses in which all ranges and blocks must be contained in. The
  internal scope of the IP space is used to define default NAT rules and BGP prefixes. 
* `ip_range` - (Optional) One or more ranges for floating IP address allocation. (Floating IP
  addresses are just IP addresses taken from the defined range) [ip_range](#ipspace-ip-range)
* `ip_range_quota` - (Optional) If you entered at least one IP Range (`ip_range`) page, enter a
  number of floating IP addresses to allocate individually. `-1` is unlimited, while `0` means that
  no IPs can be allocated.
* `ip_prefix` - (Optional) One or more IP prefixes (blocks) [ip_prefix](#ipspace-ip-prefix)
* `external_scope` - (Optional) The external scope defines the total span of IP addresses to which the IP
  space has access, for example the internet or a WAN. The external scope of the IP space is used to
  define default NAT rules and BGP prefixes (e.g. 0.0.0.0/24)
* `route_advertisement_enabled` - (default `false`) Toggle on the route advertisement option to
  enable advertising networks with IP prefixes from this IP space
* ``

<a id="ipspace-ip-range"></a>

* `start_address` - (Required) - Start IP address of a range
* `end_address` - (Required) - End IP address of a range

<a id="ipspace-ip-prefix"></a>

* `default_quota` 
* `prefix` 

<a id="ipspace-ip-prefix-prefix"></a>

Defines blocks of IPs. Blocks must fall into subnets defined in `internal_scope` and not clash with
IP ranges defined in `ip_range` 

* `first_ip` - (Required) - First IP of the prefix
* `prefix_length` - (Required) - Prefix length
* `prefix_count` - (Required) - Number of prefixes 

```hcl
 ip_prefix {
  default_quota = 2

  prefix {
    first_ip      = "192.168.1.100"
    prefix_length = 30
    prefix_count  = 4
  }

  prefix {
    first_ip      = "192.168.1.200"
    prefix_length = 30
    prefix_count  = 4
  }
}
```

## Attribute Reference

The following attributes are exported on this resource:

* `max_virtual_services` - Maximum number of virtual services this NSX-T ALB Service Engine Group can run


## Importing

~> The current implementation of Terraform import can only import resources into the state.
It does not generate configuration. [More information.](https://www.terraform.io/docs/import/)

An existing NSX-T ALB Service Engine Group configuration can be [imported][docs-import] into this resource
via supplying path for it. An example is
below:

[docs-import]: https://www.terraform.io/docs/import/

```
terraform import vcd_ip_space.imported ip-space-name
```

The above would import the `ip-space-name` IP Space is defined at provider
level.
