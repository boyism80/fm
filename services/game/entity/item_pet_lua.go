package entity

import lua "github.com/yuin/gopher-lua"

func (*Pet) LuaTypeName() string { return "LuaPet" }
func (*Pet) LuaBuiltinFuncs() map[string]lua.LGFunction {
	return map[string]lua.LGFunction{}
}
func (p *Pet) String() string       { return p.LuaTypeName() }
func (p *Pet) Type() lua.LValueType { return lua.LTUserData }
