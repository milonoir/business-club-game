package message

import (
	"encoding/json"
)

// errorMessage represents a server Error message.
type errorMessage struct {
	err string
}

// NewError returns a new Message of Error kind.
func NewError(err string) Message {
	return errorMessage{
		err: err,
	}
}

// Type implements the Message interface.
func (m errorMessage) Type() Kind {
	return Error
}

// Payload implements the Message interface.
func (m errorMessage) Payload() any {
	return m.err
}

// MarshalJSON implements the json.Marshaler interface.
func (m errorMessage) MarshalJSON() ([]byte, error) {
	b := base{
		Kind: Error,
		Data: []byte(m.err),
	}
	return json.Marshal(b)
}
