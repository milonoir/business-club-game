package main

import (
	"log/slog"
	"os"

	"github.com/milonoir/business-club-game/internal/client"
)

func main() {
	f, err := os.OpenFile("bc-client.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	defer func() { _ = f.Close() }()

	logger := slog.New(slog.NewJSONHandler(f, &slog.HandlerOptions{Level: slog.LevelDebug})).With("app", "bc-client")

	if err = client.NewApplication(logger).Run(); err != nil {
		logger.Error("run application", "error", err)
	}
}
