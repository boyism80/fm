package entity

import lua "github.com/yuin/gopher-lua"

func (*MiscItem) LuaTypeName() string { return "LuaMiscItem" }
func (*MiscItem) LuaBuiltinFuncs() map[string]lua.LGFunction {
	return map[string]lua.LGFunction{}
}
func (m *MiscItem) String() string       { return m.LuaTypeName() }
func (m *MiscItem) Type() lua.LValueType { return lua.LTUserData }
