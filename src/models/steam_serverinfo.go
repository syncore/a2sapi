package models

// steam_serverinfo.go - Model for server info returned by an A2S_INFO query

// SteamServerInfo represents the original information returned by a direct
// A2S_INFO query of a given host.
type SteamServerInfo struct {
	Protocol      int             `json:"protocol"`
	Name          string          `json:"serverName"`
	Map           string          `json:"map"`
	Folder        string          `json:"gameDir"`
	Game          string          `json:"game"`
	GameTypeShort string          `json:"gameTypeShort"` // custom field for sorting
	GameTypeFull  string          `json:"gameTypeFull"`  // custom field for sorting
	ID            int16           `json:"steamApp"`
	Players       int16           `json:"players"`
	MaxPlayers    int16           `json:"maxPlayers"`
	Bots          int16           `json:"bots"`
	ServerType    string          `json:"serverType"`
	Environment   string          `json:"serverOS"`
	Visibility    int16           `json:"private"`
	VAC           int16           `json:"antiCheat"`
	Version       string          `json:"serverVersion"`
	ExtraData     *SteamExtraData `json:"extra"`
}

// SteamExtraData represents the original extra data field, if present returned
// by a direct A2S_INFO query of a given host.
type SteamExtraData struct {
	Port         int16  `json:"gamePort"`
	SteamID      uint64 `json:"serverSteamID"`
	SourceTVPort int16  `json:"sourceTvProxyPort"`
	SourceTVName string `json:"sourceTvProxyName"`
	Keywords     string `json:"keywords"`
	GameID       uint64 `json:"steamAppID"`
}
