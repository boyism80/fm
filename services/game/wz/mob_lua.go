package wz

import (
	"github.com/boyism80/fm/core/luax"
	lua "github.com/yuin/gopher-lua"
)

const luaWzMobTypeName = "LuaWzMob"

func (*Mob) LuaTypeName() string { return luaWzMobTypeName }

func (*Mob) LuaBuiltinFuncs() map[string]lua.LGFunction {
	return map[string]lua.LGFunction{
		"id": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			m, ok := ud.Value.(*Mob)
			if !ok || m == nil {
				L.ArgError(1, "WzMob expected")
				return 0
			}
			L.Push(lua.LNumber(m.ID))
			return 1
		},
		"level": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			m, ok := ud.Value.(*Mob)
			if !ok || m == nil {
				L.ArgError(1, "WzMob expected")
				return 0
			}
			L.Push(lua.LNumber(m.Level))
			return 1
		},
		"max_hp": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			m, ok := ud.Value.(*Mob)
			if !ok || m == nil {
				L.ArgError(1, "WzMob expected")
				return 0
			}
			L.Push(lua.LNumber(m.MaxHP))
			return 1
		},
		"max_mp": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			m, ok := ud.Value.(*Mob)
			if !ok || m == nil {
				L.ArgError(1, "WzMob expected")
				return 0
			}
			L.Push(lua.LNumber(m.MaxMP))
			return 1
		},
		"exp": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			m, ok := ud.Value.(*Mob)
			if !ok || m == nil {
				L.ArgError(1, "WzMob expected")
				return 0
			}
			L.Push(lua.LNumber(m.EXP))
			return 1
		},
		"boss": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			m, ok := ud.Value.(*Mob)
			if !ok || m == nil {
				L.ArgError(1, "WzMob expected")
				return 0
			}
			L.Push(lua.LBool(m.Boss))
			return 1
		},
		"elem_resist": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			m, ok := ud.Value.(*Mob)
			if !ok || m == nil {
				L.ArgError(1, "WzMob expected")
				return 0
			}
			er := L.NewTable()
			for k, v := range m.ElemResist {
				er.RawSetString(k, lua.LNumber(v))
			}
			L.Push(er)
			return 1
		},
	}
}

func (m *Mob) String() string       { return m.LuaTypeName() }
func (m *Mob) Type() lua.LValueType { return lua.LTUserData }

var _ luax.Luable = (*Mob)(nil)
