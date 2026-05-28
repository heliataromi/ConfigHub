package scraper

import (
    "fmt"
    "html"
    "net/http"
    "strings"
    "time"

    "ConfigHub/internal/extractor"

    "github.com/PuerkitoBio/goquery"
)

// ScrapeChannel fetches a Telegram web preview and extracts configs from the last 24 hours
func ScrapeChannel(channel string) (extractor.Configs, error) {
    url := fmt.Sprintf("https://t.me/s/%s", channel)

    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        return extractor.Configs{}, err
    }

    // Prevent Telegram from blocking the request by mimicking a browser
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

    doc, err := goquery.NewDocumentFromReader(resp.Body)
    if err != nil {
        return extractor.Configs{}, err
    }

    var rawTexts []string
    now := time.Now()

    // Find all message containers on the page
    doc.Find(".tgme_widget_message").Each(func(i int, s *goquery.Selection) {
        // Extract the timestamp from the <time datetime="..."> tag
        timeTag := s.Find(".tgme_widget_message_date time")
        datetimeAttr, exists := timeTag.Attr("datetime")
        if !exists {
            return // Skip if no time is found
        }

        // Telegram formats time in RFC3339 (e.g., 2026-05-27T10:00:00+00:00)
        msgTime, err := time.Parse(time.RFC3339, datetimeAttr)
        if err != nil {
            return
        }

        // Check if the message is within the last 24 hours
        if now.Sub(msgTime) <= 24*time.Hour {
            // Get the message text HTML
            htmlContent, _ := s.Find(".tgme_widget_message_text").Html()

            // Replace HTML line breaks with actual newlines
            textContent := strings.ReplaceAll(htmlContent, "<br/>", "\n")
            textContent = strings.ReplaceAll(textContent, "<br>", "\n")

            // Unescape HTML entities (converts &amp; to &, &quot; to ", etc.)
            textContent = html.UnescapeString(textContent)

            rawTexts = append(rawTexts, textContent)
        }
    })

    // Combine all valid message texts and extract the configs
    fullText := strings.Join(rawTexts, "\n")
    return extractor.Extract(fullText), nil
}
