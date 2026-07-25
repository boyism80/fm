package luax

import (
	"fmt"
	"sync"

	"github.com/boyism80/fm/core/async"
	lua "github.com/yuin/gopher-lua"
)

type threadScript struct {
	path string
	mod  *lua.LTable
}

var threadScripts = struct {
	sync.RWMutex
	m map[*lua.LState]threadScript
}{
	m: make(map[*lua.LState]threadScript),
}

func setThreadScript(thread *lua.LState, path string, mod *lua.LTable) {
	if thread == nil {
		return
	}
	threadScripts.Lock()
	defer threadScripts.Unlock()
	threadScripts.m[thread] = threadScript{path: path, mod: mod}
}

func clearThreadScript(thread *lua.LState) {
	if thread == nil {
		return
	}
	threadScripts.Lock()
	defer threadScripts.Unlock()
	delete(threadScripts.m, thread)
}

func threadFunc(thread *lua.LState, name string) *lua.LFunction {
	if thread == nil || name == "" {
		return nil
	}
	threadScripts.RLock()
	meta, ok := threadScripts.m[thread]
	threadScripts.RUnlock()
	if !ok {
		return nil
	}
	return ModuleFunc(meta.mod, name)
}

func HasFunc(thread *lua.LState, name string) bool {
	return threadFunc(thread, name) != nil
}

func Close(thread *lua.LState) {
	if thread == nil {
		return
	}
	ClearConfiguration(thread)
	clearThreadScript(thread)
	thread.Close()
}

func shouldAutoClose(thread *lua.LState) bool {
	cfg, ok := GetConfiguration(thread)
	if !ok {
		return true
	}
	return !cfg.KeepAlive
}

func ResultValues(value interface{}) []lua.LValue {
	vals, _ := value.([]lua.LValue)
	return vals
}

func completeCall(thread *lua.LState, values []lua.LValue, err error) {
	cfg, ok := GetConfiguration(thread)
	if !ok || cfg.CallPromise == nil {
		return
	}
	promise := cfg.CallPromise
	cfg.CallPromise = nil
	SetConfiguration(thread, cfg)
	if err != nil {
		promise.SetError(err)
		return
	}
	promise.SetResult(values)
}

func NewThread(root *lua.LState, path string) (*lua.LState, error) {
	if root == nil {
		return nil, fmt.Errorf("nil lua root")
	}
	mod, err := LoadModule(root, path)
	if err != nil {
		return nil, err
	}
	co, _ := root.NewThread()
	setThreadScript(co, path, mod)
	return co, nil
}

func Call(thread *lua.LState, hook string, args ...interface{}) (lua.LValue, error) {
	if thread == nil {
		return nil, fmt.Errorf("nil lua thread")
	}
	fn := threadFunc(thread, hook)
	if fn == nil {
		if shouldAutoClose(thread) {
			Close(thread)
		}
		return nil, nil
	}
	lvArgs, err := toLValues(thread, args)
	if err != nil {
		return nil, err
	}
	thread.Push(fn)
	for _, lv := range lvArgs {
		thread.Push(lv)
	}
	if err := thread.PCall(len(lvArgs), 1, nil); err != nil {
		if shouldAutoClose(thread) {
			Close(thread)
		}
		return nil, fmt.Errorf("%s: %w", hook, err)
	}
	ret := thread.Get(-1)
	thread.Pop(1)
	if shouldAutoClose(thread) {
		Close(thread)
	}
	return ret, nil
}

func CallAsync(root *lua.LState, thread *lua.LState, hook string, args ...interface{}) *async.Promise {
	promise := async.NewDeferred()
	if root == nil || thread == nil {
		promise.SetError(fmt.Errorf("nil lua root/thread"))
		return promise
	}
	fn := threadFunc(thread, hook)
	if fn == nil {
		promise.SetResult(nil)
		if shouldAutoClose(thread) {
			Close(thread)
		}
		return promise
	}
	cfg, _ := GetConfiguration(thread)
	cfg.CallPromise = promise
	SetConfiguration(thread, cfg)
	_, _, err := resumeFn(root, thread, fn, hook, args...)
	if err != nil && !promise.Completed() {
		promise.SetError(err)
	}
	return promise
}

func Resume(root *lua.LState, thread *lua.LState, args ...interface{}) (lua.ResumeState, []lua.LValue, error) {
	return resumeFn(root, thread, nil, "", args...)
}

func resumeFn(root *lua.LState, thread *lua.LState, fn *lua.LFunction, hook string, args ...interface{}) (lua.ResumeState, []lua.LValue, error) {
	if root == nil || thread == nil {
		return lua.ResumeOK, nil, fmt.Errorf("nil lua root/thread")
	}
	lvArgs, err := toLValues(thread, args)
	if err != nil {
		return lua.ResumeOK, nil, err
	}
	if fn != nil {
		thread.Push(fn)
	}
	state, resumeErr, values := root.Resume(thread, fn, lvArgs...)
	if resumeErr != nil {
		err := resumeErr
		if hook != "" {
			err = fmt.Errorf("%s: %w", hook, resumeErr)
		}
		completeCall(thread, nil, err)
		if shouldAutoClose(thread) {
			Close(thread)
		}
		return lua.ResumeOK, nil, err
	}
	if state != lua.ResumeYield {
		completeCall(thread, values, nil)
		if shouldAutoClose(thread) {
			Close(thread)
		}
	}
	return state, values, nil
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
		case int:
			out[i] = lua.LNumber(v)
		case int16:
			out[i] = lua.LNumber(v)
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
