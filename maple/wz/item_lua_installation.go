package wz

import lua "github.com/yuin/gopher-lua"

func (*Installation) LuaTypeName() string { return luaWzInstallationTypeName }
func (*Installation) LuaBuiltinFuncs() map[string]lua.LGFunction {
	return map[string]lua.LGFunction{
		"id": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			pushWzItemID(L, ud.Value)
			return 1
		},
	}
}

func (ins *Installation) String() string       { return ins.LuaTypeName() }
func (ins *Installation) Type() lua.LValueType { return lua.LTUserData }
