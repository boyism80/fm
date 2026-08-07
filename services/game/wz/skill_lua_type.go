package wz

import (
	"github.com/boyism80/fm/core/luax"
	lua "github.com/yuin/gopher-lua"
)

const luaWzSkillTypeName = "LuaWzSkill"

func (*Skill) LuaTypeName() string { return luaWzSkillTypeName }

func (*Skill) LuaBuiltinFuncs() map[string]lua.LGFunction {
	return map[string]lua.LGFunction{
		"id": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			s, ok := ud.Value.(*Skill)
			if !ok || s == nil {
				L.ArgError(1, "WzSkill expected")
				return 0
			}
			L.Push(lua.LNumber(s.ID))
			return 1
		},
		"skill_id": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			s, ok := ud.Value.(*Skill)
			if !ok || s == nil {
				L.ArgError(1, "WzSkill expected")
				return 0
			}
			L.Push(lua.LNumber(s.ID))
			return 1
		},
		"max_level": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			s, ok := ud.Value.(*Skill)
			if !ok || s == nil {
				L.ArgError(1, "WzSkill expected")
				return 0
			}
			L.Push(lua.LNumber(s.MaxLevel))
			return 1
		},
		"master_level": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			s, ok := ud.Value.(*Skill)
			if !ok || s == nil {
				L.ArgError(1, "WzSkill expected")
				return 0
			}
			L.Push(lua.LNumber(s.MasterLevel))
			return 1
		},
		"invisible": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			s, ok := ud.Value.(*Skill)
			if !ok || s == nil {
				L.ArgError(1, "WzSkill expected")
				return 0
			}
			L.Push(lua.LBool(s.Invisible))
			return 1
		},
		"time_limited": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			s, ok := ud.Value.(*Skill)
			if !ok || s == nil {
				L.ArgError(1, "WzSkill expected")
				return 0
			}
			L.Push(lua.LBool(s.TimeLimited))
			return 1
		},
		"combat_orders": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			s, ok := ud.Value.(*Skill)
			if !ok || s == nil {
				L.ArgError(1, "WzSkill expected")
				return 0
			}
			L.Push(lua.LBool(s.CombatOrders))
			return 1
		},
		"elem_attr": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			s, ok := ud.Value.(*Skill)
			if !ok || s == nil {
				L.ArgError(1, "WzSkill expected")
				return 0
			}
			L.Push(lua.LString(s.ElemAttr))
			return 1
		},
		"effects": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			s, ok := ud.Value.(*Skill)
			if !ok || s == nil {
				L.ArgError(1, "WzSkill expected")
				return 0
			}
			effectsTable := L.NewTable()
			if s.LevelData != nil {
				for level, levelData := range s.LevelData {
					if levelData == nil {
						continue
					}
					effectsTable.RawSetInt(level, levelData.ToLuaTable(L))
				}
			}
			L.Push(effectsTable)
			return 1
		},
	}
}

func (s *Skill) String() string       { return s.LuaTypeName() }
func (s *Skill) Type() lua.LValueType { return lua.LTUserData }

var _ luax.Luable = (*Skill)(nil)
