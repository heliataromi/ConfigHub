package validator

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"ConfigHub/internal/extractor"
)

var (
	uuidRegex = regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)
	hexRegex  = regexp.MustCompile(`^[0-9a-fA-F]{32}$`)
)

// ValidateConfig checks if a proxy link is structurally valid and not a dummy/broken configuration.
// Returns (true, "") if valid, or (false, rejectReason) if rejected.
func ValidateConfig(link string) (bool, string) {
	link = strings.TrimSpace(link)
	if link == "" {
		return false, "empty config link"
	}

	// 1. Check for channel banner / dummy timestamp configs
	if isBannerOrTimestamp(link) {
		return false, "channel banner or timestamp announcement node"
	}

	// 2. Validate by protocol scheme
	if strings.HasPrefix(link, "vless://") {
		return validateVLESS(link)
	} else if strings.HasPrefix(link, "vmess://") {
		return validateVMess(link)
	} else if strings.HasPrefix(link, "trojan://") {
		return validateTrojan(link)
	} else if strings.HasPrefix(link, "ss://") {
		return validateSS(link)
	} else if strings.HasPrefix(link, "ssr://") {
		return validateSSR(link)
	} else if strings.HasPrefix(link, "tuic://") {
		return validateTUIC(link)
	} else if strings.HasPrefix(link, "hy2://") || strings.HasPrefix(link, "hysteria2://") || strings.HasPrefix(link, "hysteria://") {
		return validateHy2(link)
	} else if strings.HasPrefix(link, "socks://") || strings.HasPrefix(link, "socks5://") {
		return validateSocks(link)
	} else if strings.HasPrefix(link, "wireguard://") || strings.HasPrefix(link, "wg://") {
		return validateWireGuard(link)
	} else if strings.HasPrefix(link, "juicity://") {
		return validateJuicity(link)
	} else if strings.HasPrefix(link, "naive+https://") || strings.HasPrefix(link, "naive+http://") {
		return validateNaive(link)
	} else if strings.HasPrefix(link, "tg://proxy?") || strings.HasPrefix(link, "https://t.me/proxy?") || strings.HasPrefix(link, "http://t.me/proxy?") {
		return validateTelegram(link)
	} else if strings.HasPrefix(link, "anytls://") {
		return validateAnyTLS(link)
	} else if strings.HasPrefix(link, "snell://") {
		return validateSnell(link)
	} else if strings.HasPrefix(link, "http://") || strings.HasPrefix(link, "https://") {
		return validateHTTPProxy(link)
	} else if strings.HasPrefix(link, "cottendns://") {
		return validateWhiteDNS(link, "cottendns")
	} else if strings.HasPrefix(link, "stormdns://") {
		return validateWhiteDNS(link, "stormdns")
	}

	return false, "unsupported proxy scheme"
}

// SanitizeConfigs validates a slice of configs and filters out all rejected ones while logging to telemetry
func SanitizeConfigs(configs extractor.Configs, source string, recordDropped func(source, candidate, reason string)) extractor.Configs {
	return extractor.Configs{
		Vmess:     filterSlice(configs.Vmess, source, recordDropped),
		Vless:     filterSlice(configs.Vless, source, recordDropped),
		Trojan:    filterSlice(configs.Trojan, source, recordDropped),
		SS:        filterSlice(configs.SS, source, recordDropped),
		SSR:       filterSlice(configs.SSR, source, recordDropped),
		TUIC:      filterSlice(configs.TUIC, source, recordDropped),
		Hy2:       filterSlice(configs.Hy2, source, recordDropped),
		Hysteria:  filterSlice(configs.Hysteria, source, recordDropped),
		Socks:     filterSlice(configs.Socks, source, recordDropped),
		WireGuard: filterSlice(configs.WireGuard, source, recordDropped),
		Juicity:   filterSlice(configs.Juicity, source, recordDropped),
		Naive:     filterSlice(configs.Naive, source, recordDropped),
		Telegram:  filterSlice(configs.Telegram, source, recordDropped),
		AnyTLS:    filterSlice(configs.AnyTLS, source, recordDropped),
		Snell:     filterSlice(configs.Snell, source, recordDropped),
		HTTP:      filterSlice(configs.HTTP, source, recordDropped),
		CottenDNS: filterSlice(configs.CottenDNS, source, recordDropped),
		StormDNS:  filterSlice(configs.StormDNS, source, recordDropped),
	}
}

func filterSlice(list []string, source string, recordDropped func(source, candidate, reason string)) []string {
	valid := make([]string, 0, len(list))
	for _, item := range list {
		if ok, reason := ValidateConfig(item); ok {
			valid = append(valid, item)
		} else {
			if recordDropped != nil {
				recordDropped(source, item, reason)
			}
		}
	}
	return valid
}

func isBannerOrTimestamp(link string) bool {
	lower := strings.ToLower(link)
	if strings.Contains(lower, "127.0.0.1:0") ||
		strings.Contains(lower, "0.0.0.0:0") ||
		strings.Contains(lower, "last%20update") ||
		strings.Contains(lower, "last update") ||
		strings.Contains(lower, "update:") ||
		strings.Contains(lower, "channel:") ||
		strings.Contains(lower, "کانال") {
		// If port is 0 or host is loopback with timestamp keywords
		if strings.Contains(link, ":0") || strings.Contains(link, "127.0.0.1") || strings.Contains(link, "0.0.0.0") {
			return true
		}
	}
	return false
}

func validateHostAndPort(host, portStr string, defaultPort int) (bool, string) {
	host = strings.TrimSpace(strings.Trim(host, "[]"))
	if host == "" {
		return false, "missing server hostname or ip"
	}

	// Reject loopback and local multicast
	if host == "127.0.0.1" || host == "0.0.0.0" || host == "localhost" || host == "::1" || host == "0000:0000:0000:0000:0000:0000:0000:0001" {
		return false, fmt.Sprintf("invalid loopback server host: %s", host)
	}

	// Check if IP is in private RFC1918 range
	if ip := net.ParseIP(host); ip != nil {
		if ip.IsLoopback() || ip.IsUnspecified() || ip.IsLinkLocalUnicast() || ip.IsPrivate() {
			return false, fmt.Sprintf("invalid private or non-routable ip: %s", host)
		}
	}

	port := defaultPort
	if portStr != "" {
		p, err := strconv.Atoi(portStr)
		if err != nil || p < 1 || p > 65535 {
			return false, fmt.Sprintf("invalid port: %s", portStr)
		}
		port = p
	}

	if port < 1 || port > 65535 {
		return false, fmt.Sprintf("invalid port number: %d", port)
	}

	return true, ""
}

func validateUUID(id string) (bool, string) {
	id = strings.TrimSpace(id)
	if id == "" {
		return false, "missing uuid"
	}
	if id == "00000000-0000-0000-0000-000000000000" || id == "00000000000000000000000000000000" {
		return false, "dummy placeholder uuid (all zeros)"
	}
	lower := strings.ToLower(id)
	if strings.Contains(lower, "your_uuid") || strings.Contains(lower, "your-uuid") || lower == "test" || lower == "uuid" {
		return false, fmt.Sprintf("placeholder uuid: %s", id)
	}
	if !uuidRegex.MatchString(id) && !hexRegex.MatchString(id) {
		return false, fmt.Sprintf("malformed uuid format: %s", id)
	}
	return true, ""
}

func validateVLESS(link string) (bool, string) {
	raw := strings.Split(link, "#")[0]
	u, err := url.Parse(raw)
	if err != nil {
		return false, "malformed vless url"
	}

	if u.User == nil {
		return false, "missing vless uuid"
	}
	if ok, reason := validateUUID(u.User.Username()); !ok {
		return false, reason
	}

	return validateHostAndPort(u.Hostname(), u.Port(), 443)
}

func validateVMess(link string) (bool, string) {
	payload := strings.TrimPrefix(link, "vmess://")
	payload = strings.TrimSpace(strings.Split(payload, "#")[0])

	decoded, err := decodeBase64Safe(payload)
	if err != nil {
		return false, "invalid vmess base64 payload"
	}

	var v map[string]interface{}
	if err := json.Unmarshal(decoded, &v); err != nil {
		return false, "invalid vmess json"
	}

	id := fmt.Sprintf("%v", v["id"])
	if ok, reason := validateUUID(id); !ok {
		return false, reason
	}

	add := fmt.Sprintf("%v", v["add"])
	portStr := fmt.Sprintf("%v", v["port"])
	return validateHostAndPort(add, portStr, 443)
}

func validateTrojan(link string) (bool, string) {
	raw := strings.Split(link, "#")[0]
	u, err := url.Parse(raw)
	if err != nil {
		return false, "malformed trojan url"
	}

	if u.User == nil || strings.TrimSpace(u.User.Username()) == "" {
		return false, "missing trojan password"
	}

	return validateHostAndPort(u.Hostname(), u.Port(), 443)
}

func validateSS(link string) (bool, string) {
	raw := strings.TrimPrefix(link, "ss://")
	raw = strings.Split(raw, "#")[0]

	var userInfo, hostPort string
	if strings.Contains(raw, "@") {
		atSplit := strings.SplitN(raw, "@", 2)
		userInfo = atSplit[0]
		if unescaped, err := url.QueryUnescape(userInfo); err == nil && unescaped != "" {
			userInfo = unescaped
		}
		rest := atSplit[1]
		if qIdx := strings.Index(rest, "?"); qIdx != -1 {
			hostPort = rest[:qIdx]
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
		// Legacy base64
		if decoded, err := decodeBase64Safe(raw); err == nil {
			decodedStr := string(decoded)
			if atIdx := strings.LastIndex(decodedStr, "@"); atIdx != -1 {
				userInfo = decodedStr[:atIdx]
				hostPort = decodedStr[atIdx+1:]
			}
		}
	}

	parts := strings.SplitN(userInfo, ":", 2)
	if len(parts) < 2 || strings.TrimSpace(parts[1]) == "" {
		return false, "missing shadowsocks cipher or password"
	}

	hostPort = strings.TrimRight(hostPort, "/")
	h, p, err := net.SplitHostPort(hostPort)
	if err != nil {
		h = hostPort
		p = "8388"
	}

	return validateHostAndPort(h, p, 8388)
}

func validateSSR(link string) (bool, string) {
	payload := strings.TrimPrefix(link, "ssr://")
	payload = strings.TrimSpace(strings.Split(payload, "#")[0])

	decoded, err := decodeBase64Safe(payload)
	if err != nil {
		return false, "invalid ssr base64 payload"
	}

	parts := strings.SplitN(string(decoded), "/?", 2)
	mainParts := strings.Split(parts[0], ":")
	if len(mainParts) < 6 {
		return false, "malformed ssr payload"
	}

	server := mainParts[0]
	port := mainParts[1]
	return validateHostAndPort(server, port, 8388)
}

func validateTUIC(link string) (bool, string) {
	raw := strings.Split(link, "#")[0]
	u, err := url.Parse(raw)
	if err != nil {
		return false, "malformed tuic url"
	}

	if u.User == nil || strings.TrimSpace(u.User.Username()) == "" {
		return false, "missing tuic credentials"
	}

	return validateHostAndPort(u.Hostname(), u.Port(), 443)
}

func validateHy2(link string) (bool, string) {
	raw := strings.Split(link, "#")[0]
	u, err := url.Parse(raw)
	if err != nil {
		return false, "malformed hysteria url"
	}

	if u.User == nil || strings.TrimSpace(u.User.Username()) == "" {
		return false, "missing hysteria password"
	}

	return validateHostAndPort(u.Hostname(), u.Port(), 443)
}

func validateSocks(link string) (bool, string) {
	raw := strings.Split(link, "#")[0]
	u, err := url.Parse(raw)
	if err != nil {
		return false, "malformed socks url"
	}

	return validateHostAndPort(u.Hostname(), u.Port(), 1080)
}

func validateWireGuard(link string) (bool, string) {
	raw := strings.Split(link, "#")[0]
	u, err := url.Parse(raw)
	if err != nil {
		return false, "malformed wireguard url"
	}

	privKey := ""
	if u.User != nil {
		privKey = u.User.Username()
	}
	if privKey == "" {
		privKey = u.Query().Get("private_key")
		if privKey == "" {
			privKey = u.Query().Get("privkey")
		}
	}
	if privKey == "" {
		return false, "missing wireguard private key"
	}

	pubKey := u.Query().Get("public_key")
	if pubKey == "" {
		pubKey = u.Query().Get("publicKey")
		if pubKey == "" {
			pubKey = u.Query().Get("pubkey")
		}
	}
	if pubKey == "" {
		return false, "missing wireguard public key"
	}

	return validateHostAndPort(u.Hostname(), u.Port(), 51820)
}

func validateTelegram(link string) (bool, string) {
	link = strings.Split(link, "#")[0]
	u, err := url.Parse(link)
	if err != nil {
		return false, "malformed telegram proxy url"
	}
	q := u.Query()
	server := strings.TrimSuffix(strings.TrimSpace(q.Get("server")), ".")
	if server == "" {
		return false, "missing telegram proxy server"
	}
	port := strings.TrimSpace(q.Get("port"))
	secret := strings.TrimSpace(q.Get("secret"))
	if secret == "" {
		return false, "missing telegram proxy secret"
	}
	if len(secret) < 20 {
		return false, "telegram proxy secret too short"
	}

	return validateHostAndPort(server, port, 443)
}

func validateJuicity(link string) (bool, string) {
	u, err := url.Parse(link)
	if err != nil {
		return false, "malformed juicity url"
	}
	if u.User == nil || strings.TrimSpace(u.User.Username()) == "" {
		return false, "missing juicity uuid/user"
	}
	return validateHostAndPort(u.Hostname(), u.Port(), 443)
}

func validateNaive(link string) (bool, string) {
	raw := strings.TrimPrefix(link, "naive+")
	u, err := url.Parse(raw)
	if err != nil {
		return false, "malformed naiveproxy url"
	}
	if u.User == nil || strings.TrimSpace(u.User.Username()) == "" {
		return false, "missing naiveproxy credentials"
	}
	return validateHostAndPort(u.Hostname(), u.Port(), 443)
}

func validateAnyTLS(link string) (bool, string) {
	u, err := url.Parse(link)
	if err != nil {
		return false, "malformed anytls url"
	}
	if u.User == nil || strings.TrimSpace(u.User.Username()) == "" {
		return false, "missing anytls uuid"
	}
	if ok, reason := validateUUID(u.User.Username()); !ok {
		return false, reason
	}
	return validateHostAndPort(u.Hostname(), u.Port(), 443)
}

func validateSnell(link string) (bool, string) {
	u, err := url.Parse(link)
	if err != nil {
		return false, "malformed snell url"
	}
	if u.User == nil || strings.TrimSpace(u.User.Username()) == "" {
		return false, "missing snell psk"
	}
	return validateHostAndPort(u.Hostname(), u.Port(), 443)
}

func validateHTTPProxy(link string) (bool, string) {
	u, err := url.Parse(link)
	if err != nil {
		return false, "malformed http proxy url"
	}
	if u.User == nil || strings.TrimSpace(u.User.Username()) == "" {
		return false, "missing http proxy credentials"
	}
	defaultPort := 8080
	if strings.HasPrefix(link, "https://") {
		defaultPort = 8443
	}
	return validateHostAndPort(u.Hostname(), u.Port(), defaultPort)
}

func validateWhiteDNS(link, scheme string) (bool, string) {
	payload := strings.TrimPrefix(link, scheme+"://")
	decoded, err := decodeBase64Safe(payload)
	if err != nil {
		return false, fmt.Sprintf("invalid %s base64 payload", scheme)
	}

	var data map[string]interface{}
	if err := json.Unmarshal(decoded, &data); err != nil {
		return false, fmt.Sprintf("invalid %s json payload", scheme)
	}

	profile, ok := data["profile"].(map[string]interface{})
	if !ok {
		return false, fmt.Sprintf("missing %s profile", scheme)
	}
	server, ok := profile["server"].(map[string]interface{})
	if !ok {
		return false, fmt.Sprintf("missing %s server config", scheme)
	}

	domainStr := strings.TrimSpace(safeString(server["domain"]))
	if domainStr == "" {
		return false, fmt.Sprintf("missing %s domain", scheme)
	}

	firstDomain := strings.TrimSpace(strings.Split(domainStr, ",")[0])
	firstDomain = strings.Trim(firstDomain, "[]")

	return validateHostAndPort(firstDomain, "443", 443)
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

func safeString(v interface{}) string {
	if v == nil {
		return ""
	}
	return fmt.Sprintf("%v", v)
}
