package network

import (
	"encoding/json"
	"strings"
)

const (
	authMsgSep = ":"
)

var (
	EmptyAuth = authMessage{}
)

// authMessage represents an Auth kind message.
type authMessage struct {
	key  string
	name string
}

// NewAuthMessageFromBytes returns a new Message of Auth kind.
func NewAuthMessageFromBytes(b []byte) Message {
	s := string(b)
	if len(s) == 0 || !strings.Contains(s, authMsgSep) {
		return EmptyAuth
	}

	split := strings.SplitN(s, authMsgSep, 2)

	return authMessage{
		key:  split[0],
		name: split[1],
	}
}

// NewAuthMessageWithName returns a new Message of Auth kind.
func NewAuthMessageWithName(key, name string) Message {
	return authMessage{
		key:  key,
		name: name,
	}
}

// Type implements the Message interface.
func (m authMessage) Type() Kind {
	return Auth
}

// Payload implements the Message interface.
func (m authMessage) Payload() any {
	return []string{m.key, m.name}
}

// MarshalJSON implements the json.Marshaler interface.
func (m authMessage) MarshalJSON() ([]byte, error) {
	b := base{
		Kind: Auth,
		Data: []byte(m.key + authMsgSep + m.name),
	}
	return json.Marshal(b)
}
