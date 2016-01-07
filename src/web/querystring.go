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
	// gametype=
	qsGetServersGameType = "gametype"
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
	// ?isNotFull= (bool)
	qsGetServersIsNotFull = "isNotFull"
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
		name: qsGetServersGameType,
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
	querystring{
		name:     qsGetServersIsNotFull,
		boolonly: true,
	},
}

// getQStringValues takes the map returned by a *http.Request URL.Query(),
// extracts and returns the values of a key defined in that map which is
// specified as a known querystring value to match.
func getQStringValues(m map[string][]string, querystring string) []string {
	var vals []string
	for k := range m {
		if strings.EqualFold(k, querystring) {
			vals = strings.Split(m[k][0], ",")
			break
		}
	}
	if vals == nil {
		return nil
	}
	// case where there's no value after query string (i.e: ?querystring=)
	if vals[0] == "" {
		return nil
	}
	return vals
}
