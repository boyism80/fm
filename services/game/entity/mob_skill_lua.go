package entity

import (
	"github.com/boyism80/fm/core/clock"
	"github.com/boyism80/fm/core/luax"

	lua "github.com/yuin/gopher-lua"
)

func (s *MobSkill) LuaTypeName() string {
	return "LuaMobSkill"
}

func (s *MobSkill) LuaBuiltinFuncs() map[string]lua.LGFunction {
	return map[string]lua.LGFunction{
		"id": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			skill, ok := ud.Value.(*MobSkill)
			if !ok || skill == nil {
				L.ArgError(1, "MobSkill expected")
				return 0
			}
			if L.GetTop() != 1 {
				L.ArgError(2, "id() is read-only")
				return 0
			}
			L.Push(lua.LNumber(skill.Slot.SkillID))
			return 1
		},
		"level": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			skill, ok := ud.Value.(*MobSkill)
			if !ok || skill == nil {
				L.ArgError(1, "MobSkill expected")
				return 0
			}
			if L.GetTop() != 1 {
				L.ArgError(2, "level() is read-only")
				return 0
			}
			L.Push(lua.LNumber(skill.Slot.Level))
			return 1
		},
		"wz": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			skill, ok := ud.Value.(*MobSkill)
			if !ok || skill == nil {
				L.ArgError(1, "MobSkill expected")
				return 0
			}
			if L.GetTop() != 1 {
				L.ArgError(2, "wz() is read-only")
				return 0
			}
			if skill.LevelData == nil || skill.LevelData.SkillWz == nil {
				L.Push(lua.LNil)
				return 1
			}
			L.Push(luax.NewLuable(L, skill.LevelData.SkillWz))
			return 1
		},
		"effect": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			skill, ok := ud.Value.(*MobSkill)
			if !ok || skill == nil {
				L.ArgError(1, "MobSkill expected")
				return 0
			}
			if L.GetTop() != 1 {
				L.ArgError(2, "effect() takes no arguments")
				return 0
			}
			if skill.LevelData == nil {
				L.Push(lua.LNil)
				return 1
			}
			L.Push(skill.LevelData.ToLuaTable(L))
			return 1
		},
		"last_used_at": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			skill, ok := ud.Value.(*MobSkill)
			if !ok || skill == nil {
				L.ArgError(1, "MobSkill expected")
				return 0
			}
			if L.GetTop() != 1 {
				L.ArgError(2, "last_used_at() is read-only")
				return 0
			}
			if skill.LastUsedAt.IsZero() {
				L.Push(lua.LNumber(0))
				return 1
			}
			L.Push(lua.LNumber(skill.LastUsedAt.UnixMilli()))
			return 1
		},
		"is_on_cooldown": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			skill, ok := ud.Value.(*MobSkill)
			if !ok || skill == nil {
				L.ArgError(1, "MobSkill expected")
				return 0
			}
			if L.GetTop() != 1 {
				L.ArgError(2, "is_on_cooldown() takes no args")
				return 0
			}
			L.Push(lua.LBool(skill.IsOnCooldown(clock.Now())))
			return 1
		},
	}
}

func (s *MobSkill) String() string {
	return s.LuaTypeName()
}

func (s *MobSkill) Type() lua.LValueType {
	return lua.LTUserData
}
