package web

// routes.go - http routes for API

import "net/http"

type route struct {
	name         string
	method       string
	path         string
	queryStrings []querystring
	handlerFunc  http.HandlerFunc
}

var apiRoutes = []route{
	// getServers
	route{
		name:         "GetServers",
		method:       "GET",
		path:         "/getServers",
		queryStrings: getServersQueryStrings,
		handlerFunc:  getServers,
	},
	// getServerID
	route{
		name:         "GetServerID",
		method:       "GET",
		path:         "/getServerID",
		queryStrings: getServerIDQueryStrings,
		handlerFunc:  getServerID,
	},
	// queryServerID
	route{
		name:         "QueryServerID",
		method:       "GET",
		path:         "/queryServerID",
		queryStrings: queryServerIDQueryStrings,
		handlerFunc:  queryServerID,
	},
	// queryServerAddr
	route{
		name:         "QueryServerAddr",
		method:       "GET",
		path:         "/queryServerAddr",
		queryStrings: queryServerAddrQueryStrings,
		handlerFunc:  queryServerAddr,
	},
}
