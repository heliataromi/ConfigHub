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
    if strings.HasPrefix(link, "ssr://") {
        return getSSRFingerprint(link)
    }

    return getStandardFingerprint(link)
}

// getStandardFingerprint parses URI schemes (vless://, trojan://, ss://)
// and normalizes them by keeping only allowed query parameters.
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
	allowed := url.Values{}

	// List of parameters strictly compared by v2rayN for standard configs
	allowedKeys := []string{
		"type", "security", "sni", "alpn", "fp", "pbk", "public_key", "publicKey",
		"sid", "short_id", "shortId", "spx", "spider_x", "spiderX", "flow", "fm",
		"encryption", "host", "path", "headerType", "seed", "mode", "extra",
		"authority", "serviceName", "servicename",
		// Parameters for Hy2, TUIC, WireGuard, SS plugins, etc.
		"obfs", "obfs-password", "obfsParam", "pinSHA256", "allow_insecure", "allowInsecure",
		"insecure", "peer", "auth", "plugin", "plugin-opts",
		"private_key", "privkey", "privatekey", "preshared_key", "presharedkey", "psk",
		"congestion_control", "udp_relay_mode", "reduce_rtt", "upmbps",
		"downmbps", "mtu", "ip", "address",
	}

	for _, key := range allowedKeys {
		if val := q.Get(key); val != "" {
			allowed.Set(key, val)
		}
	}

	// Normalize Host casing
	u.Host = strings.ToLower(u.Host)

	// Normalize User (UUID) for VLESS/VMESS if present, though Trojan might use it for password.
	// To be safe, we only lowercase host to avoid merging case-sensitive Trojan/SS passwords.
	if u.Scheme == "vless" && u.User != nil {
		u.User = url.User(strings.ToLower(u.User.Username()))
	}

	// u.Query().Encode() automatically sorts the keys alphabetically!
	u.RawQuery = allowed.Encode()
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

	// Safely extract core networking fields and normalize casing where applicable
	add := strings.ToLower(safeString(v["add"]))
	port := safeString(v["port"])
	id := strings.ToLower(safeString(v["id"]))
	net := strings.ToLower(safeString(v["net"]))
	tls := strings.ToLower(safeString(v["tls"]))
	path := safeString(v["path"])
	host := strings.ToLower(safeString(v["host"]))
	sni := safeString(v["sni"])
	alpn := safeString(v["alpn"])
	scy := safeString(v["scy"])

	// Create a unique fingerprint based on connection details compared by v2rayN
	return fmt.Sprintf("vmess://%s@%s:%s?net=%s&tls=%s&path=%s&host=%s&sni=%s&alpn=%s&scy=%s", 
		id, add, port, net, tls, path, host, sni, alpn, scy)
}

// safeString converts JSON interface{} to a string safely, as JSON numbers (like port) might parse as floats
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

	// Fix standard Base64 URL encoding and padding issues
	payload = strings.ReplaceAll(payload, "-", "+")
	payload = strings.ReplaceAll(payload, "_", "/")
	if pad := len(payload) % 4; pad != 0 {
		payload += strings.Repeat("=", 4-pad)
	}

	decoded, err := base64.StdEncoding.DecodeString(payload)
	if err != nil {
		return link // Fallback if invalid base64
	}

	decodedStr := string(decoded)
	parts := strings.SplitN(decodedStr, "/?", 2)
	if len(parts) != 2 {
		return link // No query params means no remarks to strip
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
