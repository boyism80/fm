package entity

import lua "github.com/yuin/gopher-lua"

func (*Weapon) LuaTypeName() string { return "LuaWeapon" }
func (*Weapon) LuaBuiltinFuncs() map[string]lua.LGFunction {
	return map[string]lua.LGFunction{}
}
func (w *Weapon) String() string       { return w.LuaTypeName() }
func (w *Weapon) Type() lua.LValueType { return lua.LTUserData }
