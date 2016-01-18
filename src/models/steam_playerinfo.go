package models

// steam_playerinfo.go - Model for player info returned by a steam A2S_PLAYER query

// SteamPlayerInfo represents a player returned by a Steam A2S_PLAYER query
type SteamPlayerInfo struct {
	Name              string  `json:"name"`
	Score             int32   `json:"score"`
	TimeConnectedSecs float32 `json:"secsConnected"`
	TimeConnectedTot  string  `json:"totalConnected"`
}

// FilteredPlayerInfo is a collection of all players on a server that actually
// exist on the server and are not bugged or stuck due to the Steam de-auth
// bug that exists in game servers for certain games (such as Quake Live)
type FilteredPlayerInfo struct {
	FilteredPlayerCount int               `json:"count"`
	FilteredPlayers     []SteamPlayerInfo `json:"players"`
}
