package game

import (
	"github.com/boyism80/fm/core/luax"
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
			sk := luax.Check[SkillBuff](L, 1)
			if L.GetTop() != 1 {
				L.ArgError(2, "causer() is read-only")
				return 0
			}
			L.Push(lua.LNumber(sk.causerID))
			return 1
		},
		"level": func(L *lua.LState) int {
			sk := luax.Check[SkillBuff](L, 1)
			if L.GetTop() != 1 {
				L.ArgError(2, "level() is read-only")
				return 0
			}
			L.Push(lua.LNumber(sk.skillLevel))
			return 1
		},
		"wz": func(L *lua.LState) int {
			sk := luax.Check[SkillBuff](L, 1)
			if L.GetTop() != 1 {
				L.ArgError(2, "wz() is read-only")
				return 0
			}
			L.Push(luax.NewLuable(L, sk.wz))
			return 1
		},
		"effect": func(L *lua.LState) int {
			sk := luax.Check[SkillBuff](L, 1)
			if L.GetTop() != 1 {
				L.ArgError(2, "effect() takes no arguments")
				return 0
			}
			ld := sk.wz.GetLevelData(int(sk.skillLevel))
			if ld == nil {
				L.Push(lua.LNil)
				return 1
			}
			L.Push(ld.ToLuaTable(L))
			return 1
		},
	}
}
