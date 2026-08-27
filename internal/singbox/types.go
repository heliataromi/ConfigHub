package singbox

// Config represents a complete Sing-box configuration file (v1.8+ / v1.10+ / v1.11+ compatible)
type Config struct {
	Log       *LogConfig               `json:"log,omitempty"`
	DNS       *DNSConfig               `json:"dns,omitempty"`
	Inbounds  []map[string]interface{} `json:"inbounds,omitempty"`
	Outbounds []map[string]interface{} `json:"outbounds"`
	Route     *RouteConfig             `json:"route,omitempty"`
}

// LogConfig defines logging settings for Sing-box
type LogConfig struct {
	Disabled  bool   `json:"disabled"`
	Level     string `json:"level"`
	Timestamp bool   `json:"timestamp"`
}

// DNSConfig represents DNS server definitions and routing rules
type DNSConfig struct {
	Servers          []DNSServer `json:"servers"`
	Rules            []DNSRule   `json:"rules,omitempty"`
	FakeIP           *FakeIP     `json:"fakeip,omitempty"`
	Strategy         string      `json:"strategy,omitempty"`
	IndependentCache bool        `json:"independent_cache,omitempty"`
}

// DNSServer defines a DNS upstream server
type DNSServer struct {
	Tag             string `json:"tag"`
	Address         string `json:"address"`
	AddressResolver string `json:"address_resolver,omitempty"`
	Detour          string `json:"detour,omitempty"`
}

// DNSRule defines a rule matching DNS queries
type DNSRule struct {
	Outbound      string   `json:"outbound,omitempty"`
	Geosite       []string `json:"geosite,omitempty"`
	DomainSuffix  []string `json:"domain_suffix,omitempty"`
	DomainKeyword []string `json:"domain_keyword,omitempty"`
	QueryType     []string `json:"query_type,omitempty"`
	Server        string   `json:"server"`
}

// FakeIP settings in Sing-box
type FakeIP struct {
	Enabled    bool   `json:"enabled"`
	Inet4Range string `json:"inet4_range,omitempty"`
	Inet6Range string `json:"inet6_range,omitempty"`
}

// RouteConfig represents route rules and final fallback outbound
type RouteConfig struct {
	AutoDetectInterface bool        `json:"auto_detect_interface"`
	Final               string      `json:"final"`
	Rules               []RouteRule `json:"rules,omitempty"`
}

// RouteRule defines a traffic routing rule
type RouteRule struct {
	Protocol      string      `json:"protocol,omitempty"`
	Port          interface{} `json:"port,omitempty"`
	GeoIP         []string    `json:"geoip,omitempty"`
	Geosite       []string    `json:"geosite,omitempty"`
	DomainSuffix  []string    `json:"domain_suffix,omitempty"`
	DomainKeyword []string    `json:"domain_keyword,omitempty"`
	Outbound      string      `json:"outbound"`
}

// ParsedProxy holds a proxy tag and its raw Sing-box outbound configuration map
type ParsedProxy struct {
	Tag  string
	Data map[string]interface{}
}
