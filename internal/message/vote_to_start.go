package message

import (
	"encoding/json"
)

// voteToStartMessage represents a VoteToStart message.
type voteToStartMessage struct {
	ready bool
}

// NewVoteToStart returns a new Message of VoteToStart kind.
func NewVoteToStart(ready bool) Message {
	return voteToStartMessage{
		ready: ready,
	}
}

// NewVoteToStartFromBytes returns a new Message of VoteToStart kind.
func NewVoteToStartFromBytes(b []byte) Message {
	return voteToStartMessage{
		ready: b[0] == 1,
	}
}

// Type implements the Message interface.
func (m voteToStartMessage) Type() Kind {
	return VoteToStart
}

// Payload implements the Message interface.
func (m voteToStartMessage) Payload() any {
	return m.ready
}

// MarshalJSON implements the json.Marshaler interface.
func (m voteToStartMessage) MarshalJSON() ([]byte, error) {
	b := base{
		Kind: VoteToStart,
		Data: []byte{0},
	}
	if m.ready {
		b.Data[0] = 1
	}
	return json.Marshal(b)
}
