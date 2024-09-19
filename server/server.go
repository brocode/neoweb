package server

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/brocode/neoweb/components"
	"github.com/brocode/neoweb/nvim"
)

func Run() {

	mux := http.NewServeMux()

	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		result, err := nvim.RunNvim()
		if err != nil {
			slog.Error("Nvim failed", "err", err)
			http.Error(w, "Nvim failed", 500)
			return
		}
		components.Hello(result).Render(r.Context(), w)
	})

	addr := ":8080"
	slog.Info("Start server", "addr", addr)

	server := &http.Server{
		Addr:     addr,
		Handler:  mux,
		ErrorLog: slog.NewLogLogger(slog.Default().Handler(), slog.LevelError),
	}

	err := server.ListenAndServe()
	if err != nil {
		slog.Error("Stop server", "Error", err)
		os.Exit(1)
	}

}
