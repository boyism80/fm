package game

import (
	"github.com/boyism80/fm/core/luax"
	"github.com/boyism80/fm/types"
	lua "github.com/yuin/gopher-lua"
)

func (s *Summon) LuaTypeName() string {
	return "LuaSummon"
}

func (s *Summon) LuaBuiltinFuncs() map[string]lua.LGFunction {
	return map[string]lua.LGFunction{
		"owner_id": func(L *lua.LState) int {
			summon := luax.Check[Summon](L, 1)
			if L.GetTop() != 1 {
				L.ArgError(2, "owner_id() is read-only")
				return 0
			}
			L.Push(lua.LNumber(summon.ownerID))
			return 1
		},
		"skill_id": func(L *lua.LState) int {
			summon := luax.Check[Summon](L, 1)
			if L.GetTop() != 1 {
				L.ArgError(2, "skill_id() is read-only")
				return 0
			}
			L.Push(lua.LNumber(summon.skillID))
			return 1
		},
		"movement_type": func(L *lua.LState) int {
			summon := luax.Check[Summon](L, 1)
			if L.GetTop() != 1 {
				L.ArgError(2, "movement_type() is read-only")
				return 0
			}
			L.Push(lua.LNumber(summon.movementType))
			return 1
		},
		"attack": func(L *lua.LState) int {
			summon := luax.Check[Summon](L, 1)
			animation := uint8(L.CheckInt(2))
			targetsTable := L.CheckTable(3)
			var targets []SummonAttackTarget
			targetsTable.ForEach(func(_ lua.LValue, v lua.LValue) {
				t, ok := v.(*lua.LTable)
				if !ok {
					return
				}
				oidLV := t.RawGetString("oid")
				dmgLV := t.RawGetString("damage")
				if oidLV.Type() != lua.LTNumber || dmgLV.Type() != lua.LTNumber {
					return
				}
				targets = append(targets, SummonAttackTarget{
					OID:    uint32(lua.LVAsNumber(oidLV)),
					Damage: uint32(lua.LVAsNumber(dmgLV)),
				})
			})
			summon.owner.listener.OnSummonAttack(summon.owner, summon, animation, targets)
			return 0
		},
		"move": func(L *lua.LState) int {
			summon := luax.Check[Summon](L, 1)
			x := int16(L.CheckInt(2))
			y := int16(L.CheckInt(3))
			start := types.Vector2[int16]{X: x, Y: y}
			summon.owner.listener.OnSummonMove(summon.owner, summon, start, nil)
			return 0
		},
		"use_skill": func(L *lua.LState) int {
			summon := luax.Check[Summon](L, 1)
			if L.GetTop() != 2 {
				L.ArgError(2, "use_skill() requires stance")
				return 0
			}
			newStance := uint8(L.CheckInt(2))
			summon.owner.listener.OnSummonSkill(summon.owner, summon, newStance)
			return 0
		},
	}
}
