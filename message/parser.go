package message

import (
	"encoding/json"
	"fmt"
)

// base is used for two thing:
// - parsing message kinds;
// - and creating server type messages.
type base struct {
	Kind Kind
	Data json.RawMessage
}

// Type implements the Message interface.
func (mb base) Type() Kind {
	return mb.Kind
}

// Payload implements the Message interface.
func (mb base) Payload() any {
	return mb.Data
}

// Parse returns a Message from a raw []byte.
func Parse(raw []byte) (Message, error) {
	var b base
	if err := json.Unmarshal(raw, &b); err != nil {
		return nil, fmt.Errorf("parse raw message: %w", err)
	}

	// NOTE: GameState message cannot be sent by clients, so it is not checked here.
	switch b.Type() {
	case VoteToStart:
		return NewVoteToStart(), nil
	case PlayCard:
		return NewUnknown(), nil
	case Buy:
		return NewUnknown(), nil
	case Sell:
		return NewUnknown(), nil
	case EndTurn:
		return NewUnknown(), nil
	default:
		return NewUnknown(), nil
	}
}
