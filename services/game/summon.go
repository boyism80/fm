package game

import (
	"fmt"
	"log"

	"github.com/boyism80/fm/core/luax"
	"github.com/boyism80/fm/maple/constant"
	"github.com/boyism80/fm/protocol/dto"
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

func (s *Summon) Move(fragments []dto.MoveFragment) {
	if s.movementType == constant.SummonMoveStationary {
		return
	}
	start := s.position
	for _, fragment := range fragments {
		if move, ok := fragment.(*dto.AbsoluteLifeMovement); ok {
			s.Place(move.Position)
		}
		s.SetStance(fragment.GetStance())
	}
	s.onMap.OnMoved(s, start)
	s.owner.listener.OnSummonMove(s.owner, s, start, fragments)
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

// CallAttackHook runs after the attack broadcast, so a summon that ends in its own hook is still on the map for it.
func (s *Summon) CallAttackHook(damages []dto.AttackPair) {
	m := s.owner.GetMap()
	if m == nil || s.skillID == 0 {
		return
	}
	root := m.GetLuaRoot()
	if root == nil {
		return
	}
	skill := s.owner.skills.Get(uint32(s.skillID))
	if skill == nil {
		return
	}
	thread, err := luax.NewThread(root, fmt.Sprintf("script/skill/%d.lua", s.skillID))
	if err != nil {
		log.Printf("summon on_attack thread %d: %v", s.skillID, err)
		return
	}
	if _, err := luax.Call(thread, "on_attack", s.owner, skill, s.owner.damagesTable(thread, damages)); err != nil {
		log.Printf("summon on_attack %d: %v", s.skillID, err)
	}
}
