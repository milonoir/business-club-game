package message

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

// TradeType represents the type of trade.
type TradeType uint8

func (t TradeType) AsString() string {
	switch t {
	case TradeBuy:
		return "Buy"
	case TradeSell:
		return "Sell"
	default:
		return "-"
	}
}

const (
	TradeBuy TradeType = iota
	TradeSell
)

// tradeStockMessage is a Message of TradeStock kind.
type tradeStockMessage struct {
	trade   TradeType
	company int
	amount  int
}

// NewTradeStock returns a new Message of TradeStock kind.
func NewTradeStock(trade TradeType, company int, amount int) Message {
	return tradeStockMessage{
		trade:   trade,
		company: company,
		amount:  amount,
	}
}

// NewTradeStockFromBytes returns a new Message of TradeStock kind.
func NewTradeStockFromBytes(b []byte) (Message, error) {
	split := strings.SplitN(string(b), separator, 3)
	if len(split) != 3 {
		return nil, fmt.Errorf("invalid trade stock message: %s", string(b))
	}

	trade, err := strconv.Atoi(split[0])
	if err != nil {
		return nil, fmt.Errorf("invalid trade stock message, parse trade type: %s, error: %w", split[0], err)
	}
	company, err := strconv.Atoi(split[1])
	if err != nil {
		return nil, fmt.Errorf("invalid trade stock message, parse company: %s, error: %w", split[1], err)
	}
	amount, err := strconv.Atoi(split[2])
	if err != nil {
		return nil, fmt.Errorf("invalid trade stock message, parse amount: %s, error: %w", split[2], err)
	}

	return tradeStockMessage{
		trade:   TradeType(trade),
		company: company,
		amount:  amount,
	}, nil
}

// Type implements the Message interface.
func (m tradeStockMessage) Type() Kind {
	return TradeStock
}

// Payload implements the Message interface.
func (m tradeStockMessage) Payload() any {
	return []any{m.trade, m.company, m.amount}
}

// MarshalJSON implements the json.Marshaler interface.
func (m tradeStockMessage) MarshalJSON() ([]byte, error) {
	b := base{
		Kind: TradeStock,
		Data: []byte(fmt.Sprintf("%d%s%d%s%d", m.trade, separator, m.company, separator, m.amount)),
	}
	return json.Marshal(b)
}
