package main

import (
	"fmt"
	"strconv"
	"strings"
)

type card struct {
	ID   int        `json:"id"`
	Mods []modifier `json:"mods"`
}

type modifier struct {
	Company int `json:"company"`
	Mod     mod `json:"mod"`
}

type mod func(int) int

func (m *mod) UnmarshalJSON(b []byte) error {
	// Playground example:
	// https://go.dev/play/p/ZCxkpdGkOFF

	s := string(b)
	parts := strings.Split(strings.Trim(s, `"`), " ")
	if len(parts) != 2 {
		return fmt.Errorf("mod must consist of two parts separated by one space: %s", s)
	}
	value, err := strconv.Atoi(parts[1])
	if err != nil {
		return fmt.Errorf("mod value must be an integer: %s", parts[1])
	}
	switch parts[0] {
	case "+":
		*m = func(p int) int {
			return p + value
		}
	case "-":
		*m = func(p int) int {
			return p - value
		}
	case "*":
		*m = func(p int) int {
			return p * value
		}
	case "=":
		*m = func(p int) int {
			return value
		}
	default:
		return fmt.Errorf("invalid mod definition: %s", s)
	}

	return nil
}
