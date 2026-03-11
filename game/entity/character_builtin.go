package entity

import (
	"time"

	"github.com/boyism80/fm/core/luax"
	"github.com/boyism80/fm/game/constant"
	"github.com/boyism80/fm/protocol/response"
	"github.com/boyism80/fm/types"
	lua "github.com/yuin/gopher-lua"
)

func (ch *Character) LuaTypeName() string {
	return "LuaCharacter"
}

func (ch *Character) LuaBuiltinFuncs() map[string]lua.LGFunction {
	return map[string]lua.LGFunction{
		"id": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok {
				L.ArgError(1, "Character expected")
				return 0
			}

			argc := L.GetTop()
			switch argc {
			case 1:
				L.Push(lua.LNumber(ch.GetID()))
				return 1
			default:
				L.ArgError(2, "id() is read-only")
				return 0
			}
		},
		"name": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok {
				L.ArgError(1, "Character expected")
				return 0
			}

			argc := L.GetTop()
			switch argc {
			case 1:
				L.Push(lua.LString(ch.name))
				return 1
			case 2:
				name := L.CheckString(2)
				ch.name = name
				return 0
			default:
				L.ArgError(2, "name() requires 0 or 1 arguments")
				return 0
			}
		},
		"level": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok {
				L.ArgError(1, "Character expected")
				return 0
			}

			argc := L.GetTop()
			switch argc {
			case 1:
				L.Push(lua.LNumber(ch.level))
				return 1
			case 2:
				level := L.CheckInt(2)
				if level < 1 {
					level = 1
				}
				if level > 200 {
					level = 200
				}
				ch.SetLevel(uint8(level))
				return 0
			default:
				L.ArgError(2, "level() requires 0 or 1 arguments")
				return 0
			}
		},
		"exp": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok {
				L.ArgError(1, "Character expected")
				return 0
			}

			argc := L.GetTop()
			switch argc {
			case 1:
				L.Push(lua.LNumber(ch.exp))
				return 1
			case 2:
				value := L.CheckNumber(2)
				if value >= 0 {
					if value < 0 {
						value = 0
					}
					ch.exp = uint32(value)
				} else {
					amount := uint32(-value)
					ch.AddExp(amount)
					return 0
				}
				return 0
			default:
				L.ArgError(2, "exp() requires 0 or 1 arguments")
				return 0
			}
		},
		"meso": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok {
				L.ArgError(1, "Character expected")
				return 0
			}

			argc := L.GetTop()
			switch argc {
			case 1:
				L.Push(lua.LNumber(ch.Meso))
				return 1
			case 2:
				value := L.CheckNumber(2)
				if value >= 0 {
					if value <= 2147483647 {
						ch.SetMeso(int32(value))
					} else {
						ch.SetMeso(2147483647)
					}
				} else {
					amount := int32(-value)
					if ch.Meso < amount {
						ch.SetMeso(0)
					} else {
						ch.SetMeso(ch.Meso - amount)
					}
				}
				if ch.Listener != nil {
					ch.Listener.OnMesoChanged(ch.Meso)
				}
				return 0
			default:
				L.ArgError(2, "meso() requires 0 or 1 arguments")
				return 0
			}
		},
		"chat": func(L *lua.LState) int {
			argc := L.GetTop()
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok {
				L.ArgError(1, "Character expected")
				return 0
			}

			message := L.CheckString(2)
			highlight := false
			if argc > 2 {
				highlight = L.CheckBool(3)
			}
			dontRecordHistory := false
			if argc > 3 {
				dontRecordHistory = L.CheckBool(4)
			}

			if ch.Listener != nil {
				ch.Listener.OnChat(message, highlight, dontRecordHistory)
			}
			return 0
		},
		"buff": func(L *lua.LState) int {
			argc := L.GetTop()
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok {
				L.ArgError(1, "Character expected")
				return 0
			}

			parseFlag := func(argument lua.LValue, argIndex int) (constant.BuffFlag, bool) {
				bfTable, ok := argument.(*lua.LTable)
				if !ok {
					L.ArgError(argIndex, "BuffFlag table expected")
					return constant.BuffFlag{}, false
				}
				maskLV := bfTable.RawGetString("mask")
				posLV := bfTable.RawGetString("position")
				if maskLV.Type() != lua.LTNumber || posLV.Type() != lua.LTNumber {
					L.ArgError(argIndex, "BuffFlag table must have numeric mask and position")
					return constant.BuffFlag{}, false
				}
				return constant.BuffFlag{
					Mask:     uint32(maskLV.(lua.LNumber)),
					Position: int(posLV.(lua.LNumber)),
				}, true
			}

			if argc == 2 {
				flag, ok := parseFlag(L.Get(2), 2)
				if !ok {
					return 0
				}
				entity := ch.Buffs.GetEntity(flag)
				if entity == nil || entity.Wz == nil {
					L.Push(lua.LNil)
					return 1
				}
				view := &SkillEntry{
					Skill:      entity.Wz,
					SkillLevel: int(entity.Level),
				}
				L.Push(luax.NewLuable(L, view))
				return 1
			}

			skillUD := L.CheckUserData(2)
			skillEntry, ok := skillUD.Value.(*SkillEntry)
			if !ok || skillEntry == nil || skillEntry.Skill == nil {
				L.ArgError(2, "SkillEntry with Wz expected")
				return 0
			}
			values := make(map[constant.BuffFlag]int32)

			switch argc {
			case 3:
				valueTable := L.CheckTable(3)
				valueTable.ForEach(func(key lua.LValue, value lua.LValue) {
					flag, ok := parseFlag(key, 3)
					if !ok {
						return
					}
					if value.Type() != lua.LTNumber {
						L.ArgError(3, "buff() values must be numbers")
						return
					}
					values[flag] = int32(value.(lua.LNumber))
				})
			case 4:
				flag, ok := parseFlag(L.Get(3), 3)
				if !ok {
					return 0
				}
				values[flag] = int32(L.CheckNumber(4))
			default:
				L.ArgError(3, "buff() requires (skill, {[flag]=value}) or (skill, flag, value)")
				return 0
			}

			if len(values) == 0 {
				L.ArgError(3, "buff() requires at least one flag-value pair")
				return 0
			}
			ch.Buffs.AddBuff(skillEntry.Skill, uint8(skillEntry.SkillLevel), values)
			return 0
		},
		"unbuff": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok {
				L.ArgError(1, "Character expected")
				return 0
			}
			var flags []constant.BuffFlag
			for i := 2; i <= L.GetTop(); i++ {
				bfTable := L.CheckTable(i)
				maskLV := bfTable.RawGetString("mask")
				posLV := bfTable.RawGetString("position")
				if maskLV.Type() != lua.LTNumber || posLV.Type() != lua.LTNumber {
					continue
				}
				flags = append(flags, constant.BuffFlag{Mask: uint32(maskLV.(lua.LNumber)), Position: int(posLV.(lua.LNumber))})
			}
			ch.Buffs.RemoveBuff(flags)
			return 0
		},
		"buff_value": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok {
				L.ArgError(1, "Character expected")
				return 0
			}

			bfTable := L.CheckTable(2)
			maskLV := bfTable.RawGetString("mask")
			posLV := bfTable.RawGetString("position")
			if maskLV.Type() != lua.LTNumber || posLV.Type() != lua.LTNumber {
				L.ArgError(2, "BuffFlag table must have numeric mask and position")
				return 0
			}

			flag := constant.BuffFlag{Mask: uint32(maskLV.(lua.LNumber)), Position: int(posLV.(lua.LNumber))}
			_, currentValue, ok := ch.Buffs.GetBuffValue(flag)
			if !ok {
				L.Push(lua.LNil)
				return 1
			}

			argc := L.GetTop()
			switch argc {
			case 2:
				L.Push(lua.LNumber(currentValue))
				return 1
			case 3:
				newValue := int32(L.CheckNumber(3))
				_, updated := ch.Buffs.SetBuffValue(flag, newValue)
				if !updated {
					L.Push(lua.LNil)
					return 1
				}
				L.Push(lua.LNumber(newValue))
				return 1
			default:
				L.ArgError(3, "buff_value() requires 1 or 2 arguments after self")
				return 0
			}
		},
		"dialog": func(L *lua.LState) int {
			argc := L.GetTop()
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok {
				L.ArgError(1, "Character expected")
				return 0
			}

			npc := 0
			if argc > 1 {
				npc = L.CheckInt(2)
			}

			message := ""
			if argc > 2 {
				message = L.CheckString(3)
			}

			prev := false
			if argc > 3 {
				prev = L.CheckBool(4)
			}
			next := false
			if argc > 4 {
				next = L.CheckBool(5)
			}

			if ch.Listener != nil {
				ch.Listener.OnDialog(uint32(npc), message, prev, next)
			}
			ch.SetCurrentDialog(L)
			return L.Yield(lua.LNumber(0))
		},
		"dialog_yes_no": func(L *lua.LState) int {
			argc := L.GetTop()
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok {
				L.ArgError(1, "Character expected")
				return 0
			}

			npc := 0
			if argc > 1 {
				npc = L.CheckInt(2)
			}

			message := ""
			if argc > 2 {
				message = L.CheckString(3)
			}

			prev := false
			if argc > 3 {
				prev = L.CheckBool(4)
			}
			next := false
			if argc > 4 {
				next = L.CheckBool(5)
			}

			if ch.Listener != nil {
				ch.Listener.OnDialogYesNo(uint32(npc), message, prev, next)
			}
			ch.SetCurrentDialog(L)
			return L.Yield(lua.LNumber(0))
		},
		"dialog_list": func(L *lua.LState) int {
			argc := L.GetTop()
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok {
				L.ArgError(1, "Character expected")
				return 0
			}

			npc := 0
			if argc > 1 {
				npc = L.CheckInt(2)
			}

			message := ""
			if argc > 2 {
				message = L.CheckString(3)
			}

			selections := []string{}
			if argc > 3 {
				tbl := L.CheckTable(4)
				tbl.ForEach(func(_, value lua.LValue) {
					if str, ok := value.(lua.LString); ok {
						selections = append(selections, string(str))
					}
				})
			}

			if ch.Listener != nil {
				ch.Listener.OnDialogList(uint32(npc), message, selections)
			}
			ch.SetCurrentDialog(L)
			return L.Yield(lua.LNumber(0))
		},
		"dialog_accept": func(L *lua.LState) int {
			argc := L.GetTop()
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok {
				L.ArgError(1, "Character expected")
				return 0
			}

			npc := 0
			if argc > 1 {
				npc = L.CheckInt(2)
			}

			message := ""
			if argc > 2 {
				message = L.CheckString(3)
			}

			enableEscape := false
			if argc > 3 {
				enableEscape = L.CheckBool(4)
			}

			if ch.Listener != nil {
				ch.Listener.OnDialogAccept(uint32(npc), message, enableEscape)
			}
			ch.SetCurrentDialog(L)
			return L.Yield(lua.LNumber(0))
		},
		"dialog_input": func(L *lua.LState) int {
			argc := L.GetTop()
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok {
				L.ArgError(1, "Character expected")
				return 0
			}

			npc := 0
			if argc > 1 {
				npc = L.CheckInt(2)
			}

			message := ""
			if argc > 2 {
				message = L.CheckString(3)
			}

			if ch.Listener != nil {
				ch.Listener.OnDialogInput(uint32(npc), message)
			}
			ch.SetCurrentDialog(L)
			return L.Yield(lua.LNumber(0))
		},
		"notice": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok {
				L.ArgError(1, "Character expected")
				return 0
			}
			message := L.CheckString(2)
			if ch.Listener != nil {
				ch.Listener.OnMessage(constant.MSG_LIGHT_BLUE_TEXT, message)
			}
			return 0
		},
		"skill": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok {
				L.ArgError(1, "Character expected")
				return 0
			}

			skillID := uint32(L.CheckInt(2))

			if ch.Skills == nil {
				L.Push(lua.LNil)
				return 1
			}

			skillEntry, exists := ch.Skills[skillID]
			if !exists || skillEntry == nil {
				L.Push(lua.LNil)
				return 1
			}

			skillUD := luax.NewLuable(L, skillEntry)
			L.Push(skillUD)
			return 1
		},
		"add_skill": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok {
				L.ArgError(1, "Character expected")
				return 0
			}

			skillID := uint32(L.CheckInt(2))
			skills := ch.Skills
			if existing, exists := skills[skillID]; exists && existing != nil {
				L.Push(luax.NewLuable(L, existing))
				return 1
			}

			if ch.Context == nil {
				L.Push(lua.LNil)
				return 1
			}
			resources := ch.Context.GetResources()
			if resources == nil {
				L.Push(lua.LNil)
				return 1
			}

			wzSkill := resources.GetSkill(skillID)
			if wzSkill == nil {
				L.Push(lua.LNil)
				return 1
			}

			masterLevel := 0
			if wzSkill.MasterLevel > 0 {
				masterLevel = wzSkill.MasterLevel
			} else if wzSkill.MaxLevel > 0 {
				masterLevel = wzSkill.MaxLevel
			}

			skillEntry := &SkillEntry{
				Skill:       wzSkill,
				SkillLevel:  0,
				MasterLevel: masterLevel,
				Expiration:  time.Time{},
				Owner:       ch,
			}
			skills[skillID] = skillEntry

			if ch.Listener != nil {
				ch.Send(&response.UpdateSkills{
					SkillID:     skillID,
					Level:       0,
					MasterLevel: int32(skillEntry.MasterLevel),
				}, types.SEND_POLICY_ENCRYPT)
			}

			L.Push(luax.NewLuable(L, skillEntry))
			return 1
		},
		"skills": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok {
				L.ArgError(1, "Character expected")
				return 0
			}
			if L.GetTop() != 1 {
				L.ArgError(2, "skills() is read-only")
				return 0
			}

			skillsTable := L.NewTable()
			if ch.Skills == nil {
				L.Push(skillsTable)
				return 1
			}

			for skillID, skillEntry := range ch.Skills {
				if skillEntry == nil {
					continue
				}
				skillUD := luax.NewLuable(L, skillEntry)
				skillsTable.RawSetInt(int(skillID), skillUD)
			}

			L.Push(skillsTable)
			return 1
		},
		"class": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok {
				L.ArgError(1, "Character expected")
				return 0
			}
			argc := L.GetTop()
			switch argc {
			case 1:
				L.Push(lua.LNumber(ch.Class))
				return 1
			case 2:
				classValue := L.CheckInt(2)
				ch.ChangeClass(uint16(classValue))
				return 0
			default:
				L.ArgError(2, "class() requires 0 or 1 arguments")
				return 0
			}
		},
		"hidden": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok {
				L.ArgError(1, "Character expected")
				return 0
			}
			argc := L.GetTop()
			switch argc {
			case 1:
				L.Push(lua.LBool(ch.IsHidden()))
				return 1
			case 2:
				hidden := L.CheckBool(2)
				ch.SetHidden(hidden)
				return 0
			default:
				L.ArgError(2, "hidden() requires 0 or 1 arguments")
				return 0
			}
		},
		"base_str": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok {
				L.ArgError(1, "Character expected")
				return 0
			}
			argc := L.GetTop()
			switch argc {
			case 1:
				L.Push(lua.LNumber(ch.BaseStats.Str))
				return 1
			case 2:
				v := L.CheckInt(2)
				if v < 0 {
					v = 0
				}
				if v > int(constant.STAT_MAX_STR_DEX_INT_LUK) {
					v = int(constant.STAT_MAX_STR_DEX_INT_LUK)
				}
				ch.BaseStats.Str = uint16(v)
				ch.notifyStatChange(constant.STAT_STR)
				return 0
			default:
				L.ArgError(2, "base_str() requires 0 or 1 arguments")
				return 0
			}
		},
		"base_dex": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok {
				L.ArgError(1, "Character expected")
				return 0
			}
			argc := L.GetTop()
			switch argc {
			case 1:
				L.Push(lua.LNumber(ch.BaseStats.Dex))
				return 1
			case 2:
				v := L.CheckInt(2)
				if v < 0 {
					v = 0
				}
				if v > int(constant.STAT_MAX_STR_DEX_INT_LUK) {
					v = int(constant.STAT_MAX_STR_DEX_INT_LUK)
				}
				ch.BaseStats.Dex = uint16(v)
				ch.notifyStatChange(constant.STAT_DEX)
				return 0
			default:
				L.ArgError(2, "base_dex() requires 0 or 1 arguments")
				return 0
			}
		},
		"base_int": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok {
				L.ArgError(1, "Character expected")
				return 0
			}
			argc := L.GetTop()
			switch argc {
			case 1:
				L.Push(lua.LNumber(ch.BaseStats.Int))
				return 1
			case 2:
				v := L.CheckInt(2)
				if v < 0 {
					v = 0
				}
				if v > int(constant.STAT_MAX_STR_DEX_INT_LUK) {
					v = int(constant.STAT_MAX_STR_DEX_INT_LUK)
				}
				ch.BaseStats.Int = uint16(v)
				ch.notifyStatChange(constant.STAT_INT)
				return 0
			default:
				L.ArgError(2, "base_int() requires 0 or 1 arguments")
				return 0
			}
		},
		"base_luk": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok {
				L.ArgError(1, "Character expected")
				return 0
			}
			argc := L.GetTop()
			switch argc {
			case 1:
				L.Push(lua.LNumber(ch.BaseStats.Luk))
				return 1
			case 2:
				v := L.CheckInt(2)
				if v < 0 {
					v = 0
				}
				if v > int(constant.STAT_MAX_STR_DEX_INT_LUK) {
					v = int(constant.STAT_MAX_STR_DEX_INT_LUK)
				}
				ch.BaseStats.Luk = uint16(v)
				ch.notifyStatChange(constant.STAT_LUK)
				return 0
			default:
				L.ArgError(2, "base_luk() requires 0 or 1 arguments")
				return 0
			}
		},
		"bonus_str": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok {
				L.ArgError(1, "Character expected")
				return 0
			}
			argc := L.GetTop()
			switch argc {
			case 1:
				L.Push(lua.LNumber(ch.BonusStats.Str))
				return 1
			case 2:
				ch.BonusStats.Str = int16(L.CheckInt(2))
				ch.notifyStatChange(constant.STAT_STR)
				return 0
			default:
				L.ArgError(2, "bonus_str() requires 0 or 1 arguments")
				return 0
			}
		},
		"bonus_dex": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok {
				L.ArgError(1, "Character expected")
				return 0
			}
			argc := L.GetTop()
			switch argc {
			case 1:
				L.Push(lua.LNumber(ch.BonusStats.Dex))
				return 1
			case 2:
				ch.BonusStats.Dex = int16(L.CheckInt(2))
				ch.notifyStatChange(constant.STAT_DEX)
				return 0
			default:
				L.ArgError(2, "bonus_dex() requires 0 or 1 arguments")
				return 0
			}
		},
		"bonus_int": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok {
				L.ArgError(1, "Character expected")
				return 0
			}
			argc := L.GetTop()
			switch argc {
			case 1:
				L.Push(lua.LNumber(ch.BonusStats.Int))
				return 1
			case 2:
				ch.BonusStats.Int = int16(L.CheckInt(2))
				ch.notifyStatChange(constant.STAT_INT)
				return 0
			default:
				L.ArgError(2, "bonus_int() requires 0 or 1 arguments")
				return 0
			}
		},
		"bonus_luk": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok {
				L.ArgError(1, "Character expected")
				return 0
			}
			argc := L.GetTop()
			switch argc {
			case 1:
				L.Push(lua.LNumber(ch.BonusStats.Luk))
				return 1
			case 2:
				ch.BonusStats.Luk = int16(L.CheckInt(2))
				ch.notifyStatChange(constant.STAT_LUK)
				return 0
			default:
				L.ArgError(2, "bonus_luk() requires 0 or 1 arguments")
				return 0
			}
		},
		"bonus_hp": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok {
				L.ArgError(1, "Character expected")
				return 0
			}
			argc := L.GetTop()
			switch argc {
			case 1:
				L.Push(lua.LNumber(ch.BonusStats.MaxHpFixed))
				return 1
			case 2:
				ch.BonusStats.MaxHpFixed = int16(L.CheckInt(2))
				if ch.Hp > ch.GetMaxHp() {
					ch.Hp = ch.GetMaxHp()
				}
				ch.notifyStatChange(constant.STAT_MAX_HP)
				return 0
			default:
				L.ArgError(2, "bonus_hp() requires 0 or 1 arguments")
				return 0
			}
		},
		"bonus_mp": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok {
				L.ArgError(1, "Character expected")
				return 0
			}
			argc := L.GetTop()
			switch argc {
			case 1:
				L.Push(lua.LNumber(ch.BonusStats.MaxMpFixed))
				return 1
			case 2:
				ch.BonusStats.MaxMpFixed = int16(L.CheckInt(2))
				if ch.Mp > ch.GetMaxMp() {
					ch.Mp = ch.GetMaxMp()
				}
				ch.notifyStatChange(constant.STAT_MAX_MP)
				return 0
			default:
				L.ArgError(2, "bonus_mp() requires 0 or 1 arguments")
				return 0
			}
		},
		"bonus_max_hp": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok {
				L.ArgError(1, "Character expected")
				return 0
			}
			argc := L.GetTop()
			switch argc {
			case 1:
				L.Push(lua.LNumber(ch.BonusStats.MaxHpFixed))
				L.Push(lua.LNumber(ch.BonusStats.MaxHpPercent))
				return 2
			case 2, 3:
				ch.BonusStats.MaxHpFixed = int16(L.CheckInt(2))
				if argc == 3 {
					ch.BonusStats.MaxHpPercent = int16(L.CheckInt(3))
				}
				if ch.Hp > ch.GetMaxHp() {
					ch.Hp = ch.GetMaxHp()
				}
				if ch.Listener != nil {
					stats := map[constant.Stat]int32{
						constant.STAT_MAX_HP: int32(ch.GetMaxHp()),
						constant.STAT_HP:     int32(ch.Hp),
					}
					ch.Listener.OnUpdateStats(stats, false)
				}
				return 0
			default:
				L.ArgError(2, "bonus_max_hp() requires 0, 1 or 2 arguments")
				return 0
			}
		},
		"bonus_max_mp": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok {
				L.ArgError(1, "Character expected")
				return 0
			}
			argc := L.GetTop()
			switch argc {
			case 1:
				L.Push(lua.LNumber(ch.BonusStats.MaxMpFixed))
				L.Push(lua.LNumber(ch.BonusStats.MaxMpPercent))
				return 2
			case 2, 3:
				ch.BonusStats.MaxMpFixed = int16(L.CheckInt(2))
				if argc == 3 {
					ch.BonusStats.MaxMpPercent = int16(L.CheckInt(3))
				}
				if ch.Mp > ch.GetMaxMp() {
					ch.Mp = ch.GetMaxMp()
				}
				if ch.Listener != nil {
					stats := map[constant.Stat]int32{
						constant.STAT_MAX_MP: int32(ch.GetMaxMp()),
						constant.STAT_MP:     int32(ch.Mp),
					}
					ch.Listener.OnUpdateStats(stats, false)
				}
				return 0
			default:
				L.ArgError(2, "bonus_max_mp() requires 0, 1 or 2 arguments")
				return 0
			}
		},
		"bonus_meso_multiplier": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok {
				L.ArgError(1, "Character expected")
				return 0
			}
			argc := L.GetTop()
			switch argc {
			case 1:
				L.Push(lua.LNumber(ch.BonusStats.MesoMultiplier))
				return 1
			case 2:
				ch.BonusStats.MesoMultiplier = int16(L.CheckInt(2))
				return 0
			default:
				L.ArgError(2, "bonus_meso_multiplier() requires 1 or 2 arguments")
				return 0
			}
		},
		"bonus_drop_rate": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok {
				L.ArgError(1, "Character expected")
				return 0
			}
			argc := L.GetTop()
			switch argc {
			case 1:
				L.Push(lua.LNumber(ch.BonusStats.DropRate))
				return 1
			case 2:
				ch.BonusStats.DropRate = int16(L.CheckInt(2))
				return 0
			default:
				L.ArgError(2, "bonus_drop_rate() requires 1 or 2 arguments")
				return 0
			}
		},
		"map": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok {
				L.ArgError(1, "Character expected")
				return 0
			}
			argc := L.GetTop()
			if argc == 1 {
				// Getter: return current map (or nil)
				m := ch.GetMap()
				if m == nil {
					L.Push(lua.LNil)
					return 1
				}
				L.Push(luax.NewLuable(L, m))
				return 1
			}
			// Setter: map(targetMap [, spawnPoint]) -> Warp
			targetMapUD := L.CheckUserData(2)
			targetMap, ok := targetMapUD.Value.(*Map)
			if !ok || targetMap == nil {
				L.ArgError(2, "Map expected")
				return 0
			}
			spawnPoint := uint8(1)
			if argc >= 3 {
				spawnPoint = uint8(L.CheckInt(3))
			}
			if err := ch.Warp(targetMap, spawnPoint); err != nil {
				L.RaiseError("warp: %v", err)
				return 0
			}
			return 0
		},
		"equipped": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok {
				L.ArgError(1, "Character expected")
				return 0
			}
			if L.GetTop() < 2 {
				L.ArgError(2, "equipped(part) requires equipment part (e.g. EquipmentPart.Weapon)")
				return 0
			}
			part := constant.EquipmentPartsType(L.CheckInt(2))
			if eq, ok := ch.Equipments[part]; ok && eq != nil {
				L.Push(luax.NewLuable(L, eq.(luax.Luable)))
				return 1
			}
			L.Push(lua.LNil)
			return 1
		},
		"item": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok {
				L.ArgError(1, "Character expected")
				return 0
			}
			if L.GetTop() < 3 {
				L.ArgError(2, "item(inv_type, slot) requires two arguments")
				return 0
			}
			invType := constant.InventoryType(L.CheckInt(2))
			slot := int16(L.CheckInt(3))
			item := ch.GetItem(invType, slot)
			if item == nil {
				L.Push(lua.LNil)
				return 1
			}
			switch v := item.(type) {
			case Equipment:
				L.Push(luax.NewLuable(L, v.(luax.Luable)))
			case *Consume:
				L.Push(luax.NewLuable(L, v))
			case *CashItem:
				L.Push(luax.NewLuable(L, v))
			case *GeneralItem:
				L.Push(luax.NewLuable(L, v))
			case *Installation:
				L.Push(luax.NewLuable(L, v))
			case *Pet:
				L.Push(luax.NewLuable(L, v))
			default:
				L.Push(luax.NewLuable(L, item.GetObject()))
			}
			return 1
		},
		"mkitem": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok {
				L.ArgError(1, "Character expected")
				return 0
			}
			if ch.Context == nil {
				L.Push(lua.LNil)
				return 1
			}
			resources := ch.Context.GetResources()
			if resources == nil {
				L.Push(lua.LNil)
				return 1
			}

			argc := L.GetTop()
			if argc < 2 {
				L.ArgError(2, "mkitem(itemIdOrName [, count]) requires at least one argument")
				return 0
			}

			var itemId uint32
			switch lv := L.Get(2).(type) {
			case lua.LString:
				id, ok := resources.NameToItem(string(lv))
				if !ok {
					L.Push(lua.LNil)
					return 1
				}
				itemId = id
			case lua.LNumber:
				itemId = uint32(lv)
			default:
				L.ArgError(2, "item id (number) or item name (string) expected")
				return 0
			}

			if _, ok := resources.Items[itemId]; !ok {
				L.Push(lua.LNil)
				return 1
			}

			count := uint16(1)
			if argc >= 3 {
				if n := L.CheckInt(3); n >= 1 {
					count = uint16(n)
				}
			}

			item, err := NewItem(itemId, count, ch.Context)
			if err != nil {
				L.Push(lua.LNil)
				return 1
			}

			addedItems, err := ch.AddItem(item, false)
			if err != nil || len(addedItems) == 0 {
				L.Push(lua.LNil)
				return 1
			}

			addedItem := addedItems[0]
			switch v := addedItem.(type) {
			case Equipment:
				L.Push(luax.NewLuable(L, v.(luax.Luable)))
			case *Consume:
				L.Push(luax.NewLuable(L, v))
			case *CashItem:
				L.Push(luax.NewLuable(L, v))
			case *GeneralItem:
				L.Push(luax.NewLuable(L, v))
			case *Installation:
				L.Push(luax.NewLuable(L, v))
			case *Pet:
				L.Push(luax.NewLuable(L, v))
			default:
				L.Push(luax.NewLuable(L, addedItem.GetObject()))
			}
			return 1
		},
		"rmitem": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok {
				L.ArgError(1, "Character expected")
				return 0
			}
			argc := L.GetTop()
			if argc < 2 {
				L.ArgError(2, "rmitem(itemIdOrName) or rmitem(inv_type, slot, count)")
				return 0
			}
			if argc == 4 {
				invType := constant.InventoryType(L.CheckInt(2))
				slot := int16(L.CheckInt(3))
				count := uint16(L.CheckInt(4))
				if count == 0 {
					L.Push(lua.LBool(false))
					return 1
				}
				ok := ch.RemoveItemCount(invType, slot, count)
				L.Push(lua.LBool(ok))
				return 1
			}
			if ch.Context == nil {
				L.Push(lua.LBool(false))
				return 1
			}
			resources := ch.Context.GetResources()
			if resources == nil {
				L.Push(lua.LBool(false))
				return 1
			}
			var itemId uint32
			switch lv := L.Get(2).(type) {
			case lua.LString:
				id, ok := resources.NameToItem(string(lv))
				if !ok {
					L.Push(lua.LBool(false))
					return 1
				}
				itemId = id
			case lua.LNumber:
				itemId = uint32(lv)
			default:
				L.ArgError(2, "item id (number) or item name (string) expected")
				return 0
			}
			removed := ch.RemoveItemByID(itemId)
			L.Push(lua.LBool(removed))
			return 1
		},
	}
}

func (ch *Character) String() string {
	return ch.LuaTypeName()
}

func (ch *Character) Type() lua.LValueType {
	return lua.LTUserData
}
