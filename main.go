package main

import (
	"flag"
	"log"
	"log/slog"
	"os"

	"github.com/brocode/neoweb/server"
)

func main() {
	// Define CLI flags
	handlerType := flag.String("log-format", "text", "Log format: text or json")
	logLevel := flag.String("log-level", "info", "Log level: debug, info, warn, error")

	// Parse the flags
	flag.Parse()

	// Determine the log level
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
	default:
		log.Fatalf("unknown log level: %s", *logLevel)
	}

	// Create the appropriate log handler based on the flag
	var handler slog.Handler
	switch *handlerType {
	case "text":
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: level})
	case "json":
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: level})
	default:
		log.Fatalf("unknown log format: %s", *handlerType)
	}

	// Set up the logger
	logger := slog.New(handler)
	slog.SetDefault(logger)

	// Start the server (replace server.Run() with your actual server starting logic)
	server.Run()
}
