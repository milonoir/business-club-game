package message

import (
	"encoding/json"
	"fmt"
)

// DealType represents the type of deal.
type DealType uint8

const (
	DealBuy DealType = iota
	DealSell
)

// Deal represents a deal journal item.
type Deal struct {
	Name    string
	Type    DealType
	Company int
	Amount  int
	Price   int
}

// journalDealMessage is a message that contains a deal journal.
type journalDealMessage struct {
	deal *Deal
}

// NewJournalDeal creates a new journal deal message.
func NewJournalDeal(deal *Deal) Message {
	return journalDealMessage{
		deal: deal,
	}
}

// NewJournalDealFromBytes creates a new journal deal message from bytes.
func NewJournalDealFromBytes(b []byte) (Message, error) {
	var deal Deal
	if err := json.Unmarshal(b, &deal); err != nil {
		return nil, fmt.Errorf("unmarshal deal: %w", err)
	}
	return journalDealMessage{
		deal: &deal,
	}, nil
}

// Type implements the Message interface.
func (m journalDealMessage) Type() Kind {
	return JournalDeal
}

// Payload implements the Message interface.
func (m journalDealMessage) Payload() any {
	return m.deal
}

// MarshalJSON implements the json.Marshaler interface.
func (m journalDealMessage) MarshalJSON() ([]byte, error) {
	db, err := json.Marshal(m.deal)
	if err != nil {
		return nil, fmt.Errorf("marshal deal: %w", err)
	}

	b := base{
		Kind: JournalDeal,
		Data: db,
	}
	return json.Marshal(b)
}
