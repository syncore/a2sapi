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
	// servers
	route{
		name:         "GetServers",
		method:       "GET",
		path:         "/servers",
		queryStrings: getServersQueryStrings,
		handlerFunc:  getServers,
	},
	// serverID
	route{
		name:         "GetServerIDs",
		method:       "GET",
		path:         "/serverIDs",
		queryStrings: getServerIDsQueryStrings,
		handlerFunc:  getServerIDs,
	},
	// query - by ID
	route{
		name:         "QueryServerID",
		method:       "GET",
		path:         "/query",
		queryStrings: queryServerIDQueryStrings,
		handlerFunc:  queryServerIDs,
	},
	// query - by address
	route{
		name:         "QueryServerAddr",
		method:       "GET",
		path:         "/query",
		queryStrings: queryServerAddrQueryStrings,
		handlerFunc:  queryServerAddrs,
	},
}
