package main

import (
	"math/rand"
)

const (
	maxTurns      = 10
	startingPrice = 150
)

type assets struct {
	Companies  []string `json:"companies"`
	ActionDeck []card   `json:"actionDeck"`
	BankDeck   []card   `json:"bankDeck"`
}

type game struct {
	StockPrices map[string]int
	Players     []player
	TurnsLeft   int

	assets
}

func (g *game) ShufflePlayers() {
	rand.Shuffle(len(g.Players), func(i, j int) {
		g.Players[i], g.Players[j] = g.Players[j], g.Players[i]
	})
}

func newGame(a assets) game {
	g := game{
		assets:      a,
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
