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

var anySchemeRegex = regexp.MustCompile(`([a-zA-Z0-9+.-]+://[^\s<"']+)`)

// AuditAndExtract extracts configs while scanning and logging dropped candidate URLs to telemetry
func AuditAndExtract(text, source string, recordDropped func(source, candidate, reason string)) Configs {
    configs := Extract(text)

    if recordDropped == nil {
        return configs
    }

    // Build a lookup set of all successfully extracted configs
    extractedSet := make(map[string]bool)
    for _, l := range configs.Vmess { extractedSet[l] = true }
    for _, l := range configs.Vless { extractedSet[l] = true }
    for _, l := range configs.Trojan { extractedSet[l] = true }
    for _, l := range configs.SS { extractedSet[l] = true }
    for _, l := range configs.SSR { extractedSet[l] = true }
    for _, l := range configs.TUIC { extractedSet[l] = true }
    for _, l := range configs.Hy2 { extractedSet[l] = true }
    for _, l := range configs.Hysteria { extractedSet[l] = true }
    for _, l := range configs.Socks { extractedSet[l] = true }
    for _, l := range configs.WireGuard { extractedSet[l] = true }

    // Scan for candidate URLs in text
    matches := anySchemeRegex.FindAllStringSubmatch(text, -1)
    for _, m := range matches {
        if len(m) > 1 {
            candidate := strings.TrimRight(m[1], trailingCutset)
            if len(candidate) <= 8 {
                continue
            }

            // If it was already successfully extracted, don't report as dropped
            if extractedSet[candidate] {
                continue
            }

            // Check if it starts with a known supported protocol
            isKnownProto := strings.HasPrefix(candidate, "vless://") || strings.HasPrefix(candidate, "vmess://") ||
                strings.HasPrefix(candidate, "trojan://") || strings.HasPrefix(candidate, "ss://") ||
                strings.HasPrefix(candidate, "ssr://") || strings.HasPrefix(candidate, "tuic://") ||
                strings.HasPrefix(candidate, "hy2://") || strings.HasPrefix(candidate, "hysteria://") ||
                strings.HasPrefix(candidate, "hysteria2://") || strings.HasPrefix(candidate, "socks://") ||
                strings.HasPrefix(candidate, "socks5://") || strings.HasPrefix(candidate, "wireguard://") ||
                strings.HasPrefix(candidate, "wg://")

            if isKnownProto {
                if strings.HasPrefix(candidate, "vmess://") && len(candidate) <= 50 {
                    recordDropped(source, candidate, "vmess payload too short or malformed")
                }
            } else {
                recordDropped(source, candidate, "unsupported or unrecognized proxy scheme")
            }
        }
    }

    return configs
}
