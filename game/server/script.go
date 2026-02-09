package server

import (
	"fmt"

	"github.com/boyism80/fm/core/luax"
	lua "github.com/yuin/gopher-lua"
)

// ExecuteScript runs the given function on a lua thread that already has the script loaded (e.g. from luax.NewThread(root, scriptPath)).
// Does not create any LState; only gets func from thread and resumes.
func (gs *GameServer) ExecuteScript(root *lua.LState, thread *lua.LState, funcName string, args ...interface{}) error {
	f := thread.GetGlobal(funcName)
	if f.Type() != lua.LTFunction {
		return fmt.Errorf("function %s not found in script", funcName)
	}

	lvArgs := make([]lua.LValue, len(args))
	for i, arg := range args {
		switch v := arg.(type) {
		case uint32:
			lvArgs[i] = lua.LNumber(v)
		case int32:
			lvArgs[i] = lua.LNumber(v)
		case uint64:
			lvArgs[i] = lua.LNumber(v)
		case int64:
			lvArgs[i] = lua.LNumber(v)
		case float64:
			lvArgs[i] = lua.LNumber(v)
		case string:
			lvArgs[i] = lua.LString(v)
		case bool:
			lvArgs[i] = lua.LBool(v)
		case luax.Luable:
			lvArgs[i] = luax.NewLuable(thread, v)
		case lua.LValue:
			lvArgs[i] = v
		default:
			return fmt.Errorf("unsupported argument type: %T", arg)
		}
	}

	_, resumeErr, _ := root.Resume(thread, f.(*lua.LFunction), lvArgs...)
	if resumeErr != nil {
		return fmt.Errorf("failed to call %s: %w", funcName, resumeErr)
	}
	return nil
}
