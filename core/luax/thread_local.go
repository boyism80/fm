package luax

import (
	"bytes"
	"runtime"
	"strconv"
	"sync"

	lua "github.com/yuin/gopher-lua"
)

var (
	threadLocalStates = sync.Map{}
	initHooks         []func(*lua.LState)
	initHooksMu       sync.Mutex
)

func getGoroutineID() uint64 {
	var buf [64]byte
	n := runtime.Stack(buf[:], false)

	stackStr := string(buf[:n])
	prefix := "goroutine "
	idx := bytes.Index([]byte(stackStr), []byte(prefix))
	if idx == -1 {
		return 0
	}

	start := idx + len(prefix)

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

func RegisterThreadLocalInitHook(fn func(*lua.LState)) {
	initHooksMu.Lock()
	defer initHooksMu.Unlock()
	initHooks = append(initHooks, fn)
}

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
