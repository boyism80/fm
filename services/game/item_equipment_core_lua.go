package game

import (
	"github.com/boyism80/fm/util"
	lua "github.com/yuin/gopher-lua"
)

func (*EquipmentCore) LuaTypeName() string { return "LuaEquipmentCore" }

func (*EquipmentCore) LuaBuiltinFuncs() map[string]lua.LGFunction {
	return map[string]lua.LGFunction{
		"enhance_chance": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			equip, ok := ud.Value.(Equipment)
			if !ok || equip == nil {
				L.Push(lua.LNil)
				return 1
			}
			core := equip.GetEquipmentCore()
			if core == nil {
				L.Push(lua.LNil)
				return 1
			}
			if L.GetTop() >= 2 {
				core.enhanceChance = uint8(util.ClampInt32(int32(L.CheckInt(2)), 0, 255))
			}
			L.Push(lua.LNumber(core.enhanceChance))
			return 1
		},
		"enhance_count": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			equip, ok := ud.Value.(Equipment)
			if !ok || equip == nil {
				L.Push(lua.LNil)
				return 1
			}
			core := equip.GetEquipmentCore()
			if core == nil {
				L.Push(lua.LNil)
				return 1
			}
			L.Push(lua.LNumber(core.enhanceCount))
			return 1
		},
	}
}

func (e *EquipmentCore) String() string       { return e.LuaTypeName() }
func (e *EquipmentCore) Type() lua.LValueType { return lua.LTUserData }
