package luax

import (
	"fmt"
	"os"
	"sync"
	"time"

	lua "github.com/yuin/gopher-lua"
)

const modulesRegistryKey = "fm.script_modules"

type compiledScript struct {
	proto   *lua.FunctionProto
	modTime time.Time
}

var (
	onCreateHooks   []func(*lua.LState)
	onCreateHooksMu sync.Mutex
	compileMu       sync.Mutex
	compiledProtos  = make(map[string]compiledScript)
	alwaysReload    = false
)

func init() {
	lua.MaxArrayIndex = 1 << 12
}

func SetAlwaysReload(enabled bool) {
	compileMu.Lock()
	defer compileMu.Unlock()
	alwaysReload = enabled
	if enabled {
		compiledProtos = make(map[string]compiledScript)
	}
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

	var modTime time.Time
	if alwaysReload {
		if info, err := os.Stat(path); err == nil {
			modTime = info.ModTime()
		}
	}
	if cached, ok := compiledProtos[path]; ok && cached.modTime.Equal(modTime) {
		if cached.proto == nil {
			return nil, fmt.Errorf("failed to compile %s: cached missing script", path)
		}
		return cached.proto, nil
	}

	fn, err := root.LoadFile(path)
	if err != nil {
		compiledProtos[path] = compiledScript{modTime: modTime}
		return nil, fmt.Errorf("failed to compile %s: %w", path, err)
	}
	if fn == nil || fn.Proto == nil {
		compiledProtos[path] = compiledScript{modTime: modTime}
		return nil, fmt.Errorf("failed to compile %s: empty proto", path)
	}
	compiledProtos[path] = compiledScript{proto: fn.Proto, modTime: modTime}
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

func LoadModule(root *lua.LState, path string) (*lua.LTable, error) {
	if root == nil {
		return nil, fmt.Errorf("nil lua root")
	}
	if path == "" {
		return nil, fmt.Errorf("empty script path")
	}
	mod, err := loadModule(root, path)
	if err != nil {
		return nil, err
	}
	tbl, ok := mod.(*lua.LTable)
	if ok == false {
		return nil, fmt.Errorf("%s: script must return a module table", path)
	}
	return tbl, nil
}

func loadModule(root *lua.LState, path string) (lua.LValue, error) {
	compileMu.Lock()
	reload := alwaysReload
	compileMu.Unlock()

	mods := modulesTable(root)
	if reload == false && mods != nil {
		if cached := mods.RawGetString(path); cached != lua.LNil {
			return cached, nil
		}
	}

	proto, err := preloadProto(root, path)
	if err != nil {
		return nil, err
	}
	root.Push(root.NewFunctionFromProto(proto))
	if err := root.PCall(0, 1, nil); err != nil {
		return nil, fmt.Errorf("script runtime error: %w", err)
	}
	mod := root.Get(-1)
	root.Pop(1)
	if mod == lua.LNil {
		mod = lua.LBool(true)
	}
	if mods != nil {
		mods.RawSetString(path, mod)
	}
	return mod, nil
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
