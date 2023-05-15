package main

import (
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/milonoir/business-club-game/internal"
)

const (
	serverPort = 8585
)

func main() {
	// Check command line arguments.
	if len(os.Args) < 2 {
		log.Fatal("missing assets file argument")
	}

	// Loading game assets.
	b, err := os.ReadFile(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	var a internal.Assets
	if err = json.Unmarshal(b, &a); err != nil {
		log.Fatalf("corrupted assets file: %s\n", err)
	}

	// Setup OS signal trap.
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	// Run server.
	srv := internal.NewServer(serverPort)
	srv.Start()

	// Catch signal.
	<-sig

	// Shutdown server gracefully.
	if err = srv.Stop(); err != nil {
		log.Fatal(err)
	}

	// Initializing the game.
	//g := newGame(a)
	//
	//fmt.Println(g.Companies)
	//fmt.Println(g.ActionDeck)
	//fmt.Println(g.BankDeck)
}
