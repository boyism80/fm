package entity

import lua "github.com/yuin/gopher-lua"

func (*Installation) LuaTypeName() string { return "LuaInstallation" }
func (*Installation) LuaBuiltinFuncs() map[string]lua.LGFunction {
	return map[string]lua.LGFunction{}
}
func (i *Installation) String() string       { return i.LuaTypeName() }
func (i *Installation) Type() lua.LValueType { return lua.LTUserData }
