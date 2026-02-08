package server

import (
	"fmt"

	"github.com/boyism80/fm/core/luax"
	lua "github.com/yuin/gopher-lua"
)

// ExecuteScript runs the given function on a lua thread that already has the script loaded (e.g. from luax.NewThread(rootState, scriptPath)).
// Does not create any LState; only gets func from thread and resumes.
func (gs *GameServer) ExecuteScript(rootState *lua.LState, luaThread *lua.LState, funcName string, args ...luax.Luable) error {
	f := luaThread.GetGlobal(funcName)
	if f.Type() != lua.LTFunction {
		return fmt.Errorf("function %s not found in script", funcName)
	}

	lvArgs := make([]lua.LValue, len(args))
	for i, arg := range args {
		lvArgs[i] = luax.NewLuable(luaThread, arg)
	}
	_, resumeErr, _ := rootState.Resume(luaThread, f.(*lua.LFunction), lvArgs...)
	if resumeErr != nil {
		return fmt.Errorf("failed to call %s: %w", funcName, resumeErr)
	}
	return nil
}
