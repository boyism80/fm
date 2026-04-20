package wz

import lua "github.com/yuin/gopher-lua"

func (*Consume) LuaTypeName() string { return luaWzConsumeTypeName }
func (*Consume) LuaBuiltinFuncs() map[string]lua.LGFunction {
	return map[string]lua.LGFunction{
		"id": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			pushWzItemID(L, ud.Value)
			return 1
		},
		"consume_type": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			pushConsumeType(L, ud.Value)
			return 1
		},
	}
}
func (c *Consume) String() string       { return c.LuaTypeName() }
func (c *Consume) Type() lua.LValueType { return lua.LTUserData }
