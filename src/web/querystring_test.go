package web

// Tests for query strings

import (
	"strings"
	"testing"
)

// TestGetQStringValues tests the getQStringValues query string extraction
// function.
func TestGetQStringValues(t *testing.T) {
	getsrvid := make(map[string][]string, 1)
	getsrvid["hosts"] = []string{
		"127.0.0.1,172.16.0.1,%2010.0.0.1",
	}
	result := getQStringValues(getsrvid, qsGetServerIDs)
	if len(result) != 3 {
		t.Fatalf("Expected 3 address strings in result, got: %d", len(result))
	}
	if !strings.EqualFold(result[0], "127.0.0.1") {
		t.Fatalf("Expected address to be: %s, got: %s",
			getsrvid["address"][0], result[0])
	}
	if !strings.EqualFold(result[1], "172.16.0.1") {
		t.Fatalf("Expected address to be: %s, got: %s",
			getsrvid["address"][0], result[1])
	}
	if !strings.EqualFold(result[2], "%2010.0.0.1") {
		t.Fatalf("Expected address to be: %s, got: %s",
			getsrvid["address"][0], result[2])
	}

	getsrvcountry := make(map[string][]string, 1)
	getsrvcountry["countries"] = []string{"US"}
	result = getQStringValues(getsrvcountry, qsGetServersCountry)
	if len(result) != 1 {
		t.Fatalf("Expected 1 country string in result, got: %d", len(result))
	}
	if !strings.EqualFold(result[0], "us") {
		t.Fatalf("Expected country to be: %s, got: %s",
			getsrvcountry["country"][0], result[1])
	}
	empty := make(map[string][]string, 1)
	empty["abc"] = []string{""}
	result = getQStringValues(empty, "abc")
	if result != nil {
		t.Fatal("Expected nil for result")
	}
}
