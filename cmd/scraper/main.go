package main

import (
    "fmt"
    "log"
    "sync"

    "ConfigHub/internal/config"
    "ConfigHub/internal/exporter"
    "ConfigHub/internal/scraper"
)

func main() {
    channels, err := config.ReadChannels("channels.txt")
    if err != nil {
        log.Fatalf("[-] Critical Error: Could not read channels.txt: %v\n", err)
    }

    if len(channels) == 0 {
        log.Fatalf("[-] No channels found. Please add some to channels.txt!")
    }

    fmt.Printf("Loaded %d channels. Starting scraper...\n", len(channels))

    uniqueVmess := make(map[string]bool)
    uniqueVless := make(map[string]bool)
    uniqueTrojan := make(map[string]bool)
    uniqueSS := make(map[string]bool)

    var mu sync.Mutex
    var wg sync.WaitGroup

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

            total := len(configs.Vmess) + len(configs.Vless) + len(configs.Trojan) + len(configs.SS)
            fmt.Printf("[+] %s: Found %d configs\n", channel, total)
        }(ch)
    }

    wg.Wait()

    fmt.Println("\n--- Saving Files ---")

    // Export the configs to the "sub" directory
    outputDirectory := "sub"
    err = exporter.WriteSubFiles(outputDirectory, uniqueVmess, uniqueVless, uniqueTrojan, uniqueSS)
    if err != nil {
        log.Fatalf("[-] Error saving subscription files: %v\n", err)
    }

    fmt.Printf("[+] Successfully saved all files to the '%s/' directory!\n", outputDirectory)
}
