package exporter

import (
    "encoding/base64"
    "fmt"
    "os"
    "path/filepath"
    "strings"
)

// mapKeysToSlice is a helper to extract the keys from our deduplication map
func mapKeysToSlice(m map[string]bool) []string {
    keys := make([]string, 0, len(m))
    for k := range m {
        keys = append(keys, k)
    }
    return keys
}

// WriteSubFiles generates standard and Base64 encoded subscription files
func WriteSubFiles(outDir string, vmess, vless, trojan, ss map[string]bool) error {
    // 1. Ensure the output directory exists
    err := os.MkdirAll(outDir, 0755)
    if err != nil {
        return fmt.Errorf("failed to create output directory: %w", err)
    }

    // 2. Convert maps to slices
    vmessSlice := mapKeysToSlice(vmess)
    vlessSlice := mapKeysToSlice(vless)
    trojanSlice := mapKeysToSlice(trojan)
    ssSlice := mapKeysToSlice(ss)

    // Combine all for the "mixed" file
    var mixedSlice []string
    mixedSlice = append(mixedSlice, vmessSlice...)
    mixedSlice = append(mixedSlice, vlessSlice...)
    mixedSlice = append(mixedSlice, trojanSlice...)
    mixedSlice = append(mixedSlice, ssSlice...)

    // 3. Define the files we want to create
    filesToCreate := map[string][]string{
        "vmess.txt":  vmessSlice,
        "vless.txt":  vlessSlice,
        "trojan.txt": trojanSlice,
        "ss.txt":     ssSlice,
        "mixed.txt":  mixedSlice,
    }

    // 4. Write the files
    for filename, configs := range filesToCreate {
        // Skip writing if there are no configs for this protocol
        if len(configs) == 0 {
            continue
        }

        // Join all configs with a newline
        rawContent := strings.Join(configs, "\n")

        // Write Raw File
        rawPath := filepath.Join(outDir, filename)
        err := os.WriteFile(rawPath, []byte(rawContent), 0644)
        if err != nil {
            return fmt.Errorf("failed to write %s: %w", filename, err)
        }

        // Write Base64 Encoded File
        // We append "_base64" to the filename (e.g., vmess_base64.txt)
        base64Name := strings.Replace(filename, ".txt", "_base64.txt", 1)
        base64Path := filepath.Join(outDir, base64Name)

        // Encode the raw text to Base64
        encodedContent := base64.StdEncoding.EncodeToString([]byte(rawContent))
        err = os.WriteFile(base64Path, []byte(encodedContent), 0644)
        if err != nil {
            return fmt.Errorf("failed to write %s: %w", base64Name, err)
        }
    }

    return nil
}
