package models

// steam_playerinfo.go - Model for player info returned by a steam A2S_PLAYER query

// SteamPlayerInfo represents a player returned by a Steam A2S_PLAYER query
type SteamPlayerInfo struct {
	Name              string  `json:"name"`
	Score             int32   `json:"score"`
	TimeConnectedSecs float32 `json:"secsConnected"`
	TimeConnectedTot  string  `json:"totalConnected"`
}

// RealPlayerInfo is a collection of all players estimated to be 'real', that is, not a bot or bugged/stuck player
type RealPlayerInfo struct {
	RealPlayerCount int                `json:"count"`
	Players         []*SteamPlayerInfo `json:"players"`
}
