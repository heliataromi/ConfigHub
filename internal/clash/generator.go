package clash

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

var countryNames = map[string]string{
	"DE": "🇩🇪 Germany",
	"US": "🇺🇸 United States",
	"FR": "🇫🇷 France",
	"NL": "🇳🇱 Netherlands",
	"GB": "🇬🇧 United Kingdom",
	"UK": "🇬🇧 United Kingdom",
	"TR": "🇹🇷 Turkey",
	"CA": "🇨🇦 Canada",
	"FI": "🇫🇮 Finland",
	"PL": "🇵🇱 Poland",
	"SE": "🇸🇪 Sweden",
	"CH": "🇨🇭 Switzerland",
	"IT": "🇮🇹 Italy",
	"ES": "🇪🇸 Spain",
	"RU": "🇷🇺 Russia",
	"SG": "🇸🇬 Singapore",
	"JP": "🇯🇵 Japan",
	"KR": "🇰🇷 South Korea",
	"HK": "🇭🇰 Hong Kong",
	"AE": "🇦🇪 UAE",
	"IN": "🇮🇳 India",
	"AT": "🇦🇹 Austria",
	"BE": "🇧🇪 Belgium",
	"NO": "🇳🇴 Norway",
	"DK": "🇩🇰 Denmark",
	"RO": "🇷🇴 Romania",
	"BG": "🇧🇬 Bulgaria",
	"CZ": "🇨🇿 Czechia",
	"UA": "🇺🇦 Ukraine",
	"IR": "🇮🇷 Iran",
	"AM": "🇦🇲 Armenia",
	"GE": "🇬🇪 Georgia",
	"AZ":       "🇦🇿 Azerbaijan",
	"KZ":       "🇰🇿 Kazakhstan",
	"IL":       "🇮🇱 Israel",
	"AU":       "🇦🇺 Australia",
	"BR":       "🇧🇷 Brazil",
	"ZA":       "🇿🇦 South Africa",
	"CDN":      "☁️ CDN & Cloud Relay",
	"IR-RELAY": "🇮🇷 Iran Domestic Relay",
}

// GenerateClashConfig converts a slice of proxy links to a full Clash Meta YAML string
func GenerateClashConfig(links []string) ([]byte, error) {
	var parsedProxies []ParsedProxy
	countryMap := make(map[string][]string)
	nameCounts := make(map[string]int)

	// 1. Parse each link
	for _, link := range links {
		// Ignore dummy/timestamp configs
		if strings.Contains(link, "127.0.0.1:0") || strings.Contains(link, "Last%20update") || strings.Contains(link, "Last update") {
			continue
		}

		proxyData, country, err := ParseConfigToClash(link)
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
			Name:        uniqueName,
			CountryCode: country,
			Data:        proxyData,
		})
	}

	if len(parsedProxies) == 0 {
		return nil, fmt.Errorf("no valid proxies to generate clash config")
	}

	// 2. Collect all node names and group by country
	var allNodeNames []string
	var rawProxies []map[string]interface{}

	for _, p := range parsedProxies {
		allNodeNames = append(allNodeNames, p.Name)
		rawProxies = append(rawProxies, p.Data)

		if p.CountryCode != "" {
			countryMap[p.CountryCode] = append(countryMap[p.CountryCode], p.Name)
		}
	}

	// 3. Build Country Auto Groups
	var countryGroups []ProxyGroup
	var countryGroupNames []string

	// Sort country codes for deterministic output
	var countryCodes []string
	for code := range countryMap {
		countryCodes = append(countryCodes, code)
	}
	sort.Strings(countryCodes)

	for _, code := range countryCodes {
		nodes := countryMap[code]
		// Only create country auto-groups for countries with at least 2 nodes
		if len(nodes) < 2 {
			continue
		}

		cName, ok := countryNames[code]
		if !ok {
			cName = fmt.Sprintf("🌐 %s", code)
		}
		groupName := fmt.Sprintf("%s (Auto)", cName)

		countryGroups = append(countryGroups, ProxyGroup{
			Name:      groupName,
			Type:      "url-test",
			URL:       "http://www.gstatic.com/generate_204",
			Interval:  300,
			Tolerance: 50,
			Proxies:   nodes,
		})
		countryGroupNames = append(countryGroupNames, groupName)
	}

	// 4. Build Primary Proxy Groups
	var mainSelectorProxies []string
	mainSelectorProxies = append(mainSelectorProxies, "⚡ AUTO (Fastest Node)", "🔄 FALLBACK (Failover)", "⚖️ LOAD-BALANCE")
	mainSelectorProxies = append(mainSelectorProxies, countryGroupNames...)
	mainSelectorProxies = append(mainSelectorProxies, "🎯 MANUAL (All Nodes)")

	var proxyGroups []ProxyGroup

	// Root SELECT group (clean, high-level menu)
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

	// Add country groups
	proxyGroups = append(proxyGroups, countryGroups...)

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
