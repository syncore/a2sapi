package steam

import (
	"strings"

	"github.com/syncore/a2sapi/src/models"
	"github.com/syncore/a2sapi/src/steam/filters"
)

type gtype struct {
	ShortName string
	LongName  string
}

var qlGameTypes = map[string]gtype{
	"0":  {"FFA", "Free For All"},
	"1":  {"Duel", "Duel"},
	"2":  {"Race", "Race"},
	"3":  {"TDM", "Team Deathmatch"},
	"4":  {"CA", "Clan Arena"},
	"5":  {"CTF", "Capture The Flag"},
	"6":  {"FCTF", "1-Flag Capture The Flag"},
	"8":  {"HAR", "Harvester"},
	"9":  {"FT", "Freeze Tag"},
	"10": {"DOM", "Domination"},
	"11": {"AD", "Attack & Defend"},
	"12": {"RR", "Red Rover"},
}

var reflexGameTypes = map[string]gtype{
	"1v1":  {"1v1", "Duel"},
	"a1v1": {"a1v1", "Arena Duel"},
	"affa": {"affa", "Arena Free For All"},
	"atdm": {"atdm", "Arena Team Deathmatch"},
	"ctf":  {"ctf", "Capture The Flag"},
	"ffa":  {"ffa", "Free For All"},
	"race": {"race", "Race"},
	"tdm":  {"tdm", "Team Deathmatch"},
}

func getGameType(game filters.Game, server models.APIServer) (shortname,
	longname string) {
	// Quake Live
	if strings.EqualFold(game.Name, filters.GameQuakeLive.Name) {
		if _, ok := server.Rules["g_gametype"]; !ok {
			return
		}
		if _, ok := qlGameTypes[server.Rules["g_gametype"]]; !ok {
			return
		}
		return qlGameTypes[server.Rules["g_gametype"]].ShortName,
			qlGameTypes[server.Rules["g_gametype"]].LongName
	}
	// Reflex
	if strings.EqualFold(game.Name, filters.GameReflex.Name) {
		k := strings.Split(server.Info.ExtraData.Keywords, ",")
		if _, ok := reflexGameTypes[strings.ToLower(k[0])]; !ok {
			return
		}
		return reflexGameTypes[k[0]].ShortName, reflexGameTypes[k[0]].LongName
	}
	return
}
