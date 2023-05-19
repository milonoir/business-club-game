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

	"github.com/gdamore/tcell/v2"
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"github.com/milonoir/business-club-game/internal/client"
	"github.com/milonoir/business-club-game/internal/message"
	"github.com/rivo/tview"
)

// THIS IS ONLY A TEST CLIENT.

//go:embed splash.ascii
var splashScreen string

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
	// Create application.
	app := tview.NewApplication()

	// Create pages.
	pages := tview.NewPages()

	// Create grid.
	mainScreen := tview.NewGrid().
		SetRows(3, 0, 3).
		SetColumns(30, 0, 30).
		SetBorders(true)

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

func main() {
	app := buildApp()

	if err := app.Run(); err != nil {
		panic(err)
	}
}
