// steamrules.go - testing steam server query for server information (A2S_RULES)

package steam

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"strings"
	"time"
)

type SteamConn struct {
	buf     [maxPacketSize]byte
	c       net.Conn
	timeout time.Duration
}

func NewSteamConn(host string, timeout int) (*SteamConn,
	error) {
	conn, err := net.DialTimeout("udp", host, time.Duration(timeout)*time.Second)
	if err != nil {
		return nil, HostConnectionError(err.Error())
	}

	return &SteamConn{
		c:       conn,
		timeout: time.Duration(timeout) * time.Second,
	}, nil
}

func (sc *SteamConn) Send(b []byte) error {
	if sc.timeout > 0 {
		sc.c.SetWriteDeadline(time.Now().Add(sc.timeout))
	}
	_, err := sc.c.Write(b)
	if err != nil {
		return DataTransmitError(err.Error())
	}
	return nil
}

func (sc *SteamConn) Recv() ([]byte, error) {
	if sc.timeout > 0 {
		sc.c.SetReadDeadline(time.Now().Add(sc.timeout))
	}

	numread, err := sc.c.Read(sc.buf[:maxPacketSize])
	if err != nil {
		return nil, DataTransmitError(err.Error())
	}
	b := make([]byte, numread)
	copy(b, sc.buf[:numread])
	return b, nil
}

func (sc *SteamConn) Close() {
	sc.c.Close()
}

// func GetRulesForServerList(servers []string) map[string]map[string]string {
// 	m := make(map[string]map[string]string)
// 	challengenumreq := []byte{0xFF, 0xFF, 0xFF, 0xFF, 0x56, 0xFF, 0xFF, 0xFF,
// 		0xFF}
// 	for _, h := range servers {
// 		sc, err := NewSteamConn(h, 3)
// 		if err != nil {
// 			continue
// 		}
// 		defer sc.Close()
// 		err = sc.Send(challengenumreq)
// 		if err != nil {
// 			continue
// 		}
// 		challengenumresp, err := sc.Recv()
// 		if err != nil {
// 			continue
// 		}
// 		if !bytes.HasPrefix(challengenumresp, []byte{0xFF, 0xFF, 0xFF, 0xFF, 0x41}) {
// 			continue
// 		}
// 		challengenum := bytes.TrimLeft(challengenumresp, headerStr)
// 		challengenum = challengenum[1:5]
// 		request := []byte{0xFF, 0xFF, 0xFF, 0xFF, 0x56}
// 		for _, b := range challengenum {
// 			request = append(request, b)
// 		}
// 		err = sc.Send(request)
// 		if err != nil {
// 			continue
// 		}
// 		unparsed, err := sc.Recv()
// 		if err != nil {
// 			continue
// 		}
// 		parsed, err := parseRuleInfo(unparsed)
// 		if err != nil {
// 			continue
// 		}
// 		m[h] = parsed
// 	}
// 	return m
// }

func getRulesInfo(host string, timeout int) ([]byte, error) {
	challengeNumReq := []byte{0xFF, 0xFF, 0xFF, 0xFF, 0x56, 0xFF, 0xFF, 0xFF,
		0xFF}
	conn, err := net.DialTimeout("udp", host, time.Duration(timeout)*time.Second)
	if err != nil {
		return nil, HostConnectionError(err.Error())
	}

	conn.SetDeadline(time.Now().Add(time.Duration(timeout) * time.Second))
	defer conn.Close()

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
	request := []byte{0xFF, 0xFF, 0xFF, 0xFF, 0x56}

	for _, b := range challengeNum {
		request = append(request, b)
	}

	//+fmt.Printf("will send: %x\n", request)

	_, err = conn.Write(request)
	if err != nil {
		return nil, DataTransmitError(err.Error())
	}

	var buf [maxPacketSize]byte

	numread, err := conn.Read(buf[:maxPacketSize])
	if err != nil {
		return nil, DataTransmitError(err.Error())
	}
	rulesInfo := make([]byte, numread)
	copy(rulesInfo, buf[:numread])

	//fmt.Printf("Server info is: %x", rulesInfo)
	return rulesInfo, nil
}

func parseRuleInfo(ruleinfo []byte) (map[string]string, error) {
	// A2S_RULES response appears to not be multi-packetted for QL

	//Data			Type			Comment
	//--------------------------------------------------------------------
	//Header  	byte  		Always equal to 'E' (0x45)
	//Rules  		short  		Number of rules in the response.

	// For every rule in "Rules" there is this chunk in the response:
	//--------------------------------------------------------------------

	// 	Data 			Type 			Comment
	// 	Name  		string  	Name of the rule.
	//	Value  		string  	Value of the rule.
	if !bytes.HasPrefix(ruleinfo, []byte{0xFF, 0xFF, 0xFF, 0xFF, 0x45}) {
		return nil, PacketHeaderError
	}

	ruleinfo = bytes.TrimLeft(ruleinfo, headerStr)

	numrules := int(binary.LittleEndian.Uint16(ruleinfo[1:3]))
	//fmt.Printf("Number of rules: %d\n", numrules)
	if numrules == 0 {
		return nil, NoRulesError
	}

	b := bytes.Split(ruleinfo[3:], []byte{0x00})
	m := make(map[string]string)

	var key string
	for i, y := range b {
		if i%2 != 1 {
			key = strings.TrimRight(string(y), "\x00")
		} else {
			m[key] = strings.TrimRight(string(b[i]), "\x00")
		}
	}

	return m, nil
}

func GetRulesWithTestData() error {
	ri := testRuleInfoData
	fmt.Printf("Rule info is: %x\n", ri)
	rules, err := parseRuleInfo(ri)
	if err != nil {
		return fmt.Errorf("Error parsing rules info with test data: %s\n", err)
	}
	for k, v := range rules {
		fmt.Printf("key: %s, value: %s\n", k, v)
	}

	return nil
}

func GetRulesWithLiveData(host string, timeout int) (map[string]string, error) {
	ri, err := getRulesInfo(host, timeout)
	if err != nil {
		return nil, err
	}
	//fmt.Printf("Rule info is: %x\n", ri)
	rules, err := parseRuleInfo(ri)
	if err != nil {
		return nil, err
	}
	// for k, v := range rules {
	// 	fmt.Printf("key: %s, value: %s\n", k, v)
	// }

	return rules, nil
}
