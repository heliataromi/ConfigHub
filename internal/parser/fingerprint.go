package parser

import (
    "encoding/base64"
    "encoding/json"
    "fmt"
    "net/url"
    "strings"
)

// GetFingerprint generates a unique string for a config, ignoring its remark/name.
func GetFingerprint(link string) string {
    if strings.HasPrefix(link, "vmess://") {
        return getVMessFingerprint(link)
    }

    return getStandardFingerprint(link)
}

// getStandardFingerprint parses URI schemes (vless://, trojan://, ss://)
// and normalizes them by sorting query parameters and dropping useless ones.
func getStandardFingerprint(link string) string {
    // Strip the remark first
    raw := strings.Split(link, "#")[0]

    u, err := url.Parse(raw)
    if err != nil {
        // If URL parsing fails, fallback to raw string
        return raw
    }

    // Extract query parameters
    q := u.Query()

    // Delete parameters that are just defaults and cause false-duplicates
    q.Del("insecure")
    q.Del("allowInsecure")

    // u.Query().Encode() automatically sorts the keys alphabetically!
    u.RawQuery = q.Encode()
    u.Fragment = "" // Ensure the remark is completely gone

    // Return the standardized URL string as the fingerprint
    return u.String()
}

// getVMessFingerprint decodes the VMess Base64 JSON and creates a unique identifier
func getVMessFingerprint(link string) string {
    payload := strings.TrimPrefix(link, "vmess://")
    payload = strings.TrimSpace(payload)

    // Fix standard Base64 URL encoding and padding issues that are common in scraped configs
    payload = strings.ReplaceAll(payload, "-", "+")
    payload = strings.ReplaceAll(payload, "_", "/")
    if pad := len(payload) % 4; pad != 0 {
        payload += strings.Repeat("=", 4-pad)
    }

    decoded, err := base64.StdEncoding.DecodeString(payload)
    if err != nil {
        // If it fails to decode for some reason, fallback to the raw link
        return link
    }

    var v map[string]interface{}
    if err := json.Unmarshal(decoded, &v); err != nil {
        return link
    }

    // Safely extract core networking fields
    add := safeString(v["add"])
    port := safeString(v["port"])
    id := safeString(v["id"])
    net := safeString(v["net"])
    tls := safeString(v["tls"])
    path := safeString(v["path"])
    host := safeString(v["host"])

    // Create a unique fingerprint based on connection details, completely ignoring the "ps" (remark) field
    return fmt.Sprintf("vmess://%s@%s:%s?net=%s&tls=%s&path=%s&host=%s", id, add, port, net, tls, path, host)
}

// safeString converts JSON interface{} to a string safely, as JSON numbers (like port) might parse as floats
func safeString(v interface{}) string {
    if v == nil {
        return ""
    }
    return fmt.Sprintf("%v", v)
}
