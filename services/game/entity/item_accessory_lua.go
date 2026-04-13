package entity

import lua "github.com/yuin/gopher-lua"

func (*Accessory) LuaTypeName() string { return "LuaAccessory" }
func (*Accessory) LuaBuiltinFuncs() map[string]lua.LGFunction {
	return map[string]lua.LGFunction{}
}
func (a *Accessory) String() string       { return a.LuaTypeName() }
func (a *Accessory) Type() lua.LValueType { return lua.LTUserData }
