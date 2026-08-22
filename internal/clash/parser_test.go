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

func TestGenerateClashConfig(t *testing.T) {
	links := []string{
		"vless://11111111-2222-3333-4444-555555555555@example.com:443?type=tcp&security=reality&pbk=testpubkey&sid=1234&sni=target.com#🇩🇪 DE - [t.me/test]",
		"trojan://password123@trojan.example.com:443?security=tls&sni=trojan.example.com#🇫🇷 FR - [t.me/test]",
	}

	yamlData, err := GenerateClashConfig(links)
	if err != nil {
		t.Fatalf("Failed to generate clash config: %v", err)
	}

	yamlStr := string(yamlData)
	if !strings.Contains(yamlStr, "🇩🇪 Germany (Auto)") {
		t.Errorf("Generated YAML missing German auto group:\n%s", yamlStr)
	}
	if !strings.Contains(yamlStr, "🇫🇷 France (Auto)") {
		t.Errorf("Generated YAML missing French auto group:\n%s", yamlStr)
	}
	if !strings.Contains(yamlStr, "⚡ AUTO (Fastest Node)") {
		t.Errorf("Generated YAML missing AUTO group:\n%s", yamlStr)
	}
	if !strings.Contains(yamlStr, "GEOIP,IR,DIRECT") {
		t.Errorf("Generated YAML missing Iran direct rule:\n%s", yamlStr)
	}

	var parsed map[string]interface{}
	if err := yaml.Unmarshal(yamlData, &parsed); err != nil {
		t.Fatalf("Generated YAML failed to unmarshal: %v", err)
	}
}

