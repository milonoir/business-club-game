package message

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewPlayCard(t *testing.T) {
	id := 1
	company := 2

	m := NewPlayCard(id, company)

	b, err := json.Marshal(m)
	require.NoError(t, err)

	pm, err := Parse(b)
	require.NoError(t, err)

	require.Equal(t, PlayCard, pm.Type())
	require.Equal(t, []int{id, company}, pm.Payload())
}

func TestNewPlayCardFromBytes(t *testing.T) {
	id := 1
	company := 2
	b := []byte("1" + separator + "2")

	m, err := NewPlayCardFromBytes(b)
	require.NoError(t, err)

	b, err = json.Marshal(m)
	require.NoError(t, err)

	pm, err := Parse(b)
	require.NoError(t, err)

	require.Equal(t, PlayCard, pm.Type())
	require.Equal(t, []int{id, company}, pm.Payload())
}
