package wz

import lua "github.com/yuin/gopher-lua"

func (*ItemCore) LuaTypeName() string { return luaWzItemCoreTypeName }
func (*ItemCore) LuaBuiltinFuncs() map[string]lua.LGFunction {
	return map[string]lua.LGFunction{
		"id": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			pushWzItemID(L, ud.Value)
			return 1
		},
	}
}
func (w *ItemCore) String() string       { return w.LuaTypeName() }
func (w *ItemCore) Type() lua.LValueType { return lua.LTUserData }
