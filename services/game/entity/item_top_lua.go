package entity

import lua "github.com/yuin/gopher-lua"

func (*Top) LuaTypeName() string { return "LuaTop" }
func (*Top) LuaBuiltinFuncs() map[string]lua.LGFunction {
	return map[string]lua.LGFunction{}
}
func (t *Top) String() string       { return t.LuaTypeName() }
func (t *Top) Type() lua.LValueType { return lua.LTUserData }
