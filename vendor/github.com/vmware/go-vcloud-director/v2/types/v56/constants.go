/*
 * Copyright 2019 VMware, Inc.  All rights reserved.  Licensed under the Apache v2 License.
 */

package types

const (
	// PublicCatalog Name
	PublicCatalog = "Public Catalog"

	// DefaultCatalog Name
	DefaultCatalog = "Default Catalog"

	// JSONMimeV57 the json mime for version 5.7 of the API
	JSONMimeV57 = "application/json;version=5.7"
	// AnyXMLMime511 the wildcard xml mime for version 5.11 of the API
	AnyXMLMime511 = "application/*+xml;version=5.11"
	AnyXMLMime    = "application/xml"
	// Version511 the 5.11 version
	Version511 = "5.11"
	// Version is the default version number
	Version = Version511
	// SoapXML mime type
	SoapXML = "application/soap+xml"
	// JSONMime
	JSONMime = "application/json"
)

const (
	// MimeOrgList mime for org list
	MimeOrgList = "application/vnd.vmware.vcloud.orgList+xml"
	// MimeOrg mime for org
	MimeOrg = "application/vnd.vmware.vcloud.org+xml"
	// MimeAdminOrg mime for admin org
	MimeAdminOrg = "application/vnd.vmware.admin.organization+xml"
	// MimeCatalog mime for catalog
	MimeCatalog = "application/vnd.vmware.vcloud.catalog+xml"
	// MimeCatalogItem mime for catalog item
	MimeCatalogItem = "application/vnd.vmware.vcloud.catalogItem+xml"
	// MimeVDC mime for a VDC
	MimeVDC = "application/vnd.vmware.vcloud.vdc+xml"
	// MimeVDC mime for a admin VDC
	MimeAdminVDC = "application/vnd.vmware.admin.vdc+xml"
	// MimeEdgeGateway mime for an Edge Gateway
	MimeEdgeGateway = "application/vnd.vmware.admin.edgeGateway+xml"
	// MimeVAppTemplate mime for a vapp template
	MimeVAppTemplate = "application/vnd.vmware.vcloud.vAppTemplate+xml"
	// MimeVApp mime for a vApp
	MimeVApp = "application/vnd.vmware.vcloud.vApp+xml"
	// MimeQueryRecords mime for the query records
	MimeQueryRecords = "application/vnd.vmware.vcloud.query.records+xml"
	// MimeAPIExtensibility mime for api extensibility
	MimeAPIExtensibility = "application/vnd.vmware.vcloud.apiextensibility+xml"
	// MimeEntity mime for vcloud entity
	MimeEntity = "application/vnd.vmware.vcloud.entity+xml"
	// MimeQueryList mime for query list
	MimeQueryList = "application/vnd.vmware.vcloud.query.queryList+xml"
	// MimeSession mime for a session
	MimeSession = "application/vnd.vmware.vcloud.session+xml"
	// MimeTask mime for task
	MimeTask = "application/vnd.vmware.vcloud.task+xml"
	// MimeError mime for error
	MimeError = "application/vnd.vmware.vcloud.error+xml"
	// MimeNetwork mime for a network
	MimeNetwork = "application/vnd.vmware.vcloud.network+xml"
	// MimeOrgVdcNetwork mime for an Org VDC network
	MimeOrgVdcNetwork = "application/vnd.vmware.vcloud.orgVdcNetwork+xml"
	//MimeDiskCreateParams mime for create independent disk
	MimeDiskCreateParams = "application/vnd.vmware.vcloud.diskCreateParams+xml"
	// Mime for VMs
	MimeVMs = "application/vnd.vmware.vcloud.vms+xml"
	// Mime for attach or detach independent disk
	MimeDiskAttachOrDetachParams = "application/vnd.vmware.vcloud.diskAttachOrDetachParams+xml"
	// Mime for Disk
	MimeDisk = "application/vnd.vmware.vcloud.disk+xml"
	// Mime for insert or eject media
	MimeMediaInsertOrEjectParams = "application/vnd.vmware.vcloud.mediaInsertOrEjectParams+xml"
	// Mime for catalog
	MimeAdminCatalog = "application/vnd.vmware.admin.catalog+xml"
	// Mime for virtual hardware section
	MimeVirtualHardwareSection = "application/vnd.vmware.vcloud.virtualHardwareSection+xml"
	// Mime for networkConnectionSection
	MimeNetworkConnectionSection = "application/vnd.vmware.vcloud.networkConnectionSection+xml"
	// Mime for Item
	MimeRasdItem = "application/vnd.vmware.vcloud.rasdItem+xml"
	// Mime for guest customization section
	MimeGuestCustomizationSection = "application/vnd.vmware.vcloud.guestCustomizationSection+xml"
	// Mime for guest customization status
	MimeGuestCustomizationStatus = "application/vnd.vmware.vcloud.guestcustomizationstatussection"
	// Mime for network config section
	MimeNetworkConfigSection = "application/vnd.vmware.vcloud.networkconfigsection+xml"
	// Mime for recompose vApp params
	MimeRecomposeVappParams = "application/vnd.vmware.vcloud.recomposeVAppParams+xml"
	// Mime for compose vApp params
	MimeComposeVappParams = "application/vnd.vmware.vcloud.composeVAppParams+xml"
	// Mime for undeploy vApp params
	MimeUndeployVappParams = "application/vnd.vmware.vcloud.undeployVAppParams+xml"
	// Mime for deploy vApp params
	MimeDeployVappParams = "application/vnd.vmware.vcloud.deployVAppParams+xml"
	// Mime for VM
	MimeVM = "application/vnd.vmware.vcloud.vm+xml"
	// Mime for instantiate vApp template params
	MimeInstantiateVappTemplateParams = "application/vnd.vmware.vcloud.instantiateVAppTemplateParams+xml"
	// Mime for product section
	MimeProductSection = "application/vnd.vmware.vcloud.productSections+xml"
	// Mime for metadata
	MimeMetaData = "application/vnd.vmware.vcloud.metadata+xml"
	// Mime for metadata value
	MimeMetaDataValue = "application/vnd.vmware.vcloud.metadata.value+xml"
	// Mime for a admin network
	MimeExtensionNetwork = "application/vnd.vmware.admin.extension.network+xml"
	// Mime for an external network
	MimeExternalNetwork = "application/vnd.vmware.admin.vmwexternalnet+xml"
	// Mime of an Org User
	MimeAdminUser = "application/vnd.vmware.admin.user+xml"
	// MimeAdminGroup specifies groups
	MimeAdminGroup = "application/vnd.vmware.admin.group+xml"
	// MimeOrgLdapSettings
	MimeOrgLdapSettings = "application/vnd.vmware.admin.organizationldapsettings+xml"
	// Mime of vApp network
	MimeVappNetwork = "application/vnd.vmware.vcloud.vAppNetwork+xml"
	// Mime of access control
	MimeControlAccess = "application/vnd.vmware.vcloud.controlAccess+xml"
	// Mime of VM capabilities
	MimeVmCapabilities = "application/vnd.vmware.vcloud.vmCapabilitiesSection+xml"
	// Mime of Vdc Compute Policy References
	MimeVdcComputePolicyReferences = "application/vnd.vmware.vcloud.vdcComputePolicyReferences+xml"
)

const (
	VMsCDResourceSubType = "vmware.cdrom.iso"
)

// https://blogs.vmware.com/vapp/2009/11/virtual-hardware-in-ovf-part-1.html

const (
	ResourceTypeOther     int = 0
	ResourceTypeProcessor int = 3
	ResourceTypeMemory    int = 4
	ResourceTypeIDE       int = 5
	ResourceTypeSCSI      int = 6
	ResourceTypeEthernet  int = 10
	ResourceTypeFloppy    int = 14
	ResourceTypeCD        int = 15
	ResourceTypeDVD       int = 16
	ResourceTypeDisk      int = 17
	ResourceTypeUSB       int = 23
)

const (
	FenceModeIsolated = "isolated"
	FenceModeBridged  = "bridged"
	FenceModeNAT      = "natRouted"
)

const (
	IPAllocationModeDHCP   = "DHCP"
	IPAllocationModeManual = "MANUAL"
	IPAllocationModeNone   = "NONE"
	IPAllocationModePool   = "POOL"
)

// NoneNetwork is a special type of network in vCD which represents a network card which is not
// attached to any network.
const (
	NoneNetwork = "none"
)

const (
	XMLNamespaceVCloud    = "http://www.vmware.com/vcloud/v1.5"
	XMLNamespaceOVF       = "http://schemas.dmtf.org/ovf/envelope/1"
	XMLNamespaceVMW       = "http://www.vmware.com/schema/ovf"
	XMLNamespaceXSI       = "http://www.w3.org/2001/XMLSchema-instance"
	XMLNamespaceRASD      = "http://schemas.dmtf.org/wbem/wscim/1/cim-schema/2/CIM_ResourceAllocationSettingData"
	XMLNamespaceVSSD      = "http://schemas.dmtf.org/wbem/wscim/1/cim-schema/2/CIM_VirtualSystemSettingData"
	XMLNamespaceExtension = "http://www.vmware.com/vcloud/extension/v1.5"
)

// NSX-V Edge gateway API endpoints
const (
	EdgeNatPath            = "/nat/config"
	EdgeCreateNatPath      = "/nat/config/rules"
	EdgeFirewallPath       = "/firewall/config"
	EdgeCreateFirewallPath = "/firewall/config/rules"
	EdgeVnicConfig         = "/vnics"
	EdgeVdcVnicConfig      = "/vdcNetworks"
	EdgeDhcpRelayPath      = "/dhcp/config/relay"
	EdgeDhcpLeasePath      = "/dhcp/leaseInfo"
	LbConfigPath           = "/loadbalancer/config/"
	LbMonitorPath          = "/loadbalancer/config/monitors/"
	LbServerPoolPath       = "/loadbalancer/config/pools/"
	LbAppProfilePath       = "/loadbalancer/config/applicationprofiles/"
	LbAppRulePath          = "/loadbalancer/config/applicationrules/"
	LbVirtualServerPath    = "/loadbalancer/config/virtualservers/"
)

// NSX-V proxied services API endpoints
const (
	NsxvIpSetServicePath = "/ipset"
)

// Guest customization statuses. These are all known possible statuses
const (
	GuestCustStatusPending       = "GC_PENDING"
	GuestCustStatusPostPending   = "POST_GC_PENDING"
	GuestCustStatusComplete      = "GC_COMPLETE"
	GuestCustStatusFailed        = "GC_FAILED"
	GuestCustStatusRebootPending = "REBOOT_PENDING"
)

// Edge gateway vNic types
const (
	EdgeGatewayVnicTypeUplink       = "uplink"
	EdgeGatewayVnicTypeInternal     = "internal"
	EdgeGatewayVnicTypeTrunk        = "trunk"
	EdgeGatewayVnicTypeSubinterface = "subinterface"
	EdgeGatewayVnicTypeAny          = "any"
)

// Names of the filters allowed in the search engine
const (
	FilterNameRegex = "name_regex" // a name, searched by regular expression
	FilterDate      = "date"       // a date expression (>|<|==|>=|<= date)
	FilterIp        = "ip"         // An IP, searched by regular expression
	FilterLatest    = "latest"     // gets the newest element
	FilterEarliest  = "earliest"   // gets the oldest element
	FilterParent    = "parent"     // matches the entity parent
	FilterParentId  = "parent_id"  // matches the entity parent ID
)

const (
	// The Qt* (Query Type) constants are the names used with Query requests to retrieve the corresponding entities
	QtVappTemplate      = "vAppTemplate"      // vApp template
	QtAdminVappTemplate = "adminVAppTemplate" // vApp template as admin
	QtEdgeGateway       = "edgeGateway"       // edge gateway
	QtOrgVdcNetwork     = "orgVdcNetwork"     // Org VDC network
	QtCatalog           = "catalog"           // catalog
	QtAdminCatalog      = "adminCatalog"      // catalog as admin
	QtCatalogItem       = "catalogItem"       // catalog item
	QtAdminCatalogItem  = "adminCatalogItem"  // catalog item as admin
	QtAdminMedia        = "adminMedia"        // media item as admin
	QtMedia             = "media"             // media item
	QtVm                = "vm"                // Virtual machine
	QtAdminVm           = "adminVM"           // Virtual machine as admin
	QtVapp              = "vApp"              // vApp
	QtAdminVapp         = "adminVApp"         // vApp as admin
)

// AdminQueryTypes returns the corresponding "admin" query type for each regular type
var AdminQueryTypes = map[string]string{
	QtEdgeGateway:   QtEdgeGateway,   // EdgeGateway query type is the same for admin and regular users
	QtOrgVdcNetwork: QtOrgVdcNetwork, // Org VDC Network query type is the same for admin and regular users
	QtVappTemplate:  QtAdminVappTemplate,
	QtCatalog:       QtAdminCatalog,
	QtCatalogItem:   QtAdminCatalogItem,
	QtMedia:         QtAdminMedia,
	QtVm:            QtAdminVm,
	QtVapp:          QtAdminVapp,
}

const (
	// Affinity and anti affinity definitions
	PolarityAffinity     = "Affinity"
	PolarityAntiAffinity = "Anti-Affinity"
)

// VmQueryFilter defines how we search VMs
type VmQueryFilter int

const (
	// VmQueryFilterAll defines a no-filter search, i.e. will return all elements
	VmQueryFilterAll VmQueryFilter = iota

	// VmQueryFilterOnlyDeployed defines a filter for deployed VMs
	VmQueryFilterOnlyDeployed

	// VmQueryFilterOnlyTemplates defines a filter for VMs inside a template
	VmQueryFilterOnlyTemplates
)

// String converts a VmQueryFilter into the corresponding filter needed by the query to get the wanted result
func (qf VmQueryFilter) String() string {
	// Makes sure that we handle out-of-range values
	if qf < VmQueryFilterAll || qf > VmQueryFilterOnlyTemplates {
		return ""
	}
	return [...]string{
		"",                      // No filter: will not remove any items
		"isVAppTemplate==false", // Will find only the deployed VMs
		"isVAppTemplate==true",  // Will find only those VM that are inside a template
	}[qf]
}

// LDAP modes for Organization
const (
	LdapModeNone   = "NONE"
	LdapModeSystem = "SYSTEM"
	LdapModeCustom = "CUSTOM"
)

// Access control modes
const (
	ControlAccessReadOnly    = "ReadOnly"
	ControlAccessReadWrite   = "Change"
	ControlAccessFullControl = "FullControl"
)

// BodyType allows to define API body types where applicable
type BodyType int

const (
	// BodyTypeXML
	BodyTypeXML BodyType = iota

	// BodyTypeJSON
	BodyTypeJSON
)

const (
	// FiqlQueryTimestampFormat is the format accepted by Cloud API time comparison operator in FIQL query filters
	FiqlQueryTimestampFormat = "2006-01-02T15:04:05.000Z"
)

// These constants allow to construct OpenAPI endpoint paths and avoid strings in code for easy replacement in future.
const (
	OpenApiPathVersion1_0_0                   = "1.0.0/"
	OpenApiEndpointRoles                      = "roles/"
	OpenApiEndpointAuditTrail                 = "auditTrail/"
	OpenApiEndpointImportableTier0Routers     = "nsxTResources/importableTier0Routers"
	OpenApiEndpointExternalNetworks           = "externalNetworks/"
	OpenApiEndpointVdcComputePolicies         = "vdcComputePolicies/"
	OpenApiEndpointVdcAssignedComputePolicies = "vdcs/%s/computePolicies"
	OpenApiEndpointEdgeGateways               = "edgeGateways/"
)

// Header keys to run operations in tenant context
const (
	// HeaderTenantContext requires the Org ID of the tenant
	HeaderTenantContext = "X-VMWARE-VCLOUD-TENANT-CONTEXT"
	// HeaderAuthContext requires the Org name of the tenant
	HeaderAuthContext = "X-VMWARE-VCLOUD-AUTH-CONTEXT"
)

const (
	// ExternalNetworkBackingTypeNsxtTier0Router defines backing type of NSX-T Tier-0 router
	ExternalNetworkBackingTypeNsxtTier0Router = "NSXT_TIER0"
	// ExternalNetworkBackingTypeNsxtVrfTier0Router defines backing type of NSX-T Tier-0 VRF router
	ExternalNetworkBackingTypeNsxtVrfTier0Router = "NSXT_VRF_TIER0"
	// ExternalNetworkBackingTypeNetwork defines vSwitch portgroup
	ExternalNetworkBackingTypeNetwork = "NETWORK"
	// ExternalNetworkBackingDvPortgroup refers distributed switch portgroup
	ExternalNetworkBackingDvPortgroup = "DV_PORTGROUP"
)
