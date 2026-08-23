package scraper

import (
	"strings"
	"testing"

	"ConfigHub/internal/extractor"
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
