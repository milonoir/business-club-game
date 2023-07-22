package network

import (
	"encoding/json"
	"fmt"
)

// base is used for two things:
// - parsing message kinds;
// - and marshaling messages.
type base struct {
	Kind Kind
	Data []byte
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

	switch b.Type() {
	case Auth:
		return NewAuthMessage(b.Data), nil
	case GameState:
		return NewUnknown(), nil
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
