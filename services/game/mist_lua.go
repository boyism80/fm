package game

import (
	"time"

	"github.com/boyism80/fm/core/luax"
	"github.com/boyism80/fm/maple/constant"
	"github.com/boyism80/fm/types"
	lua "github.com/yuin/gopher-lua"
)

// Parse reads create_mist(skill, duration_ms, type, bounds[, poison_delay_ms[, poison_tick_multiplier]]) after the skill.
func (spec *mistSpec) Parse(L *lua.LState) {
	spec.duration = time.Duration(max(L.CheckInt(3), 0)) * time.Millisecond
	spec.mistType = constant.MistType(L.CheckInt(4))
	bounds := L.CheckTable(5)
	spec.bounds = types.Rect[int32]{
		Left:   int32(lua.LVAsNumber(bounds.RawGetString("left"))),
		Top:    int32(lua.LVAsNumber(bounds.RawGetString("top"))),
		Right:  int32(lua.LVAsNumber(bounds.RawGetString("right"))),
		Bottom: int32(lua.LVAsNumber(bounds.RawGetString("bottom"))),
	}
	spec.poisonDelay = mistPoisonDelay
	if L.GetTop() >= 6 && L.CheckInt(6) > 0 {
		spec.poisonDelay = time.Duration(L.CheckInt(6)) * time.Millisecond
	}
	spec.poisonTickMultiplier = 1.0
	if L.GetTop() >= 7 && L.CheckNumber(7) > 0 {
		spec.poisonTickMultiplier = float64(L.CheckNumber(7))
	}
}

func (mist *Mist) LuaTypeName() string {
	return "LuaMist"
}

func (mist *Mist) String() string {
	return mist.LuaTypeName()
}

func (mist *Mist) Type() lua.LValueType {
	return lua.LTUserData
}

func (mist *Mist) LuaBuiltinFuncs() map[string]lua.LGFunction {
	return map[string]lua.LGFunction{
		"causer": func(L *lua.LState) int {
			m := luax.Check[Mist](L, 1)
			if L.GetTop() != 1 {
				L.ArgError(2, "causer() is read-only")
				return 0
			}
			L.Push(lua.LNumber(m.causer))
			return 1
		},
		"level": func(L *lua.LState) int {
			m := luax.Check[Mist](L, 1)
			if L.GetTop() != 1 {
				L.ArgError(2, "level() is read-only")
				return 0
			}
			L.Push(lua.LNumber(m.skillLevel))
			return 1
		},
		"wz": func(L *lua.LState) int {
			m := luax.Check[Mist](L, 1)
			if L.GetTop() != 1 {
				L.ArgError(2, "wz() is read-only")
				return 0
			}
			if m.mobLevelData != nil {
				if m.mobLevelData.SkillWz != nil {
					L.Push(luax.NewLuable(L, m.mobLevelData.SkillWz))
					return 1
				}
				L.Push(lua.LNil)
				return 1
			}
			if m.skillWz == nil {
				L.Push(lua.LNil)
				return 1
			}
			L.Push(luax.NewLuable(L, m.skillWz))
			return 1
		},
		"effect": func(L *lua.LState) int {
			m := luax.Check[Mist](L, 1)
			if L.GetTop() != 1 {
				L.ArgError(2, "effect() takes no arguments")
				return 0
			}
			if m.mobLevelData != nil {
				L.Push(m.mobLevelData.ToLuaTable(L))
				return 1
			}
			if m.skillWz == nil {
				L.Push(lua.LNil)
				return 1
			}
			ld := m.skillWz.GetLevelData(int(m.skillLevel))
			if ld == nil {
				L.Push(lua.LNil)
				return 1
			}
			L.Push(ld.ToLuaTable(L))
			return 1
		},
		"from_mob": func(L *lua.LState) int {
			m := luax.Check[Mist](L, 1)
			if L.GetTop() != 1 {
				L.ArgError(2, "from_mob() takes no arguments")
				return 0
			}
			L.Push(lua.LBool(m.mobMist))
			return 1
		},
		"poison_tick_multiplier": func(L *lua.LState) int {
			m := luax.Check[Mist](L, 1)
			if L.GetTop() != 1 {
				L.ArgError(2, "poison_tick_multiplier() is read-only")
				return 0
			}
			L.Push(lua.LNumber(m.poisonTickMultiplier))
			return 1
		},
	}
}
