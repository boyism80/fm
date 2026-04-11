package luax

import (
	"sync"

	"github.com/asynkron/protoactor-go/actor"
	lua "github.com/yuin/gopher-lua"
)

// threadPIDs maps each Lua thread (coroutine) LState to the actor PID that runs it.
// Keys are *lua.LState; only the thread used for script execution is registered.
var threadPIDs = struct {
	sync.RWMutex
	m map[*lua.LState]*actor.PID
}{
	m: make(map[*lua.LState]*actor.PID),
}

// SetThreadPID associates the given Lua thread with an actor PID.
// Call before resuming the thread (e.g. right before ExecuteScript runs).
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

// GetThreadPID returns the actor PID for the given Lua thread, or nil.
// Safe to call from any builtin (e.g. sleep) that receives the current thread L.
func GetThreadPID(L *lua.LState) *actor.PID {
	if L == nil {
		return nil
	}
	threadPIDs.RLock()
	defer threadPIDs.RUnlock()
	return threadPIDs.m[L]
}

// ClearThreadPID removes the association for the given thread.
// Call when the script ends (ResumeOK or error).
func ClearThreadPID(L *lua.LState) {
	if L == nil {
		return
	}
	threadPIDs.Lock()
	defer threadPIDs.Unlock()
	delete(threadPIDs.m, L)
}
