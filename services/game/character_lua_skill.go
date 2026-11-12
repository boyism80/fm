package game

import (
	"time"

	"github.com/boyism80/fm/core/luax"
	"github.com/boyism80/fm/maple/constant"
	"github.com/boyism80/fm/maple/wz"
	lua "github.com/yuin/gopher-lua"
)

func (ch *Character) luaSkillFuncs() map[string]lua.LGFunction {
	return map[string]lua.LGFunction{
		"buff": func(L *lua.LState) int {
			argc := L.GetTop()
			ch := luax.Check[Character](L, 1)

			getFlagFromTable := func(argument lua.LValue, argIndex int) (constant.BuffFlag, bool) {
				fields, ok := luax.ParseTable(L, argument, argIndex)
				if !ok {
					return constant.BuffFlag{}, false
				}
				maskLV := fields[lua.LString("mask")]
				posLV := fields[lua.LString("position")]
				if maskLV == nil || posLV == nil || maskLV.Type() != lua.LTNumber || posLV.Type() != lua.LTNumber {
					L.ArgError(argIndex, "BuffFlag table must have numeric mask and position")
					return constant.BuffFlag{}, false
				}
				return constant.BuffFlag{
					Mask:     uint32(maskLV.(lua.LNumber)),
					Position: int(posLV.(lua.LNumber)),
				}, true
			}
			getValuesFromTable := func(argument lua.LValue, argIndex int) (map[constant.BuffFlag]int32, bool) {
				fieldMap, ok := luax.ParseTable(L, argument, argIndex)
				if !ok {
					return nil, false
				}
				if len(fieldMap) == 0 {
					L.ArgError(argIndex, "buff() requires at least one flag-value pair")
					return nil, false
				}
				values := make(map[constant.BuffFlag]int32, len(fieldMap))
				for keyLV, valueLV := range fieldMap {
					flag, ok := getFlagFromTable(keyLV, argIndex)
					if !ok {
						return nil, false
					}
					if valueLV == nil || valueLV.Type() != lua.LTNumber {
						L.ArgError(argIndex, "buff() values must be numbers")
						return nil, false
					}
					values[flag] = int32(valueLV.(lua.LNumber))
				}
				return values, true
			}
			isSingleFlagTable := func(argument lua.LValue, argIndex int) (bool, bool) {
				fields, ok := luax.ParseTable(L, argument, argIndex)
				if !ok {
					return false, false
				}
				maskLV := fields[lua.LString("mask")]
				posLV := fields[lua.LString("position")]
				if maskLV == nil || posLV == nil {
					return false, true
				}
				if maskLV.Type() != lua.LTNumber || posLV.Type() != lua.LTNumber {
					L.ArgError(argIndex, "BuffFlag table must have numeric mask and position")
					return false, false
				}
				return true, true
			}
			// buff(skill) gives the skill's WZ stats alone.
			if ud, ok := L.Get(2).(*lua.LUserData); ok && argc == 2 {
				skillEntry, ok := ud.Value.(*SkillEntry)
				if !ok || skillEntry == nil || skillEntry.wz == nil {
					L.ArgError(2, "SkillEntry with Wz expected")
					return 0
				}
				ld := skillEntry.wz.GetLevelData(skillEntry.Level())
				if ld == nil {
					return 0
				}
				if values := ld.StatBuffs(); len(values) > 0 {
					ch.buffs.AddBuff(skillEntry.wz, ld.Time, uint8(skillEntry.Level()), ch.GetID(), values, true)
				}
				return 0
			}
			if argc == 2 {
				if L.Get(2).Type() != lua.LTTable {
					L.ArgError(2, "buff() getter requires BuffFlag table")
					return 0
				}
				flagFields, ok := luax.ParseTable(L, L.Get(2), 2)
				if !ok {
					return 0
				}
				maskLV := flagFields[lua.LString("mask")]
				posLV := flagFields[lua.LString("position")]
				if maskLV == nil || posLV == nil || maskLV.Type() != lua.LTNumber || posLV.Type() != lua.LTNumber {
					L.ArgError(2, "BuffFlag table must have numeric mask and position")
					return 0
				}
				flag := constant.BuffFlag{
					Mask:     uint32(maskLV.(lua.LNumber)),
					Position: int(posLV.(lua.LNumber)),
				}
				buff := ch.buffs.GetEntity(flag)
				switch entity := buff.(type) {
				case *SkillBuff:
					if entity == nil {
						L.Push(lua.LNil)
						return 1
					}
					L.Push(luax.NewLuable(L, entity))
					return 1
				case *ItemBuff:
					if entity == nil {
						L.Push(lua.LNil)
						return 1
					}
					L.Push(luax.NewLuable(L, entity))
					return 1
				default:
					L.Push(lua.LNil)
					return 1
				}
			} else {
				if argc < 3 {
					L.ArgError(3, "buff() requires (skill, {[flag]=value}), (skill, {[flag]=value}, option), (skill, flag, value), (skill, flag, value, option), (consume, durationMs, {[flag]=value}), or (consume, durationMs, flag, value)")
					return 0
				}

				arg2UD := L.CheckUserData(2)

				var skillEntry *SkillEntry
				var consumeWz *wz.Consume
				var values map[constant.BuffFlag]int32
				baseDuration := time.Duration(0)
				duration := time.Duration(0)
				offset := 2

				if se, ok := arg2UD.Value.(*SkillEntry); ok {
					if se == nil {
						L.ArgError(2, "SkillEntry with Wz expected")
						return 0
					}
					skillEntry = se
					if ld := skillEntry.wz.GetLevelData(skillEntry.Level()); ld != nil {
						baseDuration = ld.Time
					}
					offset++
				} else if cw, ok := arg2UD.Value.(*wz.Consume); ok {
					if cw == nil {
						L.ArgError(2, "wz.Consume expected")
						return 0
					}
					consumeWz = cw
					offset++
					durationMs := L.CheckInt(offset)
					duration = time.Duration(durationMs) * time.Millisecond
					offset++
				} else {
					L.ArgError(2, "SkillEntry with Wz or wz.Consume expected")
					return 0
				}

				flagsArg := L.Get(offset)
				if flagsArg.Type() != lua.LTTable {
					L.ArgError(offset, "buff() requires BuffFlag table or flag-value table")
					return 0
				}

				isFlag, ok := isSingleFlagTable(flagsArg, offset)
				if !ok {
					return 0
				}

				if isFlag {
					offset++
					valueIndex := offset
					if valueIndex > argc {
						if skillEntry != nil {
							L.ArgError(valueIndex, "buff() requires value for (skill, flag, value)")
						} else {
							L.ArgError(valueIndex, "buff() requires value for (consume, durationMs, flag, value)")
						}
						return 0
					}
					flag, ok := getFlagFromTable(flagsArg, offset-1)
					if !ok {
						return 0
					}
					values = map[constant.BuffFlag]int32{flag: int32(L.CheckNumber(valueIndex))}
					offset++
				} else {
					parsedValues, ok := getValuesFromTable(flagsArg, offset)
					if !ok {
						return 0
					}
					values = parsedValues
					// A skill's table buff carries its WZ stats too; a flag the table names keeps the table's value.
					if skillEntry != nil {
						if ld := skillEntry.wz.GetLevelData(skillEntry.Level()); ld != nil {
							for flag, value := range ld.StatBuffs() {
								if _, set := values[flag]; !set {
									values[flag] = value
								}
							}
						}
					}
				}

				if skillEntry != nil {
					duration = baseDuration
					if offset <= argc {
						optionLV := L.Get(offset)
						if optionLV.Type() != lua.LTNil {
							optionFields, ok := luax.ParseTable(L, optionLV, offset)
							if !ok {
								return 0
							}
							timeLV := optionFields[lua.LString("time")]
							if timeLV != nil && timeLV.Type() != lua.LTNil {
								if timeLV.Type() != lua.LTNumber {
									L.ArgError(offset, "buff() option.time must be number")
									return 0
								}
								sec := float64(timeLV.(lua.LNumber))
								if sec < 0 {
									L.ArgError(offset, "buff() option.time must be >= 0")
									return 0
								}
								if sec == 0 {
									duration = 0
								} else {
									duration = time.Duration(sec * float64(time.Second))
								}
							}
						}
					}
					ch.buffs.AddBuff(skillEntry.wz, duration, uint8(skillEntry.Level()), ch.GetID(), values, true)
				} else {
					ch.buffs.AddItemBuff(consumeWz, duration, values, true, true)
				}
				return 0
			}
		},
		"ride": func(L *lua.LState) int {
			ch := luax.Check[Character](L, 1)
			skillEntry := luax.Check[SkillEntry](L, 2)
			if skillEntry.wz == nil {
				L.ArgError(2, "SkillEntry with Wz expected")
				return 0
			}
			L.Push(lua.LBool(ch.mount.Ride(skillEntry, int32(L.CheckInt(3))) == nil))
			return 1
		},
		"dismount": func(L *lua.LState) int {
			ch := luax.Check[Character](L, 1)
			ch.mount.Dismount()
			return 0
		},
		"unbuff": func(L *lua.LState) int {
			ch := luax.Check[Character](L, 1)
			var flags []constant.BuffFlag
			for i := 2; i <= L.GetTop(); i++ {
				lv := L.Get(i)
				switch v := lv.(type) {
				case lua.LNumber:
					ch.buffs.CancelBySource(int32(L.CheckInt(i)))
				case *lua.LTable:
					fields, ok := luax.ParseTable(L, v, i)
					if !ok {
						return 0
					}
					maskLV := fields[lua.LString("mask")]
					posLV := fields[lua.LString("position")]
					if maskLV == nil || posLV == nil || maskLV.Type() != lua.LTNumber || posLV.Type() != lua.LTNumber {
						continue
					}
					flags = append(flags, constant.BuffFlag{
						Mask:     uint32(maskLV.(lua.LNumber)),
						Position: int(posLV.(lua.LNumber)),
					})
				default:
					L.ArgError(i, "BuffFlag table or skill id (number) expected")
					return 0
				}
			}
			if len(flags) > 0 {
				ch.buffs.RemoveBuff(flags)
			}
			return 0
		},
		"buffs": func(L *lua.LState) int {
			ch := luax.Check[Character](L, 1)

			includeItem := true
			includeSkill := true
			if L.GetTop() >= 2 {
				lv := L.Get(2)
				if lv != lua.LNil {
					buffType := constant.BuffType(lua.LVAsNumber(lv))
					switch buffType {
					case constant.BuffTypeAll:
					case constant.BuffTypeItem:
						includeSkill = false
					case constant.BuffTypeSkill:
						includeItem = false
					default:
						L.ArgError(2, "buffs filter must be nil or BuffType")
						return 0
					}
				}
			}

			tbl := L.NewTable()
			idx := 1
			for _, entity := range ch.buffs.Entities() {
				if entity == nil {
					continue
				}
				if includeSkill {
					if skillBuff, ok := entity.(*SkillBuff); ok && skillBuff != nil {
						tbl.RawSetInt(idx, luax.NewLuable(L, skillBuff))
						idx++
						continue
					}
				}
				if includeItem {
					if itemBuff, ok := entity.(*ItemBuff); ok && itemBuff != nil {
						tbl.RawSetInt(idx, luax.NewLuable(L, itemBuff))
						idx++
					}
				}
			}
			L.Push(tbl)
			return 1
		},
		"summons": func(L *lua.LState) int {
			ch := luax.Check[Character](L, 1)
			tbl := L.NewTable()
			idx := 1
			for _, s := range ch.summons.All() {
				tbl.RawSetInt(idx, luax.NewLuable(L, s))
				idx++
			}
			L.Push(tbl)
			return 1
		},
		"buff_value": func(L *lua.LState) int {
			ch := luax.Check[Character](L, 1)

			flagFields, ok := luax.ParseTable(L, L.CheckTable(2), 2)
			if !ok {
				return 0
			}
			maskLV := flagFields[lua.LString("mask")]
			posLV := flagFields[lua.LString("position")]
			if maskLV == nil || posLV == nil || maskLV.Type() != lua.LTNumber || posLV.Type() != lua.LTNumber {
				L.ArgError(2, "BuffFlag table must have numeric mask and position")
				return 0
			}

			flag := constant.BuffFlag{Mask: uint32(maskLV.(lua.LNumber)), Position: int(posLV.(lua.LNumber))}
			_, currentValue, ok := ch.buffs.GetBuffValue(flag)
			if !ok {
				L.Push(lua.LNil)
				return 1
			}

			argc := L.GetTop()
			if argc == 2 {
				L.Push(lua.LNumber(currentValue))
				return 1
			} else if argc == 3 {
				newValue := int32(L.CheckNumber(3))
				_, updated := ch.buffs.SetBuffValue(flag, newValue)
				if !updated {
					L.Push(lua.LNil)
					return 1
				}
				L.Push(lua.LNumber(newValue))
				return 1
			} else {
				L.ArgError(3, "buff_value() requires 1 or 2 arguments after self")
				return 0
			}
		},
		"create_summon": func(L *lua.LState) int {
			ch := luax.Check[Character](L, 1)

			skillID := uint32(L.CheckInt(2))
			skillLevel := uint8(L.CheckInt(3))
			durationMs := L.CheckInt(4)
			if durationMs < 0 {
				durationMs = 0
			}
			movType := constant.SummonMovementType(L.CheckInt(5))
			summonType := constant.SummonType(L.CheckInt(6))

			duration := time.Duration(durationMs) * time.Millisecond
			s := ch.summons.Spawn(constant.SkillID(skillID), skillLevel, movType, summonType, duration)
			if s == nil {
				L.Push(lua.LNil)
				return 1
			}
			L.Push(luax.NewLuable(L, s))
			return 1
		},
		"skill_area": func(L *lua.LState) int {
			ch := luax.Check[Character](L, 1)
			skill := luax.Check[SkillEntry](L, 2)
			ld := skill.wz.GetLevelData(skill.Level())
			if ld == nil {
				L.Push(lua.LNil)
				return 1
			}
			x, y := int32(ch.position.X), int32(ch.position.Y)
			// An attack reports the point it used; one far from the character is not trusted.
			const reach = 500
			if dx, dy := int32(L.OptInt(3, int(x)))-x, int32(L.OptInt(4, int(y)))-y; dx >= -reach && dx <= reach && dy >= -reach && dy <= reach {
				x, y = x+dx, y+dy
			}
			area := ld.Area(x, y, ch.stance%2 != 0)
			tbl := L.NewTable()
			tbl.RawSetString("left", lua.LNumber(area.Left))
			tbl.RawSetString("top", lua.LNumber(area.Top))
			tbl.RawSetString("right", lua.LNumber(area.Right))
			tbl.RawSetString("bottom", lua.LNumber(area.Bottom))
			L.Push(tbl)
			return 1
		},
		"create_mist": func(L *lua.LState) int {
			ch := luax.Check[Character](L, 1)
			skill := luax.Check[SkillEntry](L, 2)
			var spec mistSpec
			spec.Parse(L)
			mist := ch.SpawnMist(skill, spec)
			if mist == nil {
				L.Push(lua.LNil)
				return 1
			}
			L.Push(luax.NewLuable(L, mist))
			return 1
		},
		"create_door": func(L *lua.LState) int {
			ch := luax.Check[Character](L, 1)
			buff := luax.Check[SkillBuff](L, 2)
			ch.doors.Spawn(buff)
			return 0
		},
		"remove_summon": func(L *lua.LState) int {
			ch := luax.Check[Character](L, 1)
			ch.summons.Remove(ch.summons.Get(constant.SkillID(L.CheckInt(2))), true)
			return 0
		},
		"remove_door": func(L *lua.LState) int {
			ch := luax.Check[Character](L, 1)
			buff := luax.Check[SkillBuff](L, 2)
			ch.doors.Remove(buff, true)
			return 0
		},
		"clear_summons": func(L *lua.LState) int {
			ch := luax.Check[Character](L, 1)
			ch.summons.Clear()
			return 0
		},
		"mount": func(L *lua.LState) int {
			ch := luax.Check[Character](L, 1)
			L.Push(lua.LNumber(ch.mount.level))
			L.Push(lua.LNumber(ch.mount.exp))
			L.Push(lua.LNumber(ch.mount.fatigue))
			L.Push(lua.LNumber(ch.mount.hp))
			_, vehicle, _ := ch.buffs.GetBuffValue(constant.BuffFlagMonsterRiding)
			L.Push(lua.LNumber(vehicle))
			return 5
		},
		"set_mount_level": func(L *lua.LState) int {
			ch := luax.Check[Character](L, 1)
			ch.mount.SetLevel(uint32(L.CheckInt(2)))
			return 0
		},
		"set_mount_fatigue": func(L *lua.LState) int {
			ch := luax.Check[Character](L, 1)
			ch.mount.SetFatigue(uint32(L.CheckInt(2)))
			return 0
		},
		"set_mount_hp": func(L *lua.LState) int {
			ch := luax.Check[Character](L, 1)
			ch.mount.hp = uint32(max(L.CheckInt(2), 0))
			return 0
		},
		"start_mount_fatigue": func(L *lua.LState) int {
			ch := luax.Check[Character](L, 1)
			ch.mount.startFatigue()
			return 0
		},
		"show_skill_cooldown": func(L *lua.LState) int {
			ch := luax.Check[Character](L, 1)
			ch.listener.OnSkillCooldown(ch, uint32(L.CheckInt(2)), uint16(L.CheckInt(3)))
			return 0
		},
		"tire_mount": func(L *lua.LState) int {
			ch := luax.Check[Character](L, 1)
			ch.mount.tire()
			return 0
		},
		"skill": func(L *lua.LState) int {
			ch := luax.Check[Character](L, 1)

			skillID := uint32(L.CheckInt(2))

			skillEntry := ch.skills.Get(skillID)
			if skillEntry == nil {
				L.Push(lua.LNil)
				return 1
			}

			skillUD := luax.NewLuable(L, skillEntry)
			L.Push(skillUD)
			return 1
		},
		// skill_at(skill_id, level) is a skill at a level the character has not learned, for a skill's effect an item
		// or an NPC gives; nothing is added to the character's skills.
		"skill_at": func(L *lua.LState) int {
			ch := luax.Check[Character](L, 1)
			wzSkill := ch.gameWorld.GetResources().GetSkill(uint32(L.CheckInt(2)))
			if wzSkill == nil {
				L.Push(lua.LNil)
				return 1
			}
			L.Push(luax.NewLuable(L, NewSkillEntry(ch, wzSkill, L.CheckInt(3), 0)))
			return 1
		},
		"add_skill": func(L *lua.LState) int {
			ch := luax.Check[Character](L, 1)

			skillID := uint32(L.CheckInt(2))
			if existing := ch.skills.Learned(skillID); existing != nil {
				L.Push(luax.NewLuable(L, existing))
				return 1
			}

			if ch.gameWorld == nil {
				L.Push(lua.LNil)
				return 1
			}
			resources := ch.gameWorld.GetResources()
			if resources == nil {
				L.Push(lua.LNil)
				return 1
			}

			wzSkill := resources.GetSkill(skillID)
			if wzSkill == nil {
				L.Push(lua.LNil)
				return 1
			}

			skillEntry := NewSkillEntry(ch, wzSkill, 1, wzSkill.DefaultMasterLevel())
			skillEntry.expiration = time.Time{}
			ch.skills.Register(skillID, skillEntry)

			L.Push(luax.NewLuable(L, skillEntry))
			return 1
		},
		"skills": func(L *lua.LState) int {
			ch := luax.Check[Character](L, 1)
			switch L.GetTop() {
			case 1:

			default:
				L.ArgError(2, "skills() is read-only")
				return 0
			}

			skillsTable := L.NewTable()
			ch.skills.ForEach(func(skillID uint32, skillEntry *SkillEntry) {
				skillUD := luax.NewLuable(L, skillEntry)
				skillsTable.RawSetInt(int(skillID), skillUD)
			})

			L.Push(skillsTable)
			return 1
		},
		"debuff": func(L *lua.LState) int {
			ch := luax.Check[Character](L, 1)
			parseDebuffFlag := func(lv lua.LValue, argIndex int) (constant.DebuffFlag, bool) {
				tbl, ok := lv.(*lua.LTable)
				if !ok {
					L.ArgError(argIndex, "DebuffFlag table expected")
					return constant.DebuffFlag{}, false
				}
				maskLV := tbl.RawGetString("mask")
				posLV := tbl.RawGetString("position")
				if maskLV.Type() != lua.LTNumber || posLV.Type() != lua.LTNumber {
					L.ArgError(argIndex, "DebuffFlag table must have numeric mask and position")
					return constant.DebuffFlag{}, false
				}
				f := constant.DebuffFlag{
					Mask:     uint32(maskLV.(lua.LNumber)),
					Position: int(posLV.(lua.LNumber)),
				}
				if debuffLV := tbl.RawGetString("debuff"); debuffLV.Type() == lua.LTNumber {
					f.DiseaseSkillID = uint16(debuffLV.(lua.LNumber))
				}
				return f, true
			}
			flag, ok := parseDebuffFlag(L.Get(2), 2)
			if !ok {
				return 0
			}
			durationMs := L.CheckNumber(3)
			if durationMs <= 0 {
				L.ArgError(3, "time must be positive")
				return 0
			}
			x := int16(1)
			skillID := flag.DiseaseSkillID
			skillLevel := uint16(1)
			if L.GetTop() >= 4 {
				x = int16(L.CheckNumber(4))
			}
			if L.GetTop() >= 5 {
				skillID = uint16(L.CheckNumber(5))
			}
			if L.GetTop() >= 6 {
				skillLevel = uint16(L.CheckNumber(6))
			}
			ch.debuffs.Give(flag, time.Duration(durationMs)*time.Millisecond, x, skillID, skillLevel)
			return 0
		},
		"has_debuff": func(L *lua.LState) int {
			ch := luax.Check[Character](L, 1)
			parseDebuffFlag := func(lv lua.LValue, argIndex int) (constant.DebuffFlag, bool) {
				tbl, ok := lv.(*lua.LTable)
				if !ok {
					L.ArgError(argIndex, "DebuffFlag table expected")
					return constant.DebuffFlag{}, false
				}
				maskLV := tbl.RawGetString("mask")
				posLV := tbl.RawGetString("position")
				if maskLV.Type() != lua.LTNumber || posLV.Type() != lua.LTNumber {
					L.ArgError(argIndex, "DebuffFlag table must have numeric mask and position")
					return constant.DebuffFlag{}, false
				}
				f := constant.DebuffFlag{
					Mask:     uint32(maskLV.(lua.LNumber)),
					Position: int(posLV.(lua.LNumber)),
				}
				if debuffLV := tbl.RawGetString("debuff"); debuffLV.Type() == lua.LTNumber {
					f.DiseaseSkillID = uint16(debuffLV.(lua.LNumber))
				}
				return f, true
			}
			flag, ok := parseDebuffFlag(L.Get(2), 2)
			if !ok {
				return 0
			}
			L.Push(lua.LBool(ch.debuffs.Has(flag)))
			return 1
		},
		"remove_debuff": func(L *lua.LState) int {
			ch := luax.Check[Character](L, 1)
			parseDebuffFlag := func(lv lua.LValue, argIndex int) (constant.DebuffFlag, bool) {
				tbl, ok := lv.(*lua.LTable)
				if !ok {
					L.ArgError(argIndex, "DebuffFlag table expected")
					return constant.DebuffFlag{}, false
				}
				maskLV := tbl.RawGetString("mask")
				posLV := tbl.RawGetString("position")
				if maskLV.Type() != lua.LTNumber || posLV.Type() != lua.LTNumber {
					L.ArgError(argIndex, "DebuffFlag table must have numeric mask and position")
					return constant.DebuffFlag{}, false
				}
				f := constant.DebuffFlag{
					Mask:     uint32(maskLV.(lua.LNumber)),
					Position: int(posLV.(lua.LNumber)),
				}
				if debuffLV := tbl.RawGetString("debuff"); debuffLV.Type() == lua.LTNumber {
					f.DiseaseSkillID = uint16(debuffLV.(lua.LNumber))
				}
				return f, true
			}
			var flags []constant.DebuffFlag
			for i := 2; i <= L.GetTop(); i++ {
				flag, ok := parseDebuffFlag(L.Get(i), i)
				if !ok {
					return 0
				}
				flags = append(flags, flag)
			}
			if len(flags) > 0 {
				ch.debuffs.Remove(flags...)
			}
			return 0
		},
		"morph": func(L *lua.LState) int {
			ch := luax.Check[Character](L, 1)
			if L.GetTop() == 1 {
				_, value, ok := ch.buffs.GetBuffValue(constant.BuffFlagMorph)
				if !ok {
					L.Push(lua.LNumber(0))
					return 1
				}
				L.Push(lua.LNumber(value))
				return 1
			}
			if L.GetTop() != 2 || lua.LVAsBool(L.Get(2)) {
				L.ArgError(2, "morph() reads the morph id; morph(false) clears it")
				return 0
			}
			ch.buffs.RemoveBuff([]constant.BuffFlag{constant.BuffFlagMorph})
			return 0
		},
	}
}
