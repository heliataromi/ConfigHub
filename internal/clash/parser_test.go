package clash

import (
	"encoding/base64"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestParseVLESS(t *testing.T) {
	link := "vless://11111111-2222-3333-4444-555555555555@example.com:443?type=tcp&security=reality&pbk=testpubkey&sid=1234&sni=target.com&fp=chrome#%F0%9F%87%A9%F0%9F%87%AA%20DE%20-%20%5Bt.me%2Ftest%5D"
	proxy, country, err := ParseConfigToClash(link)
	if err != nil {
		t.Fatalf("Failed to parse VLESS: %v", err)
	}

	if country != "DE" {
		t.Errorf("Expected country DE, got %s", country)
	}
	if proxy["type"] != "vless" {
		t.Errorf("Expected type vless, got %v", proxy["type"])
	}
	if proxy["server"] != "example.com" {
		t.Errorf("Expected server example.com, got %v", proxy["server"])
	}
	if proxy["servername"] != "target.com" {
		t.Errorf("Expected servername target.com, got %v", proxy["servername"])
	}
	if reality, ok := proxy["reality-opts"].(map[string]interface{}); !ok || reality["public-key"] != "testpubkey" {
		t.Errorf("Expected reality-opts public-key testpubkey, got %v", proxy["reality-opts"])
	}
}

func TestParseVMess(t *testing.T) {
	vmessJSON := `{"v":"2","ps":"🇺🇸 US - [t.me/test]","add":"1.2.3.4","port":443,"id":"11111111-2222-3333-4444-555555555555","aid":0,"scy":"auto","net":"ws","type":"none","host":"host.com","path":"/ws","tls":"tls","sni":"sni.com"}`
	encoded := "vmess://" + base64.StdEncoding.EncodeToString([]byte(vmessJSON))

	proxy, country, err := ParseConfigToClash(encoded)
	if err != nil {
		t.Fatalf("Failed to parse VMess: %v", err)
	}

	if country != "US" {
		t.Errorf("Expected country US, got %s", country)
	}
	if proxy["type"] != "vmess" {
		t.Errorf("Expected type vmess, got %v", proxy["type"])
	}
	if proxy["server"] != "1.2.3.4" {
		t.Errorf("Expected server 1.2.3.4, got %v", proxy["server"])
	}
	if proxy["network"] != "ws" {
		t.Errorf("Expected network ws, got %v", proxy["network"])
	}
	if wsOpts, ok := proxy["ws-opts"].(map[string]interface{}); !ok || wsOpts["path"] != "/ws" {
		t.Errorf("Expected ws-opts path /ws, got %v", proxy["ws-opts"])
	}
}

func TestParseTrojan(t *testing.T) {
	link := "trojan://password123@trojan.example.com:443?security=tls&sni=trojan.example.com&type=ws&path=%2Ftrojan-ws#%F0%9F%87%AB%F0%9F%87%B7%20FR%20-%20%5Bt.me%2Ftest%5D"
	proxy, country, err := ParseConfigToClash(link)
	if err != nil {
		t.Fatalf("Failed to parse Trojan: %v", err)
	}

	if country != "FR" {
		t.Errorf("Expected country FR, got %s", country)
	}
	if proxy["type"] != "trojan" {
		t.Errorf("Expected type trojan, got %v", proxy["type"])
	}
	if proxy["password"] != "password123" {
		t.Errorf("Expected password password123, got %v", proxy["password"])
	}
}

func TestParseSS(t *testing.T) {
	userInfo := base64.StdEncoding.EncodeToString([]byte("aes-256-gcm:password123"))
	link := "ss://" + userInfo + "@ss.example.com:8388#%F0%9F%87%B3%F0%9F%87%B1%20NL%20-%20%5Bt.me%2Ftest%5D"
	proxy, country, err := ParseConfigToClash(link)
	if err != nil {
		t.Fatalf("Failed to parse SS: %v", err)
	}

	if country != "NL" {
		t.Errorf("Expected country NL, got %s", country)
	}
	if proxy["type"] != "ss" {
		t.Errorf("Expected type ss, got %v", proxy["type"])
	}
	if proxy["cipher"] != "aes-256-gcm" {
		t.Errorf("Expected cipher aes-256-gcm, got %v", proxy["cipher"])
	}
	if proxy["password"] != "password123" {
		t.Errorf("Expected password password123, got %v", proxy["password"])
	}
}

func TestParseSS2022(t *testing.T) {
	// 2022 cipher with URL-encoded characters in key
	link := "ss://2022-blake3-chacha20-poly1305:vIFQQpkS5WcGd%2F5Sfpv3pt72Il7ReRZ3Z5FAzGM6e74=@129.151.74.15:61312#%F0%9F%87%AC%F0%9F%87%A7%20GB"
	proxy, country, err := ParseConfigToClash(link)
	if err != nil {
		t.Fatalf("Failed to parse SS 2022: %v", err)
	}

	if country != "GB" {
		t.Errorf("Expected country GB, got %s", country)
	}
	if proxy["password"] != "vIFQQpkS5WcGd/5Sfpv3pt72Il7ReRZ3Z5FAzGM6e74=" {
		t.Errorf("Expected unescaped base64 key, got %v", proxy["password"])
	}
}

func TestParseHy2(t *testing.T) {
	link := "hy2://pass123@hy2.example.com:443?sni=hy2.example.com&insecure=1#%F0%9F%87%B9%F0%9F%87%B7%20TR%20-%20%5Bt.me%2Ftest%5D"
	proxy, country, err := ParseConfigToClash(link)
	if err != nil {
		t.Fatalf("Failed to parse Hy2: %v", err)
	}

	if country != "TR" {
		t.Errorf("Expected country TR, got %s", country)
	}
	if proxy["type"] != "hysteria2" {
		t.Errorf("Expected type hysteria2, got %v", proxy["type"])
	}
	if proxy["skip-cert-verify"] != true {
		t.Errorf("Expected skip-cert-verify true, got %v", proxy["skip-cert-verify"])
	}
}

func TestParseTUIC(t *testing.T) {
	link := "tuic://uuid-1234:pass123@tuic.example.com:443?congestion_control=bbr&alpn=h3#%F0%9F%87%AC%F0%9F%87%A7%20GB%20-%20%5Bt.me%2Ftest%5D"
	proxy, country, err := ParseConfigToClash(link)
	if err != nil {
		t.Fatalf("Failed to parse TUIC: %v", err)
	}

	if country != "GB" {
		t.Errorf("Expected country GB, got %s", country)
	}
	if proxy["type"] != "tuic" {
		t.Errorf("Expected type tuic, got %v", proxy["type"])
	}
	if proxy["congestion-controller"] != "bbr" {
		t.Errorf("Expected congestion-controller bbr, got %v", proxy["congestion-controller"])
	}
}

func TestParseWireGuard_QueryParamKey(t *testing.T) {
	link := "wireguard://example.com:51820?private_key=cGFzc3dvcmQxMjM=&public_key=cHVibGljMTIzNDU=&ip=10.0.0.2#%F0%9F%87%A9%F0%9F%87%AA%20DE%20-%20%5Bt.me%2Ftest%5D"
	proxy, country, err := ParseConfigToClash(link)
	if err != nil {
		t.Fatalf("Failed to parse WireGuard with query private_key: %v", err)
	}

	if country != "DE" {
		t.Errorf("Expected country DE, got %s", country)
	}
	if proxy["type"] != "wireguard" {
		t.Errorf("Expected type wireguard, got %v", proxy["type"])
	}
	if proxy["private-key"] != "cGFzc3dvcmQxMjM=" {
		t.Errorf("Expected private-key cGFzc3dvcmQxMjM=, got %v", proxy["private-key"])
	}
	if proxy["public-key"] != "cHVibGljMTIzNDU=" {
		t.Errorf("Expected public-key cHVibGljMTIzNDU=, got %v", proxy["public-key"])
	}
}

func TestParseVLESS_RealityAliases(t *testing.T) {
	link := "vless://11111111-2222-3333-4444-555555555555@example.com:443?type=tcp&security=reality&public_key=aliaspubkey&short_id=5678&sni=target.com#%F0%9F%87%AB%F0%9F%87%B7%20FR"
	proxy, country, err := ParseConfigToClash(link)
	if err != nil {
		t.Fatalf("Failed to parse VLESS with reality aliases: %v", err)
	}

	if country != "FR" {
		t.Errorf("Expected country FR, got %s", country)
	}
	reality, ok := proxy["reality-opts"].(map[string]interface{})
	if !ok {
		t.Fatalf("reality-opts missing or invalid")
	}
	if reality["public-key"] != "aliaspubkey" {
		t.Errorf("Expected public-key aliaspubkey, got %v", reality["public-key"])
	}
	if reality["short-id"] != "5678" {
		t.Errorf("Expected short-id 5678, got %v", reality["short-id"])
	}
}

func TestParseSS_ChachaAlias(t *testing.T) {
	userInfo := base64.StdEncoding.EncodeToString([]byte("chacha20-poly1305:password123"))
	link := "ss://" + userInfo + "@ss.example.com:8388#%F0%9F%87%B3%F0%9F%87%B1%20NL"
	proxy, country, err := ParseConfigToClash(link)
	if err != nil {
		t.Fatalf("Failed to parse SS with chacha20-poly1305: %v", err)
	}

	if country != "NL" {
		t.Errorf("Expected country NL, got %s", country)
	}
	if proxy["cipher"] != "chacha20-ietf-poly1305" {
		t.Errorf("Expected normalized cipher chacha20-ietf-poly1305, got %v", proxy["cipher"])
	}
}

func TestParseSS_IPv6(t *testing.T) {
	userInfo := base64.StdEncoding.EncodeToString([]byte("aes-256-gcm:password123"))
	link := "ss://" + userInfo + "@[2001:db8::1]:8388#%F0%9F%87%A9%F0%9F%87%AA%20DE"
	proxy, _, err := ParseConfigToClash(link)
	if err != nil {
		t.Fatalf("Failed to parse IPv6 SS: %v", err)
	}

	if proxy["server"] != "2001:db8::1" {
		t.Errorf("Expected server 2001:db8::1, got %v", proxy["server"])
	}
	if proxy["port"] != 8388 {
		t.Errorf("Expected port 8388, got %v", proxy["port"])
	}
}

func TestGenerateClashConfig(t *testing.T) {
	links := []string{
		"vless://11111111-2222-3333-4444-555555555555@example.com:443?type=tcp&security=reality&pbk=testpubkey&sid=1234&sni=target.com#🇩🇪 DE - [t.me/test]",
		"vless://22222222-3333-4444-5555-666666666666@example2.com:443?type=tcp&security=reality&pbk=testpubkey&sid=1234&sni=target2.com#🇩🇪 DE - [t.me/test2]",
		"trojan://password123@trojan.example.com:443?security=tls&sni=trojan.example.com#🇫🇷 FR - [t.me/test]",
		"trojan://password456@trojan2.example.com:443?security=tls&sni=trojan2.example.com#🇫🇷 FR - [t.me/test2]",
		"vless://33333333-4444-5555-6666-777777777777@104.16.1.1:443?type=ws&security=tls#☁️ CDN/RELAY - [t.me/test]",
		"vless://44444444-5555-6666-7777-888888888888@104.16.1.2:443?type=ws&security=tls#☁️ CDN/RELAY - [t.me/test2]",
		"vless://55555555-6666-7777-8888-999999999999@ir.relay.example:443?type=tcp#🇮🇷 IR-RELAY - [t.me/test]",
		"vless://66666666-7777-8888-9999-000000000000@ir2.relay.example:443?type=tcp#🇮🇷 IR-RELAY - [t.me/test2]",
	}

	yamlData, err := GenerateClashConfig(links)
	if err != nil {
		t.Fatalf("Failed to generate clash config: %v", err)
	}

	yamlStr := string(yamlData)
	if !strings.Contains(yamlStr, "⚡ AUTO (Fastest Node)") {
		t.Errorf("Generated YAML missing AUTO group:\n%s", yamlStr)
	}
	if !strings.Contains(yamlStr, "🔄 FALLBACK (Failover)") {
		t.Errorf("Generated YAML missing FALLBACK group:\n%s", yamlStr)
	}
	if !strings.Contains(yamlStr, "⚖️ LOAD-BALANCE") {
		t.Errorf("Generated YAML missing LOAD-BALANCE group:\n%s", yamlStr)
	}
	if !strings.Contains(yamlStr, "🎯 MANUAL (All Nodes)") {
		t.Errorf("Generated YAML missing MANUAL group:\n%s", yamlStr)
	}
	if !strings.Contains(yamlStr, "GEOIP,IR,DIRECT") {
		t.Errorf("Generated YAML missing Iran direct rule:\n%s", yamlStr)
	}

	var parsed map[string]interface{}
	if err := yaml.Unmarshal(yamlData, &parsed); err != nil {
		t.Fatalf("Generated YAML failed to unmarshal: %v", err)
	}
}

func TestParseSSR(t *testing.T) {
	// ssr://1.2.3.4:8388:origin:aes-256-cfb:plain:cGFzc3dvcmQ=/?remarks=8J-HqfCfh6ogREUgLSBbdC5tZS90ZXN0XQ==&obfsparam=b2Jmcw==&protoparam=cHJvdG8=
	payload := base64.RawURLEncoding.EncodeToString([]byte("1.2.3.4:8388:origin:aes-256-cfb:plain:" + base64.RawURLEncoding.EncodeToString([]byte("password123")) + "/?remarks=" + base64.RawURLEncoding.EncodeToString([]byte("🇩🇪 DE - [t.me/test]")) + "&obfsparam=" + base64.RawURLEncoding.EncodeToString([]byte("obfsparam1")) + "&protoparam=" + base64.RawURLEncoding.EncodeToString([]byte("protoparam1"))))
	link := "ssr://" + payload

	proxy, country, err := ParseConfigToClash(link)
	if err != nil {
		t.Fatalf("Failed to parse SSR: %v", err)
	}

	if country != "DE" {
		t.Errorf("Expected country DE, got %s", country)
	}
	if proxy["type"] != "ssr" {
		t.Errorf("Expected type ssr, got %v", proxy["type"])
	}
	if proxy["server"] != "1.2.3.4" {
		t.Errorf("Expected server 1.2.3.4, got %v", proxy["server"])
	}
	if proxy["cipher"] != "aes-256-cfb" {
		t.Errorf("Expected cipher aes-256-cfb, got %v", proxy["cipher"])
	}
	if proxy["password"] != "password123" {
		t.Errorf("Expected password password123, got %v", proxy["password"])
	}
	if proxy["obfs-param"] != "obfsparam1" {
		t.Errorf("Expected obfs-param obfsparam1, got %v", proxy["obfs-param"])
	}
	if proxy["protocol-param"] != "protoparam1" {
		t.Errorf("Expected protocol-param protoparam1, got %v", proxy["protocol-param"])
	}
}

func TestParseSS_Plugins(t *testing.T) {
	// 1. v2ray-plugin
	userBase64 := base64.StdEncoding.EncodeToString([]byte("aes-256-gcm:pass123"))
	v2rayLink := "ss://" + userBase64 + "@1.2.3.4:8388?plugin=v2ray-plugin%3Bhost%3Dexample.com%3Bpath%3D%2Fws%3Btls%3D1#%F0%9F%87%A9%F0%9F%87%AA%20DE"
	proxyV2, _, err := ParseConfigToClash(v2rayLink)
	if err != nil {
		t.Fatalf("Failed to parse SS with v2ray-plugin: %v", err)
	}
	if proxyV2["plugin"] != "v2ray-plugin" {
		t.Errorf("Expected plugin v2ray-plugin, got %v", proxyV2["plugin"])
	}

	// 2. obfs-local
	obfsLink := "ss://" + userBase64 + "@1.2.3.4:8388?plugin=obfs-local%3Bobfs%3Dhttp%3Bobfs-host%3Dexample.com#%F0%9F%87%A9%F0%9F%87%AA%20DE"
	proxyObfs, _, err := ParseConfigToClash(obfsLink)
	if err != nil {
		t.Fatalf("Failed to parse SS with obfs: %v", err)
	}
	if proxyObfs["plugin"] != "obfs" {
		t.Errorf("Expected plugin obfs, got %v", proxyObfs["plugin"])
	}
}

func TestParseHysteria(t *testing.T) {
	link := "hysteria://1.2.3.4:443?auth=pass123&peer=example.com&upmbps=50&downmbps=100&alpn=h3&obfs=obfspass#%F0%9F%87%A9%F0%9F%87%AA%20DE"
	proxy, country, err := ParseConfigToClash(link)
	if err != nil {
		t.Fatalf("Failed to parse Hysteria: %v", err)
	}

	if country != "DE" {
		t.Errorf("Expected country DE, got %s", country)
	}
	if proxy["type"] != "hysteria" {
		t.Errorf("Expected type hysteria, got %v", proxy["type"])
	}
	if proxy["auth_str"] != "pass123" {
		t.Errorf("Expected auth_str pass123, got %v", proxy["auth_str"])
	}
	if proxy["up"] != "50" || proxy["down"] != "100" {
		t.Errorf("Expected up 50 down 100, got up=%v down=%v", proxy["up"], proxy["down"])
	}
}

func TestParseSocks(t *testing.T) {
	link := "socks5://user1:pass123@1.2.3.4:1080#%F0%9F%87%A9%F0%9F%87%AA%20DE"
	proxy, country, err := ParseConfigToClash(link)
	if err != nil {
		t.Fatalf("Failed to parse Socks: %v", err)
	}

	if country != "DE" {
		t.Errorf("Expected country DE, got %s", country)
	}
	if proxy["type"] != "socks5" {
		t.Errorf("Expected type socks5, got %v", proxy["type"])
	}
	if proxy["username"] != "user1" || proxy["password"] != "pass123" {
		t.Errorf("Expected user1:pass123, got user=%v pass=%v", proxy["username"], proxy["password"])
	}
}

func TestParseTransports_Clash(t *testing.T) {
	// 1. VLESS gRPC
	grpcLink := "vless://11111111-2222-3333-4444-555555555555@example.com:443?type=grpc&serviceName=mygrpcservice&security=tls#🇩🇪 DE"
	proxyGRPC, _, err := ParseConfigToClash(grpcLink)
	if err != nil {
		t.Fatalf("Failed to parse VLESS gRPC: %v", err)
	}
	if proxyGRPC["network"] != "grpc" {
		t.Errorf("Expected network grpc, got %v", proxyGRPC["network"])
	}
	if grpcOpts, ok := proxyGRPC["grpc-opts"].(map[string]interface{}); !ok || grpcOpts["grpc-service-name"] != "mygrpcservice" {
		t.Errorf("Expected grpc-service-name mygrpcservice, got %v", proxyGRPC["grpc-opts"])
	}

	// 2. VLESS HTTPUpgrade
	upgradeLink := "vless://11111111-2222-3333-4444-555555555555@example.com:443?type=httpupgrade&path=%2Fupgrade&host=example.com#🇩🇪 DE"
	proxyUpgrade, _, err := ParseConfigToClash(upgradeLink)
	if err != nil {
		t.Fatalf("Failed to parse VLESS HTTPUpgrade: %v", err)
	}
	if proxyUpgrade["network"] != "httpupgrade" {
		t.Errorf("Expected network httpupgrade, got %v", proxyUpgrade["network"])
	}

	// 3. VLESS H2
	h2Link := "vless://11111111-2222-3333-4444-555555555555@example.com:443?type=h2&path=%2Fh2path&host=example.com#🇩🇪 DE"
	proxyH2, _, err := ParseConfigToClash(h2Link)
	if err != nil {
		t.Fatalf("Failed to parse VLESS H2: %v", err)
	}
	if proxyH2["network"] != "h2" {
		t.Errorf("Expected network h2, got %v", proxyH2["network"])
	}

	// 4. VMess gRPC
	vmessGRPC := `{"v":"2","ps":"🇩🇪 DE","add":"1.2.3.4","port":443,"id":"11111111-2222-3333-4444-555555555555","aid":0,"net":"grpc","path":"grpc-service","tls":"tls"}`
	proxyVmessGRPC, _, err := ParseConfigToClash("vmess://" + base64.StdEncoding.EncodeToString([]byte(vmessGRPC)))
	if err != nil {
		t.Fatalf("Failed to parse VMess gRPC: %v", err)
	}
	if proxyVmessGRPC["network"] != "grpc" {
		t.Errorf("Expected VMess network grpc, got %v", proxyVmessGRPC["network"])
	}

	// 5. Trojan gRPC
	trojanGRPC := "trojan://mypass@1.2.3.4:443?type=grpc&serviceName=trojangrpc&security=tls#🇩🇪 DE"
	proxyTrojanGRPC, _, err := ParseConfigToClash(trojanGRPC)
	if err != nil {
		t.Fatalf("Failed to parse Trojan gRPC: %v", err)
	}
	if proxyTrojanGRPC["network"] != "grpc" {
		t.Errorf("Expected Trojan network grpc, got %v", proxyTrojanGRPC["network"])
	}
}

func TestParseWireGuard_AdvancedParams_Clash(t *testing.T) {
	wgLink := "wireguard://cGFzc3dvcmQ=@1.2.3.4:51820?public_key=cHVibGlj&ip=10.0.0.2&mtu=1420&reserved=1,2,3&preshared_key=cHJlc2hhcmVk#🇩🇪 DE"
	proxyWG, _, err := ParseConfigToClash(wgLink)
	if err != nil {
		t.Fatalf("Failed to parse WireGuard advanced: %v", err)
	}

	if proxyWG["mtu"] != 1420 {
		t.Errorf("Expected mtu 1420, got %v", proxyWG["mtu"])
	}
	if proxyWG["preshared-key"] != "cHJlc2hhcmVk" {
		t.Errorf("Expected preshared-key cHJlc2hhcmVk, got %v", proxyWG["preshared-key"])
	}
	if reserved, ok := proxyWG["reserved"].([]int); !ok || len(reserved) != 3 {
		t.Errorf("Expected reserved [1,2,3], got %v", proxyWG["reserved"])
	}
}



