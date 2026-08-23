package geoip

import (
	"testing"
)

func TestGetFlag(t *testing.T) {
	tests := []struct {
		iso      string
		expected string
	}{
		{"DE", "🇩🇪"},
		{"US", "🇺🇸"},
		{"FR", "🇫🇷"},
		{"NL", "🇳🇱"},
		{"IR", "🇮🇷"},
		{"UNK", "🏳️"},
		{"", "🏳️"},
	}

	for _, tt := range tests {
		t.Run(tt.iso, func(t *testing.T) {
			if got := getFlag(tt.iso); got != tt.expected {
				t.Errorf("getFlag(%q) = %q, want %q", tt.iso, got, tt.expected)
			}
		})
	}
}

func TestGetCountry_NilDBAndCaching(t *testing.T) {
	// Query with nil DB should safely return UNK / 🏳️
	iso, flag := GetCountry("1.2.3.4", nil)
	if iso != "UNK" || flag != "🏳️" {
		t.Errorf("Expected UNK/🏳️ with nil DB, got %s/%s", iso, flag)
	}

	// Repeated query should hit cache
	iso2, flag2 := GetCountry("1.2.3.4", nil)
	if iso2 != "UNK" || flag2 != "🏳️" {
		t.Errorf("Expected cached UNK/🏳️, got %s/%s", iso2, flag2)
	}
}
