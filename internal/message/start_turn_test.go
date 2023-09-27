package message

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewStartTurn(t *testing.T) {
	m := NewStartTurn()

	require.Equal(t, StartTurn, m.Type())
	require.Nil(t, m.Payload())
}
