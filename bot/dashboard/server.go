package dashboard

import (
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"net/http"
	"sync"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/common/config"
	"github.com/boyism80/fm/maple/wz"
)

//go:embed web
var web embed.FS

// Server runs one integration run at a time and streams its state to the dashboard page.
type Server struct {
	cfg       *config.Bot
	wz        *wz.Resources
	root      *actor.RootContext
	dir       string
	mu        sync.Mutex
	run       *Run
	runner    *actor.PID
	listeners map[chan []byte]bool
}

func New(cfg *config.Bot, resources *wz.Resources, root *actor.RootContext, dir string) *Server {
	return &Server{
		cfg:       cfg,
		wz:        resources,
		root:      root,
		dir:       dir,
		listeners: make(map[chan []byte]bool),
	}
}

func (s *Server) ListenAndServe(addr string) error {
	files, err := fs.Sub(web, "web")
	if err != nil {
		return err
	}
	static := http.FileServerFS(files)
	mux := http.NewServeMux()
	// Embedded files carry no modification time, so the page would keep a stale copy after a rebuild.
	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-cache")
		static.ServeHTTP(w, r)
	})
	mux.HandleFunc("GET /api/suites", s.suites)
	mux.HandleFunc("GET /api/events", s.events)
	mux.HandleFunc("GET /api/runs", s.runs)
	mux.HandleFunc("POST /api/runs", s.start)
	mux.HandleFunc("POST /api/runs/current/stop", s.stop)
	mux.HandleFunc("GET /api/runs/{id}", s.history)
	return http.ListenAndServe(addr, mux)
}

func (s *Server) suites(w http.ResponseWriter, _ *http.Request) {
	suites, err := s.catalog()
	if err != nil {
		reply(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	reply(w, http.StatusOK, map[string]any{"seats": s.cfg.Seats, "loaded_at": time.Now(), "suites": suites})
}

func (s *Server) events(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "streaming unsupported", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")

	updates := make(chan []byte, 1)
	s.mu.Lock()
	s.listeners[updates] = true
	first := s.state()
	s.mu.Unlock()
	defer func() {
		s.mu.Lock()
		delete(s.listeners, updates)
		s.mu.Unlock()
	}()

	fmt.Fprintf(w, "data: %s\n\n", first)
	flusher.Flush()
	for {
		select {
		case <-r.Context().Done():
			return
		case data := <-updates:
			fmt.Fprintf(w, "data: %s\n\n", data)
			flusher.Flush()
		}
	}
}

// state must be called with mu held.
func (s *Server) state() []byte {
	data, _ := json.Marshal(map[string]any{"now": time.Now(), "seats": s.cfg.Seats, "run": s.run})
	return data
}

// broadcast must be called with mu held. A slow page only misses the states between its reads.
func (s *Server) broadcast() {
	data := s.state()
	for updates := range s.listeners {
		select {
		case <-updates:
		default:
		}
		updates <- data
	}
}

func reply(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}
