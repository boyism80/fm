package entity

import (
	"github.com/boyism80/fm/protocol/dto"
	"github.com/boyism80/fm/protocol/response"
	"github.com/boyism80/fm/services/game/constant"
	"github.com/boyism80/fm/types"
)

type Summon struct {
	LifeCore
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

func (s *Summon) GetRole() constant.CharacterRole {
	if s.Owner != nil {
		return s.Owner.Role
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
		OwnerID:      s.Owner.GetID(),
		OID:          s.OID,
		SkillID:      s.SkillID,
		SkillLevel:   s.SkillLevel,
		Position:     s.Position,
		MovementType: s.MovementType,
		SummonType:   s.SummonType,
		Animated:     false,
	}, types.SEND_POLICY_ENCRYPT)
}

type SummonAttackTarget struct {
	OID    uint32
	Damage uint32
}

func (s *Summon) Spawn(animated bool) {
	m := s.Owner.GetMap()
	if m == nil {
		return
	}
	pos := s.Owner.Position
	s.Position = pos
	m.AddSummon(s)
}

func (s *Summon) Remove(animated bool) {
	s.Owner.RemoveSummon(s, animated)
}

func (s *Summon) Move(start types.Vector2[int16], movements []dto.MoveFragment) {
	s.Owner.Listener.OnSummonMove(s.Owner, s, start, movements)
}

func (s *Summon) Attack(animation uint8, targets []SummonAttackTarget) {
	s.Owner.Listener.OnSummonAttack(s.Owner, s, animation, targets)
}

func (s *Summon) UseSkill(newStance uint8) {
	s.Owner.Listener.OnSummonSkill(s.Owner, s, newStance)
}

func (s *Summon) TakeDamage(unknown uint8, damage uint32, monsterIdFrom uint32) {
	s.Owner.Listener.OnSummonDamaged(s.Owner, s, unknown, damage, monsterIdFrom)
}
