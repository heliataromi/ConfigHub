package clash

// Config represents a complete Clash / Clash Meta (Mihomo) configuration file
type Config struct {
	Port               int                      `yaml:"port,omitempty"`
	SocksPort          int                      `yaml:"socks-port,omitempty"`
	MixedPort          int                      `yaml:"mixed-port,omitempty"`
	AllowLan           bool                     `yaml:"allow-lan"`
	Mode               string                   `yaml:"mode"`
	LogLevel           string                   `yaml:"log-level"`
	IPv6               bool                     `yaml:"ipv6"`
	ExternalController string                   `yaml:"external-controller,omitempty"`
	DNS                DNSConfig                `yaml:"dns"`
	Proxies            []map[string]interface{} `yaml:"proxies"`
	ProxyGroups        []ProxyGroup             `yaml:"proxy-groups"`
	Rules              []string                 `yaml:"rules"`
}

// DNSConfig represents Clash DNS settings
type DNSConfig struct {
	Enable       bool     `yaml:"enable"`
	IPv6         bool     `yaml:"ipv6"`
	Listen       string   `yaml:"listen,omitempty"`
	EnhancedMode string   `yaml:"enhanced-mode"`
	FakeIPRange  string   `yaml:"fake-ip-range,omitempty"`
	Nameserver   []string `yaml:"nameserver"`
}

// ProxyGroup represents a Clash proxy group (select, url-test, fallback, load-balance)
type ProxyGroup struct {
	Name      string   `yaml:"name"`
	Type      string   `yaml:"type"`
	URL       string   `yaml:"url,omitempty"`
	Interval  int      `yaml:"interval,omitempty"`
	Tolerance int      `yaml:"tolerance,omitempty"`
	Strategy  string   `yaml:"strategy,omitempty"`
	Proxies   []string `yaml:"proxies"`
}

// ParsedProxy holds the raw map for Clash YAML and the detected country code
type ParsedProxy struct {
	Name        string
	CountryCode string
	Data        map[string]interface{}
}
