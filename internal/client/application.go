package client

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"

	"github.com/gobwas/ws"
	"github.com/milonoir/business-club-game/internal/client/ui"
	"github.com/milonoir/business-club-game/internal/network"
	"github.com/rivo/tview"
	"golang.org/x/exp/slog"
)

const (
	titlePageName = "titlePage"
	lobbyPageName = "lobbyPage"
	gamePageName  = "gamePage"
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
	action    *ui.ActionList
	version   *ui.VersionPanel
	srvStatus *ui.ServerStatusPanel

	// Lobby screen.

	// Title screen.
	title *ui.TitlePanel
	login *ui.LoginForm

	cp *ui.CompanyProvider

	server network.Connection
	l      *slog.Logger
}

func NewApplication(l *slog.Logger) *Application {
	a := &Application{
		cp: ui.NewCompanyProvider(),
		l:  l,
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
				return
			}
			// Successful login.
			a.pages.SwitchToPage(gamePageName)
		},
		func() {
			a.Stop()
		},
	)
	titlePage.AddItem(a.login.GetForm(), 0, 1, true)

	// -------------------------------------------------- Lobby page

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
	bottomRow := tview.NewGrid()
	bottomRow.
		SetColumns(0, 0, 0).
		SetRows(9)
	gamePage.AddItem(bottomRow, 2, 0, 1, 3, 1, 1, true)

	// Action list.
	// TODO: how to change focus?
	a.action = ui.NewActionList(a.cp)
	bottomRow.AddItem(a.action.GetList(), 0, 1, 1, 1, 1, 1, true)

	// Game version widget.
	a.version = ui.NewVersionPanel()
	gamePage.AddItem(a.version.GetTextView(), 4, 0, 1, 1, 1, 1, false)

	// Server status widget.
	a.srvStatus = ui.NewServerStatus()
	gamePage.AddItem(a.srvStatus.GetTextView(), 4, 1, 1, 2, 1, 1, false)

	// Setup pages and set pages as root.
	a.pages.
		AddPage(titlePageName, titlePage, true, true).
		AddPage(gamePageName, gamePage, true, false)

	a.app.SetRoot(a.pages, true)
}

func (a *Application) GetApplication() *tview.Application {
	return a.app
}

func (a *Application) Run() error {
	return a.app.Run()
}

func (a *Application) Stop() {
	if a.server != nil {
		_ = a.server.Close()
	}
	a.app.Stop()
}

func (a *Application) connect(data *ui.LoginData) error {
	var (
		ctx  = context.Background()
		err  error
		conn net.Conn
	)

	a.l.Info("connecting to server", "host", data.Host, "port", data.Port, "key", data.AuthKey, "tls", data.TLS)
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
	go a.responder(data, a.server.Inbox())
	a.l.Info("connection established")

	a.l.Info("sending auth key", "key", data.AuthKey)
	auth := network.EmptyAuth
	if data.AuthKey != "" {
		auth = network.NewAuthMessage([]byte(data.AuthKey))
	}
	return a.server.Send(auth)
}

func (a *Application) responder(data *ui.LoginData, incoming <-chan network.Message) {
	for msg := range incoming {
		a.l.Debug("received message", "type", msg.Type(), "payload", msg.Payload())
		switch msg.Type() {
		case network.Auth:
			a.handleAuth(data, msg.Payload().(string))
		}
	}
}

func (a *Application) handleAuth(data *ui.LoginData, key string) {
	if key == "" {
		auth := network.EmptyAuth
		if data.AuthKey != "" {
			auth = network.NewAuthMessage([]byte(data.AuthKey))
		}
		if err := a.server.Send(auth); err != nil {
			a.l.Error("send auth", "error", err)
		}
	} else {
		a.l.Info("received auth key", "key", key)
	}
}
