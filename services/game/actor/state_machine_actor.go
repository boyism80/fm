package actor

import (
	"fmt"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/scheduler"
	c_actor "github.com/boyism80/fm/core/actor"
	"github.com/boyism80/fm/core/luax"
	"github.com/boyism80/fm/services/game/entity"
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
	attachDone    bool          // all map attaches have been acked
	deferred      []interface{} // messages queued before attachDone
	attachPending int           // outstanding AttachStateMachine acks
	detachPending int           // outstanding DetachStateMachine acks
	detaching     bool          // beginDetach has been called
	stopped       bool          // finalizeStop has run; only packets are still forwarded
}

type StateMachineActor struct {
	GameLogicActor
	StateMachine   *entity.StateMachine
	luaRoot        *lua.LState
	timeoutVersion uint64
	timeoutCancel  scheduler.CancelFunc
	namedSchedules map[string]*namedSchedule
	lc             smLifecycle
}

func NewStateMachineActor(sm *entity.StateMachine, gameWorld entity.GameWorld) *StateMachineActor {
	a := &StateMachineActor{
		StateMachine:   sm,
		namedSchedules: make(map[string]*namedSchedule),
	}
	a.GameWorld = gameWorld
	a.maps = func() []*entity.Map {
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
	case *entity.BootstrapStateMachine:
		a.beginCreate(ctx)
	case *entity.ApplyStateMachineCreateMaps:
		a.registerCreateMaps(ctx, msg)
	case *entity.FinishStateMachineCreate:
		a.finishCreate(ctx)
	case *entity.EnterStateMachinePlayer:
		if a.lc.attachDone {
			a.handleEnterPlayer(ctx, msg)
		} else {
			a.lc.deferred = append(a.lc.deferred, msg)
		}
	case *entity.StartStateMachine:
		if a.lc.attachDone {
			a.handleStart(ctx)
		} else {
			a.lc.deferred = append(a.lc.deferred, msg)
		}
	case *entity.CallStateMachineHook:
		if a.lc.attachDone {
			a.callHook(ctx, msg.Hook, msg.Args...)
		} else {
			a.lc.deferred = append(a.lc.deferred, msg)
		}
	case *entity.LeaveStateMachinePlayer:
		if msg != nil && a.StateMachine != nil {
			a.StateMachine.LeavePlayer(ctx, msg.Character, msg.WarpLeaver, msg.Reason)
		}
	case *entity.ScheduleStateMachineTimeout:
		if a.scheduler != nil && a.StateMachine != nil {
			a.timeoutVersion++
			if a.timeoutCancel != nil {
				a.timeoutCancel()
			}
			ms := a.StateMachine.LimitTimer(msg.Milliseconds)
			deadline := time.Now().Add(time.Duration(ms) * time.Millisecond)
			a.StateMachine.SetTimeoutDeadline(deadline)
			a.timeoutCancel = a.scheduler.SendOnce(time.Duration(ms)*time.Millisecond, ctx.Self(), &entity.StateMachineTimeout{
				Version: a.timeoutVersion,
			})
			if ms > 0 {
				a.StateMachine.BroadcastClock(int32(ms / 1000))
			}
		}
		if msg.ReplyTo != nil {
			ctx.Send(msg.ReplyTo, &entity.StateMachineTimerAck{})
		}
	case *entity.CancelStateMachineTimeout:
		a.timeoutVersion++
		if a.timeoutCancel != nil {
			a.timeoutCancel()
			a.timeoutCancel = nil
		}
		if a.StateMachine != nil {
			a.StateMachine.SetTimeoutDeadline(time.Time{})
		}
		if msg.ReplyTo != nil {
			ctx.Send(msg.ReplyTo, &entity.StateMachineTimerAck{})
		}
	case *entity.StateMachineTimeout:
		if msg.Version == a.timeoutVersion {
			if a.StateMachine != nil {
				a.StateMachine.SetTimeoutDeadline(time.Time{})
			}
			a.callHook(ctx, "on_scheduled_timeout")
		}
	case *entity.ScheduleStateMachineAfter:
		a.scheduleHook(ctx, msg)
	case *entity.ScheduleStateMachineCron:
		a.scheduleCron(ctx, msg)
	case *entity.CancelStateMachineNamedSchedule:
		if msg.ID == "" {
			a.cancelAllNamedSchedules()
		} else {
			a.cancelNamedSchedule(msg.ID)
		}
		if msg.ReplyTo != nil {
			ctx.Send(msg.ReplyTo, &entity.StateMachineTimerAck{})
		}
	case *entity.StateMachineNamedTimeout:
		a.handleNamedTimeout(ctx, msg)
	case *AttachStateMachineAck:
		a.handleAttachAck(ctx, msg)
	case *entity.StopStateMachine:
		a.beginDetach(ctx)
	case *DetachStateMachineAck:
		a.handleDetachAck(ctx)
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

func (a *StateMachineActor) scheduleHook(ctx actor.Context, msg *entity.ScheduleStateMachineAfter) {
	if msg == nil {
		return
	}
	if a.scheduler == nil || msg.ID == "" || msg.Hook == "" || msg.Milliseconds <= 0 {
		if msg.ReplyTo != nil {
			ctx.Send(msg.ReplyTo, &entity.StateMachineTimerAck{})
		}
		return
	}
	a.cancelNamedSchedule(msg.ID)
	entry := &namedSchedule{hook: msg.Hook}
	entry.version++
	entry.cancel = a.scheduler.SendOnce(time.Duration(a.StateMachine.LimitTimer(msg.Milliseconds))*time.Millisecond, ctx.Self(), &entity.StateMachineNamedTimeout{
		ID:      msg.ID,
		Version: entry.version,
	})
	a.namedSchedules[msg.ID] = entry
	if msg.ReplyTo != nil {
		ctx.Send(msg.ReplyTo, &entity.StateMachineTimerAck{})
	}
}

func (a *StateMachineActor) scheduleCron(ctx actor.Context, msg *entity.ScheduleStateMachineCron) {
	if msg == nil {
		return
	}
	if a.scheduler == nil || msg.ID == "" || msg.Hook == "" || msg.Expr == "" {
		if msg.ReplyTo != nil {
			ctx.Send(msg.ReplyTo, &entity.StateMachineTimerAck{})
		}
		return
	}
	sched, err := entity.ParseStateMachineCron(msg.Expr)
	if err != nil {
		fmt.Printf("state machine cron %s: %v\n", msg.Expr, err)
		if msg.ReplyTo != nil {
			ctx.Send(msg.ReplyTo, &entity.StateMachineTimerAck{})
		}
		return
	}
	a.cancelNamedSchedule(msg.ID)
	entry := &namedSchedule{hook: msg.Hook, cron: sched}
	a.namedSchedules[msg.ID] = entry
	a.scheduleNamedCron(ctx, msg.ID, entry)
	if msg.ReplyTo != nil {
		ctx.Send(msg.ReplyTo, &entity.StateMachineTimerAck{})
	}
}

func (a *StateMachineActor) scheduleNamedCron(ctx actor.Context, id string, entry *namedSchedule) {
	if a.scheduler == nil || entry == nil || entry.cron == nil {
		return
	}
	delay := entity.NextStateMachineCronDelay(entry.cron, time.Now())
	if delay <= 0 {
		return
	}
	entry.version++
	entry.cancel = a.scheduler.SendOnce(delay, ctx.Self(), &entity.StateMachineNamedTimeout{
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

func (a *StateMachineActor) handleNamedTimeout(ctx actor.Context, msg *entity.StateMachineNamedTimeout) {
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
	if a.StateMachine == nil || a.StateMachine.Group == nil {
		return
	}
	if a.luaRoot == nil {
		a.luaRoot = luax.NewState()
	}
	thread, err := luax.NewThread(a.luaRoot, a.StateMachine.Group.ScriptPath)
	if err != nil {
		a.abortCreate(ctx, err.Error())
		return
	}
	if luax.HasFunc(thread, "on_create") == false {
		luax.Close(thread)
		a.abortCreate(ctx, "on_create is required")
		return
	}
	luax.SetConfiguration(thread, luax.Configuration{
		ActorContext: ctx,
		ActorPID:     ctx.Self(),
	})

	root := ctx.ActorSystem().Root
	self := ctx.Self()
	gw := a.GameWorld
	luax.CallAsync(ctx, a.luaRoot, thread, "on_create", a.StateMachine).Then(func(result interface{}) (interface{}, error) {
		var ms entity.MapSystem
		if gw != nil {
			ms = gw.GetMapSystem()
		}
		vals := luax.ResultValues(result)
		var maps []*entity.Map
		var err error
		if len(vals) == 0 || vals[0] == nil || vals[0] == lua.LNil {
			err = fmt.Errorf("on_create must return a non-empty map array")
		} else if tbl, ok := vals[0].(*lua.LTable); !ok {
			err = fmt.Errorf("on_create must return a table")
		} else {
			maps, err = entity.ParseCreateMaps(tbl, ms)
		}
		msg := &entity.ApplyStateMachineCreateMaps{}
		if err != nil {
			msg.Err = err.Error()
		} else {
			msg.Maps = maps
		}
		root.Send(self, msg)
		return nil, nil
	}).OnError(func(err error) {
		root.Send(self, &entity.ApplyStateMachineCreateMaps{Err: err.Error()})
	})
}

func (a *StateMachineActor) registerCreateMaps(ctx actor.Context, msg *entity.ApplyStateMachineCreateMaps) {
	if msg == nil {
		a.abortCreate(ctx, "nil create maps")
		return
	}
	if msg.Err != "" {
		a.abortCreate(ctx, msg.Err)
		return
	}
	if a.StateMachine == nil || a.StateMachine.Disposed() {
		a.abortCreate(ctx, "state machine disposed")
		return
	}
	if len(msg.Maps) == 0 {
		a.abortCreate(ctx, "on_create must return a non-empty map array")
		return
	}
	sm := a.StateMachine
	a.lc.attachPending = 0
	for _, m := range msg.Maps {
		already, err := sm.ReserveMap(m)
		if err != nil {
			a.abortCreate(ctx, err.Error())
			return
		}
		if already {
			continue
		}
		home := m.HomeActorPID()
		if home == nil {
			a.abortCreate(ctx, fmt.Sprintf("map %d has no home actor", m.GetMapID()))
			return
		}
		a.lc.attachPending++
		ctx.Send(home, &AttachStateMachine{
			StateMachine: sm,
			ReplyTo:      ctx.Self(),
			MapID:        m.GetMapID(),
		})
	}
	if a.lc.attachPending == 0 {
		a.finishCreate(ctx)
	}
}

func (a *StateMachineActor) abortCreate(ctx actor.Context, reason string) {
	name := ""
	if a.StateMachine != nil && a.StateMachine.Group != nil {
		name = a.StateMachine.Group.Name
	}
	if reason != "" {
		fmt.Printf("state machine %s create failed: %s\n", name, reason)
	}
	a.lc.deferred = nil
	if a.StateMachine != nil {
		a.StateMachine.AbortStart()
	}
	a.beginDetach(ctx)
}

func (a *StateMachineActor) finishCreate(ctx actor.Context) {
	if a.lc.attachDone {
		return
	}
	if a.StateMachine != nil && !a.StateMachine.HasRegisteredMaps() {
		a.abortCreate(ctx, "no maps registered")
		return
	}
	a.lc.attachDone = true
	deferred := a.lc.deferred
	a.lc.deferred = nil
	for _, msg := range deferred {
		switch m := msg.(type) {
		case *entity.EnterStateMachinePlayer:
			a.handleEnterPlayer(ctx, m)
		case *entity.StartStateMachine:
			a.handleStart(ctx)
		case *entity.CallStateMachineHook:
			if m != nil {
				a.callHook(ctx, m.Hook, m.Args...)
			}
		}
	}
}

func (a *StateMachineActor) handleEnterPlayer(ctx actor.Context, msg *entity.EnterStateMachinePlayer) {
	if msg == nil || msg.Character == nil || a.StateMachine == nil {
		return
	}
	a.StateMachine.Register(msg.Character)
	a.callHook(ctx, "on_player_enter", msg.Character)
}

func (a *StateMachineActor) handleStart(ctx actor.Context) {
	a.callHook(ctx, "on_start")
}

func (a *StateMachineActor) beginDetach(ctx actor.Context) {
	if a.lc.detaching {
		return
	}
	a.lc.detaching = true
	a.timeoutVersion++
	if a.timeoutCancel != nil {
		a.timeoutCancel()
		a.timeoutCancel = nil
	}
	a.cancelAllNamedSchedules()
	if a.StateMachine != nil {
		a.StateMachine.SetTimeoutDeadline(time.Time{})
	}
	root := ctx.ActorSystem().Root
	self := ctx.Self()
	sm := a.StateMachine
	a.lc.detachPending = 0
	for _, m := range sm.MapList() {
		if m == nil {
			continue
		}
		// A closed instance has stopped its actor, so it would never ack.
		if m.Closed() {
			continue
		}
		home := m.HomeActorPID()
		if home == nil {
			continue
		}
		a.lc.detachPending++
		root.Send(home, &DetachStateMachine{StateMachine: sm, ReplyTo: self, MapID: m.GetMapID()})
	}
	if a.lc.detachPending == 0 && a.lc.attachPending == 0 {
		a.finalizeStop()
	}
}

func (a *StateMachineActor) handleAttachAck(ctx actor.Context, msg *AttachStateMachineAck) {
	if msg == nil {
		return
	}
	if a.lc.attachPending > 0 {
		a.lc.attachPending--
	}

	if a.lc.detaching {
		if msg.OK && msg.Map.Closed() == false {
			if home := msg.Map.HomeActorPID(); home != nil {
				a.lc.detachPending++
				ctx.Send(home, &DetachStateMachine{StateMachine: a.StateMachine, ReplyTo: ctx.Self(), MapID: msg.MapID})
			}
		}
		if a.lc.detachPending == 0 && a.lc.attachPending == 0 {
			a.finalizeStop()
		}
		return
	}

	if msg.OK == false {
		a.abortCreate(ctx, fmt.Sprintf("attach to map %d failed", msg.MapID))
		return
	}
	if a.lc.attachPending == 0 {
		a.finishCreate(ctx)
	}
}

func (a *StateMachineActor) handleDetachAck(_ actor.Context) {
	if a.lc.detachPending > 0 {
		a.lc.detachPending--
	}
	if a.lc.detachPending == 0 && a.lc.attachPending == 0 {
		a.finalizeStop()
	}
}

func (a *StateMachineActor) finalizeStop() {
	sm := a.StateMachine
	if sm == nil {
		return
	}
	a.lc.stopped = true
	sm.ClearMaps()
	sm.ReleaseMaps()
	if sm.Group != nil {
		sm.Group.RemoveMachine(sm)
	}
	if a.GameWorld != nil {
		a.GameWorld.StopStateMachineActor(sm)
	}
}

func (a *StateMachineActor) callHook(ctx actor.Context, hook string, args ...interface{}) {
	if a.StateMachine == nil || a.StateMachine.Group == nil || hook == "" {
		return
	}
	if a.StateMachine.Disposed() {
		return
	}
	if a.luaRoot == nil {
		a.luaRoot = luax.NewState()
	}
	thread, err := luax.NewThread(a.luaRoot, a.StateMachine.Group.ScriptPath)
	if err != nil {
		return
	}
	if luax.HasFunc(thread, hook) == false {
		luax.Close(thread)
		return
	}
	luax.SetConfiguration(thread, luax.Configuration{
		ActorContext: ctx,
		ActorPID:     ctx.Self(),
	})
	luax.CallAsync(ctx, a.luaRoot, thread, hook, append([]interface{}{a.StateMachine}, args...)...).OnError(func(err error) {
		fmt.Printf("state machine %s hook %s: %v\n", a.StateMachine.Group.Name, hook, err)
	})
}
