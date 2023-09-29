package message

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewJournalDeal(t *testing.T) {
	deal := &Deal{
		Name:    "player 3",
		Type:    DealBuy,
		Company: 3,
		Amount:  15,
		Price:   3250,
	}

	m := NewJournalDeal(deal)

	b, err := json.Marshal(m)
	require.NoError(t, err)

	pm, err := Parse(b)
	require.NoError(t, err)

	require.Equal(t, JournalDeal, pm.Type())
	require.Equal(t, deal, pm.Payload())
}

func TestNewJournalDealFromBytes(t *testing.T) {
	b := []byte(`{
	"Name": "player 3",
	"Type": 0,
	"Company": 3,
	"Amount": 15,
	"Price": 3250
}`)

	pm, err := NewJournalDealFromBytes(b)
	require.NoError(t, err)
	require.Equal(t, JournalDeal, pm.Type())

	deal := &Deal{
		Name:    "player 3",
		Type:    DealBuy,
		Company: 3,
		Amount:  15,
		Price:   3250,
	}
	require.Equal(t, deal, pm.Payload())
}
