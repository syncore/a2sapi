package constants

// test_constants.go - Test-related constants (and a few variables)

import "path"

const (
	// TestTempDirectory specifies the temporary directory for test-related files.
	TestTempDirectory = "test_temp"
	// TestConfigFilename specifies the name of the test configuration file.
	TestConfigFilename = "test.conf"
	// TestServerDbFilename specifies the name of the server database file used in
	// tests.
	TestServerDbFilename = "servers_test.sqlite"
)

var (
	// TestConfigFilePath represents the OS-independent full path to the config file.
	TestConfigFilePath = path.Join(TestTempDirectory, TestConfigFilename)
	// TestServerDbFilePath represents the OS-independent full path to the server
	// database file used in tests.
	TestServerDbFilePath = path.Join(TestTempDirectory, TestServerDbFilename)

	// TestServerDumpJSON is the JSON used for the server dump when performing
	// tests.
	TestServerDumpJSON = []byte(`
{"retrievalDate":"Sat Dec 26 23:08:14 2015 EST","timestamp":1451189294,
"serverCount":3,"servers":[{"serverId":1029,"address":"54.172.5.67:25801",
"ip":"54.172.5.67","port":25801,"location":{"countryName":"United States",
"countryCode":"US","region":"North America","state":"VA"},
"info":{"protocol":17,"serverName":"TurboPixel Appreciation Society (Official) #1",
"map":"xfdm2","gameDir":"base","game":"Reflex","players":5,"maxPlayers":8,
"serverType":"dedicated","serverOS":"Windows","antiCheat":1,"serverVersion":"0.38.2",
"extra":{"gamePort":25800,"serverSteamId":90098615517053960,"sourceTvProxyPort":0,
"sourceTvProxyName":"","keywords":"atdm||62|1","steamAppId":328070}},
"players":[{"name":"KovaaK","score":92,"secsConnected":4317.216,
"totalConnected":"1h11m57s"},{"name":"Sharqosity","score":42,
"secsConnected":3428.6987,"totalConnected":"57m8s"},{"name":"dhaK","score":42,
"secsConnected":1730.0668,
"totalConnected":"28m50s"},{"name":"yoo","score":45,"secsConnected":467.6571,
"totalConnected":"7m47s"},{"name":"twitch.tv/liveanton - SANE","score":75,
"secsConnected":452.20792,"totalConnected":"7m32s"}],
"realPlayers":{"count":5,"players":[{"name":"KovaaK","score":92,
"secsConnected":4317.216,"totalConnected":"1h11m57s"},{"name":"Sharqosity",
"score":42,"secsConnected":3428.6987,"totalConnected":"57m8s"},{"name":"dhaK",
"score":42,"secsConnected":1730.0668,"totalConnected":"28m50s"},{"name":"yoo",
"score":45,"secsConnected":467.6571,"totalConnected":"7m47s"},
{"name":"twitch.tv/liveanton - SANE","score":75,"secsConnected":452.20792,
"totalConnected":"7m32s"}]},"rules":{}},{"serverId":360,"address":"192.211.62.11:27960",
"ip":"192.211.62.11","port":27960,"location":{"countryName":"United States",
"countryCode":"US","region":"North America","state":"TX"},"info":{"protocol":17,
"serverName":"exile.syncore.org | US-Central #1 | Competitive","map":"overkill",
"gameDir":"baseq3","game":"Clan Arena","maxPlayers":16,"serverType":"dedicated",
"serverOS":"Linux","antiCheat":1,"serverVersion":"1066","extra":{"gamePort":27960,
"serverSteamId":90098677041473542,"sourceTvProxyPort":0,"sourceTvProxyName":"",
"keywords":"clanarena,minqlx,syncore,texas,central,newmaps","steamAppId":282440}},
"players":[],"realPlayers":{"count":0,"players":[]},"rules":{"capturelimit":"8",
"dmflags":"28","fraglimit":"50","g_adCaptureScoreBonus":"3","g_adElimScoreBonus":"2",
"g_adTouchScoreBonus":"1","g_blueScore":"","g_customSettings":"0","g_factory":"ca",
"g_factoryTitle":"Clan Arena","g_freezeRoundDelay":"4000","g_gameState":"PRE_GAME",
"g_gametype":"4","g_gravity":"800",
"g_instaGib":"0","g_itemHeight":"35","g_itemTimers":"1","g_levelStartTime":"1451179049",
"g_loadout":"0","g_needpass":"0","g_overtime":"0","g_quadDamageFactor":"3","g_redScore":"",
"g_roundWarmupDelay":"10000","g_startingHealth":"200","g_teamForceBalance":"1",
"g_teamSizeMin":"1","g_timeoutCount":"0","g_voteFlags":"0","g_weaponRespawn":"5",
"mapname":"overkill","mercylimit":"0","protocol":"91","roundlimit":"10",
"roundtimelimit":"180","scorelimit":"150",
"sv_hostname":"exile.syncore.org | US-Central #1 | Competitive",
"sv_maxclients":"16","sv_privateClients":"0","teamsize":"4","timelimit":"0",
"version":"1066 linux-x64 Dec 17 2015 15:36:49"}},{"serverId":746,
"address":"45.55.168.160:27960","ip":"45.55.168.160","port":27960,
"location":{"countryName":"United States","countryCode":"US","region":"North America",
"state":"NY"},"info":{"protocol":17,
"serverName":"triton.syncore.org | US-East #1 | Competitive","map":"overkill",
"gameDir":"baseq3","game":"Clan Arena","maxPlayers":16,"serverType":"dedicated",
"serverOS":"Linux","antiCheat":1,"serverVersion":"1066","extra":{"gamePort":27960,
"serverSteamId":90098677079644165,
"sourceTvProxyPort":0,"sourceTvProxyName":"",
"keywords":"clanarena,minqlx,syncore,newyork,east,newmaps","steamAppId":282440}},
"players":[],"realPlayers":{"count":0,"players":[]},"rules":{"capturelimit":"8",
"dmflags":"28","fraglimit":"50","g_adCaptureScoreBonus":"3","g_adElimScoreBonus":"2",
"g_adTouchScoreBonus":"1","g_blueScore":"","g_customSettings":"0","g_factory":"ca",
"g_factoryTitle":"Clan Arena","g_freezeRoundDelay":"4000","g_gameState":"PRE_GAME",
"g_gametype":"4",
"g_gravity":"800","g_instaGib":"0","g_itemHeight":"35","g_itemTimers":"1",
"g_levelStartTime":"1451179064","g_loadout":"0","g_needpass":"0","g_overtime":"0",
"g_quadDamageFactor":"3","g_redScore":"","g_roundWarmupDelay":"10000",
"g_startingHealth":"200","g_teamForceBalance":"1","g_teamSizeMin":"1",
"g_timeoutCount":"0","g_voteFlags":"0","g_weaponRespawn":"5","mapname":"overkill",
"mercylimit":"0","protocol":"91","roundlimit":"10","roundtimelimit":"180",
"scorelimit":"150","sv_hostname":"triton.syncore.org | US-East #1 | Competitive",
"sv_maxclients":"16","sv_privateClients":"0","teamsize":"4","timelimit":"0",
"version":"1066 linux-x64 Dec 17 2015 15:36:49"}}],"failedCount":0,"failedServers":[]}
`)
)
