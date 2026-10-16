package game

import (
	"github.com/boyism80/fm/core/luax"
	lua "github.com/yuin/gopher-lua"
)

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
			ud := L.CheckUserData(1)
			m, ok := ud.Value.(*Mist)
			if !ok {
				L.ArgError(1, "Mist expected")
				return 0
			}
			if L.GetTop() != 1 {
				L.ArgError(2, "causer() is read-only")
				return 0
			}
			L.Push(lua.LNumber(m.causer))
			return 1
		},
		"level": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			m, ok := ud.Value.(*Mist)
			if !ok {
				L.ArgError(1, "Mist expected")
				return 0
			}
			if L.GetTop() != 1 {
				L.ArgError(2, "level() is read-only")
				return 0
			}
			L.Push(lua.LNumber(m.skillLevel))
			return 1
		},
		"wz": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			m, ok := ud.Value.(*Mist)
			if !ok {
				L.ArgError(1, "Mist expected")
				return 0
			}
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
			ud := L.CheckUserData(1)
			m, ok := ud.Value.(*Mist)
			if !ok {
				L.ArgError(1, "Mist expected")
				return 0
			}
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
			ud := L.CheckUserData(1)
			m, ok := ud.Value.(*Mist)
			if !ok {
				L.ArgError(1, "Mist expected")
				return 0
			}
			if L.GetTop() != 1 {
				L.ArgError(2, "from_mob() takes no arguments")
				return 0
			}
			L.Push(lua.LBool(m.mobMist))
			return 1
		},
		"poison_tick_multiplier": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			m, ok := ud.Value.(*Mist)
			if !ok {
				L.ArgError(1, "Mist expected")
				return 0
			}
			if L.GetTop() != 1 {
				L.ArgError(2, "poison_tick_multiplier() is read-only")
				return 0
			}
			L.Push(lua.LNumber(m.poisonTickMultiplier))
			return 1
		},
	}
}
