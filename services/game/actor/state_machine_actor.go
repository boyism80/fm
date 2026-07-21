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
	pendingHooks   []*entity.CallStateMachineHook
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
	case *entity.EnterStateMachinePlayers:
		for _, ch := range a.StateMachine.Players() {
			a.callHook(ctx, "on_player_entry", a.StateMachine, ch)
		}
	case *entity.CallStateMachineHook:
		if a.attachComplete {
			a.callHook(ctx, msg.Hook, msg.Args...)
		} else {
			a.pendingHooks = append(a.pendingHooks, msg)
		}
	case *entity.ScheduleStateMachineTimeout:
		if a.scheduler != nil {
			a.timeoutVersion++
			if a.timeoutCancel != nil {
				a.timeoutCancel()
			}
			a.timeoutCancel = a.scheduler.SendOnce(time.Duration(msg.Milliseconds)*time.Millisecond, ctx.Self(), &entity.StateMachineTimeout{
				Version: a.timeoutVersion,
			})
		}
	case *entity.StateMachineTimeout:
		if msg.Version == a.timeoutVersion {
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
		a.pendingHooks = nil
		a.StateMachine.AbortStart()
		a.beginDetach(ctx)
		return
	}
	a.attachComplete = true
	a.callHook(ctx, "on_setup", a.StateMachine)
	for _, msg := range a.pendingHooks {
		if msg != nil {
			a.callHook(ctx, msg.Hook, msg.Args...)
		}
	}
	a.pendingHooks = nil
}

func (a *StateMachineActor) beginDetach(ctx actor.Context) {
	if a.detaching {
		return
	}
	a.detaching = true
	a.timeoutVersion++
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
		return
	}
	if thread.GetGlobal(hook).Type() != lua.LTFunction {
		luax.Close(thread)
		return
	}
	luax.SetConfiguration(thread, luax.Configuration{
		ActorContext: ctx,
		ActorPID:     ctx.Self(),
	})
	root := ctx.ActorSystem().Root
	self := ctx.Self()
	luax.CallAsync(a.luaRoot, thread, hook, args...).Then(func(interface{}) (interface{}, error) {
		if hook == "on_setup" && root != nil {
			root.Send(self, &entity.EnterStateMachinePlayers{})
		}
		return nil, nil
	}).OnError(func(err error) {
		fmt.Printf("state machine %s hook %s: %v\n", a.StateMachine.Group.Name, hook, err)
	})
}
