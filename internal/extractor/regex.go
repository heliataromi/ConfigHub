package extractor

import (
    "regexp"
)

var (
    vmessRegex  = regexp.MustCompile(`\bvmess://[A-Za-z0-9+/=\-_]+`)
    vlessRegex  = regexp.MustCompile(`\bvless://[^\s<"']+`)
    trojanRegex = regexp.MustCompile(`\btrojan://[^\s<"']+`)
    ssRegex     = regexp.MustCompile(`\bss://[^\s<"']+`)
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
    var validVmess []string

    // Filter out malformed/garbage VMess links (like isolated UUIDs)
    rawVmess := vmessRegex.FindAllString(text, -1)
    for _, v := range rawVmess {
        if len(v) > 50 {
            validVmess = append(validVmess, v)
        }
    }

    return Configs{
        Vmess:  validVmess,
        Vless:  vlessRegex.FindAllString(text, -1),
        Trojan: trojanRegex.FindAllString(text, -1),
        SS:     ssRegex.FindAllString(text, -1),
    }
}
