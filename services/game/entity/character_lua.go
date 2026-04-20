package entity

import (
	"log"
	"time"

	"github.com/boyism80/fm/core/luax"
	"github.com/boyism80/fm/protocol/response"
	"github.com/boyism80/fm/services/game/constant"
	"github.com/boyism80/fm/services/game/wz"
	"github.com/boyism80/fm/types"
	lua "github.com/yuin/gopher-lua"
)

func luaValuesToInterfaces(L *lua.LState, from, to int) []interface{} {
	if from > to {
		return nil
	}
	out := make([]interface{}, 0, to-from+1)
	for i := from; i <= to; i++ {
		lv := L.Get(i)
		switch v := lv.(type) {
		case lua.LNumber:
			out = append(out, float64(v))
		case lua.LString:
			out = append(out, string(v))
		case lua.LBool:
			out = append(out, bool(v))
		case *lua.LNilType:
			out = append(out, nil)
		default:
			out = append(out, nil)
		}
	}
	return out
}

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
		"gender": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok {
				L.ArgError(1, "Character expected")
				return 0
			}
			if L.GetTop() != 1 {
				L.ArgError(2, "gender() is read-only")
				return 0
			}
			L.Push(lua.LNumber(ch.GetGender()))
			return 1
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
		"ability_point": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok {
				L.ArgError(1, "Character expected")
				return 0
			}
			argc := L.GetTop()
			switch argc {
			case 1:
				L.Push(lua.LNumber(ch.AbilityPoint))
				return 1
			case 2, 3:
				v := L.CheckInt(2)
				if v < 0 {
					v = 0
				}
				notify := argc == 2 || L.ToBool(3)
				ch.SetAbilityPoint(uint16(v), notify)
				return 0
			default:
				L.ArgError(2, "ability_point() requires 0, 1 or 2 arguments")
				return 0
			}
		},
		"skill_point": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok {
				L.ArgError(1, "Character expected")
				return 0
			}
			argc := L.GetTop()
			switch argc {
			case 1:
				L.Push(lua.LNumber(ch.SkillPoint))
				return 1
			case 2, 3:
				v := L.CheckInt(2)
				if v < 0 {
					v = 0
				}
				notify := argc == 2 || L.ToBool(3)
				ch.SetSkillPoint(uint16(v), notify)
				return 0
			default:
				L.ArgError(2, "skill_point() requires 0, 1 or 2 arguments")
				return 0
			}
		},
		"base_hp": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok {
				L.ArgError(1, "Character expected")
				return 0
			}
			argc := L.GetTop()
			switch argc {
			case 1:
				L.Push(lua.LNumber(ch.BaseHp))
				return 1
			case 2, 3:
				v := L.CheckInt(2)
				if v < 0 {
					v = 0
				}
				notify := argc == 2 || L.ToBool(3)
				ch.SetBaseHp(uint32(v), notify)
				return 0
			default:
				L.ArgError(2, "base_hp() requires 0, 1 or 2 arguments")
				return 0
			}
		},
		"base_mp": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok {
				L.ArgError(1, "Character expected")
				return 0
			}
			argc := L.GetTop()
			switch argc {
			case 1:
				L.Push(lua.LNumber(ch.BaseMp))
				return 1
			case 2, 3:
				v := L.CheckInt(2)
				if v < 0 {
					v = 0
				}
				notify := argc == 2 || L.ToBool(3)
				ch.SetBaseMp(uint32(v), notify)
				return 0
			default:
				L.ArgError(2, "base_mp() requires 0, 1 or 2 arguments")
				return 0
			}
		},
		"update_stats": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok {
				L.ArgError(1, "Character expected")
				return 0
			}
			tbl := L.CheckTable(2)
			stats := make(map[constant.Stat]int32)
			tbl.ForEach(func(k lua.LValue, v lua.LValue) {
				if v.Type() != lua.LTNumber {
					return
				}
				stat := constant.Stat(lua.LVAsNumber(v))
				if val, ok := ch.GetStatValue(stat); ok {
					stats[stat] = val
				}
			})
			if len(stats) > 0 {
				ch.Listener.OnUpdateStats(ch, stats, true)
			}
			return 0
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
				value := float64(L.CheckNumber(2))
				if value < 0 {
					L.ArgError(2, "exp() setter requires a non-negative value; compute new total in script")
					return 0
				}
				maxU := float64(^uint32(0))
				var newExp uint32
				if value > maxU {
					newExp = ^uint32(0)
				} else {
					newExp = uint32(value)
				}
				cur := ch.exp
				if newExp == cur {
					return 0
				}
				if newExp > cur {
					ch.AddExp(newExp - cur)
					return 0
				}
				ch.exp = newExp
				ch.Listener.OnUpdateStats(ch, map[constant.Stat]int32{
					constant.STAT_EXP: int32(ch.exp),
				}, false)
				return 0
			default:
				L.ArgError(2, "exp() getter: exp(); setter: exp(value)")
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
				if value < 0 {
					L.ArgError(2, "meso() setter requires a non-negative value; compute new total in script")
					return 0
				}
				if value <= 2147483647 {
					ch.SetMeso(int32(value))
				} else {
					ch.SetMeso(2147483647)
				}
				ch.Listener.OnMesoChanged(ch, ch.Meso)
				return 0
			default:
				L.ArgError(2, "meso() getter: meso(); setter: meso(value)")
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
			dontRecordHistory := false
			switch argc {
			case 4:
				dontRecordHistory = L.CheckBool(4)
				fallthrough
			case 3:
				highlight = L.CheckBool(3)
			}

			ch.Listener.OnChat(ch, message, highlight, dontRecordHistory)
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
				buff := ch.Buffs.GetEntity(flag)
				skillBuff, ok := buff.(*SkillBuff)
				if !ok || skillBuff == nil {
					L.Push(lua.LNil)
					return 1
				}
				L.Push(luax.NewLuable(L, skillBuff))
				return 1
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
					if ld := skillEntry.Wz.GetLevelData(skillEntry.Level()); ld != nil {
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
					ch.Buffs.AddBuff(skillEntry.Wz, duration, uint8(skillEntry.Level()), ch.GetID(), values)
				} else {
					ch.Buffs.AddItemBuff(consumeWz, duration, values)
				}
				return 0
			}
		},
		"ridding": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok {
				L.ArgError(1, "Character expected")
				return 0
			}
			skillUD := L.CheckUserData(2)
			skillEntry, ok := skillUD.Value.(*SkillEntry)
			if !ok || skillEntry == nil || skillEntry.Wz == nil {
				L.ArgError(2, "SkillEntry with Wz expected")
				return 0
			}
			mountID := int32(L.CheckInt(3))
			if mountID <= 0 {
				L.ArgError(3, "mount_id must be > 0")
				return 0
			}
			ch.Buffs.AddBuff(
				skillEntry.Wz,
				time.Duration(0),
				uint8(skillEntry.Level()),
				ch.GetID(),
				map[constant.BuffFlag]int32{constant.BuffFlagMonsterRiding: mountID},
			)
			return 0
		},
		"show_buff_effect": func(L *lua.LState) int {
			argc := L.GetTop()
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok {
				L.ArgError(1, "Character expected")
				return 0
			}

			skillUD := L.CheckUserData(2)
			skillEntry, ok := skillUD.Value.(*SkillEntry)
			if !ok || skillEntry == nil {
				L.ArgError(2, "SkillEntry expected")
				return 0
			}

			effectID := uint8(1)
			if argc >= 3 {
				effectID = uint8(L.CheckInt(3))
			}
			skillLevelInt := skillEntry.Level()
			if skillLevelInt < 0 {
				skillLevelInt = 0
			}
			skillLevel := uint8(skillLevelInt)

			ch.Listener.OnShowBuffEffect(ch, effectID, skillEntry.Wz.ID, skillLevel, nil)
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
				lv := L.Get(i)
				switch v := lv.(type) {
				case lua.LNumber:
					ch.Buffs.RemoveSkillBuff(uint32(L.CheckInt(i)))
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
				ch.Buffs.RemoveBuff(flags)
			}
			return 0
		},
		"summons": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok {
				L.ArgError(1, "Character expected")
				return 0
			}
			tbl := L.NewTable()
			idx := 1
			for _, s := range ch.GetSummons() {
				if s != nil {
					tbl.RawSetInt(idx, luax.NewLuable(L, s))
					idx++
				}
			}
			L.Push(tbl)
			return 1
		},
		"buff_value": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok {
				L.ArgError(1, "Character expected")
				return 0
			}

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
			_, currentValue, ok := ch.Buffs.GetBuffValue(flag)
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
				_, updated := ch.Buffs.SetBuffValue(flag, newValue)
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
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok {
				L.ArgError(1, "Character expected")
				return 0
			}

			skillID := uint32(L.CheckInt(2))
			skillLevel := uint8(L.CheckInt(3))
			durationMs := L.CheckInt(4)
			if durationMs < 0 {
				durationMs = 0
			}
			movType := constant.SummonMovementType(L.CheckInt(5))
			summonType := constant.SummonType(L.CheckInt(6))

			pos := ch.Position
			if L.GetTop() >= 7 {
				lv := L.Get(7)
				if lv != lua.LNil {
					posTbl, tableOK := lv.(*lua.LTable)
					if !tableOK {
						L.ArgError(7, "position table or nil expected")
						return 0
					}
					var x, y int16
					if lx := posTbl.RawGetInt(1); lx != lua.LNil {
						x = int16(lua.LVAsNumber(lx))
					} else if lx := posTbl.RawGetString("x"); lx != lua.LNil {
						x = int16(lua.LVAsNumber(lx))
					}
					if ly := posTbl.RawGetInt(2); ly != lua.LNil {
						y = int16(lua.LVAsNumber(ly))
					} else if ly := posTbl.RawGetString("y"); ly != lua.LNil {
						y = int16(lua.LVAsNumber(ly))
					}
					pos = types.Point[int16]{X: x, Y: y}
				}
			}

			duration := time.Duration(durationMs) * time.Millisecond
			s := ch.SpawnSummon(constant.SkillID(skillID), skillLevel, movType, summonType, pos, duration)
			if s == nil {
				L.Push(lua.LNil)
				return 1
			}
			L.Push(luax.NewLuable(L, s))
			return 1
		},
		"create_mist": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok {
				L.ArgError(1, "Character expected")
				return 0
			}
			skillUD := L.CheckUserData(2)
			skill, ok := skillUD.Value.(*SkillEntry)
			if !ok || skill == nil {
				L.ArgError(2, "Skill expected")
				return 0
			}
			durationMs := L.CheckInt(3)
			if durationMs < 0 {
				durationMs = 0
			}
			mistType := constant.MistType(L.CheckInt(4))
			boundsTable := L.CheckTable(5)
			left := int32(lua.LVAsNumber(boundsTable.RawGetString("left")))
			top := int32(lua.LVAsNumber(boundsTable.RawGetString("top")))
			right := int32(lua.LVAsNumber(boundsTable.RawGetString("right")))
			bottom := int32(lua.LVAsNumber(boundsTable.RawGetString("bottom")))
			duration := time.Duration(durationMs) * time.Millisecond
			initialDelay := time.Duration(0)
			if L.GetTop() >= 6 {
				initialDelayMs := L.CheckInt(6)
				if initialDelayMs > 0 {
					initialDelay = time.Duration(initialDelayMs) * time.Millisecond
				}
			}
			poisonTickMultiplier := 1.0
			if L.GetTop() >= 7 {
				poisonTickMultiplier = float64(L.CheckNumber(7))
				if poisonTickMultiplier <= 0 {
					poisonTickMultiplier = 1.0
				}
			}
			bounds := types.Rect[int32]{Left: left, Top: top, Right: right, Bottom: bottom}
			mist := ch.SpawnMist(skill, ch.Position, mistType, bounds, duration, initialDelay, poisonTickMultiplier)
			if mist == nil {
				L.Push(lua.LNil)
				return 1
			}
			L.Push(luax.NewLuable(L, mist))
			return 1
		},
		"create_door": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok {
				L.ArgError(1, "Character expected")
				return 0
			}
			skillID := uint32(L.CheckInt(2))
			ch.SpawnDoor(constant.SkillID(skillID))
			return 0
		},
		"remove_summon": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok {
				L.ArgError(1, "Character expected")
				return 0
			}
			skillID := uint32(L.CheckInt(2))
			for _, s := range ch.summons {
				if s != nil && s.SkillID == constant.SkillID(skillID) {
					ch.RemoveSummon(s, true)
					break
				}
			}
			return 0
		},
		"remove_door": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok {
				L.ArgError(1, "Character expected")
				return 0
			}
			skillID := uint32(L.CheckInt(2))
			ch.RemoveDoorBySkill(constant.SkillID(skillID), true)
			return 0
		},
		"clear_summons": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok {
				L.ArgError(1, "Character expected")
				return 0
			}
			ch.ClearSummons()
			return 0
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
			message := ""
			prev := false
			next := false
			switch argc {
			case 5:
				next = L.CheckBool(5)
				fallthrough
			case 4:
				prev = L.CheckBool(4)
				fallthrough
			case 3:
				message = L.CheckString(3)
				fallthrough
			case 2:
				npc = L.CheckInt(2)
			}

			ch.Listener.OnDialog(ch, uint32(npc), message, prev, next)
			ch.SetDialog(L)
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
			message := ""
			prev := false
			next := false
			switch argc {
			case 5:
				next = L.CheckBool(5)
				fallthrough
			case 4:
				prev = L.CheckBool(4)
				fallthrough
			case 3:
				message = L.CheckString(3)
				fallthrough
			case 2:
				npc = L.CheckInt(2)
			}

			ch.Listener.OnDialogYesNo(ch, uint32(npc), message, prev, next)
			ch.SetDialog(L)
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
			message := ""
			selections := []string{}
			switch argc {
			case 4:
				tbl := L.CheckTable(4)
				tbl.ForEach(func(_, value lua.LValue) {
					if str, ok := value.(lua.LString); ok {
						selections = append(selections, string(str))
					}
				})
				fallthrough
			case 3:
				message = L.CheckString(3)
				fallthrough
			case 2:
				npc = L.CheckInt(2)
			}

			ch.Listener.OnDialogList(ch, uint32(npc), message, selections)
			ch.SetDialog(L)
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
			message := ""
			enableEscape := false
			switch argc {
			case 4:
				enableEscape = L.CheckBool(4)
				fallthrough
			case 3:
				message = L.CheckString(3)
				fallthrough
			case 2:
				npc = L.CheckInt(2)
			}

			ch.Listener.OnDialogAccept(ch, uint32(npc), message, enableEscape)
			ch.SetDialog(L)
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
			message := ""
			switch argc {
			case 3:
				message = L.CheckString(3)
				fallthrough
			case 2:
				npc = L.CheckInt(2)
			}

			ch.Listener.OnDialogInput(ch, uint32(npc), message)
			ch.SetDialog(L)
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
			ch.Listener.OnMessage(ch, constant.MSG_LIGHT_BLUE_TEXT, message)
			return 0
		},
		"role": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok {
				L.ArgError(1, "Character expected")
				return 0
			}
			argc := L.GetTop()
			switch argc {
			case 1:
				L.Push(lua.LNumber(ch.Role))
				return 1
			case 2:
				ch.Role = constant.CharacterRole(L.CheckInt(2))
				return 0
			default:
				L.ArgError(2, "role() requires 0 or 1 arguments")
				return 0
			}
		},
		"invincible": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok {
				L.ArgError(1, "Character expected")
				return 0
			}
			argc := L.GetTop()
			switch argc {
			case 1:
				L.Push(lua.LBool(ch.GetInvincible()))
				return 1
			case 2:
				ch.SetInvincible(L.CheckBool(2))
				return 0
			default:
				L.ArgError(2, "invincible() requires 0 or 1 arguments")
				return 0
			}
		},
		"script": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok {
				L.ArgError(1, "Character expected")
				return 0
			}
			argc := L.GetTop()
			if argc < 3 {
				L.ArgError(2, "script(scriptPath, funcName, ...) requires at least 2 arguments")
				return 0
			}
			scriptPath := L.CheckString(2)
			funcName := L.CheckString(3)
			args := make([]interface{}, 0, argc-2)
			args = append(args, ch)
			for i := 4; i <= argc; i++ {
				if v, ok := luax.LValueToInterface(L.Get(i)); ok {
					args = append(args, v)
				}
			}
			cfg, ok := luax.GetConfiguration(L)
			if !ok || cfg.ActorContext == nil {
				L.RaiseError("script: thread has no actor PID (call from command context)")
				return 0
			}
			pid := cfg.ActorContext.Self()
			if pid == nil {
				L.RaiseError("script: thread has no actor PID (call from command context)")
				return 0
			}
			mapInstance := ch.GetMap()
			if mapInstance == nil {
				L.RaiseError("script: character map not found")
				return 0
			}
			root := mapInstance.GetLuaRoot()
			if root == nil {
				L.RaiseError("script: root lua state not found")
				return 0
			}
			result, thread, err := luax.Call(root, scriptPath, funcName, args...)
			if thread != nil {
				thread.Close()
			}
			if err != nil {
				L.RaiseError("script: %v", err)
				return 0
			}
			if result != nil {
				L.Push(result)
				return 1
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

			skillEntry := ch.Skills.Get(skillID)
			if skillEntry == nil {
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
			if existing := ch.Skills.Get(skillID); existing != nil {
				L.Push(luax.NewLuable(L, existing))
				return 1
			}

			if ch.GameWorld == nil {
				L.Push(lua.LNil)
				return 1
			}
			resources := ch.GameWorld.GetResources()
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

			skillEntry := NewSkillEntry(ch, wzSkill, 1, masterLevel)
			skillEntry.Expiration = time.Time{}
			ch.Skills.Register(skillID, skillEntry)

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
			switch L.GetTop() {
			case 1:

			default:
				L.ArgError(2, "skills() is read-only")
				return 0
			}

			skillsTable := L.NewTable()
			ch.Skills.ForEach(func(skillID uint32, skillEntry *SkillEntry) {
				skillUD := luax.NewLuable(L, skillEntry)
				skillsTable.RawSetInt(int(skillID), skillUD)
			})

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
		"chair": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok {
				L.ArgError(1, "Character expected")
				return 0
			}
			L.Push(lua.LNumber(ch.Chair))
			return 1
		},
		"stance": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok {
				L.ArgError(1, "Character expected")
				return 0
			}
			argc := L.GetTop()
			switch argc {
			case 1:
				L.Push(lua.LNumber(ch.Stance))
				return 1
			case 2:
				ch.Stance = uint8(L.CheckInt(2))
				return 0
			default:
				L.ArgError(2, "stance() requires 0 or 1 argument")
				return 0
			}
		},
		"stance_of": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok {
				L.ArgError(1, "Character expected")
				return 0
			}
			group := constant.StanceKind(L.CheckInt(2))
			members, ok := constant.StanceGroupByKind[group]
			if ok {
				for _, v := range members {
					if ch.Stance == v {
						L.Push(lua.LTrue)
						return 1
					}
				}
				L.Push(lua.LFalse)
				return 1
			}
			L.Push(lua.LBool(ch.Stance == uint8(group)))
			return 1
		},
		"class_of": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok {
				L.ArgError(1, "Character expected")
				return 0
			}
			classCode := uint16(L.CheckInt(2))
			L.Push(lua.LBool(ch.ClassOf(classCode)))
			return 1
		},
		"homing": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok || ch == nil {
				L.ArgError(1, "Character expected")
				return 0
			}
			if L.GetTop() != 1 {
				L.ArgError(2, "homing() takes no arguments")
				return 0
			}
			if ch.HomingTargetOID == nil {
				L.Push(lua.LNil)
				return 1
			}
			mapInstance := ch.GetMap()
			if mapInstance == nil {
				L.Push(lua.LNil)
				return 1
			}
			mob := mapInstance.GetMob(*ch.HomingTargetOID)
			if mob == nil {
				L.Push(lua.LNil)
				return 1
			}
			L.Push(luax.NewLuable(L, mob))
			return 1
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
		"bonus_watk": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok {
				L.ArgError(1, "Character expected")
				return 0
			}
			argc := L.GetTop()
			switch argc {
			case 1:
				L.Push(lua.LNumber(ch.BonusStats.Watk))
				return 1
			case 2:
				ch.BonusStats.Watk = int16(L.CheckInt(2))
				return 0
			default:
				L.ArgError(2, "bonus_watk() requires 0 or 1 arguments")
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
				L.Push(lua.LNumber(ch.GetBonusHp()))
				L.Push(lua.LNumber(ch.BonusStats.MaxHpPercent))
				return 2
			case 2:
				ch.SetBonusHp(int32(L.CheckInt(2)), true)
				return 0
			case 3:
				ch.SetBonusHp(int32(L.CheckInt(2)), false)
				ch.SetMaxHpPercent(int16(L.CheckInt(3)), true)
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
				L.Push(lua.LNumber(ch.GetBonusMp()))
				L.Push(lua.LNumber(ch.BonusStats.MaxMpPercent))
				return 2
			case 2:
				ch.SetBonusMp(int32(L.CheckInt(2)), true)
				return 0
			case 3:
				ch.SetBonusMp(int32(L.CheckInt(2)), false)
				ch.SetMaxMpPercent(int16(L.CheckInt(3)), true)
				return 0
			default:
				L.ArgError(2, "bonus_max_mp() requires 0, 1 or 2 arguments")
				return 0
			}
		},
		"bonus_max_hp_fixed": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok {
				L.ArgError(1, "Character expected")
				return 0
			}
			argc := L.GetTop()
			switch argc {
			case 1:
				L.Push(lua.LNumber(ch.GetBonusHp()))
				return 1
			case 2, 3:
				v := int32(L.CheckInt(2))
				notify := argc == 2 || L.ToBool(3)
				ch.SetBonusHp(v, notify)
				return 0
			default:
				L.ArgError(2, "bonus_max_hp_fixed() requires 0, 1 or 2 arguments")
				return 0
			}
		},
		"bonus_max_hp_ratio": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok {
				L.ArgError(1, "Character expected")
				return 0
			}
			argc := L.GetTop()
			switch argc {
			case 1:
				L.Push(lua.LNumber(ch.BonusStats.MaxHpPercent))
				return 1
			case 2, 3:
				p := int16(L.CheckInt(2))
				notify := argc == 2 || L.ToBool(3)
				ch.SetMaxHpPercent(p, notify)
				return 0
			default:
				L.ArgError(2, "bonus_max_hp_ratio() requires 0, 1 or 2 arguments")
				return 0
			}
		},
		"bonus_max_mp_fixed": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok {
				L.ArgError(1, "Character expected")
				return 0
			}
			argc := L.GetTop()
			switch argc {
			case 1:
				L.Push(lua.LNumber(ch.GetBonusMp()))
				return 1
			case 2, 3:
				v := int32(L.CheckInt(2))
				notify := argc == 2 || L.ToBool(3)
				ch.SetBonusMp(v, notify)
				return 0
			default:
				L.ArgError(2, "bonus_max_mp_fixed() requires 0, 1 or 2 arguments")
				return 0
			}
		},
		"bonus_max_mp_ratio": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok {
				L.ArgError(1, "Character expected")
				return 0
			}
			argc := L.GetTop()
			switch argc {
			case 1:
				L.Push(lua.LNumber(ch.BonusStats.MaxMpPercent))
				return 1
			case 2, 3:
				p := int16(L.CheckInt(2))
				notify := argc == 2 || L.ToBool(3)
				ch.SetMaxMpPercent(p, notify)
				return 0
			default:
				L.ArgError(2, "bonus_max_mp_ratio() requires 0, 1 or 2 arguments")
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
		"channel_drop_rate": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok {
				L.ArgError(1, "Character expected")
				return 0
			}
			if ch.GameWorld == nil {
				L.Push(lua.LNumber(1))
				return 1
			}
			L.Push(lua.LNumber(ch.GameWorld.GetDropRate()))
			return 1
		},
		"bonus_exp_rate": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok {
				L.ArgError(1, "Character expected")
				return 0
			}
			argc := L.GetTop()
			switch argc {
			case 1:
				L.Push(lua.LNumber(ch.BonusStats.ExpRate))
				return 1
			case 2:
				ch.BonusStats.ExpRate = int16(L.CheckInt(2))
				return 0
			default:
				L.ArgError(2, "bonus_exp_rate() requires 1 or 2 arguments")
				return 0
			}
		},
		"potion_heal_rate": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok {
				L.ArgError(1, "Character expected")
				return 0
			}
			argc := L.GetTop()
			switch argc {
			case 1:
				L.Push(lua.LNumber(ch.BonusStats.PotionHealRate))
				return 1
			case 2:
				ch.BonusStats.PotionHealRate = int16(L.CheckInt(2))
				return 0
			default:
				L.ArgError(2, "potion_heal_rate() requires 1 or 2 arguments")
				return 0
			}
		},
		"potion_duration_rate": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok {
				L.ArgError(1, "Character expected")
				return 0
			}
			argc := L.GetTop()
			switch argc {
			case 1:
				L.Push(lua.LNumber(ch.BonusStats.PotionDurationRate))
				return 1
			case 2:
				ch.BonusStats.PotionDurationRate = int16(L.CheckInt(2))
				return 0
			default:
				L.ArgError(2, "potion_duration_rate() requires 1 or 2 arguments")
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
			switch argc {
			case 1:
				m := ch.GetMap()
				if m == nil {
					L.Push(lua.LNil)
					return 1
				}
				L.Push(luax.NewLuable(L, m))
				return 1
			case 2, 3:
				var targetMap *Map
				arg2 := L.Get(2)
				switch v := arg2.(type) {
				case *lua.LUserData:
					var ok bool
					targetMap, ok = v.Value.(*Map)
					if !ok || targetMap == nil {
						L.ArgError(2, "Map or map name (string) expected")
						return 0
					}
				case lua.LString:
					if ch.GameWorld == nil {
						L.RaiseError("map: no context to resolve map name")
						return 0
					}
					resources := ch.GameWorld.GetResources()
					if resources == nil {
						L.RaiseError("map: no resources to resolve map name")
						return 0
					}
					mapId, ok := resources.NameToMap(string(v))
					if !ok {
						L.RaiseError("map: unknown map name %q", string(v))
						return 0
					}
					targetMap = ch.GameWorld.GetMap(mapId)
					if targetMap == nil {
						L.RaiseError("map: map %q (id %d) not found", string(v), mapId)
						return 0
					}
				default:
					L.ArgError(2, "Map or map name (string) expected")
					return 0
				}
				spawnPoint := uint8(1)
				if argc == 3 {
					spawnPoint = uint8(L.CheckInt(3))
				}
				if err := ch.Warp(targetMap, spawnPoint); err != nil {
					L.RaiseError("warp: %v", err)
					return 0
				}
				return 0
			default:
				L.ArgError(2, "map() getter: 0 args; setter: map or name (string), optional spawnPoint")
				return 0
			}
		},
		"show_magnet": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok {
				L.ArgError(1, "Character expected")
				return 0
			}
			mobUD := L.CheckUserData(2)
			if mobUD.Value == nil {
				L.ArgError(2, "Mob expected")
				return 0
			}
			mob, ok := mobUD.Value.(*Mob)
			if !ok {
				L.ArgError(2, "Mob expected")
				return 0
			}
			success := L.CheckBool(3)
			var successByte uint8
			if success {
				successByte = 1
			}
			m := ch.GetMap()
			if m != nil {
				ch.Broadcast(&response.ShowMagnet{MobID: mob.OID, Success: successByte}, &ObjectBroadcastOption{})
			}
			return 0
		},
		"debuff": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok {
				L.ArgError(1, "Character expected")
				return 0
			}
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
				L.ArgError(3, "duration_ms must be positive")
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
			ch.GiveDebuff(flag, time.Duration(durationMs)*time.Millisecond, x, skillID, skillLevel)
			return 0
		},
		"has_debuff": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok {
				L.ArgError(1, "Character expected")
				return 0
			}
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
			L.Push(lua.LBool(ch.HasDebuff(flag)))
			return 1
		},
		"remove_debuff": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok {
				L.ArgError(1, "Character expected")
				return 0
			}
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
				ch.RemoveDebuff(flags...)
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
			switch L.GetTop() {
			case 1:
				L.ArgError(2, "equipped(part) requires equipment part (e.g. EquipmentPart.Weapon)")
				return 0
			default:
				part := constant.EquipmentPartsType(L.CheckInt(2))
				if eq, ok := ch.Equipments[part]; ok && eq != nil {
					L.Push(luax.NewLuable(L, eq.(luax.Luable)))
					return 1
				}
				L.Push(lua.LNil)
				return 1
			}
		},
		"equip": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok {
				L.ArgError(1, "Character expected")
				return 0
			}
			if L.GetTop() < 2 {
				L.ArgError(2, "equip(slot | item | name) requires one argument")
				return 0
			}
			var slot int16
			switch lv := L.Get(2).(type) {
			case lua.LNumber:
				slot = int16(lv)
			case lua.LString:
				if ch.GameWorld == nil {
					L.Push(lua.LBool(false))
					return 1
				}
				resources := ch.GameWorld.GetResources()
				if resources == nil {
					L.Push(lua.LBool(false))
					return 1
				}
				itemID, ok := resources.NameToItem(string(lv))
				if !ok {
					L.Push(lua.LBool(false))
					return 1
				}
				invType, slots := ch.FindSlots(itemID)
				if invType != constant.INVENTORY_TYPE_EQUIPMENT || len(slots) == 0 {
					L.Push(lua.LBool(false))
					return 1
				}
				slot = slots[0]
			default:
				if itemUD, ok := L.Get(2).(*lua.LUserData); ok && itemUD.Value != nil {
					if item, ok := itemUD.Value.(Item); ok {
						var found bool
						slot, found = ch.FindSlot(constant.INVENTORY_TYPE_EQUIPMENT, item)
						if !found {
							L.Push(lua.LBool(false))
							return 1
						}
						break
					}
				}
				L.ArgError(2, "equip(slot | item | name): slot (number), item (equipment), or name (string) expected")
				return 0
			}
			err := ch.Equip(slot)
			L.Push(lua.LBool(err == nil))
			return 1
		},
		"unequip": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok {
				L.ArgError(1, "Character expected")
				return 0
			}
			if L.GetTop() < 2 {
				L.ArgError(2, "unequip(parts | item | name) requires one argument")
				return 0
			}
			var parts constant.EquipmentPartsType
			switch lv := L.Get(2).(type) {
			case lua.LNumber:
				parts = constant.EquipmentPartsType(lv)
			case lua.LString:
				if ch.GameWorld == nil {
					L.Push(lua.LBool(false))
					return 1
				}
				resources := ch.GameWorld.GetResources()
				if resources == nil {
					L.Push(lua.LBool(false))
					return 1
				}
				itemID, ok := resources.NameToItem(string(lv))
				if !ok {
					L.Push(lua.LBool(false))
					return 1
				}
				var found bool
				for p, eq := range ch.Equipments {
					if eq != nil && eq.GetModel().GetID() == itemID {
						parts = p
						found = true
						break
					}
				}
				if !found {
					L.Push(lua.LBool(false))
					return 1
				}
			default:
				if itemUD, ok := L.Get(2).(*lua.LUserData); ok && itemUD.Value != nil {
					if item, ok := itemUD.Value.(Item); ok {
						var found bool
						for p, eq := range ch.Equipments {
							if eq == item {
								parts = p
								found = true
								break
							}
						}
						if !found {
							L.Push(lua.LBool(false))
							return 1
						}
						break
					}
				}
				L.ArgError(2, "unequip(parts | item | name): parts (EquipmentPart), item (equipment), or name (string) expected")
				return 0
			}
			err := ch.Unequip(parts)
			L.Push(lua.LBool(err == nil))
			return 1
		},
		"item": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok {
				L.ArgError(1, "Character expected")
				return 0
			}
			switch L.GetTop() {
			case 2:
				var itemID uint32
				switch lv := L.Get(2).(type) {
				case lua.LNumber:
					itemID = uint32(lv)
				case lua.LString:
					if ch.GameWorld == nil {
						L.Push(L.NewTable())
						return 1
					}
					resources := ch.GameWorld.GetResources()
					if resources == nil {
						L.Push(L.NewTable())
						return 1
					}
					id, ok := resources.NameToItem(string(lv))
					if !ok {
						L.Push(L.NewTable())
						return 1
					}
					itemID = id
				default:
					L.ArgError(2, "item id (number) or item name (string) expected")
					return 0
				}
				invType := constant.GetInventoryTypeByItemID(itemID)
				slotTbl := L.NewTable()
				if inven := ch.Inventory[invType]; inven != nil {
					for slot, item := range inven.Items {
						if item != nil && item.GetModel().GetID() == itemID {
							L.Push(luax.NewLuable(L, item))
							slotTbl.RawSet(lua.LNumber(slot), L.Get(-1))
							L.Pop(1)
						}
					}
				}
				L.Push(slotTbl)
				return 1
			case 3:
				invType := constant.InventoryType(L.CheckInt(2))
				slot := int16(L.CheckInt(3))
				item := ch.GetItem(invType, slot)
				if item == nil {
					L.Push(lua.LNil)
				} else {
					L.Push(luax.NewLuable(L, item))
				}
				return 1
			default:
				L.ArgError(2, "item(item_id) or item(inven_type, slot) expected")
				return 0
			}
		},
		"items": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok {
				L.ArgError(1, "Character expected")
				return 0
			}
			argc := L.GetTop()
			switch argc {
			case 1:
				result := L.NewTable()
				for invType, inven := range ch.Inventory {
					if inven == nil {
						continue
					}
					invTbl := L.NewTable()
					for slot, item := range inven.Items {
						if item == nil {
							continue
						}
						L.Push(luax.NewLuable(L, item))
						invTbl.RawSet(lua.LNumber(slot), L.Get(-1))
						L.Pop(1)
					}
					result.RawSet(lua.LNumber(invType), invTbl)
				}
				L.Push(result)
				return 1
			default:
				invType := constant.InventoryType(L.CheckInt(2))
				slotTbl := L.NewTable()
				if inven := ch.Inventory[invType]; inven != nil {
					for slot, item := range inven.Items {
						if item == nil {
							continue
						}
						L.Push(luax.NewLuable(L, item))
						slotTbl.RawSet(lua.LNumber(slot), L.Get(-1))
						L.Pop(1)
					}
				}
				L.Push(slotTbl)
				return 1
			}
		},
		"mkitem": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok {
				L.ArgError(1, "Character expected")
				return 0
			}
			if ch.GameWorld == nil {
				L.Push(lua.LNil)
				return 1
			}
			resources := ch.GameWorld.GetResources()
			if resources == nil {
				L.Push(lua.LNil)
				return 1
			}

			argc := L.GetTop()
			count := uint16(1)
			switch argc {
			case 1:
				L.ArgError(2, "mkitem(itemIdOrName [, count]) requires at least one argument")
				return 0
			case 2:

			default:
				if n := L.CheckInt(3); n >= 1 {
					count = uint16(n)
				}
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

			_ = count
			model := resources.Items[itemId]
			if model == nil {
				L.Push(lua.LNil)
				return 1
			}
			item, err := NewItem(itemId, count, ch.GameWorld)
			if err != nil {
				L.Push(lua.LNil)
				return 1
			}

			added, err := ch.AddItem(item, true)
			if err != nil {
				L.Push(lua.LNil)
				return 1
			}

			L.Push(luax.NewLuable(L, added[0]))
			return 1
		},
		"item_wz": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok {
				L.ArgError(1, "Character expected")
				return 0
			}
			if ch.GameWorld == nil {
				L.Push(lua.LNil)
				return 1
			}
			resources := ch.GameWorld.GetResources()
			if resources == nil {
				L.Push(lua.LNil)
				return 1
			}

			argc := L.GetTop()
			count := uint16(1)
			switch argc {
			case 1:
				L.ArgError(2, "item_wz(itemIdOrName [, count]) requires at least one argument")
				return 0
			case 2:

			default:
				if n := L.CheckInt(3); n >= 1 {
					count = uint16(n)
				}
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

			item, err := NewItem(itemId, count, ch.GameWorld)
			if err != nil {
				L.Push(lua.LNil)
				return 1
			}

			consume, ok := item.(*Consume)
			if !ok || consume == nil {
				L.Push(lua.LNil)
				return 1
			}

			L.Push(luax.NewLuable(L, consume))
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
			switch argc {
			case 1:
				L.ArgError(2, "rmitem(itemIdOrName [, count]) or rmitem(inv_type, slot, count)")
				return 0
			case 4:
				invType := constant.InventoryType(L.CheckInt(2))
				slot := int16(L.CheckInt(3))
				count := uint16(L.CheckInt(4))
				if count == 0 {
					L.Push(lua.LBool(false))
					return 1
				}
				ok := ch.RemoveItem(invType, slot, count)
				L.Push(lua.LBool(ok))
				return 1
			case 2, 3:
				if ch.GameWorld == nil {
					L.Push(lua.LBool(false))
					return 1
				}
				resources := ch.GameWorld.GetResources()
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
				count := uint16(1)
				if argc == 3 {
					if n := L.CheckInt(3); n < 1 {
						L.Push(lua.LBool(false))
						return 1
					}
					count = uint16(L.CheckInt(3))
				}
				removed := ch.RemoveByItemIDCount(itemId, count)
				L.Push(lua.LBool(removed))
				return 1
			default:
				L.ArgError(2, "rmitem(itemIdOrName [, count]) or rmitem(inv_type, slot, count)")
				return 0
			}
		},
		"mktimer": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok {
				L.ArgError(1, "Character expected")
				return 0
			}
			argc := L.GetTop()
			if argc < 5 {
				L.ArgError(2, "mktimer(key, intervalMs, repeat, function [, ...args])")
				return 0
			}
			key := L.CheckString(2)
			intervalMs := L.CheckInt(3)
			if intervalMs < 1 {
				L.ArgError(3, "intervalMs must be positive")
				return 0
			}
			repeat := L.CheckBool(4)
			fn := L.CheckFunction(5)
			args := luaValuesToInterfaces(L, 6, argc)
			callback := func() {
				m := ch.GetMap()
				if m == nil {
					return
				}
				root := m.GetLuaRoot()
				if root == nil {
					return
				}
				fullArgs := make([]interface{}, 0, len(args)+1)
				fullArgs = append(fullArgs, ch)
				fullArgs = append(fullArgs, args...)
				if _, err := luax.CallFunction(root, fn, fullArgs...); err != nil {
					log.Printf("RunCharacterTimer %s: %v", key, err)
				}
			}
			added := ch.AddTimerWithCallback(key, time.Duration(intervalMs)*time.Millisecond, repeat, callback)
			L.Push(lua.LBool(added))
			return 1
		},
		"rmtimer": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok {
				L.ArgError(1, "Character expected")
				return 0
			}
			key := L.CheckString(2)
			removed := ch.RemoveTimer(key)
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
