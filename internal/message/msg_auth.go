package message

import (
	"encoding/json"
)

var (
	EmptyAuth = NewAuthMessage(nil)
)

// authMessage represents an Auth kind message.
type authMessage struct {
	key string
}

// NewAuthMessage returns new Message of Auth kind.
func NewAuthMessage(b []byte) Message {
	return authMessage{
		key: string(b),
	}
}

// Type implements the Message interface.
func (m authMessage) Type() Kind {
	return Auth
}

// Payload implements the Message interface.
func (m authMessage) Payload() any {
	return m.key
}

// MarshalJSON implements the json.Marshaler interface.
func (m authMessage) MarshalJSON() ([]byte, error) {
	b := base{
		Kind: Auth,
		Data: []byte(m.key),
	}
	return json.Marshal(b)
}
