package web

// server.go - web server for API

import (
	"fmt"
	"log"
	"net/http"
	"steamtest/src/util"
)

func Start() {
	cfg, err := util.ReadConfig()
	if err != nil {
		log.Fatalf("Unable to read configuration to start web server for API: %s",
			err)
	}

	r := newRouter()

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", cfg.APIWebPort), r))
}
