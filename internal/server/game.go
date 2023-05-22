package server

import (
	"math/rand"

	"github.com/milonoir/business-club-game/internal/game"
)

const (
	maxTurns      = 10
	startingPrice = 150
	maxPrice      = 400
)

type Assets struct {
	Companies  []string    `json:"companies"`
	ActionDeck []game.Card `json:"actionDeck"`
	BankDeck   []game.Card `json:"bankDeck"`
}

type Game struct {
	StockPrices map[string]int
	Players     []Player
	TurnsLeft   int

	Assets
}

func NewGame(a Assets) Game {
	g := Game{
		Assets:      a,
		TurnsLeft:   maxTurns,
		StockPrices: make(map[string]int, len(a.Companies)),
	}

	// Set starting price for each company.
	for _, company := range a.Companies {
		g.StockPrices[company] = startingPrice
	}

	// Shuffle action and bank decks.
	rand.Shuffle(len(g.ActionDeck), func(i, j int) {
		g.ActionDeck[i], g.ActionDeck[j] = g.ActionDeck[j], g.ActionDeck[i]
	})
	rand.Shuffle(len(g.BankDeck), func(i, j int) {
		g.BankDeck[i], g.BankDeck[j] = g.BankDeck[j], g.BankDeck[i]
	})

	return g
}

func (g *Game) shufflePlayers() {
	rand.Shuffle(len(g.Players), func(i, j int) {
		g.Players[i], g.Players[j] = g.Players[j], g.Players[i]
	})
}

func (g *Game) start() {
	// Game ends after maxTurns.
	for turn := 0; turn < maxTurns; turn++ {
		// Shuffle players, then iterate over players.
		g.shufflePlayers()
		for _, p := range g.Players {
			_ = p
			// Get player action card. Validate action. Retry if invalid.
			// Update game state.
			// Get player transaction. Validate transaction. Retry if invalid.
		}
		// Perform bank action.
		// Update game state.
	}
}

func (g *Game) getPlayerCard(p *Player) *game.Card {
	return nil
}

func (g *Game) drawBankCard() *game.Card {
	return nil
}

func (g *Game) update(c *game.Card) {

}
