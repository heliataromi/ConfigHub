package singbox

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net"
	"net/url"
	"strconv"
	"strings"
)

// ParseConfigToSingbox parses any supported proxy URL into a Sing-box outbound data map
func ParseConfigToSingbox(link string) (map[string]interface{}, string, error) {
	link = strings.TrimSpace(link)
	if strings.HasPrefix(link, "vless://") {
		return parseVLESS(link)
	} else if strings.HasPrefix(link, "vmess://") {
		return parseVMess(link)
	} else if strings.HasPrefix(link, "trojan://") {
		return parseTrojan(link)
	} else if strings.HasPrefix(link, "ss://") {
		return parseSS(link)
	} else if strings.HasPrefix(link, "hy2://") || strings.HasPrefix(link, "hysteria2://") {
		return parseHy2(link)
	} else if strings.HasPrefix(link, "hysteria://") {
		return parseHysteria(link)
	} else if strings.HasPrefix(link, "tuic://") {
		return parseTUIC(link)
	} else if strings.HasPrefix(link, "socks://") || strings.HasPrefix(link, "socks5://") {
		return parseSocks(link)
	} else if strings.HasPrefix(link, "wireguard://") || strings.HasPrefix(link, "wg://") {
		return parseWireGuard(link)
	}
	return nil, "", fmt.Errorf("unsupported protocol scheme for sing-box")
}

// extractTag extracts the remark/tag from URL fragment
func extractTag(remark string) string {
	remark = strings.TrimSpace(remark)
	if unescaped, err := url.QueryUnescape(remark); err == nil && unescaped != "" {
		remark = unescaped
	}
	if remark == "" {
		remark = "Proxy"
	}
	return remark
}

func parseVLESS(link string) (map[string]interface{}, string, error) {
	u, err := url.Parse(link)
	if err != nil {
		return nil, "", err
	}

	server := u.Hostname()
	portStr := u.Port()
	if portStr == "" {
		portStr = "443"
	}
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return nil, "", fmt.Errorf("invalid port: %s", portStr)
	}

	uuid := ""
	if u.User != nil {
		uuid = u.User.Username()
	}
	if server == "" || uuid == "" {
		return nil, "", fmt.Errorf("missing server or uuid in vless link")
	}

	q := u.Query()
	network := q.Get("type")
	if network == "" {
		network = q.Get("net")
	}
	if network == "" {
		network = "tcp"
	}

	security := strings.ToLower(q.Get("security"))
	sni := q.Get("sni")
	if sni == "" {
		sni = q.Get("host")
	}
	if sni == "" && (security == "tls" || security == "reality") {
		sni = server
	}

	tag := extractTag(u.Fragment)

	outbound := map[string]interface{}{
		"type":            "vless",
		"tag":             tag,
		"server":          server,
		"server_port":     port,
		"uuid":            uuid,
		"packet_encoding": "xudp",
	}

	flow := q.Get("flow")
	if flow != "" {
		outbound["flow"] = flow
	}

	// TLS & Reality configuration
	if security == "tls" || security == "reality" {
		tlsConfig := map[string]interface{}{
			"enabled": true,
		}
		if sni != "" {
			tlsConfig["server_name"] = sni
		}
		if alpn := q.Get("alpn"); alpn != "" {
			tlsConfig["alpn"] = splitAndTrim(alpn, ",")
		}
		if q.Get("allowInsecure") == "1" || q.Get("insecure") == "1" || strings.EqualFold(q.Get("allowInsecure"), "true") {
			tlsConfig["insecure"] = true
		}

		fp := q.Get("fp")
		if fp == "" {
			fp = "chrome"
		}
		tlsConfig["utls"] = map[string]interface{}{
			"enabled":     true,
			"fingerprint": fp,
		}

		if security == "reality" {
			pbk := q.Get("pbk")
			if pbk == "" {
				pbk = q.Get("public_key")
			}
			if pbk == "" {
				pbk = q.Get("publicKey")
			}

			sid := q.Get("sid")
			if sid == "" {
				sid = q.Get("short_id")
			}
			if sid == "" {
				sid = q.Get("shortId")
			}

			tlsConfig["reality"] = map[string]interface{}{
				"enabled":    true,
				"public_key": pbk,
				"short_id":   sid,
			}
		}

		outbound["tls"] = tlsConfig
	}

	// Transport configuration
	if network == "ws" || network == "websocket" {
		path := q.Get("path")
		if path == "" {
			path = "/"
		}
		host := q.Get("host")
		if host == "" {
			host = sni
		}

		transport := map[string]interface{}{
			"type": "ws",
			"path": path,
		}
		if host != "" {
			transport["headers"] = map[string]string{
				"Host": host,
			}
		}
		outbound["transport"] = transport
	} else if network == "grpc" {
		serviceName := q.Get("serviceName")
		if serviceName == "" {
			serviceName = q.Get("service_name")
		}
		transport := map[string]interface{}{
			"type": "grpc",
		}
		if serviceName != "" {
			transport["service_name"] = serviceName
		}
		outbound["transport"] = transport
	} else if network == "httpupgrade" {
		path := q.Get("path")
		if path == "" {
			path = "/"
		}
		host := q.Get("host")
		if host == "" {
			host = sni
		}
		transport := map[string]interface{}{
			"type": "httpupgrade",
			"path": path,
		}
		if host != "" {
			transport["host"] = host
		}
		outbound["transport"] = transport
	}

	return outbound, tag, nil
}

func parseVMess(link string) (map[string]interface{}, string, error) {
	raw := strings.TrimPrefix(link, "vmess://")
	decoded, err := decodeBase64Safe(raw)
	if err != nil {
		return nil, "", fmt.Errorf("failed to decode vmess base64: %w", err)
	}

	var data map[string]interface{}
	if err := json.Unmarshal(decoded, &data); err != nil {
		return nil, "", fmt.Errorf("failed to parse vmess json: %w", err)
	}

	server, _ := data["add"].(string)
	portVal := data["port"]
	port := 0
	switch v := portVal.(type) {
	case float64:
		port = int(v)
	case string:
		port, _ = strconv.Atoi(v)
	}

	uuid, _ := data["id"].(string)
	if server == "" || port == 0 || uuid == "" {
		return nil, "", fmt.Errorf("invalid vmess json fields")
	}

	tagStr, _ := data["ps"].(string)
	tag := extractTag(tagStr)

	security, _ := data["scy"].(string)
	if security == "" {
		security = "auto"
	}

	outbound := map[string]interface{}{
		"type":            "vmess",
		"tag":             tag,
		"server":          server,
		"server_port":     port,
		"uuid":            uuid,
		"security":        security,
		"alter_id":        0,
		"packet_encoding": "packetaddr",
	}

	tlsVal, _ := data["tls"].(string)
	sni, _ := data["sni"].(string)
	host, _ := data["host"].(string)
	if sni == "" {
		sni = host
	}
	if sni == "" && strings.EqualFold(tlsVal, "tls") {
		sni = server
	}

	if strings.EqualFold(tlsVal, "tls") {
		tlsConfig := map[string]interface{}{
			"enabled": true,
		}
		if sni != "" {
			tlsConfig["server_name"] = sni
		}
		if alpn, ok := data["alpn"].(string); ok && alpn != "" {
			tlsConfig["alpn"] = splitAndTrim(alpn, ",")
		}
		if fp, ok := data["fp"].(string); ok && fp != "" {
			tlsConfig["utls"] = map[string]interface{}{
				"enabled":     true,
				"fingerprint": fp,
			}
		}
		outbound["tls"] = tlsConfig
	}

	netVal, _ := data["net"].(string)
	netVal = strings.ToLower(netVal)

	if netVal == "ws" || netVal == "websocket" {
		path, _ := data["path"].(string)
		if path == "" {
			path = "/"
		}
		transport := map[string]interface{}{
			"type": "ws",
			"path": path,
		}
		if host != "" {
			transport["headers"] = map[string]string{
				"Host": host,
			}
		}
		outbound["transport"] = transport
	} else if netVal == "grpc" {
		serviceName, _ := data["path"].(string)
		transport := map[string]interface{}{
			"type": "grpc",
		}
		if serviceName != "" {
			transport["service_name"] = serviceName
		}
		outbound["transport"] = transport
	} else if netVal == "httpupgrade" {
		path, _ := data["path"].(string)
		if path == "" {
			path = "/"
		}
		transport := map[string]interface{}{
			"type": "httpupgrade",
			"path": path,
		}
		if host != "" {
			transport["host"] = host
		}
		outbound["transport"] = transport
	}

	return outbound, tag, nil
}

func parseTrojan(link string) (map[string]interface{}, string, error) {
	u, err := url.Parse(link)
	if err != nil {
		return nil, "", err
	}

	server := u.Hostname()
	portStr := u.Port()
	if portStr == "" {
		portStr = "443"
	}
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return nil, "", fmt.Errorf("invalid port: %s", portStr)
	}

	password := ""
	if u.User != nil {
		password = u.User.Username()
	}
	if server == "" || password == "" {
		return nil, "", fmt.Errorf("missing server or password in trojan link")
	}

	tag := extractTag(u.Fragment)
	q := u.Query()
	sni := q.Get("sni")
	if sni == "" {
		sni = q.Get("peer")
	}
	if sni == "" {
		sni = server
	}

	outbound := map[string]interface{}{
		"type":        "trojan",
		"tag":         tag,
		"server":      server,
		"server_port": port,
		"password":    password,
		"tls": map[string]interface{}{
			"enabled":     true,
			"server_name": sni,
		},
	}

	if alpn := q.Get("alpn"); alpn != "" {
		outbound["tls"].(map[string]interface{})["alpn"] = splitAndTrim(alpn, ",")
	}
	if q.Get("allowInsecure") == "1" || q.Get("insecure") == "1" || strings.EqualFold(q.Get("allowInsecure"), "true") {
		outbound["tls"].(map[string]interface{})["insecure"] = true
	}

	network := q.Get("type")
	if network == "" {
		network = q.Get("net")
	}
	if network == "ws" || network == "websocket" {
		path := q.Get("path")
		if path == "" {
			path = "/"
		}
		host := q.Get("host")
		if host == "" {
			host = sni
		}
		transport := map[string]interface{}{
			"type": "ws",
			"path": path,
		}
		if host != "" {
			transport["headers"] = map[string]string{
				"Host": host,
			}
		}
		outbound["transport"] = transport
	} else if network == "grpc" {
		serviceName := q.Get("serviceName")
		if serviceName == "" {
			serviceName = q.Get("service_name")
		}
		transport := map[string]interface{}{
			"type": "grpc",
		}
		if serviceName != "" {
			transport["service_name"] = serviceName
		}
		outbound["transport"] = transport
	}

	return outbound, tag, nil
}

func parseSS(link string) (map[string]interface{}, string, error) {
	raw := strings.TrimPrefix(link, "ss://")
	parts := strings.SplitN(raw, "#", 2)
	mainPart := parts[0]
	tag := "Proxy"
	if len(parts) > 1 {
		tag = extractTag(parts[1])
	}

	var method, password, server string
	var port int

	if strings.Contains(mainPart, "@") {
		// SIP002: base64(method:password)@server:port
		atSplit := strings.SplitN(mainPart, "@", 2)
		userInfo := atSplit[0]
		if decodedUser, err := decodeBase64Safe(userInfo); err == nil {
			userInfo = string(decodedUser)
		}
		userParts := strings.SplitN(userInfo, ":", 2)
		if len(userParts) == 2 {
			method = userParts[0]
			password = userParts[1]
		}

		rest := atSplit[1]
		if qIdx := strings.Index(rest, "?"); qIdx != -1 {
			rest = rest[:qIdx]
		}
		rest = strings.TrimRight(rest, "/")
		h, p, err := net.SplitHostPort(rest)
		if err != nil {
			return nil, "", err
		}
		server = strings.Trim(h, "[]")
		port, _ = strconv.Atoi(p)
	} else {
		// Legacy: base64(method:password@server:port)
		decoded, err := decodeBase64Safe(mainPart)
		if err != nil {
			return nil, "", err
		}
		decodedStr := string(decoded)
		atIdx := strings.LastIndex(decodedStr, "@")
		if atIdx == -1 {
			return nil, "", fmt.Errorf("invalid legacy ss format")
		}
		userInfo := decodedStr[:atIdx]
		userParts := strings.SplitN(userInfo, ":", 2)
		if len(userParts) == 2 {
			method = userParts[0]
			password = userParts[1]
		}
		hostPort := decodedStr[atIdx+1:]
		h, p, err := net.SplitHostPort(hostPort)
		if err != nil {
			return nil, "", err
		}
		server = strings.Trim(h, "[]")
		port, _ = strconv.Atoi(p)
	}

	if server == "" || port == 0 || method == "" || password == "" {
		return nil, "", fmt.Errorf("incomplete ss parameters")
	}

	outbound := map[string]interface{}{
		"type":        "shadowsocks",
		"tag":         tag,
		"server":      server,
		"server_port": port,
		"method":      method,
		"password":    password,
	}

	return outbound, tag, nil
}

func parseHy2(link string) (map[string]interface{}, string, error) {
	cleanLink := link
	if strings.HasPrefix(cleanLink, "hysteria2://") {
		cleanLink = "hy2://" + strings.TrimPrefix(cleanLink, "hysteria2://")
	}

	u, err := url.Parse(cleanLink)
	if err != nil {
		return nil, "", err
	}

	server := u.Hostname()
	portStr := u.Port()
	if portStr == "" {
		portStr = "443"
	}
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return nil, "", fmt.Errorf("invalid port: %s", portStr)
	}

	password := ""
	if u.User != nil {
		password = u.User.Username()
	}

	tag := extractTag(u.Fragment)
	q := u.Query()
	sni := q.Get("sni")
	if sni == "" {
		sni = server
	}

	outbound := map[string]interface{}{
		"type":        "hysteria2",
		"tag":         tag,
		"server":      server,
		"server_port": port,
		"password":    password,
		"tls": map[string]interface{}{
			"enabled":     true,
			"server_name": sni,
			"insecure":    true,
			"alpn":        []string{"h3"},
		},
	}

	if q.Get("insecure") == "0" || strings.EqualFold(q.Get("insecure"), "false") {
		outbound["tls"].(map[string]interface{})["insecure"] = false
	}

	// Obfs
	obfsType := q.Get("obfs")
	obfsPass := q.Get("obfs-password")
	if obfsType != "" && obfsPass != "" {
		outbound["obfs"] = map[string]interface{}{
			"type":     obfsType,
			"password": obfsPass,
		}
	}

	return outbound, tag, nil
}

func parseHysteria(link string) (map[string]interface{}, string, error) {
	u, err := url.Parse(link)
	if err != nil {
		return nil, "", err
	}

	server := u.Hostname()
	portStr := u.Port()
	if portStr == "" {
		portStr = "443"
	}
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return nil, "", fmt.Errorf("invalid port: %s", portStr)
	}

	q := u.Query()
	auth := q.Get("auth")
	tag := extractTag(u.Fragment)
	sni := q.Get("peer")
	if sni == "" {
		sni = server
	}

	outbound := map[string]interface{}{
		"type":        "hysteria",
		"tag":         tag,
		"server":      server,
		"server_port": port,
		"auth_str":    auth,
		"protocol":    "udp",
		"tls": map[string]interface{}{
			"enabled":     true,
			"server_name": sni,
			"insecure":    true,
			"alpn":        []string{"h3"},
		},
	}

	if obfs := q.Get("obfs"); obfs != "" {
		outbound["obfs"] = obfs
	}

	return outbound, tag, nil
}

func parseTUIC(link string) (map[string]interface{}, string, error) {
	u, err := url.Parse(link)
	if err != nil {
		return nil, "", err
	}

	server := u.Hostname()
	portStr := u.Port()
	if portStr == "" {
		portStr = "443"
	}
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return nil, "", fmt.Errorf("invalid port: %s", portStr)
	}

	uuid := ""
	password := ""
	if u.User != nil {
		uuid = u.User.Username()
		password, _ = u.User.Password()
	}

	tag := extractTag(u.Fragment)
	q := u.Query()
	sni := q.Get("sni")
	if sni == "" {
		sni = server
	}

	cc := q.Get("congestion_controller")
	if cc == "" {
		cc = "bbr"
	}

	udpRelay := q.Get("udp_relay_mode")
	if udpRelay == "" {
		udpRelay = "native"
	}

	outbound := map[string]interface{}{
		"type":                  "tuic",
		"tag":                   tag,
		"server":                server,
		"server_port":           port,
		"uuid":                  uuid,
		"password":              password,
		"congestion_controller": cc,
		"udp_relay_mode":        udpRelay,
		"zero_rtt_handshake":    false,
		"tls": map[string]interface{}{
			"enabled":     true,
			"server_name": sni,
			"insecure":    true,
			"alpn":        []string{"h3"},
		},
	}

	if q.Get("allow_insecure") == "0" || strings.EqualFold(q.Get("allow_insecure"), "false") {
		outbound["tls"].(map[string]interface{})["insecure"] = false
	}

	return outbound, tag, nil
}

func parseSocks(link string) (map[string]interface{}, string, error) {
	u, err := url.Parse(link)
	if err != nil {
		return nil, "", err
	}

	server := u.Hostname()
	portStr := u.Port()
	if portStr == "" {
		portStr = "1080"
	}
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return nil, "", fmt.Errorf("invalid port: %s", portStr)
	}

	username := ""
	password := ""
	if u.User != nil {
		username = u.User.Username()
		password, _ = u.User.Password()
	}

	tag := extractTag(u.Fragment)

	outbound := map[string]interface{}{
		"type":        "socks",
		"tag":         tag,
		"server":      server,
		"server_port": port,
	}
	if username != "" {
		outbound["username"] = username
	}
	if password != "" {
		outbound["password"] = password
	}

	return outbound, tag, nil
}

func parseWireGuard(link string) (map[string]interface{}, string, error) {
	u, err := url.Parse(link)
	if err != nil {
		return nil, "", err
	}

	server := u.Hostname()
	portStr := u.Port()
	if portStr == "" {
		portStr = "51820"
	}
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return nil, "", fmt.Errorf("invalid port: %s", portStr)
	}

	privateKey := ""
	if u.User != nil {
		privateKey = u.User.Username()
	}

	q := u.Query()
	publicKey := q.Get("publickey")
	if publicKey == "" {
		publicKey = q.Get("public_key")
	}

	ip := q.Get("ip")
	if ip == "" {
		ip = q.Get("address")
	}
	if ip == "" {
		ip = "10.0.0.2/32"
	}

	tag := extractTag(u.Fragment)

	outbound := map[string]interface{}{
		"type":            "wireguard",
		"tag":             tag,
		"server":          server,
		"server_port":     port,
		"local_address":   []string{ip},
		"private_key":     privateKey,
		"peer_public_key": publicKey,
		"mtu":             1420,
	}

	if reservedStr := q.Get("reserved"); reservedStr != "" {
		parts := strings.Split(reservedStr, ",")
		if len(parts) == 3 {
			var resInts []int
			for _, p := range parts {
				if val, err := strconv.Atoi(strings.TrimSpace(p)); err == nil {
					resInts = append(resInts, val)
				}
			}
			if len(resInts) == 3 {
				outbound["reserved"] = resInts
			}
		}
	}

	psk := q.Get("preshared_key")
	if psk == "" {
		psk = q.Get("presharedkey")
	}
	if psk == "" {
		psk = q.Get("psk")
	}
	if psk != "" {
		outbound["pre_shared_key"] = psk
	}

	return outbound, tag, nil
}

func decodeBase64Safe(s string) ([]byte, error) {
	s = strings.TrimSpace(s)
	s = strings.ReplaceAll(s, "-", "+")
	s = strings.ReplaceAll(s, "_", "/")
	switch len(s) % 4 {
	case 2:
		s += "=="
	case 3:
		s += "="
	}
	return base64.StdEncoding.DecodeString(s)
}

func splitAndTrim(s, sep string) []string {
	var result []string
	for _, part := range strings.Split(s, sep) {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}
