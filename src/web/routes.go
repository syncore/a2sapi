package web

import "net/http"

// routes.go - http routes for API

type route struct {
	name         string
	method       string
	path         string
	queryStrings []querystring
	handlerFunc  http.HandlerFunc
}

type querystring struct {
	name     string
	required bool
}

const (
	// /getServerID?address=
	getServerIDQueryStr = "address"
	// /queryServerID?id=
	queryServerIDQueryStr = "id"
	// /queryServerAddr?address=
	queryServerAddrQueryStr = "address"
)

var apiRoutes = []route{
	route{
		name:   "GetServers",
		method: "GET",
		path:   "/getServers",
		// queryStrings: []querystring{
		// 	querystring{
		// 		name:     getServerIDQueryStr,
		// 		required: true,
		// 	},
		// },
		handlerFunc: getServers,
	},
	route{
		name:   "GetServerID",
		method: "GET",
		path:   "/getServerID",
		queryStrings: []querystring{
			querystring{
				name:     getServerIDQueryStr,
				required: true,
			},
		},
		handlerFunc: getServerID,
	},
	route{
		name:   "QueryServerID",
		method: "GET",
		path:   "/queryServerID",
		queryStrings: []querystring{
			querystring{
				name:     queryServerIDQueryStr,
				required: true,
			},
		},
		handlerFunc: queryServerID,
	},
	route{
		name:   "QueryServerAddr",
		method: "GET",
		path:   "/queryServerAddr",
		queryStrings: []querystring{
			querystring{
				name:     queryServerAddrQueryStr,
				required: true,
			},
		},
		handlerFunc: queryServerAddr,
	},
}
