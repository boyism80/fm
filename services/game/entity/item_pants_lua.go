package entity

import lua "github.com/yuin/gopher-lua"

func (*Pants) LuaTypeName() string { return "LuaPants" }
func (*Pants) LuaBuiltinFuncs() map[string]lua.LGFunction {
	return map[string]lua.LGFunction{}
}
func (p *Pants) String() string       { return p.LuaTypeName() }
func (p *Pants) Type() lua.LValueType { return lua.LTUserData }
