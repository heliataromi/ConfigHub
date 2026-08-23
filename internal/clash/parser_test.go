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

