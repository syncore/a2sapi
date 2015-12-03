package web

// handlers.go - Handler functions for API

import (
	"encoding/json"
	"fmt"
	"net"
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
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	hosts := getQStrValues(r.URL.Query(), getServerIDQueryStr)
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

func queryServerID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	ids := getQStrValues(r.URL.Query(), queryServerIDQueryStr)
	util.WriteDebug("queryServerID: ids length: %d", len(ids))
	util.WriteDebug("queryServerID: ids are: %s", ids)

	if ids[0] == "" {
		w.WriteHeader(http.StatusNotFound)
		util.WriteDebug("queryServerID: Got empty query. Ignoring.")
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

	queryServerIDRetriever(w, ids)
}

func queryServerAddr(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

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
	if !cfg.AllowDirectUserQueries {
		w.WriteHeader(http.StatusNotFound)
		// TODO: error json
		if err := json.NewEncoder(w).Encode(models.GetDefaultServerList()); err != nil {
			util.LogWebError(err)
			return
		}
		return
	}

	addresses := getQStrValues(r.URL.Query(), queryServerAddrQueryStr)
	util.WriteDebug("addresses length: %d", len(addresses))
	util.WriteDebug("addresses are: %s", addresses)

	if addresses[0] == "" {
		w.WriteHeader(http.StatusNotFound)
		util.WriteDebug("queryServerAddr: Got empty address query. Ignoring.")
		if err := json.NewEncoder(w).Encode(models.GetDefaultServerList()); err != nil {
			util.LogWebError(err)
			return
		}
		return
	}

	var parsedaddresses []string
	for _, addr := range addresses {
		ip, port, err := net.SplitHostPort(addr)
		if err != nil {
			continue
		}
		parsedaddresses = append(parsedaddresses, fmt.Sprintf("%s:%s", ip, port))
	}

	if len(parsedaddresses) == 0 {
		w.WriteHeader(http.StatusNotFound)
		util.WriteDebug("queryServerAddr: No valid addresses for query. Ignoring.")
		if err := json.NewEncoder(w).Encode(models.GetDefaultServerList()); err != nil {
			util.LogWebError(err)
			return
		}
		return
	}

	if len(parsedaddresses) > cfg.MaximumHostsPerAPIQuery {
		util.WriteDebug("Maximum number of allowed API query hosts exceeded, truncating")
		parsedaddresses = parsedaddresses[:cfg.MaximumHostsPerAPIQuery]
	}
	queryServerAddrRetriever(w, parsedaddresses)
}
