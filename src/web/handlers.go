package web

// handlers.go - Handler functions for API

import (
	"encoding/json"
	"fmt"
	"net/http"
	"steamtest/src/web/models"
	"strings"
)

func getServerIDs(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("URL: %s\tPATH: %s\n", r.URL, r.URL.Path)
	fmt.Printf("Query: %v\n", r.URL.Query())
	var hosts []string
	// already matched by MatcherFunc; this allows case-insensitive query str lookups
	for k := range r.URL.Query() {
		if strings.EqualFold(k, getServerIDsQueryStr) {
			hosts = strings.Split(r.URL.Query().Get(k), ",")
			break
		}
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	for _, v := range hosts {
		fmt.Printf("vars host slice value: %s\n", v)
		// basically require at least 2 octets
		if len(v) < 4 {
			w.WriteHeader(http.StatusBadRequest)
			if err := json.NewEncoder(w).Encode(models.GetDefaultServerID()); err != nil {
				return
			}
			return
		}
	}
	getServerIDsRetriever(w, hosts)
}
