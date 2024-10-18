package main

import (
	"flag"
	"log"
	"log/slog"
	"os"

	"github.com/brocode/neoweb/config"
	"github.com/brocode/neoweb/server"
)

func main() {
	configFile := flag.String("config-file", "config.hcl", "Config file location")

	config, err := config.ParseConfig(*configFile)
	if err != nil {
		log.Fatalf("Failed to parse config %v", err)
	}

	flag.Parse()

	var level slog.Level
	switch config.Log.Level {
	case "debug":
		level = slog.LevelDebug
	case "info":
		level = slog.LevelInfo
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	default:
		log.Fatalf("unknown log level: %s", config.Log.Level)
	}

	var handler slog.Handler
	switch config.Log.Format {
	case "text":
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: level})
	case "json":
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: level})
	default:
		log.Fatalf("unknown log format: %s", config.Log.Format)
	}

	logger := slog.New(handler)
	slog.SetDefault(logger)

	server := server.NewServer(config)
	defer server.Close()

	server.Start()
}
