package message

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewTradeStockFromBytes(t *testing.T) {
	trade := TradeBuy
	company := 1
	amount := 5
	b := []byte("0" + separator + "1" + separator + "5")

	m, err := NewTradeStockFromBytes(b)
	require.NoError(t, err)

	b, err = json.Marshal(m)
	require.NoError(t, err)

	pm, err := Parse(b)
	require.NoError(t, err)

	require.Equal(t, TradeStock, pm.Type())
	require.Equal(t, []any{trade, company, amount}, pm.Payload())
}
