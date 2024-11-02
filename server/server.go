package server

import (
	"embed"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/brocode/neoweb/components"
	"github.com/brocode/neoweb/config"
	"github.com/brocode/neoweb/key"
	"github.com/brocode/neoweb/nvimwrapper"
	"github.com/brocode/neoweb/server/middleware"
)

//go:embed static
var staticFs embed.FS

type Server struct {
	nw     *nvimwrapper.NvimWrapper
	config *config.Config
}

func NewServer(config *config.Config) *Server {
	return &Server{
		nw:     nil,
		config: config,
	}
}

func (s *Server) Close() {
	if s.nw != nil {
		s.nw.Close()
	}
}

func (s *Server) getRoot(w http.ResponseWriter, r *http.Request) {
	err := components.AdminPage(s.config, s.nw != nil).Render(r.Context(), w)
	if err != nil {
		slog.Error("Failed to render root", "err", err)
	}
}

func (s *Server) getEditor(w http.ResponseWriter, r *http.Request) {
	if s.nw == nil {
		err := components.ErrorPage("no nvim connection").Render(r.Context(), w)
		if err != nil {
			slog.Error("Failed to render error page", "err", err)
		}
		return
	}
	result, err := s.nw.Render()
	if err != nil {
		slog.Error("Nvim failed", "err", err)
		http.Error(w, "Nvim failed", 500)

		renderErr := components.ErrorPage("nvim render failed").Render(r.Context(), w)
		if renderErr != nil {
			slog.Error("Failed to render error page", "err", renderErr)
		}
		return
	}
	err = components.EditorPage(result).Render(r.Context(), w)
	if err != nil {
		slog.Error("Failed to render response", "error", err)
	}
}

func (s *Server) killNvim(_ http.ResponseWriter, _ *http.Request) {
	if s.nw != nil {
		s.nw.Close()
	}
	s.nw = nil
}

func (s *Server) launchNvim(_ http.ResponseWriter, _ *http.Request) {
	if s.nw != nil {
		slog.Warn("Already connected to nvim, closing first")
		s.nw.Close()
	}
	nvimWrapper, err := nvimwrapper.Spawn(&s.config.Nvim)
	if err != nil {
		slog.Error("Failed to spawn neovim", "Error", err)

		// TODO this sleep is so that async stuff can print errors.
		// Need to remove the entire nvim spawn handling anyway
		time.Sleep(1 * time.Second)
		os.Exit(1)
	}

	err = nvimWrapper.OpenFile("demo.md")
	if err != nil {
		slog.Error("Failed to open file", "Error", err)
		os.Exit(1)
	}
	s.nw = nvimWrapper
}

func (s *Server) postKeypress(w http.ResponseWriter, r *http.Request) {
	if s.nw == nil {
		slog.Error("Failed to keypress", "reason", "nvim not connected")
	}
	var keyPress key.KeyPress
	err := json.NewDecoder(r.Body).Decode(&keyPress)
	if err != nil {
		http.Error(w, "Failed to unmarshall request", 400)
		return
	}

	s.nw.SendKey(keyPress)
}

func (s *Server) postPaste(w http.ResponseWriter, r *http.Request) {
	if s.nw == nil {
		slog.Error("Failed to paste", "reason", "nvim not connected")
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request", 400)
		return
	}
	text := string(body)

	err = s.nw.Paste(text)
	if err != nil {
		http.Error(w, "Failed to paste text", 500)
		return
	}

}
func (s *Server) getEvents(w http.ResponseWriter, r *http.Request) {
	if s.nw == nil {
		slog.Error("Failed to get events", "reason", "nvim not connected")
	}
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	// Flush the headers immediately
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported", http.StatusInternalServerError)
		return
	}

	err := s.nw.RenderOnFlush(r.Context(), func(result nvimwrapper.NvimResult) error {
		fmt.Fprintf(w, "event: render\n")
		fmt.Fprintf(w, "data:")
		err := components.Editor(result).Render(r.Context(), w)
		if err != nil {
			return fmt.Errorf("Failed to render response: %w", err)
		}
		fmt.Fprintf(w, "\n\n")
		flusher.Flush()
		return nil
	})
	if err != nil {
		slog.Error("Failed to render on flush", "err", err)
	}

	slog.Info("Events client disconnected")
}

func (s *Server) Start() {

	mux := http.NewServeMux()

	mux.Handle("GET /static/", middleware.CacheWhileServerIsRunning(middleware.GzipMiddleware(http.FileServer(http.FS(staticFs)))))

	mux.Handle("GET /", middleware.GzipMiddleware(http.HandlerFunc(s.getRoot)))
	mux.Handle("GET /editor", middleware.GzipMiddleware(http.HandlerFunc(s.getEditor)))

	mux.HandleFunc("POST /editor/launch", s.launchNvim)
	mux.HandleFunc("POST /editor/close", s.killNvim)

	mux.HandleFunc("POST /editor/keypress", s.postKeypress)

	mux.HandleFunc("POST /editor/paste", s.postPaste)

	mux.HandleFunc("GET /editor/events", s.getEvents)

	addr := s.config.Server.ListenAddr
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
