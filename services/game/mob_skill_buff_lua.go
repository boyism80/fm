package game

import (
	"github.com/boyism80/fm/core/luax"
	lua "github.com/yuin/gopher-lua"
)

func (e *MobBuff) LuaTypeName() string {
	return "LuaMobSkillBuff"
}

func (e *MobBuff) String() string {
	return e.LuaTypeName()
}

func (e *MobBuff) Type() lua.LValueType {
	return lua.LTUserData
}

func (e *MobBuff) LuaBuiltinFuncs() map[string]lua.LGFunction {
	return map[string]lua.LGFunction{
		"causer": func(L *lua.LState) int {
			mb := luax.Check[MobBuff](L, 1)
			if L.GetTop() != 1 {
				L.ArgError(2, "causer() is read-only")
				return 0
			}
			L.Push(lua.LNumber(mb.causer))
			return 1
		},
		"level": func(L *lua.LState) int {
			mb := luax.Check[MobBuff](L, 1)
			if L.GetTop() != 1 {
				L.ArgError(2, "level() is read-only")
				return 0
			}
			L.Push(lua.LNumber(mb.skillLevel))
			return 1
		},
		"wz": func(L *lua.LState) int {
			mb := luax.Check[MobBuff](L, 1)
			if L.GetTop() != 1 {
				L.ArgError(2, "wz() is read-only")
				return 0
			}
			if mb.wz == nil {
				L.Push(lua.LNil)
				return 1
			}
			L.Push(luax.NewLuable(L, mb.wz))
			return 1
		},
		"effect": func(L *lua.LState) int {
			mb := luax.Check[MobBuff](L, 1)
			if L.GetTop() != 1 {
				L.ArgError(2, "effect() takes no arguments")
				return 0
			}
			if mb.wz == nil {
				L.Push(lua.LNil)
				return 1
			}
			ld := mb.wz.GetLevelData(int(mb.skillLevel))
			if ld == nil {
				L.Push(lua.LNil)
				return 1
			}
			L.Push(ld.ToLuaTable(L))
			return 1
		},
	}
}
