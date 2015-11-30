package web

import "net/http"

// routes.go - http routes for API

type route struct {
	name        string
	method      string
	path        string
	queryString string
	handlerFunc http.HandlerFunc
}

const (
	getServerIDPath     = "/getServerIDs"
	getServerIDQueryStr = "address"
	queryServerPath     = "/queryServer"
	queryServerQueryStr = "ids"
)

type routes []route

var apiRoutes = routes{
	route{
		name:        "GetServerIDs",
		method:      "GET",
		path:        getServerIDPath,
		queryString: getServerIDQueryStr,
		handlerFunc: getServerID,
	},
	route{
		name:        "QueryServer",
		method:      "GET",
		path:        queryServerPath,
		queryString: queryServerQueryStr,
		handlerFunc: queryServer,
	},
}
