package entity

import lua "github.com/yuin/gopher-lua"

func (*Shield) LuaTypeName() string { return "LuaShield" }
func (*Shield) LuaBuiltinFuncs() map[string]lua.LGFunction {
	return map[string]lua.LGFunction{}
}
func (s *Shield) String() string       { return s.LuaTypeName() }
func (s *Shield) Type() lua.LValueType { return lua.LTUserData }
