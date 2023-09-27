package message

import (
	"encoding/json"
)

// startTurnMessage represents a StartTurn message.
type startTurnMessage struct{}

// NewStartTurn returns a new Message of StartTurn kind.
func NewStartTurn() Message {
	return startTurnMessage{}
}

// Type implements the Message interface.
func (m startTurnMessage) Type() Kind {
	return StartTurn
}

// Payload implements the Message interface.
func (m startTurnMessage) Payload() any {
	return nil
}

// MarshalJSON implements the json.Marshaler interface.
func (m startTurnMessage) MarshalJSON() ([]byte, error) {
	b := base{
		Kind: StartTurn,
		Data: nil,
	}
	return json.Marshal(b)
}
