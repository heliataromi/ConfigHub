package parser

import (
    "encoding/base64"
    "encoding/json"
    "fmt"
    "net"
    "net/url"
    "strconv"
    "strings"
)

// GetFingerprint generates a canonical unique identifier for a config to deduplicate accurately.
func GetFingerprint(link string) string {
    link = strings.TrimSpace(link)
    if strings.HasPrefix(link, "vmess://") {
        return getVMessFingerprint(link)
    }
    if strings.HasPrefix(link, "ssr://") {
        return getSSRFingerprint(link)
    }
    if strings.HasPrefix(link, "ss://") {
        return getSSFingerprint(link)
    }
    if strings.HasPrefix(link, "tg://proxy?") || strings.HasPrefix(link, "https://t.me/proxy?") || strings.HasPrefix(link, "http://t.me/proxy?") {
        return getTelegramFingerprint(link)
    }

    return getStandardFingerprint(link)
}

// aliasMap maps various parameter alias names to their canonical form
var aliasMap = map[string]string{
    "public_key":         "pbk",
    "publickey":          "pbk",
    "publicKey":          "pbk",
    "short_id":           "sid",
    "shortid":            "sid",
    "shortId":            "sid",
    "spider_x":           "spx",
    "spiderx":            "spx",
    "spiderX":            "spx",
    "serviceName":        "servicename",
    "service_name":       "servicename",
    "allow_insecure":     "allowInsecure",
    "allowinsecure":      "allowInsecure",
    "insecure":           "allowInsecure",
    "private_key":        "privkey",
    "privatekey":         "privkey",
    "secret_key":         "privkey",
    "preshared_key":      "psk",
    "presharedkey":       "psk",
    "peer":               "sni",
    "client-fingerprint": "fp",
}

// getStandardFingerprint parses URI schemes (vless://, trojan://, tuic://, hy2://, etc.)
// and canonicalizes query parameters, casing, and default ports.
func getStandardFingerprint(link string) string {
    raw := strings.Split(link, "#")[0]

    u, err := url.Parse(raw)
    if err != nil {
        return raw
    }

    q := u.Query()
    allowed := url.Values{}

    allowedKeys := []string{
        "type", "security", "sni", "alpn", "fp", "pbk", "sid", "spx", "flow", "fm",
        "encryption", "host", "path", "headerType", "seed", "mode", "extra",
        "authority", "servicename", "obfs", "obfs-password", "obfsParam", "pinSHA256",
        "allowInsecure", "auth", "plugin", "plugin-opts", "privkey", "psk",
        "congestion_control", "udp_relay_mode", "reduce_rtt", "upmbps",
        "downmbps", "mtu", "ip", "address",
    }

    // Canonicalize query parameter keys and values
    for k, vals := range q {
        if len(vals) == 0 {
            continue
        }
        val := vals[0]
        canonicalKey := k
        if mapped, exists := aliasMap[k]; exists {
            canonicalKey = mapped
        }

        // Normalize specific values
        switch canonicalKey {
        case "type", "security", "fp", "flow", "mode", "headerType":
            val = strings.ToLower(strings.TrimSpace(val))
        case "allowInsecure":
            if val == "1" || strings.EqualFold(val, "true") {
                val = "1"
            } else {
                val = "0"
            }
        case "path":
            val = strings.TrimSpace(val)
            if val == "" && (q.Get("type") == "ws" || q.Get("type") == "xhttp") {
                val = "/"
            }
        }

        allowed.Set(canonicalKey, val)
    }

    // Filter to only compare recognized keys
    filtered := url.Values{}
    for _, key := range allowedKeys {
        if val := allowed.Get(key); val != "" {
            filtered.Set(key, val)
        }
    }

    // Host & Port normalization
    host := u.Hostname()
    port := u.Port()
    if port == "" && host != "" {
        // Assign standard default ports when omitted
        sec := strings.ToLower(filtered.Get("security"))
        if sec == "tls" || sec == "reality" || u.Scheme == "trojan" || u.Scheme == "hy2" || u.Scheme == "tuic" {
            port = "443"
        } else if u.Scheme == "socks" || u.Scheme == "socks5" {
            port = "1080"
        } else if u.Scheme == "wireguard" || u.Scheme == "wg" {
            port = "51820"
        } else {
            port = "80"
        }
    }
    if host != "" {
        u.Host = strings.ToLower(strings.Trim(host, "[]")) + ":" + port
    }

    // User normalization (UUID) for VLESS/WireGuard
    if u.User != nil {
        if u.Scheme == "vless" {
            u.User = url.User(strings.ToLower(u.User.Username()))
        }
    }

    u.RawQuery = filtered.Encode()
    u.Fragment = ""

    return u.String()
}

// getSSFingerprint normalizes both SIP002 and legacy Shadowsocks formats into a canonical fingerprint
func getSSFingerprint(link string) string {
    raw := strings.TrimPrefix(link, "ss://")
    raw = strings.Split(raw, "#")[0]

    var userInfo, hostPort, query string
    if strings.Contains(raw, "@") {
        atSplit := strings.SplitN(raw, "@", 2)
        userInfo = atSplit[0]
        if unescaped, err := url.QueryUnescape(userInfo); err == nil && unescaped != "" {
            userInfo = unescaped
        }
        rest := atSplit[1]

        if qIdx := strings.Index(rest, "?"); qIdx != -1 {
            hostPort = rest[:qIdx]
            query = rest[qIdx+1:]
        } else {
            hostPort = strings.TrimRight(rest, "/")
        }

        if decoded, err := decodeBase64Safe(userInfo); err == nil {
            decodedStr := string(decoded)
            if strings.Contains(decodedStr, ":") {
                userInfo = decodedStr
            }
        }
    } else {
        // Legacy base64(method:password@host:port)
        legacyRaw := raw
        if unescaped, err := url.QueryUnescape(legacyRaw); err == nil && unescaped != "" {
            legacyRaw = unescaped
        }
        if decoded, err := decodeBase64Safe(legacyRaw); err == nil {
            decodedStr := string(decoded)
            if atIdx := strings.LastIndex(decodedStr, "@"); atIdx != -1 {
                userInfo = decodedStr[:atIdx]
                hostPort = decodedStr[atIdx+1:]
            }
        }
    }

    userParts := strings.SplitN(userInfo, ":", 2)
    if len(userParts) < 2 {
        return getStandardFingerprint(link)
    }

    cipher := strings.ToLower(strings.TrimSpace(userParts[0]))
    password := userParts[1]
    if cipher == "chacha20-poly1305" {
        cipher = "chacha20-ietf-poly1305"
    }

    server := hostPort
    port := "8388"
    if h, p, err := net.SplitHostPort(hostPort); err == nil {
        server = h
        port = p
    } else {
        server = strings.Trim(hostPort, "[]")
    }
    server = strings.ToLower(strings.Trim(server, "[]"))

    normQuery := ""
    if query != "" {
        if q, err := url.ParseQuery(query); err == nil {
            normQuery = "?" + q.Encode()
        }
    }

    return fmt.Sprintf("ss://%s:%s@%s:%s%s", cipher, password, server, port, normQuery)
}

// getVMessFingerprint decodes the VMess Base64 JSON and creates a canonical fingerprint
func getVMessFingerprint(link string) string {
    payload := strings.TrimPrefix(link, "vmess://")
    payload = strings.TrimSpace(payload)

    decoded, err := decodeBase64Safe(payload)
    if err != nil {
        return link
    }

    var v map[string]interface{}
    if err := json.Unmarshal(decoded, &v); err != nil {
        return link
    }

    add := strings.ToLower(strings.Trim(safeString(v["add"]), "[]"))
    port := safeString(v["port"])
    if p, err := strconv.Atoi(port); err == nil && p > 0 {
        port = strconv.Itoa(p)
    } else if port == "" {
        port = "443"
    }

    id := strings.ToLower(strings.TrimSpace(safeString(v["id"])))
    netType := strings.ToLower(strings.TrimSpace(safeString(v["net"])))
    if netType == "" {
        netType = "tcp"
    }

    tls := strings.ToLower(strings.TrimSpace(safeString(v["tls"])))
    if tls == "" || tls == "none" {
        tls = "none"
    }

    path := strings.TrimSpace(safeString(v["path"]))
    if path == "" && (netType == "ws" || netType == "http" || netType == "xhttp") {
        path = "/"
    }

    host := strings.ToLower(strings.TrimSpace(safeString(v["host"])))
    sni := strings.ToLower(strings.TrimSpace(safeString(v["sni"])))
    alpn := strings.ToLower(strings.TrimSpace(safeString(v["alpn"])))
    scy := strings.ToLower(strings.TrimSpace(safeString(v["scy"])))
    if scy == "" {
        scy = "auto"
    }

    return fmt.Sprintf("vmess://%s@%s:%s?net=%s&tls=%s&path=%s&host=%s&sni=%s&alpn=%s&scy=%s",
        id, add, port, netType, tls, path, host, sni, alpn, scy)
}

// safeString converts JSON interface{} to a string safely
func safeString(v interface{}) string {
    if v == nil {
        return ""
    }
    return fmt.Sprintf("%v", v)
}

// getSSRFingerprint decodes the SSR Base64, removes remarks/group, and re-encodes
func getSSRFingerprint(link string) string {
    payload := strings.TrimPrefix(link, "ssr://")
    payload = strings.TrimSpace(payload)

    decoded, err := decodeBase64Safe(payload)
    if err != nil {
        return link
    }

    decodedStr := string(decoded)
    parts := strings.SplitN(decodedStr, "/?", 2)
    if len(parts) != 2 {
        return link
    }

    q, err := url.ParseQuery(parts[1])
    if err != nil {
        return link
    }

    allowed := url.Values{}
    allowedKeys := []string{"obfsparam", "protoparam"}
    for _, key := range allowedKeys {
        if val := q.Get(key); val != "" {
            allowed.Set(key, val)
        }
    }

    newPayload := parts[0]
    encodedQuery := allowed.Encode()
    if encodedQuery != "" {
        newPayload += "/?" + encodedQuery
    }

    newBase64 := base64.RawURLEncoding.EncodeToString([]byte(newPayload))
    return "ssr://" + newBase64
}

func getTelegramFingerprint(link string) string {
    link = strings.Split(link, "#")[0]
    u, err := url.Parse(link)
    if err != nil {
        return link
    }
    q := u.Query()
    server := strings.ToLower(strings.TrimSuffix(strings.TrimSpace(q.Get("server")), "."))
    port := strings.TrimSpace(q.Get("port"))
    secret := strings.ToLower(strings.TrimSpace(q.Get("secret")))
    return fmt.Sprintf("tg://proxy?server=%s:%s&secret=%s", server, port, secret)
}
