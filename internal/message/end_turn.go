package message

import (
	"encoding/json"
)

// endTurnMessage represents an EndTurn message.
type endTurnMessage struct{}

// NewEndTurn returns a new Message of EndTurn kind.
func NewEndTurn() Message {
	return endTurnMessage{}
}

// Type implements the Message interface.
func (m endTurnMessage) Type() Kind {
	return EndTurn
}

// Payload implements the Message interface.
func (m endTurnMessage) Payload() any {
	return nil
}

// MarshalJSON implements the json.Marshaler interface.
func (m endTurnMessage) MarshalJSON() ([]byte, error) {
	b := base{
		Kind: EndTurn,
		Data: nil,
	}
	return json.Marshal(b)
}
