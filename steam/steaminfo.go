// steaminfo.go - testing steam server query for info (A2S_INFO)

package steam

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"steamtest/util"
	"time"
)

type ServerInfo struct {
	Protocol    int        `json:"protocol"`
	Name        string     `json:"serverName"`
	Map         string     `json:"map"`
	Folder      string     `json:"gameDir"`
	Game        string     `json:"game"`
	ID          int16      `json:"steamApp"`
	Players     int16      `json:"players"`
	MaxPlayers  int16      `json:"maxPlayers"`
	Bots        int16      `json:"bots"`
	ServerType  string     `json:"serverType"`
	Environment string     `json:"serverOs"`
	Visibility  int16      `json:"private"`
	VAC         int16      `json:"antiCheat"`
	Version     string     `json:"serverVersion"`
	ExtraData   *extraData `json:"extra"`
}

type extraData struct {
	Port         int16  `json:"gamePort"`
	SteamID      uint64 `json:"serverSteamId"`
	SourceTVPort int16  `json:"sourceTvProxyPort"`
	SourceTVName string `json:"sourceTvProxyName"`
	Keywords     string `json:"keywords"`
	GameID       uint64 `json:"steamAppId"`
}

func getServerInfo(host string, timeout int) ([]byte, error) {
	request := []byte{
		0xFF, 0xFF, 0xFF, 0xFF,
		0x54, 0x53, 0x6F, 0x75, 0x72,
		0x63, 0x65, 0x20, 0x45, 0x6E,
		0x67, 0x69, 0x6E, 0x65, 0x20,
		0x51, 0x75, 0x65, 0x72, 0x79,
		0x00,
	}

	conn, err := net.DialTimeout("udp", host, time.Duration(timeout)*time.Second)
	if err != nil {
		return nil, HostConnectionError(err.Error())
	}
	defer conn.Close()
	conn.SetDeadline(time.Now().Add(time.Duration(timeout-1) * time.Second))

	_, err = conn.Write(request)
	if err != nil {
		return nil, DataTransmitError(err.Error())
	}

	var buf [maxPacketSize]byte
	numread, err := conn.Read(buf[:maxPacketSize])
	if err != nil {
		return nil, DataTransmitError(err.Error())
	}
	serverInfo := make([]byte, numread)
	copy(serverInfo, buf[:numread])

	if !bytes.HasPrefix(serverInfo, []byte{0xFF, 0xFF, 0xFF, 0xFF, 0x49}) {
		fmt.Printf("Server info response header is invalid\n")
		return nil, PacketHeaderError
	}

	//fmt.Printf("Server info reply from server: %x\n", serverInfo)

	return serverInfo, nil
}

func parseServerInfo(serverinfo []byte) (*ServerInfo, error) {
	if !bytes.HasPrefix(serverinfo, []byte{0xFF, 0xFF, 0xFF, 0xFF, 0x49}) {
		return nil, PacketHeaderError
	}

	serverinfo = bytes.TrimLeft(serverinfo, headerStr)

	// no info (should usually not happen)
	if len(serverinfo) <= 1 {
		return nil, NoInfoError
	}

	serverinfo = serverinfo[1:] // 0x49
	protocol := int(serverinfo[0])
	serverinfo = serverinfo[1:]

	name := util.ReadTillNul(serverinfo)
	serverinfo = serverinfo[len(name)+1:]
	mapname := util.ReadTillNul(serverinfo)
	serverinfo = serverinfo[len(mapname)+1:]
	folder := util.ReadTillNul(serverinfo)
	serverinfo = serverinfo[len(folder)+1:]
	game := util.ReadTillNul(serverinfo)
	serverinfo = serverinfo[len(game)+1:]
	id := int16(binary.LittleEndian.Uint16(serverinfo[:2]))
	serverinfo = serverinfo[2:]
	if id >= 2400 && id <= 2412 {
		return nil, fmt.Errorf("The Ship servers are not supported")
	}
	players := int16(serverinfo[0])
	serverinfo = serverinfo[1:]
	maxplayers := int16(serverinfo[0])
	serverinfo = serverinfo[1:]
	bots := int16(serverinfo[0])
	serverinfo = serverinfo[1:]
	servertype := string(serverinfo[0])
	serverinfo = serverinfo[1:]
	environment := string(serverinfo[0])
	serverinfo = serverinfo[1:]
	visibility := int16(serverinfo[0])
	serverinfo = serverinfo[1:]
	vac := int16(serverinfo[0])
	serverinfo = serverinfo[1:]
	version := util.ReadTillNul(serverinfo)
	serverinfo = serverinfo[len(version)+1:]

	// extra data flags
	var port int16
	var steamid uint64
	var sourcetvport int16
	var sourcetvname string
	var keywords string
	var gameid uint64
	edf := serverinfo[0]
	serverinfo = serverinfo[1:]
	if edf != 0x00 {
		if edf&0x80 > 0 {
			port = int16(binary.LittleEndian.Uint16(serverinfo[:2]))
			serverinfo = serverinfo[2:]
		}
		if edf&0x10 > 0 {
			steamid = binary.LittleEndian.Uint64(serverinfo[:8])
			serverinfo = serverinfo[8:]
		}
		if edf&0x40 > 0 {
			sourcetvport = int16(binary.LittleEndian.Uint16(serverinfo[:2]))
			serverinfo = serverinfo[2:]
			sourcetvname = util.ReadTillNul(serverinfo)
			serverinfo = serverinfo[len(sourcetvname)+1:]
		}
		if edf&0x20 > 0 {
			keywords = util.ReadTillNul(serverinfo)
			serverinfo = serverinfo[len(keywords)+1:]
		}
		if edf&0x01 > 0 {
			gameid = binary.LittleEndian.Uint64(serverinfo[:8])
			serverinfo = serverinfo[len(serverinfo):]
		}

	}

	// format a few ambiguous values
	if environment == "l" {
		environment = "Linux"
	}
	if environment == "w" {
		environment = "Windows"
	}
	if servertype == "d" {
		servertype = "dedicated"
	}
	if servertype == "l" {
		servertype = "listen"
	}

	return &ServerInfo{
		Protocol:    protocol,
		Name:        name,
		Map:         mapname,
		Folder:      folder,
		Game:        game,
		ID:          id,
		Players:     players,
		MaxPlayers:  maxplayers,
		Bots:        bots,
		ServerType:  servertype,
		Environment: environment,
		Visibility:  visibility,
		VAC:         vac,
		Version:     version,
		ExtraData: &extraData{
			Port:         port,
			SteamID:      steamid,
			SourceTVPort: sourcetvport,
			SourceTVName: sourcetvname,
			Keywords:     keywords,
			GameID:       gameid,
		},
	}, nil
}

func GetServerInfoWithLiveData(host string, timeout int) (*ServerInfo, error) {
	si, err := getServerInfo(host, timeout)
	if err != nil {
		return nil, err
	}
	//fmt.Printf("Server info is: %x\n", si)
	serverinfo, err := parseServerInfo(si)
	if err != nil {
		return nil, err
	}
	// fmt.Printf(`protocol:%d, name:%s, mapname:%s, folder:%s, game:%s, id:%d,
	//  players:%d, maxplayers:%d, bots:%d, servertype:%s, environment:%s,
	//  visibility:%d, vac:%d, version:%s, port:%d, steamid:%d, sourcetvport:%d,
	//  sourcetvname:%s, keywords:%s, gameid:%d`,
	// 	serverinfo.Protocol, serverinfo.Name, serverinfo.Map, serverinfo.Folder,
	// 	serverinfo.Game, serverinfo.ID, serverinfo.Players, serverinfo.MaxPlayers,
	// 	serverinfo.Bots, serverinfo.ServerType, serverinfo.Environment,
	// 	serverinfo.Visibility, serverinfo.VAC, serverinfo.Version,
	// 	serverinfo.ExtraData.Port, serverinfo.ExtraData.SteamID,
	// 	serverinfo.ExtraData.SourceTVPort, serverinfo.ExtraData.SourceTVName,
	// 	serverinfo.ExtraData.Keywords, serverinfo.ExtraData.GameID)

	return serverinfo, nil
}

func GetServerInfoWithTestData() error {
	si := testServerInfoData
	serverinfo, err := parseServerInfo(si)
	if err != nil {
		return err
	}
	fmt.Printf(`protocol:%d, name:%s, mapname:%s, folder:%s, game:%s, id:%d,
	 players:%d, maxplayers:%d, bots:%d, servertype:%s, environment:%s,
	 visibility:%d, vac:%d, version:%s, port:%d, steamid:%d, sourcetvport:%d,
	 sourcetvname:%s, keywords:%s, gameid:%d`,
		serverinfo.Protocol, serverinfo.Name, serverinfo.Map, serverinfo.Folder,
		serverinfo.Game, serverinfo.ID, serverinfo.Players, serverinfo.MaxPlayers,
		serverinfo.Bots, serverinfo.ServerType, serverinfo.Environment,
		serverinfo.Visibility, serverinfo.VAC, serverinfo.Version,
		serverinfo.ExtraData.Port, serverinfo.ExtraData.SteamID,
		serverinfo.ExtraData.SourceTVPort, serverinfo.ExtraData.SourceTVName,
		serverinfo.ExtraData.Keywords, serverinfo.ExtraData.GameID)
	return nil
}
