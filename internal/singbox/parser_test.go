package singbox

import (
	"testing"
)

func TestParseConfigToSingbox(t *testing.T) {
	tests := []struct {
		name         string
		link         string
		expectedType string
		expectErr    bool
	}{
		{
			name:         "VLESS Reality TCP",
			link:         "vless://11111111-2222-3333-4444-555555555555@1.2.3.4:443?security=reality&sni=example.com&fp=chrome&pbk=fakePubkey123&sid=1234abcd&flow=xtls-rprx-vision#🇩🇪%20DE%20-%20[t.me/test]",
			expectedType: "vless",
			expectErr:    false,
		},
		{
			name:         "VLESS WS TLS",
			link:         "vless://11111111-2222-3333-4444-555555555555@1.2.3.4:443?type=ws&security=tls&sni=cdn.example.com&path=%2Fvless-ws&host=cdn.example.com#🇺🇸%20US",
			expectedType: "vless",
			expectErr:    false,
		},
		{
			name:         "VMess WS TLS",
			link:         "vmess://eyJhZGQiOiIxLjIuMy40IiwicG9ydCI6NDQzLCJpZCI6IjExMTExMTExLTIyMjItMzMzMy00NDQ0LTU1NTU1NTU1NTU1NSIsImFpZCI6MCwic2N5IjoiYXV0byIsIm5ldCI6IndzIiwidHlwZSI6Im5vbmUiLCJob3N0IjoiY2RuLmV4YW1wbGUuY29tIiwicGF0aCI6Ii92bWVzcyIsInRscyI6InRscyIsInBzIjoi8J+HqfCfh6YgREUgLSBbY2hhbm5lbF0ifQ==",
			expectedType: "vmess",
			expectErr:    false,
		},
		{
			name:         "Trojan WS TLS",
			link:         "trojan://myPassword@1.2.3.4:443?security=tls&sni=trojan.example.com&type=ws&path=%2Ftrojan-ws#🇫🇷%20FR",
			expectedType: "trojan",
			expectErr:    false,
		},
		{
			name:         "Hysteria 2",
			link:         "hy2://myHy2Pass@1.2.3.4:8443?sni=hy2.example.com&obfs=salamander&obfs-password=secret#🇳🇱%20NL",
			expectedType: "hysteria2",
			expectErr:    false,
		},
		{
			name:         "TUIC",
			link:         "tuic://11111111-2222-3333-4444-555555555555:myTuicPass@1.2.3.4:8443?sni=tuic.example.com&congestion_controller=bbr#🇬🇧%20GB",
			expectedType: "tuic",
			expectErr:    false,
		},
		{
			name:         "Shadowsocks SIP002",
			link:         "ss://YWVzLTI1Ni1nY206cGFzc3dvcmRAMS4yLjMuNDo4Mzg4#🇯🇵%20JP",
			expectedType: "shadowsocks",
			expectErr:    false,
		},
		{
			name:         "WireGuard",
			link:         "wireguard://aGVsbG9wcml2YXRla2V5MTIzNA==@1.2.3.4:51820?publickey=aGVsbG9wdWJsaWNrZXkxMjM0&ip=10.0.0.2%2F32#🇨🇦%20CA",
			expectedType: "wireguard",
			expectErr:    false,
		},
		{
			name:         "Unsupported Scheme",
			link:         "ftp://1.2.3.4:21",
			expectedType: "",
			expectErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			outbound, tag, err := ParseConfigToSingbox(tt.link)
			if tt.expectErr {
				if err == nil {
					t.Fatalf("expected error for %s, got nil", tt.link)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error for %s: %v", tt.link, err)
			}
			if outbound["type"] != tt.expectedType {
				t.Errorf("expected type %s, got %v", tt.expectedType, outbound["type"])
			}
			if tag == "" {
				t.Errorf("expected non-empty tag, got empty")
			}
		})
	}
}

func TestVLESSRealityDetails(t *testing.T) {
	link := "vless://11111111-2222-3333-4444-555555555555@1.2.3.4:443?security=reality&sni=example.com&fp=chrome&pbk=myPublicKey123&sid=abcd1234&flow=xtls-rprx-vision#TestReality"
	outbound, tag, err := ParseConfigToSingbox(link)
	if err != nil {
		t.Fatalf("failed to parse: %v", err)
	}

	if tag != "TestReality" {
		t.Errorf("expected tag TestReality, got %s", tag)
	}

	if outbound["flow"] != "xtls-rprx-vision" {
		t.Errorf("expected flow xtls-rprx-vision, got %v", outbound["flow"])
	}

	tlsMap, ok := outbound["tls"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected tls map in outbound")
	}

	realityMap, ok := tlsMap["reality"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected reality map in tls")
	}

	if realityMap["public_key"] != "myPublicKey123" {
		t.Errorf("expected public_key myPublicKey123, got %v", realityMap["public_key"])
	}
	if realityMap["short_id"] != "abcd1234" {
		t.Errorf("expected short_id abcd1234, got %v", realityMap["short_id"])
	}
}
