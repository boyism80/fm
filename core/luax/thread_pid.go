package luax

import (
	"sync"

	"github.com/asynkron/protoactor-go/actor"
	lua "github.com/yuin/gopher-lua"
)

type Configuration struct {
	ActorContext actor.Context
	KeepAlive    bool
}

var threadConfig = struct {
	sync.RWMutex
	m map[*lua.LState]Configuration
}{
	m: make(map[*lua.LState]Configuration),
}

func SetConfiguration(L *lua.LState, cfg Configuration) {
	if L == nil {
		return
	}
	threadConfig.Lock()
	defer threadConfig.Unlock()
	threadConfig.m[L] = cfg
}

func GetConfiguration(L *lua.LState) (Configuration, bool) {
	if L == nil {
		return Configuration{}, false
	}
	threadConfig.RLock()
	defer threadConfig.RUnlock()
	cfg, ok := threadConfig.m[L]
	return cfg, ok
}

func ClearConfiguration(L *lua.LState) {
	if L == nil {
		return
	}
	threadConfig.Lock()
	defer threadConfig.Unlock()
	delete(threadConfig.m, L)
}

func SetThreadPID(L *lua.LState, pid *actor.PID) {

	if pid == nil {
		return
	}
}

func GetThreadPID(L *lua.LState) *actor.PID {
	cfg, ok := GetConfiguration(L)
	if !ok {
		return nil
	}
	if cfg.ActorContext == nil {
		return nil
	}
	return cfg.ActorContext.Self()
}

func SetThreadActorContext(L *lua.LState, ctx actor.Context) {
	cfg, _ := GetConfiguration(L)
	cfg.ActorContext = ctx
	SetConfiguration(L, cfg)
}

func GetThreadActorContext(L *lua.LState) actor.Context {
	cfg, ok := GetConfiguration(L)
	if !ok {
		return nil
	}
	return cfg.ActorContext
}

func ClearThreadPID(L *lua.LState) {

	ClearConfiguration(L)
}
