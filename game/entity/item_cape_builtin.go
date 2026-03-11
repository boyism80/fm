package entity

import lua "github.com/yuin/gopher-lua"

func (*Cape) LuaTypeName() string { return "LuaCape" }
func (*Cape) LuaBuiltinFuncs() map[string]lua.LGFunction {
	return map[string]lua.LGFunction{}
}
func (c *Cape) String() string       { return c.LuaTypeName() }
func (c *Cape) Type() lua.LValueType { return lua.LTUserData }
