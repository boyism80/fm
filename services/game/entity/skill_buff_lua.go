package entity

import (
	lua "github.com/yuin/gopher-lua"
)

func (b *SkillBuff) LuaTypeName() string {
	return "LuaSkillBuff"
}

func (b *SkillBuff) String() string {
	return b.LuaTypeName()
}

func (b *SkillBuff) Type() lua.LValueType {
	return lua.LTUserData
}

func (b *SkillBuff) LuaBuiltinFuncs() map[string]lua.LGFunction {
	return map[string]lua.LGFunction{
		"causer": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			sk, ok := ud.Value.(*SkillBuff)
			if !ok {
				L.ArgError(1, "SkillBuff expected")
				return 0
			}
			if L.GetTop() != 1 {
				L.ArgError(2, "causer() is read-only")
				return 0
			}
			L.Push(lua.LNumber(sk.CauserID))
			return 1
		},
		"level": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			sk, ok := ud.Value.(*SkillBuff)
			if !ok {
				L.ArgError(1, "SkillBuff expected")
				return 0
			}
			if L.GetTop() != 1 {
				L.ArgError(2, "level() is read-only")
				return 0
			}
			L.Push(lua.LNumber(sk.SkillLevel))
			return 1
		},
		"wz": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			sk, ok := ud.Value.(*SkillBuff)
			if !ok {
				L.ArgError(1, "SkillBuff expected")
				return 0
			}
			if L.GetTop() != 1 {
				L.ArgError(2, "wz() is read-only")
				return 0
			}
			wzTable := skillToLuaTable(L, sk.Wz)
			L.Push(wzTable)
			return 1
		},
		"effect": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			sk, ok := ud.Value.(*SkillBuff)
			if !ok {
				L.ArgError(1, "SkillBuff expected")
				return 0
			}
			if L.GetTop() != 1 {
				L.ArgError(2, "effect() takes no arguments")
				return 0
			}
			ld := sk.Wz.GetLevelData(int(sk.SkillLevel))
			if ld == nil {
				L.Push(lua.LNil)
				return 1
			}
			L.Push(skillLevelDataToLuaTable(L, ld))
			return 1
		},
	}
}
