package cluster

import (
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"net/http"
	"sync"
	"time"

	"github.com/boyism80/fm/common/config"
	"github.com/boyism80/fm/maple/wz"
)

//go:embed web
var web embed.FS

// Master serves the dashboard, takes agent connections and hands the suites of one run at a time to them.
// Agents run the bots; the master runs none.
type Master struct {
	cfg       *config.Bot
	wz        *wz.Resources
	dir       string
	token     string
	mu        sync.Mutex
	run       *Run
	agents    []*agentConn
	nextBase  int
	nextTag   int
	listeners map[chan []byte]bool
}

func NewMaster(cfg *config.Bot, resources *wz.Resources, dir, token string) *Master {
	return &Master{
		cfg:       cfg,
		wz:        resources,
		dir:       dir,
		token:     token,
		listeners: make(map[chan []byte]bool),
	}
}

func (m *Master) ListenAndServe(addr string) error {
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
	mux.HandleFunc("GET /api/suites", m.suites)
	mux.HandleFunc("GET /api/events", m.events)
	mux.HandleFunc("GET /api/runs", m.runs)
	mux.HandleFunc("POST /api/runs", m.start)
	mux.HandleFunc("POST /api/runs/current/stop", m.stop)
	mux.HandleFunc("GET /api/runs/{id}", m.history)
	mux.HandleFunc("POST /api/agents/{id}/quit", m.quit)
	mux.HandleFunc("GET /ws/agent", m.accept)

	go func() {
		for range time.Tick(heartbeatInterval) {
			m.mu.Lock()
			m.broadcast()
			m.mu.Unlock()
		}
	}()
	return http.ListenAndServe(addr, mux)
}

func (m *Master) suites(w http.ResponseWriter, _ *http.Request) {
	suites, err := m.catalog()
	if err != nil {
		reply(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	reply(w, http.StatusOK, map[string]any{"loaded_at": time.Now(), "suites": suites})
}

func (m *Master) events(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "streaming unsupported", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")

	updates := make(chan []byte, 1)
	m.mu.Lock()
	m.listeners[updates] = true
	first := m.state()
	m.mu.Unlock()
	defer func() {
		m.mu.Lock()
		delete(m.listeners, updates)
		m.mu.Unlock()
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
func (m *Master) state() []byte {
	agents := make([]agentRow, 0, len(m.agents))
	for _, a := range m.agents {
		agents = append(agents, a.row())
	}
	data, _ := json.Marshal(map[string]any{"now": time.Now(), "run": m.run, "agents": agents})
	return data
}

// broadcast must be called with mu held. A slow page only misses the states between its reads.
func (m *Master) broadcast() {
	data := m.state()
	for updates := range m.listeners {
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
