package scraper

import (
    "encoding/base64"
    "fmt"
    "io"
    "net/http"
    "strings"
    "time"

    "ConfigHub/internal/extractor"
)

// ScrapeSubscription fetches a subscription link, attempts to base64 decode it,
// and extracts the configs.
func ScrapeSubscription(url string) (extractor.Configs, error) {
    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        return extractor.Configs{}, err
    }

    req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")

    client := &http.Client{Timeout: 10 * time.Second}
    resp, err := client.Do(req)
    if err != nil {
        return extractor.Configs{}, err
    }
    defer resp.Body.Close()

    if resp.StatusCode != 200 {
        return extractor.Configs{}, fmt.Errorf("received status code %d", resp.StatusCode)
    }

    bodyBytes, err := io.ReadAll(resp.Body)
    if err != nil {
        return extractor.Configs{}, err
    }

    bodyText := string(bodyBytes)
    
    decodedText := decodeBase64(bodyText)
    if decodedText != "" {
        configs := extractor.Extract(decodedText)
        if hasConfigs(configs) {
            return configs, nil
        }
    }

    // Fallback to plain text
    return extractor.Extract(bodyText), nil
}

func decodeBase64(s string) string {
    s = strings.TrimSpace(s)
    
    // Clean up typical whitespace issues in sub links
    s = strings.ReplaceAll(s, "\n", "")
    s = strings.ReplaceAll(s, "\r", "")
    s = strings.ReplaceAll(s, " ", "")
    
    // Try Standard Encoding
    decoded, err := base64.StdEncoding.DecodeString(s)
    if err == nil {
        return string(decoded)
    }
    
    // Try Raw Encoding (no padding)
    decoded, err = base64.RawStdEncoding.DecodeString(s)
    if err == nil {
        return string(decoded)
    }
    
    // Try URL Encoding
    decoded, err = base64.URLEncoding.DecodeString(s)
    if err == nil {
        return string(decoded)
    }
    
    // Try Raw URL Encoding
    decoded, err = base64.RawURLEncoding.DecodeString(s)
    if err == nil {
        return string(decoded)
    }

    // Attempt to manually fix padding for StdEncoding
    pad := len(s) % 4
    if pad > 0 {
        s += strings.Repeat("=", 4-pad)
        decoded, err = base64.StdEncoding.DecodeString(s)
        if err == nil {
            return string(decoded)
        }
    }

    return ""
}

func hasConfigs(c extractor.Configs) bool {
    return c.Count() > 0
}
