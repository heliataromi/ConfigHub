package main

import (
    "fmt"
    "log"
    "sync"

    "ConfigHub/internal/config"
    "ConfigHub/internal/exporter"
    "ConfigHub/internal/parser" // Import our new parser
    "ConfigHub/internal/scraper"
)

// helper to convert map values to a slice
func getMapValues(m map[string]string) []string {
    var s []string
    for _, v := range m {
        s = append(s, v)
    }
    return s
}

func main() {
    channels, err := config.ReadChannels("channels.txt")
    if err != nil {
        log.Fatalf("[-] Critical Error: Could not read channels.txt: %v\n", err)
    }

    if len(channels) == 0 {
        log.Fatalf("[-] No channels found. Please add some to channels.txt!")
    }

    fmt.Printf("Loaded %d channels. Starting scraper...\n", len(channels))

    // Changed to map[string]string
    // Key = Semantic Fingerprint, Value = Original Raw Config
    uniqueVmess := make(map[string]string)
    uniqueVless := make(map[string]string)
    uniqueTrojan := make(map[string]string)
    uniqueSS := make(map[string]string)

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
            // Use the parser to generate fingerprints
            for _, c := range configs.Vmess {
                fp := parser.GetFingerprint(c)
                if _, exists := uniqueVmess[fp]; !exists {
                    uniqueVmess[fp] = c
                }
            }
            for _, c := range configs.Vless {
                fp := parser.GetFingerprint(c)
                if _, exists := uniqueVless[fp]; !exists {
                    uniqueVless[fp] = c
                }
            }
            for _, c := range configs.Trojan {
                fp := parser.GetFingerprint(c)
                if _, exists := uniqueTrojan[fp]; !exists {
                    uniqueTrojan[fp] = c
                }
            }
            for _, c := range configs.SS {
                fp := parser.GetFingerprint(c)
                if _, exists := uniqueSS[fp]; !exists {
                    uniqueSS[fp] = c
                }
            }
            mu.Unlock()

            totalFound := len(configs.Vmess) + len(configs.Vless) + len(configs.Trojan) + len(configs.SS)
            fmt.Printf("[+] %s: Found %d configs\n", channel, totalFound)
        }(ch)
    }

    wg.Wait()

    // Extract the pure slices of unique configs
    finalVmess := getMapValues(uniqueVmess)
    finalVless := getMapValues(uniqueVless)
    finalTrojan := getMapValues(uniqueTrojan)
    finalSS := getMapValues(uniqueSS)

    fmt.Println("\n--- Deduplication Complete ---")
    fmt.Printf("VMess:  %d unique\n", len(finalVmess))
    fmt.Printf("VLESS:  %d unique\n", len(finalVless))
    fmt.Printf("Trojan: %d unique\n", len(finalTrojan))
    fmt.Printf("SS:     %d unique\n", len(finalSS))

    // Pass the slices to the exporter
    outputDirectory := "sub"
    err = exporter.WriteSubFiles(outputDirectory, finalVmess, finalVless, finalTrojan, finalSS)
    if err != nil {
        log.Fatalf("[-] Error saving files: %v\n", err)
    }

    fmt.Printf("[+] Successfully saved highly-deduplicated files to '%s/'!\n", outputDirectory)
}
