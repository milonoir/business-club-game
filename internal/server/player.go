package server

import (
	"github.com/milonoir/business-club-game/internal/game"
)

type Player struct {
	// Networking properties
	conn  *connection
	key   string
	ready bool

	// Game properties
	cash    int
	stocks  map[string]int
	actions []game.Card
}
