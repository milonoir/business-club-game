package message

import (
	"encoding/json"
	"fmt"
)

// TradeType represents the type of trade.
type TradeType uint8

const (
	TradeBuy TradeType = iota
	TradeSell
)

// Trade represents a trade journal item.
type Trade struct {
	Name    string
	Type    TradeType
	Company int
	Amount  int
	Price   int
}

// journalTradeMessage is a message that contains a trade journal.
type journalTradeMessage struct {
	trade *Trade
}

// NewJournalTrade creates a new journal trade message.
func NewJournalTrade(trade *Trade) Message {
	return journalTradeMessage{
		trade: trade,
	}
}

// NewJournalTradeFromBytes creates a new journal trade message from bytes.
func NewJournalTradeFromBytes(b []byte) (Message, error) {
	var trade Trade
	if err := json.Unmarshal(b, &trade); err != nil {
		return nil, fmt.Errorf("unmarshal trade: %w", err)
	}
	return journalTradeMessage{
		trade: &trade,
	}, nil
}

// Type implements the Message interface.
func (m journalTradeMessage) Type() Kind {
	return JournalTrade
}

// Payload implements the Message interface.
func (m journalTradeMessage) Payload() any {
	return m.trade
}

// MarshalJSON implements the json.Marshaler interface.
func (m journalTradeMessage) MarshalJSON() ([]byte, error) {
	tb, err := json.Marshal(m.trade)
	if err != nil {
		return nil, fmt.Errorf("marshal trade: %w", err)
	}

	b := base{
		Kind: JournalTrade,
		Data: tb,
	}
	return json.Marshal(b)
}
