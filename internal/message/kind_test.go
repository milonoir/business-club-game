package message

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestKind_UnmarshalJSON(t *testing.T) {
	testcases := []struct {
		name    string
		raw     []byte
		expKind Kind
		expErr  error
	}{
		{
			name:    "error",
			raw:     []byte(`{"Kind":"Error"}`),
			expKind: Error,
		},
		{
			name:    "key exchange",
			raw:     []byte(`{"Kind":"KeyExchange"}`),
			expKind: KeyExchange,
		},
		{
			name:    "game state",
			raw:     []byte(`{"Kind":"StateUpdate"}`),
			expKind: StateUpdate,
		},
		{
			name:    "vote to start",
			raw:     []byte(`{"Kind":"VoteToStart"}`),
			expKind: VoteToStart,
		},
		{
			name:    "start turn",
			raw:     []byte(`{"Kind":"StartTurn"}`),
			expKind: StartTurn,
		},
		{
			name:    "end turn",
			raw:     []byte(`{"Kind":"EndTurn"}`),
			expKind: EndTurn,
		},
		{
			name:    "play a card",
			raw:     []byte(`{"Kind":"PlayCard"}`),
			expKind: PlayCard,
		},
		{
			name:    "trade stocks",
			raw:     []byte(`{"Kind":"TradeStock"}`),
			expKind: TradeStock,
		},
		{
			name:    "journal action",
			raw:     []byte(`{"Kind":"JournalAction"}`),
			expKind: JournalAction,
		},
		{
			name:    "journal trade",
			raw:     []byte(`{"Kind":"JournalTrade"}`),
			expKind: JournalTrade,
		},
		{
			name:    "unknown kind",
			raw:     []byte(`{"Kind":"Foobar"}`),
			expKind: Unknown,
		},
		{
			name:   "unmarshal error",
			raw:    []byte(`{kind:2}`),
			expErr: errors.New("invalid character 'k' looking for beginning of object key string"),
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			var v struct {
				Kind Kind
			}
			err := json.Unmarshal(tc.raw, &v)
			if tc.expErr != nil {
				require.EqualError(t, err, tc.expErr.Error())
				return
			}
			require.NoError(t, err)
			require.Equal(t, tc.expKind, v.Kind)
		})
	}
}
