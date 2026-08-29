package main

import (
	"testing"
)

func TestProcessAndRename(t *testing.T) {
	configMap := map[string]ConfigItem{
		"fp1": {
			Raw:     "vless://7aecb4a5-f0a4-32a0-aabe-c9d5241e313f@191.40.32.96:10018?security=reality&type=tcp#OldName",
			Channel: "t.me/testchan",
			Source:  "channel",
		},
		"fp2": {
			Raw:     "trojan://mypass@1.2.3.4:443?security=tls#OldName",
			Channel: "sub.example.com",
			Source:  "subscription",
		},
	}

	results := processAndRename(configMap, nil)
	if len(results) != 2 {
		t.Fatalf("Expected 2 renamed configs, got %d", len(results))
	}

	for _, res := range results {
		if res.URL == "" {
			t.Errorf("Expected non-empty URL")
		}
		if res.Source != "channel" && res.Source != "subscription" {
			t.Errorf("Unexpected source: %s", res.Source)
		}
	}
}

func TestProcessAndRename_Empty(t *testing.T) {
	results := processAndRename(map[string]ConfigItem{}, nil)
	if len(results) != 0 {
		t.Errorf("Expected 0 results for empty map, got %d", len(results))
	}
}
