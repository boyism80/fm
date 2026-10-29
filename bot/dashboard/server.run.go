package dashboard

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	botactor "github.com/boyism80/fm/bot/actor"
)

const (
	RunRunning  = "running"
	RunStopping = "stopping"
	RunDone     = "done"

	SuitePending = "pending"
	SuiteRunning = "running"
	SuitePassed  = "passed"
	SuiteFailed  = "failed"
	SuiteSkipped = "skipped"
	SuiteStopped = "stopped"
)

type Run struct {
	ID         string      `json:"id"`
	State      string      `json:"state"`
	Seats      int         `json:"seats"`
	StartedAt  time.Time   `json:"started_at"`
	FinishedAt time.Time   `json:"finished_at,omitzero"`
	Code       int         `json:"code"`
	Error      string      `json:"error,omitempty"`
	Suites     []*RunSuite `json:"suites"`
}

type RunSuite struct {
	Name      string    `json:"name"`
	Title     string    `json:"title"`
	Serial    bool      `json:"serial"`
	State     string    `json:"state"`
	Seat      int       `json:"seat,omitempty"`
	StartedAt time.Time `json:"started_at,omitzero"`
	Seconds   float64   `json:"seconds,omitempty"`
	Infra     bool      `json:"infra,omitempty"`
	Failures  []string  `json:"failures,omitempty"`
	Log       []string  `json:"log,omitempty"`
}

type startRequest struct {
	Suites []string `json:"suites"`
	Seats  int      `json:"seats"`
}

func (s *Server) start(w http.ResponseWriter, r *http.Request) {
	var req startRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		reply(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	if len(req.Suites) == 0 {
		reply(w, http.StatusBadRequest, map[string]string{"error": "실행할 스위트를 고르세요"})
		return
	}
	if req.Seats == 0 {
		req.Seats = s.cfg.Seats
	}
	if req.Seats < 1 || req.Seats > s.cfg.Seats {
		reply(w, http.StatusBadRequest, map[string]string{"error": "seat 수는 1 이상 config seats 이하여야 합니다"})
		return
	}
	suites, err := s.catalog()
	if err != nil {
		reply(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	titles := make(map[string]string, len(suites))
	for _, suite := range suites {
		titles[suite.Name] = suite.Title
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	if s.run != nil && s.run.State != RunDone {
		reply(w, http.StatusConflict, map[string]string{"error": "이미 실행 중인 테스트가 있습니다"})
		return
	}

	cfg := *s.cfg
	cfg.Seats = req.Seats
	done := make(chan int, 1)
	s.run = &Run{State: RunRunning, Seats: req.Seats, StartedAt: time.Now()}
	s.runner = s.root.Spawn(actor.PropsFromProducer(func() actor.Actor {
		return botactor.NewRunnerActor(&cfg, s.wz, "", req.Suites, "", func(msg any) { s.watch(msg, titles) }, done)
	}))
	go s.finish(s.run, s.runner, done)
	s.broadcast()
	reply(w, http.StatusAccepted, map[string]string{"state": RunRunning})
}

func (s *Server) stop(w http.ResponseWriter, _ *http.Request) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.run == nil || s.run.State != RunRunning {
		reply(w, http.StatusConflict, map[string]string{"error": "실행 중인 테스트가 없습니다"})
		return
	}
	s.run.State = RunStopping
	s.root.Send(s.runner, &botactor.RunStop{})
	s.broadcast()
	reply(w, http.StatusAccepted, map[string]string{"state": RunStopping})
}

func (s *Server) watch(msg any, titles map[string]string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	switch msg := msg.(type) {
	case *botactor.RunStarted:
		s.run.ID = msg.ID
		for _, suite := range slices.Concat(msg.Parallel, msg.Serial) {
			s.run.Suites = append(s.run.Suites, &RunSuite{Name: suite.Name, Title: titles[suite.Name], Serial: suite.Serial, State: SuitePending})
		}
	case *botactor.SuiteStarted:
		suite := s.find(msg.Suite)
		if suite == nil {
			return
		}
		suite.State, suite.Seat, suite.StartedAt = SuiteRunning, msg.Seat, time.Now()
	case *botactor.SuiteFinished:
		suite := s.find(msg.Suite)
		if suite == nil {
			return
		}
		switch {
		case msg.Stopped:
			suite.State = SuiteStopped
		case !msg.Passed:
			suite.State = SuiteFailed
		case msg.Skipped:
			suite.State = SuiteSkipped
		default:
			suite.State = SuitePassed
		}
		suite.Title, suite.Infra, suite.Failures, suite.Log = msg.Name, msg.Infra, msg.Failures, msg.Log
		suite.Seconds = msg.Elapsed.Seconds()
	}
	s.broadcast()
}

// find must be called with mu held.
func (s *Server) find(name string) *RunSuite {
	for _, suite := range s.run.Suites {
		if suite.Name == name {
			return suite
		}
	}
	return nil
}

func (s *Server) finish(run *Run, runner *actor.PID, done <-chan int) {
	code := <-done
	s.root.Stop(runner)

	s.mu.Lock()
	defer s.mu.Unlock()
	run.State, run.Code, run.FinishedAt = RunDone, code, time.Now()
	for _, suite := range run.Suites {
		if suite.State == SuitePending {
			suite.State = SuiteStopped
		}
	}
	if run.ID == "" {
		run.Error = "registry.lua를 읽지 못했거나 고른 스위트가 registry에 없습니다"
		s.broadcast()
		return
	}
	if err := s.save(run); err != nil {
		log.Printf("dashboard: save run %s: %v", run.ID, err)
	}
	s.broadcast()
}

func (s *Server) save(run *Run) error {
	dir := filepath.Join(s.dir, run.ID)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(run, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(dir, "summary.json"), data, 0o644)
}

type runRow struct {
	ID         string    `json:"id"`
	StartedAt  time.Time `json:"started_at"`
	FinishedAt time.Time `json:"finished_at"`
	Seats      int       `json:"seats"`
	Code       int       `json:"code"`
	Total      int       `json:"total"`
	Passed     int       `json:"passed"`
	Failed     int       `json:"failed"`
	Skipped    int       `json:"skipped"`
	Stopped    int       `json:"stopped"`
}

func (s *Server) runs(w http.ResponseWriter, _ *http.Request) {
	entries, err := os.ReadDir(s.dir)
	if err != nil && !os.IsNotExist(err) {
		reply(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	rows := []runRow{}
	for _, entry := range entries {
		data, err := os.ReadFile(filepath.Join(s.dir, entry.Name(), "summary.json"))
		if err != nil {
			continue
		}
		var run Run
		if err := json.Unmarshal(data, &run); err != nil {
			continue
		}
		row := runRow{ID: run.ID, StartedAt: run.StartedAt, FinishedAt: run.FinishedAt, Seats: run.Seats, Code: run.Code, Total: len(run.Suites)}
		for _, suite := range run.Suites {
			switch suite.State {
			case SuitePassed:
				row.Passed++
			case SuiteFailed:
				row.Failed++
			case SuiteSkipped:
				row.Skipped++
			case SuiteStopped:
				row.Stopped++
			}
		}
		rows = append(rows, row)
	}
	slices.SortFunc(rows, func(a, b runRow) int { return b.StartedAt.Compare(a.StartedAt) })
	reply(w, http.StatusOK, rows)
}

func (s *Server) history(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" || strings.Trim(id, "0123456789abcdefghijklmnopqrstuvwxyz") != "" {
		reply(w, http.StatusBadRequest, map[string]string{"error": "잘못된 실행 ID"})
		return
	}
	data, err := os.ReadFile(filepath.Join(s.dir, id, "summary.json"))
	if err != nil {
		reply(w, http.StatusNotFound, map[string]string{"error": "실행 기록이 없습니다"})
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_, _ = w.Write(data)
}
