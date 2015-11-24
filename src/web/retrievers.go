package web

// retrievers.go - Bridge between http requests and database layer

import (
	"encoding/json"
	"net/http"
	"steamtest/src/db"
	"steamtest/src/util"
	"steamtest/src/web/models"
)

func getServerIDsRetriever(w http.ResponseWriter, hosts []string) {
	m := make(chan map[int64]string, len(hosts))
	sdb, err := db.OpenServerDB()
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		if err := json.NewEncoder(w).Encode(models.GetDefaultServerID()); err != nil {
			util.LogWebError(err)
			return
		}
		return
	}
	go db.GetIDsForAPIQuery(m, sdb, hosts)
	ids := <-m
	if len(ids) > 0 {
		serverID := &models.ServerID{}
		for k, v := range ids {
			s := &models.Server{
				ID:   k,
				Host: v,
			}
			serverID.Servers = append(serverID.Servers, s)
		}
		serverID.ServerCount = len(serverID.Servers)

		if err := json.NewEncoder(w).Encode(serverID); err != nil {
			util.LogWebError(err)
			return
		}
	} else {
		w.WriteHeader(http.StatusNotFound)
		if err := json.NewEncoder(w).Encode(models.GetDefaultServerID()); err != nil {
			util.LogWebError(err)
			return
		}
	}
}
