package scraper

import (
	"fmt"
	"html"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"ConfigHub/internal/extractor"

	"github.com/PuerkitoBio/goquery"
)

func TestHTMLTagStripping(t *testing.T) {
	htmlSample := `
<blockquote><code>vless://7aecb4a5-f0a4-32a0-aabe-c9d5241e313f@191.40.32.96:10018?security=reality&type=tcp#AppSoonerNode</code></blockquote>
<div><pre>hy2://pass1@1.2.3.4:443?sni=example.com#Hy2Node</pre></div>
`
	text := strings.ReplaceAll(htmlSample, "<br/>", "\n")
	text = strings.ReplaceAll(text, "<br>", "\n")
	text = strings.ReplaceAll(text, "</p>", "\n")
	text = strings.ReplaceAll(text, "</div>", "\n")
	text = strings.ReplaceAll(text, "</blockquote>", "\n")
	text = htmlTagRegex.ReplaceAllString(text, " ")

	configs := extractor.Extract(text)

	if len(configs.Vless) != 1 {
		t.Fatalf("Expected 1 VLESS config from code block, got %d", len(configs.Vless))
	}
	if configs.Vless[0] != "vless://7aecb4a5-f0a4-32a0-aabe-c9d5241e313f@191.40.32.96:10018?security=reality&type=tcp#AppSoonerNode" {
		t.Errorf("Unexpected VLESS link: %s", configs.Vless[0])
	}
	if len(configs.Hy2) != 1 {
		t.Fatalf("Expected 1 Hy2 config from pre block, got %d", len(configs.Hy2))
	}
}

func TestExtractInlineHyperlinksAndButtons(t *testing.T) {
	sampleHTML := `
<div class="tgme_widget_message js-widget_message" data-post="iRoProxy/57532">
  <div class="tgme_widget_message_text js-message_text">
    چهار کار دنیا به زور نمی شود:<br/>
    <a href="https://t.me/proxy?server=silnet.varfootball.co.uk&amp;amp;port=2053&amp;amp;secret=eeNEgYdJvXrFGRMCIMJdCQ">پروکسی ۱</a> |
    <a href="https://t.me/proxy?server=iro.varfootball2.co.uk&amp;port=2053&amp;secret=eeNEgYdJvXrFGRMCIMJdCQ">پروکسی ۲</a> |
    <a href="https://t.me/socks?server=proxy.socks.example.com&amp;port=1080&amp;user=u1&amp;pass=p1">ساکس تلگرام</a>
  </div>
  <div class="tgme_widget_message_inline_keyboard">
    <a class="tgme_widget_message_inline_button url_button" href="tg://proxy?server=Suodbo.co.uk&amp;port=443&amp;secret=eef0eeb0bd9adc4fd4a93994ee3b2a216b63646e2e79656b74616e65742e636f6d">
      <span class="tgme_widget_message_inline_button_text">Connect</span>
    </a>
  </div>
</div>
`
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(sampleHTML))
	if err != nil {
		t.Fatalf("Failed to parse HTML: %v", err)
	}

	var rawTexts []string
	doc.Find(".tgme_widget_message").Each(func(i int, s *goquery.Selection) {
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
	})

	fullText := strings.Join(rawTexts, "\n")
	configs := extractor.Extract(fullText)

	if len(configs.Telegram) != 4 {
		t.Fatalf("Expected 4 extracted Telegram proxies, got %d: %v", len(configs.Telegram), configs.Telegram)
	}
}

func TestScrapeChannel_MockedServer(t *testing.T) {
	recentTime := time.Now().UTC().Format(time.RFC3339)
	oldTime := time.Now().Add(-48 * time.Hour).UTC().Format(time.RFC3339)

	sampleHTML := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<body>
  <!-- Recent message with text and photo caption -->
  <div class="tgme_widget_message" data-post="test/1">
    <div class="tgme_widget_message_date"><time datetime="%s">Recent</time></div>
    <div class="tgme_widget_message_text">
      vless://7aecb4a5-f0a4-32a0-aabe-c9d5241e313f@191.40.32.96:10018?security=reality&type=tcp#Node1
    </div>
    <div class="tgme_widget_message_caption">
      hy2://pass1@1.2.3.4:443?sni=example.com#PhotoCaptionNode
    </div>
    <div class="tgme_widget_message_inline_keyboard">
      <a class="tgme_widget_message_inline_button" href="tg://proxy?server=Suodbo.co.uk&amp;port=443&amp;secret=eef0eeb0bd9adc4fd4a93994ee3b2a216b63646e2e79656b74616e65742e636f6d">Connect</a>
    </div>
  </div>

  <!-- Old message (> 24h) should be skipped -->
  <div class="tgme_widget_message" data-post="test/2">
    <div class="tgme_widget_message_date"><time datetime="%s">Old</time></div>
    <div class="tgme_widget_message_text">
      trojan://oldpass@1.2.3.4:443#OldNode
    </div>
  </div>
</body>
</html>
`, recentTime, oldTime)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(sampleHTML))
	}))
	defer server.Close()

	configs, err := ScrapeChannel(server.URL)
	if err != nil {
		t.Fatalf("ScrapeChannel failed: %v", err)
	}

	if len(configs.Vless) != 1 {
		t.Errorf("Expected 1 VLESS config, got %d", len(configs.Vless))
	}
	if len(configs.Hy2) != 1 {
		t.Errorf("Expected 1 Hy2 config, got %d", len(configs.Hy2))
	}
	if len(configs.Telegram) != 1 {
		t.Errorf("Expected 1 Telegram proxy config, got %d", len(configs.Telegram))
	}
	if len(configs.Trojan) != 0 {
		t.Errorf("Expected 0 Trojan configs (older than 24h), got %d", len(configs.Trojan))
	}
}

func TestScrapeChannel_HTTPError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	_, err := ScrapeChannel(server.URL)
	if err == nil {
		t.Errorf("Expected error on HTTP 404, got nil")
	}
}

