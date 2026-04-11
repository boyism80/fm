package entity

import lua "github.com/yuin/gopher-lua"

func (*Consume) LuaTypeName() string { return "LuaConsume" }
func (*Consume) LuaBuiltinFuncs() map[string]lua.LGFunction {
	return map[string]lua.LGFunction{}
}
func (c *Consume) String() string       { return c.LuaTypeName() }
func (c *Consume) Type() lua.LValueType { return lua.LTUserData }
