// filters.go - steam master server filters
package filters

import "fmt"

type ServerRegion []byte
type ServerFilter []byte

var (
	// Regions
	SrUsEastCoast  ServerRegion = []byte{0x00}
	SrUsWestCoast  ServerRegion = []byte{0x01}
	SrSouthAmerica ServerRegion = []byte{0x02}
	SrEurope       ServerRegion = []byte{0x03}
	SrAsia         ServerRegion = []byte{0x04}
	SrAustralia    ServerRegion = []byte{0x05}
	SrMiddleEast   ServerRegion = []byte{0x06}
	SrAfrica       ServerRegion = []byte{0x07}
	SrAll          ServerRegion = []byte{0xFF}

	// --------------------- "Constant" filters ---------------------
	// Dedicated servers
	SfDedicated ServerFilter = []byte("\\type\\d")
	// Servers using anti-cheat technology (VAC, but maybe others as well)
	SfSecure ServerFilter = []byte("\\secure\\1")
	// Servers running on a Linux platform
	SfLinux ServerFilter = []byte("\\linux\\1")
	// Servers that are not empty
	SfNotEmpty ServerFilter = []byte("\\empty\\1")
	// Servers that are not full
	SfNotFull ServerFilter = []byte("\\full\\1")
	// Servers that spectator proxies
	SfSpectatorProxy ServerFilter = []byte("\\proxy\\1")
	// Servers that are empty
	SfEmpty ServerFilter = []byte("\\noplayers\\1")
	// Servers that are whitelisted
	SfWhitelisted ServerFilter = []byte("\\white\\1")
	// Return only one server for each unique IP address matched
	SfOneUniquePerIP ServerFilter = []byte("\\collapse_addr_hash\\1")
	// ALL servers
	SfAll ServerFilter = []byte{0x00}

	// ----------------Filters the take variable input ----------------
	// \appid\[appid] - Servers that are running game [appid]
	AppIdFilter = func(val string) ServerFilter {
		return []byte(fmt.Sprintf("\\appid\\%s", val))
	}
	// \gameaddr\[ip]Return only servers on the specified IP address
	// (port supported and optional)
	GameAddrFilter = func(val string) ServerFilter {
		return []byte(fmt.Sprintf("\\gameaddr\\%s", val))
	}
	// \gamedata\[tag,...] - Servers with all of the given tag(s) in their
	//'hidden' tags (L4D2)
	GameDataFilter = func(val string) ServerFilter {
		return []byte(fmt.Sprintf("\\gamedata\\%s", val))
	}
	// \gamedataor\[tag,...] - Servers with any of the given tag(s) in their
	// 'hidden' tags (L4D2)
	GameDataOrFilter = func(val string) ServerFilter {
		return []byte(fmt.Sprintf("\\gamedataor\\%s", val))
	}
	// \gamedir\[mod] - Servers running the specified modification (ex. cstrike)
	GameDirFilter = func(val string) ServerFilter {
		return []byte(fmt.Sprintf("\\gamedir\\%s", val))
	}
	// \gametype\[tag,...] - Servers with all of the given tag(s) in sv_tags
	GameTypeFilter = func(val string) ServerFilter {
		return []byte(fmt.Sprintf("\\gametype\\%s", val))
	}
	// \name_match\[hostname] - Servers with their hostname matching [hostname]
	// (can use * as a wildcard)
	NameMatchFilter = func(val string) ServerFilter {
		return []byte(fmt.Sprintf("\\name_match\\%s", val))
	}
	// \nand\[x] - A special filter, specifies that servers matching all of the
	// following [x] conditions should not be returned
	NAndFilter = func(val string) ServerFilter {
		return []byte(fmt.Sprintf("\\nand\\%s", val))
	}
	// \nor\[x] - A special filter, specifies that servers matching any of the
	//following [x] conditions should not be returned
	NOrFilter = func(val string) ServerFilter {
		return []byte(fmt.Sprintf("\\nor\\%s", val))
	}
	// \napp\[appid] - Servers that are NOT running game [appid]
	// (This was introduced to block Left 4 Dead games from the Steam Server Browser
	NAppIdFilter = func(val string) ServerFilter {
		return []byte(fmt.Sprintf("\\nappid\\%s", val))
	}
	// \map\[map] - Servers running the specified map (ex. cs_italy)
	MapFilter = func(val string) ServerFilter {
		return []byte(fmt.Sprintf("\\map\\%s", val))
	}
	// \version_match\[version] - Servers running version [version]
	// (can use * as a wildcard)
	VersionMatchFilter = func(val string) ServerFilter {
		return []byte(fmt.Sprintf("\\version_match\\%s", val))
	}

	// --------------------- A few games ---------------------
	// Additional "Source Engine Games" can be added from:
	// https://developer.valvesoftware.com/wiki/Steam_Application_IDs
	GameCsGo      = AppIdFilter("730")
	GameQuakeLive = AppIdFilter("282440")
	GameReflex    = AppIdFilter("328070")
	GameTF2       = AppIdFilter("440")
)
