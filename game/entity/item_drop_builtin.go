package entity

import lua "github.com/yuin/gopher-lua"

func (d *Drop) LuaTypeName() string { return "LuaDrop" }
func (d *Drop) LuaBuiltinFuncs() map[string]lua.LGFunction {
	return map[string]lua.LGFunction{}
}
func (d *Drop) String() string       { return d.LuaTypeName() }
func (d *Drop) Type() lua.LValueType { return lua.LTUserData }
