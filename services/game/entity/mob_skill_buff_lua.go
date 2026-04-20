package entity

import (
	lua "github.com/yuin/gopher-lua"
)

func (e *MobSkillBuff) LuaTypeName() string {
	return "LuaMobSkillBuff"
}

func (e *MobSkillBuff) String() string {
	return e.LuaTypeName()
}

func (e *MobSkillBuff) Type() lua.LValueType {
	return lua.LTUserData
}

func (e *MobSkillBuff) LuaBuiltinFuncs() map[string]lua.LGFunction {
	return map[string]lua.LGFunction{
		"causer": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			mb, ok := ud.Value.(*MobSkillBuff)
			if !ok {
				L.ArgError(1, "MobSkillBuff expected")
				return 0
			}
			if L.GetTop() != 1 {
				L.ArgError(2, "causer() is read-only")
				return 0
			}
			L.Push(lua.LNumber(mb.Causer))
			return 1
		},
		"level": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			mb, ok := ud.Value.(*MobSkillBuff)
			if !ok {
				L.ArgError(1, "MobSkillBuff expected")
				return 0
			}
			if L.GetTop() != 1 {
				L.ArgError(2, "level() is read-only")
				return 0
			}
			L.Push(lua.LNumber(mb.SkillLevel))
			return 1
		},
		"wz": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			mb, ok := ud.Value.(*MobSkillBuff)
			if !ok {
				L.ArgError(1, "MobSkillBuff expected")
				return 0
			}
			if L.GetTop() != 1 {
				L.ArgError(2, "wz() is read-only")
				return 0
			}
			wzTable := skillToLuaTable(L, mb.Wz)
			L.Push(wzTable)
			return 1
		},
		"effect": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			mb, ok := ud.Value.(*MobSkillBuff)
			if !ok {
				L.ArgError(1, "MobSkillBuff expected")
				return 0
			}
			if L.GetTop() != 1 {
				L.ArgError(2, "effect() takes no arguments")
				return 0
			}
			ld := mb.Wz.GetLevelData(int(mb.SkillLevel))
			if ld == nil {
				L.Push(lua.LNil)
				return 1
			}
			L.Push(skillLevelDataToLuaTable(L, ld))
			return 1
		},
	}
}
