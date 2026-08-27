package main

import (
    "fmt"
    "log"
    "net/url"
    "os"
    "path/filepath"
    "strings"
    "sync"

    "ConfigHub/internal/config"
    "ConfigHub/internal/exporter"
    "ConfigHub/internal/geoip"
    "ConfigHub/internal/parser"
    "ConfigHub/internal/scraper"
    "ConfigHub/internal/telemetry"
    "ConfigHub/internal/validator"

    "github.com/oschwald/maxminddb-golang"
)

// ConfigItem holds the raw link and the channel it was found in
type ConfigItem struct {
    Raw     string
    Channel string
    Source  string
}

func main() {
    // 0. Load optional local .env configuration
    config.LoadEnv(".env")

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
        "juicity":   make(map[string]ConfigItem),
        "naive":     make(map[string]ConfigItem),
        "telegram":  make(map[string]ConfigItem),
        "anytls":    make(map[string]ConfigItem),
        "snell":     make(map[string]ConfigItem),
        "http":      make(map[string]ConfigItem),
        "cottendns": make(map[string]ConfigItem),
        "stormdns":  make(map[string]ConfigItem),
    }

    var mu sync.Mutex
    var wg sync.WaitGroup

    // Helper function to deduplicate configs with channel source priority and telemetry tracking
    processConfigs := func(configs []string, proto string, channel string, source string) {
        for _, c := range configs {
            fp := parser.GetFingerprint(c)
            if existing, exists := uniqueConfigs[proto][fp]; !exists {
                uniqueConfigs[proto][fp] = ConfigItem{Raw: c, Channel: channel, Source: source}
            } else {
                if source == "channel" && existing.Source != "channel" {
                    // Channel source takes priority over external subscription
                    telemetry.Global.RecordDuplicate(proto, fp, channel, c, existing.Channel, existing.Raw)
                    uniqueConfigs[proto][fp] = ConfigItem{Raw: c, Channel: channel, Source: source}
                } else {
                    // Current config is duplicate of existing retained config
                    telemetry.Global.RecordDuplicate(proto, fp, existing.Channel, existing.Raw, channel, c)
                }
            }
        }
    }

    // 3. Scrape Concurrently with Bounded Worker Pools
    channelSem := make(chan struct{}, 6) // Max 6 concurrent requests to Telegram
    for _, ch := range channels {
        wg.Add(1)
        go func(channel string) {
            defer wg.Done()
            channelSem <- struct{}{}
            defer func() { <-channelSem }()

            configs, err := scraper.ScrapeChannel(channel)
            if err != nil {
                fmt.Printf("[-] Error scraping channel %s: %v\n", channel, err)
                return
            }

            channelName := "t.me/" + channel
            configs = validator.SanitizeConfigs(configs, channelName, telemetry.Global.RecordDropped)

            mu.Lock()
            processConfigs(configs.Vmess, "vmess", channelName, "channel")
            processConfigs(configs.Vless, "vless", channelName, "channel")
            processConfigs(configs.Trojan, "trojan", channelName, "channel")
            processConfigs(configs.SS, "ss", channelName, "channel")
            processConfigs(configs.SSR, "ssr", channelName, "channel")
            processConfigs(configs.TUIC, "tuic", channelName, "channel")
            processConfigs(configs.Hy2, "hy2", channelName, "channel")
            processConfigs(configs.Hysteria, "hysteria", channelName, "channel")
            processConfigs(configs.Socks, "socks", channelName, "channel")
            processConfigs(configs.WireGuard, "wireguard", channelName, "channel")
            processConfigs(configs.Juicity, "juicity", channelName, "channel")
            processConfigs(configs.Naive, "naive", channelName, "channel")
            processConfigs(configs.Telegram, "telegram", channelName, "channel")
            processConfigs(configs.AnyTLS, "anytls", channelName, "channel")
            processConfigs(configs.Snell, "snell", channelName, "channel")
            processConfigs(configs.HTTP, "http", channelName, "channel")
            processConfigs(configs.CottenDNS, "cottendns", channelName, "channel")
            processConfigs(configs.StormDNS, "stormdns", channelName, "channel")
            mu.Unlock()
            fmt.Printf("[+] %s: Scraped (%d configs)\n", channelName, configs.Count())
        }(ch)
    }

    subSem := make(chan struct{}, 8) // Max 8 concurrent requests for subscriptions
    for _, sub := range subscriptions {
        wg.Add(1)
        go func(subURL string) {
            defer wg.Done()
            subSem <- struct{}{}
            defer func() { <-subSem }()

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
                } else if parsedURL.Hostname() != "" {
                    shortName = parsedURL.Hostname()
                }
            }
            if shortName == subURL && len(subURL) > 30 {
                shortName = subURL[:30] + "..."
            }

            configs = validator.SanitizeConfigs(configs, shortName, telemetry.Global.RecordDropped)

            mu.Lock()
            processConfigs(configs.Vmess, "vmess", shortName, "subscription")
            processConfigs(configs.Vless, "vless", shortName, "subscription")
            processConfigs(configs.Trojan, "trojan", shortName, "subscription")
            processConfigs(configs.SS, "ss", shortName, "subscription")
            processConfigs(configs.SSR, "ssr", shortName, "subscription")
            processConfigs(configs.TUIC, "tuic", shortName, "subscription")
            processConfigs(configs.Hy2, "hy2", shortName, "subscription")
            processConfigs(configs.Hysteria, "hysteria", shortName, "subscription")
            processConfigs(configs.Socks, "socks", shortName, "subscription")
            processConfigs(configs.WireGuard, "wireguard", shortName, "subscription")
            processConfigs(configs.Juicity, "juicity", shortName, "subscription")
            processConfigs(configs.Naive, "naive", shortName, "subscription")
            processConfigs(configs.Telegram, "telegram", shortName, "subscription")
            processConfigs(configs.AnyTLS, "anytls", shortName, "subscription")
            processConfigs(configs.Snell, "snell", shortName, "subscription")
            processConfigs(configs.HTTP, "http", shortName, "subscription")
            processConfigs(configs.CottenDNS, "cottendns", shortName, "subscription")
            processConfigs(configs.StormDNS, "stormdns", shortName, "subscription")
            mu.Unlock()
            fmt.Printf("[+] %s: Scraped (%d configs)\n", shortName, configs.Count())
        }(sub)
    }

    wg.Wait()

    // 4. Rename configs using GeoIP & Channel ID (Concurrent Worker Pool)
    fmt.Println("\n[*] Resolving IPs and Applying GeoIP Names. This may take a moment...")

    finalConfigs := make(map[string][]exporter.RenamedConfig)
    for proto, cmap := range uniqueConfigs {
        renamed := processAndRename(cmap, db)
        if len(renamed) > 0 {
            finalConfigs[proto] = renamed
        }
    }

    // 5. Export Files
    fmt.Println("\n--- Deduplication & Renaming Complete ---")
    uniqueCounts := make(map[string]int)
    for proto, configs := range finalConfigs {
        count := len(configs)
        uniqueCounts[proto] = count
        fmt.Printf("%-10s: %d unique\n", strings.ToUpper(proto), count)
    }

    err = exporter.WriteSubFiles("sub", finalConfigs)
    if err != nil {
        log.Fatalf("[-] Error saving files: %v\n", err)
    }
    fmt.Println("[+] Successfully exported Geo-named files to 'sub/'!")

    // 6. Export Telemetry & Observability Reports (Only if EXPORT_TELEMETRY == "true" or "1")
    exportEnv := strings.ToLower(os.Getenv("EXPORT_TELEMETRY"))
    if exportEnv == "true" || exportEnv == "1" {
        totalUnique := 0
        for _, count := range uniqueCounts {
            totalUnique += count
        }

        reportsDir := "reports"
        if err := os.MkdirAll(reportsDir, 0755); err != nil {
            fmt.Printf("[-] Warning: Failed to create reports directory: %v\n", err)
        } else {
            telemetryPath := filepath.Join(reportsDir, "telemetry.json")
            duplicatesPath := filepath.Join(reportsDir, "duplicates.json")

            if err := telemetry.Global.ExportReport(telemetryPath, uniqueCounts); err != nil {
                fmt.Printf("[-] Warning: Failed to export telemetry report: %v\n", err)
            } else {
                fmt.Printf("[+] Successfully exported observability report to '%s'!\n", telemetryPath)
            }

            if err := telemetry.Global.ExportDuplicates(duplicatesPath, 0, totalUnique); err != nil {
                fmt.Printf("[-] Warning: Failed to export duplicates inspection report: %v\n", err)
            } else {
                fmt.Printf("[+] Successfully exported duplicate inspection report to '%s'!\n", duplicatesPath)
            }
        }
    }
    telemetry.Global.PrintConsoleSummary(uniqueCounts)
}

// processAndRename iterates through the map and renames each config concurrently
func processAndRename(configMap map[string]ConfigItem, db *maxminddb.Reader) []exporter.RenamedConfig {
    items := make([]ConfigItem, 0, len(configMap))
    for _, item := range configMap {
        items = append(items, item)
    }

    results := make([]exporter.RenamedConfig, len(items))
    var renameWg sync.WaitGroup
    workers := 16
    sem := make(chan struct{}, workers)

    for i, item := range items {
        renameWg.Add(1)
        go func(idx int, ci ConfigItem) {
            defer renameWg.Done()
            sem <- struct{}{}
            defer func() { <-sem }()

            renamed := parser.RenameConfig(ci.Raw, ci.Channel, db)
            results[idx] = exporter.RenamedConfig{URL: renamed, Source: ci.Source}
        }(i, item)
    }

    renameWg.Wait()

    var finalResults []exporter.RenamedConfig
    for _, res := range results {
        if res.URL != "" {
            finalResults = append(finalResults, res)
        }
    }
    return finalResults
}
