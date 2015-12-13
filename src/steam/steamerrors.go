package steam

import (
	"errors"
	"fmt"
)

// Errors
var (
	// ErrHostConnection is an error related to the establishment of a connection.
	ErrHostConnection = func(msg string) error {
		return fmt.Errorf("Steam: host connection error: %s", msg)
	}
	// ErrDataTransmit is an error related to sending data to a connection.
	ErrDataTransmit = func(msg string) error {
		return fmt.Errorf("Steam: data transmission error: %s", msg)
	}
	// ErrMultiPacketTransmit is an error related to sending data to a connection
	//  in the multi-packet context of A2S_RULES.
	ErrMultiPacketTransmit = func(msg string) error {
		return fmt.Errorf("Steam: multi-packet data transmission error: %s", msg)
	}
	// ErrChallengeResponse is an error thrown for an invalid challense response
	// header.
	ErrChallengeResponse = errors.New("Steam: invalid challenge response header")

	// ErrPacketHeader is an error thrown upon detection of an invalid packet header.
	ErrPacketHeader = errors.New("Steam: invalid packet header")

	// ErrMultiPacketDuplicate is an error thrown when a duplicate packet is
	// detected int he multi-packet context of A2S_RULES.
	ErrMultiPacketDuplicate = errors.New(
		"Steam: multi-packet: duplicate packet detected")

	// ErrMultiPacketIDMismatch is an error thrown in the context of multi-packet
	// A2S_RULES when the current packet ID does match the packet ID for the batch
	// of multiple packets currently being processed.
	ErrMultiPacketIDMismatch = errors.New(
		"Steam: multi-packet error: packet ID mismatch")

	// ErrMultiPacketNumExceeded is an error thrown in the A2S_RULES multi-packet
	// context when the current packet's number is greater than the total number of
	// packets to be parsed within the current batch.
	ErrMultiPacketNumExceeded = errors.New(
		"Steam: multi-packet error: packet number greater than total")

	// ErrNoPlayers is a generic error thrown when a server is empty.
	ErrNoPlayers = errors.New("Steam: server contains no players")

	// ErrNoRules is a generic error thrown when no A2S_RULES data could be parsed
	// for the given server.
	ErrNoRules = errors.New("Steam: no A2S_RULES for server")

	// ErrNoInfo is a generic error thrown when no A2S_INFO could be parsed for the
	// given server.
	ErrNoInfo = errors.New("Steam: no A2S_INFO for server")
)
