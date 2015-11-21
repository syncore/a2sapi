package steam

// steamrules.go - testing steam server query for server information (A2S_RULES)

import (
	"bytes"
	"encoding/binary"
	"net"
	"steamtest/src/util"
	"strings"
	"sync"
	"time"
)

func getRulesInfo(host string, timeout int) ([]byte, error) {
	conn, err := net.DialTimeout("udp", host, time.Duration(timeout)*time.Second)
	if err != nil {
		// TODO: simplify this log + return into just a return
		util.LogAppError(ErrHostConnection(err.Error()))
		return nil, ErrHostConnection(err.Error())
	}

	conn.SetDeadline(time.Now().Add(time.Duration(timeout-1) * time.Second))
	defer conn.Close()

	_, err = conn.Write(rulesChallengeReq)
	if err != nil {
		// TODO: simplify this log + return into just a return
		util.LogAppError(ErrDataTransmit(err.Error()))
		return nil, ErrDataTransmit(err.Error())
	}

	challengeNumResp := make([]byte, maxPacketSize)
	_, err = conn.Read(challengeNumResp)
	if err != nil {
		// TODO: simplify this log + return into just a return
		util.LogAppError(ErrDataTransmit(err.Error()))
		return nil, ErrDataTransmit(err.Error())
	}
	if !bytes.HasPrefix(challengeNumResp, expectedRulesRespHeader) {
		// TODO: simplify this log + return into just a return
		util.LogAppError(ErrChallengeResponse)
		return nil, ErrChallengeResponse
	}

	challengeNum := bytes.TrimLeft(challengeNumResp, headerStr)
	challengeNum = challengeNum[1:5]
	request := []byte{0xFF, 0xFF, 0xFF, 0xFF, 0x56}

	for _, b := range challengeNum {
		request = append(request, b)
	}

	_, err = conn.Write(request)
	if err != nil {
		// TODO: simplify this log + return into just a return
		util.LogAppError(ErrDataTransmit(err.Error()))
		return nil, ErrDataTransmit(err.Error())
	}

	var buf [maxPacketSize]byte
	numread, err := conn.Read(buf[:maxPacketSize])
	if err != nil {
		// TODO: simplify this log + return into just a return
		util.LogAppError(ErrDataTransmit(err.Error()))
		return nil, ErrDataTransmit(err.Error())
	}
	rulesInfo := make([]byte, numread)
	copy(rulesInfo, buf[:numread])

	return rulesInfo, nil
}

func parseRuleInfo(ruleinfo []byte) (map[string]string, error) {
	// TODO: handle multi-packetted responses for games that use them
	if !bytes.HasPrefix(ruleinfo, expectedRuleChunkHeader) {
		// TODO: simplify this log + return into just a return
		util.LogAppError(ErrPacketHeader)
		return nil, ErrPacketHeader
	}

	ruleinfo = bytes.TrimLeft(ruleinfo, headerStr)
	numrules := int(binary.LittleEndian.Uint16(ruleinfo[1:3]))

	if numrules == 0 {
		return nil, ErrNoRules
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

func RetryFailedRulesReq(failed []string,
	retrycount int) map[string]map[string]string {

	m := make(map[string]map[string]string)
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
				r, err := GetRulesForServer(h, QueryTimeout)
				if err != nil {
					if err != ErrNoRules {
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

func GetRulesForServer(host string, timeout int) (map[string]string, error) {
	// Caller will log. Return err instead of wrapped util.LogAppError so as not
	// to interfere with custom error types that need to be analyzed when
	// determining if retry needs to be done.
	ri, err := getRulesInfo(host, timeout)
	if err != nil {
		return nil, err
	}

	rules, err := parseRuleInfo(ri)
	if err != nil {
		return nil, err
	}

	return rules, nil
}
