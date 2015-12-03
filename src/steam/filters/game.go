package filters

import "strings"

// game.go - Steam game-to-appid and A2S rule ignore mappings

// Game This struct represents a queryable Steam game, including its application ID
// and whether particular A2S requests need to be ignored when querying.
type Game struct {
	Name  string
	AppID uint64
	// Some games (i.e. newer/beta ones) do not have all 3 of A2S_INFO,PLAYER,RULES
	// any of these ignore values set to true will skip that request when querying
	IgnoreRules   bool
	IgnorePlayers bool
	IgnoreInfo    bool
}

type games []*Game

// A few games, additional games can be added from https://steamdb.info/apps/
var (
	// GameCsGo Counter-Strike: GO
	GameCsGo = &Game{
		Name:          "CSGO",
		AppID:         730,
		IgnoreRules:   false,
		IgnorePlayers: false,
		IgnoreInfo:    false,
	}
	// GameQuakeLive Quake Live
	GameQuakeLive = &Game{
		Name:          "QuakeLive",
		AppID:         282440,
		IgnoreRules:   false,
		IgnorePlayers: false,
		IgnoreInfo:    false,
	}
	// GameReflex Reflex
	GameReflex = &Game{
		Name:          "Reflex",
		AppID:         328070,
		IgnoreRules:   true, // Reflex does not implement A2S_RULES
		IgnorePlayers: false,
		IgnoreInfo:    false,
	}
	// GameTF2 Team Fortress 2
	GameTF2 = &Game{
		Name:          "TF2",
		AppID:         440,
		IgnoreRules:   false,
		IgnorePlayers: false,
		IgnoreInfo:    false,
	}
	// GameUnspecified Unspecified game for direct server queries, if enabled;
	// if unspecified games actually ignore some A2S requests there will be issues.
	GameUnspecified = &Game{
		Name:          "Unspecified",
		AppID:         0,
		IgnoreRules:   false,
		IgnorePlayers: false,
		IgnoreInfo:    false,
	}

	gamelist = games{
		GameCsGo,
		GameQuakeLive,
		GameReflex,
		GameTF2,
		GameUnspecified,
	}
)

func (g *Game) String() string {
	return g.Name
}

// GetGameByName searches the list of pre-defined games and returns a pointer to
// a Game struct based on the name of the game.
func GetGameByName(name string) *Game {
	for _, g := range gamelist {
		if strings.EqualFold(name, g.Name) {
			return g
		}
	}
	return GameUnspecified
}

// GetGameByAppID searches the list of pre-defined games and returns a pointer to
// a Game struct based on the AppID of the game.
func GetGameByAppID(appid uint64) *Game {
	for _, g := range gamelist {
		if appid == g.AppID {
			return g
		}
	}
	return GameUnspecified
}

// NewGame specifies a new game, including its name, Steam application-ID, and
// whether A2S_RULES, A2S_PLAYERS, and/or AS2_INFO requests should be ignored
// when performing a query.
func NewGame(name string, appid uint64, ignoreRules, ignorePlayers,
	ignoreInfo bool) *Game {
	return &Game{
		Name:          name,
		AppID:         appid,
		IgnoreRules:   ignoreRules,
		IgnorePlayers: ignorePlayers,
		IgnoreInfo:    ignoreInfo,
	}
}
