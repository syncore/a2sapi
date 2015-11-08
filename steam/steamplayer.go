// steamplayer.go - testing steam server query for players (A2S_PLAYER)
package steam

import (
	"bytes"
	"encoding/binary"
	"math"
	"net"
	"time"
)

type PlayerInfo struct {
	Name              string  `json:"name"`
	Score             int32   `json:"score"`
	TimeConnectedSecs float32 `json:"secsConnected"`
	TimeConnectedTot  string  `json:"totalConnected"`
}

func getPlayerInfo(host string, timeout int) ([]byte, error) {
	challengeNumReq := []byte{0xFF, 0xFF, 0xFF, 0xFF, 0x55, 0xFF, 0xFF, 0xFF,
		0xFF}

	conn, err := net.DialTimeout("udp", host, time.Duration(timeout)*time.Second)
	if err != nil {
		return nil, HostConnectionError(err.Error())
	}

	defer conn.Close()
	conn.SetDeadline(time.Now().Add(time.Duration(timeout-1) * time.Second))

	_, err = conn.Write(challengeNumReq)
	if err != nil {
		return nil, DataTransmitError(err.Error())
	}

	challengeNumResp := make([]byte, maxPacketSize)

	_, err = conn.Read(challengeNumResp)
	if err != nil {
		return nil, DataTransmitError(err.Error())
	}
	if !bytes.HasPrefix(challengeNumResp, []byte{0xFF, 0xFF, 0xFF, 0xFF, 0x41}) {
		return nil, ChallengeResponseError
	}
	challengeNum := bytes.TrimLeft(challengeNumResp, headerStr)
	challengeNum = challengeNum[1:5]

	//fmt.Printf("Reply from server: %x\n", challengeNumResp)
	//fmt.Printf("Challenge number is: %x\n", challengeNum)
	request := []byte{0xFF, 0xFF, 0xFF, 0xFF, 0x55}

	for _, b := range challengeNum {
		request = append(request, b)
	}

	//fmt.Printf("will send: %x\n", request)

	_, err = conn.Write(request)
	if err != nil {
		return nil, DataTransmitError(err.Error())
	}
	var buf [maxPacketSize]byte
	numread, err := conn.Read(buf[:maxPacketSize])
	if err != nil {
		return nil, DataTransmitError(err.Error())
	}
	pi := make([]byte, numread)
	copy(pi, buf[:numread])

	//fmt.Printf("Player info is: %x", playerInfo)
	//return playerInfo, nil

	return pi, nil
}

func parsePlayerInfo(unparsed []byte) ([]*PlayerInfo, error) {
	//Data			Type			Comment
	//--------------------------------------------------------------------
	//Header  	byte  		Always equal to 'D' (0x44)
	//Players  	byte  		Number of players whose information was gathered.

	//For every player in "Players" there is this chunk in the response:
	//--------------------------------------------------------------------

	//	Data			Type		Comment
	//	Index 		byte  	Index of player chunk starting from 0.
	//	Name  		string 	Name of the player.
	//	Score  		long  	Player's score (usually "frags" or "kills".)
	//	Duration  float  	Time (in seconds) player has been connected to the server

	if !bytes.HasPrefix(unparsed, []byte{0xFF, 0xFF, 0xFF, 0xFF, 0x44}) {
		return nil, PacketHeaderError
	}
	unparsed = bytes.TrimLeft(unparsed, headerStr)

	numplayers := int(unparsed[1])
	//fmt.Printf("Number of players: %d\n", numplayers)
	if numplayers == 0 {
		return nil, NoPlayersError
	}

	//fmt.Printf("New trimmed player info slice is: %x\n", unparsed)
	players := []*PlayerInfo{}

	// index 0 = '44' | 1 = 'numplayers' byte | 2 = player 1 separator byte '00'
	// | 3 = start of player 1 name; additional player start indexes are player separator + 1
	startidx := 3
	var b []byte
	for i := 0; i < numplayers; i++ {
		if i == 0 {
			b = unparsed[startidx:]
		} else {
			b = b[startidx+1:]
		}
		nul := bytes.IndexByte(b, 0x00)
		name := b[:nul]              // string (variable length)
		score := b[nul+1 : nul+5]    // long (4 bytes)
		duration := b[nul+5 : nul+9] // float (4 bytes)
		startidx = nul + 9

		seconds, timeformatted := getDuration(duration)
		players = append(players, &PlayerInfo{
			Name:              string(name),
			Score:             int32(binary.LittleEndian.Uint32(score)),
			TimeConnectedSecs: seconds,
			TimeConnectedTot:  timeformatted,
		})
	}

	return players, nil
}

func getDuration(bytes []byte) (float32, string) {
	bits := binary.LittleEndian.Uint32(bytes)
	f := math.Float32frombits(bits)
	s := time.Duration(int64(f)) * time.Second
	return f, s.String()
}

// func GetPlayersWithTestData() error {
// 	pi := testPlayerInfoData
// 	players, err := parsePlayerInfo(pi)
// 	if err != nil {
// 		return fmt.Errorf("Error parsing player info: %s\n", err)
// 	}
// 	for _, p := range players {
// 		fmt.Printf("Name: %s, Score: %d, Connected for: %s\n",
// 			p.Name, p.Score, p.TimeConnectedTot)
// 	}
// 	return nil
// }

func GetPlayersWithLiveData(host string, timeout int) ([]*PlayerInfo, error) {
	//var open bool
	pi, err := getPlayerInfo(host, timeout)
	if err != nil {
		return nil, err
	}
	//fmt.Printf("Player info is: %x\n", pi)
	players, err := parsePlayerInfo(pi)
	if err != nil {
		return nil, err
	}
	// }
	// for _, p := range players {
	// 	fmt.Printf("Name: %s, Score: %d, Connected for: %s\n",
	// 		p.Name, p.Score, p.TimeConnectedTot)
	// }
	return players, nil
}
