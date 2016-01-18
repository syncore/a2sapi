package steam

import (
	"a2sapi/src/models"
	"a2sapi/src/steam/filters"
	"a2sapi/src/test"
	"strings"
	"testing"
)

var testData a2sData

func init() {
	test.SetupEnvironment()
	hostsgames := make(map[string]filters.Game, 2)
	hostsgames["54.172.5.67:25801"] = filters.GameReflex
	hostsgames["192.211.62.11:27960"] = filters.GameQuakeLive

	info := make(map[string]models.SteamServerInfo, 2)
	info["54.172.5.67:25801"] = models.SteamServerInfo{
		Protocol:    17,
		Name:        "TurboPixel Appreciation Society (Official) #1",
		Map:         "xfdm2",
		Folder:      "base",
		Game:        "Reflex",
		Players:     6,
		MaxPlayers:  8,
		Bots:        0,
		ServerType:  "dedicated",
		Environment: "Windows",
		VAC:         1,
		Version:     "0.38.2",
		ExtraData: models.SteamExtraData{
			Port:         25800,
			SteamID:      90098615517053960,
			SourceTVPort: 0,
			SourceTVName: "",
			Keywords:     "atdm||62|1",
			GameID:       328070,
		},
	}
	info["192.211.62.11:27960"] = models.SteamServerInfo{
		Protocol:    17,
		Name:        "exile.syncore.org | US-Central #1 | Competitive",
		Map:         "overkill",
		Folder:      "baseq3",
		Game:        "Clan Arena",
		Players:     0,
		MaxPlayers:  16,
		Bots:        0,
		ServerType:  "dedicated",
		Environment: "Linux",
		VAC:         1,
		Version:     "1066",
		ExtraData: models.SteamExtraData{
			Port:         27960,
			SteamID:      90098677041473542,
			SourceTVPort: 0,
			SourceTVName: "",
			Keywords:     "clanarena,minqlx,syncore,texas,central,newmaps",
			GameID:       282440,
		},
	}
	rules := make(map[string]map[string]string, 2)
	rules["54.172.5.67:25801"] = nil
	r := make(map[string]string, 42)
	r["dmflags"] = "28"
	r["fraglimit"] = "50"
	r["g_adCaptureScoreBonus"] = "3"
	r["g_adElimScoreBonus"] = "2"
	r["g_adTouchScoreBonus"] = "1"
	r["g_blueScore"] = ""
	r["g_customSettings"] = "0"
	r["g_factory"] = "ca"
	r["g_factoryTitle"] = "Clan Arena"
	r["g_freezeRoundDelay"] = "4000"
	r["g_gameState"] = "PRE_GAME"
	r["g_gametype"] = "4"
	r["g_gravity"] = "800"
	r["g_instaGib"] = "0"
	r["g_itemHeight"] = "35"
	r["g_itemTimers"] = "1"
	r["g_levelStartTime"] = "1451179049"
	r["g_loadout"] = "0"
	r["g_needpass"] = "0"
	r["g_overtime"] = "0"
	r["g_quadDamageFactor"] = "3"
	r["g_redScore"] = ""
	r["g_roundWarmupDelay"] = "10000"
	r["g_startingHealth"] = "200"
	r["g_teamForceBalance"] = "1"
	r["g_teamSizeMin"] = "1"
	r["g_timeoutCount"] = "0"
	r["g_voteFlags"] = "0"
	r["g_weaponRespawn"] = "5"
	r["mapname"] = "overkill"
	r["mercylimit"] = "0"
	r["protocol"] = "91"
	r["roundlimit"] = "10"
	r["roundtimelimit"] = "180"
	r["scorelimit"] = "150"
	r["sv_hostname"] = "exile.syncore.org | US-Central #1 | Competitive"
	r["sv_maxclients"] = "16"
	r["sv_privateClients"] = "0"
	r["teamsize"] = "4"
	r["timelimit"] = "0"
	r["version"] = "1066 linux-x64 Dec 17 2015 15:36:49"
	rules["192.211.62.11:27960"] = r

	players := make(map[string][]models.SteamPlayerInfo, 2)
	players["54.172.5.67:25801"] = []models.SteamPlayerInfo{
		models.SteamPlayerInfo{
			Name:              "KovaaK",
			Score:             92,
			TimeConnectedSecs: 4317.216,
			TimeConnectedTot:  "1h11m57s",
		},
		models.SteamPlayerInfo{
			Name:              "Sharqosity",
			Score:             42,
			TimeConnectedSecs: 3428.6987,
			TimeConnectedTot:  "57m8s",
		},
		models.SteamPlayerInfo{
			Name:              "dhaK",
			Score:             42,
			TimeConnectedSecs: 1730.0668,
			TimeConnectedTot:  "28m50s",
		},
		models.SteamPlayerInfo{
			Name:              "yoo",
			Score:             45,
			TimeConnectedSecs: 467.6571,
			TimeConnectedTot:  "7m47s",
		},
		models.SteamPlayerInfo{
			Name:              "twitch.tv/liveanton - SANE",
			Score:             75,
			TimeConnectedSecs: 452.20792,
			TimeConnectedTot:  "7m32s",
		},
		models.SteamPlayerInfo{
			Name:              "ObviouslyBuggedPlayer",
			Score:             0,
			TimeConnectedSecs: 24120.2000,
			TimeConnectedTot:  "6h42m2s",
		},
	}
	players["192.211.62.11:27960"] = nil
	testData = a2sData{
		HostsGames: hostsgames,
		Info:       info,
		Rules:      rules,
		Players:    players,
	}
}

func TestBuildServerList(t *testing.T) {
	asl, err := buildServerList(testData, false)
	if err != nil {
		t.Fatalf("Unexpected error occurred when building server list.")
	}
	if len(asl.Servers) != 2 {
		t.Fatalf("Expected 2 servers, got: %d", len(asl.Servers))
	}
	// Slice not guaranteed to be in order
	var reflexServer models.APIServer
	var qlServer models.APIServer
	if asl.Servers[0].Info.ExtraData.GameID == 282440 {
		qlServer = asl.Servers[0]
		reflexServer = asl.Servers[1]
	} else {
		qlServer = asl.Servers[1]
		reflexServer = asl.Servers[0]
	}
	if reflexServer.Info.Players != 6 {
		t.Fatalf("Expected Reflex server to contain 6 players, got: %d",
			reflexServer.Info.Players)
	}
	if qlServer.Info.ExtraData.GameID != 282440 {
		t.Fatalf("Expected Quake Live server's steam game id to be 282440, got: %d",
			qlServer.Info.ExtraData.GameID)
	}
	if len(qlServer.Players) != 0 {
		t.Fatalf("Expected Quake Live server to contain no players, got: %v",
			len(qlServer.Players))
	}
}

func TestRemoveBuggedPlayers(t *testing.T) {
	buggedRemoved := removeBuggedPlayers(testData.Players["54.172.5.67:25801"])
	if len(buggedRemoved.FilteredPlayers) != 5 {
		t.Fatalf("Expected 5 players after bugged player removal, got: %d",
			len(buggedRemoved.FilteredPlayers))
	}
	for _, player := range buggedRemoved.FilteredPlayers {
		if strings.EqualFold(player.Name, "ObviouslyBuggedPlayer") {
			t.Fatalf("Filtered player list should not contain bugged player")
		}
	}
}
