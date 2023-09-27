package message

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

// sellStockMessage represents a Sell message.
type sellStockMessage struct {
	company int
	amount  int
}

// NewSellStockMessageFromBytes returns a new Message of Sell kind.
func NewSellStockMessageFromBytes(b []byte) (Message, error) {
	split := strings.SplitN(string(b), separator, 2)
	if len(split) != 2 {
		return nil, fmt.Errorf("invalid sell stock message: %s", string(b))
	}

	company, err := strconv.Atoi(split[0])
	if err != nil {
		return nil, fmt.Errorf("invalid sell stock message, parse company: %s, error: %w", split[0], err)
	}
	amount, err := strconv.Atoi(split[1])
	if err != nil {
		return nil, fmt.Errorf("invalid sell stock message, parse amount: %s, error: %w", split[1], err)
	}

	return sellStockMessage{
		company: company,
		amount:  amount,
	}, nil
}

// Type implements the Message interface.
func (m sellStockMessage) Type() Kind {
	return Sell
}

// Payload implements the Message interface.
func (m sellStockMessage) Payload() any {
	return []int{m.company, m.amount}
}

// MarshalJSON implements the json.Marshaler interface.
func (m sellStockMessage) MarshalJSON() ([]byte, error) {
	b := base{
		Kind: Sell,
		Data: []byte(fmt.Sprintf("%d%s%d", m.company, separator, m.amount)),
	}
	return json.Marshal(b)
}
