package entity

import lua "github.com/yuin/gopher-lua"

func (*Cap) LuaTypeName() string { return "LuaCap" }
func (*Cap) LuaBuiltinFuncs() map[string]lua.LGFunction {
	return map[string]lua.LGFunction{}
}
func (c *Cap) String() string       { return c.LuaTypeName() }
func (c *Cap) Type() lua.LValueType { return lua.LTUserData }
