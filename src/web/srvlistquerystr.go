package web

import (
	"steamtest/src/models"
	"strings"
)

// slquerystr.go - server list query string filter operations

type slQueryFilter struct {
	filterType   string
	filterValues []string
}

// getLocationFilters determines whether a query string has valid region, country,
// and/or state filters and retrieves them if so as well as extracting their values.
func getLocationFilters(m map[string][]string) []slQueryFilter {
	var qfilters []slQueryFilter
	for key := range m {
		if strings.EqualFold(key, getServersCountryQStr) {
			countries := getQStrValues(m, key)
			if len(countries) > 0 {
				qfilters = append(qfilters,
					slQueryFilter{filterType: getServersCountryQStr, filterValues: countries})
			}
		}
		if strings.EqualFold(key, getServersRegionQStr) {
			regions := getQStrValues(m, key)
			if len(regions) > 0 {
				qfilters = append(qfilters,
					slQueryFilter{filterType: getServersRegionQStr, filterValues: regions})
			}
		}
		if strings.EqualFold(key, getServersStateQStr) {
			states := getQStrValues(m, key)
			if len(states) > 0 {
				qfilters = append(qfilters,
					slQueryFilter{filterType: getServersStateQStr, filterValues: states})
			}
		}
	}
	return qfilters
}

// filterByLocation takes the location filters and the last retrieved server list
// and creates a new server list based on the matched location filter information.
func filterByLocation(sqf []slQueryFilter,
	sl *models.APIServerList) *models.APIServerList {
	doFilter := func(filter string, searchvals []string,
		servers []*models.APIServer) []*models.APIServer {
		var matched []*models.APIServer
		var search string
		for _, srv := range servers {
			switch filter {
			case getServersRegionQStr:
				search = srv.CountryInfo.Continent
			case getServersCountryQStr:
				search = srv.CountryInfo.CountryCode
			case getServersStateQStr:
				search = srv.CountryInfo.State
			}
			for _, val := range searchvals {
				if strings.EqualFold(search, val) {
					matched = append(matched, srv)
				}
			}
		}
		return matched
	}
	filtered := sl.Servers
	for _, s := range sqf {
		if s.filterType == getServersRegionQStr {
			filtered = doFilter(getServersRegionQStr, s.filterValues, filtered)
		}
		if s.filterType == getServersCountryQStr {
			filtered = doFilter(getServersCountryQStr, s.filterValues, filtered)
		}
		if s.filterType == getServersStateQStr {
			filtered = doFilter(getServersStateQStr, s.filterValues, filtered)
		}
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
