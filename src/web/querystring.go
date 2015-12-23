package web

import "strings"

// querystring.go - URL query string definitions and helper functions

type querystring struct {
	name     string
	boolonly bool
	required bool
}

type slQueryFilter struct {
	name      string
	needsbool bool
	values    []string
}

// query string names
const (
	// getServerID:
	// ?address=
	qsGetServerID = "address"

	// queryServerID:
	// ?id=
	qsQueryServerID = "id"

	// /queryServerAddr:
	// ?address
	qsQueryServerAddr = "address"

	// getServers:
	// ?country=
	qsGetServersCountry = "country"
	// ?region=
	qsGetServersRegion = "region"
	// ?state=
	qsGetServersState = "state"
	// info filtering
	// ?serverName=
	qsGetServersName = "serverName"
	// ?map=
	qsGetServersMap = "map"
	// ?gamedir=
	qsGetServersGameDir = "gameDir"
	// ?game=
	qsGetServersGame = "game"
	// ?serverType=
	qsGetServersType = "serverType"
	// ?serverOS=
	qsGetServersOS = "serverOS"
	// ?serverVersion=
	qsGetServersVersion = "serverVersion"
	// ?serverKeywords=
	qsGetServersKeywords = "serverKeywords"
	// ?hasPlayers= (bool)
	qsGetServersHasPlayers = "hasPlayers"
	// ?hasBots= (bool)
	qsGetServersHasBots = "hasBots"
	// ?hasPassword= (bool)
	qsGetServersHasPassword = "hasPassword"
	// ?hasAntiCheat= (bool)
	qsGetServersHasAntiCheat = "hasAntiCheat"
)

// getServerID query strings
var getServerIDQueryStrings = []querystring{
	querystring{
		name:     qsGetServerID,
		required: true,
	},
}

// queryServerID query strings
var queryServerIDQueryStrings = []querystring{
	querystring{
		name:     qsQueryServerID,
		required: true,
	},
}

// queryServerAddr query strings
var queryServerAddrQueryStrings = []querystring{
	querystring{
		name:     qsQueryServerAddr,
		required: true,
	},
}

// getServers query strings
var getServersQueryStrings = []querystring{
	querystring{
		name: qsGetServersCountry,
	},
	querystring{
		name: qsGetServersRegion,
	},
	querystring{
		name: qsGetServersState,
	},
	querystring{
		name: qsGetServersName,
	},
	querystring{
		name: qsGetServersMap,
	},
	querystring{
		name: qsGetServersGameDir,
	},
	querystring{
		name: qsGetServersGame,
	},
	querystring{
		name: qsGetServersType,
	},
	querystring{
		name: qsGetServersOS,
	},
	querystring{
		name: qsGetServersVersion,
	},
	querystring{
		name: qsGetServersKeywords,
	},
	querystring{
		name:     qsGetServersHasPlayers,
		boolonly: true,
	},
	querystring{
		name:     qsGetServersHasBots,
		boolonly: true,
	},
	querystring{
		name:     qsGetServersHasPassword,
		boolonly: true,
	},
	querystring{
		name:     qsGetServersHasAntiCheat,
		boolonly: true,
	},
}

// getQStringValues takes the map returned by URL
func getQStringValues(m map[string][]string, querystring string) []string {
	var vals []string
	for k := range m {
		if strings.EqualFold(k, querystring) {
			vals = strings.Split(m[k][0], ",")
			break
		}
	}
	// case where there's no value after query string (i.e: ?querystring=)
	if vals[0] == "" {
		return nil
	}
	return vals
}
