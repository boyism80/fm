package game

import (
	"github.com/boyism80/fm/core/luax"

	lua "github.com/yuin/gopher-lua"
)

func (a *Alliance) LuaTypeName() string {
	return "LuaAlliance"
}

func (a *Alliance) LuaBuiltinFuncs() map[string]lua.LGFunction {
	return map[string]lua.LGFunction{
		"id": func(L *lua.LState) int {
			alliance := luax.Check[Alliance](L, 1)
			if L.GetTop() != 1 {
				L.ArgError(2, "id() is read-only")
				return 0
			}
			L.Push(lua.LNumber(alliance.allianceID))
			return 1
		},
		"name": func(L *lua.LState) int {
			alliance := luax.Check[Alliance](L, 1)
			if L.GetTop() != 1 {
				L.ArgError(2, "name() is read-only")
				return 0
			}
			L.Push(lua.LString(alliance.name))
			return 1
		},
		"capacity": func(L *lua.LState) int {
			alliance := luax.Check[Alliance](L, 1)
			if L.GetTop() != 1 {
				L.ArgError(2, "capacity() is read-only")
				return 0
			}
			L.Push(lua.LNumber(alliance.capacity))
			return 1
		},
		"notice": func(L *lua.LState) int {
			alliance := luax.Check[Alliance](L, 1)
			if L.GetTop() != 1 {
				L.ArgError(2, "notice() is read-only")
				return 0
			}
			L.Push(lua.LString(alliance.notice))
			return 1
		},
	}
}

func (a *Alliance) String() string {
	return a.LuaTypeName()
}

func (a *Alliance) Type() lua.LValueType {
	return lua.LTUserData
}
