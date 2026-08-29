package scraper

import (
	"encoding/base64"
	"net/http"
	"net/http/httptest"
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

func TestDecodeBase64_Variants(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Standard base64",
			input:    base64.StdEncoding.EncodeToString([]byte("vless://test@1.1.1.1:443#Standard")),
			expected: "vless://test@1.1.1.1:443#Standard",
		},
		{
			name:     "Raw Std base64 (unpadded)",
			input:    base64.RawStdEncoding.EncodeToString([]byte("vmess://test@1.1.1.1:443#RawStd")),
			expected: "vmess://test@1.1.1.1:443#RawStd",
		},
		{
			name:     "URL encoding",
			input:    base64.URLEncoding.EncodeToString([]byte("trojan://test@1.1.1.1:443#URL")),
			expected: "trojan://test@1.1.1.1:443#URL",
		},
		{
			name:     "Raw URL encoding",
			input:    base64.RawURLEncoding.EncodeToString([]byte("ss://test@1.1.1.1:443#RawURL")),
			expected: "ss://test@1.1.1.1:443#RawURL",
		},
		{
			name:     "With newlines and spaces",
			input:    " \n" + base64.StdEncoding.EncodeToString([]byte("hy2://test@1.1.1.1:443#Spaced")) + "\r\n ",
			expected: "hy2://test@1.1.1.1:443#Spaced",
		},
		{
			name:     "Broken base64",
			input:    "!!!NOT_BASE64@@@",
			expected: "",
		},
		{
			name:     "Empty input",
			input:    "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := decodeBase64(tt.input)
			if got != tt.expected {
				t.Errorf("decodeBase64(%q) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}

func TestScrapeSubscription(t *testing.T) {
	tests := []struct {
		name          string
		statusCode    int
		responseBody  string
		expectError   bool
		expectedCount int
	}{
		{
			name:          "200 OK with plain text configs",
			statusCode:    http.StatusOK,
			responseBody:  "vless://7aecb4a5-f0a4-32a0-aabe-c9d5241e313f@191.40.32.96:10018?security=reality&type=tcp#Node1\ntrojan://pass@1.1.1.1:443#Node2",
			expectError:   false,
			expectedCount: 2,
		},
		{
			name:          "200 OK with base64 encoded subscription",
			statusCode:    http.StatusOK,
			responseBody:  base64.StdEncoding.EncodeToString([]byte("ss://YWVzLTI1Ni1nY206cGFzc3dvcmQ=@1.1.1.1:8388#Node1\nhy2://pass@1.1.1.1:443#Node2")),
			expectError:   false,
			expectedCount: 2,
		},
		{
			name:          "200 OK with empty body",
			statusCode:    http.StatusOK,
			responseBody:  "",
			expectError:   false,
			expectedCount: 0,
		},
		{
			name:          "404 Not Found error",
			statusCode:    http.StatusNotFound,
			responseBody:  "Not Found",
			expectError:   true,
			expectedCount: 0,
		},
		{
			name:          "500 Internal Server Error",
			statusCode:    http.StatusInternalServerError,
			responseBody:  "Server Error",
			expectError:   true,
			expectedCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.statusCode)
				_, _ = w.Write([]byte(tt.responseBody))
			}))
			defer server.Close()

			configs, err := ScrapeSubscription(server.URL)
			if (err != nil) != tt.expectError {
				t.Fatalf("ScrapeSubscription() error = %v, expectError = %v", err, tt.expectError)
			}
			if configs.Count() != tt.expectedCount {
				t.Errorf("ScrapeSubscription() count = %d, want %d", configs.Count(), tt.expectedCount)
			}
		})
	}
}
