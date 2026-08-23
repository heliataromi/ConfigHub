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
	}

	err := WriteSubFiles(tmpDir, configsMap)
	if err != nil {
		t.Fatalf("WriteSubFiles failed: %v", err)
	}

	// 1. Verify standard files
	for _, name := range []string{"vless.txt", "vless_base64.txt", "mixed.txt", "clash.yaml"} {
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
}
