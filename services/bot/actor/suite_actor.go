package actor

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync/atomic"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/common/config"
	"github.com/boyism80/fm/core/luax"
	"github.com/boyism80/fm/services/bot/bot"
	"github.com/boyism80/fm/services/bot/luamarshal"
	"github.com/boyism80/fm/services/game/wz"
	lua "github.com/yuin/gopher-lua"
)

type SuiteActor struct {
	cfg      *config.Bot
	wz       *wz.Resources
	suite    *Suite
	seat     int
	runID    string
	botSeq   *atomic.Int32
	name     string
	bots     []*bot.Bot
	entered  int
	actors   *actor.RootContext
	self     *actor.PID
	waiters  []*waiter
	nextWait uint64
	hooks    map[string]func(b *bot.Bot, pkt any)
	moving   map[*bot.Bot]func(err error)
	entering map[*bot.Bot]bool
	started  time.Time
	deadline *time.Timer
	failures []string
	report   *os.File
	infra    bool
	skipped  bool
	ended    bool
	finished bool
	closed   bool

	L        *lua.LState
	def      *lua.LTable
	marshal  *luamarshal.Marshal
	threads  map[*lua.LState]func(values []lua.LValue, err error)
	sleeping map[*lua.LState]bool
	ctxUD    *lua.LUserData
	botUDs   map[*bot.Bot]*lua.LUserData
}

func NewSuiteActor(cfg *config.Bot, resources *wz.Resources, suite *Suite, seat int, runID string, botSeq *atomic.Int32) *SuiteActor {
	return &SuiteActor{
		cfg:      cfg,
		wz:       resources,
		suite:    suite,
		seat:     seat,
		runID:    runID,
		botSeq:   botSeq,
		name:     suite.Name,
		hooks:    make(map[string]func(b *bot.Bot, pkt any)),
		moving:   make(map[*bot.Bot]func(err error)),
		entering: make(map[*bot.Bot]bool),
		marshal:  luamarshal.New(),
		threads:  make(map[*lua.LState]func(values []lua.LValue, err error)),
		sleeping: make(map[*lua.LState]bool),
		botUDs:   make(map[*bot.Bot]*lua.LUserData),
	}
}

func (a *SuiteActor) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *actor.Started:
		a.start(ctx)
	case *actor.Stopped:
		if a.L != nil {
			a.L.Close()
		}
		if a.report != nil {
			_ = a.report.Close()
		}
	case *BotEntered:
		a.botEntered(msg)
	case *ChannelEntered:
		a.channelEntered(msg)
	case *PacketReceived:
		a.packetReceived(msg)
	case *Disconnected:
		a.disconnected(msg)
	case *WaitTimeout:
		a.waitTimeout(msg)
	case *SleepDone:
		if a.sleeping[msg.Thread] == false {
			return
		}
		delete(a.sleeping, msg.Thread)
		a.wake(msg.Thread)
	case *SuiteDeadline:
		a.Fail(fmt.Sprintf("suite deadline %dms", a.timeoutMs()))
		a.cleanup()
	case *CleanupDone:
		a.close(ctx)
	}
}

func (a *SuiteActor) start(ctx actor.Context) {
	a.actors = ctx.ActorSystem().Root
	a.self = ctx.Self()
	a.started = time.Now()
	a.deadline = time.AfterFunc(time.Duration(a.timeoutMs())*time.Millisecond, func() {
		a.actors.Send(a.self, &SuiteDeadline{})
	})

	a.L = luax.NewState()
	a.register()
	path := filepath.Join(a.cfg.ScriptDir, a.suite.Name+"_test.lua")
	if err := a.L.DoFile(path); err != nil {
		a.Fail(fmt.Sprintf("load %s: %v", path, err))
		a.cleanup()
		return
	}
	if a.def == nil {
		a.Fail(fmt.Sprintf("%s does not call test_suite", path))
		a.cleanup()
		return
	}
	if name, ok := a.def.RawGetString("name").(lua.LString); ok {
		a.name = string(name)
	}

	if fn, ok := a.def.RawGetString("should_skip").(*lua.LFunction); ok {
		skip, err := luax.CallFunction(a.L, fn, a.ctxUD)
		if err != nil {
			a.Fail(fmt.Sprintf("should_skip: %v", err))
			a.cleanup()
			return
		}
		if lua.LVAsBool(skip) {
			a.skipped = true
			a.cleanup()
			return
		}
	}

	count := 1
	if n, ok := a.def.RawGetString("bot_count").(lua.LNumber); ok {
		count = int(n)
	}
	for i := 0; i < count; i++ {
		b, err := bot.New(a.cfg, a.runID, int(a.botSeq.Add(1)-1))
		if err != nil {
			a.Fail(err.Error())
			a.cleanup()
			return
		}
		a.bots = append(a.bots, b)
		a.botUDs[b] = a.newUserData(b, "bot")
	}
	log.Printf("[%s] seat %d starts with %d bots", a.name, a.seat, len(a.bots))

	for i, b := range a.bots {
		a.entering[b] = true
		go func() {
			a.actors.Send(a.self, &BotEntered{Bot: i, Err: b.Enter()})
		}()
	}
}

func (a *SuiteActor) botEntered(msg *BotEntered) {
	b := a.bots[msg.Bot]
	delete(a.entering, b)
	if a.finished {
		a.closeLate(b)
		return
	}
	if msg.Err != nil {
		if errors.Is(msg.Err, bot.ErrConnect) {
			a.infra = true
		}
		a.Fail(fmt.Sprintf("%s 접속 실패: %v", b.ID, msg.Err))
		a.cleanup()
		return
	}

	a.listen(msg.Bot)
	a.entered++
	if a.entered < len(a.bots) {
		return
	}

	fn, ok := a.def.RawGetString("on_initialize").(*lua.LFunction)
	if ok == false {
		a.runScenarios(1)
		return
	}
	a.spawn(fn, func(values []lua.LValue, err error) {
		if err != nil {
			a.Fail(fmt.Sprintf("on_initialize: %v", err))
			a.end()
			return
		}
		if len(values) > 0 && values[0] == lua.LFalse {
			a.Fail("on_initialize returned false")
			a.end()
			return
		}
		a.runScenarios(1)
	}, a.ctxUD)
}

func (a *SuiteActor) listen(i int) {
	a.bots[i].Listen(func(gen int, opcode uint16, body []byte) {
		a.actors.Send(a.self, &PacketReceived{Bot: i, Gen: gen, Opcode: opcode, Body: body})
	}, func(gen int, err error) {
		a.actors.Send(a.self, &Disconnected{Bot: i, Gen: gen, Err: err})
	})
}

func (a *SuiteActor) channelEntered(msg *ChannelEntered) {
	b := a.bots[msg.Bot]
	done := a.moving[b]
	delete(a.moving, b)
	if a.finished {
		a.closeLate(b)
		return
	}
	if msg.Err == nil {
		a.listen(msg.Bot)
	}
	done(msg.Err)
}

func (a *SuiteActor) runScenarios(i int) {
	scenarios, _ := a.def.RawGetString("scenarios").(*lua.LTable)
	if scenarios == nil || i > scenarios.Len() {
		a.end()
		return
	}

	a.callHook("on_scenario_started", a.ctxUD, lua.LNumber(i))
	failures := len(a.failures)
	done := func(passed bool) {
		a.callHook("on_scenario_finished", a.ctxUD, lua.LNumber(i), lua.LBool(passed))
		if passed == false {
			if len(a.failures) == failures {
				a.Fail(fmt.Sprintf("시나리오 %d 실패", i))
			}
			a.end()
			return
		}
		a.runScenarios(i + 1)
	}

	switch sc := scenarios.RawGetInt(i).(type) {
	case *lua.LFunction:
		a.runScenario(sc, done)
	case *lua.LTable:
		a.runParallel(sc, done)
	default:
		a.Fail(fmt.Sprintf("시나리오 %d: function or { parallel = ... } expected", i))
		a.end()
	}
}

func (a *SuiteActor) runScenario(fn *lua.LFunction, done func(bool)) {
	a.spawn(fn, func(values []lua.LValue, err error) {
		if err != nil {
			a.Fail(err.Error())
			done(false)
			return
		}
		done(len(values) > 0 && values[0] == lua.LTrue)
	}, a.ctxUD)
}

func (a *SuiteActor) runParallel(tbl *lua.LTable, done func(bool)) {
	queues, ok := tbl.RawGetString("parallel").(*lua.LTable)
	if ok == false || queues.Len() == 0 {
		a.Fail("parallel block needs a non-empty parallel list")
		done(false)
		return
	}

	left, passed := queues.Len(), true
	for q := 1; q <= queues.Len(); q++ {
		var fns []*lua.LFunction
		switch entry := queues.RawGetInt(q).(type) {
		case *lua.LFunction:
			fns = append(fns, entry)
		case *lua.LTable:
			for j := 1; j <= entry.Len(); j++ {
				if fn, ok := entry.RawGetInt(j).(*lua.LFunction); ok {
					fns = append(fns, fn)
				}
			}
		}
		a.runQueue(q, fns, 0, func(ok bool) {
			passed = passed && ok
			left--
			if left == 0 {
				done(passed)
			}
		})
	}
}

func (a *SuiteActor) runQueue(q int, fns []*lua.LFunction, j int, done func(bool)) {
	if j >= len(fns) {
		done(true)
		return
	}
	a.callHook("on_parallel_scenario_started", a.ctxUD, lua.LNumber(q), lua.LNumber(j+1))
	a.runScenario(fns[j], func(passed bool) {
		a.callHook("on_parallel_scenario_finished", a.ctxUD, lua.LNumber(q), lua.LNumber(j+1), lua.LBool(passed))
		if passed == false {
			done(false)
			return
		}
		a.runQueue(q, fns, j+1, done)
	})
}

func (a *SuiteActor) end() {
	if a.ended {
		return
	}
	a.ended = true
	fn, ok := a.def.RawGetString("on_finished").(*lua.LFunction)
	if ok == false {
		a.cleanup()
		return
	}
	a.spawn(fn, func(_ []lua.LValue, err error) {
		if err != nil {
			a.Fail(fmt.Sprintf("on_finished: %v", err))
		}
		a.cleanup()
	}, a.ctxUD)
}

func (a *SuiteActor) callHook(name string, args ...lua.LValue) {
	fn, ok := a.def.RawGetString(name).(*lua.LFunction)
	if ok == false {
		return
	}
	values := make([]interface{}, len(args))
	for i, v := range args {
		values[i] = v
	}
	if _, err := luax.CallFunction(a.L, fn, values...); err != nil {
		a.Fail(fmt.Sprintf("%s: %v", name, err))
	}
}

func (a *SuiteActor) timeoutMs() int {
	if a.suite.TimeoutMs > 0 {
		return a.suite.TimeoutMs
	}
	return a.cfg.SuiteTimeoutMs
}

func (a *SuiteActor) Fail(msg string) bool {
	a.failures = append(a.failures, msg)
	return false
}
