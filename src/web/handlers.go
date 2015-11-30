package web

// handlers.go - Handler functions for API

import (
	"encoding/json"
	"net/http"
	"steamtest/src/models"
	"steamtest/src/util"
	"strings"
)

func getQStrValues(m map[string][]string, querystring string) []string {
	var vals []string
	for k := range m {
		if strings.EqualFold(k, querystring) {
			vals = strings.Split(m[k][0], ",")
			break
		}
	}
	return vals
}

func getServerID(w http.ResponseWriter, r *http.Request) {
	util.WriteDebug("URL: %s\tPATH: %s", r.URL, r.URL.Path)
	util.WriteDebug("Query: %v", r.URL.Query())

	hosts := getQStrValues(r.URL.Query(), getServerIDQueryStr)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	for _, v := range hosts {
		util.WriteDebug("host slice values: %s", v)
		// basically require at least 2 octets
		if len(v) < 4 {
			w.WriteHeader(http.StatusBadRequest)
			if err := json.NewEncoder(w).Encode(models.GetDefaultServerID()); err != nil {
				return
			}
			return
		}
	}
	getServerIDRetriever(w, hosts)
}

func queryServer(w http.ResponseWriter, r *http.Request) {
	util.WriteDebug("URL: %s\tPATH: %s", r.URL, r.URL.Path)
	util.WriteDebug("Query: %v", r.URL.Query())

	ids := getQStrValues(r.URL.Query(), queryServerQueryStr)
	if ids[0] == "" {
		w.WriteHeader(http.StatusNotFound)
		util.WriteDebug("Got empty query. Ignoring.")
		if err := json.NewEncoder(w).Encode(models.GetDefaultServerList()); err != nil {
			util.LogWebError(err)
			return
		}
		return
	}
	cfg, err := util.ReadConfig()
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		util.LogAppError(err)
		if err := json.NewEncoder(w).Encode(models.GetDefaultServerList()); err != nil {
			util.LogWebError(err)
			return
		}
		return
	}
	if len(ids) > cfg.MaximumHostsPerAPIQuery {
		util.WriteDebug("Maximum number of allowed API query hosts exceeded, truncating")
		ids = ids[:cfg.MaximumHostsPerAPIQuery]
	}
	util.WriteDebug("ids length: %d", len(ids))
	util.WriteDebug("ids are: %s", ids)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	util.WriteDebug("id slice values: %s", ids)
	queryServerRetriever(w, ids)
}
