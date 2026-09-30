package entity

import (
	"github.com/boyism80/fm/core/async"
	"github.com/boyism80/fm/core/luax"
	lua "github.com/yuin/gopher-lua"
)

func luaYieldPromise(L *lua.LState, gw GameWorld, promise *async.Promise, resume func(result interface{}, err error) []lua.LValue) int {
	if L == nil || promise == nil || gw == nil {
		return 0
	}
	cfg, ok := luax.GetConfiguration(L)
	if !ok || cfg.ActorContext == nil {
		return 0
	}
	pid := cfg.ActorContext.Self()
	root := L.Parent
	if pid == nil || root == nil {
		return 0
	}
	thread := L
	finish := func(result interface{}, err error) {
		var args []lua.LValue
		if resume != nil {
			args = resume(result, err)
		}
		gw.ResumeLua(pid, root, thread, args)
	}
	promise.Then(func(result interface{}) (interface{}, error) {
		finish(result, nil)
		return nil, nil
	}).OnError(func(err error) {
		finish(nil, err)
	})
	return L.Yield(lua.LNil)
}
