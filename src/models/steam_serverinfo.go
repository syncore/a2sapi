package models

// steam_serverinfo.go - Model for server info returned by an A2S_INFO query

type SteamServerInfo struct {
	Protocol    int             `json:"protocol"`
	Name        string          `json:"serverName"`
	Map         string          `json:"map"`
	Folder      string          `json:"gameDir"`
	Game        string          `json:"game"`
	ID          int16           `json:"steamApp"`
	Players     int16           `json:"players"`
	MaxPlayers  int16           `json:"maxPlayers"`
	Bots        int16           `json:"bots"`
	ServerType  string          `json:"serverType"`
	Environment string          `json:"serverOs"`
	Visibility  int16           `json:"private"`
	VAC         int16           `json:"antiCheat"`
	Version     string          `json:"serverVersion"`
	ExtraData   *SteamExtraData `json:"extra"`
}

type SteamExtraData struct {
	Port         int16  `json:"gamePort"`
	SteamID      uint64 `json:"serverSteamId"`
	SourceTVPort int16  `json:"sourceTvProxyPort"`
	SourceTVName string `json:"sourceTvProxyName"`
	Keywords     string `json:"keywords"`
	GameID       uint64 `json:"steamAppId"`
}
