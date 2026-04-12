package luax

import (
	"sync"

	"github.com/asynkron/protoactor-go/actor"
	lua "github.com/yuin/gopher-lua"
)

var threadPIDs = struct {
	sync.RWMutex
	m map[*lua.LState]*actor.PID
}{
	m: make(map[*lua.LState]*actor.PID),
}

func SetThreadPID(L *lua.LState, pid *actor.PID) {
	if L == nil {
		return
	}
	threadPIDs.Lock()
	defer threadPIDs.Unlock()
	if pid == nil {
		delete(threadPIDs.m, L)
		return
	}
	threadPIDs.m[L] = pid
}

func GetThreadPID(L *lua.LState) *actor.PID {
	if L == nil {
		return nil
	}
	threadPIDs.RLock()
	defer threadPIDs.RUnlock()
	return threadPIDs.m[L]
}

func ClearThreadPID(L *lua.LState) {
	if L == nil {
		return
	}
	threadPIDs.Lock()
	defer threadPIDs.Unlock()
	delete(threadPIDs.m, L)
}
