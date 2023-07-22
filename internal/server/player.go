package server

import (
	"github.com/milonoir/business-club-game/internal/game"
	"github.com/milonoir/business-club-game/internal/network"
)

type Player struct {
	// Networking properties
	conn  network.Connection
	key   string
	ready bool

	// Game properties
	cash    int
	stocks  map[string]int
	actions []game.Card
}
