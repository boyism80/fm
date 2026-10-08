package entity

import (
	"log"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/core/luax"
	pconst "github.com/boyism80/fm/protocol/constant"
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
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

func parseLuaExchangeSide(L *lua.LState, lv lua.LValue, argIndex int, nameToItem func(string) (uint32, bool)) (ExchangeSide, bool) {
	if lv == nil || lv == lua.LNil {
		return ExchangeSide{}, true
	}
	tbl, ok := lv.(*lua.LTable)
	if !ok {
		L.ArgError(argIndex, "exchange side table or nil expected")
		return ExchangeSide{}, false
	}

	side := ExchangeSide{}

	if itemLV := tbl.RawGetString("item"); itemLV != lua.LNil {
		itemTbl, ok := itemLV.(*lua.LTable)
		if !ok {
			L.ArgError(argIndex, "exchange side.item must be a table")
			return ExchangeSide{}, false
		}
		side.Items = make(map[uint32]uint16)
		failed := false
		itemTbl.ForEach(func(k, v lua.LValue) {
			if failed {
				return
			}
			countN, ok := v.(lua.LNumber)
			if !ok {
				L.ArgError(argIndex, "exchange side.item values must be numbers")
				failed = true
				return
			}
			if countN < 0 {
				L.ArgError(argIndex, "exchange side.item count must be non-negative")
				failed = true
				return
			}
			if countN == 0 {
				return
			}
			if countN > 65535 {
				L.ArgError(argIndex, "exchange side.item count exceeds uint16")
				failed = true
				return
			}
			count := uint16(countN)

			var itemID uint32
			switch key := k.(type) {
			case lua.LNumber:
				if key < 0 {
					L.ArgError(argIndex, "exchange side.item id must be non-negative")
					failed = true
					return
				}
				itemID = uint32(key)
			case lua.LString:
				if nameToItem == nil {
					L.ArgError(argIndex, "exchange item name lookup unavailable")
					failed = true
					return
				}
				id, ok := nameToItem(string(key))
				if !ok {
					L.ArgError(argIndex, "unknown item name in exchange side.item")
					failed = true
					return
				}
				itemID = id
			default:
				L.ArgError(argIndex, "exchange side.item keys must be item id or name")
				failed = true
				return
			}
			side.Items[itemID] += count
		})
		if failed {
			return ExchangeSide{}, false
		}
	}

	if mesoLV := tbl.RawGetString("meso"); mesoLV != lua.LNil {
		n, ok := mesoLV.(lua.LNumber)
		if !ok || n < 0 {
			L.ArgError(argIndex, "exchange side.meso must be a non-negative number")
			return ExchangeSide{}, false
		}
		if n > 2147483647 {
			L.ArgError(argIndex, "exchange side.meso exceeds int32")
			return ExchangeSide{}, false
		}
		side.Meso = int32(n)
	}

	if expLV := tbl.RawGetString("exp"); expLV != lua.LNil {
		n, ok := expLV.(lua.LNumber)
		if !ok || n < 0 {
			L.ArgError(argIndex, "exchange side.exp must be a non-negative number")
			return ExchangeSide{}, false
		}
		side.Exp = uint32(n)
	}

	if popLV := tbl.RawGetString("population"); popLV != lua.LNil {
		n, ok := popLV.(lua.LNumber)
		if !ok || n < 0 {
			L.ArgError(argIndex, "exchange side.population must be a non-negative number")
			return ExchangeSide{}, false
		}
		if n > 2147483647 {
			L.ArgError(argIndex, "exchange side.population exceeds int32")
			return ExchangeSide{}, false
		}
		side.Population = int32(n)
	}

	if randomLV := tbl.RawGetString("random"); randomLV != lua.LNil {
		side.Randomize = lua.LVAsBool(randomLV)
	}

	if bonusLV := tbl.RawGetString("bonus"); bonusLV != lua.LNil {
		bonusTbl, ok := bonusLV.(*lua.LTable)
		if !ok {
			L.ArgError(argIndex, "exchange side.bonus must be a table")
			return ExchangeSide{}, false
		}
		side.Bonus = map[string]int16{}
		bonusTbl.ForEach(func(key, value lua.LValue) {
			name, ok := key.(lua.LString)
			if !ok {
				return
			}
			n, ok := value.(lua.LNumber)
			if !ok {
				return
			}
			side.Bonus[string(name)] = int16(n)
		})
	}

	return side, true
}

func (ch *Character) LuaTypeName() string {
	return "LuaCharacter"
}

func (ch *Character) LuaBuiltinFuncs() map[string]lua.LGFunction {
	return map[string]lua.LGFunction{
		"__eq": func(L *lua.LState) int {
			leftUD := L.CheckUserData(1)
			rightUD := L.CheckUserData(2)
			leftCh, leftOk := leftUD.Value.(*Character)
			rightCh, rightOk := rightUD.Value.(*Character)
			if !leftOk || !rightOk || leftCh == nil || rightCh == nil {
				L.Push(lua.LFalse)
				return 1
			}
			L.Push(lua.LBool(leftCh.GetID() == rightCh.GetID()))
			return 1
		},
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
			if L.GetTop() == 1 {
				L.Push(lua.LNumber(ch.GetGender()))
				return 1
			}
			ch.SetGender(uint8(L.CheckInt(2)))
			return 0
		},
		"hair": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok {
				L.ArgError(1, "Character expected")
				return 0
			}
			if L.GetTop() == 1 {
				L.Push(lua.LNumber(ch.GetHair()))
				return 1
			}
			ch.SetHair(uint32(L.CheckInt(2)))
			return 0
		},
		"face": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok {
				L.ArgError(1, "Character expected")
				return 0
			}
			if L.GetTop() == 1 {
				L.Push(lua.LNumber(ch.GetFace()))
				return 1
			}
			ch.SetFace(uint32(L.CheckInt(2)))
			return 0
		},
		"skin": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok {
				L.ArgError(1, "Character expected")
				return 0
			}
			if L.GetTop() == 1 {
				L.Push(lua.LNumber(ch.GetSkinColor()))
				return 1
			}
			ch.SetSkin(uint8(L.CheckInt(2)))
			return 0
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
				oldLevel := ch.level
				ch.SetLevel(uint8(level))
				if ch.level > oldLevel && ch.Quests != nil {
					ch.Quests.RunAutoTriggers(nil, AutoQuestTriggerLevelUp, 0)
				}
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
				L.Push(lua.LNumber(ch.Points.AP))
				return 1
			case 2, 3:
				v := L.CheckInt(2)
				if v < 0 {
					v = 0
				}
				notify := argc == 2 || L.ToBool(3)
				ch.Points.SetAP(uint16(v), notify)
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
				L.Push(lua.LNumber(ch.Points.SP))
				return 1
			case 2, 3:
				v := L.CheckInt(2)
				if v < 0 {
					v = 0
				}
				notify := argc == 2 || L.ToBool(3)
				ch.Points.SetSP(uint16(v), notify)
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
					constant.StatEXP: int32(ch.exp),
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
				L.Push(lua.LNumber(ch.Inventory.Meso))
				return 1
			case 2:
				value := L.CheckNumber(2)
				if value < 0 {
					L.ArgError(2, "meso() setter requires a non-negative value; compute new total in script")
					return 0
				}
				if value <= 2147483647 {
					ch.Inventory.SetMeso(int32(value))
				} else {
					ch.Inventory.SetMeso(2147483647)
				}
				ch.Listener.OnMesoChanged(ch, ch.Inventory.Meso)
				return 0
			default:
				L.ArgError(2, "meso() getter: meso(); setter: meso(value)")
				return 0
			}
		},
		"population": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok {
				L.ArgError(1, "Character expected")
				return 0
			}

			argc := L.GetTop()
			switch argc {
			case 1:
				L.Push(lua.LNumber(ch.Stats.Population))
				return 1
			case 2:
				ch.Stats.setPopulation(int32(L.CheckInt(2)))
				return 0
			default:
				L.ArgError(2, "population() getter: population(); setter: population(value)")
				return 0
			}
		},
		"exchange": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok || ch == nil {
				L.ArgError(1, "Character expected")
				return 0
			}
			argc := L.GetTop()
			if argc < 2 || argc > 3 {
				L.ArgError(2, "exchange(cost [, reward]) requires 1 or 2 side tables")
				return 0
			}

			var nameToItem func(string) (uint32, bool)
			if ch.GameWorld != nil {
				if resources := ch.GameWorld.GetResources(); resources != nil {
					nameToItem = resources.NameToItem
				}
			}

			cost, ok := parseLuaExchangeSide(L, L.Get(2), 2, nameToItem)
			if !ok {
				return 0
			}
			var reward ExchangeSide
			if argc >= 3 {
				reward, ok = parseLuaExchangeSide(L, L.Get(3), 3, nameToItem)
				if !ok {
					return 0
				}
			}

			result := ch.Exchange(ExchangeSpec{Cost: cost, Reward: reward})
			L.Push(lua.LNumber(result))
			return 1
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
					ch.Buffs.AddBuff(skillEntry.Wz, duration, uint8(skillEntry.Level()), ch.GetID(), values, true)
				} else {
					ch.Buffs.AddItemBuff(consumeWz, duration, values, true, true)
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
				true,
			)
			return 0
		},
		"show_skill_effect": func(L *lua.LState) int {
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

			skillLevelInt := skillEntry.Level()
			if skillLevelInt < 0 {
				skillLevelInt = 0
			}
			skillLevel := uint8(skillLevelInt)
			if argc >= 3 {
				effectType := L.CheckInt(3)
				if effectType == int(pconst.SkillEffectTypeCast) || effectType == int(pconst.SkillEffectTypeAffected) {
					ch.Listener.OnShowSelfSkillEffect(ch, pconst.SkillEffectType(effectType), skillEntry.Wz.ID, skillLevel, nil)
					return 0
				}
				L.ArgError(3, "show_skill_effect() supports SkillEffectType.Cast or SkillEffectType.Affected")
				return 0
			}

			ch.Listener.OnShowSelfSkillEffect(ch, pconst.SkillEffectTypeCast, skillEntry.Wz.ID, skillLevel, nil)
			return 0
		},
		"show_effect": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok {
				L.ArgError(1, "Character expected")
				return 0
			}
			effectType := response.EffectType(L.CheckInt(2))
			switch effectType {
			case response.EffectTypeLevelUp,
				response.EffectTypeClassChange,
				response.EffectTypeQuestCompletion,
				response.EffectTypeRegisterCard,
				response.EffectTypeItemLevelUp:
				ch.Listener.OnShowSelfEffect(ch, effectType)
				return 0
			default:
				L.ArgError(2, "show_effect() supports EffectType.LevelUp, ClassChange, QuestCompletion, RegisterCard, or ItemLevelUp")
				return 0
			}
		},
		"show_quest_completion": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok || ch == nil {
				L.ArgError(1, "Character expected")
				return 0
			}
			if L.GetTop() != 2 {
				L.ArgError(2, "show_quest_completion(quest_id) requires quest id")
				return 0
			}
			questID := uint32(L.CheckInt(2))
			if ch.Listener != nil {
				ch.Listener.OnShowQuestCompletion(ch, questID)
			}
			return 0
		},
		"play_sound": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok || ch == nil {
				L.ArgError(1, "Character expected")
				return 0
			}
			argc := L.GetTop()
			if argc < 2 || argc > 3 {
				L.ArgError(2, "play_sound(name [, broadcast])")
				return 0
			}
			sound := L.CheckString(2)
			broadcast := false
			if argc >= 3 {
				broadcast = L.CheckBool(3)
			}
			if ch.Listener != nil {
				ch.Listener.OnPlaySound(ch, sound, broadcast)
			}
			return 0
		},
		"play_portal_sound": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok || ch == nil {
				L.ArgError(1, "Character expected")
				return 0
			}
			if L.GetTop() != 1 {
				L.ArgError(2, "play_portal_sound() takes no arguments")
				return 0
			}
			if ch.Listener != nil {
				ch.Listener.OnPlayPortalSound(ch)
			}
			return 0
		},
		"show_dragon_blood_effect": func(L *lua.LState) int {
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
			skillLevelInt := skillEntry.Level()
			if skillLevelInt < 0 {
				skillLevelInt = 0
			}
			ch.Listener.OnShowSelfDragonBloodEffect(ch, skillEntry.Wz.ID, uint8(skillLevelInt))
			return 0
		},
		"show_hp_healed_effect": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok {
				L.ArgError(1, "Character expected")
				return 0
			}
			amount := int32(L.CheckInt(2))
			ch.Listener.OnShowSelfHPHealedEffect(ch, amount)
			return 0
		},
		"show_reward_item_animation": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok {
				L.ArgError(1, "Character expected")
				return 0
			}
			itemID := uint32(L.CheckInt(2))
			effectText := L.CheckString(3)
			ch.Listener.OnShowSelfRewardItemAnimation(ch, itemID, effectText)
			return 0
		},
		"show_item_maker_success_effect": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok {
				L.ArgError(1, "Character expected")
				return 0
			}
			ch.Listener.OnShowSelfItemMakerSuccessEffect(ch)
			return 0
		},
		"show_crafting_effect": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok {
				L.ArgError(1, "Character expected")
				return 0
			}
			effectText := L.CheckString(2)
			timeValue := int32(L.CheckInt(3))
			modeValue := int32(L.CheckInt(4))
			ch.Listener.OnShowSelfCraftingEffect(ch, effectText, timeValue, modeValue)
			return 0
		},
		"show_dice_effect": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok {
				L.ArgError(1, "Character expected")
				return 0
			}
			effectID := int32(L.CheckInt(2))
			skillUD := L.CheckUserData(3)
			skillEntry, ok := skillUD.Value.(*SkillEntry)
			if !ok || skillEntry == nil || skillEntry.Wz == nil {
				L.ArgError(3, "SkillEntry with Wz expected")
				return 0
			}
			skillLevelInt := skillEntry.Level()
			if skillLevelInt < 0 {
				skillLevelInt = 0
			}
			ch.Listener.OnShowSelfDiceEffect(ch, effectID, skillEntry.Wz.ID, uint8(skillLevelInt))
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
					ch.Buffs.CancelBySource(int32(L.CheckInt(i)))
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
		"buffs": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok || ch == nil {
				L.ArgError(1, "Character expected")
				return 0
			}

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
			for _, entity := range ch.Buffs.Entities() {
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
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok {
				L.ArgError(1, "Character expected")
				return 0
			}
			tbl := L.NewTable()
			idx := 1
			for _, s := range ch.Summons.All() {
				tbl.RawSetInt(idx, luax.NewLuable(L, s))
				idx++
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
			s := ch.Summons.Spawn(constant.SkillID(skillID), skillLevel, movType, summonType, pos, duration)
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
			buff, ok := L.CheckUserData(2).Value.(*SkillBuff)
			if !ok {
				L.ArgError(2, "SkillBuff expected")
				return 0
			}
			ch.Doors.Spawn(buff)
			return 0
		},
		"remove_summon": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok {
				L.ArgError(1, "Character expected")
				return 0
			}
			ch.Summons.Remove(ch.Summons.Get(constant.SkillID(L.CheckInt(2))), true)
			return 0
		},
		"remove_door": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok {
				L.ArgError(1, "Character expected")
				return 0
			}
			buff, ok := L.CheckUserData(2).Value.(*SkillBuff)
			if !ok {
				L.ArgError(2, "SkillBuff expected")
				return 0
			}
			ch.Doors.Remove(buff, true)
			return 0
		},
		"clear_summons": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok {
				L.ArgError(1, "Character expected")
				return 0
			}
			ch.Summons.Clear()
			return 0
		},
		"dialog": func(L *lua.LState) int {
			argc := L.GetTop()
			if argc == 1 {
				switch L.Get(1).Type() {
				case lua.LTBool:
					b := lua.LVAsBool(L.Get(1))
					L.SetTop(0)
					L.Push(lua.LBool(b))
					return 1
				case lua.LTNil:
					L.SetTop(0)
					L.Push(lua.LBool(false))
					return 1
				}
			}
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
				switch npcArg := L.Get(2).(type) {
				case lua.LNumber:
					npc = int(npcArg)
				case *lua.LUserData:
					if n, ok := npcArg.Value.(*Npc); ok && n.Wz != nil && n.Wz.BaseSpawn != nil {
						npc = int(n.Wz.BaseSpawn.ID)
					}
				}
			}
			ch.Listener.OnDialog(ch, uint32(npc), message, prev, next)
			ch.Dialog.Ask(L, constant.DialogTypeDefault, 0)
			return L.Yield()
		},
		"dialog_yes_no": func(L *lua.LState) int {
			argc := L.GetTop()
			if argc == 1 {
				switch L.Get(1).Type() {
				case lua.LTBool:
					b := lua.LVAsBool(L.Get(1))
					L.SetTop(0)
					L.Push(lua.LBool(b))
					return 1
				case lua.LTNil:
					L.SetTop(0)
					L.Push(lua.LBool(false))
					return 1
				}
			}
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
				switch npcArg := L.Get(2).(type) {
				case lua.LNumber:
					npc = int(npcArg)
				case *lua.LUserData:
					if n, ok := npcArg.Value.(*Npc); ok && n.Wz != nil && n.Wz.BaseSpawn != nil {
						npc = int(n.Wz.BaseSpawn.ID)
					}
				}
			}
			ch.Listener.OnDialogYesNo(ch, uint32(npc), message, prev, next)
			ch.Dialog.Ask(L, constant.DialogTypeYesNo, 0)
			return L.Yield()
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
				switch npcArg := L.Get(2).(type) {
				case lua.LNumber:
					npc = int(npcArg)
				case *lua.LUserData:
					if n, ok := npcArg.Value.(*Npc); ok && n.Wz != nil && n.Wz.BaseSpawn != nil {
						npc = int(n.Wz.BaseSpawn.ID)
					}
				}
			}

			ch.Listener.OnDialogList(ch, uint32(npc), message, selections)
			ch.Dialog.Ask(L, constant.DialogTypeList, len(selections))
			return L.Yield(lua.LNumber(0))
		},
		"dialog_style": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok {
				L.ArgError(1, "Character expected")
				return 0
			}
			if L.GetTop() != 4 {
				L.ArgError(2, "dialog_style(npc, text, styles) requires npc, text, and style ids")
				return 0
			}
			npc := 0
			switch npcArg := L.Get(2).(type) {
			case lua.LNumber:
				npc = int(npcArg)
			case *lua.LUserData:
				if n, ok := npcArg.Value.(*Npc); ok && n.Wz != nil && n.Wz.BaseSpawn != nil {
					npc = int(n.Wz.BaseSpawn.ID)
				}
			}
			styles := []uint32{}
			L.CheckTable(4).ForEach(func(_, value lua.LValue) {
				if n, ok := value.(lua.LNumber); ok {
					styles = append(styles, uint32(n))
				}
			})
			ch.Listener.OnDialogStyle(ch, uint32(npc), L.CheckString(3), styles)
			ch.Dialog.Ask(L, constant.DialogTypeStyle, len(styles))
			return L.Yield(lua.LNumber(0))
		},
		"dialog_accept": func(L *lua.LState) int {
			argc := L.GetTop()
			if argc == 1 {
				switch L.Get(1).Type() {
				case lua.LTBool:
					b := lua.LVAsBool(L.Get(1))
					L.SetTop(0)
					L.Push(lua.LBool(b))
					return 1
				case lua.LTNil:
					L.SetTop(0)
					L.Push(lua.LBool(false))
					return 1
				}
			}
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
				switch npcArg := L.Get(2).(type) {
				case lua.LNumber:
					npc = int(npcArg)
				case *lua.LUserData:
					if n, ok := npcArg.Value.(*Npc); ok && n.Wz != nil && n.Wz.BaseSpawn != nil {
						npc = int(n.Wz.BaseSpawn.ID)
					}
				}
			}
			ch.Listener.OnDialogAccept(ch, uint32(npc), message, enableEscape)
			if enableEscape {
				ch.Dialog.Ask(L, constant.DialogTypeAcceptEscape, 0)
			} else {
				ch.Dialog.Ask(L, constant.DialogTypeAccept, 0)
			}
			return L.Yield()
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
				switch npcArg := L.Get(2).(type) {
				case lua.LNumber:
					npc = int(npcArg)
				case *lua.LUserData:
					if n, ok := npcArg.Value.(*Npc); ok && n.Wz != nil && n.Wz.BaseSpawn != nil {
						npc = int(n.Wz.BaseSpawn.ID)
					}
				}
			}

			ch.Listener.OnDialogInput(ch, uint32(npc), message)
			ch.Dialog.Ask(L, constant.DialogTypeInput, 0)
			return L.Yield(lua.LNumber(0))
		},
		"message": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok {
				L.ArgError(1, "Character expected")
				return 0
			}
			text := L.CheckString(2)
			msgType := constant.MsgLightBlueText
			scope := constant.MessageScopeSelf
			ear := false
			top := L.GetTop()
			if top >= 3 && L.Get(3) != lua.LNil {
				msgType = constant.ServerMessageType(L.CheckInt(3))
			}
			if top >= 4 && L.Get(4) != lua.LNil {
				scope = constant.MessageScope(L.CheckInt(4))
			}
			if top >= 5 {
				ear = L.ToBool(5)
			}
			switch scope {
			case constant.MessageScopeMap:
				ch.MapMessage(msgType, text)
			case constant.MessageScopeChannel:
				if ch.GameWorld != nil {
					ch.GameWorld.BroadcastNotice(msgType, text, 0, false)
				}
			case constant.MessageScopeWorld:
				cfg, ok := luax.GetConfiguration(L)
				if !ok || cfg.ActorContext == nil || ch.GameWorld == nil {
					L.Push(lua.LFalse)
					return 1
				}
				promise := ch.Listener.BroadcastNoticeAsync(cfg.ActorContext, ch, msgType, text, ear)
				if promise == nil {
					L.Push(lua.LFalse)
					return 1
				}
				return LuaYieldPromise(L, ch.GameWorld, promise, func(_ *internal.BroadcastNoticeReply, err error) []lua.LValue {
					if err != nil {
						return []lua.LValue{lua.LFalse}
					}
					return []lua.LValue{lua.LTrue}
				})
			default:
				ch.Listener.OnMessage(ch, msgType, text)
			}
			L.Push(lua.LTrue)
			return 1
		},
		"marriage": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok {
				L.ArgError(1, "Character expected")
				return 0
			}
			if ch.Wedding.Marriage == nil {
				L.Push(lua.LNil)
				return 1
			}
			L.Push(luax.NewLuable(L, ch.Wedding.Marriage))
			return 1
		},
		"reserve_wedding": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok {
				L.ArgError(1, "Character expected")
				return 0
			}
			cfg, ok := luax.GetConfiguration(L)
			if !ok || cfg.ActorContext == nil {
				L.Push(lua.LFalse)
				return 1
			}
			promise, err := ch.Wedding.Reserve(cfg.ActorContext, uint32(L.CheckNumber(2)))
			if err != nil {
				L.Push(lua.LFalse)
				return 1
			}
			return LuaYieldPromise(L, ch.GameWorld, promise, func(result bool, err error) []lua.LValue {
				if err != nil {
					return []lua.LValue{lua.LFalse}
				}
				return []lua.LValue{lua.LBool(result)}
			})
		},
		"break_engagement": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok {
				L.ArgError(1, "Character expected")
				return 0
			}
			cfg, ok := luax.GetConfiguration(L)
			if !ok || cfg.ActorContext == nil {
				L.Push(lua.LFalse)
				return 1
			}
			promise, err := ch.Wedding.BreakEngagement(cfg.ActorContext)
			if err != nil {
				L.Push(lua.LFalse)
				return 1
			}
			return LuaYieldPromise(L, ch.GameWorld, promise, func(result bool, err error) []lua.LValue {
				if err != nil {
					return []lua.LValue{lua.LFalse}
				}
				return []lua.LValue{lua.LBool(result)}
			})
		},
		"request_divorce": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok {
				L.ArgError(1, "Character expected")
				return 0
			}
			cfg, ok := luax.GetConfiguration(L)
			if !ok || cfg.ActorContext == nil {
				L.Push(lua.LNil)
				return 1
			}
			promise, err := ch.Wedding.RequestDivorce(cfg.ActorContext)
			if err != nil {
				L.Push(lua.LNil)
				return 1
			}
			return LuaYieldPromise(L, ch.GameWorld, promise, func(result string, err error) []lua.LValue {
				if err != nil {
					return []lua.LValue{lua.LNil}
				}
				return []lua.LValue{lua.LString(result)}
			})
		},
		"cancel_divorce": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok {
				L.ArgError(1, "Character expected")
				return 0
			}
			cfg, ok := luax.GetConfiguration(L)
			if !ok || cfg.ActorContext == nil {
				L.Push(lua.LFalse)
				return 1
			}
			promise, err := ch.Wedding.CancelDivorce(cfg.ActorContext)
			if err != nil {
				L.Push(lua.LFalse)
				return 1
			}
			return LuaYieldPromise(L, ch.GameWorld, promise, func(canceled bool, err error) []lua.LValue {
				return []lua.LValue{lua.LBool(err == nil && canceled)}
			})
		},
		"invited_to": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok {
				L.ArgError(1, "Character expected")
				return 0
			}
			L.Push(lua.LBool(ch.Wedding.InvitedTo(uint32(L.CheckNumber(2)))))
			return 1
		},
		"open_wedding_wishlist": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok {
				L.ArgError(1, "Character expected")
				return 0
			}
			ch.Listener.OnWeddingWishlistInput(ch)
			return 0
		},
		"open_wedding_gift": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok {
				L.ArgError(1, "Character expected")
				return 0
			}
			receiverID := uint32(L.CheckNumber(2))
			wishes := make([]string, 0)
			if tbl, ok := L.Get(3).(*lua.LTable); ok {
				tbl.ForEach(func(_, value lua.LValue) {
					wishes = append(wishes, value.String())
				})
			}
			ch.Wedding.OpenGift(receiverID, wishes)
			return 0
		},
		"open_wedding_gift_box": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok {
				L.ArgError(1, "Character expected")
				return 0
			}
			cfg, ok := luax.GetConfiguration(L)
			if !ok || cfg.ActorContext == nil {
				L.Push(lua.LFalse)
				return 1
			}
			return LuaYieldPromise(L, ch.GameWorld, ch.Wedding.OpenGiftBox(cfg.ActorContext), func(result bool, err error) []lua.LValue {
				if err != nil {
					return []lua.LValue{lua.LFalse}
				}
				return []lua.LValue{lua.LBool(result)}
			})
		},
		"open_storage": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if ok == false {
				L.ArgError(1, "Character expected")
				return 0
			}
			npcID := uint32(L.CheckInt(2))
			storeFee := int32(L.OptInt(3, 100))
			takeOutFee := int32(L.OptInt(4, 0))
			cfg, ok := luax.GetConfiguration(L)
			if ok == false || cfg.ActorContext == nil {
				L.ArgError(1, "actor context not available")
				return 0
			}
			ch.OpenStorage(cfg.ActorContext, npcID, storeFee, takeOutFee)
			return 0
		},
		"open_store_bank": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if ok == false {
				L.ArgError(1, "Character expected")
				return 0
			}
			npcID := uint32(L.CheckInt(2))
			cfg, ok := luax.GetConfiguration(L)
			if ok == false || cfg.ActorContext == nil {
				L.ArgError(1, "actor context not available")
				return 0
			}
			ch.StoreBank.Open(cfg.ActorContext, npcID)
			return 0
		},
		"open_duey": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if ok == false {
				L.ArgError(1, "Character expected")
				return 0
			}
			cfg, ok := luax.GetConfiguration(L)
			if ok == false || cfg.ActorContext == nil {
				L.ArgError(1, "actor context not available")
				return 0
			}
			ch.Duey.Open(cfg.ActorContext, false)
			return 0
		},
		"pet_id": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if ok == false {
				L.ArgError(1, "Character expected")
				return 0
			}
			if ch.Pets.Active == nil {
				L.Push(lua.LNumber(0))
				return 1
			}
			L.Push(lua.LNumber(ch.Pets.Active.Item.GetModel().GetID()))
			return 1
		},
		"pet_closeness": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if ok == false {
				L.ArgError(1, "Character expected")
				return 0
			}
			if ch.Pets.Active == nil {
				L.Push(lua.LNumber(0))
				return 1
			}
			L.Push(lua.LNumber(ch.Pets.Active.Item.Closeness))
			return 1
		},
		"add_pet_closeness": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if ok == false {
				L.ArgError(1, "Character expected")
				return 0
			}
			if ch.Pets.Active == nil {
				L.Push(lua.LFalse)
				return 1
			}
			ch.Pets.Active.AddCloseness(L.CheckInt(2))
			L.Push(lua.LTrue)
			return 1
		},
		"rename_pet": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if ok == false {
				L.ArgError(1, "Character expected")
				return 0
			}
			if ch.Pets.Active == nil {
				L.Push(lua.LFalse)
				return 1
			}
			L.Push(lua.LBool(ch.Pets.Active.Rename(L.CheckString(2)) == nil))
			return 1
		},
		"feed_pet_cash": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if ok == false {
				L.ArgError(1, "Character expected")
				return 0
			}
			if ch.Pets.Active == nil {
				L.Push(lua.LFalse)
				return 1
			}
			L.Push(lua.LBool(ch.Pets.Active.FeedCash(uint32(L.CheckInt(2))) == nil))
			return 1
		},
		"change_pet_skill": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if ok == false {
				L.ArgError(1, "Character expected")
				return 0
			}
			sn := uint64(L.CheckNumber(2))
			L.Push(lua.LBool(ch.Pets.ChangeSkill(sn, uint32(L.CheckInt(3))) == nil))
			return 1
		},
		"expired_pets": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if ok == false {
				L.ArgError(1, "Character expected")
				return 0
			}
			tbl := L.NewTable()
			for slot, pet := range ch.Pets.Expired() {
				entry := L.NewTable()
				entry.RawSetString("slot", lua.LNumber(slot))
				entry.RawSetString("id", lua.LNumber(pet.GetModel().GetID()))
				entry.RawSetString("name", lua.LString(pet.Name))
				tbl.Append(entry)
			}
			L.Push(tbl)
			return 1
		},
		"revive_pet": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if ok == false {
				L.ArgError(1, "Character expected")
				return 0
			}
			L.Push(lua.LBool(ch.Pets.Revive(int16(L.CheckInt(2))) == nil))
			return 1
		},
		"open_quick_delivery": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if ok == false {
				L.ArgError(1, "Character expected")
				return 0
			}
			L.Push(lua.LBool(ch.Duey.OpenQuick() == nil))
			return 1
		},
		"send_parcel": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if ok == false {
				L.ArgError(1, "Character expected")
				return 0
			}
			recipient := L.CheckString(2)
			senderName := L.CheckString(3)
			itemID := uint32(L.OptInt(4, 0))
			count := uint16(L.OptInt(5, 1))
			meso := int32(L.OptInt(6, 0))
			message := L.OptString(7, "")
			cfg, ok := luax.GetConfiguration(L)
			if ok == false || cfg.ActorContext == nil {
				L.ArgError(1, "actor context not available")
				return 0
			}
			L.Push(lua.LBool(ch.Duey.SendFromSystem(cfg.ActorContext, recipient, senderName, itemID, count, meso, message) == nil))
			return 1
		},
		"add_cash": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if ok == false {
				L.ArgError(1, "Character expected")
				return 0
			}
			nxCash := int32(L.CheckInt(2))
			maplePoint := int32(L.OptInt(3, 0))
			cfg, ok := luax.GetConfiguration(L)
			if ok == false || cfg.ActorContext == nil {
				L.ArgError(1, "actor context not available")
				return 0
			}
			ch.AddCash(cfg.ActorContext, nxCash, maplePoint)
			return 0
		},
		"create_cash_coupons": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if ok == false {
				L.ArgError(1, "Character expected")
				return 0
			}
			kind := internal.CashCouponKind(L.CheckInt(2))
			value := uint32(L.CheckInt(3))
			count := uint32(L.OptInt(4, 1))
			cfg, ok := luax.GetConfiguration(L)
			if ok == false || cfg.ActorContext == nil {
				L.ArgError(1, "actor context not available")
				return 0
			}
			ch.CreateCashCoupons(cfg.ActorContext, kind, value, count)
			return 0
		},
		"show_instruction": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok {
				L.ArgError(1, "Character expected")
				return 0
			}
			if L.GetTop() < 2 || L.GetTop() > 4 {
				L.ArgError(2, "show_instruction(text[, width[, height]])")
				return 0
			}
			text := L.CheckString(2)
			width := uint16(0)
			height := uint16(0)
			if L.GetTop() >= 3 {
				width = uint16(L.CheckInt(3))
			}
			if L.GetTop() >= 4 {
				height = uint16(L.CheckInt(4))
			}
			ch.Send(&response.Hint{
				Text:   text,
				Width:  width,
				Height: height,
			}, types.SEND_POLICY_ENCRYPT)
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
		"dojo_energy": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok {
				L.ArgError(1, "Character expected")
				return 0
			}
			switch L.GetTop() {
			case 1:
				L.Push(lua.LNumber(ch.DojoEnergy()))
				return 1
			case 2:
				ch.SetDojoEnergy(L.CheckInt(2))
				return 0
			default:
				L.ArgError(2, "dojo_energy() requires 0 or 1 arguments")
				return 0
			}
		},
		"instant_kill": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok {
				L.ArgError(1, "Character expected")
				return 0
			}
			argc := L.GetTop()
			switch argc {
			case 1:
				L.Push(lua.LBool(ch.GM.InstantKill))
				return 1
			case 2:
				ch.GM.InstantKill = L.CheckBool(2)
				return 0
			default:
				L.ArgError(2, "instant_kill() requires 0 or 1 arguments")
				return 0
			}
		},
		"player_mode": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok {
				L.ArgError(1, "Character expected")
				return 0
			}
			argc := L.GetTop()
			switch argc {
			case 1:
				L.Push(lua.LBool(ch.GM.PlayerMode))
				return 1
			case 2:
				ch.GM.PlayerMode = L.CheckBool(2)
				return 0
			default:
				L.ArgError(2, "player_mode() requires 0 or 1 arguments")
				return 0
			}
		},
		"timer_limit": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok {
				L.ArgError(1, "Character expected")
				return 0
			}
			argc := L.GetTop()
			switch argc {
			case 1:
				L.Push(lua.LNumber(ch.GM.TimerLimit))
				return 1
			case 2:
				ch.GM.TimerLimit = uint32(L.CheckNumber(2))
				return 0
			default:
				L.ArgError(2, "timer_limit() requires 0 or 1 arguments")
				return 0
			}
		},
		"open_npc": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok {
				L.ArgError(1, "Character expected")
				return 0
			}
			if L.GetTop() != 2 {
				L.ArgError(2, "open_npc(npc) requires npc")
				return 0
			}
			cfg, ok := luax.GetConfiguration(L)
			if !ok || cfg.ActorContext == nil {
				L.RaiseError("open_npc: thread has no actor context")
				return 0
			}
			var npc *Npc
			if ud, ok := L.Get(2).(*lua.LUserData); ok {
				npc, _ = ud.Value.(*Npc)
			}
			if npc == nil {
				npcID, ok := LuaCheckNpcID(L, 2)
				if ok == false {
					return 0
				}
				npc = ch.GetMap().NpcByTemplate(npcID)
				if npc == nil {
					npc = &Npc{Wz: &wz.NpcSpawn{BaseSpawn: &wz.BaseSpawn{ID: npcID}}}
				}
			}
			if err := ch.OpenNpc(cfg.ActorContext, npc); err != nil {
				L.RaiseError("open_npc: %v", err)
				return 0
			}
			return 0
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
			thread, err := luax.NewThread(root, scriptPath)
			if err != nil {
				L.RaiseError("script: %v", err)
				return 0
			}
			luax.SetConfiguration(thread, luax.Configuration{
				ActorContext: cfg.ActorContext,
				ActorPID:     mapInstance.LogicActorPID(),
			})
			p := luax.CallAsync(cfg.ActorContext, root, thread, funcName, args...)
			if !p.Completed() {
				return 0
			}
			vals, err := p.Result()
			if err != nil {
				L.RaiseError("script: %v", err)
				return 0
			}
			luax.Close(thread)
			if len(vals) == 0 {
				return 0
			}
			L.Push(vals[0])
			return 1
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

			skillEntry := NewSkillEntry(ch, wzSkill, 1, wzSkill.DefaultMasterLevel())
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
			classCode := constant.ClassType(L.CheckInt(2))
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
		"party": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok || ch == nil {
				L.ArgError(1, "Character expected")
				return 0
			}
			if L.GetTop() != 1 {
				L.ArgError(2, "party() takes no arguments")
				return 0
			}
			pid := ch.Party.ID()
			if pid == nil || ch.GameWorld == nil {
				L.Push(lua.LNil)
				return 1
			}
			p := ch.GameWorld.GetPartySystem().Get(*pid)
			if p == nil {
				L.Push(lua.LNil)
				return 1
			}
			L.Push(luax.NewLuable(L, p))
			return 1
		},
		"state_machine": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok || ch == nil {
				L.ArgError(1, "Character expected")
				return 0
			}
			sm := ch.StateMachine()
			if sm == nil {
				L.Push(lua.LNil)
				return 1
			}
			L.Push(luax.NewLuable(L, sm))
			return 1
		},
		"try_party_quest": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok || ch == nil {
				L.ArgError(1, "Character expected")
				return 0
			}
			questID := uint32(L.CheckInt(2))
			ch.TryPartyQuest(questID)
			return 0
		},
		"end_party_quest": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok || ch == nil {
				L.ArgError(1, "Character expected")
				return 0
			}
			questID := uint32(L.CheckInt(2))
			ch.EndPartyQuest(questID)
			return 0
		},
		"carnival_team": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok || ch == nil {
				L.ArgError(1, "Character expected")
				return 0
			}
			team := ch.CarnivalTeam()
			if team == nil {
				L.Push(lua.LNil)
			} else {
				L.Push(luax.NewLuable(L, team))
			}
			return 1
		},
		"carnival_summon": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok || ch == nil {
				L.ArgError(1, "Character expected")
				return 0
			}
			ch.Listener.OnCarnivalSummon(ch, constant.CarnivalTab(L.CheckInt(2)), uint8(L.CheckInt(3)), ch.GetName())
			return 0
		},
		"save_location": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok || ch == nil {
				L.ArgError(1, "Character expected")
				return 0
			}
			ch.SavedLocations.Save(L.CheckString(2))
			return 0
		},
		"saved_location": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok || ch == nil {
				L.ArgError(1, "Character expected")
				return 0
			}
			mapID, ok := ch.SavedLocations.Find(L.CheckString(2))
			if !ok {
				L.Push(lua.LNil)
				return 1
			}
			L.Push(lua.LNumber(mapID))
			return 1
		},
		"clear_saved_location": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok || ch == nil {
				L.ArgError(1, "Character expected")
				return 0
			}
			ch.SavedLocations.Clear(L.CheckString(2))
			return 0
		},
		"records": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok || ch == nil {
				L.ArgError(1, "Character expected")
				return 0
			}
			L.Push(luax.NewLuable(L, ch.Records))
			return 1
		},
		"account_records": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok || ch == nil {
				L.ArgError(1, "Character expected")
				return 0
			}
			L.Push(luax.NewLuable(L, ch.AccountRecords))
			return 1
		},
		"max_hp": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok || ch == nil {
				L.ArgError(1, "Character expected")
				return 0
			}
			if L.GetTop() == 1 {
				L.Push(lua.LNumber(ch.GetMaxHp()))
				return 1
			}
			ch.SetBaseHp(uint32(L.CheckInt(2)), true)
			return 0
		},
		"max_mp": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok || ch == nil {
				L.ArgError(1, "Character expected")
				return 0
			}
			if L.GetTop() == 1 {
				L.Push(lua.LNumber(ch.GetMaxMp()))
				return 1
			}
			ch.SetBaseMp(uint32(L.CheckInt(2)), true)
			return 0
		},
		"quest": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok || ch == nil {
				L.ArgError(1, "Character expected")
				return 0
			}
			if L.GetTop() != 2 {
				L.ArgError(2, "quest(quest_id) requires quest id")
				return 0
			}
			questID := uint32(L.CheckInt(2))
			qp := ch.LuaQuest(questID)
			if qp == nil {
				L.Push(lua.LNil)
				return 1
			}
			L.Push(luax.NewLuable(L, qp))
			return 1
		},
		"run_quest_hook": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok || ch == nil {
				L.ArgError(1, "Character expected")
				return 0
			}
			if L.GetTop() != 3 {
				L.ArgError(2, "run_quest_hook(quest_id, hook) requires quest id and hook")
				return 0
			}
			questID := uint32(L.CheckInt(2))
			hook := L.CheckString(3)
			ch.RunQuestHook(questID, hook)
			return 0
		},
		"quests": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok || ch == nil {
				L.ArgError(1, "Character expected")
				return 0
			}
			if L.GetTop() != 1 {
				L.ArgError(2, "quests() takes no arguments")
				return 0
			}
			tbl := L.NewTable()
			if ch.Quests != nil {
				ch.Quests.ForEach(func(questID uint32, qp *Quest) {
					if qp == nil || qp.Status == QuestStatusNotStarted {
						return
					}
					tbl.RawSet(lua.LNumber(questID), luax.NewLuable(L, qp))
				})
			}
			L.Push(tbl)
			return 1
		},
		"start_quest": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok || ch == nil {
				L.ArgError(1, "Character expected")
				return 0
			}
			if L.GetTop() != 3 {
				L.ArgError(2, "start_quest(quest_id, npc) requires quest id and npc")
				return 0
			}
			questID := uint32(L.CheckInt(2))
			npcID, ok := LuaCheckNpcID(L, 3)
			if ok == false {
				return 0
			}
			if ch.GameWorld == nil {
				L.Push(lua.LNil)
				return 1
			}
			_, err := ch.Quests.Start(questID, QuestPhaseOpts{NpcID: &npcID})
			if err != nil {
				L.Push(lua.LNil)
				return 1
			}
			qp := ch.Quests.Get(questID)
			if qp == nil {
				L.Push(lua.LNil)
				return 1
			}
			L.Push(luax.NewLuable(L, qp))
			return 1
		},
		"clear_quests": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok {
				L.ArgError(1, "Character expected")
				return 0
			}
			if L.GetTop() == 1 {
				L.Push(lua.LNumber(ch.Quests.ClearAll()))
				return 1
			}
			if L.GetTop() == 2 {
				questID := uint32(L.CheckInt(2))
				cleared := 0
				if ch.Quests.Clear(questID) {
					cleared = 1
				}
				L.Push(lua.LNumber(cleared))
				return 1
			}
			L.ArgError(2, "clear_quests() or clear_quests(quest_id)")
			return 0
		},
		"completed_quest_count": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok || ch == nil {
				L.ArgError(1, "Character expected")
				return 0
			}
			if L.GetTop() != 1 {
				L.ArgError(2, "completed_quest_count() takes no arguments")
				return 0
			}
			L.Push(lua.LNumber(ch.Quests.CompletedCount()))
			return 1
		},
		"clear_inventory": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok {
				L.ArgError(1, "Character expected")
				return 0
			}
			if L.GetTop() != 1 {
				L.ArgError(2, "clear_inventory() takes no arguments")
				return 0
			}
			L.Push(lua.LNumber(ch.Inventory.ClearInventory()))
			return 1
		},
		"register_card": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok {
				L.ArgError(1, "Character expected")
				return 0
			}
			cardID := uint32(L.CheckInt(2))
			card, ok := ch.itemModel(cardID).(*wz.Consume)
			if ok == false || card.MonsterBook == false {
				L.Push(lua.LFalse)
				return 1
			}
			for range L.OptInt(3, 1) {
				ch.RegisterCard(cardID)
			}
			L.Push(lua.LTrue)
			return 1
		},
		"reset_monster_book": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok {
				L.ArgError(1, "Character expected")
				return 0
			}
			ch.ResetMonsterBook()
			return 0
		},
		"reset_teleport_stones": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok {
				L.ArgError(1, "Character expected")
				return 0
			}
			ch.TeleportStones.Reset()
			return 0
		},
		"guild": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok || ch == nil {
				L.ArgError(1, "Character expected")
				return 0
			}
			if L.GetTop() != 1 {
				L.ArgError(2, "guild() takes no arguments")
				return 0
			}
			guildID, inGuild := ch.Guild.ID()
			if !inGuild || ch.GameWorld == nil {
				L.Push(lua.LNil)
				return 1
			}
			g := ch.GameWorld.GetGuildSystem().Get(guildID)
			if g == nil {
				L.Push(lua.LNil)
				return 1
			}
			L.Push(luax.NewLuable(L, g))
			return 1
		},
		"show_guild_ranking": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok || ch == nil {
				L.ArgError(1, "Character expected")
				return 0
			}
			npc := 0
			switch npcArg := L.Get(2).(type) {
			case lua.LNumber:
				npc = int(npcArg)
			case *lua.LUserData:
				if n, ok := npcArg.Value.(*Npc); ok && n.Wz != nil && n.Wz.BaseSpawn != nil {
					npc = int(n.Wz.BaseSpawn.ID)
				}
			}
			cfg, ok := luax.GetConfiguration(L)
			if ok == false || cfg.ActorContext == nil || ch.GameWorld == nil {
				return 0
			}
			ch.GameWorld.GetGuildSystem().ShowRankingAsync(cfg.ActorContext, ch, uint32(npc))
			return 0
		},
		"create_alliance": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok || ch == nil {
				L.ArgError(1, "Character expected")
				return 0
			}
			return luaCreateAlliance(L, ch)
		},
		"disband_alliance": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok || ch == nil {
				L.ArgError(1, "Character expected")
				return 0
			}
			return luaDisbandAlliance(L, ch)
		},
		"inc_alliance_capacity": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok || ch == nil {
				L.ArgError(1, "Character expected")
				return 0
			}
			return luaIncAllianceCapacity(L, ch)
		},
		"generic_guild_message": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok || ch == nil {
				L.ArgError(1, "Character expected")
				return 0
			}
			code := pconst.GuildResponseCode(L.CheckInt(2))
			err := ch.Send(&response.GuildMessage{
				Code: code,
			}, types.SEND_POLICY_ENCRYPT)
			if err != nil {
				L.Push(lua.LNumber(constant.GuildCreateResultSendFailed))
				return 1
			}
			if code == pconst.GuildResponseEmblemDialog {
				return 0
			}
			ch.Dialog.Set(L)
			return L.Yield(lua.LNumber(0))
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
		"in_dialog": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok {
				L.ArgError(1, "Character expected")
				return 0
			}
			L.Push(lua.LBool(ch.Dialog.Thread() != nil))
			return 1
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
				L.Push(lua.LNumber(ch.Stats.Base.Str))
				return 1
			case 2:
				v := L.CheckInt(2)
				if v < 0 {
					v = 0
				}
				if v > int(constant.StatMaxStrDexIntLuk) {
					v = int(constant.StatMaxStrDexIntLuk)
				}
				ch.Stats.Base.Str = uint16(v)
				ch.notifyStatChange(constant.StatStr)
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
				L.Push(lua.LNumber(ch.Stats.Base.Dex))
				return 1
			case 2:
				v := L.CheckInt(2)
				if v < 0 {
					v = 0
				}
				if v > int(constant.StatMaxStrDexIntLuk) {
					v = int(constant.StatMaxStrDexIntLuk)
				}
				ch.Stats.Base.Dex = uint16(v)
				ch.notifyStatChange(constant.StatDex)
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
				L.Push(lua.LNumber(ch.Stats.Base.Int))
				return 1
			case 2:
				v := L.CheckInt(2)
				if v < 0 {
					v = 0
				}
				if v > int(constant.StatMaxStrDexIntLuk) {
					v = int(constant.StatMaxStrDexIntLuk)
				}
				ch.Stats.Base.Int = uint16(v)
				ch.notifyStatChange(constant.StatInt)
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
				L.Push(lua.LNumber(ch.Stats.Base.Luk))
				return 1
			case 2:
				v := L.CheckInt(2)
				if v < 0 {
					v = 0
				}
				if v > int(constant.StatMaxStrDexIntLuk) {
					v = int(constant.StatMaxStrDexIntLuk)
				}
				ch.Stats.Base.Luk = uint16(v)
				ch.notifyStatChange(constant.StatLuk)
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
				L.Push(lua.LNumber(ch.Stats.Bonus.Str))
				return 1
			case 2:
				ch.Stats.Bonus.Str = int16(L.CheckInt(2))
				ch.notifyStatChange(constant.StatStr)
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
				L.Push(lua.LNumber(ch.Stats.Bonus.Dex))
				return 1
			case 2:
				ch.Stats.Bonus.Dex = int16(L.CheckInt(2))
				ch.notifyStatChange(constant.StatDex)
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
				L.Push(lua.LNumber(ch.Stats.Bonus.Int))
				return 1
			case 2:
				ch.Stats.Bonus.Int = int16(L.CheckInt(2))
				ch.notifyStatChange(constant.StatInt)
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
				L.Push(lua.LNumber(ch.Stats.Bonus.Luk))
				return 1
			case 2:
				ch.Stats.Bonus.Luk = int16(L.CheckInt(2))
				ch.notifyStatChange(constant.StatLuk)
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
				L.Push(lua.LNumber(ch.Stats.Bonus.Watk))
				return 1
			case 2:
				ch.Stats.Bonus.Watk = int16(L.CheckInt(2))
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
				L.Push(lua.LNumber(ch.Stats.Bonus.MaxHpPercent))
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
				L.Push(lua.LNumber(ch.Stats.Bonus.MaxMpPercent))
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
				L.Push(lua.LNumber(ch.Stats.Bonus.MaxHpPercent))
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
				L.Push(lua.LNumber(ch.Stats.Bonus.MaxMpPercent))
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
				L.Push(lua.LNumber(ch.Stats.Bonus.MesoMultiplier))
				return 1
			case 2:
				ch.Stats.Bonus.MesoMultiplier = int16(L.CheckInt(2))
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
				L.Push(lua.LNumber(ch.Stats.Bonus.DropRate))
				return 1
			case 2:
				ch.Stats.Bonus.DropRate = int16(L.CheckInt(2))
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
				L.Push(lua.LNumber(ch.Stats.Bonus.ExpRate))
				return 1
			case 2:
				ch.Stats.Bonus.ExpRate = int16(L.CheckInt(2))
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
				L.Push(lua.LNumber(ch.Stats.Bonus.PotionHealRate))
				return 1
			case 2:
				ch.Stats.Bonus.PotionHealRate = int16(L.CheckInt(2))
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
				L.Push(lua.LNumber(ch.Stats.Bonus.PotionDurationRate))
				return 1
			case 2:
				ch.Stats.Bonus.PotionDurationRate = int16(L.CheckInt(2))
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
			case 2, 3, 4:
				var targetMap *Map
				arg2 := L.Get(2)
				switch v := arg2.(type) {
				case *lua.LUserData:
					var ok bool
					targetMap, ok = v.Value.(*Map)
					if !ok || targetMap == nil {
						L.ArgError(2, "Map, map name (string), or map id (number) expected")
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
					targetMap = ch.GameWorld.GetMapSystem().Find(ch.StateMachine(), mapId)
					if targetMap == nil {
						L.RaiseError("map: map %q (id %d) not found", string(v), mapId)
						return 0
					}
				case lua.LNumber:
					if ch.GameWorld == nil {
						L.RaiseError("map: no context to resolve map id")
						return 0
					}
					mapId := uint32(v)
					targetMap = ch.GameWorld.GetMapSystem().Find(ch.StateMachine(), mapId)
					if targetMap == nil {
						L.RaiseError("map: map id %d not found", mapId)
						return 0
					}
				default:
					L.ArgError(2, "Map, map name (string), or map id (number) expected")
					return 0
				}
				spawnPoint := uint8(0)
				optsIndex := 4
				switch v := L.Get(3).(type) {
				case *lua.LNilType:
				case *lua.LTable:
					optsIndex = 3
				case lua.LNumber:
					spawnPoint = uint8(v)
				case lua.LString:
					portal := targetMap.FindPortalByName(string(v))
					if portal == nil {
						L.RaiseError("map: portal %q not found in map %d", string(v), targetMap.GetMapID())
						return 0
					}
					spawnPoint = portal.Wz.ID
				default:
					L.ArgError(3, "spawn point id (number), portal name (string), or options (table) expected")
					return 0
				}

				relocate := false
				var callback *lua.LFunction
				if opts, ok := L.Get(optsIndex).(*lua.LTable); ok {
					relocate = lua.LVAsBool(opts.RawGetString("relocate"))
					callback, _ = opts.RawGetString("callback").(*lua.LFunction)
				}

				currentMap := ch.GetMap()
				if relocate && currentMap != nil && currentMap.GetMapID() == targetMap.GetMapID() {
					err := ch.Relocate(spawnPoint)
					if err != nil {
						L.RaiseError("warp: %v", err)
						return 0
					}
					if callback != nil {
						L.CallByParam(lua.P{Fn: callback, NRet: 0, Protect: false}, luax.NewLuable(L, ch))
					}
					return 0
				}

				var onEnter func(actor.Context)
				if callback != nil {
					detached, err := luax.Detach(L, callback)
					if err != nil {
						L.RaiseError("map: callback %v", err)
						return 0
					}
					onEnter = func(ctx actor.Context) {
						m := ch.GetMap()
						if m != targetMap {
							return
						}
						root := m.EnsureLuaRoot(ctx)
						fn, err := detached.Attach(root)
						if err == nil {
							err = luax.Spawn(root, fn, luax.Configuration{ActorContext: ctx, ActorPID: m.LogicActorPID()}, ch)
						}
						if err != nil {
							log.Printf("map callback: %v", err)
						}
					}
				}
				cfg, _ := luax.GetConfiguration(L)
				err := ch.GameWorld.GetMapSystem().Warp(cfg.ActorContext, ch, targetMap, spawnPoint, onEnter)
				if err != nil {
					L.RaiseError("warp: %v", err)
				}
				return 0
			default:
				L.ArgError(2, "map() getter: 0 args; setter: map, name (string), or id (number), optional spawnPoint or portal name, optional { relocate, callback }")
				return 0
			}
		},
		"spawn_point": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok {
				L.ArgError(1, "Character expected")
				return 0
			}
			spawnID := ch.GetSpawnPoint()
			L.Push(lua.LNumber(spawnID))
			m := ch.GetMap()
			if m == nil {
				return 1
			}
			portal := m.FindPortal(spawnID)
			if portal == nil || portal.Wz == nil {
				return 1
			}
			L.Push(lua.LString(portal.Wz.Name))
			L.Push(lua.LNumber(portal.Wz.Position.X))
			L.Push(lua.LNumber(portal.Wz.Position.Y))
			return 4
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
			ch.Debuffs.Give(flag, time.Duration(durationMs)*time.Millisecond, x, skillID, skillLevel)
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
			L.Push(lua.LBool(ch.Debuffs.Has(flag)))
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
				ch.Debuffs.Remove(flags...)
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
				if eq, ok := ch.Inventory.Equipped[part]; ok && eq != nil {
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
				invType, slots := ch.Inventory.FindSlots(itemID)
				if invType != constant.InventoryTypeEquipment || len(slots) == 0 {
					L.Push(lua.LBool(false))
					return 1
				}
				slot = slots[0]
			default:
				if itemUD, ok := L.Get(2).(*lua.LUserData); ok && itemUD.Value != nil {
					if item, ok := itemUD.Value.(Item); ok {
						var found bool
						slot, found = ch.Inventory.FindSlot(constant.InventoryTypeEquipment, item)
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
			item := ch.Inventory.Tabs[constant.InventoryTypeEquipment].Items[slot]
			if item == nil {
				L.Push(lua.LBool(false))
				return 1
			}
			parts := constant.EquipmentParts(item.GetModel().GetID())
			if len(parts) == 0 {
				L.Push(lua.LBool(false))
				return 1
			}
			err := ch.Inventory.Equip(slot, parts[0])
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
				for p, eq := range ch.Inventory.Equipped {
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
						for p, eq := range ch.Inventory.Equipped {
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
			err := ch.Inventory.Unequip(parts)
			L.Push(lua.LBool(err == nil))
			return 1
		},
		"empty_slots": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok {
				L.ArgError(1, "Character expected")
				return 0
			}
			if L.GetTop() != 2 {
				L.ArgError(2, "empty_slots(inventory_type) requires inventory type")
				return 0
			}
			inven := ch.Inventory.Tabs[constant.InventoryType(L.CheckInt(2))]
			if inven == nil {
				L.Push(lua.LNumber(0))
				return 1
			}
			L.Push(lua.LNumber(inven.EmptySlotCount()))
			return 1
		},
		"put_key": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok {
				L.ArgError(1, "Character expected")
				return 0
			}
			if L.GetTop() != 4 {
				L.ArgError(2, "put_key(slot, type, action) requires slot, type, and action")
				return 0
			}
			ch.BindKey(L.CheckInt(2), byte(L.CheckInt(3)), int32(L.CheckInt(4)))
			return 0
		},
		"morph": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok {
				L.ArgError(1, "Character expected")
				return 0
			}
			if L.GetTop() == 1 {
				_, value, ok := ch.Buffs.GetBuffValue(constant.BuffFlagMorph)
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
			ch.Buffs.RemoveBuff([]constant.BuffFlag{constant.BuffFlagMorph})
			return 0
		},
		"use_item": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok {
				L.ArgError(1, "Character expected")
				return 0
			}
			if L.GetTop() != 2 {
				L.ArgError(2, "use_item(item_id) requires an item id")
				return 0
			}
			L.Push(lua.LBool(ch.UseItemEffect(uint32(L.CheckInt(2)))))
			return 1
		},
		"buddy_capacity": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok {
				L.ArgError(1, "Character expected")
				return 0
			}
			if L.GetTop() == 1 {
				L.Push(lua.LNumber(ch.Buddies.Capacity()))
				return 1
			}
			ch.Buddies.SetCapacity(uint32(L.CheckInt(2)))
			return 0
		},
		"sync_item": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok {
				L.ArgError(1, "Character expected")
				return 0
			}
			if L.GetTop() != 3 {
				L.ArgError(2, "sync_item(inventory_type, slot)")
				return 0
			}
			invType := constant.InventoryType(L.CheckInt(2))
			slot := int16(L.CheckInt(3))
			item := ch.Inventory.GetItem(invType, slot)
			if item == nil {
				return 0
			}
			ch.Listener.OnInventorySlotUpdated(ch, invType, slot, item)
			return 0
		},
		"mob_drops": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok {
				L.ArgError(1, "Character expected")
				return 0
			}
			if L.GetTop() != 2 {
				L.ArgError(2, "mob_drops(mob_id) requires a mob id")
				return 0
			}
			t := L.NewTable()
			if ch.GameWorld == nil {
				L.Push(t)
				return 1
			}
			resources := ch.GameWorld.GetResources()
			if resources == nil {
				L.Push(t)
				return 1
			}
			mobDrops := resources.MobDrops[uint32(L.CheckInt(2))]
			for i, d := range mobDrops {
				row := L.NewTable()
				row.RawSetString("item", lua.LNumber(d.Item))
				row.RawSetString("money", lua.LNumber(d.Money))
				row.RawSetString("prob", lua.LNumber(d.Prob))
				row.RawSetString("min", lua.LNumber(d.Min))
				row.RawSetString("max", lua.LNumber(d.Max))
				row.RawSetString("quest", lua.LNumber(resources.QuestItems[d.Item]))
				t.RawSetInt(i+1, row)
			}
			L.Push(t)
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
				if inven := ch.Inventory.Tabs[invType]; inven != nil {
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
				item := ch.Inventory.GetItem(invType, slot)
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
				for invType, inven := range ch.Inventory.Tabs {
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
				if inven := ch.Inventory.Tabs[invType]; inven != nil {
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
			isRandomizeStats := false
			switch argc {
			case 1:
				L.ArgError(2, "mkitem(itemIdOrName [, count] [, isRandomizeStats]) requires at least one argument")
				return 0
			case 2:
			case 3:
				if n := int(L.CheckInt(3)); n >= 1 {
					count = uint16(n)
				}
			default:
				if n := L.CheckInt(3); n >= 1 {
					count = uint16(n)
				}
				if argc >= 4 {
					isRandomizeStats = lua.LVAsBool(L.Get(4))
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

			if resources.Items[itemId] == nil {
				L.Push(lua.LNil)
				return 1
			}
			item, err := NewItem(itemId, count, ch.GameWorld)
			if err != nil {
				L.Push(lua.LNil)
				return 1
			}
			if isRandomizeStats {
				if eq, ok := item.(Equipment); ok {
					if em, ok := eq.GetModel().(wz.Equipment); ok {
						eq.GetEquipmentCore().RandomizeStats(em)
					}
				}
			}

			added, err := ch.Inventory.AddItem(item, true)
			if err != nil {
				L.Push(lua.LNil)
				return 1
			}

			L.Push(luax.NewLuable(L, added[0]))
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
				ok := ch.Inventory.RemoveItem(invType, slot, count)
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
				removed := ch.Inventory.RemoveByItemIDCount(itemId, count)
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
			detached, err := luax.Detach(L, L.CheckFunction(5))
			if err != nil {
				L.RaiseError("mktimer: %v", err)
				return 0
			}
			args := luaValuesToInterfaces(L, 6, argc)
			// The timer follows the character across maps; each map runs its own copy of fn.
			var fnRoot *lua.LState
			var fn *lua.LFunction
			callback := func() {
				m := ch.GetMap()
				if m == nil {
					return
				}
				root := m.GetLuaRoot()
				if root == nil {
					return
				}
				if root != fnRoot {
					attached, err := detached.Attach(root)
					if err != nil {
						log.Printf("RunObjectTimer %s: %v", key, err)
						return
					}
					fnRoot = root
					fn = attached
				}
				fullArgs := make([]interface{}, 0, len(args)+1)
				fullArgs = append(fullArgs, ch)
				fullArgs = append(fullArgs, args...)
				if _, err := luax.CallFunction(root, fn, fullArgs...); err != nil {
					log.Printf("RunObjectTimer %s: %v", key, err)
				}
			}
			added := ch.AddTimer(key, time.Duration(intervalMs)*time.Millisecond, repeat, callback)
			L.Push(lua.LBool(added))
			return 1
		},
		"clock": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok {
				L.ArgError(1, "Character expected")
				return 0
			}
			seconds := L.CheckInt(2)
			fn := L.OptFunction(3, nil)
			ch.RemoveTimer(clockTimer)
			ch.Listener.OnClock(ch, int32(seconds))
			if fn == nil {
				return 0
			}

			// The client drops the clock with the field, so the timer never outlives this map or its Lua root.
			added := ch.AddTimer(clockTimer, time.Duration(seconds)*time.Second, false, func() {
				root := ch.GetMap().GetLuaRoot()
				if root == nil {
					return
				}
				if _, err := luax.CallFunction(root, fn, ch); err != nil {
					log.Printf("clock: %v", err)
				}
			})
			if added == false {
				L.RaiseError("clock: timer not added")
			}
			return 0
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
