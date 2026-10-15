package game

import (
	"github.com/asynkron/protoactor-go/actor"
	c_actor "github.com/boyism80/fm/core/actor"
	"github.com/boyism80/fm/core/luax"
	lua "github.com/yuin/gopher-lua"
)

type luaMapCall struct {
	gw GameWorld
}

func (c luaMapCall) callerPID(L *lua.LState, actorCtx actor.Context) *actor.PID {
	if actorCtx != nil {
		return actorCtx.Self()
	}
	cfg, _ := luax.GetConfiguration(L)
	return cfg.ActorPID
}

func (c luaMapCall) InvokeAwait(L *lua.LState, actorCtx actor.Context, targetPID *actor.PID, run MapCallFunc) int {
	if c.gw == nil || L == nil || targetPID == nil || run == nil {
		return 0
	}
	callerPID := c.callerPID(L, actorCtx)
	if callerPID == nil {
		return 0
	}
	msg := &MapCall{Run: run}
	if actorCtx != nil && callerPID.Equal(targetPID) {
		if runner, ok := actorCtx.Actor().(MapCallRunner); ok {
			runner.RunMapCall(actorCtx, msg)
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
	msg.ReplyTo = callerPID
	msg.Root = root
	msg.Thread = L
	c.gw.Systems().Get[*c_actor.ActorRegistry]().Root().Send(targetPID, msg)
	return L.Yield(lua.LNil)
}

func (c luaMapCall) InvokeAwaitAsync(L *lua.LState, actorCtx actor.Context, targetPID *actor.PID, run MapCallAsyncFunc) int {
	if c.gw == nil || L == nil || targetPID == nil || run == nil {
		return 0
	}
	callerPID := c.callerPID(L, actorCtx)
	root := L.Parent
	if callerPID == nil || root == nil {
		return 0
	}
	msg := &MapCallAsync{
		Run:     run,
		ReplyTo: callerPID,
		Root:    root,
		Thread:  L,
	}
	if actorCtx != nil && callerPID.Equal(targetPID) {
		if runner, ok := actorCtx.Actor().(MapCallRunner); ok {
			runner.RunMapCallAsync(actorCtx, msg)
			return L.Yield(lua.LNil)
		}
	}
	c.gw.Systems().Get[*c_actor.ActorRegistry]().Root().Send(targetPID, msg)
	return L.Yield(lua.LNil)
}
