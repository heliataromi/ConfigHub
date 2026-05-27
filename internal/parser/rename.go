package parser

import (
    "encoding/base64"
    "encoding/json"
    "fmt"
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
    // e.g., vless://uuid@1.2.3.4:443?type=ws#old-remark
    parts := strings.Split(link, "@")
    if len(parts) < 2 {
        return link
    }

    // Extract everything after @
    hostAndQuery := parts[1]

    // The address is the part before the colon (port)
    address := strings.Split(hostAndQuery, ":")[0]

    // Get Geo Data
    iso, flag := geoip.GetCountry(address, db)
    newName := fmt.Sprintf("%s %s - [%s]", flag, iso, channel)

    // Strip the old remark (anything after #) and append the new one
    baseLink := strings.Split(link, "#")[0]
    return fmt.Sprintf("%s#%s", baseLink, newName)
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
