package message

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewJournalTrade(t *testing.T) {
	trade := &Trade{
		Name:    "player 3",
		Type:    TradeBuy,
		Company: 3,
		Amount:  15,
		Price:   3250,
	}

	m := NewJournalTrade(trade)

	b, err := json.Marshal(m)
	require.NoError(t, err)

	pm, err := Parse(b)
	require.NoError(t, err)

	require.Equal(t, JournalTrade, pm.Type())
	require.Equal(t, trade, pm.Payload())
}

func TestNewJournalTradeFromBytes(t *testing.T) {
	b := []byte(`{
	"Name": "player 3",
	"Type": 0,
	"Company": 3,
	"Amount": 15,
	"Price": 3250
}`)

	pm, err := NewJournalTradeFromBytes(b)
	require.NoError(t, err)
	require.Equal(t, JournalTrade, pm.Type())

	trade := &Trade{
		Name:    "player 3",
		Type:    TradeBuy,
		Company: 3,
		Amount:  15,
		Price:   3250,
	}
	require.Equal(t, trade, pm.Payload())
}
