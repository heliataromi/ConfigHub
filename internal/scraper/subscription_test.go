package scraper

import (
	"encoding/base64"
	"testing"

	"ConfigHub/internal/extractor"
)

func TestHasConfigs(t *testing.T) {
	tests := []struct {
		name     string
		configs  extractor.Configs
		expected bool
	}{
		{
			name:     "Empty",
			configs:  extractor.Configs{},
			expected: false,
		},
		{
			name: "Hy2 only",
			configs: extractor.Configs{
				Hy2: []string{"hy2://pass@example.com:443"},
			},
			expected: true,
		},
		{
			name: "TUIC only",
			configs: extractor.Configs{
				TUIC: []string{"tuic://uuid:pass@example.com:443"},
			},
			expected: true,
		},
		{
			name: "WireGuard only",
			configs: extractor.Configs{
				WireGuard: []string{"wireguard://privkey@example.com:51820?public_key=pubkey"},
			},
			expected: true,
		},
		{
			name: "SSR only",
			configs: extractor.Configs{
				SSR: []string{"ssr://base64encodedpayload"},
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := hasConfigs(tt.configs); got != tt.expected {
				t.Errorf("hasConfigs() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestDecodeBase64_Hy2Subscription(t *testing.T) {
	rawContent := "hy2://pass1@example.com:443?sni=example.com#Hy2Node1\nhy2://pass2@example.com:443?sni=example.com#Hy2Node2"
	encoded := base64.StdEncoding.EncodeToString([]byte(rawContent))

	decoded := decodeBase64(encoded)
	if decoded == "" {
		t.Fatalf("decodeBase64 failed on valid base64 subscription")
	}

	configs := extractor.Extract(decoded)
	if len(configs.Hy2) != 2 {
		t.Errorf("Expected 2 Hy2 configs from decoded sub, got %d", len(configs.Hy2))
	}
	if !hasConfigs(configs) {
		t.Errorf("hasConfigs returned false for Hy2 decoded configs")
	}
}
