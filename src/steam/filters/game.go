package filters

// game.go - Steam game-to-appid operations, A2S ignore mappings, and game list.

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/syncore/a2sapi/src/constants"
	"github.com/syncore/a2sapi/src/util"
)

// Game represents a queryable Steam game, including its application ID
// and whether particular A2S requests need to be ignored when querying.
type Game struct {
	Name  string `json:"name"`
	AppID uint64 `json:"appID"`
	// Some games (i.e. newer/beta ones) do not have all 3 of A2S_INFO,PLAYER,RULES
	// any of these ignore values set to true will skip that request when querying
	IgnoreRules   bool `json:"ignoreRules"`
	IgnorePlayers bool `json:"ignorePlayers"`
	IgnoreInfo    bool `json:"ignoreInfo"`
}

// GameList represents the list of games.
type GameList struct {
	Games []Game `json:"games"`
}

// A few default games, additional games can be added from https://steamdb.info/apps/
var (
	// GameAlienSwarm Alien Swarm
	GameAlienSwarm = Game{
		Name:          "AlienSwarm",
		AppID:         630,
		IgnoreRules:   false,
		IgnorePlayers: false,
		IgnoreInfo:    false,
	}
	// GameARMA3 ARMA 3
	GameARMA3 = Game{
		Name:          "ARMA3",
		AppID:         107410,
		IgnoreRules:   false,
		IgnorePlayers: false,
		IgnoreInfo:    false,
	}
	// GameARKSurvivalEvolved ARK: Survival Evolved
	GameARKSurvivalEvolved = Game{
		Name:          "ARKSurvivalEvolved",
		AppID:         346110,
		IgnoreRules:   false,
		IgnorePlayers: false,
		IgnoreInfo:    false,
	}
	// GameCsGo Counter-Strike: GO
	GameCsGo = Game{
		Name:          "CSGO",
		AppID:         730,
		IgnoreRules:   true, // CSGO no longer sends rules as of 1.32.3.0 (02/21/14)
		IgnorePlayers: false,
		IgnoreInfo:    false,
	}
	// GameCSSource Counter-Strike: Source
	GameCSSource = Game{
		Name:          "CSSource",
		AppID:         240,
		IgnoreRules:   false,
		IgnorePlayers: false,
		IgnoreInfo:    false,
	}
	// GameDayZ DayZ
	GameDayZ = Game{
		Name:          "DayZ",
		AppID:         221100,
		IgnoreRules:   false,
		IgnorePlayers: false,
		IgnoreInfo:    false,
	}
	// GameGarrysMod Garry's Mod
	GameGarrysMod = Game{
		Name:          "GarrysMod",
		AppID:         4000,
		IgnoreRules:   false,
		IgnorePlayers: false,
		IgnoreInfo:    false,
	}
	// GameHL2DM Half-Life 2: Deathmatch
	GameHL2DM = Game{
		Name:          "HL2DM",
		AppID:         320,
		IgnoreRules:   false,
		IgnorePlayers: false,
		IgnoreInfo:    false,
	}
	// GameL4D2 Left 4 Dead 2
	GameL4D2 = Game{
		Name:          "L4D2",
		AppID:         550,
		IgnoreRules:   false,
		IgnorePlayers: false,
		IgnoreInfo:    false,
	}
	// GameOpposingForce Half-Life: Opposing Force
	GameOpposingForce = Game{
		Name:          "OpposingForce",
		AppID:         50,
		IgnoreRules:   false,
		IgnorePlayers: false,
		IgnoreInfo:    false,
	}
	// GameQuakeLive Quake Live
	GameQuakeLive = Game{
		Name:          "QuakeLive",
		AppID:         282440,
		IgnoreRules:   false,
		IgnorePlayers: false,
		IgnoreInfo:    false,
	}
	// GameReflex Reflex
	GameReflex = Game{
		Name:          "Reflex",
		AppID:         328070,
		IgnoreRules:   true, // Reflex does not implement A2S_RULES
		IgnorePlayers: false,
		IgnoreInfo:    false,
	}
	// GameRust Rust
	GameRust = Game{
		Name:          "Rust",
		AppID:         252490,
		IgnoreRules:   false,
		IgnorePlayers: false,
		IgnoreInfo:    false,
	}
	// GameTF2 Team Fortress 2
	GameTF2 = Game{
		Name:          "TF2",
		AppID:         440,
		IgnoreRules:   false,
		IgnorePlayers: false,
		IgnoreInfo:    false,
	}
	// GameUnspecified Unspecified game for direct server queries, if enabled;
	// if unspecified games actually ignore some A2S requests there will be issues.
	// This is intentionally left out of the defaultGames GameList struct so it
	//is not user-selectable in the configuration creation.
	GameUnspecified = Game{
		Name:          "Unspecified",
		AppID:         0,
		IgnoreRules:   false,
		IgnorePlayers: false,
		IgnoreInfo:    false,
	}

	defaultGames = GameList{
		Games: []Game{
			GameAlienSwarm,
			GameARMA3,
			GameARKSurvivalEvolved,
			GameCsGo,
			GameCSSource,
			GameDayZ,
			GameGarrysMod,
			GameHL2DM,
			GameL4D2,
			GameOpposingForce,
			GameQuakeLive,
			GameReflex,
			GameRust,
			GameTF2,
		},
	}
	highServerCountGames = GameList{
		Games: []Game{
			GameARKSurvivalEvolved,
			GameARMA3,
			GameCsGo,
			GameGarrysMod,
			GameL4D2,
			GameTF2,
		},
	}
)

func (g *Game) String() string {
	return g.Name
}

// GetGameNames returns a slice of strings containing the games' names.
func GetGameNames() []string {
	var names []string
	for _, g := range ReadGames() {
		names = append(names, g.Name)
	}
	return names
}

// GetGameByName searches the list of pre-defined games and returns a a Game
// struct based on the name of the game.
func GetGameByName(name string) Game {
	for _, g := range ReadGames() {
		if strings.EqualFold(name, g.Name) {
			return g
		}
	}
	return GameUnspecified
}

// GetGameByAppID searches the list of pre-defined games and returns a Game struct
// based on the AppID of the game.
func GetGameByAppID(appid uint64) Game {
	for _, g := range ReadGames() {
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
	ignoreInfo bool) Game {
	return Game{
		Name:          name,
		AppID:         appid,
		IgnoreRules:   ignoreRules,
		IgnorePlayers: ignorePlayers,
		IgnoreInfo:    ignoreInfo,
	}
}

// ReadGames reads the game file from disk and returns a slice to a pointer of
// Game structs if successful, otherwise panics.
func ReadGames() []Game {
	var f *os.File
	var err error
	f, err = os.Open(constants.GameFileFullPath)
	if err != nil {
		// try to create
		DumpDefaultGames()
		// re-open
		f, err = os.Open(constants.GameFileFullPath)
		if err != nil {
			panic(fmt.Sprintf("Error reading games file file: %s\n", err))
		}
	}
	defer f.Close()
	r := bufio.NewReader(f)
	d := json.NewDecoder(r)
	games := GameList{}
	if err := d.Decode(&games); err != nil {
		panic(fmt.Sprintf("Error decoding games file file: %s\n", err))
	}
	return games.Games
}

// DumpDefaultGames writes the default struct containing the default games to disk
// on success, otherwise panics.
func DumpDefaultGames() {
	if err := util.WriteJSONConfig(defaultGames, constants.ConfigDirectory,
		constants.GameFileFullPath); err != nil {
		panic(err)
	}
}

// IsValidGame determines whether the specified game exists within the list of
// games and returns true if it does, otherwise false.
func IsValidGame(name string) bool {
	for _, g := range ReadGames() {
		if strings.EqualFold(name, g.Name) {
			return true
		}
	}
	return false
}

// HasHighServerCount determines if the specified game is in the list of games
// that are known to return more than 6930 servers, which is the value at which
// Valve begins to throttle future responses from the master server.
func HasHighServerCount(name string) bool {
	for _, g := range highServerCountGames.Games {
		if strings.EqualFold(name, g.Name) {
			return true
		}
	}
	return false
}
