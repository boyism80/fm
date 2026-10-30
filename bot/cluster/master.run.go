package cluster

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"time"
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
	SuiteLost    = "lost"
)

type Run struct {
	ID         string      `json:"id"`
	State      string      `json:"state"`
	Retry      int         `json:"retry"`
	StartedAt  time.Time   `json:"started_at"`
	FinishedAt time.Time   `json:"finished_at,omitzero"`
	Suites     []*RunSuite `json:"suites"`
}

type RunSuite struct {
	Name      string    `json:"name"`
	Title     string    `json:"title"`
	Serial    bool      `json:"serial"`
	TimeoutMs int       `json:"timeout_ms"`
	State     string    `json:"state"`
	Agent     string    `json:"agent,omitempty"`
	Seat      int       `json:"seat,omitempty"`
	Attempts  int       `json:"attempts"`
	StartedAt time.Time `json:"started_at,omitzero"`
	Seconds   float64   `json:"seconds,omitempty"`
	Infra     bool      `json:"infra,omitempty"`
	Failures  []string  `json:"failures,omitempty"`
	Log       []string  `json:"log,omitempty"`
}

func (run *Run) find(name string) *RunSuite {
	for _, suite := range run.Suites {
		if suite.Name == name {
			return suite
		}
	}
	return nil
}

func (run *Run) next(serial bool) *RunSuite {
	for _, suite := range run.Suites {
		if suite.State == SuitePending && suite.Serial == serial {
			return suite
		}
	}
	return nil
}

func (run *Run) count(state string) int {
	n := 0
	for _, suite := range run.Suites {
		if suite.State == state {
			n++
		}
	}
	return n
}

// lose puts back the suites a lost agent was running, up to Retry times each.
func (run *Run) lose(a *agentConn) {
	for _, name := range a.busy {
		suite := run.find(name)
		if suite == nil || suite.State != SuiteRunning || suite.Agent != a.name {
			continue
		}
		if run.State == RunRunning && suite.Attempts <= run.Retry {
			suite.State, suite.Agent, suite.Seat = SuitePending, "", 0
			continue
		}
		suite.State = SuiteLost
		suite.Seconds = time.Since(suite.StartedAt).Seconds()
		suite.Failures = append(suite.Failures, "Agent "+a.name+" 연결 끊김")
	}
}

type startRequest struct {
	Suites []string `json:"suites"`
	Retry  int      `json:"retry"`
}

func (m *Master) start(w http.ResponseWriter, r *http.Request) {
	var req startRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		reply(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	if len(req.Suites) == 0 {
		reply(w, http.StatusBadRequest, map[string]string{"error": "실행할 스위트를 고르세요"})
		return
	}
	catalog, err := m.catalog()
	if err != nil {
		reply(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	m.mu.Lock()
	defer m.mu.Unlock()
	if m.run != nil && m.run.State != RunDone {
		reply(w, http.StatusConflict, map[string]string{"error": "이미 실행 중인 테스트가 있습니다"})
		return
	}
	if !slices.ContainsFunc(m.agents, func(a *agentConn) bool { return !a.lost }) {
		reply(w, http.StatusConflict, map[string]string{"error": "연결된 Agent가 없습니다"})
		return
	}

	// The run ID names the bots, so two runs never share one.
	id := time.Now().Unix()
	if m.run != nil {
		if last, err := strconv.ParseInt(m.run.ID, 36, 64); err == nil && id <= last {
			id = last + 1
		}
	}
	run := &Run{ID: strconv.FormatInt(id, 36), State: RunRunning, Retry: max(0, req.Retry), StartedAt: time.Now()}
	for _, serial := range []bool{false, true} {
		for _, suite := range catalog {
			if suite.Serial == serial && slices.Contains(req.Suites, suite.Name) {
				run.Suites = append(run.Suites, &RunSuite{Name: suite.Name, Title: suite.Title, Serial: suite.Serial, TimeoutMs: suite.TimeoutMs, State: SuitePending})
			}
		}
	}
	if len(run.Suites) == 0 {
		reply(w, http.StatusBadRequest, map[string]string{"error": "고른 스위트가 registry.lua에 없습니다"})
		return
	}
	m.run = run
	log.Printf("master: run %s starts %d suites", run.ID, len(run.Suites))
	m.advance()
	m.broadcast()
	reply(w, http.StatusAccepted, map[string]string{"id": run.ID})
}

func (m *Master) stop(w http.ResponseWriter, _ *http.Request) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.run == nil || m.run.State != RunRunning {
		reply(w, http.StatusConflict, map[string]string{"error": "실행 중인 테스트가 없습니다"})
		return
	}
	m.run.State = RunStopping
	for _, suite := range m.run.Suites {
		if suite.State == SuitePending {
			suite.State = SuiteStopped
		}
	}
	for _, a := range m.agents {
		if !a.lost && len(a.busy) > 0 {
			a.send(Message{Type: MessageStop, Run: m.run.ID})
		}
	}
	m.advance()
	m.broadcast()
	reply(w, http.StatusAccepted, map[string]string{"state": RunStopping})
}

// advance must be called with mu held. Parallel suites fill every free seat; a serial suite runs alone, after
// every parallel suite, because it uses state the whole server shares.
func (m *Master) advance() {
	run := m.run
	if run == nil || run.State == RunDone {
		return
	}
	if run.State == RunRunning {
		for _, a := range m.agents {
			for seat := 1; seat <= a.seats && !a.lost; seat++ {
				suite := run.next(false)
				if suite == nil {
					break
				}
				if a.busy[seat] == "" {
					m.assign(run, a, seat, suite)
				}
			}
		}
		if run.count(SuiteRunning) == 0 && run.next(false) == nil {
			if suite := run.next(true); suite != nil {
				index := slices.IndexFunc(m.agents, func(a *agentConn) bool { return !a.lost })
				if index >= 0 {
					m.assign(run, m.agents[index], 1, suite)
				}
			}
		}
	}
	if run.count(SuiteRunning) == 0 && (run.State == RunStopping || run.count(SuitePending) == 0) {
		m.finish(run)
	}
}

// assign must be called with mu held.
func (m *Master) assign(run *Run, a *agentConn, seat int, suite *RunSuite) {
	suite.State, suite.Agent, suite.Seat, suite.StartedAt = SuiteRunning, a.name, a.base+seat, time.Now()
	suite.Attempts++
	a.busy[seat] = suite.Name
	a.send(Message{
		Type:      MessageRunSuite,
		Run:       run.ID,
		Bots:      run.ID + strconv.FormatInt(int64(a.tag%36), 36),
		Suite:     suite.Name,
		Serial:    suite.Serial,
		TimeoutMs: suite.TimeoutMs,
		Seat:      suite.Seat,
	})
}

// finish must be called with mu held.
func (m *Master) finish(run *Run) {
	run.State, run.FinishedAt = RunDone, time.Now()
	log.Printf("master: run %s finished: %d passed, %d failed, %d lost", run.ID, run.count(SuitePassed), run.count(SuiteFailed), run.count(SuiteLost))
	if err := m.save(run); err != nil {
		log.Printf("master: save run %s: %v", run.ID, err)
	}
}

func (m *Master) save(run *Run) error {
	dir := filepath.Join(m.dir, run.ID)
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
	Agents     int       `json:"agents"`
	Total      int       `json:"total"`
	Passed     int       `json:"passed"`
	Failed     int       `json:"failed"`
	Skipped    int       `json:"skipped"`
	Stopped    int       `json:"stopped"`
	Lost       int       `json:"lost"`
}

func (m *Master) runs(w http.ResponseWriter, _ *http.Request) {
	entries, err := os.ReadDir(m.dir)
	if err != nil && !os.IsNotExist(err) {
		reply(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	rows := []runRow{}
	for _, entry := range entries {
		data, err := os.ReadFile(filepath.Join(m.dir, entry.Name(), "summary.json"))
		if err != nil {
			continue
		}
		var run Run
		if err := json.Unmarshal(data, &run); err != nil {
			continue
		}
		agents := make(map[string]bool)
		for _, suite := range run.Suites {
			if suite.Agent != "" {
				agents[suite.Agent] = true
			}
		}
		rows = append(rows, runRow{
			ID:         run.ID,
			StartedAt:  run.StartedAt,
			FinishedAt: run.FinishedAt,
			Agents:     len(agents),
			Total:      len(run.Suites),
			Passed:     run.count(SuitePassed),
			Failed:     run.count(SuiteFailed),
			Skipped:    run.count(SuiteSkipped),
			Stopped:    run.count(SuiteStopped),
			Lost:       run.count(SuiteLost),
		})
	}
	slices.SortFunc(rows, func(a, b runRow) int { return b.StartedAt.Compare(a.StartedAt) })
	reply(w, http.StatusOK, rows)
}

func (m *Master) history(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" || strings.Trim(id, "0123456789abcdefghijklmnopqrstuvwxyz") != "" {
		reply(w, http.StatusBadRequest, map[string]string{"error": "잘못된 실행 ID"})
		return
	}
	data, err := os.ReadFile(filepath.Join(m.dir, id, "summary.json"))
	if err != nil {
		reply(w, http.StatusNotFound, map[string]string{"error": "실행 기록이 없습니다"})
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_, _ = w.Write(data)
}
