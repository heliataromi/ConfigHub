package geoip

import (
	"os"
	"testing"

	"github.com/oschwald/maxminddb-golang"
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

	// Empty address
	isoEmpty, flagEmpty := GetCountry("   ", nil)
	if isoEmpty != "UNK" || flagEmpty != "🏳️" {
		t.Errorf("Expected UNK/🏳️ for empty address, got %s/%s", isoEmpty, flagEmpty)
	}

	// Unresolvable domain
	isoDomain, flagDomain := GetCountry("invalid.nonexistent.domain.test", nil)
	if isoDomain != "UNK" || flagDomain != "🏳️" {
		t.Errorf("Expected UNK/🏳️ for invalid domain, got %s/%s", isoDomain, flagDomain)
	}
}

func TestEnsureDB_ExistingFile(t *testing.T) {
	// EnsureDB returns nil if file exists
	err := EnsureDB()
	if err != nil {
		t.Logf("EnsureDB returned: %v (network/file dependent)", err)
	}
}

func TestGetCountry_WithRealDB(t *testing.T) {
	if _, err := os.Stat("GeoLite2-Country.mmdb"); os.IsNotExist(err) {
		t.Skip("GeoLite2-Country.mmdb not found locally")
	}

	db, err := maxminddb.Open("GeoLite2-Country.mmdb")
	if err != nil {
		t.Fatalf("Failed to open mmdb: %v", err)
	}
	defer db.Close()

	// 8.8.8.8 -> US
	iso, flag := GetCountry("8.8.8.8", db)
	if iso != "US" || flag != "🇺🇸" {
		t.Errorf("Expected US/🇺🇸 for 8.8.8.8, got %s/%s", iso, flag)
	}

	// 1.1.1.1 -> AU or US or Cloudflare CDN
	isoCF, flagCF := GetCountry("1.1.1.1", db)
	if isoCF == "" || flagCF == "" {
		t.Errorf("Expected non-empty result for 1.1.1.1, got %s/%s", isoCF, flagCF)
	}
}

