package entity

import lua "github.com/yuin/gopher-lua"

func (*GeneralItem) LuaTypeName() string { return "LuaGeneralItem" }
func (*GeneralItem) LuaBuiltinFuncs() map[string]lua.LGFunction {
	return map[string]lua.LGFunction{}
}
func (g *GeneralItem) String() string       { return g.LuaTypeName() }
func (g *GeneralItem) Type() lua.LValueType { return lua.LTUserData }
