package server

import (
	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/core/luax"
	g_actor "github.com/boyism80/fm/services/game/actor"
	lua "github.com/yuin/gopher-lua"
)

type luaMapCall struct {
	gs *GameServer
}

func (c luaMapCall) callerPID(L *lua.LState, actorCtx actor.Context) *actor.PID {
	if actorCtx != nil {
		return actorCtx.Self()
	}
	cfg, _ := luax.GetConfiguration(L)
	return cfg.MapActorPID
}

func (c luaMapCall) gameLogic(a actor.Actor) *g_actor.GameLogicActor {
	switch v := a.(type) {
	case *g_actor.MapActor:
		return &v.GameLogicActor
	case *g_actor.StateMachineActor:
		return &v.GameLogicActor
	}
	return nil
}

func (c luaMapCall) InvokeAwait(L *lua.LState, actorCtx actor.Context, targetPID *actor.PID, run g_actor.MapCallFunc) int {
	if c.gs == nil || L == nil || targetPID == nil || run == nil {
		return 0
	}
	callerPID := c.callerPID(L, actorCtx)
	if callerPID == nil {
		return 0
	}
	msg := &g_actor.MapCall{Run: run}
	if actorCtx != nil && callerPID.Equal(targetPID) {
		if logic := c.gameLogic(actorCtx.Actor()); logic != nil {
			handler := g_actor.MapCallHandler{}
			handler.Handle(actorCtx, logic, msg)
			for _, value := range msg.Values {
				L.Push(value)
			}
			return len(msg.Values)
		}
	}
	root := L.Parent
	if root == nil {
		return 0
	}
	rootContext := c.gs.GetRootContext()
	if rootContext == nil {
		return 0
	}
	msg.ReplyTo = callerPID
	msg.Root = root
	msg.Thread = L
	rootContext.Send(targetPID, msg)
	return L.Yield(lua.LNil)
}

func (c luaMapCall) InvokeAwaitAsync(L *lua.LState, actorCtx actor.Context, targetPID *actor.PID, run g_actor.MapCallAsyncFunc) int {
	if c.gs == nil || L == nil || targetPID == nil || run == nil {
		return 0
	}
	callerPID := c.callerPID(L, actorCtx)
	root := L.Parent
	if callerPID == nil || root == nil {
		return 0
	}
	rootContext := c.gs.GetRootContext()
	if rootContext == nil {
		return 0
	}
	msg := &g_actor.MapCallAsync{
		Run:     run,
		ReplyTo: callerPID,
		Root:    root,
		Thread:  L,
	}
	if actorCtx != nil && callerPID.Equal(targetPID) {
		if logic := c.gameLogic(actorCtx.Actor()); logic != nil {
			handler := g_actor.MapCallAsyncHandler{}
			handler.Handle(actorCtx, logic, msg)
			return L.Yield(lua.LNil)
		}
	}
	rootContext.Send(targetPID, msg)
	return L.Yield(lua.LNil)
}
