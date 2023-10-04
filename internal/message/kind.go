package message

import (
	"errors"
	"fmt"
	"strconv"
)

// Kind is the type of the message.
type Kind uint8

// UnmarshalJSON implements the json.Unmarshaler interface.
func (k *Kind) UnmarshalJSON(b []byte) error {
	if b == nil {
		return errors.New("kind cannot be nil")
	}

	v, err := strconv.ParseUint(string(b), 10, 8)
	if err != nil {
		return fmt.Errorf("cannot parse kind: %w", err)
	}

	if kk := Kind(v); kk > JournalTrade {
		*k = Unknown
	} else {
		*k = kk
	}

	return nil
}

const (
	// Unknown is used for unidentified messages.
	Unknown Kind = iota

	// Error is a server type message that contains an error.
	Error

	// KeyExchange is used for sending/receiving reconnect keys.
	KeyExchange

	// StateUpdate is a server type message that contains the up-to-date game state sent to clients.
	StateUpdate

	// VoteToStart is a client type message that represents client readiness.
	VoteToStart

	// StartTurn is a server type message that signals a client that their turn has started.
	StartTurn

	// EndTurn is a client type message when a player wants to end their turn.
	EndTurn

	// PlayCard is a client type message when a player wants to play a card.
	PlayCard

	// TradeStock is a client type message when a player wants to trade stocks.
	TradeStock

	// JournalAction is a server type message that contains an action journal message.
	JournalAction

	// JournalTrade is a server type message that contains a trade journal message.
	JournalTrade
)
