package message

import (
	"encoding/json"
	"strings"
)

const (
	keyExchangeSeparator = ":"
)

var (
	EmptyKeyExchange = keyExchangeMessage{}
)

// keyExchangeMessage represents an KeyExchange kind message.
type keyExchangeMessage struct {
	key  string
	name string
}

// NewKeyExchangeFromBytes returns a new Message of KeyExchange kind.
func NewKeyExchangeFromBytes(b []byte) Message {
	s := string(b)
	if len(s) == 0 || !strings.Contains(s, keyExchangeSeparator) {
		return EmptyKeyExchange
	}

	split := strings.SplitN(s, keyExchangeSeparator, 2)

	return keyExchangeMessage{
		key:  split[0],
		name: split[1],
	}
}

// NewKeyExchangeWithName returns a new Message of KeyExchange kind.
func NewKeyExchangeWithName(key, name string) Message {
	return keyExchangeMessage{
		key:  key,
		name: name,
	}
}

// Type implements the Message interface.
func (m keyExchangeMessage) Type() Kind {
	return KeyExchange
}

// Payload implements the Message interface.
func (m keyExchangeMessage) Payload() any {
	return []string{m.key, m.name}
}

// MarshalJSON implements the json.Marshaler interface.
func (m keyExchangeMessage) MarshalJSON() ([]byte, error) {
	b := base{
		Kind: KeyExchange,
		Data: []byte(m.key + keyExchangeSeparator + m.name),
	}
	return json.Marshal(b)
}
