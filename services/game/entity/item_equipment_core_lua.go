package entity

import (
	"github.com/boyism80/fm/services/game/wz"
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
				core.EnhanceChance = uint8(util.ClampInt32(int32(L.CheckInt(2)), 0, 255))
			}
			L.Push(lua.LNumber(core.EnhanceChance))
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
			if L.GetTop() >= 2 {
				core.EnhanceCount = uint8(util.ClampInt32(int32(L.CheckInt(2)), 0, 255))
			}
			L.Push(lua.LNumber(core.EnhanceCount))
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
			model, ok := core.Wz.(wz.Equipment)
			if !ok || model == nil {
				L.Push(lua.LNil)
				return 1
			}
			if core.BonusStats == nil {
				core.BonusStats = &EquipmentBonusStats{}
			}
			ability := model.GetAbility()
			key := L.CheckString(2)
			switch key {
			case "str":
				L.Push(lua.LNumber(int32(ability.Str) + int32(core.BonusStats.Str)))
			case "dex":
				L.Push(lua.LNumber(int32(ability.Dex) + int32(core.BonusStats.Dex)))
			case "int":
				L.Push(lua.LNumber(int32(ability.Int) + int32(core.BonusStats.Int)))
			case "luk":
				L.Push(lua.LNumber(int32(ability.Luk) + int32(core.BonusStats.Luk)))
			case "max_hp":
				L.Push(lua.LNumber(int32(ability.MaxHP) + int32(core.BonusStats.MaxHP)))
			case "max_mp":
				L.Push(lua.LNumber(int32(ability.MaxMP) + int32(core.BonusStats.MaxMP)))
			case "pad":
				L.Push(lua.LNumber(int32(ability.PAD) + int32(core.BonusStats.PAD)))
			case "mad":
				L.Push(lua.LNumber(int32(ability.MAD) + int32(core.BonusStats.MAD)))
			case "pdd":
				L.Push(lua.LNumber(int32(ability.PDD) + int32(core.BonusStats.PDD)))
			case "mdd":
				L.Push(lua.LNumber(int32(ability.MDD) + int32(core.BonusStats.MDD)))
			case "acc":
				L.Push(lua.LNumber(int32(ability.ACC) + int32(core.BonusStats.ACC)))
			case "avoid":
				L.Push(lua.LNumber(int32(ability.Avoid) + int32(core.BonusStats.Avoid)))
			case "hands":
				L.Push(lua.LNumber(int32(ability.Hands) + int32(core.BonusStats.Hands)))
			case "speed":
				L.Push(lua.LNumber(int32(ability.Speed) + int32(core.BonusStats.Speed)))
			case "jump":
				L.Push(lua.LNumber(int32(ability.Jump) + int32(core.BonusStats.Jump)))
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
			if core.BonusStats == nil {
				core.BonusStats = &EquipmentBonusStats{}
			}
			tbl := L.CheckTable(2)
			addInt := func(key string, cur int16) int16 {
				lv := tbl.RawGetString(key)
				if num, ok := lv.(lua.LNumber); ok {
					return int16(util.ClampInt32(int32(cur)+int32(num), -32768, 32767))
				}
				return cur
			}
			core.BonusStats.Str = addInt("str", core.BonusStats.Str)
			core.BonusStats.Dex = addInt("dex", core.BonusStats.Dex)
			core.BonusStats.Int = addInt("int", core.BonusStats.Int)
			core.BonusStats.Luk = addInt("luk", core.BonusStats.Luk)
			core.BonusStats.MaxHP = addInt("max_hp", core.BonusStats.MaxHP)
			core.BonusStats.MaxMP = addInt("max_mp", core.BonusStats.MaxMP)
			core.BonusStats.PAD = addInt("pad", core.BonusStats.PAD)
			core.BonusStats.MAD = addInt("mad", core.BonusStats.MAD)
			core.BonusStats.PDD = addInt("pdd", core.BonusStats.PDD)
			core.BonusStats.MDD = addInt("mdd", core.BonusStats.MDD)
			core.BonusStats.ACC = addInt("acc", core.BonusStats.ACC)
			core.BonusStats.Avoid = addInt("avoid", core.BonusStats.Avoid)
			core.BonusStats.Hands = addInt("hands", core.BonusStats.Hands)
			core.BonusStats.Speed = addInt("speed", core.BonusStats.Speed)
			core.BonusStats.Jump = addInt("jump", core.BonusStats.Jump)
			L.Push(lua.LBool(true))
			return 1
		},
	}
}

func (e *EquipmentCore) String() string       { return e.LuaTypeName() }
func (e *EquipmentCore) Type() lua.LValueType { return lua.LTUserData }
