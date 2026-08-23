package exporter

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestWriteSubFiles_WithCountrySubs(t *testing.T) {
	tmpDir := t.TempDir()

	configsMap := map[string][]RenamedConfig{
		"vless": {
			{URL: "vless://uuid1@1.2.3.4:443?type=tcp#🇩🇪 DE - [t.me/test]", Source: "channel"},
			{URL: "vless://uuid2@1.2.3.5:443?type=tcp#🇺🇸 US - [t.me/test2]", Source: "channel"},
			{URL: "vless://uuid3@1.2.3.6:443?type=tcp#☁️ CDN/RELAY - [t.me/test3]", Source: "subscription"},
		},
		"trojan": {
			{URL: "trojan://pass1@1.2.3.7:443?security=tls#🇩🇪 DE - [t.me/test]", Source: "channel"},
		},
		"telegram": {
			{URL: "tg://proxy?server=1.2.3.8&port=443&secret=ee1603010200010001fc030386e24c3add6d656469612e737465616d706f77657265642e636f6d#%F0%9F%87%A9%F0%9F%87%AA%20DE%20-%20%5Bt.me%2Ftest%5D", Source: "channel"},
		},
		"cottendns": {
			{URL: "cottendns://eyJzY2hlbWEiOiJ3aGl0ZWRucy5wcm9maWxlIiwidmVyc2lvbiI6MSwicHJvZmlsZSI6eyJuYW1lIjoiVFIiLCJzZXJ2ZXIiOnsiZG9tYWluIjoidHJ1ZW1haWwuY29tIiwiZW5jcnlwdGlvbl9rZXkiOiJrZXkxIiwiZW5jcnlwdGlvbl9tZXRob2QiOjF9fX0=", Source: "channel"},
		},
		"stormdns": {
			{URL: "stormdns://eyJzY2hlbWEiOiJ3aGl0ZWRucy5wcm9maWxlIiwidmVyc2lvbiI6MSwicHJvZmlsZSI6eyJuYW1lIjoiVVMiLCJzZXJ2ZXIiOnsiZG9tYWluIjoidXNhaG9zdC5jb20iLCJlbmNyeXB0aW9uX2tleSI6ImtleTIiLCJlbmNyeXB0aW9uX21ldGhvZCI6MX19fQ==", Source: "channel"},
		},
	}

	err := WriteSubFiles(tmpDir, configsMap)
	if err != nil {
		t.Fatalf("WriteSubFiles failed: %v", err)
	}

	// 1. Verify standard files
	for _, name := range []string{"vless.txt", "vless_base64.txt", "telegram.txt", "telegram_base64.txt", "mtproto.txt", "mtproto_base64.txt", "cottendns.txt", "stormdns.txt", "whitedns.txt", "mixed.txt", "clash.yaml"} {
		path := filepath.Join(tmpDir, name)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			t.Errorf("Expected standard file %s to exist", name)
		}
	}

	// 2. Verify country sub files in sub/countries/
	countriesDir := filepath.Join(tmpDir, "countries")
	for _, countryFile := range []string{"de.txt", "de_base64.txt", "us.txt", "us_base64.txt", "cdn.txt", "cdn_base64.txt"} {
		path := filepath.Join(countriesDir, countryFile)
		data, err := os.ReadFile(path)
		if err != nil {
			t.Errorf("Expected country file %s to exist: %v", countryFile, err)
			continue
		}
		if !strings.Contains(string(data), "Last%20update") && !strings.HasSuffix(countryFile, "_base64.txt") {
			t.Errorf("Expected timestamp header in %s", countryFile)
		}
	}

	// Check that de.txt contains both vless and trojan German nodes
	deContent, _ := os.ReadFile(filepath.Join(countriesDir, "de.txt"))
	if !strings.Contains(string(deContent), "vless://uuid1") || !strings.Contains(string(deContent), "trojan://pass1") {
		t.Errorf("Expected de.txt to contain both German nodes, got:\n%s", string(deContent))
	}

	// Verify telegram proxies and whitedns protocols are excluded from mixed.txt
	mixedContent, _ := os.ReadFile(filepath.Join(tmpDir, "mixed.txt"))
	if strings.Contains(string(mixedContent), "tg://proxy?") || strings.Contains(string(mixedContent), "cottendns://") || strings.Contains(string(mixedContent), "stormdns://") {
		t.Errorf("Telegram or WhiteDNS proxies must not be present in mixed.txt:\n%s", string(mixedContent))
	}

	// Verify telegram.txt and whitedns.txt do NOT contain dummy vless header
	tgContent, _ := os.ReadFile(filepath.Join(tmpDir, "telegram.txt"))
	if strings.Contains(string(tgContent), "vless://00000000-0000") {
		t.Errorf("telegram.txt must not contain dummy vless config:\n%s", string(tgContent))
	}
	if !strings.Contains(string(tgContent), "tg://proxy?server=1.2.3.8") {
		t.Errorf("telegram.txt missing telegram proxy link:\n%s", string(tgContent))
	}

	whitednsContent, _ := os.ReadFile(filepath.Join(tmpDir, "whitedns.txt"))
	if strings.Contains(string(whitednsContent), "vless://00000000-0000") {
		t.Errorf("whitedns.txt must not contain dummy vless config:\n%s", string(whitednsContent))
	}
	if !strings.Contains(string(whitednsContent), "cottendns://") || !strings.Contains(string(whitednsContent), "stormdns://") {
		t.Errorf("whitedns.txt must contain both CottenDNS and StormDNS configs:\n%s", string(whitednsContent))
	}
}
