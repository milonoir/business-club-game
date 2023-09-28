package message

import (
	"encoding/json"
	"fmt"
)

// base is used for two things:
// - parsing message kinds;
// - and marshaling messages.
type base struct {
	Kind Kind
	Data []byte
}

// Type implements the Message interface.
func (mb base) Type() Kind {
	return mb.Kind
}

// Payload implements the Message interface.
func (mb base) Payload() any {
	return mb.Data
}

// Parse returns a Message from a raw []byte.
func Parse(raw []byte) (Message, error) {
	var b base
	if err := json.Unmarshal(raw, &b); err != nil {
		return nil, fmt.Errorf("parse raw message: %w", err)
	}

	switch b.Type() {
	case Error:
		return NewError(string(b.Data)), nil
	case KeyExchange:
		return NewKeyExchangeFromBytes(b.Data), nil
	case StateUpdate:
		return NewStateUpdateFromBytes(b.Data)
	case VoteToStart:
		return NewVoteToStartFromBytes(b.Data), nil
	case StartTurn:
		return NewStartTurn(), nil
	case EndTurn:
		return NewEndTurn(), nil
	case PlayCard:
		return NewPlayCardFromBytes(b.Data)
	case Buy:
		return NewBuyStockMessageFromBytes(b.Data)
	case Sell:
		return NewSellStockMessageFromBytes(b.Data)
	case QueryCompanies:
		return NewQueryCompaniesFromBytes(b.Data)
	case QueryPlayers:
		return NewUnknown(), nil
	default:
		return NewUnknown(), nil
	}
}
