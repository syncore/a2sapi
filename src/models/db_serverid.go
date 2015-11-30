package models

// db_serverid.go - Model for host, id, and game returned by server DB

type DbServer struct {
	ID   int64  `json:"serverID"`
	Game string `json:"game"`
	Host string `json:"host"`
}

type DbServerID struct {
	ServerCount int         `json:"serverCount"`
	Servers     []*DbServer `json:"servers"`
}

func GetDefaultServerID() *DbServerID {
	return &DbServerID{
		ServerCount: 0,
		Servers:     make([]*DbServer, 0),
	}
}
