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
	case cash < 100:
		return 1
	case cash < 1_000:
		return 2
	case cash < 10_000:
		return 3
	case cash < 100_000:
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
	case amount < 5:
		return 1
	case amount < 50:
		return 2
	case amount < 500:
		return 3
	case amount < 5_000:
		return 4
	default:
		return 5
	}
}
