package network

import (
	"encoding/json"
	"strings"
)

const (
	keyExMsgSep = ":"
)

var (
	EmptyKeyEx = keyExMessage{}
)

// keyExMessage represents an KeyEx kind message.
type keyExMessage struct {
	key  string
	name string
}

// NewKeyExMessageFromBytes returns a new Message of KeyEx kind.
func NewKeyExMessageFromBytes(b []byte) Message {
	s := string(b)
	if len(s) == 0 || !strings.Contains(s, keyExMsgSep) {
		return EmptyKeyEx
	}

	split := strings.SplitN(s, keyExMsgSep, 2)

	return keyExMessage{
		key:  split[0],
		name: split[1],
	}
}

// NewKeyExMessageWithName returns a new Message of KeyEx kind.
func NewKeyExMessageWithName(key, name string) Message {
	return keyExMessage{
		key:  key,
		name: name,
	}
}

// Type implements the Message interface.
func (m keyExMessage) Type() Kind {
	return KeyEx
}

// Payload implements the Message interface.
func (m keyExMessage) Payload() any {
	return []string{m.key, m.name}
}

// MarshalJSON implements the json.Marshaler interface.
func (m keyExMessage) MarshalJSON() ([]byte, error) {
	b := base{
		Kind: KeyEx,
		Data: []byte(m.key + keyExMsgSep + m.name),
	}
	return json.Marshal(b)
}
