package web

// retrievers.go - Bridge between http requests and database (and potentially other) layers

import (
	"encoding/json"
	"net/http"
	"steamtest/src/db"
	"steamtest/src/models"
	"steamtest/src/steam"
	"steamtest/src/util"
)

func getServerIDRetriever(w http.ResponseWriter, hosts []string) {
	m := make(chan *models.DbServerID, 1)
	sdb, err := db.OpenServerDB()
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		util.LogWebError(err)
		if err := json.NewEncoder(w).Encode(models.GetDefaultServerID()); err != nil {
			w.WriteHeader(http.StatusNotFound)
			util.LogWebError(err)
			return
		}
		return
	}
	defer sdb.Close()
	go db.GetIDsAPIQuery(m, sdb, hosts)
	ids := <-m
	if len(ids.Servers) > 0 {
		if err := json.NewEncoder(w).Encode(ids); err != nil {
			w.WriteHeader(http.StatusNotFound)
			util.LogWebError(err)
			return
		}
	} else {
		w.WriteHeader(http.StatusNotFound)
		if err := json.NewEncoder(w).Encode(models.GetDefaultServerID()); err != nil {
			w.WriteHeader(http.StatusNotFound)
			util.LogWebError(err)
			return
		}
	}
}

func queryServerIDRetriever(w http.ResponseWriter, ids []string) {
	s := make(chan map[string]string, len(ids))
	sdb, err := db.OpenServerDB()
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		util.LogWebError(err)
		if err := json.NewEncoder(w).Encode(models.GetDefaultServerList()); err != nil {
			w.WriteHeader(http.StatusNotFound)
			util.LogWebError(err)
		}
		return
	}
	defer sdb.Close()
	db.GetHostsAndGameFromIDAPIQuery(s, sdb, ids)
	hostsgames := <-s
	if len(hostsgames) == 0 {
		w.WriteHeader(http.StatusNotFound)
		if err := json.NewEncoder(w).Encode(models.GetDefaultServerList()); err != nil {
			util.LogWebError(err)
		}
		return
	}
	serverlist, err := steam.Query(hostsgames)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		util.LogWebError(err)
		if err := json.NewEncoder(w).Encode(models.GetDefaultServerList()); err != nil {
			w.WriteHeader(http.StatusNotFound)
			util.LogWebError(err)
			return
		}
		return
	}
	if err := json.NewEncoder(w).Encode(serverlist); err != nil {
		w.WriteHeader(http.StatusNotFound)
		util.LogWebError(err)
	}
}

func queryServerAddrRetriever(w http.ResponseWriter, addresses []string) {
	serverlist, err := steam.DirectQuery(addresses)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		util.LogWebError(err)
		if err := json.NewEncoder(w).Encode(models.GetDefaultServerList()); err != nil {
			w.WriteHeader(http.StatusNotFound)
			util.LogWebError(err)
			return
		}
		return
	}
	if err := json.NewEncoder(w).Encode(serverlist); err != nil {
		w.WriteHeader(http.StatusNotFound)
		util.LogWebError(err)
	}
}
