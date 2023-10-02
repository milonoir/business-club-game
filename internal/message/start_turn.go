package message

import (
	"encoding/json"

	"github.com/milonoir/business-club-game/internal/game"
)

// startTurnMessage represents a StartTurn message.
type startTurnMessage struct {
	phase game.TurnPhase
}

// NewStartTurn returns a new Message of StartTurn kind.
func NewStartTurn(phase game.TurnPhase) Message {
	return startTurnMessage{phase: phase}
}

// NewStartTurnFromBytes returns a new Message of StartTurn kind.
func NewStartTurnFromBytes(b []byte) Message {
	return startTurnMessage{phase: game.TurnPhase(b[0])}
}

// Type implements the Message interface.
func (m startTurnMessage) Type() Kind {
	return StartTurn
}

// Payload implements the Message interface.
func (m startTurnMessage) Payload() any {
	return m.phase
}

// MarshalJSON implements the json.Marshaler interface.
func (m startTurnMessage) MarshalJSON() ([]byte, error) {
	b := base{
		Kind: StartTurn,
		Data: []byte{byte(m.phase)},
	}
	return json.Marshal(b)
}
