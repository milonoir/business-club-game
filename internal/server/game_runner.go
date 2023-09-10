package server

import (
	"log/slog"
	"math/rand"

	"github.com/milonoir/business-club-game/internal/game"
	"github.com/milonoir/business-club-game/internal/message"
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

		for i, key := range order {
			// Send state update to all players.
			g.sendStateUpdate(order, turn, i)

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

func (g *gameRunner) sendStateUpdate(order []string, turn, currentPlayer int) {
	state := &message.GameState{
		Started:       true,
		Companies:     g.assets.Companies,
		StockPrices:   g.stockPrices,
		Turn:          turn,
		PlayerOrder:   order,
		CurrentPlayer: currentPlayer,
	}

	// Build player states first.
	playerStates := make(map[string]game.Player, game.MaxPlayers)
	keys := g.players.keys()
	for _, key := range keys {
		p, _ := g.players.get(key)
		playerStates[key] = game.Player{
			Name:   p.Name(),
			Cash:   p.Cash(),
			Stocks: p.Stocks(),
			Hand:   p.Hand(),
		}
	}

	// Reuse game state.
	for _, key := range keys {
		// Separate player and opponents.
		state.Player = playerStates[key]
		opps := make([]game.Player, 0, game.MaxPlayers-1)
		for k, p := range playerStates {
			if k != key {
				opps = append(opps, p)
			}
		}
		state.Opponents = opps

		// Send state update to player.
		p, _ := g.players.get(key)
		if p.Conn().IsAlive() {
			if err := p.Conn().Send(message.NewStateUpdate(state)); err != nil {
				g.l.Error("send game state update", "error", err, "remote_addr", p.Conn().RemoteAddress())
			}
		}
	}
}
