package entity

import (
	"time"

	"github.com/boyism80/fm/core/luax"
	lua "github.com/yuin/gopher-lua"
)

func (s *SkillEntry) LuaTypeName() string {
	return "LuaSkill"
}

func (s *SkillEntry) LuaBuiltinFuncs() map[string]lua.LGFunction {
	return map[string]lua.LGFunction{
		"level": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			skill, ok := ud.Value.(*SkillEntry)
			if !ok {
				L.ArgError(1, "Skill expected")
				return 0
			}

			argc := L.GetTop()
			switch argc {
			case 1:
				L.Push(lua.LNumber(skill.Level()))
				return 1
			case 2:
				skill.SetLevel(L.CheckInt(2))
				return 0
			case 3:
				skill.SetLevelAndMaster(L.CheckInt(2), L.CheckInt(3))
				return 0
			default:
				L.ArgError(2, "level() requires 1, 2 or 3 arguments")
				return 0
			}
		},
		"master_level": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			skill, ok := ud.Value.(*SkillEntry)
			if !ok {
				L.ArgError(1, "Skill expected")
				return 0
			}

			argc := L.GetTop()
			switch argc {
			case 1:
				L.Push(lua.LNumber(skill.MasterLevel))
				return 1
			case 2:
				skill.SetMasterLevel(L.CheckInt(2))
				return 0
			default:
				L.ArgError(2, "master_level() requires 0 or 1 arguments")
				return 0
			}
		},
		"wz": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			skill, ok := ud.Value.(*SkillEntry)
			if !ok {
				L.ArgError(1, "Skill expected")
				return 0
			}
			if skill.Wz == nil {
				L.Push(lua.LNil)
				return 1
			}
			L.Push(luax.NewLuable(L, skill.Wz))
			return 1
		},
		"effect": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			skill, ok := ud.Value.(*SkillEntry)
			if !ok {
				L.ArgError(1, "Skill expected")
				return 0
			}
			if L.GetTop() != 1 {
				L.ArgError(2, "effect() takes no arguments")
				return 0
			}
			ld := skill.Wz.GetLevelData(skill.Level())
			if ld == nil {
				L.Push(lua.LNil)
				return 1
			}
			L.Push(ld.ToLuaTable(L))
			return 1
		},
		"is_cooling": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			skill, ok := ud.Value.(*SkillEntry)
			if !ok {
				L.ArgError(1, "Skill expected")
				return 0
			}
			L.Push(lua.LBool(skill.CooldownRemaining() > 0))
			return 1
		},
		"cooldown": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			skill, ok := ud.Value.(*SkillEntry)
			if !ok {
				L.ArgError(1, "Skill expected")
				return 0
			}
			argc := L.GetTop()
			switch argc {
			case 1:
				L.Push(lua.LNumber(skill.CooldownRemaining().Milliseconds()))
				return 1
			case 2:
				ms := L.CheckNumber(2)
				if ms <= 0 {
					skill.ClearCooldown()
				} else {
					skill.StartCooldown(time.Duration(ms) * time.Millisecond)
				}
				return 0
			default:
				L.ArgError(2, "cooldown() requires 0 or 1 arguments")
				return 0
			}
		},
	}
}

func (s *SkillEntry) String() string {
	return s.LuaTypeName()
}

func (s *SkillEntry) Type() lua.LValueType {
	return lua.LTUserData
}
