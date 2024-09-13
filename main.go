package main

import (
	"log/slog"
	"os"

	"github.com/brocode/neoweb/server"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	server.Run()
}
