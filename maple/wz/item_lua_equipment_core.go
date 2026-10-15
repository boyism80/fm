package wz

import lua "github.com/yuin/gopher-lua"

func (*EquipmentCore) LuaTypeName() string { return luaWzEquipmentTypeName }
func (*EquipmentCore) LuaBuiltinFuncs() map[string]lua.LGFunction {
	return map[string]lua.LGFunction{
		"id": func(L *lua.LState) int {
			ud := L.CheckUserData(1)

			pushWzItemID(L, ud.Value)
			return 1
		},
		"enhance_chance": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			equip, ok := ud.Value.(Equipment)
			if !ok || equip == nil {
				L.Push(lua.LNil)
				return 1
			}
			L.Push(lua.LNumber(equip.GetEnhanceChance()))
			return 1
		},
	}
}

func (e *EquipmentCore) String() string       { return e.LuaTypeName() }
func (e *EquipmentCore) Type() lua.LValueType { return lua.LTUserData }
