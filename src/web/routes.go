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
	getServerIDsPath     = "/getServerIDs"
	getServerIDsQueryStr = "address"
)

type routes []route

var apiRoutes = routes{
	route{
		name:        "GetServerIDs",
		method:      "GET",
		path:        getServerIDsPath,
		queryString: getServerIDsQueryStr,
		handlerFunc: getServerIDs,
	},
}
