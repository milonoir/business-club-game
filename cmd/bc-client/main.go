package main

import (
	"context"
	"crypto/tls"
	_ "embed"
	"encoding/json"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"github.com/milonoir/business-club-game/internal/client"
	"github.com/milonoir/business-club-game/internal/game"
	"github.com/milonoir/business-club-game/internal/message"
	"github.com/rivo/tview"
)

// THIS IS ONLY A TEST CLIENT.

//go:embed splash.ascii
var splashScreen string

//go:embed sample_cards.json
var cardsJson string

func main_old() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ws.DefaultDialer.TLSConfig = &tls.Config{
		InsecureSkipVerify: true,
	}

	// TLS
	//conn, _, _, err := ws.DefaultDialer.Dial(ctx, "wss://localhost:8585")
	//if err != nil {
	//	log.Fatal(err)
	//}

	// Non-TLS
	conn, _, _, err := ws.DefaultDialer.Dial(ctx, "ws://localhost:8585")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	var wg sync.WaitGroup
	done := make(chan struct{})

	wg.Add(1)
	go func() {
		defer wg.Done()
		responder(conn, done)
	}()

	// Setup OS signal trap.
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	// Catch signal.
	<-sig

	close(done)
	wg.Wait()

	//if err = wsutil.WriteClientMessage(conn, ws.OpText, []byte("hello")); err != nil {
	//	log.Fatal(err)
	//}
}

func responder(conn net.Conn, done chan struct{}) {
	for {
		select {
		case <-done:
			return
		default:
		}

		m := make([]wsutil.Message, 0, 5)
		m, err := wsutil.ReadServerMessage(conn, m)
		if err != nil {
			log.Printf("read error: %+v", err)
		}

		for i, msg := range m {
			log.Printf("#%d, opcode: %v, payload: %s", i+1, msg.OpCode, msg.Payload)

			if msg.OpCode == ws.OpPing {
				if err = wsutil.WriteClientMessage(conn, ws.OpPong, msg.Payload); err != nil {
					log.Printf("write error: %+v", err)
				}
			} else {
				parseAndRespond(conn, msg.Payload)
			}
		}
	}
}

func parseAndRespond(conn net.Conn, raw []byte) {
	msg, err := message.Parse(raw)
	if err != nil {
		log.Printf("parse error: %+v", err)
		return
	}

	switch msg.Type() {
	case message.Auth:
		key := msg.Payload().(string)
		if key == "" {
			sendEmptyAuth(conn)
			return
		}
		log.Printf("received auth key: %s", key)
	}
}

func sendEmptyAuth(conn net.Conn) {
	bb, _ := json.Marshal(message.EmptyAuth)
	if err := wsutil.WriteClientMessage(conn, ws.OpText, bb); err != nil {
		log.Printf("write error: %+v", err)
	}
}

func buildApp() *tview.Application {
	cp := client.NewCompanyProvider()
	cp.AddCompany("Amfora", "blue")
	cp.AddCompany("Domus", "orange")
	cp.AddCompany("Piért", "yellow")
	cp.AddCompany("Skála-Coop", "red")

	pp := client.NewPlayerProvider([]string{"Xenial Xerus", "Bionic Beaver", "Focal Fossa", "Jammy Jellyfish"})

	// Create application.
	app := tview.NewApplication()

	// Create pages.
	pages := tview.NewPages()

	// Create grid.
	mainScreen := tview.NewGrid().
		SetColumns(0, 0, 0).
		SetRows(0, 22, 0, 1, 1)

	// Top row of the main screen.
	topRow := tview.NewGrid()
	topRow.
		SetColumns(0, 96, 0, 25, 0).
		SetRows(9)
	mainScreen.AddItem(topRow, 0, 0, 1, 3, 1, 1, false)

	// Standings panel.
	standings := client.NewStandingsPanel(cp)
	topRow.AddItem(standings.GetGrid(), 0, 1, 1, 1, 1, 1, false)

	// TEST ONLY.
	standings.SetPlayerNames("Jammy Jellyfish", []string{"Xenial Xerus", "Bionic Beaver", "Focal Fossa"})
	standings.RefreshCompanyNames()
	standings.SetPrices(30, 190, 330, 60)
	standings.PlayerUpdate(5, 8, 2, 10, 5000)
	standings.OpponentUpdate(0, 0, 1, 2, 0, 2, false, false)
	standings.OpponentUpdate(1, 5, 2, 0, 1, 5, false, false)
	standings.OpponentUpdate(2, 3, 1, 0, 1, 1, false, false)

	// Turn widget.
	turns := client.NewTurnPanel()
	topRow.AddItem(turns.GetTextView(), 0, 3, 1, 1, 1, 1, false)

	// TEST ONLY.
	turns.SetMaxTurns(15)
	turns.NewTurn(pp.Players())
	turns.NextPlayer()

	// Bottom row of the main screen.
	bottomRow := tview.NewGrid()
	bottomRow.
		SetColumns(0, 0, 0).
		SetRows(9)
	mainScreen.AddItem(bottomRow, 2, 0, 1, 3, 1, 1, true)

	// Action list.
	var cards []*game.Card
	if err := json.Unmarshal([]byte(cardsJson), &cards); err != nil {
		panic(err)
	}
	actions := client.NewActionList(cp)
	bottomRow.AddItem(actions.GetList(), 0, 1, 1, 1, 1, 1, true)
	actions.Update(cards)

	// Game version widget.
	ver := client.NewVersionPanel()
	mainScreen.AddItem(ver.GetTextView(), 4, 0, 1, 1, 1, 1, false)

	// TEST ONLY.
	ver.SetVersion("0.1")

	// Server status widget.
	status := client.NewServerStatus()
	mainScreen.AddItem(status.GetTextView(), 4, 1, 1, 2, 1, 1, false)

	// TEST ONLY.
	status.SetHost("localhost:8585")
	status.SetAuthKey("ab3tesjk4")
	status.SetConnection(true)

	// Middle row of the main screen.
	middleRow := tview.NewGrid()
	middleRow.
		SetColumns(96, 0).
		SetRows(22)
	mainScreen.AddItem(middleRow, 1, 0, 1, 3, 1, 1, false)

	// Stock price graph panel.
	graphs := client.NewGraphPanel(cp)
	middleRow.AddItem(graphs.GetGrid(), 0, 0, 1, 1, 1, 1, false)
	go func() {
		graphs.Add(10, 0, 60, 290)
		time.Sleep(time.Second)
		graphs.Add(20, 220, 80, 0)
		time.Sleep(time.Second)
		graphs.Add(30, 120, 230, 0)
		time.Sleep(time.Second)
		graphs.Add(40, 60, 390, 40)
		time.Sleep(time.Second)
		graphs.Add(50, 280, 190, 0)
		time.Sleep(time.Second)
		graphs.Add(60, 0, 10, 10)
		time.Sleep(time.Second)
		graphs.Add(70, 10, 0, 190)
		time.Sleep(time.Second)
		graphs.Add(80, 20, 0, 240)
		time.Sleep(time.Second)
		graphs.Add(90, 70, 90, 140)
		time.Sleep(time.Second)
		graphs.Add(100, 210, 170, 340)
		time.Sleep(time.Second)
		graphs.Add(300, 210, 190, 370)
	}()

	// History panel.
	history := client.NewHistoryPanel(cp)
	middleRow.AddItem(history.GetTextView(), 0, 1, 1, 1, 1, 1, false)

	// TEST ONLY>
	history.AddItem(&client.ActionItem{
		ActorType: client.ActorBank,
		Mod:       &(cards[0].Mods[0]),
		NewPrice:  120,
	})
	history.AddItem(&client.ActionItem{
		ActorType: client.ActorBank,
		Mod:       &(cards[0].Mods[1]),
		NewPrice:  40,
	})
	history.AddItem(&client.ActionItem{
		ActorType: client.ActorPlayer,
		Name:      pp.Players()[0],
		Mod:       &(cards[1].Mods[1]),
		NewPrice:  370,
	})
	history.AddItem(&client.DealItem{
		Name:         pp.Players()[0],
		Type:         client.DealBuy,
		CompanyIndex: 3,
		Amount:       20,
		Price:        10,
	})

	// Title screen.
	title := tview.NewTextView().
		SetTextAlign(tview.AlignLeft).
		SetTextColor(tcell.ColorYellow).
		SetText(splashScreen)
	title.
		SetBorderPadding(6, 1, 10, 1)

	// Login form.
	login := client.NewLoginForm(
		func(data *client.LoginData) {
			pages.SwitchToPage("main")
		},
		func() {
			app.Stop()
		},
	)

	// Welcome screen.
	welcome := tview.NewFlex().
		AddItem(title, 0, 3, false).
		AddItem(login.GetForm(), 0, 1, true)

	// Add main widgets to pages.
	pages.
		AddPage("welcome", welcome, true, true).
		AddPage("main", mainScreen, true, false)

	// Set application root primitive.
	app.SetRoot(pages, true)

	return app
}

func refresh(app *tview.Application) {
	tick := time.NewTicker(500 * time.Microsecond)
	for {
		select {
		case <-tick.C:
			app.Draw()
		}
	}
}

func main() {
	app := buildApp()

	go refresh(app)

	if err := app.Run(); err != nil {
		panic(err)
	}
}
