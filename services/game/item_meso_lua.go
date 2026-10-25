package game

import (
	"github.com/boyism80/fm/core/luax"
	lua "github.com/yuin/gopher-lua"
)

func (m *Meso) LuaTypeName() string { return "LuaMeso" }

func (m *Meso) LuaBuiltinFuncs() map[string]lua.LGFunction {
	return map[string]lua.LGFunction{
		"count": func(L *lua.LState) int {
			meso := luax.Check[Meso](L, 1)
			if L.GetTop() != 1 {
				L.ArgError(2, "count() is read-only")
				return 0
			}
			L.Push(lua.LNumber(meso.count))
			return 1
		},
	}
}

func (m *Meso) String() string       { return m.LuaTypeName() }
func (m *Meso) Type() lua.LValueType { return lua.LTUserData }
