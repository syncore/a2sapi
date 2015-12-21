package web

// handlers.go - Handler functions for API

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"steamtest/src/config"
	"steamtest/src/logger"
	"steamtest/src/models"
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
	// case where there's no value after query string (i.e: ?querystring=)
	if vals[0] == "" {
		return nil
	}
	return vals
}

func dumpFileAsMasterList(filename string) *models.APIServerList {
	f, err := os.Open(filename)
	if err != nil {
		logger.LogAppErrorf("Unable to open test API server dump file: %s", err)
		return nil
	}
	defer f.Close()
	r := bufio.NewReader(f)
	d := json.NewDecoder(r)
	ml := &models.APIServerList{}
	if err := d.Decode(ml); err != nil {
		logger.LogAppErrorf("Unable to decode test API server dump as json: %s", err)
		return nil
	}
	return ml
}

func getServers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	var ml *models.APIServerList
	if config.ReadConfig().DebugConfig.ServerDumpFileAsMasterList {
		ml = dumpFileAsMasterList(config.ReadConfig().DebugConfig.ServerDumpFilename)
	} else {
		ml = models.MasterList
	}

	lf := getLocationFilters(r.URL.Query())
	logger.WriteDebug("location filters: %v", lf)
	ml = filterByLocation(lf, ml)

	if err := json.NewEncoder(w).Encode(ml); err != nil {
		w.WriteHeader(http.StatusNotFound)
		logger.LogWebError(err)
		return
	}

}

func getServerID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	hosts := getQStrValues(r.URL.Query(), getServerIDQStr)
	for _, v := range hosts {
		logger.WriteDebug("host slice values: %s", v)
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
	ids := getQStrValues(r.URL.Query(), queryServerIDQStr)
	logger.WriteDebug("queryServerID: ids length: %d", len(ids))
	logger.WriteDebug("queryServerID: ids are: %s", ids)

	if ids[0] == "" {
		w.WriteHeader(http.StatusNotFound)
		logger.WriteDebug("queryServerID: Got empty query. Ignoring.")
		if err := json.NewEncoder(w).Encode(models.GetDefaultServerList()); err != nil {
			logger.LogWebError(err)
			return
		}
		return
	}
	cfg := config.ReadConfig()
	if len(ids) > cfg.WebConfig.MaximumHostsPerAPIQuery {
		logger.WriteDebug("Maximum number of allowed API query hosts exceeded, truncating")
		ids = ids[:cfg.WebConfig.MaximumHostsPerAPIQuery]
	}

	queryServerIDRetriever(w, ids)
}

func queryServerAddr(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	cfg := config.ReadConfig()
	if !cfg.WebConfig.AllowDirectUserQueries {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, `{"error":"Not allowed"}`)
		return
	}

	addresses := getQStrValues(r.URL.Query(), queryServerAddrQStr)
	logger.WriteDebug("addresses length: %d", len(addresses))
	logger.WriteDebug("addresses are: %s", addresses)

	if addresses[0] == "" {
		w.WriteHeader(http.StatusNotFound)
		logger.WriteDebug("queryServerAddr: Got empty address query. Ignoring.")
		if err := json.NewEncoder(w).Encode(models.GetDefaultServerList()); err != nil {
			logger.LogWebError(err)
			return
		}
		return
	}

	var parsedaddresses []string
	for _, addr := range addresses {
		host, err := net.ResolveTCPAddr("tcp4", addr)
		if err != nil {
			continue
		}
		parsedaddresses = append(parsedaddresses, fmt.Sprintf("%s:%d", host.IP, host.Port))
	}

	if len(parsedaddresses) == 0 {
		w.WriteHeader(http.StatusNotFound)
		logger.WriteDebug("queryServerAddr: No valid addresses for query. Ignoring.")
		if err := json.NewEncoder(w).Encode(models.GetDefaultServerList()); err != nil {
			logger.LogWebError(err)
			return
		}
		return
	}

	if len(parsedaddresses) > cfg.WebConfig.MaximumHostsPerAPIQuery {
		logger.WriteDebug("Maximum number of allowed API query hosts exceeded, truncating")
		parsedaddresses = parsedaddresses[:cfg.WebConfig.MaximumHostsPerAPIQuery]
	}
	queryServerAddrRetriever(w, parsedaddresses)
}
