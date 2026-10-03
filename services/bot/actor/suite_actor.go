package actor

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/common/config"
	"github.com/boyism80/fm/services/bot/bot"
	"github.com/boyism80/fm/services/game/wz"
)

type SuiteActor struct {
	cfg      *config.Bot
	wz       *wz.Resources
	suite    *Suite
	seat     int
	bots     []*bot.Bot
	entered  int
	root     *actor.RootContext
	self     *actor.PID
	waiters  []*waiter
	nextWait uint64
	hooks    map[string]func(b *bot.Bot, pkt any)
	started  time.Time
	deadline *time.Timer
	failures []string
	infra    bool
	finished bool
}

func NewSuiteActor(cfg *config.Bot, resources *wz.Resources, suite *Suite, seat int, bots []*bot.Bot) *SuiteActor {
	return &SuiteActor{
		cfg:   cfg,
		wz:    resources,
		suite: suite,
		seat:  seat,
		bots:  bots,
		hooks: make(map[string]func(b *bot.Bot, pkt any)),
	}
}

func (a *SuiteActor) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *actor.Started:
		a.start(ctx)
	case *BotEntered:
		a.botEntered(msg)
	case *PacketReceived:
		a.packetReceived(msg)
	case *Disconnected:
		a.disconnected(msg)
	case *WaitTimeout:
		a.waitTimeout(msg)
	case *SleepDone:
		if a.finished == false {
			msg.Done()
		}
	case *SuiteDeadline:
		a.Fail(fmt.Sprintf("suite deadline %dms", a.cfg.SuiteTimeoutMs))
		a.finish()
	case *CleanupDone:
		a.close(ctx)
	}
}

func (a *SuiteActor) start(ctx actor.Context) {
	a.root = ctx.ActorSystem().Root
	a.self = ctx.Self()
	a.started = time.Now()
	a.deadline = time.AfterFunc(time.Duration(a.cfg.SuiteTimeoutMs)*time.Millisecond, func() {
		a.root.Send(a.self, &SuiteDeadline{})
	})
	log.Printf("[%s] seat %d starts with %d bots", a.suite.Name, a.seat, len(a.bots))

	for i, b := range a.bots {
		go func() {
			a.root.Send(a.self, &BotEntered{Bot: i, Err: b.Enter()})
		}()
	}
}

func (a *SuiteActor) botEntered(msg *BotEntered) {
	if a.finished {
		return
	}
	b := a.bots[msg.Bot]
	if msg.Err != nil {
		if errors.Is(msg.Err, bot.ErrConnect) {
			a.infra = true
		}
		a.Fail(fmt.Sprintf("%s 접속 실패: %v", b.ID, msg.Err))
		a.finish()
		return
	}

	b.Listen(func(gen int, opcode uint16, body []byte) {
		a.root.Send(a.self, &PacketReceived{Bot: msg.Bot, Gen: gen, Opcode: opcode, Body: body})
	}, func(gen int, err error) {
		a.root.Send(a.self, &Disconnected{Bot: msg.Bot, Gen: gen, Err: err})
	})
	a.entered++
	if a.entered < len(a.bots) {
		return
	}
	a.runScenarios(0)
}

func (a *SuiteActor) runScenarios(i int) {
	if a.finished {
		return
	}
	if i >= len(a.suite.Scenarios) {
		a.finish()
		return
	}
	sc := a.suite.Scenarios[i]
	failures := len(a.failures)
	a.runScenario(sc, func(passed bool) {
		if a.finished {
			return
		}
		if passed == false {
			if len(a.failures) == failures {
				a.Fail(fmt.Sprintf("시나리오 %s 실패", sc.Name))
			}
			a.finish()
			return
		}
		a.runScenarios(i + 1)
	})
}

func (a *SuiteActor) runScenario(sc Scenario, done func(bool)) {
	if len(sc.Parallel) == 0 {
		sc.Run(a, done)
		return
	}
	left, passed := len(sc.Parallel), true
	for _, p := range sc.Parallel {
		a.runScenario(p, func(ok bool) {
			passed = passed && ok
			left--
			if left == 0 {
				done(passed)
			}
		})
	}
}

func (a *SuiteActor) Fail(msg string) bool {
	a.failures = append(a.failures, msg)
	return false
}
