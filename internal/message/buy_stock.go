package message

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

// buyStockMessage represents a Buy message.
type buyStockMessage struct {
	company int
	amount  int
}

// NewBuyStockMessageFromBytes returns a new Message of Buy kind.
func NewBuyStockMessageFromBytes(b []byte) (Message, error) {
	split := strings.SplitN(string(b), separator, 2)
	if len(split) != 2 {
		return nil, fmt.Errorf("invalid buy stock message: %s", string(b))
	}

	company, err := strconv.Atoi(split[0])
	if err != nil {
		return nil, fmt.Errorf("invalid buy stock message, parse company: %s, error: %w", split[0], err)
	}
	amount, err := strconv.Atoi(split[1])
	if err != nil {
		return nil, fmt.Errorf("invalid buy stock message, parse amount: %s, error: %w", split[1], err)
	}

	return buyStockMessage{
		company: company,
		amount:  amount,
	}, nil
}

// Type implements the Message interface.
func (m buyStockMessage) Type() Kind {
	return Buy
}

// Payload implements the Message interface.
func (m buyStockMessage) Payload() any {
	return []int{m.company, m.amount}
}

// MarshalJSON implements the json.Marshaler interface.
func (m buyStockMessage) MarshalJSON() ([]byte, error) {
	b := base{
		Kind: Buy,
		Data: []byte(fmt.Sprintf("%d%s%d", m.company, separator, m.amount)),
	}
	return json.Marshal(b)
}
