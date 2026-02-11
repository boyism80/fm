package luax

import (
	"fmt"

	"github.com/asynkron/protoactor-go/actor"
	lua "github.com/yuin/gopher-lua"
)

// ExecuteScript runs the given function on a lua thread. mapActorPID is required (nil returns error).
// Sets thread PID before first resume; clears it when script ends (ResumeOK or error), not on yield.
func ExecuteScript(root *lua.LState, thread *lua.LState, pid *actor.PID, funcName string, args ...interface{}) (lua.ResumeState, error) {
	if pid == nil {
		return lua.ResumeOK, fmt.Errorf("script requires map actor PID")
	}
	f := thread.GetGlobal(funcName)
	if f.Type() != lua.LTFunction {
		return lua.ResumeOK, fmt.Errorf("function %s not found in script", funcName)
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
		case Luable:
			lvArgs[i] = NewLuable(thread, v)
		case lua.LValue:
			lvArgs[i] = v
		default:
			return lua.ResumeOK, fmt.Errorf("unsupported argument type: %T", arg)
		}
	}

	SetThreadPID(thread, pid)
	state, resumeErr, _ := root.Resume(thread, f.(*lua.LFunction), lvArgs...)
	if state == lua.ResumeOK || resumeErr != nil {
		ClearThreadPID(thread)
	}
	if resumeErr != nil {
		return lua.ResumeOK, fmt.Errorf("failed to call %s: %w", funcName, resumeErr)
	}
	return state, nil
}
