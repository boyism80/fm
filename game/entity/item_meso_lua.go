package entity

import (
	lua "github.com/yuin/gopher-lua"
)

func (m *Meso) LuaTypeName() string { return "LuaMeso" }

func (m *Meso) LuaBuiltinFuncs() map[string]lua.LGFunction {
	return map[string]lua.LGFunction{
		"count": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			meso, ok := ud.Value.(*Meso)
			if !ok {
				L.ArgError(1, "Meso expected")
				return 0
			}
			if L.GetTop() != 1 {
				L.ArgError(2, "count() is read-only")
				return 0
			}
			L.Push(lua.LNumber(meso.Count))
			return 1
		},
	}
}

func (m *Meso) String() string       { return m.LuaTypeName() }
func (m *Meso) Type() lua.LValueType { return lua.LTUserData }
