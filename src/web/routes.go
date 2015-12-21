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
	getServerIDQStr = "address"
	// /queryServerID?id=
	queryServerIDQStr = "id"
	// /queryServerAddr?address=
	queryServerAddrQStr = "address"

	// getServers filters:

	// getServers?country=
	getServersCountryQStr = "country"
	// getServers?region=
	getServersRegionQStr = "region"
	// getServers?state=
	getServersStateQStr = "state"
)

var apiRoutes = []route{
	// getServers
	route{
		name:   "GetServers",
		method: "GET",
		path:   "/getServers",
		queryStrings: []querystring{
			querystring{
				name:     getServersCountryQStr,
				required: false,
			},
			querystring{
				name:     getServersRegionQStr,
				required: false,
			},
			querystring{
				name:     getServersStateQStr,
				required: false,
			},
		},
		handlerFunc: getServers,
	},
	// getServerID
	route{
		name:   "GetServerID",
		method: "GET",
		path:   "/getServerID",
		queryStrings: []querystring{
			querystring{
				name:     getServerIDQStr,
				required: true,
			},
		},
		handlerFunc: getServerID,
	},
	// queryServerID
	route{
		name:   "QueryServerID",
		method: "GET",
		path:   "/queryServerID",
		queryStrings: []querystring{
			querystring{
				name:     queryServerIDQStr,
				required: true,
			},
		},
		handlerFunc: queryServerID,
	},
	// queryServerAddr
	route{
		name:   "QueryServerAddr",
		method: "GET",
		path:   "/queryServerAddr",
		queryStrings: []querystring{
			querystring{
				name:     queryServerAddrQStr,
				required: true,
			},
		},
		handlerFunc: queryServerAddr,
	},
}
