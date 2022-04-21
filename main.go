package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// Check command line arguments.
	if len(os.Args) < 2 {
		log.Fatal("missing assets file argument")
	}

	// Loading game assets.
	b, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	var a assets
	if err = json.Unmarshal(b, &a); err != nil {
		log.Fatalf("corrupted assets file: %s\n", err)
	}

	// Initializing the RNG.
	rand.Seed(time.Now().UnixNano())

	// Capture OS signals.
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	// Run server.
	srv := newServer(8585)
	srv.start(sig)

	// Initializing the game.
	//g := newGame(a)
	//
	//// --- DEBUG ---
	//fmt.Println(g.Companies)
	//fmt.Println(g.ActionDeck)
	//fmt.Println(g.BankDeck)
}
