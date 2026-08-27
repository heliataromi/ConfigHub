package singbox

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestGenerateSingboxConfig(t *testing.T) {
	links := []string{
		"vless://00000000-0000-0000-0000-000000000000@127.0.0.1:0?type=tcp#Last%20update", // Dummy should be skipped
		"vless://11111111-2222-3333-4444-555555555555@1.2.3.4:443?security=reality&sni=example.com&pbk=fakeKey&sid=1234#🇩🇪%20DE%20-%20[channel]",
		"vless://22222222-3333-4444-5555-666666666666@1.2.3.5:443?security=reality&sni=example.com&pbk=fakeKey&sid=1234#🇩🇪%20DE%20-%20[channel]", // Duplicate tag, should be renamed with 02
		"hy2://myPass@1.2.3.6:8443?sni=hy2.example.com#🇺🇸%20US%20-%20[channel]",
	}

	data, err := GenerateSingboxConfig(links)
	if err != nil {
		t.Fatalf("GenerateSingboxConfig returned error: %v", err)
	}

	if len(data) == 0 {
		t.Fatalf("generated data is empty")
	}

	// Validate JSON parseability
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		t.Fatalf("failed to unmarshal generated sing-box config: %v\nJSON:\n%s", err, string(data))
	}

	// Check DNS settings
	if cfg.DNS == nil || len(cfg.DNS.Servers) == 0 {
		t.Errorf("expected DNS servers to be configured")
	}
	if cfg.DNS.FakeIP == nil || !cfg.DNS.FakeIP.Enabled {
		t.Errorf("expected FakeIP to be enabled")
	}

	// Check Route rules
	if cfg.Route == nil || len(cfg.Route.Rules) == 0 {
		t.Errorf("expected Route rules to be configured")
	}
	if cfg.Route.Final != "🔰 PROXY" {
		t.Errorf("expected Route Final to be 🔰 PROXY, got %s", cfg.Route.Final)
	}

	// Check Outbounds
	var tags []string
	for _, ob := range cfg.Outbounds {
		if tag, ok := ob["tag"].(string); ok {
			tags = append(tags, tag)
		}
	}

	expectedTags := []string{
		"🔰 PROXY",
		"⚡ AUTO (Fastest Node)",
		"🔄 FALLBACK (Failover)",
		"⚖️ LOAD-BALANCE",
		"🎯 MANUAL (All Nodes)",
		"🇩🇪 DE - [channel] 01",
		"🇩🇪 DE - [channel] 02",
		"🇺🇸 US - [channel] 01",
		"direct",
		"block",
		"dns-out",
	}

	for _, expected := range expectedTags {
		found := false
		for _, tag := range tags {
			if tag == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("expected outbound tag %q not found in %v", expected, tags)
		}
	}

	// Ensure no dummy node was included
	for _, tag := range tags {
		if strings.Contains(tag, "Last update") || strings.Contains(tag, "Last%20update") {
			t.Errorf("dummy node was not filtered out: %s", tag)
		}
	}
}

func TestGenerateSingboxConfigEmpty(t *testing.T) {
	_, err := GenerateSingboxConfig([]string{})
	if err == nil {
		t.Errorf("expected error for empty links, got nil")
	}
}
