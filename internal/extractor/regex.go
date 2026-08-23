package extractor

import (
    "regexp"
    "strings"
)

var (
    vmessRegex    = regexp.MustCompile(`(?:^|[^a-zA-Z])(vmess://[A-Za-z0-9+/=\-_]+)`)
    vlessRegex    = regexp.MustCompile(`(?:^|[^a-zA-Z])(vless://[^\s<"']+)`)
    trojanRegex   = regexp.MustCompile(`(?:^|[^a-zA-Z])(trojan://[^\s<"']+)`)
    ssRegex       = regexp.MustCompile(`(?:^|[^a-zA-Z])(ss://[^\s<"']+)`)
    ssrRegex      = regexp.MustCompile(`(?:^|[^a-zA-Z])(ssr://[A-Za-z0-9+/=\-_]+)`)
    tuicRegex     = regexp.MustCompile(`(?:^|[^a-zA-Z])(tuic://[^\s<"']+)`)
    hy2Regex      = regexp.MustCompile(`(?:^|[^a-zA-Z])((?:hy2|hysteria2)://[^\s<"']+)`)
    hysteriaRegex = regexp.MustCompile(`(?:^|[^a-zA-Z])(hysteria://[^\s<"']+)`)
    socksRegex    = regexp.MustCompile(`(?:^|[^a-zA-Z])(socks[45]?://[^\s<"']+)`)
    wgRegex       = regexp.MustCompile(`(?:^|[^a-zA-Z])((?:wireguard|wg)://[^\s<"']+)`)
)

// trailingCutset defines characters commonly appended by markdown, HTML, or punctuation
const trailingCutset = ".,;!?) \r\n\t\"'>]`*~_\\"

// Configs holds the categorized extracted links
type Configs struct {
    Vmess     []string
    Vless     []string
    Trojan    []string
    SS        []string
    SSR       []string
    TUIC      []string
    Hy2       []string
    Hysteria  []string
    Socks     []string
    WireGuard []string
}

// Count returns the total number of valid configs across all protocols
func (c Configs) Count() int {
    return len(c.Vmess) + len(c.Vless) + len(c.Trojan) + len(c.SS) +
        len(c.SSR) + len(c.TUIC) + len(c.Hy2) + len(c.Hysteria) +
        len(c.Socks) + len(c.WireGuard)
}

func extractRegex(text string, re *regexp.Regexp) []string {
    matches := re.FindAllStringSubmatch(text, -1)
    var valid []string
    for _, m := range matches {
        if len(m) > 1 {
            // m[1] contains the actual URL
            link := strings.TrimRight(m[1], trailingCutset)
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
            v := strings.TrimRight(m[1], trailingCutset)
            if len(v) > 50 {
                validVmess = append(validVmess, v)
            }
        }
    }

    return Configs{
        Vmess:     validVmess,
        Vless:     extractRegex(text, vlessRegex),
        Trojan:    extractRegex(text, trojanRegex),
        SS:        extractRegex(text, ssRegex),
        SSR:       extractRegex(text, ssrRegex),
        TUIC:      extractRegex(text, tuicRegex),
        Hy2:       extractRegex(text, hy2Regex),
        Hysteria:  extractRegex(text, hysteriaRegex),
        Socks:     extractRegex(text, socksRegex),
        WireGuard: extractRegex(text, wgRegex),
    }
}

// AuditAndExtract extracts configs while scanning and logging dropped candidate lines to telemetry
func AuditAndExtract(text, source string, recordDropped func(source, candidate, reason string)) Configs {
    configs := Extract(text)

    if recordDropped == nil {
        return configs
    }

    // Scan lines to audit dropped/unsupported candidate URLs
    lines := strings.Split(text, "\n")
    for _, line := range lines {
        line = strings.TrimSpace(line)
        if strings.Contains(line, "://") && len(line) > 8 {
            matched := strings.HasPrefix(line, "vless://") || strings.HasPrefix(line, "vmess://") ||
                strings.HasPrefix(line, "trojan://") || strings.HasPrefix(line, "ss://") ||
                strings.HasPrefix(line, "ssr://") || strings.HasPrefix(line, "tuic://") ||
                strings.HasPrefix(line, "hy2://") || strings.HasPrefix(line, "hysteria://") ||
                strings.HasPrefix(line, "hysteria2://") || strings.HasPrefix(line, "socks://") ||
                strings.HasPrefix(line, "socks5://") || strings.HasPrefix(line, "wireguard://") ||
                strings.HasPrefix(line, "wg://")

            if !matched {
                // Suspicious link with unsupported or custom scheme
                recordDropped(source, line, "unsupported or unrecognized proxy scheme")
            } else if strings.HasPrefix(line, "vmess://") && len(line) <= 50 {
                recordDropped(source, line, "vmess payload too short or malformed")
            }
        }
    }

    return configs
}
