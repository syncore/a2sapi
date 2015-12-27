package web

// Tests for handler functions for API
// Essentially these are just simple tests for the status codes; the actual tests that target the data retrieved by these endpoints is performed in the tests for the db package.

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gorilla/mux"
)

const (
	testPort = 40081
)

var (
	testUrlBase = fmt.Sprintf("http://:%d", testPort)
)

func formatUrl(path string) string {
	return fmt.Sprintf("%s/%s", testUrlBase, path)
}

func init() {
	// need base directory
	err := os.Chdir("../..")
	if err != nil {
		panic("Unable to change directory for tests")
	}

	go func() {
		r := mux.NewRouter().StrictSlash(true)
		for _, ar := range apiRoutes {
			var handler http.Handler
			handler = ar.handlerFunc

			r.Methods(ar.method).
				MatcherFunc(pathQStrToLowerMatcherFunc(r, ar.path, ar.queryStrings,
				getRequiredQryStringCount(ar.queryStrings))).
				Name(ar.name).
				Handler(http.TimeoutHandler(handler,
				time.Duration(7)*time.Second,
				`{"error":"Timeout"}`))
		}
		err := http.ListenAndServe(fmt.Sprintf(":%d", testPort), r)
		if err != nil {
			panic("Unable to start web server")
		}
	}()
}

func TestGetServers(t *testing.T) {
	r, _ := http.NewRequest("GET", formatUrl("getServers"), nil)
	w := httptest.NewRecorder()
	getServers(w, r)
	if w.Code != http.StatusOK {
		t.Errorf("Expected status code: %v for GetServers handler; got: %v", http.StatusOK, w.Code)
	}
	if len(w.Body.Bytes()) == 0 {
		t.Errorf("Response body should not be empty")
	}
}

func TestGetServerID(t *testing.T) {
	r, _ := http.NewRequest("GET", formatUrl("getServerID?addr=127.0.0.1:65534"), nil)
	w := httptest.NewRecorder()
	getServerID(w, r)
	// this actually should be 404
	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status code: %v for GetServerID handler; got: %v", http.StatusNotFound, w.Code)
	}
	if len(w.Body.Bytes()) == 0 {
		t.Errorf("GetServerID handler body should not be empty")
	}
}

func TestQueryServerID(t *testing.T) {
	r, _ := http.NewRequest("GET", formatUrl("queryServerID?id=788593993848"), nil)
	w := httptest.NewRecorder()
	queryServerID(w, r)
	// should be 404
	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status code %v for queryServerID handler; got: %v", http.StatusNotFound, w.Code)
	}
	if len(w.Body.Bytes()) == 0 {
		t.Errorf("queryServerID handler body should not be empty")
	}
}

func TestQueryServerAddr(t *testing.T) {
	r1, _ := http.NewRequest("GET", formatUrl("queryServerAddr?address=127.0.0.1:65534"), nil)
	w1 := httptest.NewRecorder()
	queryServerAddr(w1, r1)
	// 200 - default server list
	if w1.Code != http.StatusOK {
		t.Errorf("Expected status code %v for queryServerAddr handler; got: %v", http.StatusNotFound, w1.Code)
	}
	if len(w1.Body.Bytes()) == 0 {
		t.Errorf("queryServerAddr handler body should not be empty")
	}
	// 404 - no addreses specified
	r2, _ := http.NewRequest("GET", formatUrl("queryServerAddr?address="), nil)
	w2 := httptest.NewRecorder()
	queryServerAddr(w2, r2)
	if w2.Code != http.StatusNotFound {
		t.Errorf("Expected status code %v for queryServerAddr handler; got: %v", http.StatusNotFound, w2.Code)
	}
	if len(w2.Body.Bytes()) == 0 {
		t.Errorf("queryServerAddr handler body should not be empty")
	}
}
