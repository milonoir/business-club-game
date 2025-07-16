package main

import (
	"flag"
	"log/slog"
	"os"

	"github.com/milonoir/business-club-game/internal/client"
)

var (
	logFile   = flag.String("log", "", "log file path")
	logLevel  = flag.String("log-level", "debug", "log level (debug, info, warn, error, fatal)")
	logFormat = flag.String("log-format", "text", "log format (json, text)")

	noOpFunc = func() error {
		return nil
	}
)

func getLogHandler() (slog.Handler, func() error) {
	if *logFile == "" {
		return slog.DiscardHandler, noOpFunc
	}

	f, err := os.OpenFile(*logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}

	var level slog.Level
	switch *logLevel {
	case "debug":
		level = slog.LevelDebug
	case "info":
		level = slog.LevelInfo
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	case "fatal":
		level = slog.LevelError
	default:
		panic("unsupported log level: " + *logLevel)
	}
	opts := &slog.HandlerOptions{Level: level}

	var handler slog.Handler
	switch *logFormat {
	case "json":
		handler = slog.NewJSONHandler(f, opts)
	case "text":
		handler = slog.NewTextHandler(f, opts)
	default:
		panic("unsupported log format: " + *logFormat)
	}

	return handler, f.Close
}

func main() {
	flag.Parse()

	handler, closeFunc := getLogHandler()
	defer func() {
		if err := closeFunc(); err != nil {
			slog.Error("failed to close log file", "error", err)
		}
	}()

	logger := slog.New(handler).With("app", "bc-client")

	if err := client.NewApplication(logger).Run(); err != nil {
		logger.Error("run application", "error", err)
	}
}
