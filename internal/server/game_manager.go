package server

import (
	"math/rand"

	"github.com/milonoir/business-club-game/internal/common"
	"github.com/milonoir/business-club-game/internal/game"
)

type gameManager struct {
	players     map[string]Player
	messages    chan signedMessage
	stockPrices [4]int
	assets      *game.Assets
}

func (g *gameManager) init() {
	// Reset stock prices.
	g.stockPrices = [4]int{common.StartingPrice, common.StartingPrice, common.StartingPrice, common.StartingPrice}

	// Shuffle action and bank decks.
	g.assets.ShuffleActionDeck()
	g.assets.ShuffleBankDeck()

	// TODO: Deal cards to players.
}

func (g *gameManager) run() {
	// Game ends after maxTurns.
	for turn := 0; turn < common.MaxTurns; turn++ {
		// Shuffle players, then iterate over players.
		g.shufflePlayers()
		for _, p := range g.players {
			_ = p
			// Get player action card. Validate action. Retry if invalid.
			// Update game state.
			// Get player transaction. Validate transaction. Retry if invalid.
		}
		// Perform bank action.
		// Update game state.
	}
}

func (g *gameManager) shufflePlayers() []string {
	var order []string

	for key := range g.players {
		order = append(order, key)
	}

	rand.Shuffle(len(order), func(i, j int) {
		order[i], order[j] = order[j], order[i]
	})

	return order
}
