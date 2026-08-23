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
