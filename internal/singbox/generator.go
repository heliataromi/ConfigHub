package singbox

import (
	"encoding/json"
	"fmt"
	"strings"
)

// GenerateSingboxConfig converts a slice of proxy links into a full Sing-box JSON profile
func GenerateSingboxConfig(links []string) ([]byte, error) {
	var allNodeTags []string
	var rawOutbounds []map[string]interface{}
	nameCounts := make(map[string]int)

	// 1. Parse each link
	for _, link := range links {
		// Ignore dummy/timestamp configs
		if strings.Contains(link, "127.0.0.1:0") || strings.Contains(link, "Last%20update") || strings.Contains(link, "Last update") {
			continue
		}

		outbound, _, err := ParseConfigToSingbox(link)
		if err != nil || outbound == nil {
			continue
		}

		baseTag, _ := outbound["tag"].(string)
		if baseTag == "" {
			baseTag = "Proxy"
		}

		nameCounts[baseTag]++
		count := nameCounts[baseTag]
		uniqueTag := fmt.Sprintf("%s %02d", baseTag, count)
		outbound["tag"] = uniqueTag

		allNodeTags = append(allNodeTags, uniqueTag)
		rawOutbounds = append(rawOutbounds, outbound)
	}

	if len(rawOutbounds) == 0 {
		return nil, fmt.Errorf("no valid proxies to generate sing-box config")
	}

	// 2. Build Selector & URLTest Group Outbounds
	mainSelectorTags := []string{
		"⚡ AUTO (Fastest Node)",
		"🔄 FALLBACK (Failover)",
		"⚖️ LOAD-BALANCE",
		"🎯 MANUAL (All Nodes)",
	}

	var groupOutbounds []map[string]interface{}

	// Root SELECT group
	groupOutbounds = append(groupOutbounds, map[string]interface{}{
		"type":      "selector",
		"tag":       "🔰 PROXY",
		"outbounds": mainSelectorTags,
		"default":   "⚡ AUTO (Fastest Node)",
	})

	// AUTO group
	groupOutbounds = append(groupOutbounds, map[string]interface{}{
		"type":      "urltest",
		"tag":       "⚡ AUTO (Fastest Node)",
		"outbounds": allNodeTags,
		"url":       "http://www.gstatic.com/generate_204",
		"interval":  "3m",
		"tolerance": 50,
	})

	// FALLBACK group
	groupOutbounds = append(groupOutbounds, map[string]interface{}{
		"type":      "urltest",
		"tag":       "🔄 FALLBACK (Failover)",
		"outbounds": allNodeTags,
		"url":       "http://www.gstatic.com/generate_204",
		"interval":  "3m",
	})

	// LOAD-BALANCE group
	groupOutbounds = append(groupOutbounds, map[string]interface{}{
		"type":      "urltest",
		"tag":       "⚖️ LOAD-BALANCE",
		"outbounds": allNodeTags,
		"url":       "http://www.gstatic.com/generate_204",
		"interval":  "3m",
	})

	// Manual Selection group containing all individual nodes
	groupOutbounds = append(groupOutbounds, map[string]interface{}{
		"type":      "selector",
		"tag":       "🎯 MANUAL (All Nodes)",
		"outbounds": allNodeTags,
	})

	// System Outbounds
	systemOutbounds := []map[string]interface{}{
		{
			"type": "direct",
			"tag":  "direct",
		},
		{
			"type": "block",
			"tag":  "block",
		},
		{
			"type": "dns",
			"tag":  "dns-out",
		},
	}

	// 3. Assemble full Outbounds list
	var finalOutbounds []map[string]interface{}
	finalOutbounds = append(finalOutbounds, groupOutbounds...)
	finalOutbounds = append(finalOutbounds, rawOutbounds...)
	finalOutbounds = append(finalOutbounds, systemOutbounds...)

	// 4. Construct complete Sing-box configuration
	cfg := Config{
		Log: &LogConfig{
			Disabled:  false,
			Level:     "info",
			Timestamp: true,
		},
		DNS: &DNSConfig{
			Servers: []DNSServer{
				{
					Tag:             "remote-dns",
					Address:         "https://8.8.8.8/dns-query",
					AddressResolver: "local-dns",
					Detour:          "🔰 PROXY",
				},
				{
					Tag:     "local-dns",
					Address: "https://1.1.1.1/dns-query",
					Detour:  "direct",
				},
				{
					Tag:     "fakeip-dns",
					Address: "fakeip",
				},
			},
			Rules: []DNSRule{
				{
					Outbound: "any",
					Server:   "local-dns",
				},
				{
					Geosite:       []string{"ir"},
					DomainSuffix:  []string{".ir"},
					DomainKeyword: []string{"shaparak"},
					Server:        "local-dns",
				},
				{
					QueryType: []string{"A", "AAAA"},
					Server:    "fakeip-dns",
				},
			},
			FakeIP: &FakeIP{
				Enabled:    true,
				Inet4Range: "198.18.0.0/15",
				Inet6Range: "fc00::/18",
			},
			Strategy:         "prefer_ipv4",
			IndependentCache: true,
		},
		Inbounds: []map[string]interface{}{
			{
				"type":        "mixed",
				"tag":         "mixed-in",
				"listen":      "127.0.0.1",
				"listen_port": 2080,
				"sniff":       true,
			},
		},
		Outbounds: finalOutbounds,
		Route: &RouteConfig{
			AutoDetectInterface: true,
			Final:               "🔰 PROXY",
			Rules: []RouteRule{
				{
					Protocol: "dns",
					Outbound: "dns-out",
				},
				{
					Port:     53,
					Outbound: "dns-out",
				},
				{
					Geosite:  []string{"category-ads-all"},
					Outbound: "block",
				},
				{
					GeoIP:         []string{"private", "ir"},
					Geosite:       []string{"ir"},
					DomainSuffix:  []string{".ir"},
					DomainKeyword: []string{"shaparak"},
					Outbound:      "direct",
				},
			},
		},
	}

	return json.MarshalIndent(&cfg, "", "  ")
}
