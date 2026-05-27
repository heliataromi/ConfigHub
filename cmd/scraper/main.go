package main

import (
    "fmt"
    "log"
    "sync"

    "ConfigHub/internal/config"
    "ConfigHub/internal/scraper"
)

func main() {
    // 1. Read the channels from the text file
    channels, err := config.ReadChannels("channels.txt")
    if err != nil {
        log.Fatalf("[-] Critical Error: Could not read channels.txt: %v\n", err)
    }

    if len(channels) == 0 {
        log.Fatalf("[-] No channels found in channels.txt. Please add some!")
    }

    fmt.Printf("Loaded %d channels from file. Starting scraper...\n", len(channels))

    // 2. Setup maps for deduplication
    uniqueVmess := make(map[string]bool)
    uniqueVless := make(map[string]bool)
    uniqueTrojan := make(map[string]bool)
    uniqueSS := make(map[string]bool)

    var mu sync.Mutex
    var wg sync.WaitGroup

    // 3. Scrape concurrently
    for _, ch := range channels {
        wg.Add(1)

        go func(channel string) {
            defer wg.Done()

            configs, err := scraper.ScrapeChannel(channel)
            if err != nil {
                fmt.Printf("[-] Failed to scrape %s: %v\n", channel, err)
                return
            }

            mu.Lock()
            for _, c := range configs.Vmess {
                uniqueVmess[c] = true
            }
            for _, c := range configs.Vless {
                uniqueVless[c] = true
            }
            for _, c := range configs.Trojan {
                uniqueTrojan[c] = true
            }
            for _, c := range configs.SS {
                uniqueSS[c] = true
            }
            mu.Unlock()

            totalFound := len(configs.Vmess) + len(configs.Vless) + len(configs.Trojan) + len(configs.SS)
            fmt.Printf("[+] %s: Found %d configs in the last 24h\n", channel, totalFound)
        }(ch)
    }

    wg.Wait()

    // 4. Print final output
    fmt.Println("\n--- Final Deduplicated Results ---")
    fmt.Printf("VMess:  %d\n", len(uniqueVmess))
    fmt.Printf("VLESS:  %d\n", len(uniqueVless))
    fmt.Printf("Trojan: %d\n", len(uniqueTrojan))
    fmt.Printf("SS:     %d\n", len(uniqueSS))

    total := len(uniqueVmess) + len(uniqueVless) + len(uniqueTrojan) + len(uniqueSS)
    fmt.Printf("Total Unique:  %d\n", total)
}
