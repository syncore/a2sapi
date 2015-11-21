package steam

import (
	"errors"
	"fmt"
)

// Errors
var (
	ErrHostConnection = func(msg string) error {
		return fmt.Errorf("Steam: host connection error: %s", msg)
	}
	ErrDataTransmit = func(msg string) error {
		return fmt.Errorf("Steam: data transmission error: %s", msg)
	}
	ErrChallengeResponse = errors.New("Steam: invalid challenge response header")
	ErrPacketHeader      = errors.New("Steam: invalid packet header")
	ErrNoPlayers         = errors.New("Steam: server contains no players")
	ErrNoRules           = errors.New("Steam: no A2S_RULES for server")
	ErrNoInfo            = errors.New("Steam: no A2S_INFO for server")
)
