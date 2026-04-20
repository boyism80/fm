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

