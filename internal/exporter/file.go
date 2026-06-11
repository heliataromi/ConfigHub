package exporter

import (
    "encoding/base64"
    "fmt"
    "os"
    "path/filepath"
    "strings"
    "time"

    ptime "github.com/yaa110/go-persian-calendar"
)

// WriteSubFiles generates standard and Base64 encoded subscription files from slices
func WriteSubFiles(outDir string, vmess, vless, trojan, ss []string) error {
    // 1. Ensure the output directory exists
    err := os.MkdirAll(outDir, 0755)
    if err != nil {
        return fmt.Errorf("failed to create output directory: %w", err)
    }

    // Combine all for the "mixed" file
    var mixed []string
    mixed = append(mixed, vmess...)
    mixed = append(mixed, vless...)
    mixed = append(mixed, trojan...)
    mixed = append(mixed, ss...)

    // 2. Define the files we want to create
    filesToCreate := map[string][]string{
        "vmess.txt":  vmess,
        "vless.txt":  vless,
        "trojan.txt": trojan,
        "ss.txt":     ss,
        "mixed.txt":  mixed,
    }

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

    return nil
}
