package message

import (
	"encoding/json"
	"fmt"

	"github.com/milonoir/business-club-game/internal/game"
)

// GameState represents the state of the game.
type GameState struct {
	Started       bool
	Readiness     []Readiness
	Turn          int
	PlayerOrder   []string
	CurrentPlayer int
	Companies     []string
	StockPrices   [4]int
	Player        game.Player
	Opponents     []game.Player
}

// Readiness represents a player's readiness.
type Readiness struct {
	Name  string
	Ready bool
}

// stateUpdateMessage is a message that contains the game state.
type stateUpdateMessage struct {
	state *GameState
}

// NewStateUpdate creates a new state update message.
func NewStateUpdate(state *GameState) Message {
	return stateUpdateMessage{
		state: state,
	}
}

// NewStateUpdateFromBytes creates a new state update message from bytes.
func NewStateUpdateFromBytes(b []byte) (Message, error) {
	var state GameState
	if err := json.Unmarshal(b, &state); err != nil {
		return nil, fmt.Errorf("unmarshal state: %w", err)
	}
	return stateUpdateMessage{
		state: &state,
	}, nil
}

// Type implements the Message interface.
func (m stateUpdateMessage) Type() Kind {
	return StateUpdate
}

// Payload implements the Message interface.
func (m stateUpdateMessage) Payload() any {
	return m.state
}

// MarshalJSON implements the json.Marshaler interface.
func (m stateUpdateMessage) MarshalJSON() ([]byte, error) {
	sb, err := json.Marshal(m.state)
	if err != nil {
		return nil, fmt.Errorf("marshal state: %w", err)
	}

	b := base{
		Kind: StateUpdate,
		Data: sb,
	}
	return json.Marshal(b)
}
