package entity

import (
	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/core/async"
	"github.com/boyism80/fm/core/luax"
	lua "github.com/yuin/gopher-lua"
)

func LuaYieldPromise[T any](L *lua.LState, gw GameWorld, promise *async.Promise[T], resume func(result T, err error) []lua.LValue) int {
	if promise == nil {
		return 0
	}
	var result T
	task := promise.Do(func(value T) error {
		result = value
		return nil
	}).Task()
	return LuaYieldTask(L, gw, task, func(err error) []lua.LValue {
		if resume == nil {
			return nil
		}
		return resume(result, err)
	})
}

func LuaYieldTask(L *lua.LState, gw GameWorld, task *async.Task, resume func(err error) []lua.LValue) int {
	if L == nil || task == nil || gw == nil {
		return 0
	}
	cfg, ok := luax.GetConfiguration(L)
	if !ok {
		return 0
	}
	var pid *actor.PID
	if cfg.ActorContext != nil {
		pid = cfg.ActorContext.Self()
	}
	if pid == nil {
		pid = cfg.ActorPID
	}
	root := L.Parent
	if pid == nil || root == nil {
		return 0
	}
	thread := L
	finish := func(err error) {
		var args []lua.LValue
		if resume != nil {
			args = resume(err)
		}
		gw.ResumeLua(pid, root, thread, args)
	}
	task.Do(func() error {
		finish(nil)
		return nil
	}).OnError(finish)
	return L.Yield(lua.LNil)
}
