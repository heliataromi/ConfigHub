package extractor

import (
	"testing"
)

func TestExtract_AllProtocols(t *testing.T) {
	sampleText := `
Here are some free configs:
vless://11111111-2222-3333-4444-555555555555@example.com:443?type=tcp&security=tls#VlessNode
vmess://eyJhZGQiOiIxLjIuMy40IiwicG9ydCI6NDQzLCJpZCI6IjExMTExMTExLTIyMjItMzMzMy00NDQ0LTU1NTU1NTU1NTU1NSIsIm5ldCI6IndzIiwidHlwZSI6Im5vbmUiLCJob3N0IjoiZXhhbXBsZS5jb20iLCJwYXRoIjoiL3dzIiwidGxzIjoidGxzIiwicHMiOiJWTWVzc05vZGUifQ==
trojan://password123@trojan.example.com:443?security=tls#TrojanNode
ss://YWVzLTI1Ni1nY206cGFzc3dvcmQ=@ss.example.com:8388#SSNode
ssr://MS4yLjMuNDo4Mzg4Om9yaWdpbjphZXMtMjU2LWNmbjpwbGFpbjpkR1Z6ZEEvP29iZnNwYXJhbT0mcmVtYXJrcz1kR1Z6ZEE=
tuic://uuid1:pass1@tuic.example.com:443?congestion_control=bbr#TuicNode
hy2://pass1@hy2.example.com:443?sni=hy2.example.com#Hy2Node
hysteria://pass1@hysteria.example.com:443?protocol=udp#HysteriaNode
socks5://user:pass@socks.example.com:1080#SocksNode
wireguard://privkey1@wg.example.com:51820?public_key=pubkey1&ip=10.0.0.2#WGNode
juicity://uuid1:pass1@juicity.example.com:443?sni=juicity.example.com#JuicityNode
naive+https://user1:pass1@naive.example.com:443#NaiveNode
https://t.me/proxy?server=1.2.3.4&port=443&secret=ee1603010200010001fc030386e24c3add6d656469612e737465616d706f77657265642e636f6d
tg://proxy?server=5.6.7.8&port=2096&secret=dd79e344818749bd7ac519130220c25d09
anytls://11111111-2222-3333-4444-555555555555@anytls.example.com:443?security=tls#AnyTLSNode
snell://psk123@snell.example.com:443?version=4#SnellNode
http://proxyuser:proxypass@httpproxy.example.com:8080#HTTPNode
cottendns://eyJzY2hlbWEiOiJ3aGl0ZWRucy5wcm9maWxlIiwidmVyc2lvbiI6MSwicHJvZmlsZSI6eyJuYW1lIjoidC5tZVwvV2hpdGVETlMgQ290dGVuRE5T8J-HufCfh7cgdGh4IHRvIExvcmRvZkNpbmRlciIsInNlcnZlciI6eyJkb21haW4iOiJ2LmFzaGVudGFqaXIuc2JzLCBjLmFzaGVudGFqaXIuc2l0ZSIsImVuY3J5cHRpb25fa2V5IjoiZTU1NGI4ZmI4ZGU4Mjc4ZDJmMTFlODcwNDA0NDI2OWEiLCJlbmNyeXB0aW9uX21ldGhvZCI6M319fQ
stormdns://eyJzY2hlbWEiOiJ3aGl0ZWRucy5wcm9maWxlIiwidmVyc2lvbiI6MSwicHJvZmlsZSI6eyJuYW1lIjoidC5tZVwvV2hpdGVETlMgIPCfh7nwn4e3ICAgdGh4IHRvIENvcmVmb3JnZSIsInNlcnZlciI6eyJkb21haW4iOiJ2LmFub255bW91cy5vYnNlcnZlciIsImVuY3J5cHRpb25fa2V5IjoiYjI3NTAzOTE5OWIxYzhjOSIsImVuY3J5cHRpb25fbWV0aG9kIjozfX19
`

	configs := Extract(sampleText)

	if len(configs.Vless) != 1 {
		t.Errorf("Expected 1 VLESS config, got %d", len(configs.Vless))
	}
	if len(configs.Vmess) != 1 {
		t.Errorf("Expected 1 VMess config, got %d", len(configs.Vmess))
	}
	if len(configs.Trojan) != 1 {
		t.Errorf("Expected 1 Trojan config, got %d", len(configs.Trojan))
	}
	if len(configs.SS) != 1 {
		t.Errorf("Expected 1 SS config, got %d", len(configs.SS))
	}
	if len(configs.SSR) != 1 {
		t.Errorf("Expected 1 SSR config, got %d", len(configs.SSR))
	}
	if len(configs.TUIC) != 1 {
		t.Errorf("Expected 1 TUIC config, got %d", len(configs.TUIC))
	}
	if len(configs.Hy2) != 1 {
		t.Errorf("Expected 1 Hy2 config, got %d", len(configs.Hy2))
	}
	if len(configs.Hysteria) != 1 {
		t.Errorf("Expected 1 Hysteria config, got %d", len(configs.Hysteria))
	}
	if len(configs.Socks) != 1 {
		t.Errorf("Expected 1 Socks config, got %d", len(configs.Socks))
	}
	if len(configs.WireGuard) != 1 {
		t.Errorf("Expected 1 WireGuard config, got %d", len(configs.WireGuard))
	}
	if len(configs.Juicity) != 1 {
		t.Errorf("Expected 1 Juicity config, got %d", len(configs.Juicity))
	}
	if len(configs.Naive) != 1 {
		t.Errorf("Expected 1 Naive config, got %d", len(configs.Naive))
	}
	if len(configs.Telegram) != 2 {
		t.Errorf("Expected 2 Telegram proxy configs, got %d", len(configs.Telegram))
	}
	if len(configs.AnyTLS) != 1 {
		t.Errorf("Expected 1 AnyTLS config, got %d", len(configs.AnyTLS))
	}
	if len(configs.Snell) != 1 {
		t.Errorf("Expected 1 Snell config, got %d", len(configs.Snell))
	}
	if len(configs.HTTP) != 1 {
		t.Errorf("Expected 1 HTTP proxy config, got %d", len(configs.HTTP))
	}
	if len(configs.CottenDNS) != 1 {
		t.Errorf("Expected 1 CottenDNS config, got %d", len(configs.CottenDNS))
	}
	if len(configs.StormDNS) != 1 {
		t.Errorf("Expected 1 StormDNS config, got %d", len(configs.StormDNS))
	}

	if count := configs.Count(); count != 19 {
		t.Errorf("Expected total count 19, got %d", count)
	}
}

func TestExtract_TrailingDelimiters(t *testing.T) {
	text := `
[vless://11111111-2222-3333-4444-555555555555@example.com:443?type=tcp#Node1]
*hy2://pass1@example.com:443?sni=example.com#Node2*
_trojan://pass2@example.com:443?security=tls#Node3_
<tuic://uuid:pass@example.com:443#Node4>
vless://11111111-2222-3333-4444-555555555555@example.com:443?type=tcp#Node5\
`
	configs := Extract(text)

	if len(configs.Vless) != 2 {
		t.Fatalf("Expected 2 VLESS configs, got %d", len(configs.Vless))
	}
	if configs.Vless[0] != "vless://11111111-2222-3333-4444-555555555555@example.com:443?type=tcp#Node1" {
		t.Errorf("Expected trimmed fragment without bracket, got %s", configs.Vless[0])
	}
	if len(configs.Hy2) != 1 || configs.Hy2[0] != "hy2://pass1@example.com:443?sni=example.com#Node2" {
		t.Errorf("Expected trimmed hy2 without star, got %v", configs.Hy2)
	}
	if len(configs.Trojan) != 1 || configs.Trojan[0] != "trojan://pass2@example.com:443?security=tls#Node3" {
		t.Errorf("Expected trimmed trojan without underscore, got %v", configs.Trojan)
	}
	if len(configs.TUIC) != 1 || configs.TUIC[0] != "tuic://uuid:pass@example.com:443#Node4" {
		t.Errorf("Expected trimmed tuic without angle bracket, got %v", configs.TUIC)
	}
}

func TestExtract_PersianTextWithConfig(t *testing.T) {
	text := "کل پیام رو کپی کن، برو تو v2ray، دکمه مثبت (  ➕  ) رو بزن و Import configs from clipboard رو انتخاب کن. همهی کانفیگها یه جا میاد تو برنامهت  vless://435bda4c-fe5e-42c9-a3ad-15334943b38a@104.17.163.59:80?security=none&type=ws&host=us3.rtacg.com&path=/#  🃏   @v2ray_Dalghak"
	c := Extract(text)
	if len(c.Vless) != 1 {
		t.Fatalf("Expected 1 VLESS config, got %d", len(c.Vless))
	}
}

func TestAuditAndExtract_NoFalsePositiveOnValidInlineConfig(t *testing.T) {
	text := "کل پیام رو کپی کن، برو تو v2ray، دکمه مثبت (  ➕  ) رو بزن و Import configs from clipboard رو انتخاب کن. همهی کانفیگها یه جا میاد تو برنامهت  vless://435bda4c-fe5e-42c9-a3ad-15334943b38a@104.17.163.59:80?security=none&type=ws&host=us3.rtacg.com&path=/#  🃏   @v2ray_Dalghak"
	
	var dropped []string
	configs := AuditAndExtract(text, "t.me/v2ray_dalghak", func(source, candidate, reason string) {
		dropped = append(dropped, candidate)
	})

	if len(configs.Vless) != 1 {
		t.Fatalf("Expected 1 extracted VLESS config, got %d", len(configs.Vless))
	}
	if len(dropped) != 0 {
		t.Fatalf("Expected 0 false positive dropped items, got %d: %v", len(dropped), dropped)
	}
}

