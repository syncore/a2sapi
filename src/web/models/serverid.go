package models

// serverid.go - Server ID model

type Server struct {
	ID   int64  `json:"serverID"`
	Host string `json:"host"`
}

type ServerID struct {
	ServerCount int       `json:"serverCount"`
	Servers     []*Server `json:"servers"`
}

func GetDefaultServerID() *ServerID {
	return &ServerID{
		ServerCount: 0,
		Servers:     make([]*Server, 0),
	}
}
