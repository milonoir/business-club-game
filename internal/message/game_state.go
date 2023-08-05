package message

import (
	"encoding/json"
	"fmt"

	"github.com/milonoir/business-club-game/internal/game"
)

type GameState struct {
	Started   bool
	Readiness []Readiness
	Companies []string
	Player    Player
	Hand      []game.Card
	Opponents []Player
}

type Readiness struct {
	Name  string
	Ready bool
}

type Player struct {
	Name string
	Portfolio
}

type Portfolio struct {
	Hidden bool
	Stocks [4]int
	Cash   int
}

type stateUpdateMessage struct {
	state *GameState
}

func NewStateUpdate(state *GameState) Message {
	return stateUpdateMessage{
		state: state,
	}
}

func NewStateUpdateFromBytes(b []byte) (Message, error) {
	var state GameState
	if err := json.Unmarshal(b, &state); err != nil {
		return nil, fmt.Errorf("unmarshal state: %w", err)
	}
	return stateUpdateMessage{
		state: &state,
	}, nil
}

func (m stateUpdateMessage) Type() Kind {
	return StateUpdate
}

func (m stateUpdateMessage) Payload() any {
	return m.state
}

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
