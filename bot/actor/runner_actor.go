package actor

import (
	"log"
	"slices"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/common/config"
	"github.com/boyism80/fm/maple/wz"
)

type RunnerActor struct {
	cfg      *config.Bot
	wz       *wz.Resources
	atlas    *Atlas
	filter   string
	names    []string
	junit    string
	watch    func(msg any)
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

// NewRunnerActor runs the suites named in names, or the suites whose name contains filter when names is empty.
// watch, when set, is called on the runner's goroutine with RunStarted, SuiteStarted and SuiteFinished.
func NewRunnerActor(cfg *config.Bot, resources *wz.Resources, filter string, names []string, junit string, watch func(msg any), done chan<- int) *RunnerActor {
	r := &RunnerActor{
		cfg:    cfg,
		wz:     resources,
		atlas:  NewAtlas(resources),
		filter: filter,
		names:  names,
		junit:  junit,
		watch:  watch,
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
		if r.watch != nil {
			r.watch(&RunStarted{ID: r.runID, Parallel: slices.Clone(r.parallel), Serial: slices.Clone(r.serial)})
		}
		r.schedule(ctx)
	case *SuiteFinished:
		r.results = append(r.results, msg)
		if r.watch != nil {
			r.watch(msg)
		}
		if !msg.Serial {
			r.seats = append(r.seats, msg.Seat)
		}
		r.running--
		r.schedule(ctx)
	case *RunStop:
		// The run ends once every running suite reports SuiteFinished.
		r.parallel, r.serial = nil, nil
		for _, child := range ctx.Children() {
			ctx.Send(child, &SuiteStop{})
		}
	}
}

func (r *RunnerActor) loadRegistry() error {
	suites, err := LoadSuites(r.cfg.ScriptDir)
	if err != nil {
		return err
	}
	for _, s := range suites {
		if len(r.names) > 0 && !slices.Contains(r.names, s.Name) {
			continue
		}
		if len(r.names) == 0 && (!strings.Contains(s.Name, r.filter) || (s.Explicit && r.filter == "")) {
			continue
		}
		if s.Serial {
			r.serial = append(r.serial, s)
		} else {
			r.parallel = append(r.parallel, s)
		}
	}
	return nil
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
	if r.watch != nil {
		r.watch(&SuiteStarted{Suite: s.Name, Seat: seat})
	}
	ctx.Spawn(actor.PropsFromProducer(func() actor.Actor {
		return NewSuiteActor(r.cfg, r.wz, r.atlas, s, seat, r.runID, &r.botSeq)
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
