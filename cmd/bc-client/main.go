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

//go:embed graph.ascii
var graphAscii string

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
		SetRows(0, 25, 0, 1)

	// Turn widget.
	turns := client.NewTurnPanel(10)
	mainScreen.AddItem(turns.GetTextView(), 0, 0, 1, 1, 1, 1, false)

	// TEST ONLY.
	turns.NewTurn(pp.Players())
	turns.NextPlayer()

	// Portfolio widget.
	portfolio := client.NewPortfolioPanel(cp)
	mainScreen.AddItem(portfolio.GetTextView(), 0, 1, 1, 1, 1, 1, false)

	// TEST ONLY.
	portfolio.Update(client.PortfolioUpdate{
		P1: 40, N1: 2,
		P2: 230, N2: 9,
		P3: 0, N3: 5,
		P4: 170, N4: 0,
		Cash: 3000,
	})

	// Opponents widget.
	opponents := client.NewOpponentsPanel(pp.OpponentsByPlayer("Bionic Beaver"), cp)
	mainScreen.AddItem(opponents.GetTable(), 0, 2, 1, 1, 1, 1, false)

	// TEST ONLY.
	opponents.Update(client.OpponentsUpdate{
		O1: client.OpponentData{
			N1: 2, N2: 0, N3: 4, N4: 0, C: 5,
		},
		O2: client.OpponentData{
			N1: 2, N2: 3, N3: 3, N4: 1, C: 2,
		},
		O3: client.OpponentData{
			N1: 2, N2: 2, N3: 0, N4: 5, C: 1,
		},
	})

	// Action list.
	var cards []*game.Card
	if err := json.Unmarshal([]byte(cardsJson), &cards); err != nil {
		panic(err)
	}
	actions := client.NewActionList(cp, cards)
	mainScreen.AddItem(actions.GetList(), 2, 1, 1, 1, 1, 1, true)

	// Game version widget.
	ver := client.NewVersionPanel()
	mainScreen.AddItem(ver.GetTextView(), 3, 0, 1, 1, 1, 1, false)

	// TEST ONLY.
	ver.SetVersion("0.1")

	// Server status widget.
	status := client.NewServerStatus("localhost:8585")
	mainScreen.AddItem(status.GetTextView(), 3, 1, 1, 2, 1, 1, false)

	// TEST ONLY.
	status.SetAuthKey("ab3tesjk4")

	// Stock price graph panel.
	graphs := client.NewGraphPanel(cp)
	mainScreen.AddItem(graphs.GetGrid(), 1, 0, 1, 2, 1, 1, false)
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
