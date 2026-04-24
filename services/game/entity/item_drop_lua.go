package entity

import lua "github.com/yuin/gopher-lua"

func (fp *FieldPlacement) LuaTypeName() string { return "LuaFieldPlacement" }
func (fp *FieldPlacement) LuaBuiltinFuncs() map[string]lua.LGFunction {
	return map[string]lua.LGFunction{}
}
func (fp *FieldPlacement) String() string       { return fp.LuaTypeName() }
func (fp *FieldPlacement) Type() lua.LValueType { return lua.LTUserData }
