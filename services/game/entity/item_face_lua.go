package entity

import lua "github.com/yuin/gopher-lua"

func (*Face) LuaTypeName() string { return "LuaFace" }
func (*Face) LuaBuiltinFuncs() map[string]lua.LGFunction {
	return map[string]lua.LGFunction{}
}
func (f *Face) String() string       { return f.LuaTypeName() }
func (f *Face) Type() lua.LValueType { return lua.LTUserData }
