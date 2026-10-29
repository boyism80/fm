package cluster

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	botactor "github.com/boyism80/fm/bot/actor"
	"github.com/boyism80/fm/common/config"
	"github.com/boyism80/fm/maple/wz"
	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
)

const heartbeatInterval = 2 * time.Second

// Agent connects to a master and runs the suites it assigns. It keeps one AgentActor across reconnects so a
// suite stopped by a lost connection still logs its bots out.
type Agent struct {
	cfg    *config.Bot
	root   *actor.RootContext
	runner *actor.PID
	master string
	token  string
	name   string
	mu     sync.Mutex
	conn   *websocket.Conn
}

func NewAgent(cfg *config.Bot, resources *wz.Resources, root *actor.RootContext, master, token, name string) *Agent {
	a := &Agent{cfg: cfg, root: root, master: master, token: token, name: name}
	a.runner = root.Spawn(actor.PropsFromProducer(func() actor.Actor {
		return botactor.NewAgentActor(cfg, resources, a.report)
	}))
	return a
}

// Run keeps a connection to the master, waiting 1, 2, 4 ... 30 seconds between attempts.
func (a *Agent) Run() {
	wait := time.Second
	for {
		connected, err := a.session()
		if connected {
			wait = time.Second
		}
		log.Printf("agent %s: %v; reconnecting in %s", a.name, err, wait)
		time.Sleep(wait)
		wait = min(wait*2, 30*time.Second)
	}
}

func (a *Agent) session() (bool, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	conn, _, err := websocket.Dial(ctx, a.master+"/ws/agent", nil)
	if err != nil {
		return false, err
	}
	defer conn.CloseNow()
	conn.SetReadLimit(1 << 20)

	host, _ := os.Hostname()
	hello := Message{Type: MessageHello, Token: a.token, Name: a.name, Host: host, PID: os.Getpid(), Seats: a.cfg.Seats}
	if err := wsjson.Write(ctx, conn, hello); err != nil {
		return false, err
	}
	var welcome Message
	if err := wsjson.Read(ctx, conn, &welcome); err != nil {
		return false, err
	}
	if welcome.Type != MessageWelcome {
		return false, fmt.Errorf("master refused: %s", welcome.Error)
	}
	log.Printf("agent %s: connected to %s as #%d, seats %d~%d", a.name, a.master, welcome.AgentID, welcome.SeatBase+1, welcome.SeatBase+a.cfg.Seats)

	a.mu.Lock()
	a.conn = conn
	a.mu.Unlock()
	// Suites cannot report to a master that is gone, so they stop and log their bots out.
	defer func() {
		a.mu.Lock()
		a.conn = nil
		a.mu.Unlock()
		a.root.Send(a.runner, &botactor.RunStop{})
	}()

	go func() {
		ticker := time.NewTicker(heartbeatInterval)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				a.send(Message{Type: MessageHeartbeat})
			}
		}
	}()

	for {
		var msg Message
		if err := wsjson.Read(ctx, conn, &msg); err != nil {
			return true, err
		}
		switch msg.Type {
		case MessageRunSuite:
			suite := &botactor.Suite{Name: msg.Suite, Serial: msg.Serial, TimeoutMs: msg.TimeoutMs}
			a.root.Send(a.runner, &botactor.AssignSuite{Run: msg.Run, Bots: msg.Bots, Suite: suite, Seat: msg.Seat})
		case MessageStop:
			a.root.Send(a.runner, &botactor.RunStop{})
		case MessageQuit:
			log.Printf("agent %s: the master asked this agent to quit", a.name)
			os.Exit(0)
		}
	}
}

func (a *Agent) send(msg Message) {
	a.mu.Lock()
	conn := a.conn
	a.mu.Unlock()
	if conn == nil {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := wsjson.Write(ctx, conn, msg); err != nil {
		log.Printf("agent %s: send %s: %v", a.name, msg.Type, err)
	}
}

func (a *Agent) report(run string, msg any) {
	switch msg := msg.(type) {
	case *botactor.SuiteStarted:
		a.send(Message{Type: MessageSuiteStarted, Run: run, Suite: msg.Suite, Seat: msg.Seat})
	case *botactor.SuiteFinished:
		state := SuitePassed
		switch {
		case msg.Stopped:
			state = SuiteStopped
		case !msg.Passed:
			state = SuiteFailed
		case msg.Skipped:
			state = SuiteSkipped
		}
		a.send(Message{
			Type:     MessageSuiteFinished,
			Run:      run,
			Suite:    msg.Suite,
			Seat:     msg.Seat,
			Title:    msg.Name,
			State:    state,
			Infra:    msg.Infra,
			Failures: msg.Failures,
			Log:      msg.Log,
			Seconds:  msg.Elapsed.Seconds(),
		})
	}
}
