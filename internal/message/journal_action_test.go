package message

import (
	"encoding/json"
	"testing"

	"github.com/milonoir/business-club-game/internal/game"
	"github.com/stretchr/testify/require"
)

func TestNewJournalAction(t *testing.T) {
	action := &Action{
		ActorType: ActorPlayer,
		Name:      "player 1",
		Mod: &game.Modifier{
			Company: 2,
			Mod:     game.NewMod("+", 100),
		},
		NewPrice: 320,
	}

	m := NewJournalAction(action)

	b, err := json.Marshal(m)
	require.NoError(t, err)

	pm, err := Parse(b)
	require.NoError(t, err)

	require.Equal(t, JournalAction, pm.Type())
	require.Equal(t, action, pm.Payload())
}

func TestNewJournalActionFromBytes(t *testing.T) {
	b := []byte(`{
	"ActorType": 0,
	"Name": "player 1",
	"Mod": {
		"Company": 2,
		"Mod": "+ 100"
	},
	"NewPrice":320
}`)

	pm, err := NewJournalActionFromBytes(b)
	require.NoError(t, err)
	require.Equal(t, JournalAction, pm.Type())

	action := &Action{
		ActorType: ActorPlayer,
		Name:      "player 1",
		Mod: &game.Modifier{
			Company: 2,
			Mod:     game.NewMod("+", 100),
		},
		NewPrice: 320,
	}
	require.Equal(t, action, pm.Payload())
}
