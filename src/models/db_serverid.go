package models

// db_serverid.go - Model for host, id, and game returned by server DB

// DbServer represents an individual server's internal ID information.
type DbServer struct {
	ID   int64  `json:"serverID"`
	Game string `json:"game"`
	Host string `json:"host"`
}

// DbServerID represents the outer struct that is retrieved from the server ID
// database.
type DbServerID struct {
	ServerCount int         `json:"serverCount"`
	Servers     []*DbServer `json:"servers"`
}

// GetDefaultServerID returns the default DbServerID outer struct when a given
// host does not have an ID that was found in the server ID database.
func GetDefaultServerID() *DbServerID {
	return &DbServerID{
		ServerCount: 0,
		Servers:     make([]*DbServer, 0),
	}
}
