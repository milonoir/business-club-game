package message

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewBuyStockMessageFromBytes(t *testing.T) {
	company := 1
	amount := 5
	b := []byte("1" + separator + "5")

	m, err := NewBuyStockMessageFromBytes(b)
	require.NoError(t, err)

	b, err = json.Marshal(m)
	require.NoError(t, err)

	pm, err := Parse(b)
	require.NoError(t, err)

	require.Equal(t, Buy, pm.Type())
	require.Equal(t, []int{company, amount}, pm.Payload())
}
