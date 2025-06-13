package builtin

import lua "github.com/yuin/gopher-lua"

type LuaLife struct {
	LuaObject
}

func (life *LuaLife) LuaTypeName() string {
	return "LuaLife"
}

func (life *LuaLife) LuaBuiltinFuncs() map[string]lua.LGFunction {
	return map[string]lua.LGFunction{
		"hp": func(L *lua.LState) int {
			return 0
		},
	}
}
