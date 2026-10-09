package actor

import (
	"log"
	"path/filepath"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/common/config"
	"github.com/boyism80/fm/core/luax"
	"github.com/boyism80/fm/services/game/wz"
	lua "github.com/yuin/gopher-lua"
)

type RunnerActor struct {
	cfg      *config.Bot
	wz       *wz.Resources
	filter   string
	junit    string
	done     chan<- int
	runID    string
	botSeq   atomic.Int32
	parallel []*Suite
	serial   []*Suite
	seats    []int
	running  int
	results  []*SuiteFinished
	started  time.Time
}

func NewRunnerActor(cfg *config.Bot, resources *wz.Resources, filter, junit string, done chan<- int) *RunnerActor {
	r := &RunnerActor{
		cfg:    cfg,
		wz:     resources,
		filter: filter,
		junit:  junit,
		done:   done,
		runID:  strconv.FormatInt(time.Now().Unix(), 36),
	}
	for seat := cfg.Seats; seat >= 1; seat-- {
		r.seats = append(r.seats, seat)
	}
	return r
}

func (r *RunnerActor) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *actor.Started:
		r.started = time.Now()
		if err := r.loadRegistry(); err != nil {
			log.Printf("registry: %v", err)
			r.done <- 2
			return
		}
		if len(r.parallel)+len(r.serial) == 0 {
			log.Printf("no suite matches filter %q", r.filter)
			r.done <- 2
			return
		}
		log.Printf("run %s: %d parallel, %d serial suites", r.runID, len(r.parallel), len(r.serial))
		r.schedule(ctx)
	case *SuiteFinished:
		r.results = append(r.results, msg)
		if !msg.Serial {
			r.seats = append(r.seats, msg.Seat)
		}
		r.running--
		r.schedule(ctx)
	}
}

func (r *RunnerActor) loadRegistry() error {
	L := luax.NewState()
	defer L.Close()
	L.SetGlobal("register_test", L.NewFunction(func(L *lua.LState) int {
		s := &Suite{Name: L.CheckString(1)}
		if opts := L.OptTable(2, nil); opts != nil {
			s.Serial = lua.LVAsBool(opts.RawGetString("serial"))
			s.Explicit = lua.LVAsBool(opts.RawGetString("explicit"))
			if n, ok := opts.RawGetString("timeout_ms").(lua.LNumber); ok {
				s.TimeoutMs = int(n)
			}
		}
		if !strings.Contains(s.Name, r.filter) {
			return 0
		}
		if s.Explicit && r.filter == "" {
			return 0
		}
		if s.Serial {
			r.serial = append(r.serial, s)
		} else {
			r.parallel = append(r.parallel, s)
		}
		return 0
	}))
	return L.DoFile(filepath.Join(r.cfg.ScriptDir, "registry.lua"))
}

func (r *RunnerActor) schedule(ctx actor.Context) {
	for len(r.parallel) > 0 && len(r.seats) > 0 {
		seat := r.seats[len(r.seats)-1]
		r.seats = r.seats[:len(r.seats)-1]
		r.start(ctx, r.parallel[0], seat)
		r.parallel = r.parallel[1:]
	}
	if r.running > 0 {
		return
	}
	if len(r.serial) > 0 {
		r.start(ctx, r.serial[0], 1)
		r.serial = r.serial[1:]
		return
	}
	r.report()
	r.done <- r.exitCode()
}

func (r *RunnerActor) start(ctx actor.Context, s *Suite, seat int) {
	r.running++
	ctx.Spawn(actor.PropsFromProducer(func() actor.Actor {
		return NewSuiteActor(r.cfg, r.wz, s, seat, r.runID, &r.botSeq)
	}))
}

func (r *RunnerActor) exitCode() int {
	code := 0
	for _, res := range r.results {
		if res.Infra {
			return 2
		}
		if !res.Passed {
			code = 1
		}
	}
	return code
}
