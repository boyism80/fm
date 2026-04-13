package entity

import (
	"github.com/boyism80/fm/services/game/constant"
	lua "github.com/yuin/gopher-lua"
)

func (*EquipmentCore) LuaTypeName() string { return "LuaEquipmentCore" }

func (*EquipmentCore) LuaBuiltinFuncs() map[string]lua.LGFunction {
	return map[string]lua.LGFunction{
		"parts": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			if ud.Value == nil {
				return 0
			}
			equip, ok := ud.Value.(Equipment)
			if !ok {
				return 0
			}
			model := equip.GetModel()
			if model == nil {
				L.Push(lua.LNil)
				return 1
			}
			parts := constant.GetEquipmentPartsType(model.GetID())
			if parts == 0 {
				L.Push(lua.LNil)
				return 1
			}
			L.Push(lua.LNumber(parts))
			return 1
		},
	}
}

func (e *EquipmentCore) String() string       { return e.LuaTypeName() }
func (e *EquipmentCore) Type() lua.LValueType { return lua.LTUserData }
