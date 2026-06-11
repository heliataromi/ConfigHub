package parser

import (
    "encoding/base64"
    "encoding/json"
    "fmt"
    "net/url"
    "strings"

    "ConfigHub/internal/geoip"

    "github.com/oschwald/maxminddb-golang"
)

// RenameConfig applies the Country Flag and Channel ID to the config
func RenameConfig(rawLink, channel string, db *maxminddb.Reader) string {
    if strings.HasPrefix(rawLink, "vmess://") {
        return renameVMess(rawLink, channel, db)
    }

    // For VLESS, Trojan, SS
    return renameStandard(rawLink, channel, db)
}

func renameStandard(link, channel string, db *maxminddb.Reader) string {
    u, err := url.Parse(link)
    if err != nil {
        // Fallback if url.Parse fails
        return link
    }

    address := u.Hostname()
    if address == "" {
        // Fallback if parsing failed to find a host
        parts := strings.Split(link, "@")
        if len(parts) >= 2 {
            hostAndQuery := parts[len(parts)-1]
            address = strings.Split(hostAndQuery, ":")[0]
        }
    }

    iso, flag := geoip.GetCountry(address, db)

    newName := fmt.Sprintf("%s %s - [%s]", flag, iso, channel)

    // url.Parse automatically handles url-encoding the fragment when String() is called.
    u.Fragment = newName
    return u.String()
}

func renameVMess(link, channel string, db *maxminddb.Reader) string {
    payload := strings.TrimPrefix(link, "vmess://")

    // URL-Safe Base64 replacement
    payload = strings.ReplaceAll(payload, "-", "+")
    payload = strings.ReplaceAll(payload, "_", "/")
    if pad := len(payload) % 4; pad != 0 {
        payload += strings.Repeat("=", 4-pad)
    }

    decoded, err := base64.StdEncoding.DecodeString(payload)
    if err != nil {
        return link
    }

    var v map[string]interface{}
    if err := json.Unmarshal(decoded, &v); err != nil {
        return link
    }

    address := safeString(v["add"])
    iso, flag := geoip.GetCountry(address, db)
    
    newName := fmt.Sprintf("%s %s - [%s]", flag, iso, channel)

    // Replace the old remark with the new one
    v["ps"] = newName

    // Re-encode to Base64
    newJSON, _ := json.Marshal(v)
    newPayload := base64.StdEncoding.EncodeToString(newJSON)
    return "vmess://" + newPayload
}
