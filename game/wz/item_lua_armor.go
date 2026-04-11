package wz

import lua "github.com/yuin/gopher-lua"

func (*Armor) LuaTypeName() string { return luaWzArmorTypeName }
func (*Armor) LuaBuiltinFuncs() map[string]lua.LGFunction {
	return map[string]lua.LGFunction{
		"id": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			pushWzItemID(L, ud.Value)
			return 1
		},
	}
}
func (a *Armor) String() string       { return a.LuaTypeName() }
func (a *Armor) Type() lua.LValueType { return lua.LTUserData }
