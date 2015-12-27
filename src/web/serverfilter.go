package web

// serverfilter.go - operations for filtering the server list based on query
// string data.

import (
	"steamtest/src/models"
	"strings"
)

// getSrvFilterFromQString determines whether a query string has values and
//  builds a filter that will be used on the server list.
func getSrvFilterFromQString(m map[string][]string, qs []querystring) []slQueryFilter {
	var qfilters []slQueryFilter
	for key := range m {
		for _, q := range qs {
			if strings.EqualFold(key, q.name) {
				vals := getQStringValues(m, key)
				if len(vals) > 0 {
					qfilters = append(qfilters, slQueryFilter{name: q.name,
						needsbool: q.boolonly, values: vals})
				}
			}
		}
	}
	return qfilters
}

func findMatches(sqf slQueryFilter,
	servers []*models.APIServer) []*models.APIServer {
	var matched []*models.APIServer
	var ssearch string
	var bsearch bool

	for _, srv := range servers {
		switch sqf.name {
		// location-based
		case qsGetServersRegion:
			ssearch = srv.CountryInfo.Continent
		case qsGetServersCountry:
			ssearch = srv.CountryInfo.CountryCode
		case qsGetServersState:
			ssearch = srv.CountryInfo.State
		// info-based
		case qsGetServersName:
			ssearch = srv.Info.Name
		case qsGetServersType:
			ssearch = srv.Info.ServerType
		case qsGetServersOS:
			ssearch = srv.Info.Environment
		case qsGetServersVersion:
			ssearch = srv.Info.Version
		case qsGetServersKeywords:
			ssearch = srv.Info.ExtraData.Keywords
		case qsGetServersHasPlayers:
			if strings.EqualFold(sqf.values[0], "true") {
				bsearch = srv.Info.Players > 0
			} else {
				bsearch = srv.Info.Players == 0
			}
		case qsGetServersHasBots:
			if strings.EqualFold(sqf.values[0], "true") {
				bsearch = srv.Info.Bots > 0
			} else {
				bsearch = srv.Info.Bots == 0
			}
		case qsGetServersHasPassword:
			if strings.EqualFold(sqf.values[0], "true") {
				bsearch = srv.Info.Visibility == 1
			} else {
				bsearch = srv.Info.Visibility == 0
			}
		case qsGetServersHasAntiCheat:
			if strings.EqualFold(sqf.values[0], "true") {
				bsearch = srv.Info.VAC == 1
			} else {
				bsearch = srv.Info.VAC == 0
			}
		}
		if sqf.needsbool {
			if strings.EqualFold(sqf.values[0], "true") && bsearch {
				matched = append(matched, srv)
			} else if strings.EqualFold(sqf.values[0], "false") && !bsearch {
				matched = append(matched, srv)
			}
		} else {
			for _, val := range sqf.values {
				if strings.EqualFold(ssearch, val) {
					matched = append(matched, srv)
				}
			}
		}
	}
	return matched
}

// filterServers takes the server filters and the last retrieved server list and
// returns a new, filtered server list based on the matched filters.
func filterServers(sqf []slQueryFilter,
	sl *models.APIServerList) *models.APIServerList {

	if sl == nil {
		return models.GetDefaultServerList()
	}

	filtered := sl.Servers
	for _, s := range sqf {
		filtered = findMatches(s, filtered)
	}
	if filtered == nil {
		// JSON empty array instead of null
		filtered = make([]*models.APIServer, 0)
	}
	sl.Servers = filtered
	sl.ServerCount = len(filtered)
	sl.FailedCount = 0
	sl.FailedServers = make([]string, 0)
	return sl
}
