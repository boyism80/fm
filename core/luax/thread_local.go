package luax

import (
	"bytes"
	"runtime"
	"strconv"
	"sync"

	lua "github.com/yuin/gopher-lua"
)

var (
	threadLocalStates = sync.Map{} // map[uint64]*lua.LState
	initHooks         []func(*lua.LState)
	initHooksMu       sync.Mutex
)

// getGoroutineID extracts goroutine ID from the current goroutine's stack trace
// This is more reliable than unsafe methods and works across Go versions
func getGoroutineID() uint64 {
	var buf [64]byte
	n := runtime.Stack(buf[:], false)

	// Stack trace format: "goroutine <id> [running]:"
	// Find "goroutine " prefix and extract the ID
	stackStr := string(buf[:n])
	prefix := "goroutine "
	idx := bytes.Index([]byte(stackStr), []byte(prefix))
	if idx == -1 {
		return 0
	}

	// Skip "goroutine " prefix
	start := idx + len(prefix)

	// Extract number until space or bracket
	var idStr []byte
	for i := start; i < len(stackStr); i++ {
		if stackStr[i] >= '0' && stackStr[i] <= '9' {
			idStr = append(idStr, stackStr[i])
		} else {
			break
		}
	}

	if len(idStr) == 0 {
		return 0
	}

	id, err := strconv.ParseUint(string(idStr), 10, 64)
	if err != nil {
		return 0
	}

	return id
}

// RegisterThreadLocalInitHook registers a hook that will be called when a new thread-local LuaState is created
func RegisterThreadLocalInitHook(fn func(*lua.LState)) {
	initHooksMu.Lock()
	defer initHooksMu.Unlock()
	initHooks = append(initHooks, fn)
}

// CleanupThreadLocalState removes the LuaState for the current goroutine
// This should be called when the goroutine is about to end
func CleanupThreadLocalState() {
	goid := getGoroutineID()
	if goid == 0 {
		return
	}
	if state, ok := threadLocalStates.LoadAndDelete(goid); ok {
		if luaState, ok := state.(*lua.LState); ok {
			luaState.Close()
		}
	}
}
