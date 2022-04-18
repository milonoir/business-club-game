package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"time"
)

func main() {
	// Loading game assets.
	if len(os.Args) < 2 {
		log.Fatal("missing assets file argument")
	}

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

	// Initializing the game.
	g := newGame(a)

	// --- DEBUG ---
	fmt.Println(g.Companies)
	fmt.Println(g.ActionDeck)
	fmt.Println(g.BankDeck)
}
