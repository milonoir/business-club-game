package client

import (
	"github.com/rivo/tview"
)

const (
	titlePageName = "titlePage"
	gamePageName  = "gamePage"
)

type Application struct {
	// Main application.
	app   *tview.Application
	pages *tview.Pages

	// Game screen.
	standings *StandingsPanel
	turn      *TurnPanel
	graph     *GraphPanel
	history   *HistoryPanel
	action    *ActionList
	version   *VersionPanel
	server    *ServerStatusPanel

	// Title screen.
	title *TitlePanel
	login *LoginForm

	cp *CompanyProvider
}

func NewApplication() *Application {
	a := &Application{
		cp: NewCompanyProvider(),
	}

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
	a.title = NewTitlePanel()
	titlePage.AddItem(a.title.GetTextView(), 0, 3, false)

	// Login form.
	a.login = NewLoginForm(
		func(data *LoginData) {
			// TODO: init server connection
			a.pages.SwitchToPage(gamePageName)
		},
		func() {
			a.app.Stop()
		},
	)
	titlePage.AddItem(a.login.GetForm(), 0, 1, true)

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
	a.standings = NewStandingsPanel(a.cp)
	topRow.AddItem(a.standings.GetGrid(), 0, 1, 1, 1, 1, 1, false)

	// Turn widget.
	a.turn = NewTurnPanel()
	topRow.AddItem(a.turn.GetTextView(), 0, 3, 1, 1, 1, 1, false)

	// Middle row of the game page.
	middleRow := tview.NewGrid()
	middleRow.
		SetColumns(96, 0).
		SetRows(22)
	gamePage.AddItem(middleRow, 1, 0, 1, 3, 1, 1, false)

	// Stock price graph panel.
	a.graph = NewGraphPanel(a.cp)
	middleRow.AddItem(a.graph.GetGrid(), 0, 0, 1, 1, 1, 1, false)

	// History panel.
	a.history = NewHistoryPanel(a.cp)
	middleRow.AddItem(a.history.GetTextView(), 0, 1, 1, 1, 1, 1, false)

	// Bottom row of the game page.
	bottomRow := tview.NewGrid()
	bottomRow.
		SetColumns(0, 0, 0).
		SetRows(9)
	gamePage.AddItem(bottomRow, 2, 0, 1, 3, 1, 1, true)

	// Action list.
	// TODO: how to change focus?
	a.action = NewActionList(a.cp)
	bottomRow.AddItem(a.action.GetList(), 0, 1, 1, 1, 1, 1, true)

	// Game version widget.
	a.version = NewVersionPanel()
	gamePage.AddItem(a.version.GetTextView(), 4, 0, 1, 1, 1, 1, false)

	// Server status widget.
	a.server = NewServerStatus()
	gamePage.AddItem(a.server.GetTextView(), 4, 1, 1, 2, 1, 1, false)

	// Setup pages and set pages as root.
	a.pages.
		AddPage(titlePageName, titlePage, true, true).
		AddPage(gamePageName, gamePage, true, false)

	a.app.SetRoot(a.pages, true)
}
