package exporter

import (
    "encoding/base64"
    "fmt"
    "os"
    "path/filepath"
    "strings"
    "time"

    "ConfigHub/internal/clash"

    ptime "github.com/yaa110/go-persian-calendar"
)

// RenamedConfig holds the final config string and its original source type
type RenamedConfig struct {
    URL    string
    Source string
}

// WriteSubFiles generates standard and Base64 encoded subscription files from a map of configs
func WriteSubFiles(outDir string, configsMap map[string][]RenamedConfig) error {
    // 1. Ensure the output directory exists
    err := os.MkdirAll(outDir, 0755)
    if err != nil {
        return fmt.Errorf("failed to create output directory: %w", err)
    }

    // Combine all for the "mixed" file
    var mixed []string
    var mixedLite []string
    filesToCreate := make(map[string][]string)
    
    for protocolName, configs := range configsMap {
        var protoConfigs []string
        for _, c := range configs {
            protoConfigs = append(protoConfigs, c.URL)
            mixed = append(mixed, c.URL)
            if c.Source == "channel" {
                mixedLite = append(mixedLite, c.URL)
            }
        }
        filesToCreate[protocolName+".txt"] = protoConfigs
    }
    filesToCreate["mixed.txt"] = mixed
    filesToCreate["mixed_lite.txt"] = mixedLite

    // Generate dummy config with Jalali timestamp
    loc := time.FixedZone("IRST", int(3.5*3600))
    pt := ptime.New(time.Now().In(loc))
    updateTimeStr := fmt.Sprintf("🐣🏁 Last update: %s", pt.Format("📅yyyy/MM/dd 🕒HH:mm"))
    dummyConfig := fmt.Sprintf("vless://00000000-0000-0000-0000-000000000000@127.0.0.1:0?type=tcp#%s", strings.ReplaceAll(updateTimeStr, " ", "%20"))

    // 3. Write the files
    for filename, configs := range filesToCreate {
        if len(configs) == 0 {
            continue // Skip empty files
        }

        // Prepend dummy config at the top
        finalConfigs := append([]string{dummyConfig}, configs...)
        rawContent := strings.Join(finalConfigs, "\n")

        // Write Raw File
        rawPath := filepath.Join(outDir, filename)
        err := os.WriteFile(rawPath, []byte(rawContent), 0644)
        if err != nil {
            return fmt.Errorf("failed to write %s: %w", filename, err)
        }

        // Write Base64 File
        base64Name := strings.Replace(filename, ".txt", "_base64.txt", 1)
        base64Path := filepath.Join(outDir, base64Name)

        encodedContent := base64.StdEncoding.EncodeToString([]byte(rawContent))
        err = os.WriteFile(base64Path, []byte(encodedContent), 0644)
        if err != nil {
            return fmt.Errorf("failed to write %s: %w", base64Name, err)
        }
    }

    // 4. Generate Clash Subscriptions (Full & Lite)
    if len(mixed) > 0 {
        clashMixedYAML, err := clash.GenerateClashConfig(mixed)
        if err == nil {
            clashPath := filepath.Join(outDir, "clash.yaml")
            if err := os.WriteFile(clashPath, clashMixedYAML, 0644); err != nil {
                return fmt.Errorf("failed to write clash.yaml: %w", err)
            }
        }
    }

    if len(mixedLite) > 0 {
        clashLiteYAML, err := clash.GenerateClashConfig(mixedLite)
        if err == nil {
            clashLitePath := filepath.Join(outDir, "clash_lite.yaml")
            if err := os.WriteFile(clashLitePath, clashLiteYAML, 0644); err != nil {
                return fmt.Errorf("failed to write clash_lite.yaml: %w", err)
            }
        }
    }

    return nil
}
