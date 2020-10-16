package types

const (
	gatewayTypeNsxT = "NSXT_BACKED"
	gatewayTypeNsxV = "NSXV_BACKED"
)

type NsxtEdgeGateway2 struct {
	Status      string `json:"status,omitempty"`
	ID          string `json:"id,omitempty"`
	Name        string `json:"name"`
	Description string `json:"description"`
	OrgVdc      struct {
		ID string `json:"id"`
	} `json:"orgVdc"`
	EdgeGatewayUplinks []struct {
		UplinkID   string `json:"uplinkId"`
		UplinkName string `json:"uplinkName"`
		Subnets    struct {
			Values []struct {
				Gateway      string      `json:"gateway"`
				PrefixLength int         `json:"prefixLength"`
				DNSSuffix    interface{} `json:"dnsSuffix"`
				DNSServer1   string      `json:"dnsServer1"`
				DNSServer2   string      `json:"dnsServer2"`
				IPRanges     struct {
					Values []struct {
						StartAddress string `json:"startAddress"`
						EndAddress   string `json:"endAddress"`
					} `json:"values"`
				} `json:"ipRanges"`
				Enabled      bool   `json:"enabled"`
				TotalIPCount int    `json:"totalIpCount"`
				UsedIPCount  int    `json:"usedIpCount"`
				PrimaryIp    string `json:"primaryIp,omitempty"`
			} `json:"values"`
		} `json:"subnets"`
		Dedicated bool `json:"dedicated"`
	} `json:"edgeGatewayUplinks"`
}

type NsxtEdgeGateway struct {
	Status                    string               `json:"status,omitempty"`
	ID                        string               `json:"id,omitempty"`
	Name                      string               `json:"name"`
	Description               string               `json:"description"`
	EdgeGatewayUplinks        []EdgeGatewayUplinks `json:"edgeGatewayUplinks"`
	DistributedRoutingEnabled bool                 `json:"distributedRoutingEnabled"`
	OrgVdcNetworkCount        int                  `json:"orgVdcNetworkCount"`
	GatewayBacking            GatewayBacking       `json:"gatewayBacking"`
	OrgVdc                    OrgVdc               `json:"orgVdc"`
	OrgRef                    OrgRef               `json:"orgRef"`
	ServiceNetworkDefinition  string               `json:"serviceNetworkDefinition"`
	EdgeClusterConfig         EdgeClusterConfig    `json:"edgeClusterConfig"`
}

type NsxtRangeValues struct {
	StartAddress string `json:"startAddress"`
	EndAddress   string `json:"endAddress"`
}
type NsxtIPRanges struct {
	Values []NsxtRangeValues `json:"values"`
}
type NsxtSubnetValues struct {
	Gateway              string       `json:"gateway"`
	PrefixLength         int          `json:"prefixLength"`
	DNSSuffix            interface{}  `json:"dnsSuffix"`
	DNSServer1           string       `json:"dnsServer1"`
	DNSServer2           string       `json:"dnsServer2"`
	IPRanges             NsxtIPRanges `json:"ipRanges"`
	Enabled              bool         `json:"enabled"`
	TotalIPCount         int          `json:"totalIpCount"`
	UsedIPCount          interface{}  `json:"usedIpCount"`
	PrimaryIP            string       `json:"primaryIp"`
	AutoAllocateIPRanges bool         `json:"autoAllocateIpRanges"`
}
type NsxtSubnets struct {
	Values []NsxtSubnetValues `json:"values"`
}
type EdgeGatewayUplinks struct {
	UplinkID                 string      `json:"uplinkId"`
	UplinkName               string      `json:"uplinkName"`
	Subnets                  NsxtSubnets `json:"subnets"`
	Connected                bool        `json:"connected"`
	QuickAddAllocatedIPCount interface{} `json:"quickAddAllocatedIpCount"`
	Dedicated                bool        `json:"dedicated"`
}

// type NetworkProvider struct {
// 	Name string `json:"name"`
// 	ID   string `json:"id"`
// }
type GatewayBacking struct {
	BackingID       string          `json:"backingId"`
	GatewayType     string          `json:"gatewayType"`
	NetworkProvider NetworkProvider `json:"networkProvider"`
}
type OrgVdc struct {
	Name string `json:"name"`
	ID   string `json:"id"`
}
type OrgRef struct {
	Name string `json:"name"`
	ID   string `json:"id"`
}
type EdgeClusterRef struct {
	Name string `json:"name"`
	ID   string `json:"id"`
}
type PrimaryEdgeCluster struct {
	EdgeClusterRef EdgeClusterRef `json:"edgeClusterRef"`
	BackingID      string         `json:"backingId"`
}
type EdgeClusterConfig struct {
	PrimaryEdgeCluster   PrimaryEdgeCluster `json:"primaryEdgeCluster"`
	SecondaryEdgeCluster interface{}        `json:"secondaryEdgeCluster"`
}
