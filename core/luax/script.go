package luax

import (
	"fmt"

	"github.com/asynkron/protoactor-go/actor"
	lua "github.com/yuin/gopher-lua"
)

func Call(root *lua.LState, scriptPath string, funcName string, args ...interface{}) (lua.LValue, *lua.LState, error) {
	thread, err := NewThread(root, scriptPath)
	if err != nil {
		return nil, nil, err
	}
	f := thread.GetGlobal(funcName)
	if f.Type() != lua.LTFunction {
		return nil, thread, nil
	}
	lvArgs, err := toLValues(thread, args)
	if err != nil {
		return nil, thread, err
	}
	thread.Push(f)
	for _, lv := range lvArgs {
		thread.Push(lv)
	}
	if err := thread.PCall(len(lvArgs), 1, nil); err != nil {
		return nil, thread, fmt.Errorf("%s: %w", funcName, err)
	}
	defer thread.Pop(1)
	return thread.Get(-1), thread, nil
}

// CallFunction invokes a Lua function value on the root state with the given arguments.
// Used when the timer callback is a Lua function (e.g. from mktimer(..., function() ... end)).
func CallFunction(root *lua.LState, fn *lua.LFunction, args ...interface{}) (lua.LValue, error) {
	lvArgs, err := toLValues(root, args)
	if err != nil {
		return nil, err
	}
	root.Push(fn)
	for _, lv := range lvArgs {
		root.Push(lv)
	}
	err = root.PCall(len(lvArgs), 1, nil)
	if err != nil {
		root.Pop(1)
		return nil, err
	}
	ret := root.Get(-1)
	root.Pop(1)
	return ret, nil
}

func toLValues(L *lua.LState, args []interface{}) ([]lua.LValue, error) {
	out := make([]lua.LValue, len(args))
	for i, arg := range args {
		if arg == nil {
			out[i] = lua.LNil
			continue
		}
		switch v := arg.(type) {
		case uint32:
			out[i] = lua.LNumber(v)
		case int32:
			out[i] = lua.LNumber(v)
		case uint64:
			out[i] = lua.LNumber(v)
		case int64:
			out[i] = lua.LNumber(v)
		case float64:
			out[i] = lua.LNumber(v)
		case string:
			out[i] = lua.LString(v)
		case bool:
			out[i] = lua.LBool(v)
		case Luable:
			if v == nil {
				out[i] = lua.LNil
			} else {
				out[i] = NewLuable(L, v)
			}
		case []Luable:
			tbl := L.NewTable()
			for idx, item := range v {
				if item == nil {
					continue
				}
				tbl.RawSetInt(idx+1, NewLuable(L, item))
			}
			out[i] = tbl
		case lua.LValue:
			out[i] = v
		default:
			if l, ok := arg.(Luable); ok {
				out[i] = NewLuable(L, l)
			} else {
				return nil, fmt.Errorf("unsupported argument type: %T", arg)
			}
		}
	}
	return out, nil
}

// Execute runs the given function on a lua thread. mapActorPID is required (nil returns error).
// Sets thread PID before first resume; clears it when script ends (ResumeOK or error), not on yield.
func Execute(root *lua.LState, thread *lua.LState, pid *actor.PID, funcName string, args ...interface{}) (lua.ResumeState, error) {
	if pid == nil {
		return lua.ResumeOK, fmt.Errorf("script requires map actor PID")
	}
	f := thread.GetGlobal(funcName)
	if f.Type() != lua.LTFunction {
		return lua.ResumeOK, fmt.Errorf("function %s not found in script", funcName)
	}
	lvArgs, err := toLValues(thread, args)
	if err != nil {
		return lua.ResumeOK, err
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
