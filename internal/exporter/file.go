package exporter

import (
    "encoding/base64"
    "fmt"
    "net/url"
    "os"
    "path/filepath"
    "strings"
    "time"

    "ConfigHub/internal/clash"
    "ConfigHub/internal/singbox"

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
    
    var whitednsConfigs []string
    for protocolName, configs := range configsMap {
        var protoConfigs []string
        for _, c := range configs {
            protoConfigs = append(protoConfigs, c.URL)
            if protocolName != "telegram" && protocolName != "cottendns" && protocolName != "stormdns" {
                mixed = append(mixed, c.URL)
                if c.Source == "channel" {
                    mixedLite = append(mixedLite, c.URL)
                }
            }
            if protocolName == "cottendns" || protocolName == "stormdns" {
                whitednsConfigs = append(whitednsConfigs, c.URL)
            }
        }
        filesToCreate[protocolName+".txt"] = protoConfigs
        if protocolName == "telegram" {
            filesToCreate["mtproto.txt"] = protoConfigs
        }
    }
    if len(whitednsConfigs) > 0 {
        filesToCreate["whitedns.txt"] = whitednsConfigs
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

        var finalConfigs []string
        if filename == "telegram.txt" || filename == "mtproto.txt" || filename == "cottendns.txt" || filename == "stormdns.txt" || filename == "whitedns.txt" {
            finalConfigs = configs
        } else {
            // Prepend dummy config at the top for standard client subscriptions
            finalConfigs = append([]string{dummyConfig}, configs...)
        }
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

    // 4. Generate Per-Country Subscriptions in sub/countries/
    countriesDir := filepath.Join(outDir, "countries")
    if err := os.MkdirAll(countriesDir, 0755); err != nil {
        return fmt.Errorf("failed to create countries directory: %w", err)
    }

    countryMap := make(map[string][]string)
    for _, link := range mixed {
        code := extractCountryCode(link)
        if code != "" {
            countryMap[code] = append(countryMap[code], link)
        }
    }

    for code, configs := range countryMap {
        if len(configs) == 0 {
            continue
        }

        finalConfigs := append([]string{dummyConfig}, configs...)
        rawContent := strings.Join(finalConfigs, "\n")

        // Raw country file: sub/countries/de.txt
        rawPath := filepath.Join(countriesDir, code+".txt")
        if err := os.WriteFile(rawPath, []byte(rawContent), 0644); err != nil {
            return fmt.Errorf("failed to write country file %s: %w", rawPath, err)
        }

        // Base64 country file: sub/countries/de_base64.txt
        base64Path := filepath.Join(countriesDir, code+"_base64.txt")
        encodedContent := base64.StdEncoding.EncodeToString([]byte(rawContent))
        if err := os.WriteFile(base64Path, []byte(encodedContent), 0644); err != nil {
            return fmt.Errorf("failed to write country base64 file %s: %w", base64Path, err)
        }
    }

    // 5. Generate Clash Subscriptions (Full & Lite)
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

    // 6. Generate Sing-box Subscriptions (Full & Lite)
    if len(mixed) > 0 {
        singboxMixedJSON, err := singbox.GenerateSingboxConfig(mixed)
        if err == nil {
            singboxPath := filepath.Join(outDir, "singbox.json")
            if err := os.WriteFile(singboxPath, singboxMixedJSON, 0644); err != nil {
                return fmt.Errorf("failed to write singbox.json: %w", err)
            }
        }
    }

    if len(mixedLite) > 0 {
        singboxLiteJSON, err := singbox.GenerateSingboxConfig(mixedLite)
        if err == nil {
            singboxLitePath := filepath.Join(outDir, "singbox_lite.json")
            if err := os.WriteFile(singboxLitePath, singboxLiteJSON, 0644); err != nil {
                return fmt.Errorf("failed to write singbox_lite.json: %w", err)
            }
        }
    }

    return nil
}

// extractCountryCode extracts the lowercase 2-letter ISO or category code from a config remark
func extractCountryCode(link string) string {
    parts := strings.Split(link, "#")
    if len(parts) < 2 {
        return ""
    }
    remark := parts[len(parts)-1]
    if unescaped, err := url.QueryUnescape(remark); err == nil && unescaped != "" {
        remark = unescaped
    }

    fields := strings.Fields(remark)
    if len(fields) >= 2 {
        tag := strings.ToUpper(strings.TrimSpace(fields[1]))
        if len(tag) == 2 && tag != "UN" {
            return strings.ToLower(tag)
        }
        if tag == "CDN/RELAY" || tag == "CDN" {
            return "cdn"
        }
        if tag == "IR-RELAY" || tag == "IRAN-RELAY" {
            return "ir_relay"
        }
    }
    return ""
}
