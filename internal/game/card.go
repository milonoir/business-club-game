package game

import (
	"fmt"
	"strconv"
	"strings"
)

// Card represents a player or bank action card.
type Card struct {
	ID   int        `json:"id"`
	Mods []Modifier `json:"mods"`
}

// Modifier represents a stock price modifier.
type Modifier struct {
	Company int  `json:"company"`
	Mod     *Mod `json:"mod"`
}

// Mod is a value modifier.
type Mod struct {
	op    string
	value int
}

// NewMod creates a new Mod.
func NewMod(op string, value int) *Mod {
	return &Mod{
		op:    op,
		value: value,
	}
}

// Calculate applies the Mod on p.
func (m *Mod) Calculate(p int) int {
	switch m.op {
	case "+":
		return p + m.value
	case "-":
		return p - m.value
	case "*":
		return p * m.value
	case "=":
		return m.value
	default:
		// This should not be possible.
		return 0
	}
}

// Op returns the Mod's operator.
func (m *Mod) Op() string {
	return m.op
}

// Value return the Mod's operand value.
func (m *Mod) Value() int {
	return m.value
}

// String implements the fmt.Stringer interface.
func (m *Mod) String() string {
	return fmt.Sprintf("%s %d", m.op, m.value)
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (m *Mod) UnmarshalJSON(b []byte) error {
	s := string(b)
	parts := strings.Split(strings.Trim(s, `"`), " ")
	if len(parts) != 2 {
		return fmt.Errorf("mod must consist of two parts separated by one space: %s", s)
	}

	value, err := strconv.Atoi(parts[1])
	if err != nil {
		return fmt.Errorf("mod value must be an integer: %s", parts[1])
	}
	m.value = value

	switch parts[0] {
	case "+", "-", "*", "=":
		m.op = parts[0]
	default:
		return fmt.Errorf("invalid mod definition: %s", s)
	}

	return nil
}

// MarshalJSON implements the json.Marshaler interface.
func (m *Mod) MarshalJSON() ([]byte, error) {
	return []byte(`"` + m.String() + `"`), nil
}
