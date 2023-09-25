package message

import (
	"encoding/json"
)

type startTurnMessage struct{}

func NewStartTurn() Message {
	return startTurnMessage{}
}

func (m startTurnMessage) Type() Kind {
	return StartTurn
}

func (m startTurnMessage) Payload() any {
	return nil
}

func (m startTurnMessage) MarshalJSON() ([]byte, error) {
	b := base{
		Kind: StartTurn,
		Data: nil,
	}
	return json.Marshal(b)
}
