package client

import (
	"github.com/gdamore/tcell/v2"
	"github.com/milonoir/business-club-game/internal/client/ui"
	"github.com/milonoir/business-club-game/internal/game"
	"github.com/milonoir/business-club-game/internal/message"
)

func (a *Application) handleKeyExchange(data *ui.LoginData, msg []string) {
	if key := msg[0]; key == "" {
		if err := a.server.Send(message.NewKeyExchangeWithName(data.ReconnectKey, data.Username)); err != nil {
			a.l.Error("send key exchange", "error", err)
		}
	} else {
		a.l.Info("received reconnect key", "key", key)
		// Update server status widget.
		a.srvStatus.SetReconnectKey(key)

		// Update login form.
		a.login.SetReconnectKey(key)
	}
}

func (a *Application) handleStateUpdate(state *message.GameState) {
	// Safety check.
	if state == nil {
		return
	}

	// Readiness update.
	if !state.Started {
		if a.gameStarted.Load() {
			a.gameStarted.Store(false)
			a.pages.SwitchToPage(lobbyPageName)
		}

		if len(state.Readiness) > 0 {
			a.lobby.Update(state.Readiness)
		}
		return
	}

	// Switch to main page if game started.
	if !a.gameStarted.Load() {
		a.gameStarted.Store(true)
		a.pages.SwitchToPage(gamePageName)
	}

	// Update CompanyProvider.
	// This should not change during the game, but it is useful when a client is reconnecting.
	a.cp.SetCompanies(state.Companies)

	// Update local game data.
	a.hand = state.Player.Hand
	a.cash = state.Player.Cash
	a.stocks = state.Player.Stocks
	a.prices = state.StockPrices

	// Update UI - turn.
	a.turn.Update(game.MaxTurns, state.Turn, state.PlayerOrder, state.CurrentPlayer)

	// Update UI - standings.
	a.standings.Update(state)

	// Update UI - graph.
	a.graph.Add(state.StockPrices)
}

func (a *Application) handleJournalAction(msg *message.Action) {
	a.history.AddAction(msg)
}

func (a *Application) handleJournalTrade(msg *message.Trade) {
	a.history.AddTrade(msg)
}

func (a *Application) handleStartTurn(phase game.TurnPhase) {
	switch phase {
	case game.ActionPhase:
		a.l.Info("starting action phase")
		a.turnActionPhase()
	case game.TradePhase:
		a.l.Info("starting trade phase")
		a.turnTradePhase()
	}
}

func (a *Application) turnActionPhase() {
	// Remove wait panel, while action list is displayed.
	a.bottomRow.RemoveItem(a.wait.GetTextView())
	defer a.bottomRow.AddItem(a.wait.GetTextView(), 0, 2, 1, 1, 1, 1, false)

	actCh := make(chan *game.Card)
	defer close(actCh)

	action := ui.NewActionList(a.cp, a.hand, func(card *game.Card) {
		actCh <- card
	})
	a.bottomRow.AddItem(action.GetList(), 0, 2, 1, 1, 1, 1, true)
	a.app.SetFocus(action.GetList())

	// Make sure action list is removed.
	defer a.bottomRow.RemoveItem(action.GetList())

	// Sync point.
	card := <-actCh
	id := card.ID

	// If card has a wildcard company, player must select a company.
	company := game.WildcardCompany
	if card.Mods[0].Company == game.WildcardCompany || card.Mods[1].Company == game.WildcardCompany {
		compCh := make(chan int)
		defer close(compCh)

		cl := ui.NewCompanyList(a.cp, func(i int) {
			compCh <- i
		})
		cl.GetList().SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			if event.Key() == tcell.KeyEscape {
				// Cancel the action by setting id to -1.
				id = -1
				compCh <- game.WildcardCompany
			}
			return event
		})
		a.bottomRow.AddItem(cl.GetList(), 0, 3, 1, 1, 1, 1, true)
		a.app.SetFocus(cl.GetList())

		// Sync point.
		company = <-compCh

		// Remove company list.
		a.bottomRow.RemoveItem(cl.GetList())
	}

	//confirmCh := make(chan bool)
	//defer close(confirmCh)
	//
	//cf := ui.NewConfirmPanel("Are you sure?", confirmCh)
	//a.bottomRow.AddItem(cf.GetModal(), 0, 3, 1, 1, 1, 1, true)
	//
	//// Sync point.
	//confirmed := <-confirmCh
	//
	//// Remove confirm modal.
	//a.bottomRow.RemoveItem(cf.GetModal())
	//
	//if !confirmed {
	//	// Cancel the action by sending a wildcard company to the channel.
	//	company = game.WildcardCompany
	//}

	// Send action to server.
	a.l.Info("sending action", "card", id, "company", company)
	if err := a.server.Send(message.NewPlayCard(id, company)); err != nil {
		a.l.Error("send action", "error", err)
	}
}

func (a *Application) turnTradePhase() {
	// Remove wait panel, while action list is displayed.
	a.bottomRow.RemoveItem(a.wait.GetTextView())
	defer a.bottomRow.AddItem(a.wait.GetTextView(), 0, 2, 1, 1, 1, 1, false)

	optCh := make(chan ui.TradeOption)
	defer close(optCh)

	tm := ui.NewTradeMenu(func(option ui.TradeOption) {
		optCh <- option
	})
	a.bottomRow.AddItem(tm.GetList(), 0, 2, 1, 1, 1, 1, true)
	a.app.SetFocus(tm.GetList())

	// Make sure trade menu is removed.
	defer a.bottomRow.RemoveItem(tm.GetList())

	// Sync point.
	option := <-optCh

	// End turn.
	if option == ui.EndTurn {
		a.bottomRow.RemoveItem(tm.GetList())
		if err := a.server.Send(message.NewEndTurn()); err != nil {
			a.l.Error("send end turn", "error", err)
		}
		return
	}

	// Select company.
	compCh := make(chan int)
	defer close(compCh)

	cl := ui.NewCompanyList(a.cp, func(i int) {
		compCh <- i
	})
	a.bottomRow.AddItem(cl.GetList(), 0, 3, 1, 1, 1, 1, true)
	a.app.SetFocus(cl.GetList())

	// Sync point.
	company := <-compCh

	// Remove company list.
	a.bottomRow.RemoveItem(cl.GetList())

	// Calculate the maximum amount of stocks that can be traded.
	var (
		tradeType message.TradeType
		maximum   int
	)
	switch option {
	case ui.Buy:
		tradeType = message.TradeBuy
		// Buying stocks is only possible if the price is greater than 0.
		if a.prices[company] > 0 {
			// This is an integer division, so the result is floored.
			maximum = a.cash / a.prices[company]
		}
	case ui.Sell:
		tradeType = message.TradeSell
		// Selling stocks is only possible if the price is greater than 0.
		if a.prices[company] > 0 {
			maximum = a.stocks[company]
		}
	}

	// Player has not enough money or stocks.
	// Send a trade with 0 amount.
	if maximum == 0 {
		if err := a.server.Send(message.NewTradeStock(tradeType, company, 0)); err != nil {
			a.l.Error("send cancelled trade", "error", err)
		}
		return
	}

	// Type in amount.
	amountCh := make(chan int)
	defer close(amountCh)

	form := ui.NewTradeForm(a.cp, tradeType, company, maximum, func(amount int) {
		amountCh <- amount
	})
	a.bottomRow.AddItem(form.GetForm(), 0, 3, 1, 1, 1, 1, true)
	a.app.SetFocus(form.GetForm())

	// Sync point.
	amount := <-amountCh

	// Remove trade form.
	a.bottomRow.RemoveItem(form.GetForm())

	// Send trade to server.
	if err := a.server.Send(message.NewTradeStock(tradeType, company, amount)); err != nil {
		a.l.Error("send trade", "error", err)
	}
}
