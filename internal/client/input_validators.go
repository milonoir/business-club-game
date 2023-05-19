package client

import (
	"strconv"
)

func PositiveIntegerValidator(max int) func(text string, ch rune) bool {
	return func(text string, ch rune) bool {
		i, err := strconv.Atoi(text)
		if err != nil {
			return false
		}
		if i < 1 {
			return false
		}
		if i > max {
			return false
		}
		return true
	}
}
