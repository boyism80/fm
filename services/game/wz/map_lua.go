package wz

import (
	"github.com/boyism80/fm/core/luax"
	lua "github.com/yuin/gopher-lua"
)

const luaWzMapTypeName = "LuaWzMap"

func (*Map) LuaTypeName() string { return luaWzMapTypeName }

func (*Map) LuaBuiltinFuncs() map[string]lua.LGFunction {
	return map[string]lua.LGFunction{
		"id": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			m, ok := ud.Value.(*Map)
			if !ok || m == nil {
				L.ArgError(1, "WzMap expected")
				return 0
			}
			L.Push(lua.LNumber(m.ID))
			return 1
		},
		"name": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			m, ok := ud.Value.(*Map)
			if !ok || m == nil {
				L.ArgError(1, "WzMap expected")
				return 0
			}
			L.Push(lua.LString(m.Name))
			return 1
		},
		"return_map_id": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			m, ok := ud.Value.(*Map)
			if !ok || m == nil {
				L.ArgError(1, "WzMap expected")
				return 0
			}
			L.Push(lua.LNumber(m.ReturnMapId))
			return 1
		},
		"town": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			m, ok := ud.Value.(*Map)
			if !ok || m == nil {
				L.ArgError(1, "WzMap expected")
				return 0
			}
			L.Push(lua.LBool(m.IsTown))
			return 1
		},
	}
}

func (m *Map) String() string       { return m.LuaTypeName() }
func (m *Map) Type() lua.LValueType { return lua.LTUserData }

var _ luax.Luable = (*Map)(nil)
