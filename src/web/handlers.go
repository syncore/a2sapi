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
)

func useDumpFileAsMasterList(filename string) *models.APIServerList {
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
		ml = useDumpFileAsMasterList(config.ReadConfig().DebugConfig.ServerDumpFilename)
	} else {
		ml = models.MasterList
	}
	// Master list is empty (i.e. during first retrieval/startup)
	if ml == nil {
		if err := json.NewEncoder(w).Encode(models.GetDefaultServerList()); err != nil {
			writeJSONEncodeError(w, err)
			return
		}
		return
	}

	srvfilters := getSrvFilterFromQString(r.URL.Query(), getServersQueryStrings)
	logger.WriteDebug("server list will be filtered with: %v", srvfilters)
	ml = filterServers(srvfilters, ml)

	if err := json.NewEncoder(w).Encode(ml); err != nil {
		writeJSONEncodeError(w, err)
		return
	}

}

func getServerID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	hosts := getQStringValues(r.URL.Query(), qsGetServerID)
	for _, v := range hosts {
		logger.WriteDebug("host slice values: %s", v)
		// basically require at least 2 octets
		if len(v) < 4 {
			w.WriteHeader(http.StatusBadRequest)
			if err := json.NewEncoder(w).Encode(models.GetDefaultServerID()); err != nil {
				writeJSONEncodeError(w, err)
				return
			}
			return
		}
	}
	getServerIDRetriever(w, hosts)
}

func queryServerID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	ids := getQStringValues(r.URL.Query(), qsQueryServerID)
	logger.WriteDebug("queryServerID: ids length: %d", len(ids))
	logger.WriteDebug("queryServerID: ids are: %s", ids)

	if ids[0] == "" {
		w.WriteHeader(http.StatusNotFound)
		logger.WriteDebug("queryServerID: Got empty query. Ignoring.")
		if err := json.NewEncoder(w).Encode(models.GetDefaultServerList()); err != nil {
			writeJSONEncodeError(w, err)
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

	addresses := getQStringValues(r.URL.Query(), qsQueryServerAddr)
	logger.WriteDebug("addresses length: %d", len(addresses))
	logger.WriteDebug("addresses are: %s", addresses)

	if addresses[0] == "" {
		w.WriteHeader(http.StatusNotFound)
		logger.WriteDebug("queryServerAddr: Got empty address query. Ignoring.")
		if err := json.NewEncoder(w).Encode(models.GetDefaultServerList()); err != nil {
			writeJSONEncodeError(w, err)
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
			writeJSONEncodeError(w, err)
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

// setNotFoundAndLog sets the error code of the underlying writer to 404 (not found)
// and internally logs the error.
func setNotFoundAndLog(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusNotFound)
	logger.LogWebError(err)
}

// writeJSONEncodeError displays a generic error message, returns an error code
// of 404 not found, and logs an error related to unsuccessful JSON encoding.
func writeJSONEncodeError(w http.ResponseWriter, err error) {
	setNotFoundAndLog(w, err)
	fmt.Fprintf(w, `{"error":"An error occurred."}`)
}
