package luax

import (
	"fmt"
	"sync"
	"time"

	lua "github.com/yuin/gopher-lua"
)

const modulesRegistryKey = "fm.script_modules"

var (
	onCreateHooks   []func(*lua.LState)
	onCreateHooksMu sync.Mutex
	compileMu       sync.Mutex
	compiledProtos  = make(map[string]*lua.FunctionProto)
	alwaysReload    = false
)

func SetAlwaysReload(enabled bool) {
	compileMu.Lock()
	defer compileMu.Unlock()
	alwaysReload = enabled
	if enabled {
		compiledProtos = make(map[string]*lua.FunctionProto)
	}
}

func AlwaysReload() bool {
	compileMu.Lock()
	defer compileMu.Unlock()
	return alwaysReload
}

type Luable interface {
	lua.LValue
	LuaTypeName() string
	LuaBuiltinFuncs() map[string]lua.LGFunction
}

func NewState() *lua.LState {
	luaState := lua.NewState()
	_ = luaState.DoString(fmt.Sprintf(
		"math.randomseed(%d); math.random(); math.random(); math.random()",
		time.Now().UnixNano(),
	))
	RegisterRequire(luaState)
	onCreateHooksMu.Lock()
	for _, hook := range onCreateHooks {
		hook(luaState)
	}
	onCreateHooksMu.Unlock()
	return luaState
}

func preloadProto(root *lua.LState, path string) (*lua.FunctionProto, error) {
	compileMu.Lock()
	defer compileMu.Unlock()

	optimized := !alwaysReload
	if proto, ok := compiledProtos[path]; ok && optimized {
		if proto == nil {
			return nil, fmt.Errorf("failed to compile %s: cached missing script", path)
		}
		return proto, nil
	}

	fn, err := root.LoadFile(path)
	if err != nil {
		if optimized {
			compiledProtos[path] = nil
		}
		return nil, fmt.Errorf("failed to compile %s: %w", path, err)
	}
	if fn == nil || fn.Proto == nil {
		if optimized {
			compiledProtos[path] = nil
		}
		return nil, fmt.Errorf("failed to compile %s: empty proto", path)
	}
	compiledProtos[path] = fn.Proto
	return fn.Proto, nil
}

func modulesTable(root *lua.LState) *lua.LTable {
	if root == nil {
		return nil
	}
	reg := root.Get(lua.RegistryIndex)
	rt, ok := reg.(*lua.LTable)
	if !ok {
		return nil
	}
	v := root.GetField(rt, modulesRegistryKey)
	if tbl, ok := v.(*lua.LTable); ok {
		return tbl
	}
	tbl := root.NewTable()
	root.SetField(rt, modulesRegistryKey, tbl)
	return tbl
}

func ClearModules(root *lua.LState) {
	if root == nil {
		return
	}
	reg := root.Get(lua.RegistryIndex)
	rt, ok := reg.(*lua.LTable)
	if !ok {
		return
	}
	root.SetField(rt, modulesRegistryKey, root.NewTable())
}

func LoadModule(root *lua.LState, path string) (*lua.LTable, error) {
	if root == nil {
		return nil, fmt.Errorf("nil lua root")
	}
	if path == "" {
		return nil, fmt.Errorf("empty script path")
	}

	compileMu.Lock()
	reload := alwaysReload
	compileMu.Unlock()
	if reload {
		clearRequireCache(root)
		mods := modulesTable(root)
		if mods != nil {
			mods.RawSetString(path, lua.LNil)
		}
	}

	mods := modulesTable(root)
	if mods != nil {
		if cached := mods.RawGetString(path); cached != lua.LNil {
			if tbl, ok := cached.(*lua.LTable); ok {
				return tbl, nil
			}
		}
	}

	proto, err := preloadProto(root, path)
	if err != nil {
		return nil, err
	}
	fn := root.NewFunctionFromProto(proto)
	root.Push(fn)
	if err := root.PCall(0, 1, nil); err != nil {
		return nil, fmt.Errorf("script runtime error: %w", err)
	}
	ret := root.Get(-1)
	root.Pop(1)
	tbl, ok := ret.(*lua.LTable)
	if !ok {
		return nil, fmt.Errorf("%s: script must return a module table", path)
	}
	if mods != nil {
		mods.RawSetString(path, tbl)
	}
	return tbl, nil
}

func ModuleFunc(mod *lua.LTable, name string) *lua.LFunction {
	if mod == nil || name == "" {
		return nil
	}
	v := mod.RawGetString(name)
	fn, ok := v.(*lua.LFunction)
	if !ok {
		return nil
	}
	return fn
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
