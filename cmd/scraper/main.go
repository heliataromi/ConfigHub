package main

import (
    "fmt"
    "log"
    "net/url"
    "strings"
    "sync"

    "ConfigHub/internal/config"
    "ConfigHub/internal/exporter"
    "ConfigHub/internal/geoip"
    "ConfigHub/internal/parser"
    "ConfigHub/internal/scraper"

    "github.com/oschwald/maxminddb-golang"
)

// ConfigItem holds the raw link and the channel it was found in
type ConfigItem struct {
    Raw     string
    Channel string
}

func main() {
    // 1. Prepare GeoIP Database
    err := geoip.EnsureDB()
    if err != nil {
        log.Fatalf("[-] Failed to download GeoIP DB: %v", err)
    }
    db, err := maxminddb.Open("GeoLite2-Country.mmdb")
    if err != nil {
        log.Fatalf("[-] Failed to open GeoIP DB: %v", err)
    }
    defer db.Close()

    // 2. Read Channels & Subscriptions
    channels, errC := config.ReadLines("channels.txt")
    if errC != nil {
        fmt.Println("[-] Warning: could not read channels.txt")
    }

    subscriptions, errS := config.ReadLines("subscriptions.txt")
    if errS != nil {
        fmt.Println("[-] Warning: could not read subscriptions.txt")
    }

    if len(channels) == 0 && len(subscriptions) == 0 {
        log.Fatalf("[-] Critical Error: Both channels.txt and subscriptions.txt are missing or empty.")
    }

    fmt.Printf("Loaded %d channels and %d subscriptions. Starting scraper...\n", len(channels), len(subscriptions))

    // Deduplication maps
    uniqueConfigs := map[string]map[string]ConfigItem{
        "vmess":     make(map[string]ConfigItem),
        "vless":     make(map[string]ConfigItem),
        "trojan":    make(map[string]ConfigItem),
        "ss":        make(map[string]ConfigItem),
        "ssr":       make(map[string]ConfigItem),
        "tuic":      make(map[string]ConfigItem),
        "hy2":       make(map[string]ConfigItem),
        "hysteria":  make(map[string]ConfigItem),
        "socks":     make(map[string]ConfigItem),
        "wireguard": make(map[string]ConfigItem),
    }

    var mu sync.Mutex
    var wg sync.WaitGroup

    // Helper function to deduplicate configs
    processConfigs := func(configs []string, proto string, channel string) {
        for _, c := range configs {
            fp := parser.GetFingerprint(c)
            if _, exists := uniqueConfigs[proto][fp]; !exists {
                uniqueConfigs[proto][fp] = ConfigItem{Raw: c, Channel: channel}
            }
        }
    }

    // 3. Scrape Concurrently
    for _, ch := range channels {
        wg.Add(1)
        go func(channel string) {
            defer wg.Done()
            configs, err := scraper.ScrapeChannel(channel)
            if err != nil {
                return
            }

            mu.Lock()
            processConfigs(configs.Vmess, "vmess", channel)
            processConfigs(configs.Vless, "vless", channel)
            processConfigs(configs.Trojan, "trojan", channel)
            processConfigs(configs.SS, "ss", channel)
            processConfigs(configs.SSR, "ssr", channel)
            processConfigs(configs.TUIC, "tuic", channel)
            processConfigs(configs.Hy2, "hy2", channel)
            processConfigs(configs.Hysteria, "hysteria", channel)
            processConfigs(configs.Socks, "socks", channel)
            processConfigs(configs.WireGuard, "wireguard", channel)
            mu.Unlock()
            fmt.Printf("[+] %s: Scraped\n", channel)
        }(ch)
    }

    for _, sub := range subscriptions {
        wg.Add(1)
        go func(subURL string) {
            defer wg.Done()
            configs, err := scraper.ScrapeSubscription(subURL)
            if err != nil {
                fmt.Printf("[-] Error scraping subscription %s: %v\n", subURL, err)
                return
            }

            // Create a short name for the subURL to use as the channel name
            shortName := subURL
            if parsedURL, err := url.Parse(subURL); err == nil {
                if parsedURL.Hostname() == "raw.githubusercontent.com" || parsedURL.Hostname() == "github.com" {
                    parts := strings.Split(strings.Trim(parsedURL.Path, "/"), "/")
                    if len(parts) > 0 {
                        shortName = "github.com/" + parts[0]
                    }
                }
            }
            if shortName == subURL && len(subURL) > 30 {
                shortName = subURL[:30] + "..."
            }

            mu.Lock()
            processConfigs(configs.Vmess, "vmess", shortName)
            processConfigs(configs.Vless, "vless", shortName)
            processConfigs(configs.Trojan, "trojan", shortName)
            processConfigs(configs.SS, "ss", shortName)
            processConfigs(configs.SSR, "ssr", shortName)
            processConfigs(configs.TUIC, "tuic", shortName)
            processConfigs(configs.Hy2, "hy2", shortName)
            processConfigs(configs.Hysteria, "hysteria", shortName)
            processConfigs(configs.Socks, "socks", shortName)
            processConfigs(configs.WireGuard, "wireguard", shortName)
            mu.Unlock()
            fmt.Printf("[+] %s: Scraped\n", shortName)
        }(sub)
    }

    wg.Wait()

    // 4. Rename configs using GeoIP & Channel ID
    fmt.Println("\n[*] Resolving IPs and Applying GeoIP Names. This may take a moment...")

    finalConfigs := make(map[string][]string)
    for proto, cmap := range uniqueConfigs {
        renamed := processAndRename(cmap, db)
        if len(renamed) > 0 {
            finalConfigs[proto] = renamed
        }
    }

    // 5. Export Files
    fmt.Println("\n--- Deduplication & Renaming Complete ---")
    for proto, configs := range finalConfigs {
        fmt.Printf("%-10s: %d unique\n", strings.ToUpper(proto), len(configs))
    }

    err = exporter.WriteSubFiles("sub", finalConfigs)
    if err != nil {
        log.Fatalf("[-] Error saving files: %v\n", err)
    }
    fmt.Println("[+] Successfully exported Geo-named files to 'sub/'!")
}

// processAndRename iterates through the map and renames each config
func processAndRename(configMap map[string]ConfigItem, db *maxminddb.Reader) []string {
    var results []string
    for _, item := range configMap {
        renamed := parser.RenameConfig(item.Raw, item.Channel, db)
        results = append(results, renamed)
    }
    return results
}
