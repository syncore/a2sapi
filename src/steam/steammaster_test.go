package steam

import (
	"strings"
	"testing"
)

func TestExtractHosts(t *testing.T) {
	data := []byte{
		0xC0, 0xD3, 0x3E, 0x0B, 0x6D, 0x38, 0x2D, 0x37, 0xA8, 0xA0, 0x6D, 0x38,
		0x68, 0xEC, 0x89, 0x14, 0x6D, 0x38, 0xD0, 0x43, 0x01, 0x43, 0x64, 0xC5,
		0x59, 0xC5, 0x30, 0xB6, 0x64, 0xC5, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	hosts, total, err := extractHosts(data)
	if err != nil {
		t.Fatalf("Unexpected error when extracting hosts")
	}
	// last host terminates list as 0.0.0.0:0 so total = total-1
	total = total - 1
	if total != 5 {
		t.Fatalf("Expected extraction of 5 total hosts, got: %d", total)
	}

	hoststrings := []string{
		"192.211.62.11:27960",
		"45.55.168.160:27960",
		"104.236.137.20:27960",
		"208.67.1.67:25797",
		"89.197.48.182:25797",
	}

	found, expected := 0, 5
	for _, h := range hosts {
		for _, hs := range hoststrings {
			if strings.EqualFold(h, hs) {
				found++
			}
		}
	}
	if found != expected {
		t.Fatalf("Expected 5 hosts to have been extracted, only got: %d", found)
	}
}

func TestParseIP(t *testing.T) {
	parsed, err := parseIP([]byte{0x59, 0xC5, 0x30, 0xB6, 0x64, 0xC5})
	if err != nil {
		t.Fatalf("Unexpected error when parsing IP")
	}
	if !strings.EqualFold(parsed, "89.197.48.182:25797") {
		t.Fatalf("Expected IP: 89.197.48.182:25797, got: %s", parsed)
	}
	parsed, err = parseIP([]byte{0x2D, 0x37, 0xA8, 0xA0, 0x6D, 0x38})
	if err != nil {
		t.Fatalf("Unexpected error when parsing IP")
	}
	if !strings.EqualFold(parsed, "45.55.168.160:27960") {
		t.Fatalf("Expected IP: 45.55.168.160:27960, got: %s", parsed)
	}
}
