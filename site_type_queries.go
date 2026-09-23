package cato_go_sdk

type SiteBgpStatusResult struct {
	BGPSession                string             `json:"BGP_Session,omitempty"`
	IncomingConnection        IncomingConnection `json:"Incoming_Connection,omitempty"`
	LocalASN                  string             `json:"Local_ASN,omitempty"`
	LocalIP                   string             `json:"Local_IP,omitempty"`
	Negotiated                Negotiated         `json:"Negotiated,omitempty"`
	OutgoingConnection        OutgoingConnection `json:"Outgoing_Connection,omitempty"`
	RIBOut                    []RIBOut           `json:"RIB_out,omitempty"`
	RemoteASN                 string             `json:"Remote_ASN,omitempty"`
	RemoteIP                  string             `json:"Remote_IP,omitempty"`
	RouterID                  string             `json:"Router_ID,omitempty"`
	RouterWeight              string             `json:"Router_Weight,omitempty"`
	AcceptDefaultRoute        bool               `json:"accept_default_route,omitempty"`
	AcceptFromNeighbor        bool               `json:"accept_from_neighbor,omitempty"`
	AdvertiseAccountRoutes    bool               `json:"advertise_account_routes,omitempty"`
	AdvertiseDefaultRoute     bool               `json:"advertise_default_route,omitempty"`
	AdvertiseSummaryRoutes    bool               `json:"advertise_summary_routes,omitempty"`
	IsCatod                   bool               `json:"is_catod,omitempty"`
	RoutesCount               string             `json:"routes_count,omitempty"`
	RoutesCountLimit          string             `json:"routes_count_limit,omitempty"`
	RoutesCountLimitExceeded  bool               `json:"routes_count_limit_exceeded,omitempty"`
	SummarizeVpnUsers         bool               `json:"summarize_vpn_users,omitempty"`
	TunnelPcapCapturedPackets string             `json:"tunnel_pcap_captured_packets,omitempty"`
	TunnelPcapEnabled         bool               `json:"tunnel_pcap_enabled,omitempty"`
}
type IncomingConnection struct {
	State     string `json:"State,omitempty"`
	Transport string `json:"Transport,omitempty"`
}
type Capabilities struct {
	As4                          string `json:"as4,omitempty"`
	EnhancedRouteRefresh         string `json:"enhanced_route_refresh,omitempty"`
	GracefulRestartAlwaysPublish bool   `json:"graceful_restart_always_publish,omitempty"`
	GracefulRestartEnabled       bool   `json:"graceful_restart_enabled,omitempty"`
	GracefulRestartPresent       bool   `json:"graceful_restart_present,omitempty"`
	GracefulRestartTimeout       string `json:"graceful_restart_timeout,omitempty"`
	MultiprotocolExt             string `json:"multiprotocol_ext,omitempty"`
	RouteRefresh                 string `json:"route_refresh,omitempty"`
}
type Negotiated struct {
	Capabilities    Capabilities `json:"Capabilities,omitempty"`
	HoldTime        string       `json:"Hold_Time,omitempty"`
	KeepalivePeriod string       `json:"Keepalive_Period,omitempty"`
}
type OutgoingConnection struct {
	State     string `json:"State,omitempty"`
	Transport string `json:"Transport,omitempty"`
}
type RIBOut struct {
	AddTime        string `json:"add_time,omitempty"`
	LastUpdateTime string `json:"last_update_time,omitempty"`
	Med            int    `json:"med,omitempty"`
	NextHop        string `json:"next_hop,omitempty"`
	PrependCount   int    `json:"prepend_count,omitempty"`
	Subnet         string `json:"subnet,omitempty"`
}
