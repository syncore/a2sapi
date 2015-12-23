package models

// steam_serverinfo.go - Model for server info returned by an A2S_INFO query

// SteamServerInfo represents the original information returned by a direct
// A2S_INFO query of a given host.
type SteamServerInfo struct {
	Protocol    int             `json:"protocol,omitempty"`
	Name        string          `json:"serverName,omitempty"`
	Map         string          `json:"map,omitempty"`
	Folder      string          `json:"gameDir,omitempty"`
	Game        string          `json:"game,omitempty"`
	ID          int16           `json:"steamApp,omitempty"`
	Players     int16           `json:"players,omitempty"`
	MaxPlayers  int16           `json:"maxPlayers,omitempty"`
	Bots        int16           `json:"bots,omitempty"`
	ServerType  string          `json:"serverType,omitempty"`
	Environment string          `json:"serverOS,omitempty"`
	Visibility  int16           `json:"private,omitempty"`
	VAC         int16           `json:"antiCheat,omitempty"`
	Version     string          `json:"serverVersion,omitempty"`
	ExtraData   *SteamExtraData `json:"extra,omitempty"`
}

// SteamExtraData represents the original extra data field, if present returned
// by a direct A2S_INFO query of a given host.
type SteamExtraData struct {
	Port         int16  `json:"gamePort"`
	SteamID      uint64 `json:"serverSteamId"`
	SourceTVPort int16  `json:"sourceTvProxyPort"`
	SourceTVName string `json:"sourceTvProxyName"`
	Keywords     string `json:"keywords"`
	GameID       uint64 `json:"steamAppId"`
}
