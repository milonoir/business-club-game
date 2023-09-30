package server

import (
	"log/slog"
	"math/rand"
	"time"

	"github.com/milonoir/business-club-game/internal/game"
	"github.com/milonoir/business-club-game/internal/message"
)

const (
	retryAttempts = 5
	retryDelay    = 2 * time.Second
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
		p.AddCash(game.StartingCash)
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
			// TODO: consider player timeout and automated actions.

			// Send state update to all players.
			g.sendStateUpdate(order, turn, i, false)

			// Signal player to start turn.
			p, _ := g.players.get(key)
			g.sendStartTurn(p)

			// Get player action card.
			g.handlePlayerAction(inbox, done, key, p)
			g.sendStateUpdate(order, turn, i, false)

			// Get player transaction. Send state update after each transaction.
			for g.handlePlayerTransaction(inbox, done, key, p) {
				g.sendStateUpdate(order, turn, i, false)
			}
		}

		// Update after last player in turn happens here.
		g.sendStateUpdate(order, turn, game.MaxPlayers+1, false)

		// Perform bank action.
		g.playCard("", g.assets.BankDeck[turn-1], game.WildcardCompany)
		g.sendStateUpdate(order, turn, game.MaxPlayers+1, false)
	}

	// Game has ended, calculate final scores, and send them to players.
	g.sendStateUpdate(nil, game.MaxTurns, game.MaxPlayers+1, true)
}

func (g *gameRunner) handlePlayerAction(inbox <-chan signedMessage, done <-chan struct{}, key string, p Player) {
	for {
		select {
		case <-done:
			return
		case msg := <-inbox:
			// Ignore messages from other players.
			if msg.Key != key {
				continue
			}

			// Ignore non-action messages.
			if msg.Msg.Type() != message.PlayCard {
				continue
			}

			// Validate action by card ID.
			pl := msg.Msg.Payload().([]int)
			id, company := pl[0], pl[1]

			for i, c := range p.Hand() {
				if c.ID == id {
					// Remove card from player hand.
					p.SetHand(append(p.Hand()[:i], p.Hand()[i+1:]...))

					// Play the card.
					g.playCard(p.Name(), c, company)

					// Acknowledge action.
					if err := retry(retryAttempts, retryDelay, func() error {
						return p.Conn().Send(message.NewPlayCard(0, 0))
					}); err != nil {
						g.l.Error("send play card ack", "error", err, "remote_addr", p.Conn().RemoteAddress())
					}

					// TODO: What if player disconnects?
				}
			}
		}
	}
}

func (g *gameRunner) playCard(playerName string, c *game.Card, company int) {
	for _, m := range c.Mods {
		cpny := m.Company
		if cpny <= game.WildcardCompany {
			cpny = company
		}

		// This should not ever happen, but log it just in case.
		if cpny < 0 || cpny > 3 {
			g.l.Error("asset error: invalid company", "chosen_company", company, "card_company", m.Company, "card_id", c.ID)
			continue
		}

		newPrice := m.Mod.Calculate(g.stockPrices[cpny])
		if newPrice < 0 {
			newPrice = 0
		} else if newPrice > game.MaxPrice {
			newPrice = game.MaxPrice
		}
		g.stockPrices[cpny] = newPrice

		// Send journal message.
		actor := message.ActorBank
		if playerName == "" {
			actor = message.ActorPlayer
		}
		action := &message.Action{
			ActorType: actor,
			Name:      playerName,
			Mod:       &m,
			NewPrice:  newPrice,
		}
		g.players.forEach(func(p Player) {
			if err := p.Conn().Send(message.NewJournalAction(action)); err != nil {
				g.l.Error("send journal action", "error", err, "remote_addr", p.Conn().RemoteAddress())
			}
		})
	}
}

func (g *gameRunner) handlePlayerTransaction(inbox <-chan signedMessage, done <-chan struct{}, key string, p Player) bool {
	for {
		select {
		case <-done:
			return false
		case msg := <-inbox:
			// Ignore messages from other players.
			if msg.Key != key {
				continue
			}

			switch msg.Msg.Type() {
			case message.EndTurn:
				// Player ends turn.
				return false
			case message.Buy:
				// Player wants to buy stocks.
				pl := msg.Msg.Payload().([]int)
				company, amount := pl[0], pl[1]
				g.playerBuyStocks(p, company, amount)
				return true
			case message.Sell:
				// Player wants to sell stocks.
				pl := msg.Msg.Payload().([]int)
				company, amount := pl[0], pl[1]
				g.playerSellStocks(p, company, amount)
				return true
			}
		}
	}
}

func (g *gameRunner) playerBuyStocks(p Player, company int, amount int) {
	if cost := g.stockPrices[company] * amount; p.Cash() >= cost {
		p.AddCash(-cost)
		p.AddStocks(company, amount)

		// Send journal message.
		deal := &message.Deal{
			Name:    p.Name(),
			Type:    message.DealBuy,
			Company: company,
			Amount:  amount,
			Price:   g.stockPrices[company],
		}
		g.players.forEach(func(pp Player) {
			if err := pp.Conn().Send(message.NewJournalDeal(deal)); err != nil {
				g.l.Error("send journal deal", "error", err, "remote_addr", pp.Conn().RemoteAddress())
			}
		})
	}
}

func (g *gameRunner) playerSellStocks(p Player, company int, amount int) {
	if p.Stocks()[company] >= amount {
		price := g.stockPrices[company]
		p.AddCash(price * amount)
		p.AddStocks(company, -amount)

		// Send journal message.
		deal := &message.Deal{
			Name:    p.Name(),
			Type:    message.DealSell,
			Company: company,
			Amount:  amount,
			Price:   price,
		}
		g.players.forEach(func(pp Player) {
			if err := pp.Conn().Send(message.NewJournalDeal(deal)); err != nil {
				g.l.Error("send journal deal", "error", err, "remote_addr", pp.Conn().RemoteAddress())
			}
		})
	}
}

func (g *gameRunner) shufflePlayers() []string {
	order := g.players.keys()
	rand.Shuffle(len(order), func(i, j int) {
		order[i], order[j] = order[j], order[i]
	})

	playerOrder := make([]string, 0, len(order))
	for _, key := range order {
		p, _ := g.players.get(key)
		playerOrder = append(playerOrder, p.Name())
	}

	return playerOrder
}

func (g *gameRunner) sendStartTurn(p Player) {
	if err := retry(retryAttempts, retryDelay, func() error {
		return p.Conn().Send(message.NewStartTurn())
	}); err != nil {
		g.l.Error("send start turn", "error", err, "remote_addr", p.Conn().RemoteAddress())
	}

	// TODO: What if player disconnects?
}

func (g *gameRunner) sendStateUpdate(order []string, turn, currentPlayer int, isFinal bool) {
	state := &message.GameState{
		Started:       true,
		Ended:         isFinal,
		Companies:     g.assets.Companies,
		StockPrices:   g.stockPrices,
		Turn:          turn,
		PlayerOrder:   order,
		CurrentPlayer: currentPlayer,
	}

	// Build player states first.
	// TODO: Maybe calculate non-final opponent data once and reuse it.
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
				// Never include opponent hand.
				o := game.Player{Name: p.Name}

				if isFinal {
					// Only include cash and stocks if game is over.
					o.Cash = p.Cash
					o.Stocks = p.Stocks
				} else {
					// Otherwise, only include "levels" of wealth.
					o.Cash = game.CashLevel(p.Cash)
					o.Stocks = [4]int{game.StockLevel(p.Stocks[0]), game.StockLevel(p.Stocks[1]), game.StockLevel(p.Stocks[2]), game.StockLevel(p.Stocks[3])}
				}

				opps = append(opps, o)
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
