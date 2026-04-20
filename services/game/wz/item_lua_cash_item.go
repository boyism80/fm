package wz

import lua "github.com/yuin/gopher-lua"

func (*CashItem) LuaTypeName() string { return luaWzCashItemTypeName }
func (*CashItem) LuaBuiltinFuncs() map[string]lua.LGFunction {
	return map[string]lua.LGFunction{
		"id": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			pushWzItemID(L, ud.Value)
			return 1
		},
	}
}

func (ci *CashItem) String() string       { return ci.LuaTypeName() }
func (ci *CashItem) Type() lua.LValueType { return lua.LTUserData }
