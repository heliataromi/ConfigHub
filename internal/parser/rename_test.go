package parser

import (
	"strings"
	"testing"
)

func TestRenameSS_SIP002AndLegacy(t *testing.T) {
	// Test SIP002 format
	sip002 := "ss://YWVzLTI1Ni1nY206cGFzc3dvcmQ=@1.2.3.4:8388#OldName"
	renamedSIP002 := RenameConfig(sip002, "t.me/test", nil)
	if !strings.HasPrefix(renamedSIP002, "ss://YWVzLTI1Ni1nY206cGFzc3dvcmQ=@1.2.3.4:8388#") {
		t.Errorf("Unexpected SIP002 prefix: %s", renamedSIP002)
	}

	// Test Legacy format
	legacy := "ss://YWVzLTI1Ni1nY206cGFzc3dvcmRAMS4yLjMuNDo4Mzg4#OldName"
	renamedLegacy := RenameConfig(legacy, "t.me/test", nil)
	if !strings.HasPrefix(renamedLegacy, "ss://YWVzLTI1Ni1nY206cGFzc3dvcmRAMS4yLjMuNDo4Mzg4#") {
		t.Errorf("Unexpected Legacy prefix: %s", renamedLegacy)
	}
}

func TestRenameSSR(t *testing.T) {
	// ssr://server:port:proto:method:obfs:pass/?obfsparam=&remarks=old
	ssrLink := "ssr://MS4yLjMuNDo4Mzg4Om9yaWdpbjphZXMtMjU2LWNmbjpwbGFpbjpkR1Z6ZEEvP29iZnNwYXJhbT0mcmVtYXJrcz1kR1Z6ZEE="
	renamedSSR := RenameConfig(ssrLink, "t.me/test", nil)
	if !strings.HasPrefix(renamedSSR, "ssr://") {
		t.Errorf("Expected ssr:// prefix, got %s", renamedSSR)
	}
	if renamedSSR == ssrLink {
		t.Errorf("Expected renamed SSR link to be modified with new remarks")
	}
}

func TestRenameTelegram(t *testing.T) {
	link := "tg://proxy?server=1.2.3.4&port=443&secret=ee1603010200010001fc030386e24c3add6d656469612e737465616d706f77657265642e636f6d"
	renamed := RenameConfig(link, "t.me/test", nil)
	if !strings.HasPrefix(renamed, "tg://proxy?") || !strings.Contains(renamed, "#") {
		t.Errorf("Expected renamed Telegram proxy link to contain # with remarks, got: %s", renamed)
	}
}

func TestRenameWhiteDNS(t *testing.T) {
	link := "cottendns://eyJzY2hlbWEiOiJ3aGl0ZWRucy5wcm9maWxlIiwidmVyc2lvbiI6MSwicHJvZmlsZSI6eyJuYW1lIjoidGVzdCIsInNlcnZlciI6eyJkb21haW4iOiJleGFtcGxlLmNvbSIsImVuY3J5cHRpb25fa2V5IjoiMTIzNCIsImVuY3J5cHRpb25fbWV0aG9kIjoxfX19"
	renamed := RenameConfig(link, "t.me/test", nil)
	if !strings.HasPrefix(renamed, "cottendns://") {
		t.Errorf("Expected cottendns:// prefix, got: %s", renamed)
	}
	if renamed == link {
		t.Errorf("Expected renamed link to be modified with new name inside payload")
	}
}

func TestRenameStandardProtocols(t *testing.T) {
	links := []struct {
		proto string
		link  string
	}{
		{"anytls", "anytls://7aecb4a5-f0a4-32a0-aabe-c9d5241e313f@1.2.3.4:443#Old"},
		{"snell", "snell://psk123@1.2.3.4:443?version=4#Old"},
		{"http", "http://user:pass@1.2.3.4:8080#Old"},
		{"juicity", "juicity://7aecb4a5-f0a4-32a0-aabe-c9d5241e313f:pass1@1.2.3.4:443#Old"},
		{"naive", "naive+https://user:pass@1.2.3.4:443#Old"},
	}

	for _, tc := range links {
		renamed := RenameConfig(tc.link, "t.me/test", nil)
		if !strings.HasPrefix(renamed, tc.proto) {
			t.Errorf("Expected prefix %s, got: %s", tc.proto, renamed)
		}
		if !strings.Contains(renamed, "t.me/test") {
			t.Errorf("Expected channel remark in %s: %s", tc.proto, renamed)
		}
	}
}
