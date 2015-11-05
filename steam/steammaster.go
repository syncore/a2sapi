// steammaster.go - testing steam master server query
package steam

import (
	"bytes"
	"fmt"
	"net"
	"steamtest/steam/filters"
	"time"
)

const masterServerHost = "hl2master.steampowered.com:27011"

func getServers(region filters.ServerRegion,
	filters ...filters.ServerFilter) ([]string, error) {

	var serverlist []string
	var c net.Conn
	var err error
	retrieved := 0
	count := 0
	complete := false
	addr := "0.0.0.0:0"

	c, err = net.DialTimeout("udp", masterServerHost,
		time.Duration(QueryTimeout)*time.Second)
	if err != nil {
		return nil, HostConnectionError(err.Error())
	}

	defer c.Close()
	c.SetDeadline(time.Now().Add(time.Duration(QueryTimeout) * time.Second))

	for !complete {
		s, err := queryMasterServer(c, addr, region, filters...)
		if err != nil {
			return nil, err
		}
		// get hosts:ports beginning after header (0xFF, 0xFF, 0xFF, 0xFF, 0x66, 0x0A)
		ips, total, err := getHosts(s[6:])
		if err != nil {
			return nil, fmt.Errorf("Error when extracting addresses: %s", err)
		}
		retrieved = retrieved + total
		fmt.Printf("%d hosts retrieved so far from master.\n", retrieved)
		for _, ip := range ips {
			if count >= maxHosts {
				fmt.Printf("Max host limit of %d reached!\n", maxHosts)
				complete = true
				break
			}
			serverlist = append(serverlist, ip)
			count++
		}

		if (serverlist[len(serverlist)-1]) != "0.0.0.0:0" {
			fmt.Printf("More ips need to be retrieved. Last ip was: %s\n",
				serverlist[len(serverlist)-1])
			addr = serverlist[len(serverlist)-1]
			fmt.Printf("Seeding next scan with host: %s\n", addr)
		} else {
			fmt.Println("Ip retrieval complete!")
			complete = true
			break
		}
	}
	// remove 0.0.0.0:0
	if serverlist[len(serverlist)-1] == "0.0.0.0:0" {
		serverlist = serverlist[:len(serverlist)-1]
	}
	return serverlist, nil
}

func getHosts(hbs []byte) ([]string, int, error) {
	var sl []string
	pos := 0
	total := 0
	//fmt.Printf("length of host byte slice: %d\n", len(hbs))
	for i := 0; i < len(hbs); i++ {
		if len(sl) > 0 && sl[len(sl)-1] == "0.0.0.0:0" {
			fmt.Printf("0.0.0.0:0 detected. Got %d total hosts.\n", total-1)
			break
		}
		if pos+6 > len(hbs) {
			fmt.Printf("Got %d total hosts.\n", total)
			break
		}

		host, err := parseIP(hbs[pos : pos+6])
		if err != nil {
			fmt.Printf("Error parsing host: %s\n", err)
		} else {
			//fmt.Printf("%s\n", host)
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
		return "", fmt.Errorf("Invalid ip byte size. Got: %d, expected 6\n", len(k))
	}
	port := int16(k[5]) | int16(k[4])<<8
	return fmt.Sprintf("%d.%d.%d.%d:%d", int(k[0]), int(k[1]), int(k[2]),
		int(k[3]), port), nil
}

func queryMasterServer(conn net.Conn, startaddress string,
	region filters.ServerRegion, filters ...filters.ServerFilter) ([]byte, error) {
	// Note: the connection is closed by the caller, do not close here, otherwise
	// Steam will continue to send the first batch of IPs and won't progress to the next batch

	startaddress = fmt.Sprintf("%s\x00", startaddress)
	addr := []byte(startaddress)
	request := []byte{0x31}
	for _, b := range region {
		request = append(request, b)
	}
	for _, b := range addr {
		request = append(request, b)
	}
	for i, f := range filters {
		for _, b := range f {
			request = append(request, b)
		}
		if i == len(filters)-1 {
			request = append(request, 0x00)
		}
	}

	_, err := conn.Write(request)
	if err != nil {
		return nil, DataTransmitError(err.Error())
	}

	var buf [maxPacketSize]byte
	numread, err := conn.Read(buf[:maxPacketSize])
	if err != nil {
		return nil, DataTransmitError(err.Error())
	}

	masterResponse := make([]byte, numread)
	copy(masterResponse, buf[:numread])

	if !bytes.HasPrefix(masterResponse, []byte{0xFF, 0xFF, 0xFF, 0xFF,
		0x66, 0x0A}) {
		return nil, PacketHeaderError
	}

	//fmt.Printf("Master server response is: %x", masterResponse)
	return masterResponse, nil
}

func GetServerListWithLiveData(region filters.ServerRegion,
	filters ...filters.ServerFilter) ([]string, error) {

	sl, err := getServers(region, filters...)
	if err != nil {
		return nil, err
	}
	// for _, v := range sl {
	// 	fmt.Printf("%s\n", v)
	// }
	fmt.Printf("*** Retrieved %d servers.\n", len(sl))

	return sl, nil
}
