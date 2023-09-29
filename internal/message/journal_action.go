package message

import (
	"encoding/json"
	"fmt"

	"github.com/milonoir/business-club-game/internal/game"
)

// ActorType represents the type of actor.
type ActorType uint8

const (
	ActorPlayer ActorType = iota
	ActorBank
)

// Action represents an action journal item.
type Action struct {
	ActorType ActorType
	Name      string
	Mod       *game.Modifier
	NewPrice  int
}

// journalActionMessage is a message that contains an action journal.
type journalActionMessage struct {
	action *Action
}

// NewJournalAction creates a new journal action message.
func NewJournalAction(action *Action) Message {
	return journalActionMessage{
		action: action,
	}
}

// NewJournalActionFromBytes creates a new journal action message from bytes.
func NewJournalActionFromBytes(b []byte) (Message, error) {
	var action Action
	if err := json.Unmarshal(b, &action); err != nil {
		return nil, fmt.Errorf("unmarshal action: %w", err)
	}
	return journalActionMessage{
		action: &action,
	}, nil
}

// Type implements the Message interface.
func (m journalActionMessage) Type() Kind {
	return JournalAction
}

// Payload implements the Message interface.
func (m journalActionMessage) Payload() any {
	return m.action
}

// MarshalJSON implements the json.Marshaler interface.
func (m journalActionMessage) MarshalJSON() ([]byte, error) {
	ab, err := json.Marshal(m.action)
	if err != nil {
		return nil, fmt.Errorf("marshal action: %w", err)
	}

	b := base{
		Kind: JournalAction,
		Data: ab,
	}
	return json.Marshal(b)
}
