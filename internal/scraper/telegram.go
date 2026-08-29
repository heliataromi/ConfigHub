package scraper

import (
    "fmt"
    "html"
    "net/http"
    "regexp"
    "strings"
    "time"

    "ConfigHub/internal/extractor"
    "ConfigHub/internal/telemetry"

    "github.com/PuerkitoBio/goquery"
)

var htmlTagRegex = regexp.MustCompile(`<[^>]+>`)

// ScrapeChannel fetches a Telegram web preview and extracts configs from the last 24 hours
func ScrapeChannel(channel string) (extractor.Configs, error) {
    startTime := time.Now()
    var finalConfigs extractor.Configs
    var finalErr error
    var statusCode int
    var msgCount int

    defer func() {
        errStr := ""
        if finalErr != nil {
            errStr = finalErr.Error()
        }
        telemetry.Global.RecordSource(telemetry.ChannelStat{
            Name:          channel,
            Type:          "channel",
            StatusCode:    statusCode,
            Duration:      time.Since(startTime),
            MessagesCount: msgCount,
            ConfigsYield:  finalConfigs.Count(),
            Error:         errStr,
        })
    }()

    url := fmt.Sprintf("https://t.me/s/%s", channel)
    if strings.HasPrefix(channel, "http://") || strings.HasPrefix(channel, "https://") {
        url = channel
    }

    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        finalErr = err
        return extractor.Configs{}, err
    }

    req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")

    client := &http.Client{Timeout: 10 * time.Second}
    resp, err := client.Do(req)
    if err != nil {
        finalErr = err
        return extractor.Configs{}, err
    }
    defer resp.Body.Close()

    statusCode = resp.StatusCode
    if resp.StatusCode != 200 {
        finalErr = fmt.Errorf("received status code %d", resp.StatusCode)
        return extractor.Configs{}, finalErr
    }

    doc, err := goquery.NewDocumentFromReader(resp.Body)
    if err != nil {
        finalErr = err
        return extractor.Configs{}, err
    }

    var rawTexts []string
    now := time.Now()

    doc.Find(".tgme_widget_message").Each(func(i int, s *goquery.Selection) {
        timeTag := s.Find(".tgme_widget_message_date time")
        datetimeAttr, exists := timeTag.Attr("datetime")
        if !exists {
            return
        }

        msgTime, err := time.Parse(time.RFC3339, datetimeAttr)
        if err != nil {
            return
        }

        if now.Sub(msgTime) <= 24*time.Hour {
            msgCount++

            // 1. Extract all hyperlink targets (inline links in text, captions, and glass buttons)
            s.Find("a[href]").Each(func(_ int, a *goquery.Selection) {
                if href, exists := a.Attr("href"); exists {
                    href = strings.TrimSpace(href)
                    href = html.UnescapeString(href)
                    if strings.Contains(href, "&amp;") {
                        href = html.UnescapeString(href)
                    }
                    if href != "" {
                        rawTexts = append(rawTexts, href)
                    }
                }
            })

            // 2. Extract visible text and captions
            extractHTML := func(elem *goquery.Selection) {
                if elem.Length() == 0 {
                    return
                }
                htmlContent, err := elem.Html()
                if err != nil || htmlContent == "" {
                    return
                }

                textContent := strings.ReplaceAll(htmlContent, "<br/>", "\n")
                textContent = strings.ReplaceAll(textContent, "<br>", "\n")
                textContent = strings.ReplaceAll(textContent, "</p>", "\n")
                textContent = strings.ReplaceAll(textContent, "</div>", "\n")
                textContent = strings.ReplaceAll(textContent, "</blockquote>", "\n")
                textContent = htmlTagRegex.ReplaceAllString(textContent, " ")
                textContent = html.UnescapeString(textContent)

                if strings.TrimSpace(textContent) != "" {
                    rawTexts = append(rawTexts, textContent)
                }
            }

            extractHTML(s.Find(".tgme_widget_message_text"))
            extractHTML(s.Find(".tgme_widget_message_caption"))
        }
    })

    fullText := strings.Join(rawTexts, "\n")
    finalConfigs = extractor.AuditAndExtract(fullText, "t.me/"+channel, telemetry.Global.RecordDropped)
    return finalConfigs, nil
}
