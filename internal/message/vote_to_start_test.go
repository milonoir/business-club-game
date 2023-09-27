package message

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewVoteToStart(t *testing.T) {
	for _, tc := range []bool{true, false} {
		m := NewVoteToStart(tc)

		b, err := json.Marshal(m)
		require.NoError(t, err)

		pm, err := Parse(b)
		require.NoError(t, err)
		require.Equal(t, VoteToStart, pm.Type())
		require.Equal(t, tc, pm.Payload())
	}
}

func TestNewVoteToStartFromBytes(t *testing.T) {
	testcases := []struct {
		input bool
		data  []byte
	}{
		{
			input: true,
			data:  []byte{1},
		},
		{
			input: false,
			data:  []byte{0},
		},
	}
	for _, tc := range testcases {
		m := NewVoteToStartFromBytes(tc.data)

		b, err := json.Marshal(m)
		require.NoError(t, err)

		pm, err := Parse(b)
		require.NoError(t, err)
		require.Equal(t, VoteToStart, pm.Type())
		require.Equal(t, tc.input, pm.Payload())
	}
}
