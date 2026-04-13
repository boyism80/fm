package entity

import lua "github.com/yuin/gopher-lua"

func (*RingEquip) LuaTypeName() string { return "LuaRingEquip" }
func (*RingEquip) LuaBuiltinFuncs() map[string]lua.LGFunction {
	return map[string]lua.LGFunction{}
}
func (r *RingEquip) String() string       { return r.LuaTypeName() }
func (r *RingEquip) Type() lua.LValueType { return lua.LTUserData }
