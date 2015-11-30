package steam

// steaminfo.go - steam server query for info (A2S_INFO)

import (
	"bytes"
	"encoding/binary"
	"net"
	"steamtest/src/models"
	"steamtest/src/util"
	"sync"
	"time"
)

func getServerInfo(host string, timeout int) ([]byte, error) {
	conn, err := net.DialTimeout("udp", host, time.Duration(timeout)*time.Second)
	if err != nil {
		// TODO: simplify this log + return into just a return
		util.LogSteamError(ErrHostConnection(err.Error()))
		return nil, ErrHostConnection(err.Error())
	}
	defer conn.Close()
	conn.SetDeadline(time.Now().Add(time.Duration(timeout-1) * time.Second))

	_, err = conn.Write(infoChallengeReq)
	if err != nil {
		// TODO: simplify this log + return into just a return
		util.LogSteamError(ErrDataTransmit(err.Error()))
		return nil, ErrDataTransmit(err.Error())
	}

	var buf [maxPacketSize]byte
	numread, err := conn.Read(buf[:maxPacketSize])
	if err != nil {
		// TODO: simplify this log + return into just a return
		util.LogSteamError(ErrDataTransmit(err.Error()))
		return nil, ErrDataTransmit(err.Error())
	}
	serverInfo := make([]byte, numread)
	copy(serverInfo, buf[:numread])

	if !bytes.HasPrefix(serverInfo, expectedInfoRespHeader) {
		// TODO: simplify this log + return into just a return
		util.LogSteamError(ErrPacketHeader)
		return nil, ErrPacketHeader
	}

	return serverInfo, nil
}

func parseServerInfo(serverinfo []byte) (*models.SteamServerInfo, error) {
	if !bytes.HasPrefix(serverinfo, expectedInfoRespHeader) {
		// TODO: simplify this log + return into just a return
		util.LogSteamError(ErrPacketHeader)
		return nil, ErrPacketHeader
	}

	serverinfo = bytes.TrimLeft(serverinfo, headerStr)

	// no info (should usually not happen)
	if len(serverinfo) <= 1 {
		util.LogSteamError(ErrNoInfo)
		return nil, ErrNoInfo
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
		return nil, util.LogSteamErrorf("The Ship servers are not supported")
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

	return &models.SteamServerInfo{
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
		ExtraData: &models.SteamExtraData{
			Port:         port,
			SteamID:      steamid,
			SourceTVPort: sourcetvport,
			SourceTVName: sourcetvname,
			Keywords:     keywords,
			GameID:       gameid,
		},
	}, nil
}

func RetryFailedInfoReq(failed []string,
	retrycount int) map[string]*models.SteamServerInfo {
	m := make(map[string]*models.SteamServerInfo)
	var f []string
	var wg sync.WaitGroup
	var mut sync.Mutex
	for i := 0; i < retrycount; i++ {
		if i == 0 {
			f = failed
		}
		wg.Add(len(f))
		for _, host := range f {
			go func(h string) {
				defer wg.Done()
				r, err := GetInfoForServer(h, QueryTimeout)
				if err != nil {
					if err != ErrNoPlayers {
						return
					}
				}
				mut.Lock()
				m[h] = r
				f = removeFailedHost(f, h)
				mut.Unlock()
			}(host)
		}
		wg.Wait()
	}
	return m
}

func GetInfoForServer(host string, timeout int) (*models.SteamServerInfo, error) {
	// Caller will log. Return err instead of wrapped util.LogSteamError so as not
	// to interfere with custom error types that need to be analyzed when
	// determining if retry needs to be done.
	si, err := getServerInfo(host, timeout)
	if err != nil {
		return nil, err
	}

	serverinfo, err := parseServerInfo(si)
	if err != nil {
		return nil, err
	}
	return serverinfo, nil
}
