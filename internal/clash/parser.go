package clash

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

// ParseConfigToClash parses any supported proxy URL into a Clash Meta proxy map
func ParseConfigToClash(link string) (map[string]interface{}, string, error) {
	link = strings.TrimSpace(link)
	if strings.HasPrefix(link, "vless://") {
		return parseVLESS(link)
	} else if strings.HasPrefix(link, "vmess://") {
		return parseVMess(link)
	} else if strings.HasPrefix(link, "trojan://") {
		return parseTrojan(link)
	} else if strings.HasPrefix(link, "ss://") {
		return parseSS(link)
	} else if strings.HasPrefix(link, "ssr://") {
		return parseSSR(link)
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
	return nil, "", fmt.Errorf("unsupported protocol scheme")
}

// extractNameAndCountry extracts the remark and ISO country code from URL fragment or remark
func extractNameAndCountry(remark string) (name string, country string) {
	remark = strings.TrimSpace(remark)
	if unescaped, err := url.QueryUnescape(remark); err == nil && unescaped != "" {
		remark = unescaped
	}
	if remark == "" {
		remark = "Proxy"
	}

	// Example format: "🇩🇪 DE - [t.me/channel]"
	parts := strings.Fields(remark)
	if len(parts) >= 2 && len(parts[1]) == 2 && strings.ToUpper(parts[1]) == parts[1] {
		country = parts[1]
	}

	return remark, country
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

	name, country := extractNameAndCountry(u.Fragment)

	proxy := map[string]interface{}{
		"name":   name,
		"type":   "vless",
		"server": server,
		"port":   port,
		"uuid":   uuid,
		"udp":    true,
	}

	if security == "tls" || security == "reality" {
		proxy["tls"] = true
		if sni != "" {
			proxy["servername"] = sni
		}
		if fp := q.Get("fp"); fp != "" {
			proxy["client-fingerprint"] = fp
		}
		if alpn := q.Get("alpn"); alpn != "" {
			proxy["alpn"] = splitAndTrim(alpn, ",")
		}
		if q.Get("allowInsecure") == "1" || q.Get("insecure") == "1" || strings.EqualFold(q.Get("allowInsecure"), "true") {
			proxy["skip-cert-verify"] = true
		}
	}

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

	spx := q.Get("spx")
	if spx == "" {
		spx = q.Get("spider_x")
	}
	if spx == "" {
		spx = q.Get("spiderX")
	}

	if security == "reality" || pbk != "" {
		proxy["tls"] = true
		realityOpts := map[string]interface{}{
			"public-key": pbk,
		}
		if sid != "" {
			realityOpts["short-id"] = sid
		}
		if spx != "" {
			realityOpts["spider-x"] = spx
		}
		proxy["reality-opts"] = realityOpts
	}

	if flow := q.Get("flow"); flow != "" {
		proxy["flow"] = flow
	}

	proxy["network"] = network
	switch network {
	case "ws":
		wsOpts := map[string]interface{}{}
		if path := q.Get("path"); path != "" {
			wsOpts["path"] = path
		} else {
			wsOpts["path"] = "/"
		}
		if host := q.Get("host"); host != "" {
			wsOpts["headers"] = map[string]interface{}{
				"Host": host,
			}
		}
		proxy["ws-opts"] = wsOpts
	case "grpc":
		grpcOpts := map[string]interface{}{}
		serviceName := q.Get("serviceName")
		if serviceName == "" {
			serviceName = q.Get("path")
		}
		if serviceName != "" {
			grpcOpts["grpc-service-name"] = serviceName
		}
		proxy["grpc-opts"] = grpcOpts
	case "h2", "http":
		h2Opts := map[string]interface{}{}
		if path := q.Get("path"); path != "" {
			h2Opts["path"] = splitAndTrim(path, ",")
		}
		if host := q.Get("host"); host != "" {
			h2Opts["host"] = splitAndTrim(host, ",")
		}
		proxy["h2-opts"] = h2Opts
	}

	return proxy, country, nil
}

func parseVMess(link string) (map[string]interface{}, string, error) {
	payload := strings.TrimPrefix(link, "vmess://")
	payload = strings.TrimSpace(payload)
	payload = strings.ReplaceAll(payload, "-", "+")
	payload = strings.ReplaceAll(payload, "_", "/")
	if pad := len(payload) % 4; pad != 0 {
		payload += strings.Repeat("=", 4-pad)
	}

	decoded, err := base64.StdEncoding.DecodeString(payload)
	if err != nil {
		return nil, "", fmt.Errorf("vmess base64 decode error: %w", err)
	}

	var v map[string]interface{}
	if err := json.Unmarshal(decoded, &v); err != nil {
		return nil, "", fmt.Errorf("vmess json decode error: %w", err)
	}

	server := toString(v["add"])
	port := toInt(v["port"], 443)
	uuid := toString(v["id"])
	aid := toInt(v["aid"], 0)
	scy := toString(v["scy"])
	if scy == "" {
		scy = "auto"
	}

	if server == "" || uuid == "" {
		return nil, "", fmt.Errorf("missing server or uuid in vmess json")
	}

	rawPs := toString(v["ps"])
	name, country := extractNameAndCountry(rawPs)

	proxy := map[string]interface{}{
		"name":     name,
		"type":     "vmess",
		"server":   server,
		"port":     port,
		"uuid":     uuid,
		"alterId":  aid,
		"cipher":   scy,
		"udp":      true,
	}

	tlsVal := strings.ToLower(toString(v["tls"]))
	isTLS := tlsVal == "tls" || tlsVal == "1" || tlsVal == "true"
	if isTLS {
		proxy["tls"] = true
		sni := toString(v["sni"])
		if sni == "" {
			sni = toString(v["host"])
		}
		if sni == "" {
			sni = server
		}
		proxy["servername"] = sni

		if fp := toString(v["fp"]); fp != "" {
			proxy["client-fingerprint"] = fp
		}
		if alpn := toString(v["alpn"]); alpn != "" {
			proxy["alpn"] = splitAndTrim(alpn, ",")
		}
		if toString(v["allowInsecure"]) == "1" || toString(v["insecure"]) == "1" {
			proxy["skip-cert-verify"] = true
		}
	}

	net := strings.ToLower(toString(v["net"]))
	if net == "" {
		net = "tcp"
	}
	proxy["network"] = net

	switch net {
	case "ws":
		wsOpts := map[string]interface{}{}
		path := toString(v["path"])
		if path == "" {
			path = "/"
		}
		wsOpts["path"] = path
		if host := toString(v["host"]); host != "" {
			wsOpts["headers"] = map[string]interface{}{
				"Host": host,
			}
		}
		proxy["ws-opts"] = wsOpts
	case "grpc":
		grpcOpts := map[string]interface{}{}
		serviceName := toString(v["path"])
		if serviceName == "" {
			serviceName = toString(v["serviceName"])
		}
		if serviceName != "" {
			grpcOpts["grpc-service-name"] = serviceName
		}
		proxy["grpc-opts"] = grpcOpts
	case "h2", "http":
		h2Opts := map[string]interface{}{}
		if path := toString(v["path"]); path != "" {
			h2Opts["path"] = splitAndTrim(path, ",")
		}
		if host := toString(v["host"]); host != "" {
			h2Opts["host"] = splitAndTrim(host, ",")
		}
		proxy["h2-opts"] = h2Opts
	}

	return proxy, country, nil
}

func parseTrojan(link string) (map[string]interface{}, string, error) {
	u, err := url.Parse(link)
	if err != nil {
		return nil, "", err
	}

	server := u.Hostname()
	port := toInt(u.Port(), 443)
	password := ""
	if u.User != nil {
		password = u.User.Username()
	}
	if unescaped, err := url.QueryUnescape(password); err == nil {
		password = unescaped
	}
	if port <= 0 || port > 65535 || server == "" || password == "" {
		return nil, "", fmt.Errorf("missing or invalid server, port, or password in trojan link")
	}

	q := u.Query()
	name, country := extractNameAndCountry(u.Fragment)

	sni := q.Get("sni")
	if sni == "" {
		sni = q.Get("peer")
	}
	if sni == "" {
		sni = server
	}

	proxy := map[string]interface{}{
		"name":     name,
		"type":     "trojan",
		"server":   server,
		"port":     port,
		"password": password,
		"udp":      true,
		"sni":      sni,
	}

	if fp := q.Get("fp"); fp != "" {
		proxy["client-fingerprint"] = fp
	}
	if alpn := q.Get("alpn"); alpn != "" {
		proxy["alpn"] = splitAndTrim(alpn, ",")
	}
	if q.Get("allowInsecure") == "1" || q.Get("insecure") == "1" || strings.EqualFold(q.Get("allowInsecure"), "true") {
		proxy["skip-cert-verify"] = true
	}

	network := q.Get("type")
	if network == "" {
		network = q.Get("net")
	}
	if network == "" {
		network = "tcp"
	}
	proxy["network"] = network

	switch network {
	case "ws":
		wsOpts := map[string]interface{}{}
		if path := q.Get("path"); path != "" {
			wsOpts["path"] = path
		} else {
			wsOpts["path"] = "/"
		}
		if host := q.Get("host"); host != "" {
			wsOpts["headers"] = map[string]interface{}{
				"Host": host,
			}
		}
		proxy["ws-opts"] = wsOpts
	case "grpc":
		grpcOpts := map[string]interface{}{}
		serviceName := q.Get("serviceName")
		if serviceName == "" {
			serviceName = q.Get("path")
		}
		if serviceName != "" {
			grpcOpts["grpc-service-name"] = serviceName
		}
		proxy["grpc-opts"] = grpcOpts
	}

	return proxy, country, nil
}

var validSSCiphers = map[string]bool{
	"aes-128-gcm":                   true,
	"aes-192-gcm":                   true,
	"aes-256-gcm":                   true,
	"aes-128-cfb":                   true,
	"aes-192-cfb":                   true,
	"aes-256-cfb":                   true,
	"aes-128-ctr":                   true,
	"aes-192-ctr":                   true,
	"aes-256-ctr":                   true,
	"rc4-md5":                       true,
	"chacha20-ietf":                 true,
	"chacha20":                      true,
	"chacha20-ietf-poly1305":        true,
	"xchacha20-ietf-poly1305":       true,
	"2022-blake3-aes-128-gcm":       true,
	"2022-blake3-aes-256-gcm":       true,
	"2022-blake3-chacha20-poly1305": true,
	"2022-blake3-chacha8-poly1305":  true,
	"none":                          true,
	"dummy":                         true,
}

func isValidSSCipher(cipher string) bool {
	return validSSCiphers[strings.ToLower(strings.TrimSpace(cipher))]
}

func parseSS(link string) (map[string]interface{}, string, error) {
	raw := strings.TrimPrefix(link, "ss://")
	parts := strings.SplitN(raw, "#", 2)
	mainPart := parts[0]
	remark := ""
	if len(parts) > 1 {
		remark = parts[1]
	}

	name, country := extractNameAndCountry(remark)

	var userInfo, hostPort, query string
	if strings.Contains(mainPart, "@") {
		// Standard SIP002: base64(method:password)@host:port/?query
		atSplit := strings.SplitN(mainPart, "@", 2)
		userInfo = atSplit[0]
		rest := atSplit[1]

		if qIdx := strings.Index(rest, "?"); qIdx != -1 {
			hostPort = rest[:qIdx]
			query = rest[qIdx+1:]
		} else {
			hostPort = strings.TrimRight(rest, "/")
		}

		// Decode userInfo if Base64
		if decoded, err := decodeBase64Safe(userInfo); err == nil {
			decodedStr := string(decoded)
			if strings.Contains(decodedStr, ":") && isPrintable(decodedStr) {
				p := strings.SplitN(decodedStr, ":", 2)
				if isValidSSCipher(p[0]) {
					userInfo = decodedStr
				}
			}
		}
	} else {
		// Legacy format: base64(method:password@host:port)
		decoded, err := decodeBase64Safe(mainPart)
		if err != nil {
			return nil, "", fmt.Errorf("ss legacy base64 decode error: %w", err)
		}
		decodedStr := string(decoded)
		if strings.Contains(decodedStr, "@") && isPrintable(decodedStr) {
			atSplit := strings.SplitN(decodedStr, "@", 2)
			userInfo = atSplit[0]
			hostPort = atSplit[1]
		} else {
			return nil, "", fmt.Errorf("invalid ss format")
		}
	}

	// Parse method & password
	userParts := strings.SplitN(userInfo, ":", 2)
	if len(userParts) < 2 {
		return nil, "", fmt.Errorf("invalid ss credentials")
	}
	cipher := userParts[0]
	password := userParts[1]

	if unescaped, err := url.QueryUnescape(cipher); err == nil {
		cipher = unescaped
	}
	if unescaped, err := url.QueryUnescape(password); err == nil {
		password = unescaped
	}

	cipher = strings.ToLower(strings.TrimSpace(cipher))
	if !isValidSSCipher(cipher) {
		return nil, "", fmt.Errorf("unsupported or invalid ss cipher: %s", cipher)
	}

	// Validate and normalize Shadowsocks 2022 keys
	if strings.HasPrefix(cipher, "2022-") {
		keys := strings.Split(password, ":")
		var normalizedKeys []string
		for _, k := range keys {
			dec, err := decodeBase64Safe(k)
			if err != nil || len(dec) == 0 {
				return nil, "", fmt.Errorf("invalid base64 key for shadowsocks 2022 (%s): %w", cipher, err)
			}
			normalizedKeys = append(normalizedKeys, base64.StdEncoding.EncodeToString(dec))
		}
		password = strings.Join(normalizedKeys, ":")
	}

	// Parse host & port
	hpParts := strings.Split(hostPort, ":")
	if len(hpParts) < 2 {
		return nil, "", fmt.Errorf("invalid ss host:port")
	}
	server := hpParts[0]
	port := toInt(hpParts[1], 8388)
	if port <= 0 || port > 65535 || server == "" {
		return nil, "", fmt.Errorf("invalid server or port in ss link")
	}

	proxy := map[string]interface{}{
		"name":     name,
		"type":     "ss",
		"server":   server,
		"port":     port,
		"cipher":   cipher,
		"password": password,
		"udp":      true,
	}

	// Parse plugin options if present
	if query != "" {
		if q, err := url.ParseQuery(query); err == nil {
			pluginStr := q.Get("plugin")
			if pluginStr != "" {
				parseSSPlugin(proxy, pluginStr)
			}
		}
	}

	return proxy, country, nil
}

func parseSSPlugin(proxy map[string]interface{}, pluginStr string) {
	parts := strings.Split(pluginStr, ";")
	pluginName := parts[0]
	opts := make(map[string]interface{})

	for _, p := range parts[1:] {
		kv := strings.SplitN(p, "=", 2)
		if len(kv) == 2 {
			k, v := kv[0], kv[1]
			switch k {
			case "obfs":
				opts["mode"] = v
			case "obfs-host", "host":
				opts["host"] = v
			case "tls":
				if v == "1" || v == "true" {
					opts["tls"] = true
				}
			default:
				opts[k] = v
			}
		}
	}

	if pluginName == "obfs-local" || pluginName == "simple-obfs" {
		proxy["plugin"] = "obfs"
		proxy["plugin-opts"] = opts
	} else if pluginName == "v2ray-plugin" {
		proxy["plugin"] = "v2ray-plugin"
		proxy["plugin-opts"] = opts
	}
}

func parseSSR(link string) (map[string]interface{}, string, error) {
	payload := strings.TrimPrefix(link, "ssr://")
	decoded, err := decodeBase64Safe(payload)
	if err != nil {
		return nil, "", err
	}

	str := string(decoded)
	mainAndQuery := strings.SplitN(str, "/?", 2)
	mainParts := strings.Split(mainAndQuery[0], ":")
	if len(mainParts) < 6 {
		return nil, "", fmt.Errorf("invalid ssr format")
	}

	server := mainParts[0]
	port := toInt(mainParts[1], 8388)
	protocol := mainParts[2]
	method := mainParts[3]
	obfs := mainParts[4]
	passDecoded, _ := decodeBase64Safe(mainParts[5])
	password := string(passDecoded)

	var remark, obfsParam, protoParam string
	if len(mainAndQuery) > 1 {
		q, _ := url.ParseQuery(mainAndQuery[1])
		if r := q.Get("remarks"); r != "" {
			if d, err := decodeBase64Safe(r); err == nil {
				remark = string(d)
			}
		}
		if op := q.Get("obfsparam"); op != "" {
			if d, err := decodeBase64Safe(op); err == nil {
				obfsParam = string(d)
			}
		}
		if pp := q.Get("protoparam"); pp != "" {
			if d, err := decodeBase64Safe(pp); err == nil {
				protoParam = string(d)
			}
		}
	}

	name, country := extractNameAndCountry(remark)

	proxy := map[string]interface{}{
		"name":           name,
		"type":           "ssr",
		"server":         server,
		"port":           port,
		"cipher":         method,
		"password":       password,
		"protocol":       protocol,
		"obfs":           obfs,
		"protocol-param": protoParam,
		"obfs-param":     obfsParam,
		"udp":            true,
	}

	return proxy, country, nil
}

func parseHy2(link string) (map[string]interface{}, string, error) {
	u, err := url.Parse(link)
	if err != nil {
		return nil, "", err
	}

	server := u.Hostname()
	port := toInt(u.Port(), 443)
	password := ""
	if u.User != nil {
		password = u.User.Username()
	}
	if unescaped, err := url.QueryUnescape(password); err == nil {
		password = unescaped
	}

	if port <= 0 || port > 65535 || server == "" || password == "" {
		return nil, "", fmt.Errorf("missing or invalid server, port, or password in hy2 link")
	}

	q := u.Query()
	name, country := extractNameAndCountry(u.Fragment)

	sni := q.Get("sni")
	if sni == "" {
		sni = server
	}

	proxy := map[string]interface{}{
		"name":     name,
		"type":     "hysteria2",
		"server":   server,
		"port":     port,
		"password": password,
		"sni":      sni,
	}

	if q.Get("insecure") == "1" || q.Get("allow_insecure") == "1" || strings.EqualFold(q.Get("insecure"), "true") {
		proxy["skip-cert-verify"] = true
	}

	if obfs := q.Get("obfs"); obfs != "" {
		proxy["obfs"] = obfs
		if obfsPass := q.Get("obfs-password"); obfsPass != "" {
			proxy["obfs-password"] = obfsPass
		}
	}

	return proxy, country, nil
}

func parseHysteria(link string) (map[string]interface{}, string, error) {
	u, err := url.Parse(link)
	if err != nil {
		return nil, "", err
	}

	server := u.Hostname()
	port := toInt(u.Port(), 443)
	if port <= 0 || port > 65535 || server == "" {
		return nil, "", fmt.Errorf("missing or invalid server/port in hysteria link")
	}

	q := u.Query()

	auth := q.Get("auth")
	if auth == "" && u.User != nil {
		auth = u.User.Username()
	}
	if unescaped, err := url.QueryUnescape(auth); err == nil {
		auth = unescaped
	}

	sni := q.Get("peer")
	if sni == "" {
		sni = q.Get("sni")
	}
	if sni == "" {
		sni = server
	}

	name, country := extractNameAndCountry(u.Fragment)

	proxy := map[string]interface{}{
		"name":     name,
		"type":     "hysteria",
		"server":   server,
		"port":     port,
		"auth_str": auth,
		"sni":      sni,
	}

	if up := q.Get("upmbps"); up != "" {
		proxy["up"] = up
	}
	if down := q.Get("downmbps"); down != "" {
		proxy["down"] = down
	}
	if alpn := q.Get("alpn"); alpn != "" {
		proxy["alpn"] = splitAndTrim(alpn, ",")
	}
	if q.Get("insecure") == "1" || strings.EqualFold(q.Get("insecure"), "true") {
		proxy["skip-cert-verify"] = true
	}
	if obfs := q.Get("obfs"); obfs != "" {
		proxy["obfs"] = obfs
	}
	if protocol := q.Get("protocol"); protocol != "" {
		proxy["protocol"] = protocol
	}

	return proxy, country, nil
}

func parseTUIC(link string) (map[string]interface{}, string, error) {
	u, err := url.Parse(link)
	if err != nil {
		return nil, "", err
	}

	server := u.Hostname()
	port := toInt(u.Port(), 443)

	uuid := ""
	password := ""
	if u.User != nil {
		uuid = u.User.Username()
		password, _ = u.User.Password()
	}
	if unescaped, err := url.QueryUnescape(uuid); err == nil {
		uuid = unescaped
	}
	if unescaped, err := url.QueryUnescape(password); err == nil {
		password = unescaped
	}

	if port <= 0 || port > 65535 || server == "" || uuid == "" {
		return nil, "", fmt.Errorf("missing or invalid server, port, or uuid in tuic link")
	}

	q := u.Query()
	name, country := extractNameAndCountry(u.Fragment)

	sni := q.Get("sni")
	if sni == "" {
		sni = server
	}

	proxy := map[string]interface{}{
		"name":     name,
		"type":     "tuic",
		"server":   server,
		"port":     port,
		"uuid":     uuid,
		"password": password,
		"sni":      sni,
	}

	if cc := q.Get("congestion_control"); cc != "" {
		proxy["congestion-controller"] = cc
	}
	if udpMode := q.Get("udp_relay_mode"); udpMode != "" {
		proxy["udp-relay-mode"] = udpMode
	}
	if alpn := q.Get("alpn"); alpn != "" {
		proxy["alpn"] = splitAndTrim(alpn, ",")
	}
	if q.Get("insecure") == "1" || q.Get("allow_insecure") == "1" || strings.EqualFold(q.Get("insecure"), "true") {
		proxy["skip-cert-verify"] = true
	}

	return proxy, country, nil
}

func parseSocks(link string) (map[string]interface{}, string, error) {
	u, err := url.Parse(link)
	if err != nil {
		return nil, "", err
	}

	server := u.Hostname()
	port := toInt(u.Port(), 1080)
	if port <= 0 || port > 65535 || server == "" {
		return nil, "", fmt.Errorf("missing or invalid server/port in socks link")
	}

	name, country := extractNameAndCountry(u.Fragment)

	proxy := map[string]interface{}{
		"name":   name,
		"type":   "socks5",
		"server": server,
		"port":   port,
		"udp":    true,
	}

	if u.User != nil {
		user := u.User.Username()
		if unescaped, err := url.QueryUnescape(user); err == nil {
			user = unescaped
		}
		proxy["username"] = user
		if pass, ok := u.User.Password(); ok {
			if unescaped, err := url.QueryUnescape(pass); err == nil {
				pass = unescaped
			}
			proxy["password"] = pass
		}
	}

	return proxy, country, nil
}

func parseWireGuard(link string) (map[string]interface{}, string, error) {
	u, err := url.Parse(link)
	if err != nil {
		return nil, "", err
	}

	server := u.Hostname()
	port := toInt(u.Port(), 51820)
	q := u.Query()

	privKey := ""
	if u.User != nil {
		privKey = u.User.Username()
	}
	if privKey == "" {
		privKey = q.Get("private_key")
	}
	if privKey == "" {
		privKey = q.Get("privatekey")
	}
	if privKey == "" {
		privKey = q.Get("privkey")
	}
	if privKey == "" {
		privKey = q.Get("secret_key")
	}
	if unescaped, err := url.PathUnescape(privKey); err == nil && unescaped != "" {
		privKey = unescaped
	}
	if unescaped, err := url.QueryUnescape(privKey); err == nil && unescaped != "" {
		privKey = unescaped
	}

	if port <= 0 || port > 65535 || server == "" || privKey == "" {
		return nil, "", fmt.Errorf("missing or invalid server, port, or private-key in wireguard link")
	}

	name, country := extractNameAndCountry(u.Fragment)

	pubKey := q.Get("publickey")
	if pubKey == "" {
		pubKey = q.Get("public_key")
	}
	if pubKey == "" {
		pubKey = q.Get("pbk")
	}

	ip := q.Get("address")
	if ip == "" {
		ip = q.Get("ip")
	}
	if ip == "" {
		ip = "172.16.0.2/32"
	}
	if strings.Contains(ip, "/") {
		ip = strings.Split(ip, "/")[0]
	}

	if pubKey == "" || privKey == "" {
		return nil, "", fmt.Errorf("missing public-key or private-key in wireguard link")
	}

	proxy := map[string]interface{}{
		"name":        name,
		"type":        "wireguard",
		"server":      server,
		"port":        port,
		"private-key": privKey,
		"public-key":  pubKey,
		"ip":          ip,
		"udp":         true,
	}

	if psk := q.Get("presharedkey"); psk != "" {
		proxy["preshared-key"] = psk
	} else if psk := q.Get("preshared_key"); psk != "" {
		proxy["preshared-key"] = psk
	}

	if resStr := q.Get("reserved"); resStr != "" {
		var resInts []int
		for _, part := range strings.Split(resStr, ",") {
			if num, err := strconv.Atoi(strings.TrimSpace(part)); err == nil {
				resInts = append(resInts, num)
			}
		}
		if len(resInts) > 0 {
			proxy["reserved"] = resInts
		}
	}

	if mtu := q.Get("mtu"); mtu != "" {
		proxy["mtu"] = toInt(mtu, 1420)
	}

	return proxy, country, nil
}

// Helpers

func decodeBase64Safe(s string) ([]byte, error) {
	s = strings.TrimSpace(s)
	s = strings.ReplaceAll(s, "-", "+")
	s = strings.ReplaceAll(s, "_", "/")
	if pad := len(s) % 4; pad != 0 {
		s += strings.Repeat("=", 4-pad)
	}
	return base64.StdEncoding.DecodeString(s)
}

func splitAndTrim(s string, sep string) []string {
	parts := strings.Split(s, sep)
	var res []string
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			res = append(res, p)
		}
	}
	return res
}

func toString(v interface{}) string {
	if v == nil {
		return ""
	}
	return fmt.Sprintf("%v", v)
}

func toInt(v interface{}, defaultVal int) int {
	if v == nil {
		return defaultVal
	}
	switch val := v.(type) {
	case int:
		return val
	case float64:
		return int(val)
	case string:
		if i, err := strconv.Atoi(val); err == nil {
			return i
		}
	}
	return defaultVal
}

func isPrintable(s string) bool {
	for _, r := range s {
		if r < 32 && r != '\t' && r != '\n' && r != '\r' {
			return false
		}
	}
	return true
}
