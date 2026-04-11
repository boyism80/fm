package entity

import (
	"time"

	"github.com/boyism80/fm/game/wz"
	lua "github.com/yuin/gopher-lua"
)

// Luable interface implementation
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

			wzTable := skillToLuaTable(L, skill.Wz)
			L.Push(wzTable)
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
			L.Push(skillLevelDataToLuaTable(L, ld))
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

// skillToLuaTable converts wz.Skill to Lua table
func skillToLuaTable(L *lua.LState, skill *wz.Skill) *lua.LTable {
	tbl := L.NewTable()

	// Basic skill properties
	tbl.RawSetString("id", lua.LNumber(skill.ID))
	tbl.RawSetString("max_level", lua.LNumber(skill.MaxLevel))
	tbl.RawSetString("master_level", lua.LNumber(skill.MasterLevel))
	tbl.RawSetString("invisible", lua.LBool(skill.Invisible))
	tbl.RawSetString("time_limited", lua.LBool(skill.TimeLimited))
	tbl.RawSetString("combat_orders", lua.LBool(skill.CombatOrders))
	if skill.ElemAttr != "" {
		tbl.RawSetString("elem_attr", lua.LString(skill.ElemAttr))
	}

	// Convert LevelData to effects table
	effectsTable := L.NewTable()
	if skill.LevelData != nil {
		for level, levelData := range skill.LevelData {
			if levelData != nil {
				levelTable := skillLevelDataToLuaTable(L, levelData)
				effectsTable.RawSetInt(level, levelTable)
			}
		}
	}
	tbl.RawSetString("effects", effectsTable)

	return tbl
}

// skillLevelDataToLuaTable converts wz.SkillLevelData to Lua table
func skillLevelDataToLuaTable(L *lua.LState, data *wz.SkillLevelData) *lua.LTable {
	tbl := L.NewTable()

	// Resource consumption
	tbl.RawSetString("mp_con", lua.LNumber(data.MPCon))
	tbl.RawSetString("hp_con", lua.LNumber(data.HPCon))
	tbl.RawSetString("money_con", lua.LNumber(data.MoneyCon))

	// Item consumption
	tbl.RawSetString("item_con", lua.LNumber(data.ItemCon))
	tbl.RawSetString("item_con_no", lua.LNumber(data.ItemConNo))
	tbl.RawSetString("item_consume", lua.LNumber(data.ItemConsume))
	tbl.RawSetString("bullet_consume", lua.LNumber(data.BulletConsume))
	tbl.RawSetString("bullet_count", lua.LNumber(data.BulletCount))

	// Damage and combat stats
	tbl.RawSetString("damage", lua.LNumber(data.Damage))
	tbl.RawSetString("damage_pc", lua.LNumber(data.DamagePC))
	tbl.RawSetString("fix_damage", lua.LNumber(data.FixDamage))
	tbl.RawSetString("critical_damage", lua.LNumber(data.CriticalDamage))
	tbl.RawSetString("attack_count", lua.LNumber(data.AttackCount))
	tbl.RawSetString("mob_count", lua.LNumber(data.MobCount))

	// Defense and evasion
	tbl.RawSetString("pad", lua.LNumber(data.PAD))
	tbl.RawSetString("mad", lua.LNumber(data.MAD))
	tbl.RawSetString("pdd", lua.LNumber(data.PDD))
	tbl.RawSetString("mdd", lua.LNumber(data.MDD))
	tbl.RawSetString("eva", lua.LNumber(data.EVA))
	tbl.RawSetString("acc", lua.LNumber(data.ACC))

	// Stat bonuses
	tbl.RawSetString("str", lua.LNumber(data.STR))
	tbl.RawSetString("hp", lua.LNumber(data.HP))
	tbl.RawSetString("mp", lua.LNumber(data.MP))
	tbl.RawSetString("jump", lua.LNumber(data.Jump))
	tbl.RawSetString("speed", lua.LNumber(data.Speed))

	// Skill properties
	tbl.RawSetString("mastery", lua.LNumber(data.Mastery))
	tbl.RawSetString("prop", lua.LNumber(data.Prop))
	tbl.RawSetString("range", lua.LNumber(data.Range))
	tbl.RawSetString("time", lua.LNumber(data.Time.Milliseconds()))
	tbl.RawSetString("cooldown", lua.LNumber(data.Cooldown.Milliseconds()))

	// Special properties
	tbl.RawSetString("morph", lua.LNumber(data.Morph))
	tbl.RawSetString("x", lua.LNumber(data.X))
	tbl.RawSetString("y", lua.LNumber(data.Y))
	tbl.RawSetString("z", lua.LNumber(data.Z))

	// Attack range (vectors)
	ltTable := L.NewTable()
	ltTable.RawSetString("x", lua.LNumber(data.LT.X))
	ltTable.RawSetString("y", lua.LNumber(data.LT.Y))
	tbl.RawSetString("lt", ltTable)

	rbTable := L.NewTable()
	rbTable.RawSetString("x", lua.LNumber(data.RB.X))
	rbTable.RawSetString("y", lua.LNumber(data.RB.Y))
	tbl.RawSetString("rb", rbTable)

	// String properties
	if data.HS != "" {
		tbl.RawSetString("hs", lua.LString(data.HS))
	}
	if data.Action != "" {
		tbl.RawSetString("action", lua.LString(data.Action))
	}

	return tbl
}
