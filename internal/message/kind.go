package message

import (
	"encoding/json"
	"errors"
	"fmt"
)

// Kind is the type of the message.
type Kind string

// UnmarshalJSON implements the json.Unmarshaler interface.
func (k *Kind) UnmarshalJSON(b []byte) error {
	if b == nil {
		return errors.New("kind cannot be nil")
	}
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return fmt.Errorf("cannot parse kind: %w", err)
	}
	switch v := Kind(s); v {
	case Unknown, Ack, Error, KeyExchange, StateUpdate, VoteToStart, StartTurn, EndTurn, PlayCard, TradeStock, JournalAction, JournalTrade:
		*k = v
	default:
		*k = Unknown
	}
	return nil
}

const (
	// Unknown is used for unidentified messages.
	Unknown Kind = "Unknown"

	// Ack is used for acknowledging messages.
	Ack Kind = "Ack"

	// Error is a server type message that contains an error.
	Error Kind = "Error"

	// KeyExchange is used for sending/receiving reconnect keys.
	KeyExchange Kind = "KeyExchange"

	// StateUpdate is a server type message that contains the up-to-date game state sent to clients.
	StateUpdate Kind = "StateUpdate"

	// VoteToStart is a client type message that represents client readiness.
	VoteToStart Kind = "VoteToStart"

	// StartTurn is a server type message that signals a client that their turn has started.
	StartTurn Kind = "StartTurn"

	// EndTurn is a client type message when a player wants to end their turn.
	EndTurn Kind = "EndTurn"

	// PlayCard is a client type message when a player wants to play a card.
	PlayCard Kind = "PlayCard"

	// TradeStock is a client type message when a player wants to trade stocks.
	TradeStock Kind = "TradeStock"

	// JournalAction is a server type message that contains an action journal message.
	JournalAction Kind = "JournalAction"

	// JournalTrade is a server type message that contains a trade journal message.
	JournalTrade Kind = "JournalTrade"
)
