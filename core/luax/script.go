package luax

import (
	"fmt"

	lua "github.com/yuin/gopher-lua"
)

func Close(thread *lua.LState) {
	if thread == nil {
		return
	}
	ClearConfiguration(thread)
	thread.Close()
}

func shouldAutoClose(thread *lua.LState) bool {
	cfg, ok := GetConfiguration(thread)
	if !ok {
		return true
	}
	return !cfg.KeepAlive
}

func Call(thread *lua.LState, funcName string, args ...interface{}) (lua.LValue, error) {
	if thread == nil {
		return nil, fmt.Errorf("nil lua thread")
	}
	f := thread.GetGlobal(funcName)
	if f.Type() != lua.LTFunction {
		return nil, nil
	}
	lvArgs, err := toLValues(thread, args)
	if err != nil {
		return nil, err
	}
	thread.Push(f)
	for _, lv := range lvArgs {
		thread.Push(lv)
	}
	if err := thread.PCall(len(lvArgs), 1, nil); err != nil {
		if shouldAutoClose(thread) {
			Close(thread)
		}
		return nil, fmt.Errorf("%s: %w", funcName, err)
	}
	ret := thread.Get(-1)
	thread.Pop(1)
	if shouldAutoClose(thread) {
		Close(thread)
	}
	return ret, nil
}

func Execute(root *lua.LState, thread *lua.LState, funcName string, args ...interface{}) (lua.ResumeState, lua.LValue, error) {
	if root == nil || thread == nil {
		return lua.ResumeOK, nil, fmt.Errorf("nil lua root/thread")
	}
	f := thread.GetGlobal(funcName)
	if f.Type() != lua.LTFunction {
		return lua.ResumeOK, nil, nil
	}
	lvArgs, err := toLValues(thread, args)
	if err != nil {
		return lua.ResumeOK, nil, err
	}
	thread.Push(f)
	state, resumeErr, values := root.Resume(thread, f.(*lua.LFunction), lvArgs...)
	if resumeErr != nil {
		if shouldAutoClose(thread) {
			Close(thread)
		}
		return lua.ResumeOK, nil, fmt.Errorf("%s: %w", funcName, resumeErr)
	}
	var ret lua.LValue
	if len(values) > 0 {
		ret = values[0]
	}
	if state != lua.ResumeYield && shouldAutoClose(thread) {
		Close(thread)
	}
	return state, ret, nil
}

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

func Resume(root *lua.LState, thread *lua.LState, args ...lua.LValue) (lua.ResumeState, error) {
	if root == nil || thread == nil {
		return lua.ResumeOK, fmt.Errorf("nil lua root/thread")
	}
	state, resumeErr, _ := root.Resume(thread, nil, args...)
	if resumeErr != nil {
		if shouldAutoClose(thread) {
			Close(thread)
		}
		return lua.ResumeOK, resumeErr
	}
	if state != lua.ResumeYield && shouldAutoClose(thread) {
		Close(thread)
	}
	return state, nil
}

func InlineCall(root *lua.LState, scriptPath string, funcName string, args ...interface{}) (lua.LValue, error) {
	thread, err := NewThread(root, scriptPath)
	if err != nil {
		return nil, err
	}
	return Call(thread, funcName, args...)
}
