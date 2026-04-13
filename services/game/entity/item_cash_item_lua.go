package entity

import lua "github.com/yuin/gopher-lua"

func (*CashItem) LuaTypeName() string { return "LuaCashItem" }
func (*CashItem) LuaBuiltinFuncs() map[string]lua.LGFunction {
	return map[string]lua.LGFunction{}
}
func (c *CashItem) String() string       { return c.LuaTypeName() }
func (c *CashItem) Type() lua.LValueType { return lua.LTUserData }
