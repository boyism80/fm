package wz

import lua "github.com/yuin/gopher-lua"

func (*EquipmentCore) LuaTypeName() string { return luaWzEquipmentTypeName }
func (*EquipmentCore) LuaBuiltinFuncs() map[string]lua.LGFunction {
	return map[string]lua.LGFunction{
		"id": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			// Equipment implements Item via embedded *ItemCore.
			pushWzItemID(L, ud.Value)
			return 1
		},
	}
}

func (e *EquipmentCore) String() string       { return e.LuaTypeName() }
func (e *EquipmentCore) Type() lua.LValueType { return lua.LTUserData }
