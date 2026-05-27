package extractor

import (
	"regexp"
)

var (
	vmessRegex  = regexp.MustCompile(`vmess://[A-Za-z0-9+/=\-_]+`)
	vlessRegex  = regexp.MustCompile(`vless://[^\s<"']+`)
	trojanRegex = regexp.MustCompile(`trojan://[^\s<"']+`)
	ssRegex     = regexp.MustCompile(`ss://[^\s<"']+`)
)

// Configs holds the categorized extracted links
type Configs struct {
	Vmess  []string
	Vless  []string
	Trojan []string
	SS     []string
}

// Extract parses raw text and extracts supported V2Ray configs
func Extract(text string) Configs {
	return Configs{
		Vmess:  vmessRegex.FindAllString(text, -1),
		Vless:  vlessRegex.FindAllString(text, -1),
		Trojan: trojanRegex.FindAllString(text, -1),
		SS:     ssRegex.FindAllString(text, -1),
	}
}
