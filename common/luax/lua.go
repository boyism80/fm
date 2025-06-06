package luax

import (
	"fmt"
	"sync"

	lua "github.com/yuin/gopher-lua"
)

var (
	onCreateHooks   []func(*lua.LState)
	onCreateHooksMu sync.Mutex
	root            *lua.LState
	compileMu       sync.Mutex
	compiledFuncs   = make(map[string]*lua.LFunction)
)

type Luable interface {
	LuaTypeName() string
	LuaBuiltinFuncs() map[string]lua.LGFunction
}

func preloadScript(path string) (*lua.LFunction, error) {
	compileMu.Lock()
	defer compileMu.Unlock()

	if fn, ok := compiledFuncs[path]; ok {
		return fn, nil
	}

	if root == nil {
		root = lua.NewState()
		onCreateHooksMu.Lock()
		for _, hook := range onCreateHooks {
			hook(root)
		}
		onCreateHooksMu.Unlock()
	}

	fn, err := root.LoadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to compile %s: %w", path, err)
	}
	compiledFuncs[path] = fn
	return fn, nil
}

func NewThread(path string) (*lua.LState, error) {
	fn, err := preloadScript(path)
	if err != nil {
		return nil, err
	}

	co, _ := root.NewThread()
	co.Push(fn)
	if err := co.PCall(0, lua.MultRet, nil); err != nil {
		return nil, fmt.Errorf("script runtime error: %w", err)
	}

	return co, nil
}

func Call(co *lua.LState, funcName string, args ...lua.LValue) (lua.ResumeState, error) {
	fn := co.GetGlobal(funcName)
	if fn.Type() != lua.LTFunction {
		return lua.ResumeYield, fmt.Errorf("on_start is not a function")
	}

	resumeState, err, _ := root.Resume(co, fn.(*lua.LFunction), args...)
	if err != nil {
		return lua.ResumeYield, err
	}
	return resumeState, nil
}

func Resume(co *lua.LState, args ...lua.LValue) (lua.ResumeState, error) {
	resumeState, err, _ := root.Resume(co, nil, args...)
	if err != nil {
		return lua.ResumeYield, err
	}
	return resumeState, nil
}

func RegisterOnCreateHook(fn func(*lua.LState)) {
	onCreateHooksMu.Lock()
	defer onCreateHooksMu.Unlock()
	onCreateHooks = append(onCreateHooks, fn)
}

func Register(L *lua.LState, name string, funcs map[string]lua.LGFunction) {
	mt := L.NewTypeMetatable(name)
	L.SetField(mt, "__index", mt)
	L.SetFuncs(mt, funcs)
}

func RegisterLuaType[T Luable](L *lua.LState) {
	var zero T
	typeName := zero.LuaTypeName()

	mt := L.NewTypeMetatable(typeName)

	L.SetField(mt, "__index", mt)

	L.SetFuncs(mt, zero.LuaBuiltinFuncs())
}

func RegisterLuaDerivedType[T Luable, B Luable](L *lua.LState) {
	var childZero T
	var parentZero B

	childName := childZero.LuaTypeName()
	parentName := parentZero.LuaTypeName()

	childMt := L.NewTypeMetatable(childName)

	parentMt := L.GetTypeMetatable(parentName)
	if parentMt == nil {
		panic(fmt.Sprintf("RegisterLuaDerivedType: 부모 메타테이블 '%s' 가 아직 등록되지 않았습니다.", parentName))
	}

	L.SetMetatable(childMt, parentMt)
	L.SetField(childMt, "__parent", parentMt)
	L.SetField(childMt, "__index", childMt)
	L.SetFuncs(childMt, childZero.LuaBuiltinFuncs())
}

func NewLuable(L *lua.LState, obj Luable) *lua.LUserData {
	ud := L.NewUserData()
	ud.Value = obj
	L.SetMetatable(ud, L.GetTypeMetatable(obj.LuaTypeName()))
	return ud
}
