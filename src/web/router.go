package web

// router.go - request router

import (
	"fmt"
	"net/http"
	"strings"

	"steamtest/src/util"

	"github.com/gorilla/mux"
)

func newRouter() *mux.Router {
	r := mux.NewRouter().StrictSlash(true)
	for _, ar := range apiRoutes {
		var handler http.Handler
		handler = ar.handlerFunc
		handler = util.LogWebRequest(handler, ar.name)

		r.Methods(ar.method).
			MatcherFunc(pathQStrToLowerMatcherFunc(r, ar.path, ar.queryString)).
			Name(ar.name).
			Handler(handler)
	}
	return r
}

// Provide case-insensitive matching for URL paths and query strings
func pathQStrToLowerMatcherFunc(router *mux.Router,
	routepath string, querystring string) func(req *http.Request,
	rt *mux.RouteMatch) bool {
	return func(req *http.Request, rt *mux.RouteMatch) bool {
		var pathok bool
		var qstrok bool
		// case-insensitive paths
		if strings.HasPrefix(strings.ToLower(req.URL.Path), strings.ToLower(routepath)) {
			fmt.Printf("PATH: %s matches route path: %s\n", req.URL.Path, routepath)
			pathok = true
		}
		//case-insensitive query strings
		// not all API routes will make use of query strings
		if querystring == "" {
			qstrok = true
		} else {
			qry := req.URL.Query()
			for key := range qry {
				fmt.Printf("key is: %s\n", key)
				if strings.EqualFold(key, querystring) {
					fmt.Printf("%s matches %s\n", key, querystring)
					qstrok = true
					break
				}
			}
		}
		return pathok && qstrok
	}
}
