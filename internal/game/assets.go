package game

import (
	"math/rand"
)

type Assets struct {
	Companies  []string `json:"companies"`
	ActionDeck []*Card  `json:"actionDeck"`
	BankDeck   []*Card  `json:"bankDeck"`
}

func (a *Assets) ShuffleActionDeck() {
	rand.Shuffle(len(a.ActionDeck), func(i, j int) {
		a.ActionDeck[i], a.ActionDeck[j] = a.ActionDeck[j], a.ActionDeck[i]
	})
}

func (a *Assets) ShuffleBankDeck() {
	rand.Shuffle(len(a.BankDeck), func(i, j int) {
		a.BankDeck[i], a.BankDeck[j] = a.BankDeck[j], a.BankDeck[i]
	})
}
