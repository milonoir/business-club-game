package network

import (
	"encoding/json"
)

type voteToStartMessage struct {
	ready bool
}

func NewVoteToStartMessage(ready bool) Message {
	return voteToStartMessage{
		ready: ready,
	}
}

func NewVoteToStartMessageFromBytes(b []byte) Message {
	return voteToStartMessage{
		ready: b[0] == 1,
	}
}

func (m voteToStartMessage) Type() Kind {
	return VoteToStart
}

func (m voteToStartMessage) Payload() any {
	return m.ready
}

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
