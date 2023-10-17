package client

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"time"

	"github.com/gobwas/ws"
	"github.com/milonoir/business-club-game/internal/client/ui"
	"github.com/milonoir/business-club-game/internal/game"
	"github.com/milonoir/business-club-game/internal/message"
	"github.com/milonoir/business-club-game/internal/network"
	"github.com/rivo/tview"
)

const (
	titlePageName = "titlePage"
	lobbyPageName = "lobbyPage"
	gamePageName  = "gamePage"
	errorPageName = "errorPage"
)

var (
	errConnectionClosed = errors.New("connection closed")
)

type Application struct {
	// Main application.
	app   *tview.Application
	pages *tview.Pages

	// Game screen.
	standings *ui.StandingsPanel
	turn      *ui.TurnPanel
	graph     *ui.GraphPanel
	history   *ui.HistoryPanel
	version   *ui.VersionPanel
	srvStatus *ui.ServerStatusPanel

	// The lower section of the game screen where player interactions are handled.
	bottomRow *tview.Grid
	actions   *ui.ActionDisplay

	// Lobby screen.
	lobby *ui.LobbyForm

	// Title screen.
	title *ui.TitlePanel
	login *ui.LoginForm

	// Company provider knows the matching colors and names for companies.
	cp *ui.CompanyProvider

	hand   []*game.Card
	cash   int
	stocks [4]int
	prices [4]int

	server network.Connection
	errCh  chan string
	l      *slog.Logger
}

func NewApplication(l *slog.Logger) *Application {
	a := &Application{
		cp:    ui.NewCompanyProvider(),
		l:     l,
		errCh: make(chan string, 10),
	}

	a.l.Info("initializing application")
	a.initUI()

	return a
}

func (a *Application) initUI() {
	// Create application.
	a.app = tview.NewApplication()

	// Create pages.
	a.pages = tview.NewPages()

	// -------------------------------------------------- Title page
	titlePage := tview.NewFlex()

	// Title panel.
	a.title = ui.NewTitlePanel()
	titlePage.AddItem(a.title.GetTextView(), 0, 3, false)

	// Login form.
	a.login = ui.NewLoginForm(
		func(data *ui.LoginData) {
			if err := a.connect(data); err != nil {
				// Making sure there is no dead connection stored.
				a.server = nil
				modal := ui.NewErrorModal(err)
				a.pages.AddPage(errorPageName, modal.GetModal(), true, true)
				modal.SetHandler(func(int, string) {
					a.pages.RemovePage(errorPageName)
				})
				return
			}
			// Successful login.
			a.lobby.Reset()
		},
		func() {
			a.Stop()
		},
	)
	titlePage.AddItem(a.login.GetForm(), 0, 1, true)

	// -------------------------------------------------- Lobby page
	lobbyPage := tview.NewFlex()

	// Title panel.
	lobbyPage.AddItem(a.title.GetTextView(), 0, 3, false)

	// Lobby form.
	a.lobby = ui.NewLobbyForm(
		func(ready bool) {
			_ = a.server.Send(message.NewVoteToStart(ready))
		},
		a.disconnect,
	)
	lobbyPage.AddItem(a.lobby.GetForm(), 0, 1, true)

	// -------------------------------------------------- Game page
	gamePage := tview.NewGrid().
		SetColumns(0, 0, 0).
		SetRows(0, 22, 0, 1, 1)

	// Top row of the game page.
	topRow := tview.NewGrid()
	topRow.
		SetColumns(0, 96, 0, 25, 0).
		SetRows(9)
	gamePage.AddItem(topRow, 0, 0, 1, 3, 1, 1, false)

	// Standings panel.
	a.standings = ui.NewStandingsPanel(a.cp)
	topRow.AddItem(a.standings.GetGrid(), 0, 1, 1, 1, 1, 1, false)

	// Turn widget.
	a.turn = ui.NewTurnPanel()
	topRow.AddItem(a.turn.GetTextView(), 0, 3, 1, 1, 1, 1, false)

	// Middle row of the game page.
	middleRow := tview.NewGrid()
	middleRow.
		SetColumns(96, 0).
		SetRows(22)
	gamePage.AddItem(middleRow, 1, 0, 1, 3, 1, 1, false)

	// Stock price graph panel.
	a.graph = ui.NewGraphPanel(a.cp)
	middleRow.AddItem(a.graph.GetGrid(), 0, 0, 1, 1, 1, 1, false)

	// History panel.
	a.history = ui.NewHistoryPanel(a.cp)
	middleRow.AddItem(a.history.GetTextView(), 0, 1, 1, 1, 1, 1, false)

	// Bottom row of the game page.
	a.bottomRow = tview.NewGrid()
	a.bottomRow.
		SetColumns(0, 0, 43, 0, 0).
		SetRows(9)
	gamePage.AddItem(a.bottomRow, 2, 0, 1, 3, 1, 1, true)

	// By default, display wait panel will be displayed.
	a.actions = ui.NewActionDisplay(a.cp)
	a.bottomRow.AddItem(a.actions.GetTextView(), 0, 2, 1, 1, 1, 1, true)

	// Game version widget.
	a.version = ui.NewVersionPanel()
	a.version.SetVersion(game.Version)
	gamePage.AddItem(a.version.GetTextView(), 4, 0, 1, 1, 1, 1, false)

	// Server status widget.
	a.srvStatus = ui.NewServerStatus()
	gamePage.AddItem(a.srvStatus.GetTextView(), 4, 1, 1, 2, 1, 1, false)

	// Setup pages and set pages as root.
	a.pages.
		AddPage(titlePageName, titlePage, true, true).
		AddPage(lobbyPageName, lobbyPage, true, false).
		AddPage(gamePageName, gamePage, true, false)

	a.app.SetRoot(a.pages, true)
}

func (a *Application) GetApplication() *tview.Application {
	return a.app
}

func (a *Application) Run() error {
	done := make(chan struct{})
	defer close(done)

	go a.redraw(done)

	return a.app.Run()
}

func (a *Application) redraw(done <-chan struct{}) {
	t := time.NewTicker(500 * time.Millisecond)
	defer t.Stop()

	for {
		select {
		case <-done:
			return
		case <-t.C:
			a.app.Draw()
		}
	}
}

func (a *Application) Stop() {
	if a.server != nil {
		_ = a.server.Close()
	}
	a.app.Stop()
}

func (a *Application) disconnect() {
	if err := a.server.Close(); err != nil {
		a.l.Error("error closing connection", "error", err)
	}
	a.server = nil
	a.l.Info("disconnected from server")
}

func (a *Application) connect(data *ui.LoginData) error {
	var (
		ctx  = context.Background()
		err  error
		conn net.Conn
	)

	a.l.Info("connecting to server", "host", data.Host, "port", data.Port, "key", data.ReconnectKey, "tls", data.TLS)
	if data.TLS {
		ws.DefaultDialer.TLSConfig = &tls.Config{
			InsecureSkipVerify: true,
		}

		conn, _, _, err = ws.DefaultDialer.Dial(ctx, fmt.Sprintf("wss://%s:%d", data.Host, data.Port))
		if err != nil {
			return err
		}
	} else {
		conn, _, _, err = ws.DefaultDialer.Dial(ctx, fmt.Sprintf("ws://%s:%d", data.Host, data.Port))
		if err != nil {
			return err
		}
	}

	a.server = network.NewClientConnection(conn, a.l)
	go a.receiver(data, a.server.Inbox())
	a.l.Info("connection established")

	a.l.Info("sending key exchange", "key", data.ReconnectKey, "username", data.Username)
	if err = a.server.Send(message.NewKeyExchangeWithName(data.ReconnectKey, data.Username)); err != nil {
		return err
	}

	// Wait and check if connection is closed for lobby being full.
	select {
	case m := <-a.errCh:
		return errors.New(m)
	case <-time.After(300 * time.Millisecond):
	}
	if !a.server.IsAlive() {
		return errConnectionClosed
	}

	// Update server status widget.
	a.srvStatus.SetHost(data.Host)
	a.srvStatus.SetConnection(true)
	if data.ReconnectKey != "" {
		a.srvStatus.SetReconnectKey(data.ReconnectKey)
	}

	go a.connectionWatcher()

	return nil
}

func (a *Application) receiver(data *ui.LoginData, incoming <-chan message.Message) {
	for msg := range incoming {
		a.l.Debug("received message", "type", msg.Type(), "payload", msg.Payload())
		switch msg.Type() {
		case message.Error:
			a.errCh <- msg.Payload().(string)
		case message.KeyExchange:
			a.handleKeyExchange(data, msg.Payload().([]string))
		case message.StateUpdate:
			a.handleStateUpdate(msg.Payload().(*message.GameState))
		case message.StartTurn:
			a.handleStartTurn(msg.Payload().(game.TurnPhase))
		case message.JournalAction:
			a.handleJournalAction(msg.Payload().(*message.Action))
		case message.JournalTrade:
			a.handleJournalTrade(msg.Payload().(*message.Trade))
		}
	}
}

func (a *Application) connectionWatcher() {
	for {
		if a.server == nil || !a.server.IsAlive() {
			a.srvStatus.SetConnection(false)

			modal := ui.NewErrorModal(errConnectionClosed)
			a.pages.AddPage(errorPageName, modal.GetModal(), true, true)
			modal.SetHandler(func(int, string) {
				a.pages.RemovePage(errorPageName)
				a.pages.SwitchToPage(titlePageName)
			})
			a.server = nil
			return
		}
		time.Sleep(time.Second)
	}
}
