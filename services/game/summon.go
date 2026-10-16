package game

import (
	"github.com/boyism80/fm/maple/constant"
	"github.com/boyism80/fm/protocol/response"
	"github.com/boyism80/fm/types"
)

type Summon struct {
	lifeCore
	owner        *Character
	ownerID      uint32
	skillID      constant.SkillID
	skillLevel   uint8
	movementType constant.SummonMovementType
	summonType   constant.SummonType
	changedMap   bool
}

func (s *Summon) GetObjectType() constant.ObjectType {
	return constant.ObjectTypeSummon
}

func (s *Summon) GetRole() constant.CharacterRole {
	if s.owner != nil {
		return s.owner.role
	}
	return constant.RoleUser
}

func (s *Summon) Is(typ constant.ObjectType) bool {
	return s.GetObjectType().Has(typ)
}

func (s *Summon) SendSpawnSyncToViewer(viewer *Character) {
	if s == nil || viewer == nil {
		return
	}
	viewer.Send(&response.SpawnSummon{
		OwnerID:      s.owner.GetID(),
		OID:          s.oid,
		SkillID:      s.skillID,
		SkillLevel:   s.skillLevel,
		Position:     s.position,
		MovementType: s.movementType,
		SummonType:   s.summonType,
		Animated:     false,
	}, types.SEND_POLICY_ENCRYPT)
}

func (s *Summon) SendDestroySyncToViewer(viewer *Character) {
	if s == nil || viewer == nil || s.owner == nil {
		return
	}
	viewer.Send(&response.RemoveSummon{
		OwnerID:  s.owner.GetID(),
		OID:      s.oid,
		Animated: false,
	}, types.SEND_POLICY_ENCRYPT)
}

type SummonAttackTarget struct {
	OID    uint32
	Damage uint32
}

func (s *Summon) Spawn(animated bool) {
	m := s.owner.GetMap()
	if m == nil {
		return
	}
	pos := s.owner.position
	s.position = pos
	m.AddSummon(s)
}

func (s *Summon) Remove(animated bool) {
	s.owner.summons.Remove(s, animated)
}

func (s *Summon) TakeDamage(unknown uint8, damage uint32, monsterIdFrom uint32) {
	if damage == 0 {
		return
	}

	if s.GetHp() <= damage {
		s.SetHp(0, false)
	} else {
		s.SetHp(s.GetHp()-damage, false)
	}
	s.owner.listener.OnSummonDamaged(s.owner, s, unknown, damage, monsterIdFrom)

	if s.GetHp() == 0 {
		s.owner.buffs.CancelBySource(int32(s.skillID))
	}
}
