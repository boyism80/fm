package wz

import lua "github.com/yuin/gopher-lua"

func (*MiscItem) LuaTypeName() string { return luaWzMiscItemTypeName }
func (*MiscItem) LuaBuiltinFuncs() map[string]lua.LGFunction {
	return map[string]lua.LGFunction{
		"id": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			pushWzItemID(L, ud.Value)
			return 1
		},
	}
}

func (mi *MiscItem) String() string       { return mi.LuaTypeName() }
func (mi *MiscItem) Type() lua.LValueType { return lua.LTUserData }
