package server

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/brocode/neoweb/components"
	"github.com/brocode/neoweb/nvimwrapper"
)

func Run() {

	nvimWrapper, err := nvimwrapper.Spawn()
	if err != nil {
		slog.Error("Failed to spawn neovim", "Error", err)
		os.Exit(1)
	}
	defer nvimWrapper.Close()

	err = nvimWrapper.OpenFile("main.go")
	if err != nil {
		slog.Error("Failed to open file", "Error", err)
		os.Exit(1)
	}

	// TODO this has to be actual input from the browser later
	err = nvimWrapper.Input("ggjj$")
	if err != nil {
		slog.Error("Failed to send initial input to neovim", "Error", err)
		os.Exit(1)
	}

	mux := http.NewServeMux()

	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		result, err := nvimWrapper.Render()
		if err != nil {
			slog.Error("Nvim failed", "err", err)
			http.Error(w, "Nvim failed", 500)
			return
		}
		err = components.Main(result).Render(r.Context(), w)
        if err != nil {
            slog.Error("Failed to render response", "error",  err)
        }
	})

	addr := ":8080"
	slog.Info("Start server", "addr", addr)

	server := &http.Server{
		Addr:     addr,
		Handler:  mux,
		ErrorLog: slog.NewLogLogger(slog.Default().Handler(), slog.LevelError),
	}

	err = server.ListenAndServe()
	if err != nil {
		slog.Error("Stop server", "Error", err)
		os.Exit(1)
	}

}
