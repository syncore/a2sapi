package web

// Tests for server filtering for getServers API endpoint

import (
	"a2sapi/src/constants"
	"a2sapi/src/models"
	"encoding/json"
	"strings"
	"testing"
)

func TestGetSrvFilterFromQString(t *testing.T) {
	query := make(map[string][]string, 3)
	query["regions"] = []string{"North America,Europe"}
	query["serverVersions"] = []string{
		"1066",
	}
	query["hasPlayers"] = []string{
		"true",
	}
	sqf := getSrvFilterFromQString(query, getServersQueryStrings)
	if len(sqf) != 3 {
		t.Fatalf("Expected server list query filters length of 3, got %d",
			len(sqf))
	}
	// Names
	found, expectedfound := 0, 3
	for _, s := range sqf {
		for key := range query {
			if strings.EqualFold(s.name, key) {
				found++
			}
		}
	}
	if found != expectedfound {
		t.Fatalf("Expected server list query filter to match %d elements, got: %d",
			expectedfound, found)
	}
	// Values
	found, expectedfound = 0, 4
	for _, s := range sqf {
		for _ = range s.values {
			found++
		}
	}
	if found != expectedfound {
		t.Fatalf("Expected server list query filter to contain %d values, got: %d",
			expectedfound, found)
	}
	// Boolean
	for _, s := range sqf {
		if strings.EqualFold(s.name, qsGetServersHasPlayers) {
			if s.needsbool {
				break
			} else {
				t.Fatalf("Expected that %v query string filter required bool",
					qsGetServersHasPlayers)
			}
		}
	}
}

func TestFindMatches(t *testing.T) {
	hasPlayersFilter := slQueryFilter{
		name:      qsGetServersHasPlayers,
		needsbool: true,
		values:    []string{"true"},
	}
	serverNameFilter := slQueryFilter{
		name:      qsGetServersName,
		needsbool: false,
		values:    []string{"syncore"},
	}
	stateFilter := slQueryFilter{
		name:      qsGetServersState,
		needsbool: false,
		values:    []string{"TX", "NY"},
	}
	src := &models.APIServerList{}
	err := json.Unmarshal(constants.TestServerDumpJSON, src)
	if err != nil {
		t.Fatalf("Failed to read test server data: %s", err)
	}
	// ?hasPlayers=true
	matches := findMatches(hasPlayersFilter, src.Servers)
	expected := 1
	if matches == nil {
		t.Fatal("Matches should not be nil")
	}
	if len(matches) != expected {
		t.Fatalf("Expected %d match(es), got: %d", expected, len(matches))
	}
	// ?serverName=syncore
	matches = findMatches(serverNameFilter, src.Servers)
	expected = 2
	if matches == nil {
		t.Fatal("Matches should not be nil")
	}
	if len(matches) != expected {
		t.Fatalf("Expected %d match(es), got: %d", expected, len(matches))
	}
	// ?state=TX,NY
	matches = findMatches(stateFilter, src.Servers)
	expected = 2
	if matches == nil {
		t.Fatal("Matches should not be nil")
	}
	if len(matches) != expected {
		t.Fatalf("Expected %d match(es), got: %d", expected, len(matches))
	}
}

func TestFilterServers(t *testing.T) {
	filters := []slQueryFilter{
		slQueryFilter{
			name:      qsGetServersHasPlayers,
			needsbool: true,
			values:    []string{"true"},
		},
		slQueryFilter{
			name:      qsGetServersName,
			needsbool: false,
			values:    []string{"pixel"},
		},
		slQueryFilter{
			name:      qsGetServersState,
			needsbool: false,
			values:    []string{"VA"},
		},
	}
	src := &models.APIServerList{}
	err := json.Unmarshal(constants.TestServerDumpJSON, src)
	if err != nil {
		t.Fatalf("Failed to read test server data: %s", err)
	}
	servers := filterServers(filters, src)
	if servers == nil {
		t.Fatal("Servers returned should not be nil")
	}
	if len(servers.Servers) != 1 {
		t.Fatalf("Expected 1 match, got: %d", len(servers.Servers))
	}
	for _, s := range servers.Servers {
		if !strings.EqualFold("54.172.5.67:25801", s.Host) {
			t.Fatalf("Expected matched host to be: 54.172.5.67:25801, got: %s",
				s.Host)
		}
	}
}
