package message

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewEndTurn(t *testing.T) {
	m := NewEndTurn()

	require.Equal(t, EndTurn, m.Type())
	require.Nil(t, m.Payload())
}
