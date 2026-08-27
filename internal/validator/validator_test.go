package validator

import (
	"testing"
)

func TestValidateConfig_ValidProtocols(t *testing.T) {
	validList := []string{
		"vless://7aecb4a5-f0a4-32a0-aabe-c9d5241e313f@191.40.32.96:10018?security=reality&type=tcp#VLESSNode",
		"vmess://eyJhZGQiOiIxOTguNTEuMjAwLjEiLCJwb3J0Ijo0NDMsImlkIjoiN2FlY2I0YTUtZjBhNC0zMmEwLWFhYmUtYzlkNTI0MWUzMTNmIiwibmV0Ijoid3MiLCJ0bHMiOiJ0bHMiLCJwcyI6IlZNZXNzIn0=",
		"trojan://mypassword@198.51.200.1:443?security=tls#TrojanNode",
		"ss://YWVzLTI1Ni1nY206cGFzc3dvcmQ=@198.51.200.1:8388#SSNode",
		"hy2://mypassword@198.51.200.1:443?sni=example.com#Hy2Node",
		"tuic://7aecb4a5-f0a4-32a0-aabe-c9d5241e313f:pass1@198.51.200.1:443#TuicNode",
		"wireguard://cGFzc3dvcmQ=@198.51.200.1:51820?public_key=cHVibGlj&ip=10.0.0.2#WGNode",
		"socks5://user:pass@198.51.200.1:1080#SocksNode",
		"juicity://7aecb4a5-f0a4-32a0-aabe-c9d5241e313f:pass1@198.51.200.1:443?sni=example.com#JuicityNode",
		"naive+https://user1:pass1@198.51.200.1:443#NaiveNode",
		"https://t.me/proxy?server=198.51.200.1&port=443&secret=ee1603010200010001fc030386e24c3add6d656469612e737465616d706f77657265642e636f6d",
		"tg://proxy?server=198.51.200.1&port=2096&secret=dd79e344818749bd7ac519130220c25d09",
		"tg://socks?server=198.51.200.1&port=1080&user=usr&pass=pwd",
		"https://t.me/socks?server=198.51.200.1&port=1080",
		"tg://http?server=198.51.200.1&port=8080&user=usr&pass=pwd",
		"https://t.me/http?server=198.51.200.1&port=8080",
		"anytls://7aecb4a5-f0a4-32a0-aabe-c9d5241e313f@198.51.200.1:443?security=tls#AnyTLSNode",
		"snell://psk123456@198.51.200.1:443?version=4#SnellNode",
		"http://user1:pass1@198.51.200.1:8080#HTTPNode",
		"cottendns://eyJzY2hlbWEiOiJ3aGl0ZWRucy5wcm9maWxlIiwidmVyc2lvbiI6MSwicHJvZmlsZSI6eyJuYW1lIjoidC5tZVwvV2hpdGVETlMgQ290dGVuRE5T8J-HufCfh7cgdGh4IHRvIExvcmRvZkNpbmRlciIsInNlcnZlciI6eyJkb21haW4iOiJ2LmFzaGVudGFqaXIuc2JzLCBjLmFzaGVudGFqaXIuc2l0ZSIsImVuY3J5cHRpb25fa2V5IjoiZTU1NGI4ZmI4ZGU4Mjc4ZDJmMTFlODcwNDA0NDI2OWEiLCJlbmNyeXB0aW9uX21ldGhvZCI6M319fQ",
		"stormdns://eyJzY2hlbWEiOiJ3aGl0ZWRucy5wcm9maWxlIiwidmVyc2lvbiI6MSwicHJvZmlsZSI6eyJuYW1lIjoidC5tZVwvV2hpdGVETlMgIPCfh7nwn4e3ICAgdGh4IHRvIENvcmVmb3JnZSIsInNlcnZlciI6eyJkb21haW4iOiJ2LmFub255bW91cy5vYnNlcnZlciIsImVuY3J5cHRpb25fa2V5IjoiYjI3NTAzOTE5OWIxYzhjOSIsImVuY3J5cHRpb25fbWV0aG9kIjozfX19",
	}

	for _, link := range validList {
		ok, reason := ValidateConfig(link)
		if !ok {
			t.Errorf("Expected valid config to pass, but failed with: %s\nConfig: %s", reason, link)
		}
	}
}

func TestValidateConfig_DummyBanners(t *testing.T) {
	dummyList := []string{
		"vless://00000000-0000-0000-0000-000000000000@127.0.0.1:0?type=tcp#🐣🏁%20Last%20update:%20📅1405/06/02%20🕒00:58",
		"trojan://test@127.0.0.1:0#Last%20Update",
		"vmess://eyJhZGQiOiIxMjcuMC4wLjEiLCJwb3J0IjowLCJpZCI6IjAwMDAwMDAwLTAwMDAtMDAwMC0wMDAwLTAwMDAwMDAwMDAwMCIsInBzIjoiTGFzdCB1cGRhdGUifQ==",
	}

	for _, link := range dummyList {
		ok, reason := ValidateConfig(link)
		if ok {
			t.Errorf("Expected dummy banner to be rejected, but passed: %s", link)
		}
		t.Logf("Banner correctly rejected: %s", reason)
	}
}

func TestValidateConfig_InvalidHostsAndPorts(t *testing.T) {
	invalidList := []struct {
		link   string
		expect string
	}{
		{"vless://7aecb4a5-f0a4-32a0-aabe-c9d5241e313f@127.0.0.1:443#Loopback", "loopback"},
		{"vless://7aecb4a5-f0a4-32a0-aabe-c9d5241e313f@192.168.1.1:443#PrivateIP", "private"},
		{"vless://7aecb4a5-f0a4-32a0-aabe-c9d5241e313f@10.0.0.1:443#PrivateIP2", "private"},
		{"trojan://pass@example.com:0#PortZero", "port"},
		{"trojan://pass@example.com:70000#PortTooHigh", "port"},
		{"trojan://pass@example.com:abc#PortNonNumeric", "port"},
	}

	for _, item := range invalidList {
		ok, reason := ValidateConfig(item.link)
		if ok {
			t.Errorf("Expected rejection, but passed: %s", item.link)
		}
		t.Logf("Correctly rejected (%s): %s", reason, item.link)
	}
}

func TestValidateConfig_BogusUUIDs(t *testing.T) {
	bogusUUIDs := []string{
		"vless://00000000-0000-0000-0000-000000000000@198.51.200.1:443#AllZeros",
		"vless://YOUR_UUID_HERE@198.51.200.1:443#Placeholder",
		"vless://short-uuid@198.51.200.1:443#MalformedLength",
		"trojan://@198.51.200.1:443#EmptyPassword",
	}

	for _, link := range bogusUUIDs {
		ok, reason := ValidateConfig(link)
		if ok {
			t.Errorf("Expected bogus UUID/password to be rejected: %s", link)
		}
		t.Logf("Correctly rejected (%s): %s", reason, link)
	}
}

func TestValidateConfig_NewProtocolsEdgeCases(t *testing.T) {
	invalidCases := []struct {
		link string
		desc string
	}{
		{"tg://proxy?server=1.2.3.4&port=443&secret=short", "Telegram secret too short (< 20)"},
		{"tg://proxy?port=443&secret=ee1603010200010001fc030386e24c3add6d656469612e737465616d706f77657265642e636f6d", "Telegram missing server"},
		{"anytls://00000000-0000-0000-0000-000000000000@1.2.3.4:443", "AnyTLS all-zero dummy UUID"},
		{"anytls://@1.2.3.4:443", "AnyTLS missing UUID"},
		{"snell://@1.2.3.4:443", "Snell missing PSK"},
		{"http://1.2.3.4:8080", "HTTP proxy missing credentials"},
		{"cottendns://invalid-base64-payload", "CottenDNS invalid base64"},
		{"stormdns://eyJzY2hlbWEiOiJ3aGl0ZWRucy5wcm9maWxlIiwicHJvZmlsZSI6eyJzZXJ2ZXIiOnt9fX0=", "StormDNS missing domain"},
	}

	for _, tc := range invalidCases {
		ok, reason := ValidateConfig(tc.link)
		if ok {
			t.Errorf("Expected %s to be rejected, but it passed: %s", tc.desc, tc.link)
		} else {
			t.Logf("Correctly rejected %s (%s): %s", tc.desc, reason, tc.link)
		}
	}
}
