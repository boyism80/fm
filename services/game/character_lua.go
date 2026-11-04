package game

import (
	"maps"
	"math"
	"time"

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
		invalid := ""
		bonusTbl.ForEach(func(key, value lua.LValue) {
			if invalid != "" {
				return
			}
			name, ok := key.(lua.LString)
			if !ok || (&EquipmentBonusStats{}).stat(string(name)) == nil {
				invalid = "exchange side.bonus: unknown stat " + key.String()
				return
			}
			n, ok := value.(lua.LNumber)
			if !ok || n < math.MinInt16 || n > math.MaxInt16 || n != lua.LNumber(math.Trunc(float64(n))) {
				invalid = "exchange side.bonus." + string(name) + " must be an int16 integer"
				return
			}
			side.Bonus[string(name)] = int16(n)
		})
		if invalid != "" {
			L.ArgError(argIndex, invalid)
			return ExchangeSide{}, false
		}
	}

	return side, true
}

func (ch *Character) LuaTypeName() string {
	return "LuaCharacter"
}

func (ch *Character) LuaBuiltinFuncs() map[string]lua.LGFunction {
	funcs := ch.luaStatFuncs()
	maps.Copy(funcs, ch.luaEffectFuncs())
	maps.Copy(funcs, ch.luaSkillFuncs())
	maps.Copy(funcs, ch.luaDialogFuncs())
	maps.Copy(funcs, ch.luaItemFuncs())
	maps.Copy(funcs, ch.luaPetFuncs())
	maps.Copy(funcs, ch.luaSocialFuncs())
	maps.Copy(funcs, ch.luaAdminFuncs())
	return funcs
}

func (ch *Character) String() string {
	return ch.LuaTypeName()
}

func (ch *Character) Type() lua.LValueType {
	return lua.LTUserData
}
