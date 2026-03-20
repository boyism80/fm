package entity

import (
	"github.com/boyism80/fm/game/constant"
	"github.com/boyism80/fm/protocol/dto"
	"github.com/boyism80/fm/types"
	lua "github.com/yuin/gopher-lua"
)

type Summon struct {
	Life
	Owner        *Character
	OwnerID      uint32
	SkillID      constant.SkillID
	SkillLevel   uint8
	MovementType constant.SummonMovementType
	SummonType   constant.SummonType
	ChangedMap   bool
}

func (s *Summon) GetObjectType() constant.ObjectType {
	return constant.ObjectTypeSummon
}

func (s *Summon) Is(typ constant.ObjectType) bool {
	return s.GetObjectType().Has(typ)
}

func (s *Summon) GetObject() *Object {
	return &s.Life.Object
}

func (s *Summon) LuaTypeName() string {
	return "LuaSummon"
}

func (s *Summon) LuaBuiltinFuncs() map[string]lua.LGFunction {
	return map[string]lua.LGFunction{
		"owner_id": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			summon, ok := ud.Value.(*Summon)
			if !ok {
				L.ArgError(1, "Summon expected")
				return 0
			}
			if L.GetTop() != 1 {
				L.ArgError(2, "owner_id() is read-only")
				return 0
			}
			L.Push(lua.LNumber(summon.OwnerID))
			return 1
		},
		"skill_id": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			summon, ok := ud.Value.(*Summon)
			if !ok {
				L.ArgError(1, "Summon expected")
				return 0
			}
			if L.GetTop() != 1 {
				L.ArgError(2, "skill_id() is read-only")
				return 0
			}
			L.Push(lua.LNumber(summon.SkillID))
			return 1
		},
		"attack": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			summon, ok := ud.Value.(*Summon)
			if !ok {
				L.ArgError(1, "Summon expected")
				return 0
			}
			if summon.Owner == nil || summon.Owner.Listener == nil {
				return 0
			}
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
			summon.Owner.Listener.OnSummonAttack(summon.Owner, summon, animation, targets)
			return 0
		},
		"move": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			summon, ok := ud.Value.(*Summon)
			if !ok {
				L.ArgError(1, "Summon expected")
				return 0
			}
			if summon.Owner == nil || summon.Owner.Listener == nil {
				return 0
			}
			x := int16(L.CheckInt(2))
			y := int16(L.CheckInt(3))
			start := types.Vector2[int16]{X: x, Y: y}
			summon.Owner.Listener.OnSummonMove(summon.Owner, summon, start, nil)
			return 0
		},
		"use_skill": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			summon, ok := ud.Value.(*Summon)
			if !ok {
				L.ArgError(1, "Summon expected")
				return 0
			}
			if L.GetTop() != 2 {
				L.ArgError(2, "use_skill() requires stance")
				return 0
			}
			newStance := uint8(L.CheckInt(2))
			summon.UseSkill(newStance)
			return 0
		},
	}
}

type SummonAttackTarget struct {
	OID    uint32
	Damage uint32
}

func (s *Summon) Spawn(animated bool) {
	if s.Owner == nil {
		return
	}
	m := s.Owner.GetMap()
	if m == nil {
		return
	}
	pos := s.Owner.Position
	_ = m.SpawnSummon(s.Owner, constant.SkillID(s.SkillID), s.SkillLevel, s.MovementType, s.SummonType, pos, 0)
}

func (s *Summon) Remove(animated bool) {
	if s.Owner == nil || s.Owner.Listener == nil {
		return
	}
	s.Owner.RemoveSummon(s)
	s.Owner.Listener.OnSummonRemove(s.Owner, s, animated)
}

func (s *Summon) Move(start types.Vector2[int16], movements []dto.MoveFragment) {
	if s.Owner == nil || s.Owner.Listener == nil {
		return
	}
	s.Owner.Listener.OnSummonMove(s.Owner, s, start, movements)
}

func (s *Summon) Attack(animation uint8, targets []SummonAttackTarget) {
	if s.Owner == nil || s.Owner.Listener == nil {
		return
	}
	s.Owner.Listener.OnSummonAttack(s.Owner, s, animation, targets)
}

func (s *Summon) UseSkill(newStance uint8) {
	if s.Owner == nil || s.Owner.Listener == nil {
		return
	}
	s.Owner.Listener.OnSummonSkill(s.Owner, s, newStance)
}

func (s *Summon) TakeDamage(unknown uint8, damage uint32, monsterIdFrom uint32) {
	if s.Owner == nil || s.Owner.Listener == nil {
		return
	}
	s.Owner.Listener.OnSummonDamaged(s.Owner, s, unknown, damage, monsterIdFrom)
}
