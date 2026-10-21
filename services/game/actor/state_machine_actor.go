package actor

import (
	"fmt"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/scheduler"
	c_actor "github.com/boyism80/fm/core/actor"
	"github.com/boyism80/fm/core/luax"
	"github.com/boyism80/fm/services/game"
	"github.com/robfig/cron/v3"
	lua "github.com/yuin/gopher-lua"
)

type namedSchedule struct {
	version uint64
	cancel  scheduler.CancelFunc
	hook    string
	cron    cron.Schedule
}

type smLifecycle struct {
	ready    bool
	deferred []interface{}
	stopping bool
	stopped  bool
}

type StateMachineActor struct {
	GameLogicActor
	StateMachine   *game.StateMachine
	luaRoot        *lua.LState
	timeoutVersion uint64
	timeoutCancel  scheduler.CancelFunc
	namedSchedules map[string]*namedSchedule
	lc             smLifecycle
}

func NewStateMachineActor(sm *game.StateMachine, gameWorld game.GameWorld) *StateMachineActor {
	a := &StateMachineActor{
		StateMachine:   sm,
		namedSchedules: make(map[string]*namedSchedule),
	}
	a.GameWorld = gameWorld
	a.maps = func() []*game.Map {
		if sm == nil {
			return nil
		}
		return sm.MapList()
	}
	return a
}

func (a *StateMachineActor) Receive(ctx actor.Context) {
	if a.lc.stopped {
		switch ctx.Message().(type) {
		case *c_actor.HandlePacket, *actor.Stopping, *actor.Stopped:
		default:
			return
		}
	}

	switch msg := ctx.Message().(type) {
	case *game.BootstrapStateMachine:
		a.beginCreate(ctx)
	case *game.CreateStateMachineMaps:
		a.createMaps(ctx, msg)
	case *game.StateMachineMapRemoved:
		a.handleMapRemoved(msg)
	case *game.EnterStateMachinePlayer:
		if a.lc.ready {
			a.handleEnterPlayer(ctx, msg)
		} else {
			a.lc.deferred = append(a.lc.deferred, msg)
		}
	case *game.StartStateMachine:
		if a.lc.ready {
			a.handleStart(ctx)
		} else {
			a.lc.deferred = append(a.lc.deferred, msg)
		}
	case *game.CallStateMachineHook:
		if a.lc.ready {
			a.callHook(ctx, msg.Hook, msg.Args...)
		} else {
			a.lc.deferred = append(a.lc.deferred, msg)
		}
	case *game.LeaveStateMachinePlayer:
		if msg != nil && a.StateMachine != nil {
			finished := a.StateMachine.LeavePlayer(ctx, msg.Character, msg.WarpLeaver, msg.Reason)
			if !finished && msg.Reason == game.StateMachineLeaveParty {
				a.callHook(ctx, "on_left_party", msg.Character)
			}
		}
	case *game.FinishStateMachine:
		if a.StateMachine != nil {
			a.StateMachine.Finish(ctx, msg.ExitMapID, msg.ExitPortal)
		}
	case *game.ScheduleStateMachineTimeout:
		if a.scheduler != nil && a.StateMachine != nil {
			a.timeoutVersion++
			if a.timeoutCancel != nil {
				a.timeoutCancel()
			}
			ms := a.StateMachine.LimitTimer(msg.Milliseconds)
			deadline := time.Now().Add(time.Duration(ms) * time.Millisecond)
			a.StateMachine.SetTimeoutDeadline(deadline)
			a.timeoutCancel = a.scheduler.SendOnce(time.Duration(ms)*time.Millisecond, ctx.Self(), &game.StateMachineTimeout{
				Version: a.timeoutVersion,
			})
			if ms > 0 {
				a.StateMachine.BroadcastClock(int32(ms / 1000))
			}
		}
		if msg.ReplyTo != nil {
			ctx.Send(msg.ReplyTo, &game.StateMachineTimerAck{})
		}
	case *game.CancelStateMachineTimeout:
		a.timeoutVersion++
		if a.timeoutCancel != nil {
			a.timeoutCancel()
			a.timeoutCancel = nil
		}
		if a.StateMachine != nil {
			a.StateMachine.SetTimeoutDeadline(time.Time{})
		}
		if msg.ReplyTo != nil {
			ctx.Send(msg.ReplyTo, &game.StateMachineTimerAck{})
		}
	case *game.StateMachineTimeout:
		if msg.Version == a.timeoutVersion {
			if a.StateMachine != nil {
				a.StateMachine.SetTimeoutDeadline(time.Time{})
			}
			a.callHook(ctx, "on_scheduled_timeout")
		}
	case *game.ScheduleStateMachineAfter:
		a.scheduleHook(ctx, msg)
	case *game.ScheduleStateMachineCron:
		a.scheduleCron(ctx, msg)
	case *game.CancelStateMachineNamedSchedule:
		if msg.ID == "" {
			a.cancelAllNamedSchedules()
		} else {
			a.cancelNamedSchedule(msg.ID)
		}
		if msg.ReplyTo != nil {
			ctx.Send(msg.ReplyTo, &game.StateMachineTimerAck{})
		}
	case *game.StateMachineNamedTimeout:
		a.handleNamedTimeout(ctx, msg)
	case *game.StopStateMachine:
		a.beginStop()
	case *actor.Restarting:
		a.StopTimers()
	case *actor.Stopped:
		a.StopTimers()
		if a.timeoutCancel != nil {
			a.timeoutCancel()
			a.timeoutCancel = nil
		}
		a.cancelAllNamedSchedules()
		if a.luaRoot != nil {
			a.luaRoot.Close()
			a.luaRoot = nil
		}
		a.GameLogicActor.Receive(ctx)
	default:
		a.GameLogicActor.Receive(ctx)
	}
}

func (a *StateMachineActor) scheduleHook(ctx actor.Context, msg *game.ScheduleStateMachineAfter) {
	if msg == nil {
		return
	}
	if a.scheduler == nil || msg.ID == "" || msg.Hook == "" || msg.Milliseconds <= 0 {
		if msg.ReplyTo != nil {
			ctx.Send(msg.ReplyTo, &game.StateMachineTimerAck{})
		}
		return
	}
	a.cancelNamedSchedule(msg.ID)
	entry := &namedSchedule{hook: msg.Hook}
	entry.version++
	entry.cancel = a.scheduler.SendOnce(time.Duration(a.StateMachine.LimitTimer(msg.Milliseconds))*time.Millisecond, ctx.Self(), &game.StateMachineNamedTimeout{
		ID:      msg.ID,
		Version: entry.version,
	})
	a.namedSchedules[msg.ID] = entry
	if msg.ReplyTo != nil {
		ctx.Send(msg.ReplyTo, &game.StateMachineTimerAck{})
	}
}

func (a *StateMachineActor) scheduleCron(ctx actor.Context, msg *game.ScheduleStateMachineCron) {
	if msg == nil {
		return
	}
	if a.scheduler == nil || msg.ID == "" || msg.Hook == "" || msg.Expr == "" {
		if msg.ReplyTo != nil {
			ctx.Send(msg.ReplyTo, &game.StateMachineTimerAck{})
		}
		return
	}
	sched, err := game.ParseStateMachineCron(msg.Expr)
	if err != nil {
		fmt.Printf("state machine cron %s: %v\n", msg.Expr, err)
		if msg.ReplyTo != nil {
			ctx.Send(msg.ReplyTo, &game.StateMachineTimerAck{})
		}
		return
	}
	a.cancelNamedSchedule(msg.ID)
	entry := &namedSchedule{hook: msg.Hook, cron: sched}
	a.namedSchedules[msg.ID] = entry
	a.scheduleNamedCron(ctx, msg.ID, entry)
	if msg.ReplyTo != nil {
		ctx.Send(msg.ReplyTo, &game.StateMachineTimerAck{})
	}
}

func (a *StateMachineActor) scheduleNamedCron(ctx actor.Context, id string, entry *namedSchedule) {
	if a.scheduler == nil || entry == nil || entry.cron == nil {
		return
	}
	delay := game.NextStateMachineCronDelay(entry.cron, time.Now())
	if delay <= 0 {
		return
	}
	entry.version++
	entry.cancel = a.scheduler.SendOnce(delay, ctx.Self(), &game.StateMachineNamedTimeout{
		ID:      id,
		Version: entry.version,
	})
}

func (a *StateMachineActor) cancelNamedSchedule(id string) {
	if a.namedSchedules == nil {
		return
	}
	entry := a.namedSchedules[id]
	if entry == nil {
		return
	}
	entry.version++
	if entry.cancel != nil {
		entry.cancel()
		entry.cancel = nil
	}
	delete(a.namedSchedules, id)
}

func (a *StateMachineActor) cancelAllNamedSchedules() {
	if a.namedSchedules == nil {
		return
	}
	for id := range a.namedSchedules {
		a.cancelNamedSchedule(id)
	}
}

func (a *StateMachineActor) handleNamedTimeout(ctx actor.Context, msg *game.StateMachineNamedTimeout) {
	if msg == nil || a.namedSchedules == nil {
		return
	}
	entry := a.namedSchedules[msg.ID]
	if entry == nil || msg.Version != entry.version {
		return
	}
	hook := entry.hook
	if entry.cron != nil {
		entry.cancel = nil
		a.scheduleNamedCron(ctx, msg.ID, entry)
	} else {
		delete(a.namedSchedules, msg.ID)
	}
	a.callHook(ctx, hook)
}

func (a *StateMachineActor) beginCreate(ctx actor.Context) {
	if a.StateMachine == nil || a.StateMachine.Group() == nil {
		return
	}
	if a.luaRoot == nil {
		a.luaRoot = luax.NewState()
	}
	thread, err := luax.NewThread(a.luaRoot, a.StateMachine.Group().ScriptPath())
	if err != nil {
		a.abortCreate(err.Error())
		return
	}
	if !luax.HasFunc(thread, "on_create") {
		luax.Close(thread)
		a.abortCreate("on_create is required")
		return
	}
	luax.SetConfiguration(thread, luax.Configuration{
		ActorContext: ctx,
		ActorPID:     ctx.Self(),
	})

	root := ctx.ActorSystem().Root
	self := ctx.Self()
	luax.CallAsync(ctx, a.luaRoot, thread, "on_create", a.StateMachine).Do(func(vals []lua.LValue) error {
		msg := &game.CreateStateMachineMaps{}
		if len(vals) == 0 || vals[0] == nil || vals[0] == lua.LNil {
			msg.Err = "on_create must return a map id array"
		} else if tbl, ok := vals[0].(*lua.LTable); !ok {
			msg.Err = "on_create must return a table"
		} else if specs, err := game.ParseCreateMaps(tbl); err != nil {
			msg.Err = err.Error()
		} else {
			msg.Specs = specs
		}
		root.Send(self, msg)
		return nil
	}).OnError(func(err error) {
		root.Send(self, &game.CreateStateMachineMaps{Err: err.Error()})
	})
}

func (a *StateMachineActor) createMaps(ctx actor.Context, msg *game.CreateStateMachineMaps) {
	if msg.Err != "" {
		a.abortCreate(msg.Err)
		return
	}
	if a.StateMachine.Disposed() {
		a.abortCreate("state machine disposed")
		return
	}
	if err := a.StateMachine.CreateMaps(msg.Specs); err != nil {
		a.abortCreate(err.Error())
		return
	}

	a.lc.ready = true
	a.callHook(ctx, "on_prepare")
	deferred := a.lc.deferred
	a.lc.deferred = nil
	for _, msg := range deferred {
		switch m := msg.(type) {
		case *game.EnterStateMachinePlayer:
			a.handleEnterPlayer(ctx, m)
		case *game.StartStateMachine:
			a.handleStart(ctx)
		case *game.CallStateMachineHook:
			a.callHook(ctx, m.Hook, m.Args...)
		}
	}
}

func (a *StateMachineActor) abortCreate(reason string) {
	fmt.Printf("state machine %s create failed: %s\n", a.StateMachine.Group().Name(), reason)
	a.lc.deferred = nil
	a.StateMachine.AbortStart()
	a.beginStop()
}

func (a *StateMachineActor) handleEnterPlayer(ctx actor.Context, msg *game.EnterStateMachinePlayer) {
	if msg == nil || msg.Character == nil || a.StateMachine == nil {
		return
	}
	a.StateMachine.Register(msg.Character)
	a.callHook(ctx, "on_player_enter", append([]interface{}{msg.Character}, msg.Args...)...)
}

func (a *StateMachineActor) handleStart(ctx actor.Context) {
	a.callHook(ctx, "on_start")
}

func (a *StateMachineActor) beginStop() {
	if a.lc.stopping {
		return
	}
	a.lc.stopping = true
	for _, m := range a.StateMachine.MapList() {
		m.ClearObjectTimers()
	}
	a.timeoutVersion++
	if a.timeoutCancel != nil {
		a.timeoutCancel()
		a.timeoutCancel = nil
	}
	a.cancelAllNamedSchedules()
	a.StateMachine.SetTimeoutDeadline(time.Time{})

	if len(a.StateMachine.MapList()) == 0 {
		a.finalizeStop()
		return
	}
	a.StateMachine.CloseMaps()
}

func (a *StateMachineActor) handleMapRemoved(msg *game.StateMachineMapRemoved) {
	if !a.StateMachine.RemoveMap(msg.Map) {
		return
	}
	if a.lc.stopping {
		a.finalizeStop()
	}
}

func (a *StateMachineActor) finalizeStop() {
	if a.lc.stopped {
		return
	}
	a.lc.stopped = true
	a.StateMachine.Group().RemoveMachine(a.StateMachine)
	a.GameWorld.Systems().Get[*game.StateMachineSystem]().StopActor(a.StateMachine)
	a.GameWorld.Systems().Get[*game.ExpeditionSystem]().EndBattle(a.StateMachine)
}

func (a *StateMachineActor) callHook(ctx actor.Context, hook string, args ...interface{}) {
	if a.StateMachine == nil || a.StateMachine.Group() == nil || hook == "" {
		return
	}
	if a.StateMachine.Disposed() {
		return
	}
	if a.luaRoot == nil {
		a.luaRoot = luax.NewState()
	}
	thread, err := luax.NewThread(a.luaRoot, a.StateMachine.Group().ScriptPath())
	if err != nil {
		return
	}
	if !luax.HasFunc(thread, hook) {
		luax.Close(thread)
		return
	}
	luax.SetConfiguration(thread, luax.Configuration{
		ActorContext: ctx,
		ActorPID:     ctx.Self(),
	})
	luax.CallAsync(ctx, a.luaRoot, thread, hook, append([]interface{}{a.StateMachine}, args...)...).OnError(func(err error) {
		fmt.Printf("state machine %s hook %s: %v\n", a.StateMachine.Group().Name(), hook, err)
	})
}
