package main

import (
	"os"
	"time"

	"github.com/milonoir/business-club-game/internal/client"
	"github.com/rivo/tview"
	"golang.org/x/exp/slog"
)

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
	f, err := os.OpenFile("bc-client.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	logger := slog.New(slog.NewJSONHandler(f, &slog.HandlerOptions{Level: slog.LevelDebug})).With("app", "bc-client")

	app := client.NewApplication(logger)

	go refresh(app.GetApplication())

	if err = app.Run(); err != nil {
		logger.Error("run application", "error", err)
	}
}
