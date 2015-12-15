package steam

// steamrules.go - steam server query for server information (A2S_RULES)

import (
	"bytes"
	"encoding/binary"
	"net"
	"sort"
	"steamtest/src/logger"
	"strings"
	"sync"
	"time"
)

func getRulesInfo(host string, timeout int) ([]byte, error) {
	conn, err := net.DialTimeout("udp", host, time.Duration(timeout)*time.Second)
	if err != nil {
		logger.LogSteamError(ErrHostConnection(err.Error()))
		return nil, ErrHostConnection(err.Error())
	}

	conn.SetDeadline(time.Now().Add(time.Duration(timeout-1) * time.Second))
	defer conn.Close()

	_, err = conn.Write(rulesChallengeReq)
	if err != nil {
		logger.LogSteamError(ErrDataTransmit(err.Error()))
		return nil, ErrDataTransmit(err.Error())
	}

	challengeNumResp := make([]byte, maxPacketSize)
	_, err = conn.Read(challengeNumResp)
	if err != nil {
		logger.LogSteamError(ErrDataTransmit(err.Error()))
		return nil, ErrDataTransmit(err.Error())
	}
	if !bytes.HasPrefix(challengeNumResp, expectedRulesRespHeader) {
		logger.LogSteamError(ErrChallengeResponse)
		return nil, ErrChallengeResponse
	}

	challengeNum := bytes.TrimLeft(challengeNumResp, headerStr)
	challengeNum = challengeNum[1:5]
	request := []byte{0xFF, 0xFF, 0xFF, 0xFF, 0x56}
	request = append(request, challengeNum...)

	_, err = conn.Write(request)
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
	var rulesInfo []byte
	if bytes.HasPrefix(buf[:maxPacketSize], multiPacketRespHeader) {
		// handle multi-packet response
		first := buf[:maxPacketSize]
		first = first[:numread]
		rulesInfo, err = handleMultiPacketResponse(conn, first)
		if err != nil {
			logger.LogSteamError(ErrDataTransmit(err.Error()))
			return nil, ErrDataTransmit(err.Error())
		}
	} else {
		rulesInfo = make([]byte, numread)
		copy(rulesInfo, buf[:numread])
	}
	return rulesInfo, nil
}

// Handle multi-packet responses for Source engine games.
func handleMultiPacketResponse(c net.Conn, firstReceived []byte) ([]byte,
	error) {
	// header: 4 bytes, 0xFFFFFFFE (already verified in caller)
	// ID: 4 bytes, signed
	// total # of packets: 1 byte, unsigned
	// current packet #, starts at zero: 1 byte, unsigned
	// size: 2 bytes, only for Orange Box Engine and Newer, signed
	// size & CRC32 sum for bzip2 compressed packets; but no longer used since late 2005

	// first 4 bytes [0:4] determine if split; we've already determined that it is
	id := int32(binary.LittleEndian.Uint32(firstReceived[4:8]))
	total := uint32(firstReceived[8])
	curNum := uint32(firstReceived[9])
	// note: size won't exist for 4 ancient appids (215,17550,17700,240 w/protocol 7)
	//size := int16(binary.LittleEndian.Uint16(firstReceived[10:12]))
	packets := make(map[uint32][]byte, total)
	packets[0] = firstReceived[12:]
	var buf [maxPacketSize]byte
	prevNum := curNum
	for {
		if curNum+1 == total {
			break
		}
		numread, err := c.Read(buf[:maxPacketSize])
		if err != nil {
			logger.LogSteamError(ErrMultiPacketTransmit(err.Error()))
			return nil, ErrMultiPacketTransmit(err.Error())
		}
		packet := buf[:maxPacketSize]
		packet = packet[:numread]
		curNum = uint32(packet[9])

		if prevNum == curNum {
			return nil, ErrMultiPacketDuplicate
		}
		prevNum = curNum

		if int32(binary.LittleEndian.Uint32(packet[4:8])) != id {
			logger.LogSteamError(ErrMultiPacketIDMismatch)
			return nil, ErrMultiPacketIDMismatch
		}
		if uint32(packet[9]) > total {
			logger.LogSteamError(ErrMultiPacketNumExceeded)
			return nil, ErrMultiPacketNumExceeded
		}
		// skip the header
		p := make([]byte, numread-12)
		copy(p, packet[12:])
		packets[curNum] = p

	}
	// sort packet keys
	pnums := make(u32slice, len(packets))
	for key := range packets {
		pnums = append(pnums, key)
	}
	pnums.Sort()

	var rules []byte
	for _, pn := range pnums {
		rules = append(rules, packets[pn]...)
	}
	return rules, nil
}

func parseRuleInfo(ruleinfo []byte) (map[string]string, error) {
	if !bytes.HasPrefix(ruleinfo, expectedRuleChunkHeader) {
		logger.LogSteamError(ErrPacketHeader)
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

// RetryFailedRulesReq retries a failed A2S_RULES request for a specified group of
// failed hosts for a total of retrycount times, returning a host to A2S_RULES
// mapping for any hosts that were successfully retried.
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

// GetRulesForServer requests A2S_RULES info for a given host within timeout seconds.
func GetRulesForServer(host string, timeout int) (map[string]string, error) {
	// Caller will log. Return err instead of wrapped logger.LogSteamError so as not
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

// u32slice attaches the methods of sort.Interface to []uint32, sorting in
// increasing order.
type u32slice []uint32

func (s u32slice) Len() int           { return len(s) }
func (s u32slice) Less(i, j int) bool { return s[i] < s[j] }
func (s u32slice) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

// Sort is a convenience method.
func (s u32slice) Sort() {
	sort.Sort(s)
}
