package web

// server.go - Web server for API

import (
	"fmt"
	"net/http"
	"steamtest/src/config"
	"steamtest/src/logger"
)

// Start listening for and responding to HTTP requests via the web server. Panics
// if unable to start.
func Start(runSilent bool) {
	cfg := config.ReadConfig()
	r := newRouter(cfg)

	if !runSilent {
		printStartInfo(cfg)
	}

	logger.LogAppInfo("Starting HTTP server on port %d", cfg.WebConfig.APIWebPort)
	err := http.ListenAndServe(fmt.Sprintf(":%d", cfg.WebConfig.APIWebPort), r)
	if err != nil {
		logger.LogAppError(err)
		panic(fmt.Sprintf("Unable to start HTTP server, error: %s\n", err))
	}
}

func printStartInfo(cfg *config.Config) {
	endpoints := make([]string, len(apiRoutes))
	for _, e := range apiRoutes {
		endpoints = append(endpoints, fmt.Sprintf("%s  ", e.path))
	}
	fmt.Printf("Starting HTTP server on port %d\n", cfg.WebConfig.APIWebPort)
	fmt.Printf("Available endpoints: %s\n", endpoints)

	if cfg.WebConfig.AllowDirectUserQueries {
		fmt.Println("Direct (non-ID based) server API queries: enabled")
	} else {
		fmt.Println("Direct (non-ID based) server API queries: disabled")
	}

	fmt.Printf("HTTP request timeout: %d seconds\n", cfg.WebConfig.APIWebTimeout)
	fmt.Printf("Maximum servers allowed per user API call: %d servers\n",
		cfg.WebConfig.MaximumHostsPerAPIQuery)
}
