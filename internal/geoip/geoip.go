package geoip

import (
    "context"
    "fmt"
    "io"
    "net"
    "net/http"
    "os"
    "strings"
    "sync"
    "time"

    "github.com/oschwald/maxminddb-golang"
)

const dbURL = "https://raw.githubusercontent.com/P3TERX/GeoLite.mmdb/download/GeoLite2-Country.mmdb"
const dbPath = "GeoLite2-Country.mmdb"

type geoResult struct {
    iso  string
    flag string
}

var geoCache sync.Map

// EnsureDB checks if the database exists; if not, it downloads it.
func EnsureDB() error {
    if _, err := os.Stat(dbPath); err == nil {
        return nil // File exists
    }

    fmt.Println("[*] Downloading GeoIP Database (this only happens once)...")
    resp, err := http.Get(dbURL)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    file, err := os.Create(dbPath)
    if err != nil {
        return err
    }
    defer file.Close()

    _, err = io.Copy(file, resp.Body)
    return err
}

// GetCountry returns the Country ISO Code and Emoji Flag for a given IP or Domain, using cached lookups when available.
func GetCountry(address string, db *maxminddb.Reader) (string, string) {
    address = strings.TrimSpace(address)
    address = strings.Trim(address, "[]")
    if address == "" {
        return "UNK", "🏳️"
    }

    if val, ok := geoCache.Load(address); ok {
        res := val.(geoResult)
        return res.iso, res.flag
    }

    ip := net.ParseIP(address)
    if ip == nil {
        // It's a domain, resolve it with a timeout
        ctx, cancel := context.WithTimeout(context.Background(), 1500*time.Millisecond)
        defer cancel()
        ips, err := net.DefaultResolver.LookupIPAddr(ctx, address)
        if err != nil || len(ips) == 0 {
            res := geoResult{iso: "UNK", flag: "🏳️"}
            geoCache.Store(address, res)
            return "UNK", "🏳️" // Unknown / Dead Domain
        }
        ip = ips[0].IP
    }

    if db == nil {
        return "UNK", "🏳️"
    }

    // Query the DB
    var record struct {
        Country struct {
            IsoCode string `maxminddb:"iso_code"`
        } `maxminddb:"country"`
    }

    err := db.Lookup(ip, &record)
    if err != nil || record.Country.IsoCode == "" {
        res := geoResult{iso: "CDN/RELAY", flag: "☁️"}
        geoCache.Store(address, res)
        return "CDN/RELAY", "☁️" // Likely a CDN or unrecognized IP
    }

    iso := record.Country.IsoCode
    flag := getFlag(iso)

    // If the server is physically in Iran, label it as a Domestic Relay
    if iso == "IR" {
        iso = "IR-RELAY"
    }

    res := geoResult{iso: iso, flag: flag}
    geoCache.Store(address, res)
    return iso, flag
}

// getFlag converts a 2-letter ISO code into an Emoji Flag
func getFlag(iso string) string {
    if len(iso) != 2 {
        return "🏳️"
    }
    iso = strings.ToUpper(iso)
    r1 := rune(iso[0]) - 'A' + 0x1F1E6
    r2 := rune(iso[1]) - 'A' + 0x1F1E6
    return string(r1) + string(r2)
}
