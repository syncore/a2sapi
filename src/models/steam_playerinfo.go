package models

// steam_playerinfo.go - Model for player info returned by a steam A2S_PLAYER query

type SteamPlayerInfo struct {
	Name              string  `json:"name"`
	Score             int32   `json:"score"`
	TimeConnectedSecs float32 `json:"secsConnected"`
	TimeConnectedTot  string  `json:"totalConnected"`
}
