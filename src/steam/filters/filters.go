package filters

// filters.go - steam master server filters
// See: https://developer.valvesoftware.com/wiki/Master_Server_Query_Protocol

import (
	"bytes"
	"fmt"
)

// SrvRegion represents a Master server region code filter
type SrvRegion []byte

// SrvFilter represents a Master server filter
type SrvFilter []byte

// Filter is our internal wrapper for a specified game, and its Master server
// region code and Master server filters
type Filter struct {
	Game    *Game
	Region  SrvRegion
	Filters []SrvFilter
}

// Regions and filters
var (
	SrUsEastCoast  SrvRegion = []byte{0x00}
	SrUsWestCoast  SrvRegion = []byte{0x01}
	SrSouthAmerica SrvRegion = []byte{0x02}
	SrEurope       SrvRegion = []byte{0x03}
	SrAsia         SrvRegion = []byte{0x04}
	SrAustralia    SrvRegion = []byte{0x05}
	SrMiddleEast   SrvRegion = []byte{0x06}
	SrAfrica       SrvRegion = []byte{0x07}
	SrAll          SrvRegion = []byte{0xFF}

	// --------------------- "Constant" filters ---------------------
	// Dedicated servers
	SfDedicated SrvFilter = []byte("\\type\\d")
	// Servers using anti-cheat technology (VAC, but maybe others as well)
	SfSecure SrvFilter = []byte("\\secure\\1")
	// Servers running on a Linux platform
	SfLinux SrvFilter = []byte("\\linux\\1")
	// Servers that are not empty
	SfNotEmpty SrvFilter = []byte("\\empty\\1")
	// Servers that are not full
	SfNotFull SrvFilter = []byte("\\full\\1")
	// Servers that spectator proxies
	SfSpectatorProxy SrvFilter = []byte("\\proxy\\1")
	// Servers that are empty
	SfEmpty SrvFilter = []byte("\\noplayers\\1")
	// Servers that are whitelisted
	SfWhitelisted SrvFilter = []byte("\\white\\1")
	// Return only one server for each unique IP address matched
	SfOneUniquePerIP SrvFilter = []byte("\\collapse_addr_hash\\1")
	// ALL servers
	SfAll SrvFilter = []byte{0x00}

	// ----------------Filters the take variable input ----------------
	// \appid\[appid] - Servers that are running game [appid]
	AppIDFilter = func(val string) SrvFilter {
		return []byte(fmt.Sprintf("\\appid\\%s", val))
	}
	// \gameaddr\[ip]Return only servers on the specified IP address
	// (port supported and optional)
	GameAddrFilter = func(val string) SrvFilter {
		return []byte(fmt.Sprintf("\\gameaddr\\%s", val))
	}
	// \gamedata\[tag,...] - Servers with all of the given tag(s) in their
	//'hidden' tags (L4D2)
	GameDataFilter = func(val string) SrvFilter {
		return []byte(fmt.Sprintf("\\gamedata\\%s", val))
	}
	// \gamedataor\[tag,...] - Servers with any of the given tag(s) in their
	// 'hidden' tags (L4D2)
	GameDataOrFilter = func(val string) SrvFilter {
		return []byte(fmt.Sprintf("\\gamedataor\\%s", val))
	}
	// \gamedir\[mod] - Servers running the specified modification (ex. cstrike)
	GameDirFilter = func(val string) SrvFilter {
		return []byte(fmt.Sprintf("\\gamedir\\%s", val))
	}
	// \gametype\[tag,...] - Servers with all of the given tag(s) in sv_tags
	GameTypeFilter = func(val string) SrvFilter {
		return []byte(fmt.Sprintf("\\gametype\\%s", val))
	}
	// \name_match\[hostname] - Servers with their hostname matching [hostname]
	// (can use * as a wildcard)
	NameMatchFilter = func(val string) SrvFilter {
		return []byte(fmt.Sprintf("\\name_match\\%s", val))
	}
	// \nand\[x] - A special filter, specifies that servers matching all of the
	// following [x] conditions should not be returned
	NAndFilter = func(val string) SrvFilter {
		return []byte(fmt.Sprintf("\\nand\\%s", val))
	}
	// \nor\[x] - A special filter, specifies that servers matching any of the
	//following [x] conditions should not be returned
	NOrFilter = func(val string) SrvFilter {
		return []byte(fmt.Sprintf("\\nor\\%s", val))
	}
	// \napp\[appid] - Servers that are NOT running game [appid]
	// (This was introduced to block Left 4 Dead games from the Steam Server Browser
	NAppIDFilter = func(val string) SrvFilter {
		return []byte(fmt.Sprintf("\\nappid\\%s", val))
	}
	// \map\[map] - Servers running the specified map (ex. cs_italy)
	MapFilter = func(val string) SrvFilter {
		return []byte(fmt.Sprintf("\\map\\%s", val))
	}
	// \version_match\[version] - Servers running version [version]
	// (can use * as a wildcard)
	VersionMatchFilter = func(val string) SrvFilter {
		return []byte(fmt.Sprintf("\\version_match\\%s", val))
	}
)

// NewFilter creates a new filter for use with a master server query based on
// a game to query, its region code, and any other additional master server filters
// that should be sent with the request to the master server.
func NewFilter(game *Game, region SrvRegion, filters []SrvFilter) *Filter {
	if filters != nil {
		for i, f := range filters {
			if bytes.HasPrefix(f, []byte("\\appid\\")) {
				filters[i] = AppIDFilter(game.AppID)
				break
			} else {
				filters = append(filters, AppIDFilter(game.AppID))
				break
			}
		}
	} else {
		filters = append(filters, AppIDFilter(game.AppID))
	}
	return &Filter{
		Game:    game,
		Region:  region,
		Filters: filters,
	}
}
