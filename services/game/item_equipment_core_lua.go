package game

import (
	"github.com/boyism80/fm/maple/wz"
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
		"total_stat": func(L *lua.LState) int {
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
			model, ok := core.wz.(wz.Equipment)
			if !ok || model == nil {
				L.Push(lua.LNil)
				return 1
			}
			if core.bonusStats == nil {
				core.bonusStats = &EquipmentBonusStats{}
			}
			ability := model.GetAbility()
			key := L.CheckString(2)
			switch key {
			case "str":
				L.Push(lua.LNumber(int32(ability.Str) + int32(core.bonusStats.Str)))
			case "dex":
				L.Push(lua.LNumber(int32(ability.Dex) + int32(core.bonusStats.Dex)))
			case "int":
				L.Push(lua.LNumber(int32(ability.Int) + int32(core.bonusStats.Int)))
			case "luk":
				L.Push(lua.LNumber(int32(ability.Luk) + int32(core.bonusStats.Luk)))
			case "max_hp":
				L.Push(lua.LNumber(int32(ability.MaxHP) + int32(core.bonusStats.MaxHP)))
			case "max_mp":
				L.Push(lua.LNumber(int32(ability.MaxMP) + int32(core.bonusStats.MaxMP)))
			case "pad":
				L.Push(lua.LNumber(int32(ability.PAD) + int32(core.bonusStats.PAD)))
			case "mad":
				L.Push(lua.LNumber(int32(ability.MAD) + int32(core.bonusStats.MAD)))
			case "pdd":
				L.Push(lua.LNumber(int32(ability.PDD) + int32(core.bonusStats.PDD)))
			case "mdd":
				L.Push(lua.LNumber(int32(ability.MDD) + int32(core.bonusStats.MDD)))
			case "acc":
				L.Push(lua.LNumber(int32(ability.ACC) + int32(core.bonusStats.ACC)))
			case "avoid":
				L.Push(lua.LNumber(int32(ability.Avoid) + int32(core.bonusStats.Avoid)))
			case "hands":
				L.Push(lua.LNumber(int32(ability.Hands) + int32(core.bonusStats.Hands)))
			case "speed":
				L.Push(lua.LNumber(int32(ability.Speed) + int32(core.bonusStats.Speed)))
			case "jump":
				L.Push(lua.LNumber(int32(ability.Jump) + int32(core.bonusStats.Jump)))
			default:
				L.Push(lua.LNil)
			}
			return 1
		},
		"add_bonus_stats": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			equip, ok := ud.Value.(Equipment)
			if !ok || equip == nil {
				L.Push(lua.LBool(false))
				return 1
			}
			core := equip.GetEquipmentCore()
			if core == nil {
				L.Push(lua.LBool(false))
				return 1
			}
			if core.bonusStats == nil {
				core.bonusStats = &EquipmentBonusStats{}
			}
			tbl := L.CheckTable(2)
			addInt := func(key string, cur int16) int16 {
				lv := tbl.RawGetString(key)
				if num, ok := lv.(lua.LNumber); ok {
					return int16(util.ClampInt32(int32(cur)+int32(num), -32768, 32767))
				}
				return cur
			}
			core.bonusStats.Str = addInt("str", core.bonusStats.Str)
			core.bonusStats.Dex = addInt("dex", core.bonusStats.Dex)
			core.bonusStats.Int = addInt("int", core.bonusStats.Int)
			core.bonusStats.Luk = addInt("luk", core.bonusStats.Luk)
			core.bonusStats.MaxHP = addInt("max_hp", core.bonusStats.MaxHP)
			core.bonusStats.MaxMP = addInt("max_mp", core.bonusStats.MaxMP)
			core.bonusStats.PAD = addInt("pad", core.bonusStats.PAD)
			core.bonusStats.MAD = addInt("mad", core.bonusStats.MAD)
			core.bonusStats.PDD = addInt("pdd", core.bonusStats.PDD)
			core.bonusStats.MDD = addInt("mdd", core.bonusStats.MDD)
			core.bonusStats.ACC = addInt("acc", core.bonusStats.ACC)
			core.bonusStats.Avoid = addInt("avoid", core.bonusStats.Avoid)
			core.bonusStats.Hands = addInt("hands", core.bonusStats.Hands)
			core.bonusStats.Speed = addInt("speed", core.bonusStats.Speed)
			core.bonusStats.Jump = addInt("jump", core.bonusStats.Jump)
			L.Push(lua.LBool(true))
			return 1
		},
	}
}

func (e *EquipmentCore) String() string       { return e.LuaTypeName() }
func (e *EquipmentCore) Type() lua.LValueType { return lua.LTUserData }
