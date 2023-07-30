package network

import (
	"encoding/json"
)

type errorMessage struct {
	err string
}

func NewErrorMessage(err string) Message {
	return errorMessage{
		err: err,
	}
}

func (m errorMessage) Type() Kind {
	return Error
}

func (m errorMessage) Payload() any {
	return m.err
}

func (m errorMessage) MarshalJSON() ([]byte, error) {
	b := base{
		Kind: Error,
		Data: []byte(m.err),
	}
	return json.Marshal(b)
}
