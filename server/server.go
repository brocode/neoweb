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
	"github.com/brocode/neoweb/key"
	"github.com/brocode/neoweb/nvimwrapper"
)

//go:embed static
var staticFs embed.FS

type Server struct {
	nw *nvimwrapper.NvimWrapper
}

func NewServer(clean bool) *Server {
	nvimWrapper, err := nvimwrapper.Spawn(clean)
	if err != nil {
		slog.Error("Failed to spawn neovim", "Error", err)
		os.Exit(1)
	}

	err = nvimWrapper.OpenFile("demo.sh")
	if err != nil {
		slog.Error("Failed to open file", "Error", err)
		os.Exit(1)
	}

	return &Server{
		nw: nvimWrapper,
	}
}

func (s *Server) Close() {
	s.nw.Close()
}

func (s *Server) getRoot(w http.ResponseWriter, r *http.Request) {
	result, err := s.nw.Render()
	if err != nil {
		slog.Error("Nvim failed", "err", err)
		http.Error(w, "Nvim failed", 500)
		return
	}
	err = components.Main(result).Render(r.Context(), w)
	if err != nil {
		slog.Error("Failed to render response", "error", err)
	}
}

func (s *Server) postKeypress(w http.ResponseWriter, r *http.Request) {
	var keyPress key.KeyPress
	err := json.NewDecoder(r.Body).Decode(&keyPress)
	if err != nil {
		http.Error(w, "Failed to unmarshall request", 400)
		return
	}

	s.nw.SendKey(keyPress)
}

func (s *Server) postPaste(w http.ResponseWriter, r *http.Request) {
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

func cacheWhileServerIsRunning(next http.Handler) http.Handler {
	startTime := time.Now().Truncate(time.Second)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ifModifiedSince := r.Header.Get("If-Modified-Since")
		if ifModifiedSince != "" {
			modTime, err := time.Parse(http.TimeFormat, ifModifiedSince)
			if err == nil {
				if !startTime.After(modTime) {
					w.WriteHeader(http.StatusNotModified)
					return
				}
			}
		}
		w.Header().Set("Last-Modified", startTime.UTC().Format(http.TimeFormat))
		next.ServeHTTP(w, r)
	})
}

func (s *Server) Start() {

	mux := http.NewServeMux()

	mux.Handle("GET /static/", cacheWhileServerIsRunning(http.FileServer(http.FS(staticFs))))

	mux.HandleFunc("GET /", s.getRoot)

	mux.HandleFunc("POST /keypress", s.postKeypress)

	mux.HandleFunc("POST /paste", s.postPaste)

	mux.HandleFunc("GET /events", s.getEvents)

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
