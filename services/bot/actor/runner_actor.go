package actor

import (
	"log"
	"strconv"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/common/config"
	"github.com/boyism80/fm/services/bot/bot"
	"github.com/boyism80/fm/services/game/wz"
)

type RunnerActor struct {
	cfg      *config.Bot
	wz       *wz.Resources
	junit    string
	done     chan<- int
	runID    string
	nextBot  int
	parallel []*Suite
	serial   []*Suite
	seats    []int
	running  int
	results  []*SuiteFinished
	started  time.Time
}

func NewRunnerActor(cfg *config.Bot, resources *wz.Resources, suites []*Suite, junit string, done chan<- int) *RunnerActor {
	r := &RunnerActor{
		cfg:   cfg,
		wz:    resources,
		junit: junit,
		done:  done,
		runID: strconv.FormatInt(time.Now().Unix(), 36),
	}
	for _, s := range suites {
		if s.Serial {
			r.serial = append(r.serial, s)
		} else {
			r.parallel = append(r.parallel, s)
		}
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
		log.Printf("run %s: %d parallel, %d serial suites", r.runID, len(r.parallel), len(r.serial))
		r.schedule(ctx)
	case *SuiteFinished:
		r.results = append(r.results, msg)
		if msg.Serial == false {
			r.seats = append(r.seats, msg.Seat)
		}
		r.running--
		r.schedule(ctx)
	}
}

func (r *RunnerActor) schedule(ctx actor.Context) {
	for len(r.parallel) > 0 && len(r.seats) > 0 {
		s := r.parallel[0]
		r.parallel = r.parallel[1:]
		seat := r.seats[len(r.seats)-1]
		if r.start(ctx, s, seat) {
			r.seats = r.seats[:len(r.seats)-1]
		}
	}
	if r.running > 0 {
		return
	}
	for len(r.serial) > 0 {
		s := r.serial[0]
		r.serial = r.serial[1:]
		if r.start(ctx, s, 1) {
			return
		}
	}
	r.report()
	r.done <- r.exitCode()
}

func (r *RunnerActor) start(ctx actor.Context, s *Suite, seat int) bool {
	bots := make([]*bot.Bot, 0, s.BotCount)
	for i := 0; i < s.BotCount; i++ {
		b, err := bot.New(r.cfg, r.runID, r.nextBot)
		if err != nil {
			r.results = append(r.results, &SuiteFinished{Name: s.Name, Failures: []string{err.Error()}})
			return false
		}
		r.nextBot++
		bots = append(bots, b)
	}
	r.running++
	ctx.Spawn(actor.PropsFromProducer(func() actor.Actor {
		return NewSuiteActor(r.cfg, r.wz, s, seat, bots)
	}))
	return true
}

func (r *RunnerActor) exitCode() int {
	code := 0
	for _, res := range r.results {
		if res.Infra {
			return 2
		}
		if res.Passed == false {
			code = 1
		}
	}
	return code
}
