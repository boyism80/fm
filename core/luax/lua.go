package luax

import (
	"fmt"
	"os"
	"sync"

	lua "github.com/yuin/gopher-lua"
)

var (
	onCreateHooks   []func(*lua.LState)
	onCreateHooksMu sync.Mutex
	compileMu       sync.Mutex
	compiledFuncs   = make(map[string]*lua.LFunction)
	useCache        = os.Getenv("GO_ENV") != "development"
	rootStates      sync.Map
)

func RegisterRootLuaState(pid string, L *lua.LState) {
	rootStates.Store(pid, L)
}

func GetRootLuaState(pid string) *lua.LState {
	v, ok := rootStates.Load(pid)
	if !ok {
		return nil
	}
	return v.(*lua.LState)
}

func UnregisterRootLuaState(pid string) {
	if v, ok := rootStates.LoadAndDelete(pid); ok {
		if L, ok := v.(*lua.LState); ok && L != nil {
			L.Close()
		}
	}
}

func init() {
	env, ok := os.LookupEnv("GO_ENV")
	if !ok || env == "" {
		env = "development"
	}
	useCache = env != "development"
}

type Luable interface {
	lua.LValue
	LuaTypeName() string
	LuaBuiltinFuncs() map[string]lua.LGFunction
}

func NewState() *lua.LState {
	luaState := lua.NewState()
	onCreateHooksMu.Lock()
	for _, hook := range onCreateHooks {
		hook(luaState)
	}
	onCreateHooksMu.Unlock()
	return luaState
}

func preloadScript(root *lua.LState, path string) (*lua.LFunction, error) {
	compileMu.Lock()
	defer compileMu.Unlock()

	if fn, ok := compiledFuncs[path]; ok && useCache {
		return fn, nil
	}

	fn, err := root.LoadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to compile %s: %w", path, err)
	}
	compiledFuncs[path] = fn
	return fn, nil
}

func NewThread(root *lua.LState, path string) (*lua.LState, error) {
	fn, err := preloadScript(root, path)
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

func RegisterOnCreateHook(fn func(*lua.LState)) {
	onCreateHooksMu.Lock()
	defer onCreateHooksMu.Unlock()
	onCreateHooks = append(onCreateHooks, fn)
}

func RegisterFunc(L *lua.LState, name string, fn lua.LGFunction) {
	L.SetGlobal(name, L.NewFunction(fn))
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

func LValueToInterface(lv lua.LValue) (interface{}, bool) {
	if lv == nil {
		return nil, true
	}
	switch v := lv.(type) {
	case *lua.LUserData:
		return v.Value, true
	case lua.LNumber:
		return float64(v), true
	case lua.LString:
		return string(v), true
	case lua.LBool:
		return bool(v), true
	case *lua.LNilType:
		return nil, true
	default:
		return nil, false
	}
}
