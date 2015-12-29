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
	"steamtest/src/constants"
	"steamtest/src/logger"
	"steamtest/src/models"
)

func useDumpFileAsMasterList(dumppath string) *models.APIServerList {
	f, err := os.Open(dumppath)
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
	var asl *models.APIServerList
	cfg := config.ReadConfig()

	if cfg.DebugConfig.ServerDumpFileAsMasterList {
		asl = useDumpFileAsMasterList(constants.DumpFileFullPath(
			cfg.DebugConfig.ServerDumpFilename))
	} else {
		asl = models.MasterList
	}
	// Empty (i.e. during first retrieval/startup)
	if asl == nil {
		writeJSONResponse(w, models.GetDefaultServerList())
		return
	}
	srvfilters := getSrvFilterFromQString(r.URL.Query(), getServersQueryStrings)
	logger.WriteDebug("server list will be filtered with: %v", srvfilters)
	list := filterServers(srvfilters, asl)
	writeJSONResponse(w, list)
}

func getServerID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	hosts := getQStringValues(r.URL.Query(), qsGetServerID)
	for _, v := range hosts {
		logger.WriteDebug("host slice values: %s", v)
		// basically require at least 2 octets
		if len(v) < 4 {
			w.WriteHeader(http.StatusBadRequest)
			writeJSONResponse(w, models.GetDefaultServerID())
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

	if ids == nil {
		w.WriteHeader(http.StatusNotFound)
		logger.WriteDebug("queryServerID: Got empty query. Ignoring.")
		writeJSONResponse(w, models.GetDefaultServerList())
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

	if addresses == nil {
		w.WriteHeader(http.StatusNotFound)
		logger.WriteDebug("queryServerAddr: Got empty address query. Ignoring.")
		writeJSONResponse(w, models.GetDefaultServerList())
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
		writeJSONResponse(w, models.GetDefaultServerList())
		return
	}

	if len(parsedaddresses) > cfg.WebConfig.MaximumHostsPerAPIQuery {
		logger.WriteDebug("Maximum number of allowed API query hosts exceeded, truncating")
		parsedaddresses = parsedaddresses[:cfg.WebConfig.MaximumHostsPerAPIQuery]
	}
	queryServerAddrRetriever(w, parsedaddresses)
}

// writeJSONResponse encodes data as JSON and writes it to w; if unsuccessful,
// the error will be logged and a generic error message will be displayed to the user.
func writeJSONResponse(w http.ResponseWriter, data interface{}) {
	if err := json.NewEncoder(w).Encode(data); err != nil {
		writeJSONEncodeError(w, err)
	}
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
