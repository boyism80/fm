package cluster

import (
	"context"
	"log"
	"net"
	"net/http"
	"slices"
	"strconv"
	"time"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
)

const (
	AgentRun  = "run"
	AgentIdle = "idle"
	AgentSlow = "slow"
	AgentLost = "lost"

	agentSlowAfter = 6 * time.Second
)

type agentConn struct {
	id    int
	name  string
	host  string
	pid   int
	tag   int
	seats int
	base  int
	busy  map[int]string
	out   chan Message
	conn  *websocket.Conn
	seen  time.Time
	lost  bool
}

type agentRow struct {
	ID       int      `json:"id"`
	Name     string   `json:"name"`
	Host     string   `json:"host"`
	PID      int      `json:"pid"`
	State    string   `json:"state"`
	Seats    int      `json:"seats"`
	SeatBase int      `json:"seat_base"`
	Suites   []string `json:"suites"`
	Seen     float64  `json:"seen"`
}

func (a *agentConn) row() agentRow {
	row := agentRow{ID: a.id, Name: a.name, Host: a.host, PID: a.pid, Seats: a.seats, SeatBase: a.base, Suites: []string{}, Seen: time.Since(a.seen).Seconds()}
	for seat := 1; seat <= a.seats; seat++ {
		if name := a.busy[seat]; name != "" {
			row.Suites = append(row.Suites, name)
		}
	}
	switch {
	case a.lost:
		row.State = AgentLost
	case time.Since(a.seen) > agentSlowAfter:
		row.State = AgentSlow
	case len(a.busy) > 0:
		row.State = AgentRun
	default:
		row.State = AgentIdle
	}
	return row
}

// send must be called with mu held. An agent that cannot keep up is dropped, which puts its suites back.
func (a *agentConn) send(msg Message) {
	select {
	case a.out <- msg:
	default:
		log.Printf("master: agent %s is not reading; closing it", a.name)
		_ = a.conn.CloseNow()
	}
}

func (m *Master) accept(w http.ResponseWriter, r *http.Request) {
	conn, err := websocket.Accept(w, r, nil)
	if err != nil {
		return
	}
	defer conn.CloseNow()
	conn.SetReadLimit(4 << 20)
	ctx := r.Context()

	var hello Message
	helloCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	err = wsjson.Read(helloCtx, conn, &hello)
	cancel()
	if err != nil || hello.Type != MessageHello || hello.Seats < 1 {
		return
	}
	if !m.admits(r, hello.Token) {
		_ = wsjson.Write(ctx, conn, Message{Type: MessageRefused, Error: "토큰이 맞지 않거나 토큰 없이 외부에서 연결했습니다"})
		return
	}

	a := m.join(hello, conn)
	defer m.leave(a)
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case msg := <-a.out:
				if err := wsjson.Write(ctx, conn, msg); err != nil {
					_ = conn.CloseNow()
					return
				}
			}
		}
	}()

	for {
		var msg Message
		if err := wsjson.Read(ctx, conn, &msg); err != nil {
			return
		}
		m.receive(a, msg)
	}
}

// admits takes an agent that knows the token. Without a token only agents on this machine connect.
func (m *Master) admits(r *http.Request, token string) bool {
	if m.token != "" {
		return token == m.token
	}
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return false
	}
	ip := net.ParseIP(host)
	return ip != nil && ip.IsLoopback()
}

// join gives an agent that reconnects under the same name and host its old ID and seats, so its instance
// slots stay the same.
func (m *Master) join(hello Message, conn *websocket.Conn) *agentConn {
	m.mu.Lock()
	defer m.mu.Unlock()

	index := slices.IndexFunc(m.agents, func(a *agentConn) bool {
		return a.lost && a.name == hello.Name && a.host == hello.Host
	})
	var a *agentConn
	if index >= 0 {
		a = m.agents[index]
	} else {
		a = &agentConn{id: len(m.agents) + 1, name: hello.Name, host: hello.Host}
		m.agents = append(m.agents, a)
	}
	// A new process numbers its bots from zero again, so it gets a bot tag no earlier process used.
	if index < 0 || a.pid != hello.PID {
		a.tag = m.nextTag
		m.nextTag++
	}
	if a.seats < hello.Seats {
		a.base = m.nextBase
		m.nextBase += hello.Seats
	}
	a.seats, a.pid, a.conn, a.lost = hello.Seats, hello.PID, conn, false
	a.busy = make(map[int]string)
	a.out = make(chan Message, 256)
	a.seen = time.Now()
	a.send(Message{Type: MessageWelcome, AgentID: a.id, SeatBase: a.base})
	log.Printf("master: agent #%d %s (%s, pid %d) joined with seats %d~%d", a.id, a.name, a.host, a.pid, a.base+1, a.base+a.seats)

	m.advance()
	m.broadcast()
	return a
}

func (m *Master) leave(a *agentConn) {
	m.mu.Lock()
	defer m.mu.Unlock()
	a.lost = true
	log.Printf("master: agent #%d %s left", a.id, a.name)
	if m.run != nil {
		m.run.lose(a)
	}
	clear(a.busy)
	m.advance()
	m.broadcast()
}

func (m *Master) receive(a *agentConn, msg Message) {
	m.mu.Lock()
	defer m.mu.Unlock()
	a.seen = time.Now()
	if msg.Type != MessageSuiteFinished {
		return
	}

	local := msg.Seat - a.base
	if a.busy[local] == msg.Suite {
		delete(a.busy, local)
	}
	if m.run == nil || m.run.ID != msg.Run {
		m.advance()
		return
	}
	suite := m.run.find(msg.Suite)
	if suite == nil || suite.State != SuiteRunning || suite.Agent != a.name {
		m.advance()
		return
	}
	suite.State, suite.Title, suite.Infra, suite.Failures, suite.Log, suite.Seconds = msg.State, msg.Title, msg.Infra, msg.Failures, msg.Log, msg.Seconds
	m.advance()
	m.broadcast()
}

func (m *Master) quit(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.PathValue("id"))
	m.mu.Lock()
	defer m.mu.Unlock()
	index := slices.IndexFunc(m.agents, func(a *agentConn) bool { return a.id == id && !a.lost })
	if index < 0 {
		reply(w, http.StatusNotFound, map[string]string{"error": "연결된 Agent가 아닙니다"})
		return
	}
	m.agents[index].send(Message{Type: MessageQuit})
	reply(w, http.StatusAccepted, map[string]string{"state": "quitting"})
}
