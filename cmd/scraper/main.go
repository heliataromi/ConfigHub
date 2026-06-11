package main

import (
    "fmt"
    "log"
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
    uniqueVmess := make(map[string]ConfigItem)
    uniqueVless := make(map[string]ConfigItem)
    uniqueTrojan := make(map[string]ConfigItem)
    uniqueSS := make(map[string]ConfigItem)

    var mu sync.Mutex
    var wg sync.WaitGroup

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
            for _, c := range configs.Vmess {
                fp := parser.GetFingerprint(c)
                if _, exists := uniqueVmess[fp]; !exists {
                    uniqueVmess[fp] = ConfigItem{Raw: c, Channel: channel}
                }
            }
            for _, c := range configs.Vless {
                fp := parser.GetFingerprint(c)
                if _, exists := uniqueVless[fp]; !exists {
                    uniqueVless[fp] = ConfigItem{Raw: c, Channel: channel}
                }
            }
            for _, c := range configs.Trojan {
                fp := parser.GetFingerprint(c)
                if _, exists := uniqueTrojan[fp]; !exists {
                    uniqueTrojan[fp] = ConfigItem{Raw: c, Channel: channel}
                }
            }
            for _, c := range configs.SS {
                fp := parser.GetFingerprint(c)
                if _, exists := uniqueSS[fp]; !exists {
                    uniqueSS[fp] = ConfigItem{Raw: c, Channel: channel}
                }
            }
            mu.Unlock()
            fmt.Printf("[+] %s: Scraped\n", channel)
        }(ch)
    }

    for _, sub := range subscriptions {
        wg.Add(1)
        go func(url string) {
            defer wg.Done()
            configs, err := scraper.ScrapeSubscription(url)
            if err != nil {
                fmt.Printf("[-] Error scraping subscription %s: %v\n", url, err)
                return
            }

            // Create a short name for the URL to use as the channel name
            shortName := url
            if len(url) > 30 {
                shortName = url[:30] + "..."
            }

            mu.Lock()
            for _, c := range configs.Vmess {
                fp := parser.GetFingerprint(c)
                if _, exists := uniqueVmess[fp]; !exists {
                    uniqueVmess[fp] = ConfigItem{Raw: c, Channel: shortName}
                }
            }
            for _, c := range configs.Vless {
                fp := parser.GetFingerprint(c)
                if _, exists := uniqueVless[fp]; !exists {
                    uniqueVless[fp] = ConfigItem{Raw: c, Channel: shortName}
                }
            }
            for _, c := range configs.Trojan {
                fp := parser.GetFingerprint(c)
                if _, exists := uniqueTrojan[fp]; !exists {
                    uniqueTrojan[fp] = ConfigItem{Raw: c, Channel: shortName}
                }
            }
            for _, c := range configs.SS {
                fp := parser.GetFingerprint(c)
                if _, exists := uniqueSS[fp]; !exists {
                    uniqueSS[fp] = ConfigItem{Raw: c, Channel: shortName}
                }
            }
            mu.Unlock()
            fmt.Printf("[+] %s: Scraped\n", shortName)
        }(sub)
    }

    wg.Wait()

    // 4. Rename configs using GeoIP & Channel ID
    fmt.Println("\n[*] Resolving IPs and Applying GeoIP Names. This may take a moment...")

    finalVmess := processAndRename(uniqueVmess, db)
    finalVless := processAndRename(uniqueVless, db)
    finalTrojan := processAndRename(uniqueTrojan, db)
    finalSS := processAndRename(uniqueSS, db)

    // 5. Export Files
    fmt.Println("\n--- Deduplication & Renaming Complete ---")
    fmt.Printf("VMess:  %d unique\n", len(finalVmess))
    fmt.Printf("VLESS:  %d unique\n", len(finalVless))
    fmt.Printf("Trojan: %d unique\n", len(finalTrojan))
    fmt.Printf("SS:     %d unique\n", len(finalSS))

    err = exporter.WriteSubFiles("sub", finalVmess, finalVless, finalTrojan, finalSS)
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
