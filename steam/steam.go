package steam

import "fmt"

const (
	headerStr       = "\xFF\xFF\xFF\xFF"
	maxHosts        = 2500 // max# hosts to retrieve; cannot be larger than 6930 as per steam
	maxPacketSize   = 1400 // specified by steam protocol
	QueryTimeout    = 3    // sec; connect, read & write timeout. Should be > 1
	QueryRetryCount = 3    // # of times to re-request rules, players, info on failure
)

type IgnoredRequest int

const (
	IgnoreRulesRequest IgnoredRequest = iota
	IgnorePlayerRequest
	IgnoreInfoRequest
)

var (
	// A2S_INFO: challenge request packet
	infoChallengeReq = []byte{
		0xFF, 0xFF, 0xFF, 0xFF,
		0x54, 0x53, 0x6F, 0x75, 0x72,
		0x63, 0x65, 0x20, 0x45, 0x6E,
		0x67, 0x69, 0x6E, 0x65, 0x20,
		0x51, 0x75, 0x65, 0x72, 0x79,
		0x00}
	// A2S_INFO: expected challenge response header
	expectedInfoRespHeader = []byte{0xFF, 0xFF, 0xFF, 0xFF, 0x49}

	// A2S_PLAYER: challenge request packet
	playerChallengeReq = []byte{0xFF, 0xFF, 0xFF, 0xFF, 0x55, 0xFF, 0xFF,
		0xFF, 0xFF}
	// A2S_PLAYER: expected challenge response header
	expectedPlayerRespHeader = []byte{0xFF, 0xFF, 0xFF, 0xFF, 0x41}
	// A2S_PLAYER: expected player chunk
	expectedPlayerChunkHeader = []byte{0xFF, 0xFF, 0xFF, 0xFF, 0x44}

	// A2S_RULES: challenge request packet
	rulesChallengeReq = []byte{0xFF, 0xFF, 0xFF, 0xFF, 0x56, 0xFF, 0xFF, 0xFF,
		0xFF}
	// A2S_RULES: expected challenge response header
	expectedRulesRespHeader = []byte{0xFF, 0xFF, 0xFF, 0xFF, 0x41}
	// A2S_RULES: expected rule chunk
	expectedRuleChunkHeader = []byte{0xFF, 0xFF, 0xFF, 0xFF, 0x45}

	// Steam master server: expected response header
	expectedMasterRespHeader = []byte{0xFF, 0xFF, 0xFF, 0xFF,
		0x66, 0x0A}
)

func removeFailedHost(failed []string, host string) []string {
	for i, v := range failed {
		if v == host {
			failed = append(failed[:i], failed[i+1:]...)
			fmt.Printf("removeFailedHost: removed: %s\n", host)
			fmt.Printf("removeFailedHost: new failed length: %d\n", len(failed))
			break
		}
	}
	return failed
}
