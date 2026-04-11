package entity

import lua "github.com/yuin/gopher-lua"

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
			L.Push(lua.LNumber(m.Causer))
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
			L.Push(lua.LNumber(m.SkillLevel))
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
			wzTable := skillToLuaTable(L, m.SkillWz)
			L.Push(wzTable)
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
			ld := m.SkillWz.GetLevelData(int(m.SkillLevel))
			if ld == nil {
				L.Push(lua.LNil)
				return 1
			}
			L.Push(skillLevelDataToLuaTable(L, ld))
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
			L.Push(lua.LNumber(m.PoisonTickMultiplier))
			return 1
		},
	}
}
