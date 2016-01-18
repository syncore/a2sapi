package steam

// steammaster.go - testing steam master server query

import (
	"a2sapi/src/config"
	"a2sapi/src/logger"
	"a2sapi/src/steam/filters"
	"bytes"
	"fmt"
	"net"
	"time"
)

// MasterQuery contains the servers returned by a query to the Steam master server.
type MasterQuery struct {
	Servers []string
}

const masterServerHost = "hl2master.steampowered.com:27011"

func getServers(filter filters.Filter) ([]string, error) {
	maxHosts := config.Config.SteamConfig.MaximumHostsToReceive
	var serverlist []string
	var c net.Conn
	var err error
	retrieved := 0
	addr := "0.0.0.0:0"

	c, err = net.DialTimeout("udp", masterServerHost,
		time.Duration(QueryTimeout)*time.Second)
	if err != nil {
		logger.LogSteamError(ErrHostConnection(err.Error()))
		return nil, ErrHostConnection(err.Error())
	}

	defer c.Close()
	c.SetDeadline(time.Now().Add(time.Duration(QueryTimeout) * time.Second))

	for {
		s, err := queryMasterServer(c, addr, filter)
		if err != nil {
			// usually timeout - Valve throttles >30 UDP packets (>6930 servers) per min
			logger.WriteDebug("Master query error, likely due to Valve throttle/timeout :%s",
				err)
			break
		}
		// get hosts:ports beginning after header (0xFF, 0xFF, 0xFF, 0xFF, 0x66, 0x0A)
		ips, total, err := extractHosts(s[6:])
		if err != nil {
			return nil, logger.LogAppErrorf("Error when extracting addresses: %s",
				err)
		}
		retrieved = retrieved + total
		if retrieved >= maxHosts {
			logger.LogSteamInfo("Max host limit of %d reached!", maxHosts)
			logger.WriteDebug("Max host limit of %d reached!", maxHosts)
			break
		}
		logger.LogSteamInfo("%d hosts retrieved so far from master.", retrieved)
		logger.WriteDebug("%d hosts retrieved so far from master.", retrieved)
		for _, ip := range ips {
			serverlist = append(serverlist, ip)
		}

		if (serverlist[len(serverlist)-1]) != "0.0.0.0:0" {
			logger.LogSteamInfo("More hosts need to be retrieved. Last IP was: %s",
				serverlist[len(serverlist)-1])
			logger.WriteDebug("More hosts need to be retrieved. Last IP was: %s",
				serverlist[len(serverlist)-1])
			addr = serverlist[len(serverlist)-1]
		} else {
			logger.LogSteamInfo("IP retrieval complete!")
			logger.WriteDebug("IP retrieval complete!")
			break
		}
	}
	// remove 0.0.0.0:0
	if serverlist[len(serverlist)-1] == "0.0.0.0:0" {
		serverlist = serverlist[:len(serverlist)-1]
	}
	return serverlist, nil
}

func extractHosts(hbs []byte) ([]string, int, error) {
	var sl []string
	pos, total := 0, 0
	for i := 0; i < len(hbs); i++ {
		if len(sl) > 0 && sl[len(sl)-1] == "0.0.0.0:0" {
			logger.LogSteamInfo("0.0.0.0:0 detected. Got %d total hosts.", total-1)
			break
		}
		if pos+6 > len(hbs) {
			logger.LogSteamInfo("Got %d total hosts.", total)
			break
		}

		host, err := parseIP(hbs[pos : pos+6])
		if err != nil {
			logger.LogAppErrorf("Error parsing host: %s", err)
		} else {
			sl = append(sl, host)
			total++
		}
		// host:port = 6 bytes
		pos = pos + 6
	}
	return sl, total, nil
}

func parseIP(k []byte) (string, error) {
	if len(k) != 6 {
		return "", logger.LogSteamErrorf("Invalid IP byte size. Got: %d, expected 6",
			len(k))
	}
	port := int16(k[5]) | int16(k[4])<<8
	return fmt.Sprintf("%d.%d.%d.%d:%d", int(k[0]), int(k[1]), int(k[2]),
		int(k[3]), port), nil
}

func queryMasterServer(conn net.Conn, startaddress string,
	filter filters.Filter) ([]byte, error) {
	// Note: the connection is closed by the caller, do not close here, otherwise
	// Steam will continue to send the first batch of IPs and won't progress to the next batch
	startaddress = fmt.Sprintf("%s\x00", startaddress)
	addr := []byte(startaddress)
	request := []byte{0x31}
	request = append(request, filter.Region...)
	request = append(request, addr...)

	for i, f := range filter.Filters {
		for _, b := range f {
			request = append(request, b)
		}
		if i == len(filter.Filters)-1 {
			request = append(request, 0x00)
		}
	}

	_, err := conn.Write(request)
	if err != nil {
		logger.LogSteamError(ErrDataTransmit(err.Error()))
		return nil, ErrDataTransmit(err.Error())
	}

	var buf [maxPacketSize]byte
	numread, err := conn.Read(buf[:maxPacketSize])
	if err != nil {
		logger.LogSteamError(ErrDataTransmit(err.Error()))
		return nil, ErrDataTransmit(err.Error())
	}

	masterResponse := make([]byte, numread)
	copy(masterResponse, buf[:numread])

	if !bytes.HasPrefix(masterResponse, expectedMasterRespHeader) {
		logger.LogSteamError(ErrPacketHeader)
		return nil, ErrPacketHeader
	}

	return masterResponse, nil
}

// NewMasterQuery initiates a new Steam Master server query for a given filter,
// returning a pointer to a MasterQuery struct containing the hosts retrieved in
// the event of success or an error in the event of failure.
func NewMasterQuery(filter filters.Filter) (MasterQuery, error) {
	sl, err := getServers(filter)
	if err != nil {
		return MasterQuery{}, err
	}
	logger.LogSteamInfo("*** Retrieved %d %s servers.", len(sl), filter.Game.Name)

	return MasterQuery{Servers: sl}, nil
}
