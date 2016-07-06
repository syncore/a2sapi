package web

// server.go - Web server for API

import (
	"fmt"
	"net/http"
	"time"

	"github.com/syncore/a2sapi/src/config"
	"github.com/syncore/a2sapi/src/logger"
)

// Start listening for and responding to HTTP requests via the web server. Panics
// if unable to start.
func Start(runSilent bool) {
	r := newRouter()

	if !runSilent {
		printStartInfo()
	}

	logger.LogAppInfo("Starting HTTP server on port %d",
		config.Config.WebConfig.APIWebPort)

	srv := http.Server{
		Addr:           fmt.Sprintf(":%d", config.Config.WebConfig.APIWebPort),
		Handler:        r,
		ReadTimeout:    30 * time.Second,
		WriteTimeout:   30 * time.Second,
		MaxHeaderBytes: 1 << 20}

	err := srv.ListenAndServe()
	if err != nil {
		logger.LogAppError(err)
		panic(fmt.Sprintf("Unable to start HTTP server, error: %s\n", err))
	}
}

func printStartInfo() {
	endpoints := make([]string, len(apiRoutes))
	for _, e := range apiRoutes {
		endpoints = append(endpoints, fmt.Sprintf("%s  ", e.path))
	}
	fmt.Printf("Starting HTTP server on port %d\n", config.Config.WebConfig.APIWebPort)
	fmt.Printf("Available endpoints: %s\n", endpoints)

	if config.Config.WebConfig.AllowDirectUserQueries {
		fmt.Println("Direct (non-ID based) server API queries: enabled")
	} else {
		fmt.Println("Direct (non-ID based) server API queries: disabled")
	}

	fmt.Printf("HTTP request timeout: %d seconds\n",
		config.Config.WebConfig.APIWebTimeout)
	fmt.Printf("Maximum servers allowed per user API call: %d servers\n",
		config.Config.WebConfig.MaximumHostsPerAPIQuery)
}
