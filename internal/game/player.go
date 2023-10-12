package game

// Player holds player data.
type Player struct {
	Name   string
	Cash   int
	Stocks [4]int
	Hand   []*Card
}

// CashLevel returns the level for the given cash amount.
func CashLevel(cash int) int {
	switch {
	case cash < 1:
		return 0
	case cash < 1_001:
		return 1
	case cash < 10_001:
		return 2
	case cash < 100_001:
		return 3
	case cash < 1_000_001:
		return 4
	default:
		return 5
	}
}

// StockLevel returns the level for the given stock amount.
func StockLevel(amount int) int {
	switch {
	case amount < 1:
		return 0
	case amount < 11:
		return 1
	case amount < 101:
		return 2
	case amount < 1_001:
		return 3
	case amount < 10_001:
		return 4
	default:
		return 5
	}
}
