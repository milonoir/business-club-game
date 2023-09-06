package server

import (
	"log/slog"
	"math/rand"

	"github.com/milonoir/business-club-game/internal/game"
)

type gameRunner struct {
	// Input
	assets  *game.Assets
	players *playerMap

	// Game state.
	stockPrices [4]int

	l *slog.Logger
}

func newGameRunner(players *playerMap, assets *game.Assets) *gameRunner {
	return &gameRunner{
		players: players,
		assets:  assets,
	}
}

func (g *gameRunner) init() {
	// Reset stock prices.
	g.stockPrices = [4]int{game.StartingPrice, game.StartingPrice, game.StartingPrice, game.StartingPrice}

	// Shuffle player and bank decks.
	g.assets.ShufflePlayerDeck()
	g.assets.ShuffleBankDeck()

	// Deal cards to players.
	i := 0
	g.players.forEach(func(p Player) {
		p.SetHand(g.assets.PlayerDeck[i : i+game.MaxTurns])
		i += game.MaxTurns
	})
}

func (g *gameRunner) run(inbox chan signedMessage, done <-chan struct{}) {
	// Initialize game state.
	g.init()

	// Game ends after maxTurns.
	for turn := 1; turn < game.MaxTurns+1; turn++ {
		// Shuffle players, then iterate over players.
		order := g.shufflePlayers()

		for _, key := range order {
			p, _ := g.players.get(key)
			_ = p
			// Signal player to start turn.
			// Get player action card. Validate action. Retry if invalid.
			// Update game state.
			// Get player transaction. Validate transaction. Retry if invalid.
			// Update game state.
			// If player ends turn, move on to the next player.
		}

		// Perform bank action.
		// Update game state.
	}

	// Game has ended, calculate final scores, and send them to players.
}

func (g *gameRunner) handlePlayerAction(inbox <-chan signedMessage, done <-chan struct{}, key string, p Player) {
	for finished := false; !finished; {
		select {
		case <-done:
			return
		case msg := <-inbox:
			if msg.Key != key {
				continue
			}

			// TODO: Validate action.

			finished = true
		}
	}
}

func (g *gameRunner) shufflePlayers() []string {
	order := g.players.keys()

	rand.Shuffle(len(order), func(i, j int) {
		order[i], order[j] = order[j], order[i]
	})

	return order
}
