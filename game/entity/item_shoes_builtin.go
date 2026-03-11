package entity

import lua "github.com/yuin/gopher-lua"

func (*Shoes) LuaTypeName() string { return "LuaShoes" }
func (*Shoes) LuaBuiltinFuncs() map[string]lua.LGFunction {
	return map[string]lua.LGFunction{}
}
func (s *Shoes) String() string       { return s.LuaTypeName() }
func (s *Shoes) Type() lua.LValueType { return lua.LTUserData }
