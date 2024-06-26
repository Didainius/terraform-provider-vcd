---
layout: "vcd"
page_title: "VMware Cloud Director: vcd_nsxt_alb_controller"
sidebar_current: "docs-vcd-datasource-nsxt-alb-controller"
description: |-
  Provides a data source to read ALB Controller for Providers. It helps to integrate VMware Cloud Director with
  Avi Load Balancer deployment. Controller instances are registered with VMware Cloud Director instance.
  Controller instances serve as a central control plane for the load-balancing services provided by Avi Load
  Balancer.
---

# vcd\_nsxt\_alb\_controller

Supported in provider *v3.4+* and VCD 10.2+ with NSX-T and ALB.

Provides a data source to read ALB Controller for Providers. It helps to integrate VMware Cloud Director with
Avi Load Balancer deployment. Controller instances are registered with VMware Cloud Director instance.
Controller instances serve as a central control plane for the load-balancing services provided by Avi Load
Balancer.

~> Only `System Administrator` can use this data source.

~> VCD 10.3.0 has a caching bug which prevents listing importable clouds immediately (retrieved using
[`vcd_nsxt_alb_importable_cloud`](/providers/vmware/vcd/latest/docs/data-sources/nsxt_alb_importable_cloud)) after ALB
Controller is created. This data should be available 15 minutes after the Controller is created.

## Example Usage

```hcl
data "vcd_nsxt_alb_controller" "first" {
  name = "avi controller"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required)  - Unique name of existing ALB Controller.

## Attribute Reference

All the arguments and attributes defined in
[`vcd_nsxt_alb_controller`](/providers/vmware/vcd/latest/docs/resources/nsxt_alb_controller) resource are available.
