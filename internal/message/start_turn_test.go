package message

import (
	"testing"

	"github.com/milonoir/business-club-game/internal/game"
	"github.com/stretchr/testify/require"
)

func TestNewStartTurn(t *testing.T) {
	for _, tc := range []game.TurnPhase{game.ActionPhase, game.TradePhase} {
		m := NewStartTurn(tc)

		require.Equal(t, StartTurn, m.Type())
		require.Equal(t, tc, m.Payload())
	}
}

func TestNewStartTurnFromBytes(t *testing.T) {
	testcases := []struct {
		input game.TurnPhase
		data  []byte
	}{
		{
			input: game.ActionPhase,
			data:  []byte{0},
		},
		{
			input: game.TradePhase,
			data:  []byte{1},
		},
	}
	for _, tc := range testcases {
		m := NewStartTurnFromBytes(tc.data)

		require.Equal(t, StartTurn, m.Type())
		require.Equal(t, tc.input, m.Payload())
	}
}
