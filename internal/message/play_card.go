package message

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

// playCardMessage represents a PlayCard kind message.
// It contains the ID of the card and an optional company ID if the card has a wildcard company.
type playCardMessage struct {
	id      int
	company int
}

// NewPlayCardFromBytes returns a new Message of PlayCard kind.
func NewPlayCard(id, company int) Message {
	return playCardMessage{
		id:      id,
		company: company,
	}
}

// NewPlayCardFromBytes returns a new Message of PlayCard kind.
func NewPlayCardFromBytes(b []byte) (Message, error) {
	split := strings.SplitN(string(b), separator, 2)
	if len(split) != 2 {
		return nil, fmt.Errorf("invalid play card message: %s", string(b))
	}

	id, err := strconv.Atoi(split[0])
	if err != nil {
		return nil, fmt.Errorf("invalid play card message, parse card id: %s, error: %w", split[0], err)
	}
	company, err := strconv.Atoi(split[1])
	if err != nil {
		return nil, fmt.Errorf("invalid play card message, parse company: %s, error: %w", split[1], err)
	}

	return playCardMessage{
		id:      id,
		company: company,
	}, nil
}

// Type implements the Message interface.
func (m playCardMessage) Type() Kind {
	return PlayCard
}

// Payload implements the Message interface.
func (m playCardMessage) Payload() any {
	return []int{m.id, m.company}
}

// MarshalJSON implements the json.Marshaler interface.
func (m playCardMessage) MarshalJSON() ([]byte, error) {
	b := base{
		Kind: PlayCard,
		Data: []byte(fmt.Sprintf("%d%s%d", m.id, separator, m.company)),
	}
	return json.Marshal(b)
}
