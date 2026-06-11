package extractor

import (
    "regexp"
    "strings"
)

var (
    vmessRegex  = regexp.MustCompile(`(?:^|[^a-zA-Z])(vmess://[A-Za-z0-9+/=\-_]+)`)
    vlessRegex  = regexp.MustCompile(`(?:^|[^a-zA-Z])(vless://[^\s<"']+)`)
    trojanRegex = regexp.MustCompile(`(?:^|[^a-zA-Z])(trojan://[^\s<"']+)`)
    ssRegex     = regexp.MustCompile(`(?:^|[^a-zA-Z])(ss://[^\s<"']+)`)
)

// Configs holds the categorized extracted links
type Configs struct {
    Vmess  []string
    Vless  []string
    Trojan []string
    SS     []string
}

func extractRegex(text string, re *regexp.Regexp) []string {
    matches := re.FindAllStringSubmatch(text, -1)
    var valid []string
    for _, m := range matches {
        if len(m) > 1 {
            // m[1] contains the actual URL
            link := strings.TrimRight(m[1], ".,;!?) \r\n\t")
            if len(link) > 10 {
                valid = append(valid, link)
            }
        }
    }
    return valid
}

// Extract parses raw text and extracts supported V2Ray configs
func Extract(text string) Configs {
    var validVmess []string

    // Filter out malformed/garbage VMess links (like isolated UUIDs)
    rawVmessMatches := vmessRegex.FindAllStringSubmatch(text, -1)
    for _, m := range rawVmessMatches {
        if len(m) > 1 {
            v := strings.TrimRight(m[1], ".,;!?) \r\n\t")
            if len(v) > 50 {
                validVmess = append(validVmess, v)
            }
        }
    }

    return Configs{
        Vmess:  validVmess,
        Vless:  extractRegex(text, vlessRegex),
        Trojan: extractRegex(text, trojanRegex),
        SS:     extractRegex(text, ssRegex),
    }
}
