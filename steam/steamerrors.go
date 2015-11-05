package steam

import (
	"errors"
	"fmt"
)

var (
	HostConnectionError = func(msg string) error {
		return fmt.Errorf("Steam: host connection error: %s\n", msg)
	}
	DataTransmitError = func(msg string) error {
		return fmt.Errorf("Steam: data transmission error: %s\n", msg)
	}
	ChallengeResponseError = errors.New("Steam: invalid challenge response header")
	PacketHeaderError      = errors.New("Steam: invalid packet header")
	NoPlayersError         = errors.New("Steam: server contains no players")
	NoRulesError           = errors.New("Steam: no A2S_RULES for server")
	NoInfoError            = errors.New("Steam: no A2S_INFO for server")
)
