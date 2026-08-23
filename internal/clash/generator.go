package clash

import (
	"fmt"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

// GenerateClashConfig converts a slice of proxy links to a full Clash Meta YAML string
func GenerateClashConfig(links []string) ([]byte, error) {
	var parsedProxies []ParsedProxy
	nameCounts := make(map[string]int)

	// 1. Parse each link
	for _, link := range links {
		// Ignore dummy/timestamp configs
		if strings.Contains(link, "127.0.0.1:0") || strings.Contains(link, "Last%20update") || strings.Contains(link, "Last update") {
			continue
		}

		proxyData, _, err := ParseConfigToClash(link)
		if err != nil || proxyData == nil {
			continue
		}

		baseName, _ := proxyData["name"].(string)
		if baseName == "" {
			baseName = "Proxy"
		}

		nameCounts[baseName]++
		count := nameCounts[baseName]
		uniqueName := fmt.Sprintf("%s %02d", baseName, count)
		proxyData["name"] = uniqueName

		parsedProxies = append(parsedProxies, ParsedProxy{
			Name: uniqueName,
			Data: proxyData,
		})
	}

	if len(parsedProxies) == 0 {
		return nil, fmt.Errorf("no valid proxies to generate clash config")
	}

	// 2. Collect all node names
	var allNodeNames []string
	var rawProxies []map[string]interface{}

	for _, p := range parsedProxies {
		allNodeNames = append(allNodeNames, p.Name)
		rawProxies = append(rawProxies, p.Data)
	}

	// 3. Build Primary Proxy Groups
	mainSelectorProxies := []string{
		"⚡ AUTO (Fastest Node)",
		"🔄 FALLBACK (Failover)",
		"⚖️ LOAD-BALANCE",
		"🎯 MANUAL (All Nodes)",
	}

	var proxyGroups []ProxyGroup

	// Root SELECT group
	proxyGroups = append(proxyGroups, ProxyGroup{
		Name:    "🔰 PROXY",
		Type:    "select",
		Proxies: mainSelectorProxies,
	})

	// AUTO group
	proxyGroups = append(proxyGroups, ProxyGroup{
		Name:      "⚡ AUTO (Fastest Node)",
		Type:      "url-test",
		URL:       "http://www.gstatic.com/generate_204",
		Interval:  300,
		Tolerance: 50,
		Proxies:   allNodeNames,
	})

	// FALLBACK group
	proxyGroups = append(proxyGroups, ProxyGroup{
		Name:     "🔄 FALLBACK (Failover)",
		Type:     "fallback",
		URL:      "http://www.gstatic.com/generate_204",
		Interval: 300,
		Proxies:  allNodeNames,
	})

	// LOAD-BALANCE group
	proxyGroups = append(proxyGroups, ProxyGroup{
		Name:     "⚖️ LOAD-BALANCE",
		Type:     "load-balance",
		URL:      "http://www.gstatic.com/generate_204",
		Interval: 300,
		Strategy: "consistent-hashing",
		Proxies:  allNodeNames,
	})

	// Manual Selection group containing all individual nodes
	proxyGroups = append(proxyGroups, ProxyGroup{
		Name:    "🎯 MANUAL (All Nodes)",
		Type:    "select",
		Proxies: allNodeNames,
	})

	// 5. Construct full Config object
	cfg := Config{
		Port:               7890,
		SocksPort:          7891,
		MixedPort:          7890,
		AllowLan:           true,
		Mode:               "rule",
		LogLevel:           "info",
		IPv6:               false,
		ExternalController: "127.0.0.1:9090",
		DNS: DNSConfig{
			Enable:       true,
			IPv6:         false,
			Listen:       "0.0.0.0:1053",
			EnhancedMode: "fake-ip",
			FakeIPRange:  "198.18.0.1/16",
			Nameserver: []string{
				"8.8.8.8",
				"1.1.1.1",
				"119.29.29.29",
				"https://dns.google/dns-query",
				"https://1.1.1.1/dns-query",
			},
		},
		Proxies:     rawProxies,
		ProxyGroups: proxyGroups,
		Rules: []string{
			"GEOIP,IR,DIRECT",
			"DOMAIN-SUFFIX,ir,DIRECT",
			"DOMAIN-KEYWORD,shaparak,DIRECT",
			"GEOSITE,category-ads-all,REJECT",
			"MATCH,🔰 PROXY",
		},
	}

	out, err := yaml.Marshal(cfg)
	if err != nil {
		return nil, err
	}

	return unescapeUnicodeYAML(out), nil
}

func unescapeUnicodeYAML(data []byte) []byte {
	// Unescape \U0001xxxx and \uxxxx in yaml output so emojis and flags are clean UTF-8
	str := string(data)
	var sb strings.Builder
	sb.Grow(len(str))

	for i := 0; i < len(str); {
		if str[i] == '\\' && i+1 < len(str) {
			if str[i+1] == 'U' && i+10 <= len(str) {
				hexStr := str[i+2 : i+10]
				if n, err := strconv.ParseUint(hexStr, 16, 32); err == nil {
					sb.WriteRune(rune(n))
					i += 10
					continue
				}
			} else if str[i+1] == 'u' && i+6 <= len(str) {
				hexStr := str[i+2 : i+6]
				if n, err := strconv.ParseUint(hexStr, 16, 16); err == nil {
					sb.WriteRune(rune(n))
					i += 6
					continue
				}
			}
		}
		sb.WriteByte(str[i])
		i++
	}

	return []byte(sb.String())
}
