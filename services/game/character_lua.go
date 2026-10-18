package game

import (
	"log"
	"math"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/core/clock"
	"github.com/boyism80/fm/core/luax"
	"github.com/boyism80/fm/maple/constant"
	"github.com/boyism80/fm/maple/wz"
	pconst "github.com/boyism80/fm/protocol/constant"
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
	"github.com/boyism80/fm/protocol/response"
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
			if countN != lua.LNumber(math.Trunc(float64(countN))) {
				L.ArgError(argIndex, "exchange side.item count must be an integer")
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

	if periodLV := tbl.RawGetString("period"); periodLV != lua.LNil {
		periodTbl, ok := periodLV.(*lua.LTable)
		if !ok {
			L.ArgError(argIndex, "exchange side.period must be a table")
			return ExchangeSide{}, false
		}
		side.Period = make(map[uint32]time.Duration)
		failed := false
		periodTbl.ForEach(func(k, v lua.LValue) {
			if failed {
				return
			}
			minutes, ok := v.(lua.LNumber)
			if !ok || minutes <= 0 {
				L.ArgError(argIndex, "exchange side.period values must be positive minutes")
				failed = true
				return
			}
			switch key := k.(type) {
			case lua.LNumber:
				side.Period[uint32(key)] = time.Duration(minutes) * time.Minute
			case lua.LString:
				if nameToItem == nil {
					L.ArgError(argIndex, "exchange item name lookup unavailable")
					failed = true
					return
				}
				id, ok := nameToItem(string(key))
				if !ok {
					L.ArgError(argIndex, "unknown item name in exchange side.period")
					failed = true
					return
				}
				side.Period[id] = time.Duration(minutes) * time.Minute
			default:
				L.ArgError(argIndex, "exchange side.period keys must be item id or name")
				failed = true
			}
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
		if n != lua.LNumber(math.Trunc(float64(n))) {
			L.ArgError(argIndex, "exchange side.meso must be an integer")
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
		if n != lua.LNumber(math.Trunc(float64(n))) {
			L.ArgError(argIndex, "exchange side.exp must be an integer")
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
		if n != lua.LNumber(math.Trunc(float64(n))) {
			L.ArgError(argIndex, "exchange side.population must be an integer")
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
				if ch.level > oldLevel && ch.quests != nil {
					ch.quests.RunAutoTriggers(nil, AutoQuestTriggerLevelUp, 0)
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
				L.Push(lua.LNumber(ch.points.ap))
				return 1
			case 2, 3:
				v := L.CheckInt(2)
				if v < 0 {
					v = 0
				}
				notify := argc == 2 || L.ToBool(3)
				ch.points.SetAP(uint16(v), notify)
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
				L.Push(lua.LNumber(ch.points.sp))
				return 1
			case 2, 3:
				v := L.CheckInt(2)
				if v < 0 {
					v = 0
				}
				notify := argc == 2 || L.ToBool(3)
				ch.points.SetSP(uint16(v), notify)
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
				L.Push(lua.LNumber(ch.baseHp))
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
				L.Push(lua.LNumber(ch.baseMp))
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
				ch.listener.OnUpdateStats(ch, stats, true)
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
				ch.expReachedAt = clock.Now()
				ch.listener.OnUpdateStats(ch, map[constant.Stat]int32{
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
				L.Push(lua.LNumber(ch.inventory.meso))
				return 1
			case 2:
				value := L.CheckNumber(2)
				if value < 0 {
					L.ArgError(2, "meso() setter requires a non-negative value; compute new total in script")
					return 0
				}
				if value <= 2147483647 {
					ch.inventory.SetMeso(int32(value))
				} else {
					ch.inventory.SetMeso(2147483647)
				}
				ch.listener.OnMesoChanged(ch, ch.inventory.meso)
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
				L.Push(lua.LNumber(ch.stats.population))
				return 1
			case 2:
				ch.stats.setPopulation(int32(L.CheckInt(2)))
				return 0
			default:
				L.ArgError(2, "population() getter: population(); setter: population(value)")
				return 0
			}
		},
		"clear_population_history": func(L *lua.LState) int {
			ch, ok := L.CheckUserData(1).Value.(*Character)
			if !ok {
				L.ArgError(1, "Character expected")
				return 0
			}
			ch.stats.populationHistory = PopulationHistory{given: map[uint32]time.Time{}}
			return 0
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
			if ch.gameWorld != nil {
				if resources := ch.gameWorld.GetResources(); resources != nil {
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

			ch.listener.OnChat(ch, message, highlight, dontRecordHistory)
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
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok {
				L.ArgError(1, "Character expected")
				return 0
			}
			skillUD := L.CheckUserData(2)
			skillEntry, ok := skillUD.Value.(*SkillEntry)
			if !ok || skillEntry == nil || skillEntry.wz == nil {
				L.ArgError(2, "SkillEntry with Wz expected")
				return 0
			}
			L.Push(lua.LBool(ch.mount.Ride(skillEntry, int32(L.CheckInt(3))) == nil))
			return 1
		},
		"dismount": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok {
				L.ArgError(1, "Character expected")
				return 0
			}
			ch.mount.Dismount()
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
					ch.listener.OnShowSelfSkillEffect(ch, pconst.SkillEffectType(effectType), skillEntry.wz.ID, skillLevel, nil)
					return 0
				}
				L.ArgError(3, "show_skill_effect() supports SkillEffectType.Cast or SkillEffectType.Affected")
				return 0
			}

			ch.listener.OnShowSelfSkillEffect(ch, pconst.SkillEffectTypeCast, skillEntry.wz.ID, skillLevel, nil)
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
				ch.listener.OnShowSelfEffect(ch, effectType)
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
			if ch.listener != nil {
				ch.listener.OnShowQuestCompletion(ch, questID)
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
			if ch.listener != nil {
				ch.listener.OnPlaySound(ch, sound, broadcast)
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
			if ch.listener != nil {
				ch.listener.OnPlayPortalSound(ch)
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
			if !ok || skillEntry == nil || skillEntry.wz == nil {
				L.ArgError(2, "SkillEntry with Wz expected")
				return 0
			}
			skillLevelInt := skillEntry.Level()
			if skillLevelInt < 0 {
				skillLevelInt = 0
			}
			ch.listener.OnShowSelfDragonBloodEffect(ch, skillEntry.wz.ID, uint8(skillLevelInt))
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
			ch.listener.OnShowSelfHPHealedEffect(ch, amount)
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
			ch.listener.OnShowSelfRewardItemAnimation(ch, itemID, effectText)
			return 0
		},
		"show_item_maker_success_effect": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok {
				L.ArgError(1, "Character expected")
				return 0
			}
			ch.listener.OnShowSelfItemMakerSuccessEffect(ch)
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
			ch.listener.OnShowSelfCraftingEffect(ch, effectText, timeValue, modeValue)
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
			if !ok || skillEntry == nil || skillEntry.wz == nil {
				L.ArgError(3, "SkillEntry with Wz expected")
				return 0
			}
			skillLevelInt := skillEntry.Level()
			if skillLevelInt < 0 {
				skillLevelInt = 0
			}
			ch.listener.OnShowSelfDiceEffect(ch, effectID, skillEntry.wz.ID, uint8(skillLevelInt))
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
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok {
				L.ArgError(1, "Character expected")
				return 0
			}
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

			pos := ch.position
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
			s := ch.summons.Spawn(constant.SkillID(skillID), skillLevel, movType, summonType, pos, duration)
			if s == nil {
				L.Push(lua.LNil)
				return 1
			}
			L.Push(luax.NewLuable(L, s))
			return 1
		},
		"skill_area": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok {
				L.ArgError(1, "Character expected")
				return 0
			}
			skill, ok := L.CheckUserData(2).Value.(*SkillEntry)
			if !ok || skill == nil {
				L.ArgError(2, "Skill expected")
				return 0
			}
			ld := skill.wz.GetLevelData(skill.Level())
			if ld == nil {
				L.Push(lua.LNil)
				return 1
			}
			area := ld.Area(int32(ch.position.X), int32(ch.position.Y), ch.stance%2 != 0)
			tbl := L.NewTable()
			tbl.RawSetString("left", lua.LNumber(area.Left))
			tbl.RawSetString("top", lua.LNumber(area.Top))
			tbl.RawSetString("right", lua.LNumber(area.Right))
			tbl.RawSetString("bottom", lua.LNumber(area.Bottom))
			L.Push(tbl)
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
			ch.doors.Spawn(buff)
			return 0
		},
		"remove_summon": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok {
				L.ArgError(1, "Character expected")
				return 0
			}
			ch.summons.Remove(ch.summons.Get(constant.SkillID(L.CheckInt(2))), true)
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
			ch.doors.Remove(buff, true)
			return 0
		},
		"clear_summons": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok {
				L.ArgError(1, "Character expected")
				return 0
			}
			ch.summons.Clear()
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
					if n, ok := npcArg.Value.(*Npc); ok && n.wz != nil && n.wz.BaseSpawn != nil {
						npc = int(n.wz.BaseSpawn.ID)
					}
				}
			}
			ch.listener.OnDialog(ch, uint32(npc), message, prev, next)
			ch.dialog.Ask(L, constant.DialogTypeDefault, 0)
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
					if n, ok := npcArg.Value.(*Npc); ok && n.wz != nil && n.wz.BaseSpawn != nil {
						npc = int(n.wz.BaseSpawn.ID)
					}
				}
			}
			ch.listener.OnDialogYesNo(ch, uint32(npc), message, prev, next)
			ch.dialog.Ask(L, constant.DialogTypeYesNo, 0)
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
					if n, ok := npcArg.Value.(*Npc); ok && n.wz != nil && n.wz.BaseSpawn != nil {
						npc = int(n.wz.BaseSpawn.ID)
					}
				}
			}

			ch.listener.OnDialogList(ch, uint32(npc), message, selections)
			ch.dialog.Ask(L, constant.DialogTypeList, len(selections))
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
				if n, ok := npcArg.Value.(*Npc); ok && n.wz != nil && n.wz.BaseSpawn != nil {
					npc = int(n.wz.BaseSpawn.ID)
				}
			}
			styles := []uint32{}
			L.CheckTable(4).ForEach(func(_, value lua.LValue) {
				if n, ok := value.(lua.LNumber); ok {
					styles = append(styles, uint32(n))
				}
			})
			ch.listener.OnDialogStyle(ch, uint32(npc), L.CheckString(3), styles)
			ch.dialog.Ask(L, constant.DialogTypeStyle, len(styles))
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
					if n, ok := npcArg.Value.(*Npc); ok && n.wz != nil && n.wz.BaseSpawn != nil {
						npc = int(n.wz.BaseSpawn.ID)
					}
				}
			}
			ch.listener.OnDialogAccept(ch, uint32(npc), message, enableEscape)
			if enableEscape {
				ch.dialog.Ask(L, constant.DialogTypeAcceptEscape, 0)
			} else {
				ch.dialog.Ask(L, constant.DialogTypeAccept, 0)
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
					if n, ok := npcArg.Value.(*Npc); ok && n.wz != nil && n.wz.BaseSpawn != nil {
						npc = int(n.wz.BaseSpawn.ID)
					}
				}
			}

			ch.listener.OnDialogInput(ch, uint32(npc), message)
			ch.dialog.Ask(L, constant.DialogTypeInput, 0)
			return L.Yield(lua.LNumber(0))
		},
		"dialog_number": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok {
				L.ArgError(1, "Character expected")
				return 0
			}
			npc := 0
			switch npcArg := L.Get(2).(type) {
			case lua.LNumber:
				npc = int(npcArg)
			case *lua.LUserData:
				if n, ok := npcArg.Value.(*Npc); ok && n.wz != nil && n.wz.BaseSpawn != nil {
					npc = int(n.wz.BaseSpawn.ID)
				}
			}
			message := L.CheckString(3)
			def := int32(L.CheckInt(4))
			min := int32(L.CheckInt(5))
			max := int32(L.CheckInt(6))
			if min > max {
				L.ArgError(5, "min must not exceed max")
				return 0
			}

			ch.listener.OnDialogNumber(ch, uint32(npc), message, def, min, max)
			ch.dialog.AskNumber(L, min, max)
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
				if m := ch.GetMap(); m != nil {
					m.MapMessage(msgType, text)
				}
			case constant.MessageScopeChannel:
				if ch.gameWorld != nil {
					ch.gameWorld.Systems().Get[*NoticeSystem]().BroadcastNotice(msgType, text, 0, false)
				}
			case constant.MessageScopeWorld:
				cfg, ok := luax.GetConfiguration(L)
				if !ok || cfg.ActorContext == nil || ch.gameWorld == nil {
					L.Push(lua.LFalse)
					return 1
				}
				promise := ch.gameWorld.Systems().Get[*InternalSystem]().BroadcastNoticeAsync(cfg.ActorContext, ch, msgType, text, ear)
				if promise == nil {
					L.Push(lua.LFalse)
					return 1
				}
				return LuaYieldPromise(L, ch.gameWorld, promise, func(_ *internal.BroadcastNoticeReply, err error) []lua.LValue {
					if err != nil {
						return []lua.LValue{lua.LFalse}
					}
					return []lua.LValue{lua.LTrue}
				})
			default:
				ch.listener.OnMessage(ch, msgType, text)
			}
			L.Push(lua.LTrue)
			return 1
		},
		"megaphone_lines": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok {
				L.ArgError(1, "Character expected")
				return 0
			}
			tbl := L.CheckTable(2)
			ear := L.ToBool(3)
			lines := make([]string, 0, 3)
			for i := 1; i <= tbl.Len(); i++ {
				lines = append(lines, lua.LVAsString(tbl.RawGetInt(i)))
			}
			cfg, ok := luax.GetConfiguration(L)
			if !ok || cfg.ActorContext == nil || len(lines) < 1 || len(lines) > 3 {
				L.Push(lua.LFalse)
				return 1
			}
			promise := ch.gameWorld.Systems().Get[*InternalSystem]().BroadcastMegaphoneAsync(cfg.ActorContext, ch, &response.Megaphone{
				Type:  constant.MsgThreeMegaphoneLines,
				Lines: lines,
				Ear:   ear,
			})
			return LuaYieldPromise(L, ch.gameWorld, promise, func(_ *internal.BroadcastNoticeReply, err error) []lua.LValue {
				return []lua.LValue{lua.LBool(err == nil)}
			})
		},
		"item_megaphone": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok {
				L.ArgError(1, "Character expected")
				return 0
			}
			megaphone := &response.Megaphone{
				Type:  constant.MsgItemMegaphone,
				Lines: []string{L.CheckString(2)},
				Ear:   L.ToBool(3),
			}
			if L.Get(4) != lua.LNil {
				invType := constant.InventoryType(L.CheckInt(4))
				slot := L.CheckInt(5)
				if invType < constant.InventoryTypeEquipment || invType > constant.InventoryTypeCash || slot < 1 || slot > 255 {
					L.Push(lua.LFalse)
					return 1
				}
				item := ch.inventory.tabs[invType].Get(uint8(slot))
				if item == nil {
					L.Push(lua.LFalse)
					return 1
				}
				megaphone.Item = item.BuildS2C()
			}
			cfg, ok := luax.GetConfiguration(L)
			if !ok || cfg.ActorContext == nil {
				L.Push(lua.LFalse)
				return 1
			}
			promise := ch.gameWorld.Systems().Get[*InternalSystem]().BroadcastMegaphoneAsync(cfg.ActorContext, ch, megaphone)
			return LuaYieldPromise(L, ch.gameWorld, promise, func(_ *internal.BroadcastNoticeReply, err error) []lua.LValue {
				return []lua.LValue{lua.LBool(err == nil)}
			})
		},
		"set_chalkboard": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok {
				L.ArgError(1, "Character expected")
				return 0
			}
			L.Push(lua.LBool(ch.SetChalkboard(L.CheckString(2))))
			return 1
		},
		"place_message_box": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok {
				L.ArgError(1, "Character expected")
				return 0
			}
			box := NewMessageBox(ch, uint32(L.CheckInt(2)), L.CheckString(3))
			L.Push(lua.LBool(box != nil && box.Place(ch)))
			return 1
		},
		"send_note": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok {
				L.ArgError(1, "Character expected")
				return 0
			}
			recipient := L.CheckString(2)
			message := L.CheckString(3)
			if !ClientText(message).Valid(200) {
				L.Push(lua.LFalse)
				return 1
			}
			if !ClientText(recipient).Valid(13) {
				ch.listener.OnNoteFailed(ch, internal.NoteResult_NOTE_RESULT_WRONG_NAME)
				L.Push(lua.LFalse)
				return 1
			}
			cfg, ok := luax.GetConfiguration(L)
			if !ok || cfg.ActorContext == nil {
				L.Push(lua.LFalse)
				return 1
			}
			promise := ch.gameWorld.Systems().Get[*InternalSystem]().SendNoteAsync(cfg.ActorContext, ch, recipient, message)
			return LuaYieldPromise(L, ch.gameWorld, promise, func(reply *internal.NotesReply, err error) []lua.LValue {
				if err != nil {
					ch.listener.OnNoteFailed(ch, internal.NoteResult_NOTE_RESULT_INBOX_FULL)
					return []lua.LValue{lua.LFalse}
				}
				if reply.Status != internal.NoteResult_NOTE_RESULT_OK {
					ch.listener.OnNoteFailed(ch, reply.Status)
					return []lua.LValue{lua.LFalse}
				}
				ch.listener.OnNoteSent(ch)
				return []lua.LValue{lua.LTrue}
			})
		},
		"marriage": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok {
				L.ArgError(1, "Character expected")
				return 0
			}
			if ch.wedding.marriage == nil {
				L.Push(lua.LNil)
				return 1
			}
			L.Push(luax.NewLuable(L, ch.wedding.marriage))
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
			promise, err := ch.wedding.Reserve(cfg.ActorContext, uint32(L.CheckNumber(2)))
			if err != nil {
				L.Push(lua.LFalse)
				return 1
			}
			return LuaYieldPromise(L, ch.gameWorld, promise, func(result bool, err error) []lua.LValue {
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
			promise, err := ch.wedding.BreakEngagement(cfg.ActorContext)
			if err != nil {
				L.Push(lua.LFalse)
				return 1
			}
			return LuaYieldPromise(L, ch.gameWorld, promise, func(result bool, err error) []lua.LValue {
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
			promise, err := ch.wedding.RequestDivorce(cfg.ActorContext)
			if err != nil {
				L.Push(lua.LNil)
				return 1
			}
			return LuaYieldPromise(L, ch.gameWorld, promise, func(result string, err error) []lua.LValue {
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
			promise, err := ch.wedding.CancelDivorce(cfg.ActorContext)
			if err != nil {
				L.Push(lua.LFalse)
				return 1
			}
			return LuaYieldPromise(L, ch.gameWorld, promise, func(canceled bool, err error) []lua.LValue {
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
			L.Push(lua.LBool(ch.wedding.InvitedTo(uint32(L.CheckNumber(2)))))
			return 1
		},
		"open_wedding_wishlist": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok {
				L.ArgError(1, "Character expected")
				return 0
			}
			ch.listener.OnWeddingWishlistInput(ch)
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
			ch.wedding.OpenGift(receiverID, wishes)
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
			return LuaYieldPromise(L, ch.gameWorld, ch.wedding.OpenGiftBox(cfg.ActorContext), func(result bool, err error) []lua.LValue {
				if err != nil {
					return []lua.LValue{lua.LFalse}
				}
				return []lua.LValue{lua.LBool(result)}
			})
		},
		"open_storage": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok {
				L.ArgError(1, "Character expected")
				return 0
			}
			npcID := uint32(L.CheckInt(2))
			storeFee := int32(L.OptInt(3, 100))
			takeOutFee := int32(L.OptInt(4, 0))
			ch.account.storage.Open(npcID, storeFee, takeOutFee)
			return 0
		},
		"open_store_bank": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok {
				L.ArgError(1, "Character expected")
				return 0
			}
			npcID := uint32(L.CheckInt(2))
			cfg, ok := luax.GetConfiguration(L)
			if !ok || cfg.ActorContext == nil {
				L.ArgError(1, "actor context not available")
				return 0
			}
			ch.storeBank.Open(cfg.ActorContext, npcID)
			return 0
		},
		"open_duey": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok {
				L.ArgError(1, "Character expected")
				return 0
			}
			cfg, ok := luax.GetConfiguration(L)
			if !ok || cfg.ActorContext == nil {
				L.ArgError(1, "actor context not available")
				return 0
			}
			ch.duey.Open(cfg.ActorContext, false)
			return 0
		},
		"pet_id": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok {
				L.ArgError(1, "Character expected")
				return 0
			}
			if ch.pets.active == nil {
				L.Push(lua.LNumber(0))
				return 1
			}
			L.Push(lua.LNumber(ch.pets.active.item.GetModel().GetID()))
			return 1
		},
		"pet_closeness": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok {
				L.ArgError(1, "Character expected")
				return 0
			}
			if ch.pets.active == nil {
				L.Push(lua.LNumber(0))
				return 1
			}
			L.Push(lua.LNumber(ch.pets.active.item.closeness))
			return 1
		},
		"add_pet_closeness": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok {
				L.ArgError(1, "Character expected")
				return 0
			}
			if ch.pets.active == nil {
				L.Push(lua.LFalse)
				return 1
			}
			ch.pets.active.AddCloseness(L.CheckInt(2))
			L.Push(lua.LTrue)
			return 1
		},
		"set_pet_closeness": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok {
				L.ArgError(1, "Character expected")
				return 0
			}
			if ch.pets.active == nil {
				L.Push(lua.LFalse)
				return 1
			}
			ch.pets.active.SetCloseness(L.CheckInt(2))
			L.Push(lua.LTrue)
			return 1
		},
		"set_pet_fullness": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok {
				L.ArgError(1, "Character expected")
				return 0
			}
			if ch.pets.active == nil {
				L.Push(lua.LFalse)
				return 1
			}
			ch.pets.active.SetFullness(L.CheckInt(2))
			L.Push(lua.LTrue)
			return 1
		},
		"hunger_pet": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok {
				L.ArgError(1, "Character expected")
				return 0
			}
			if ch.pets.active == nil {
				L.Push(lua.LFalse)
				return 1
			}
			ch.pets.active.hunger()
			L.Push(lua.LTrue)
			return 1
		},
		"mount": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok {
				L.ArgError(1, "Character expected")
				return 0
			}
			L.Push(lua.LNumber(ch.mount.level))
			L.Push(lua.LNumber(ch.mount.exp))
			L.Push(lua.LNumber(ch.mount.fatigue))
			L.Push(lua.LNumber(ch.mount.hp))
			_, vehicle, _ := ch.buffs.GetBuffValue(constant.BuffFlagMonsterRiding)
			L.Push(lua.LNumber(vehicle))
			return 5
		},
		"set_mount_level": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok {
				L.ArgError(1, "Character expected")
				return 0
			}
			ch.mount.SetLevel(uint32(L.CheckInt(2)))
			return 0
		},
		"set_mount_fatigue": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok {
				L.ArgError(1, "Character expected")
				return 0
			}
			ch.mount.SetFatigue(uint32(L.CheckInt(2)))
			return 0
		},
		"set_mount_hp": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok {
				L.ArgError(1, "Character expected")
				return 0
			}
			ch.mount.hp = uint32(max(L.CheckInt(2), 0))
			return 0
		},
		"start_mount_fatigue": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok {
				L.ArgError(1, "Character expected")
				return 0
			}
			ch.mount.startFatigue()
			return 0
		},
		"show_skill_cooldown": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok {
				L.ArgError(1, "Character expected")
				return 0
			}
			ch.listener.OnSkillCooldown(ch, uint32(L.CheckInt(2)), uint16(L.CheckInt(3)))
			return 0
		},
		"tire_mount": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok {
				L.ArgError(1, "Character expected")
				return 0
			}
			ch.mount.tire()
			return 0
		},
		"pet_level": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok {
				L.ArgError(1, "Character expected")
				return 0
			}
			if ch.pets.active == nil {
				L.Push(lua.LNumber(0))
				return 1
			}
			L.Push(lua.LNumber(ch.pets.active.item.level))
			return 1
		},
		"pet_fullness": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok {
				L.ArgError(1, "Character expected")
				return 0
			}
			if ch.pets.active == nil {
				L.Push(lua.LNumber(0))
				return 1
			}
			L.Push(lua.LNumber(ch.pets.active.item.fullness))
			return 1
		},
		"rename_pet": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok {
				L.ArgError(1, "Character expected")
				return 0
			}
			if ch.pets.active == nil {
				L.Push(lua.LFalse)
				return 1
			}
			L.Push(lua.LBool(ch.pets.active.Rename(L.CheckString(2)) == nil))
			return 1
		},
		"feed_pet_cash": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok {
				L.ArgError(1, "Character expected")
				return 0
			}
			if ch.pets.active == nil {
				L.Push(lua.LFalse)
				return 1
			}
			L.Push(lua.LBool(ch.pets.active.FeedCash(uint32(L.CheckInt(2))) == nil))
			return 1
		},
		"change_pet_skill": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok {
				L.ArgError(1, "Character expected")
				return 0
			}
			sn := uint64(L.CheckNumber(2))
			L.Push(lua.LBool(ch.pets.ChangeSkill(sn, uint32(L.CheckInt(3))) == nil))
			return 1
		},
		"expired_pets": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok {
				L.ArgError(1, "Character expected")
				return 0
			}
			tbl := L.NewTable()
			for slot, pet := range ch.pets.Expired() {
				entry := L.NewTable()
				entry.RawSetString("slot", lua.LNumber(slot))
				entry.RawSetString("id", lua.LNumber(pet.GetModel().GetID()))
				entry.RawSetString("name", lua.LString(pet.name))
				tbl.Append(entry)
			}
			L.Push(tbl)
			return 1
		},
		"revive_pet": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok {
				L.ArgError(1, "Character expected")
				return 0
			}
			L.Push(lua.LBool(ch.pets.Revive(int16(L.CheckInt(2))) == nil))
			return 1
		},
		"open_quick_parcel": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok {
				L.ArgError(1, "Character expected")
				return 0
			}
			L.Push(lua.LBool(ch.duey.OpenQuick() == nil))
			return 1
		},
		"send_parcel": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok {
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
			if !ok || cfg.ActorContext == nil {
				L.ArgError(1, "actor context not available")
				return 0
			}
			L.Push(lua.LBool(ch.duey.SendFromSystem(cfg.ActorContext, recipient, senderName, itemID, count, meso, message) == nil))
			return 1
		},
		"add_cash": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok {
				L.ArgError(1, "Character expected")
				return 0
			}
			nxCash := int32(L.CheckInt(2))
			maplePoint := int32(L.OptInt(3, 0))
			cfg, ok := luax.GetConfiguration(L)
			if !ok || cfg.ActorContext == nil {
				L.ArgError(1, "actor context not available")
				return 0
			}
			ch.AddCash(cfg.ActorContext, nxCash, maplePoint)
			return 0
		},
		"create_cash_coupons": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok {
				L.ArgError(1, "Character expected")
				return 0
			}
			kind := internal.CashCouponKind(L.CheckInt(2))
			value := uint32(L.CheckInt(3))
			count := uint32(L.OptInt(4, 1))
			cfg, ok := luax.GetConfiguration(L)
			if !ok || cfg.ActorContext == nil {
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
				L.Push(lua.LNumber(ch.role))
				return 1
			case 2:
				ch.role = constant.CharacterRole(L.CheckInt(2))
				return 0
			default:
				L.ArgError(2, "role() requires 0 or 1 arguments")
				return 0
			}
		},
		"account_role": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok {
				L.ArgError(1, "Character expected")
				return 0
			}
			role := constant.CharacterRole(L.CheckInt(2))
			cfg, ok := luax.GetConfiguration(L)
			if !ok || cfg.ActorContext == nil {
				L.Push(lua.LFalse)
				return 1
			}
			return LuaYieldPromise(L, ch.gameWorld, ch.gameWorld.Systems().Get[*InternalSystem]().SetAccountRoleAsync(cfg.ActorContext, ch, role), func(_ *internal.SetAccountRoleReply, err error) []lua.LValue {
				if err != nil {
					return []lua.LValue{lua.LFalse}
				}
				ch.role = role
				return []lua.LValue{lua.LTrue}
			})
		},
		"disconnect": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok {
				L.ArgError(1, "Character expected")
				return 0
			}
			ch.listener.Disconnect(ch)
			return 0
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
				L.Push(lua.LNumber(ch.dojo.energy))
				return 1
			case 2:
				ch.dojo.SetEnergy(L.CheckInt(2))
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
				L.Push(lua.LBool(ch.gm.instantKill))
				return 1
			case 2:
				ch.gm.instantKill = L.CheckBool(2)
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
				L.Push(lua.LBool(ch.gm.playerMode))
				return 1
			case 2:
				ch.gm.playerMode = L.CheckBool(2)
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
				L.Push(lua.LNumber(ch.gm.timerLimit))
				return 1
			case 2:
				ch.gm.timerLimit = uint32(L.CheckNumber(2))
				return 0
			default:
				L.ArgError(2, "timer_limit() requires 0 or 1 arguments")
				return 0
			}
		},
		"lie_detector_answer": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok {
				L.ArgError(1, "Character expected")
				return 0
			}
			m := ch.GetMap()
			if m == nil {
				L.Push(lua.LNil)
				return 1
			}
			target := m.FindPlayerByName(L.CheckString(2))
			if target == nil || target.lieDetector.testerID == 0 {
				L.Push(lua.LNil)
				return 1
			}
			L.Push(lua.LString(target.lieDetector.answer))
			return 1
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
				if !ok {
					return 0
				}
				npc = ch.GetMap().NpcByTemplate(npcID)
				if npc == nil {
					npc = &Npc{wz: &wz.NpcSpawn{BaseSpawn: &wz.BaseSpawn{ID: npcID}}}
				}
			}
			ch.dialog.Close()
			if err := ch.dialog.Open(cfg.ActorContext, npc); err != nil {
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

			skillEntry := ch.skills.Get(skillID)
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
			ch.skills.ForEach(func(skillID uint32, skillEntry *SkillEntry) {
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
				L.Push(lua.LNumber(ch.class))
				return 1
			case 2:
				classValue := L.CheckInt(2)
				ch.ChangeClass(constant.ClassType(classValue))
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
			L.Push(lua.LNumber(ch.chair))
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
				L.Push(lua.LNumber(ch.stance))
				return 1
			case 2:
				ch.stance = uint8(L.CheckInt(2))
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
					if ch.stance == v {
						L.Push(lua.LTrue)
						return 1
					}
				}
				L.Push(lua.LFalse)
				return 1
			}
			L.Push(lua.LBool(ch.stance == uint8(group)))
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
			L.Push(lua.LBool(ch.class.Of(classCode)))
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
			if ch.skills.homingTargetOID == nil {
				L.Push(lua.LNil)
				return 1
			}
			mapInstance := ch.GetMap()
			if mapInstance == nil {
				L.Push(lua.LNil)
				return 1
			}
			mob := mapInstance.GetMob(*ch.skills.homingTargetOID)
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
			pid := ch.party.ID()
			if pid == nil || ch.gameWorld == nil {
				L.Push(lua.LNil)
				return 1
			}
			p := ch.gameWorld.Systems().Get[*PartySystem]().Get(*pid)
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
			if ch.stateMachine == nil {
				L.Push(lua.LNil)
				return 1
			}
			L.Push(luax.NewLuable(L, ch.stateMachine))
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
			ch.quests.RunHook(questID, "on_quest_try")
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
			ch.quests.RunHook(questID, "on_quest_end")
			return 0
		},
		"carnival_team": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok || ch == nil {
				L.ArgError(1, "Character expected")
				return 0
			}
			team := ch.carnival.Team()
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
			ch.listener.OnCarnivalSummon(ch, constant.CarnivalTab(L.CheckInt(2)), uint8(L.CheckInt(3)), ch.GetName())
			return 0
		},
		"save_location": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok || ch == nil {
				L.ArgError(1, "Character expected")
				return 0
			}
			ch.savedLocations.Save(L.CheckString(2))
			return 0
		},
		"saved_location": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok || ch == nil {
				L.ArgError(1, "Character expected")
				return 0
			}
			mapID, ok := ch.savedLocations.Find(L.CheckString(2))
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
			ch.savedLocations.Clear(L.CheckString(2))
			return 0
		},
		"records": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok || ch == nil {
				L.ArgError(1, "Character expected")
				return 0
			}
			L.Push(luax.NewLuable(L, ch.records))
			return 1
		},
		"account_records": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok || ch == nil {
				L.ArgError(1, "Character expected")
				return 0
			}
			L.Push(luax.NewLuable(L, ch.account.records))
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
			qp := ch.quests.LuaQuest(questID)
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
			ch.quests.RunHook(questID, hook)
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
			if ch.quests != nil {
				ch.quests.ForEach(func(questID uint32, qp *Quest) {
					if qp == nil || qp.status == QuestStatusNotStarted {
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
			if !ok {
				return 0
			}
			if ch.gameWorld == nil {
				L.Push(lua.LNil)
				return 1
			}
			_, err := ch.quests.Start(questID, QuestPhaseOpts{NpcID: &npcID})
			if err != nil {
				L.Push(lua.LNil)
				return 1
			}
			qp := ch.quests.Get(questID)
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
				L.Push(lua.LNumber(ch.quests.ClearAll()))
				return 1
			}
			if L.GetTop() == 2 {
				questID := uint32(L.CheckInt(2))
				cleared := 0
				if ch.quests.Clear(questID) {
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
			L.Push(lua.LNumber(ch.quests.CompletedCount()))
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
			L.Push(lua.LNumber(ch.inventory.ClearInventory()))
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
			if !ok || !card.MonsterBook {
				L.Push(lua.LFalse)
				return 1
			}
			for range L.OptInt(3, 1) {
				ch.monsterBook.Register(cardID)
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
			ch.monsterBook.Reset()
			return 0
		},
		"reset_teleport_stones": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok {
				L.ArgError(1, "Character expected")
				return 0
			}
			ch.teleportStones.Reset()
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
			guildID, inGuild := ch.guild.ID()
			if !inGuild || ch.gameWorld == nil {
				L.Push(lua.LNil)
				return 1
			}
			g := ch.gameWorld.Systems().Get[*GuildSystem]().Get(guildID)
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
				if n, ok := npcArg.Value.(*Npc); ok && n.wz != nil && n.wz.BaseSpawn != nil {
					npc = int(n.wz.BaseSpawn.ID)
				}
			}
			cfg, ok := luax.GetConfiguration(L)
			if !ok || cfg.ActorContext == nil || ch.gameWorld == nil {
				return 0
			}
			ch.gameWorld.Systems().Get[*GuildSystem]().ShowRankingAsync(cfg.ActorContext, ch, uint32(npc))
			return 0
		},
		"gain_guild_gp": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok || ch == nil {
				L.ArgError(1, "Character expected")
				return 0
			}
			amount := L.CheckInt(2)
			cfg, ok := luax.GetConfiguration(L)
			if !ok || cfg.ActorContext == nil {
				return 0
			}
			ch.gameWorld.Systems().Get[*GuildSystem]().GainGPAsync(cfg.ActorContext, ch, int32(amount))
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
			ch.listener.OnGuildMessage(ch, pconst.GuildResponseCode(L.CheckInt(2)))
			return 0
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
		"family": func(L *lua.LState) int {
			ch, ok := L.CheckUserData(1).Value.(*Character)
			if !ok {
				L.ArgError(1, "Character expected")
				return 0
			}
			self := ch.family.Self()
			table := L.NewTable()
			table.RawSetString("leader", lua.LNumber(ch.family.tree.Leader().CharacterID()))
			table.RawSetString("senior", lua.LNumber(self.SeniorID()))
			table.RawSetString("juniors", lua.LNumber(len(ch.family.tree.Juniors(ch.GetID()))))
			table.RawSetString("reputation", lua.LNumber(self.Reputation()))
			table.RawSetString("total_reputation", lua.LNumber(self.TotalReputation()))
			L.Push(table)
			return 1
		},
		"in_dialog": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok {
				L.ArgError(1, "Character expected")
				return 0
			}
			L.Push(lua.LBool(ch.dialog.Thread() != nil))
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
				L.Push(lua.LNumber(ch.stats.base.str))
				return 1
			case 2:
				v := L.CheckInt(2)
				if v < 0 {
					v = 0
				}
				if v > int(constant.StatMaxStrDexIntLuk) {
					v = int(constant.StatMaxStrDexIntLuk)
				}
				ch.stats.base.str = uint16(v)
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
				L.Push(lua.LNumber(ch.stats.base.dex))
				return 1
			case 2:
				v := L.CheckInt(2)
				if v < 0 {
					v = 0
				}
				if v > int(constant.StatMaxStrDexIntLuk) {
					v = int(constant.StatMaxStrDexIntLuk)
				}
				ch.stats.base.dex = uint16(v)
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
				L.Push(lua.LNumber(ch.stats.base.int))
				return 1
			case 2:
				v := L.CheckInt(2)
				if v < 0 {
					v = 0
				}
				if v > int(constant.StatMaxStrDexIntLuk) {
					v = int(constant.StatMaxStrDexIntLuk)
				}
				ch.stats.base.int = uint16(v)
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
				L.Push(lua.LNumber(ch.stats.base.luk))
				return 1
			case 2:
				v := L.CheckInt(2)
				if v < 0 {
					v = 0
				}
				if v > int(constant.StatMaxStrDexIntLuk) {
					v = int(constant.StatMaxStrDexIntLuk)
				}
				ch.stats.base.luk = uint16(v)
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
				L.Push(lua.LNumber(ch.stats.bonus.str))
				return 1
			case 2:
				ch.stats.bonus.str = int16(L.CheckInt(2))
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
				L.Push(lua.LNumber(ch.stats.bonus.dex))
				return 1
			case 2:
				ch.stats.bonus.dex = int16(L.CheckInt(2))
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
				L.Push(lua.LNumber(ch.stats.bonus.int))
				return 1
			case 2:
				ch.stats.bonus.int = int16(L.CheckInt(2))
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
				L.Push(lua.LNumber(ch.stats.bonus.luk))
				return 1
			case 2:
				ch.stats.bonus.luk = int16(L.CheckInt(2))
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
				L.Push(lua.LNumber(ch.stats.bonus.watk))
				return 1
			case 2:
				ch.stats.bonus.watk = int16(L.CheckInt(2))
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
				L.Push(lua.LNumber(ch.stats.bonus.maxHpPercent))
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
				L.Push(lua.LNumber(ch.stats.bonus.maxMpPercent))
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
				L.Push(lua.LNumber(ch.stats.bonus.maxHpPercent))
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
				L.Push(lua.LNumber(ch.stats.bonus.maxMpPercent))
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
				L.Push(lua.LNumber(ch.stats.bonus.mesoMultiplier))
				return 1
			case 2:
				ch.stats.bonus.mesoMultiplier = int16(L.CheckInt(2))
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
				L.Push(lua.LNumber(ch.stats.bonus.dropRate))
				return 1
			case 2:
				ch.stats.bonus.dropRate = int16(L.CheckInt(2))
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
			if ch.gameWorld == nil {
				L.Push(lua.LNumber(1))
				return 1
			}
			L.Push(lua.LNumber(ch.gameWorld.GetDropRate()))
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
				L.Push(lua.LNumber(ch.stats.bonus.expRate))
				return 1
			case 2:
				ch.stats.bonus.expRate = int16(L.CheckInt(2))
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
				L.Push(lua.LNumber(ch.stats.bonus.potionHealRate))
				return 1
			case 2:
				ch.stats.bonus.potionHealRate = int16(L.CheckInt(2))
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
				L.Push(lua.LNumber(ch.stats.bonus.potionDurationRate))
				return 1
			case 2:
				ch.stats.bonus.potionDurationRate = int16(L.CheckInt(2))
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
					if ch.gameWorld == nil {
						L.RaiseError("map: no context to resolve map name")
						return 0
					}
					resources := ch.gameWorld.GetResources()
					if resources == nil {
						L.RaiseError("map: no resources to resolve map name")
						return 0
					}
					mapId, ok := resources.NameToMap(string(v))
					if !ok {
						L.RaiseError("map: unknown map name %q", string(v))
						return 0
					}
					targetMap = ch.gameWorld.Systems().Get[*MapSystem]().Find(ch.stateMachine, mapId)
					if targetMap == nil {
						L.RaiseError("map: map %q (id %d) not found", string(v), mapId)
						return 0
					}
				case lua.LNumber:
					if ch.gameWorld == nil {
						L.RaiseError("map: no context to resolve map id")
						return 0
					}
					mapId := uint32(v)
					targetMap = ch.gameWorld.Systems().Get[*MapSystem]().Find(ch.stateMachine, mapId)
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
					spawnPoint = portal.wz.ID
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
						root := m.LuaRoot(ctx)
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
				err := ch.gameWorld.Systems().Get[*MapSystem]().Warp(cfg.ActorContext, ch, targetMap, spawnPoint, onEnter)
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
			if portal == nil || portal.wz == nil {
				return 1
			}
			L.Push(lua.LString(portal.wz.Name))
			L.Push(lua.LNumber(portal.wz.Position.X))
			L.Push(lua.LNumber(portal.wz.Position.Y))
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
				ch.Broadcast(&response.ShowMagnet{MobID: mob.oid, Success: successByte}, &ObjectBroadcastOption{})
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
			ch.debuffs.Give(flag, time.Duration(durationMs)*time.Millisecond, x, skillID, skillLevel)
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
			L.Push(lua.LBool(ch.debuffs.Has(flag)))
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
				ch.debuffs.Remove(flags...)
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
				if eq, ok := ch.inventory.equipped[part]; ok && eq != nil {
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
				if ch.gameWorld == nil {
					L.Push(lua.LBool(false))
					return 1
				}
				resources := ch.gameWorld.GetResources()
				if resources == nil {
					L.Push(lua.LBool(false))
					return 1
				}
				itemID, ok := resources.NameToItem(string(lv))
				if !ok {
					L.Push(lua.LBool(false))
					return 1
				}
				invType, slots := ch.inventory.FindSlots(itemID)
				if invType != constant.InventoryTypeEquipment || len(slots) == 0 {
					L.Push(lua.LBool(false))
					return 1
				}
				slot = slots[0]
			default:
				if itemUD, ok := L.Get(2).(*lua.LUserData); ok && itemUD.Value != nil {
					if item, ok := itemUD.Value.(Item); ok {
						var found bool
						slot, found = ch.inventory.FindSlot(constant.InventoryTypeEquipment, item)
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
			item := ch.inventory.tabs[constant.InventoryTypeEquipment].items[slot]
			if item == nil {
				L.Push(lua.LBool(false))
				return 1
			}
			model, ok := item.GetModel().(wz.Equipment)
			if !ok {
				L.Push(lua.LBool(false))
				return 1
			}
			parts := model.GetParts()
			if len(parts) == 0 {
				L.Push(lua.LBool(false))
				return 1
			}
			err := ch.inventory.Equip(slot, parts[0])
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
				if ch.gameWorld == nil {
					L.Push(lua.LBool(false))
					return 1
				}
				resources := ch.gameWorld.GetResources()
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
				for p, eq := range ch.inventory.equipped {
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
						for p, eq := range ch.inventory.equipped {
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
			err := ch.inventory.Unequip(parts)
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
			inven := ch.inventory.tabs[constant.InventoryType(L.CheckInt(2))]
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
			ch.keyLayout.SetKey(L.CheckInt(2), byte(L.CheckInt(3)), int32(L.CheckInt(4)))
			ch.listener.OnKeyMap(ch)
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
				L.Push(lua.LNumber(ch.buddies.Capacity()))
				return 1
			}
			ch.buddies.SetCapacity(uint32(L.CheckInt(2)))
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
			item := ch.inventory.GetItem(invType, slot)
			if item == nil {
				return 0
			}
			ch.listener.OnInventorySlotUpdated(ch, invType, slot, item)
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
			if ch.gameWorld == nil {
				L.Push(t)
				return 1
			}
			resources := ch.gameWorld.GetResources()
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
					if ch.gameWorld == nil {
						L.Push(L.NewTable())
						return 1
					}
					resources := ch.gameWorld.GetResources()
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
				slotTbl := L.NewTable()
				model := ch.itemModel(itemID)
				if model == nil {
					L.Push(slotTbl)
					return 1
				}
				if inven := ch.inventory.tabs[model.GetInventoryType()]; inven != nil {
					for slot, item := range inven.items {
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
				item := ch.inventory.GetItem(invType, slot)
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
				for invType, inven := range ch.inventory.tabs {
					if inven == nil {
						continue
					}
					invTbl := L.NewTable()
					for slot, item := range inven.items {
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
				if inven := ch.inventory.tabs[invType]; inven != nil {
					for slot, item := range inven.items {
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
			if ch.gameWorld == nil {
				L.Push(lua.LNil)
				return 1
			}
			resources := ch.gameWorld.GetResources()
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
			item, err := NewItem(itemId, count, ch.gameWorld)
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

			added, err := ch.inventory.AddItem(item, true)
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
				ok := ch.inventory.RemoveItem(invType, slot, count)
				L.Push(lua.LBool(ok))
				return 1
			case 2, 3:
				if ch.gameWorld == nil {
					L.Push(lua.LBool(false))
					return 1
				}
				resources := ch.gameWorld.GetResources()
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
				removed := ch.inventory.RemoveByItemIDCount(itemId, count)
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
			ch.listener.OnClock(ch, int32(seconds))
			if fn == nil {
				return 0
			}

			// The client drops the clock with the field, so the timer never outlives this map or its Lua root.
			added := ch.AddTimer(clockTimer, time.Duration(seconds)*time.Second, false, func() {
				m := ch.GetMap()
				if m == nil {
					return
				}
				root := m.GetLuaRoot()
				if root == nil {
					return
				}
				if _, err := luax.CallFunction(root, fn, ch); err != nil {
					log.Printf("clock: %v", err)
				}
			})
			if !added {
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
