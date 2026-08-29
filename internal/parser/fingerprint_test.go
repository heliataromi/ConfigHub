package parser

import (
	"encoding/base64"
	"testing"
)

func TestGetFingerprint_RealityAliases(t *testing.T) {
	link1 := "vless://uuid1@1.2.3.4:443?security=reality&pbk=KEY123&sid=SID123&type=tcp#Name1"
	link2 := "vless://UUID1@1.2.3.4:443?security=reality&public_key=KEY123&short_id=SID123&type=tcp#Name2"
	link3 := "vless://uuid1@1.2.3.4:443?security=reality&publicKey=KEY123&shortId=SID123&type=tcp#Name3"

	fp1 := GetFingerprint(link1)
	fp2 := GetFingerprint(link2)
	fp3 := GetFingerprint(link3)

	if fp1 != fp2 {
		t.Fatalf("Fingerprints do not match between pbk and public_key:\n  fp1: %s\n  fp2: %s", fp1, fp2)
	}
	if fp1 != fp3 {
		t.Fatalf("Fingerprints do not match between pbk and publicKey:\n  fp1: %s\n  fp3: %s", fp1, fp3)
	}
}

func TestGetFingerprint_ShadowsocksLegacyAndSIP002(t *testing.T) {
	// SIP002 format: ss://base64(aes-256-gcm:pass123)@1.2.3.4:8388#Remark1
	userBase64 := base64.StdEncoding.EncodeToString([]byte("aes-256-gcm:pass123"))
	sip002 := "ss://" + userBase64 + "@1.2.3.4:8388#Remark1"

	// Legacy format: ss://base64(aes-256-gcm:pass123@1.2.3.4:8388)#Remark2
	legacyBase64 := base64.StdEncoding.EncodeToString([]byte("aes-256-gcm:pass123@1.2.3.4:8388"))
	legacy := "ss://" + legacyBase64 + "#Remark2"

	fp1 := GetFingerprint(sip002)
	fp2 := GetFingerprint(legacy)

	if fp1 != fp2 {
		t.Fatalf("Fingerprints do not match between SIP002 and Legacy SS:\n  fp1: %s\n  fp2: %s", fp1, fp2)
	}
}

func TestGetFingerprint_VMessDefaults(t *testing.T) {
	// VMess 1 with empty net and empty tls
	json1 := `{"v":"2","ps":"Node1","add":"1.2.3.4","port":"443","id":"uuid-123","aid":"0","net":"","type":"none","host":"","path":"","tls":""}`
	link1 := "vmess://" + base64.StdEncoding.EncodeToString([]byte(json1))

	// VMess 2 with explicit net=tcp and tls=none
	json2 := `{"v":"2","ps":"Node2","add":"1.2.3.4","port":443,"id":"UUID-123","aid":"0","net":"tcp","type":"none","host":"","path":"","tls":"none"}`
	link2 := "vmess://" + base64.StdEncoding.EncodeToString([]byte(json2))

	fp1 := GetFingerprint(link1)
	fp2 := GetFingerprint(link2)

	if fp1 != fp2 {
		t.Fatalf("VMess fingerprints do not match defaults:\n  fp1: %s\n  fp2: %s", fp1, fp2)
	}
}

func TestGetFingerprint_DefaultPorts(t *testing.T) {
	linkExplicit := "trojan://pass1@example.com:443?security=tls&sni=example.com#Name1"
	linkImplicit := "trojan://pass1@example.com?security=tls&sni=example.com#Name2"

	fp1 := GetFingerprint(linkExplicit)
	fp2 := GetFingerprint(linkImplicit)

	if fp1 != fp2 {
		t.Fatalf("Default port 443 normalization failed:\n  fp1: %s\n  fp2: %s", fp1, fp2)
	}
}

func TestGetFingerprint_TelegramProxy(t *testing.T) {
	tgLink := "tg://proxy?server=1.2.3.4&port=443&secret=EE1603010200010001fc030386e24c3add6d656469612e737465616d706f77657265642e636f6d#Remark1"
	httpsLink := "https://t.me/proxy?server=1.2.3.4.&port=443&secret=ee1603010200010001fc030386e24c3add6d656469612e737465616d706f77657265642e636f6d"

	fp1 := GetFingerprint(tgLink)
	fp2 := GetFingerprint(httpsLink)

	if fp1 != fp2 {
		t.Fatalf("Telegram proxy fingerprint normalization failed:\n  fp1: %s\n  fp2: %s", fp1, fp2)
	}

	socksTg := "tg://socks?server=1.2.3.4&port=1080&user=usr&pass=pwd#Remark1"
	socksHTTPS := "https://t.me/socks?server=1.2.3.4.&port=1080&user=usr&pass=pwd"
	fpSocks1 := GetFingerprint(socksTg)
	fpSocks2 := GetFingerprint(socksHTTPS)

	if fpSocks1 != fpSocks2 {
		t.Fatalf("Telegram socks fingerprint normalization failed:\n  fp1: %s\n  fp2: %s", fpSocks1, fpSocks2)
	}

	httpTg := "tg://http?server=1.2.3.4&port=8080&user=usr&pass=pwd#Remark1"
	httpHTTPS := "https://t.me/http?server=1.2.3.4.&port=8080&user=usr&pass=pwd"
	fpHTTP1 := GetFingerprint(httpTg)
	fpHTTP2 := GetFingerprint(httpHTTPS)

	if fpHTTP1 != fpHTTP2 {
		t.Fatalf("Telegram http fingerprint normalization failed:\n  fp1: %s\n  fp2: %s", fpHTTP1, fpHTTP2)
	}
}

func TestGetFingerprint_WhiteDNS(t *testing.T) {
	cotten1 := "cottendns://eyJzY2hlbWEiOiJ3aGl0ZWRucy5wcm9maWxlIiwidmVyc2lvbiI6MSwicHJvZmlsZSI6eyJuYW1lIjoiTmFtZTEiLCJzZXJ2ZXIiOnsiZG9tYWluIjoidXNhLmFsaWFzLmNvbSIsImVuY3J5cHRpb25fa2V5IjoiQUJDREVGIiwiZW5jcnlwdGlvbl9tZXRob2QiOjF9fX0="
	cotten2 := "cottendns://eyJzY2hlbWEiOiJ3aGl0ZWRucy5wcm9maWxlIiwidmVyc2lvbiI6MSwicHJvZmlsZSI6eyJuYW1lIjoiTmFtZTIiLCJzZXJ2ZXIiOnsiZG9tYWluIjoiVVNBLkFMSUFTLkNPTSIsImVuY3J5cHRpb25fa2V5IjoiYWJjZGVmIiwiZW5jcnlwdGlvbl9tZXRob2QiOjF9fX0="

	fp1 := GetFingerprint(cotten1)
	fp2 := GetFingerprint(cotten2)

	if fp1 != fp2 {
		t.Fatalf("WhiteDNS fingerprint normalization failed:\n  fp1: %s\n  fp2: %s", fp1, fp2)
	}
}

func TestGetFingerprint_SSR(t *testing.T) {
	// ssr://1.2.3.4:8388:origin:aes-256-cfb:plain:pass/?obfsparam=&remarks=name1
	payload1 := base64.RawURLEncoding.EncodeToString([]byte("1.2.3.4:8388:origin:aes-256-cfb:plain:pass/?obfsparam=&remarks=name1"))
	ssr1 := "ssr://" + payload1

	// same node with different remark name2
	payload2 := base64.RawURLEncoding.EncodeToString([]byte("1.2.3.4:8388:origin:aes-256-cfb:plain:pass/?obfsparam=&remarks=name2"))
	ssr2 := "ssr://" + payload2

	fp1 := GetFingerprint(ssr1)
	fp2 := GetFingerprint(ssr2)

	if fp1 != fp2 {
		t.Fatalf("SSR fingerprint normalization failed (remarks should be ignored):\n  fp1: %s\n  fp2: %s", fp1, fp2)
	}
}

