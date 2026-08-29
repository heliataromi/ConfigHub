package parser

import (
    "encoding/base64"
    "encoding/json"
    "fmt"
    "net"
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
    if strings.HasPrefix(rawLink, "ssr://") {
        return renameSSR(rawLink, channel, db)
    }
    if strings.HasPrefix(rawLink, "ss://") {
        return renameSS(rawLink, channel, db)
    }
    if strings.HasPrefix(rawLink, "tg://proxy?") || strings.HasPrefix(rawLink, "https://t.me/proxy?") || strings.HasPrefix(rawLink, "http://t.me/proxy?") ||
        strings.HasPrefix(rawLink, "tg://socks?") || strings.HasPrefix(rawLink, "https://t.me/socks?") || strings.HasPrefix(rawLink, "http://t.me/socks?") ||
        strings.HasPrefix(rawLink, "tg://http?") || strings.HasPrefix(rawLink, "https://t.me/http?") || strings.HasPrefix(rawLink, "http://t.me/http?") {
        return renameTelegram(rawLink, channel, db)
    }
    if strings.HasPrefix(rawLink, "cottendns://") {
        return renameWhiteDNS(rawLink, "cottendns", channel, db)
    }
    if strings.HasPrefix(rawLink, "stormdns://") {
        return renameWhiteDNS(rawLink, "stormdns", channel, db)
    }

    // For VLESS, Trojan, TUIC, Hy2, Hysteria, Socks, WireGuard, Juicity, Naive, AnyTLS, Snell, HTTP
    return renameStandard(rawLink, channel, db)
}

func renameStandard(link, channel string, db *maxminddb.Reader) string {
    u, err := url.Parse(link)
    if err != nil {
        return link
    }

    address := u.Hostname()
    if address == "" {
        parts := strings.Split(link, "@")
        if len(parts) >= 2 {
            hostAndQuery := parts[len(parts)-1]
            if h, _, err := net.SplitHostPort(hostAndQuery); err == nil {
                address = h
            } else {
                address = strings.Split(hostAndQuery, ":")[0]
            }
        }
    }

    address = strings.Trim(address, "[]")
    iso, flag := geoip.GetCountry(address, db)
    newName := fmt.Sprintf("%s %s - [%s]", flag, iso, channel)

    u.Fragment = newName
    return u.String()
}

func renameSS(link, channel string, db *maxminddb.Reader) string {
    raw := strings.TrimPrefix(link, "ss://")
    parts := strings.SplitN(raw, "#", 2)
    mainPart := parts[0]

    address := ""
    if strings.Contains(mainPart, "@") {
        // SIP002 format: base64(method:pass)@host:port
        atSplit := strings.SplitN(mainPart, "@", 2)
        userInfo := atSplit[0]
        if unescaped, err := url.QueryUnescape(userInfo); err == nil && unescaped != "" {
            userInfo = unescaped
        }
        rest := atSplit[1]
        if qIdx := strings.Index(rest, "?"); qIdx != -1 {
            rest = rest[:qIdx]
        }
        rest = strings.TrimRight(rest, "/")
        if h, _, err := net.SplitHostPort(rest); err == nil {
            address = h
        } else {
            address = rest
        }
    } else {
        // Legacy format: base64(method:pass@host:port)
        legacyPart := mainPart
        if unescaped, err := url.QueryUnescape(legacyPart); err == nil && unescaped != "" {
            legacyPart = unescaped
        }
        if decoded, err := decodeBase64Safe(legacyPart); err == nil {
            decodedStr := string(decoded)
            if atIdx := strings.LastIndex(decodedStr, "@"); atIdx != -1 {
                hostPort := decodedStr[atIdx+1:]
                if h, _, err := net.SplitHostPort(hostPort); err == nil {
                    address = h
                } else {
                    address = hostPort
                }
            }
        }
    }

    address = strings.Trim(address, "[]")
    iso, flag := geoip.GetCountry(address, db)
    newName := fmt.Sprintf("%s %s - [%s]", flag, iso, channel)

    return fmt.Sprintf("ss://%s#%s", mainPart, url.QueryEscape(newName))
}

func renameSSR(link, channel string, db *maxminddb.Reader) string {
    payload := strings.TrimPrefix(link, "ssr://")
    decoded, err := decodeBase64Safe(payload)
    if err != nil {
        return link
    }

    str := string(decoded)
    parts := strings.SplitN(str, "/?", 2)
    mainParts := strings.Split(parts[0], ":")
    if len(mainParts) < 6 {
        return link
    }

    server := strings.Trim(mainParts[0], "[]")
    iso, flag := geoip.GetCountry(server, db)
    newName := fmt.Sprintf("%s %s - [%s]", flag, iso, channel)
    newRemarksBase64 := base64.RawURLEncoding.EncodeToString([]byte(newName))

    q := url.Values{}
    if len(parts) > 1 {
        if parsedQuery, err := url.ParseQuery(parts[1]); err == nil {
            q = parsedQuery
        }
    }
    q.Set("remarks", newRemarksBase64)

    newPayload := parts[0] + "/?" + q.Encode()
    return "ssr://" + base64.RawURLEncoding.EncodeToString([]byte(newPayload))
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

    address := strings.Trim(safeString(v["add"]), "[]")
    iso, flag := geoip.GetCountry(address, db)
    
    newName := fmt.Sprintf("%s %s - [%s]", flag, iso, channel)

    // Replace the old remark with the new one
    v["ps"] = newName

    // Re-encode to Base64
    newJSON, _ := json.Marshal(v)
    newPayload := base64.StdEncoding.EncodeToString(newJSON)
    return "vmess://" + newPayload
}

func decodeBase64Safe(s string) ([]byte, error) {
    s = strings.TrimSpace(s)
    s = strings.ReplaceAll(s, "-", "+")
    s = strings.ReplaceAll(s, "_", "/")
    if pad := len(s) % 4; pad != 0 {
        s += strings.Repeat("=", 4-pad)
    }
    return base64.StdEncoding.DecodeString(s)
}

func renameTelegram(link, channel string, db *maxminddb.Reader) string {
    link = strings.Split(link, "#")[0]
    u, err := url.Parse(link)
    if err != nil {
        return link
    }
    q := u.Query()
    server := strings.TrimSuffix(strings.TrimSpace(q.Get("server")), ".")
    server = strings.Trim(server, "[]")

    isSocks := strings.HasPrefix(link, "tg://socks?") || strings.HasPrefix(link, "https://t.me/socks?") || strings.HasPrefix(link, "http://t.me/socks?")
    isHTTP := strings.HasPrefix(link, "tg://http?") || strings.HasPrefix(link, "https://t.me/http?") || strings.HasPrefix(link, "http://t.me/http?")

    port := q.Get("port")
    if port == "" {
        if isSocks {
            port = "1080"
        } else {
            port = "443"
        }
    }

    secret := q.Get("secret")
    if secret != "" {
        return fmt.Sprintf("tg://proxy?server=%s&port=%s&secret=%s", server, port, secret)
    }

    user := q.Get("user")
    pass := q.Get("pass")
    authParams := ""
    if user != "" {
        authParams += fmt.Sprintf("&user=%s", url.QueryEscape(user))
    }
    if pass != "" {
        authParams += fmt.Sprintf("&pass=%s", url.QueryEscape(pass))
    }

    if isHTTP {
        return fmt.Sprintf("tg://http?server=%s&port=%s%s", server, port, authParams)
    }

    return fmt.Sprintf("tg://socks?server=%s&port=%s%s", server, port, authParams)
}

func renameWhiteDNS(rawLink, scheme, channel string, db *maxminddb.Reader) string {
    payload := strings.TrimPrefix(rawLink, scheme+"://")
    decoded, err := decodeBase64Safe(payload)
    if err != nil {
        return rawLink
    }

    var data map[string]interface{}
    if err := json.Unmarshal(decoded, &data); err != nil {
        return rawLink
    }

    profile, ok := data["profile"].(map[string]interface{})
    if !ok {
        return rawLink
    }
    server, ok := profile["server"].(map[string]interface{})
    if !ok {
        return rawLink
    }

    domainStr := safeString(server["domain"])
    firstDomain := strings.TrimSpace(strings.Split(domainStr, ",")[0])
    firstDomain = strings.Trim(firstDomain, "[]")

    iso, flag := geoip.GetCountry(firstDomain, db)
    newName := fmt.Sprintf("%s %s - [%s]", flag, iso, channel)
    profile["name"] = newName

    newJSON, err := json.Marshal(data)
    if err != nil {
        return rawLink
    }
    return scheme + "://" + base64.StdEncoding.EncodeToString(newJSON)
}
