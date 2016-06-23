package web

// handlers.go - Handler functions for API

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"

	"github.com/syncore/a2sapi/src/config"
	"github.com/syncore/a2sapi/src/constants"
	"github.com/syncore/a2sapi/src/logger"
	"github.com/syncore/a2sapi/src/models"
)

func compressGzip(hf http.HandlerFunc, shouldCompress bool) http.Handler {
	if !shouldCompress {
		return hf
	}
	return GzipHandler(hf)
}

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

	if config.Config.DebugConfig.ServerDumpFileAsMasterList {
		asl = useDumpFileAsMasterList(constants.DumpFileFullPath(
			config.Config.DebugConfig.ServerDumpFilename))
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

func getServerIDs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	hosts := getQStringValues(r.URL.Query(), qsGetServerIDs)
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

func queryServerIDs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	ids := getQStringValues(r.URL.Query(), qsQueryServerIDs)
	logger.WriteDebug("queryServerID: ids length: %d", len(ids))
	logger.WriteDebug("queryServerID: ids are: %s", ids)

	if ids == nil {
		w.WriteHeader(http.StatusOK)
		logger.WriteDebug("queryServerID: Got empty query. Ignoring.")
		writeJSONResponse(w, models.GetDefaultServerList())
		return
	}
	if len(ids) > config.Config.WebConfig.MaximumHostsPerAPIQuery {
		logger.WriteDebug("Maximum number of allowed API query hosts exceeded, truncating")
		ids = ids[:config.Config.WebConfig.MaximumHostsPerAPIQuery]
	}

	queryServerIDRetriever(w, ids)
}

func queryServerAddrs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	if !config.Config.WebConfig.AllowDirectUserQueries {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w,
			`{"error": {"code": 400,"message": "Direct server queries are disabled. Use the %s parameter."}}`,
			qsQueryServerIDs)
		return
	}
	addresses := getQStringValues(r.URL.Query(), qsQueryServerAddrs)
	logger.WriteDebug("addresses length: %d", len(addresses))
	logger.WriteDebug("addresses are: %s", addresses)

	if addresses == nil {
		w.WriteHeader(http.StatusOK)
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
		w.WriteHeader(http.StatusOK)
		logger.WriteDebug("queryServerAddr: No valid addresses for query. Ignoring.")
		writeJSONResponse(w, models.GetDefaultServerList())
		return
	}

	if len(parsedaddresses) > config.Config.WebConfig.MaximumHostsPerAPIQuery {
		logger.WriteDebug("Maximum number of allowed API query hosts exceeded, truncating")
		parsedaddresses = parsedaddresses[:config.Config.WebConfig.MaximumHostsPerAPIQuery]
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
	fmt.Fprintf(w, `{"error": {"code": 400,"message": "JSON encoding error."}}`)
}
