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
			raw:     []byte(`{"Kind":1}`),
			expKind: Error,
		},
		{
			name:    "key exchange",
			raw:     []byte(`{"Kind":2}`),
			expKind: KeyExchange,
		},
		{
			name:    "game state",
			raw:     []byte(`{"Kind":3}`),
			expKind: StateUpdate,
		},
		{
			name:    "vote to start",
			raw:     []byte(`{"Kind":4}`),
			expKind: VoteToStart,
		},
		{
			name:    "play a card",
			raw:     []byte(`{"Kind":5}`),
			expKind: PlayCard,
		},
		{
			name:    "buy stocks",
			raw:     []byte(`{"Kind":6}`),
			expKind: Buy,
		},
		{
			name:    "sell stocks",
			raw:     []byte(`{"Kind":7}`),
			expKind: Sell,
		},
		{
			name:    "end turn",
			raw:     []byte(`{"Kind":8}`),
			expKind: EndTurn,
		},
		{
			name:    "no kind",
			raw:     []byte(`{"Type":6}`),
			expKind: Unknown,
		},
		{
			name:    "unknown kind",
			raw:     []byte(`{"Kind":71}`),
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
