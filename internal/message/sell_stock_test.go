package message

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewSellStockMessageFromBytes(t *testing.T) {
	company := 1
	amount := 2
	b := []byte("1:2")

	m, err := NewSellStockMessageFromBytes(b)
	require.NoError(t, err)

	b, err = json.Marshal(m)
	require.NoError(t, err)

	pm, err := Parse(b)
	require.NoError(t, err)

	require.Equal(t, Sell, pm.Type())
	require.Equal(t, []int{company, amount}, pm.Payload())
}
