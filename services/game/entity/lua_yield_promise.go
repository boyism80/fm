package entity

import (
	"github.com/boyism80/fm/core/async"
	"github.com/boyism80/fm/core/luax"
	lua "github.com/yuin/gopher-lua"
)

func luaYieldPromise(L *lua.LState, gw GameWorld, promise *async.Promise) int {
	if L == nil || promise == nil {
		return 0
	}
	cfg, ok := luax.GetConfiguration(L)
	if !ok || cfg.ActorContext == nil || gw == nil {
		return 0
	}
	pid := cfg.ActorContext.Self()
	root := L.Parent
	if pid == nil || root == nil {
		return 0
	}
	thread := L
	promise.Then(func(interface{}) (interface{}, error) {
		gw.ResumeLua(pid, root, thread, nil)
		return nil, nil
	}).OnError(func(error) {
		gw.ResumeLua(pid, root, thread, nil)
	})
	return L.Yield(lua.LNil)
}
