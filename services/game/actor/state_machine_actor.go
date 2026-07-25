package actor

import (
	"fmt"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/scheduler"
	"github.com/boyism80/fm/core/luax"
	"github.com/boyism80/fm/services/game/entity"
	lua "github.com/yuin/gopher-lua"
)

type StateMachineActor struct {
	GameLogicActor
	StateMachine   *entity.StateMachine
	luaRoot        *lua.LState
	timeoutVersion uint64
	timeoutCancel  scheduler.CancelFunc
	pendingAttach  int
	attachFailed   bool
	attachComplete bool
	pendingMsgs    []interface{}
	pendingDetach  int
	detaching      bool
}

func NewStateMachineActor(sm *entity.StateMachine, gameWorld entity.GameWorld) *StateMachineActor {
	a := &StateMachineActor{StateMachine: sm}
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
	switch msg := ctx.Message().(type) {
	case *entity.BootstrapStateMachine:
		a.beginAttach(ctx)
	case *AttachStateMachineAck:
		a.handleAttachAck(ctx, msg)
	case *entity.FinishStateMachineCreate:
		a.finishCreate(ctx)
	case *entity.EnterStateMachinePlayer:
		if a.attachComplete {
			a.handleEnterPlayer(ctx, msg)
		} else {
			a.pendingMsgs = append(a.pendingMsgs, msg)
		}
	case *entity.StartStateMachine:
		if a.attachComplete {
			a.handleStart(ctx)
		} else {
			a.pendingMsgs = append(a.pendingMsgs, msg)
		}
	case *entity.CallStateMachineHook:
		if a.attachComplete {
			a.callHook(ctx, msg.Hook, msg.Args...)
		} else {
			a.pendingMsgs = append(a.pendingMsgs, msg)
		}
	case *entity.LeaveStateMachinePlayer:
		if msg != nil && a.StateMachine != nil {
			a.StateMachine.LeavePlayer(ctx, msg.Character, msg.WarpLeaver)
		}
	case *entity.ScheduleStateMachineTimeout:
		if a.scheduler != nil {
			a.timeoutVersion++
			if a.timeoutCancel != nil {
				a.timeoutCancel()
			}
			deadline := time.Now().Add(time.Duration(msg.Milliseconds) * time.Millisecond)
			a.StateMachine.SetTimeoutDeadline(deadline)
			a.timeoutCancel = a.scheduler.SendOnce(time.Duration(msg.Milliseconds)*time.Millisecond, ctx.Self(), &entity.StateMachineTimeout{
				Version: a.timeoutVersion,
			})
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
			a.callHook(ctx, "on_scheduled_timeout", a.StateMachine)
		}
	case *entity.StopStateMachine:
		a.beginDetach(ctx)
	case *DetachStateMachineAck:
		a.handleDetachAck(ctx)
	case *actor.Stopped:
		a.StopTimers()
		if a.timeoutCancel != nil {
			a.timeoutCancel()
			a.timeoutCancel = nil
		}
		if a.luaRoot != nil {
			a.luaRoot.Close()
			a.luaRoot = nil
		}
		a.GameLogicActor.Receive(ctx)
	default:
		a.GameLogicActor.Receive(ctx)
	}
}

func (a *StateMachineActor) beginAttach(ctx actor.Context) {
	sm := a.StateMachine
	if sm == nil || sm.Group == nil || a.GameWorld == nil {
		return
	}
	root := ctx.ActorSystem().Root
	self := ctx.Self()
	a.pendingAttach = 0
	a.attachFailed = false
	for _, id := range sm.Group.DeclaredMaps() {
		m := a.GameWorld.GetMapSystem().Get(id)
		if m == nil {
			a.attachFailed = true
			continue
		}
		home := m.HomeActorPID()
		if home == nil {
			a.attachFailed = true
			continue
		}
		a.pendingAttach++
		root.Send(home, &AttachStateMachine{StateMachine: sm, ReplyTo: self, MapID: id})
	}
	if a.pendingAttach == 0 {
		a.finishAttach(ctx)
	}
}

func (a *StateMachineActor) handleAttachAck(ctx actor.Context, msg *AttachStateMachineAck) {
	if msg == nil {
		return
	}
	if msg.OK {
		m := a.GameWorld.GetMapSystem().Get(msg.MapID)
		a.StateMachine.RecordMap(msg.MapID, m)
	} else {
		a.attachFailed = true
	}
	if a.pendingAttach > 0 {
		a.pendingAttach--
	}
	if a.pendingAttach == 0 {
		a.finishAttach(ctx)
	}
}

func (a *StateMachineActor) finishAttach(ctx actor.Context) {
	if a.attachFailed {
		a.pendingMsgs = nil
		a.StateMachine.AbortStart()
		a.beginDetach(ctx)
		return
	}
	if a.StateMachine == nil || a.StateMachine.Group == nil {
		a.finishCreate(ctx)
		return
	}
	if a.luaRoot == nil {
		a.luaRoot = luax.NewState()
	}
	thread, err := luax.NewThread(a.luaRoot, a.StateMachine.Group.ScriptPath)
	if err != nil {
		a.finishCreate(ctx)
		return
	}
	if !luax.HasFunc(thread, "on_create") {
		luax.Close(thread)
		a.finishCreate(ctx)
		return
	}
	luax.Close(thread)
	a.callHook(ctx, "on_create", a.StateMachine)
}

func (a *StateMachineActor) finishCreate(ctx actor.Context) {
	if a.attachComplete {
		return
	}
	a.attachComplete = true
	pending := a.pendingMsgs
	a.pendingMsgs = nil
	for _, msg := range pending {
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
	a.callHook(ctx, "on_player_enter", a.StateMachine, msg.Character)
}

func (a *StateMachineActor) handleStart(ctx actor.Context) {
	a.callHook(ctx, "on_start", a.StateMachine)
}

func (a *StateMachineActor) beginDetach(ctx actor.Context) {
	if a.detaching {
		return
	}
	a.detaching = true
	a.timeoutVersion++
	if a.timeoutCancel != nil {
		a.timeoutCancel()
		a.timeoutCancel = nil
	}
	if a.StateMachine != nil {
		a.StateMachine.SetTimeoutDeadline(time.Time{})
	}
	root := ctx.ActorSystem().Root
	self := ctx.Self()
	sm := a.StateMachine
	a.pendingDetach = 0
	for _, m := range sm.MapList() {
		if m == nil {
			continue
		}
		home := m.HomeActorPID()
		if home == nil {
			continue
		}
		a.pendingDetach++
		root.Send(home, &DetachStateMachine{StateMachine: sm, ReplyTo: self, MapID: m.GetMapID()})
	}
	if a.pendingDetach == 0 {
		a.finalizeStop()
	}
}

func (a *StateMachineActor) handleDetachAck(_ actor.Context) {
	if a.pendingDetach > 0 {
		a.pendingDetach--
	}
	if a.pendingDetach == 0 {
		a.finalizeStop()
	}
}

func (a *StateMachineActor) finalizeStop() {
	sm := a.StateMachine
	if sm == nil {
		return
	}
	sm.ClearMaps()
	if sm.Group != nil {
		sm.Group.RemoveMachine(sm.ID)
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
		if hook == "on_create" {
			a.finishCreate(ctx)
		}
		return
	}
	if !luax.HasFunc(thread, hook) {
		luax.Close(thread)
		if hook == "on_create" {
			a.finishCreate(ctx)
		}
		return
	}
	luax.SetConfiguration(thread, luax.Configuration{
		ActorContext: ctx,
		ActorPID:     ctx.Self(),
	})
	root := ctx.ActorSystem().Root
	self := ctx.Self()
	luax.CallAsync(a.luaRoot, thread, hook, args...).Then(func(interface{}) (interface{}, error) {
		if hook == "on_create" && root != nil {
			root.Send(self, &entity.FinishStateMachineCreate{})
		}
		return nil, nil
	}).OnError(func(err error) {
		fmt.Printf("state machine %s hook %s: %v\n", a.StateMachine.Group.Name, hook, err)
		if hook == "on_create" && root != nil {
			root.Send(self, &entity.FinishStateMachineCreate{})
		}
	})
}
