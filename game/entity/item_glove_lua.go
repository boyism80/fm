package entity

import lua "github.com/yuin/gopher-lua"

func (*Glove) LuaTypeName() string { return "LuaGlove" }
func (*Glove) LuaBuiltinFuncs() map[string]lua.LGFunction {
	return map[string]lua.LGFunction{}
}
func (g *Glove) String() string       { return g.LuaTypeName() }
func (g *Glove) Type() lua.LValueType { return lua.LTUserData }
