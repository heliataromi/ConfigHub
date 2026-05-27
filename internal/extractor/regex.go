package extractor

import (
    "regexp"
)

var (
    // Matches standard and URL-Safe Base64
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
    var validVmess []string

    // Filter out malformed/garbage VMess links
    rawVmess := vmessRegex.FindAllString(text, -1)
    for _, v := range rawVmess {
        // A real Base64 JSON payload for VMess is always much longer than a UUID.
        // "vmess://" is 8 chars, plus a valid base64 JSON is usually 60+ chars.
        // Dropping anything under 50 characters eliminates the isolated UUID bugs.
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
