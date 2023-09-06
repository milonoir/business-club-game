package game

// Player holds player data.
type Player struct {
	Name   string
	Cash   int
	Stocks [4]int
	Hand   []*Card
}
