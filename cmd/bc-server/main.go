package main

import (
	"encoding/json"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/milonoir/business-club-game/internal/game"
	"github.com/milonoir/business-club-game/internal/server"
)

const (
	serverPort = 8585
)

func main() {
	// Setup logger.
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})).With("app", "bc-server")

	// Check command line arguments.
	if len(os.Args) < 2 {
		log.Fatal("missing assets file argument")
	}

	// Loading game assets.
	b, err := os.ReadFile(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	var a game.Assets
	if err = json.Unmarshal(b, &a); err != nil {
		log.Fatalf("corrupted assets file: %s\n", err)
	}

	// Initializing the game.
	//g := server.NewGame(a)
	//
	//fmt.Println(g.Companies)
	//fmt.Println(g.PlayerDeck)
	//fmt.Println(g.BankDeck)

	// Setup OS signal trap.
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	// Run server.
	srv := server.NewServer(serverPort, logger)
	srv.Start(&a)

	// Catch signal.
	<-sig

	// Shutdown server gracefully.
	if err = srv.Stop(); err != nil {
		log.Fatal(err)
	}
}
