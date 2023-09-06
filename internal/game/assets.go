package game

import (
	"math/rand"
)

type Assets struct {
	Companies  []string `json:"companies"`
	PlayerDeck []*Card  `json:"playerDeck"`
	BankDeck   []*Card  `json:"bankDeck"`
}

func (a *Assets) ShufflePlayerDeck() {
	rand.Shuffle(len(a.PlayerDeck), func(i, j int) {
		a.PlayerDeck[i], a.PlayerDeck[j] = a.PlayerDeck[j], a.PlayerDeck[i]
	})
}

func (a *Assets) ShuffleBankDeck() {
	rand.Shuffle(len(a.BankDeck), func(i, j int) {
		a.BankDeck[i], a.BankDeck[j] = a.BankDeck[j], a.BankDeck[i]
	})
}
