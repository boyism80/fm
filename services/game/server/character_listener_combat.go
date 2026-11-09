package server

import (
	"time"

	"github.com/boyism80/fm/maple/constant"
	"github.com/boyism80/fm/protocol/dto"
	"github.com/boyism80/fm/protocol/response"
	"github.com/boyism80/fm/services/game"
	"github.com/boyism80/fm/types"
)

func (l *CharacterListenerImpl) OnControlMoveMob(ch *game.Character, mob *game.Mob, moveId uint16, enabledSkill bool, mp uint16, skillId uint32, skillLevel uint8) {
	ch.Send(&response.ControlMoveMob{
		OID:          mob.GetOID(),
		MoveId:       moveId,
		EnabledSkill: enabledSkill,
		MP:           mp,
		SkillId:      uint8(skillId),
		SkillLevel:   skillLevel,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnStopControlMob(ch *game.Character, mob *game.Mob) {
	ch.Send(&response.StopControlMob{OID: mob.GetOID()}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnShowMobHp(ch *game.Character, mob *game.Mob, percentage uint8) {
	ch.Send(&response.ShowMobHp{
		OID:        mob.GetOID(),
		Percentage: percentage,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnMobMoved(ch *game.Character, mob *game.Mob, isAggroed bool, centerSplit int8, skill1 uint8, skill2 uint8, skill3 uint8, skill4 uint8, startPoint types.Vector2[int16], movements []dto.MoveFragment) {
	mapInstance := mob.GetMap()
	if mapInstance == nil {
		return
	}

	movePacket := &response.MoveMob{
		IsAggroed:   isAggroed,
		CenterSplit: centerSplit,
		Skill1:      skill1,
		Skill2:      skill2,
		Skill3:      skill3,
		Skill4:      skill4,
		OID:         mob.GetOID(),
		StartPoint:  startPoint,
		Movements:   movements,
	}

	controllerID := uint32(0)
	if ch != nil {
		controllerID = ch.GetID()
	}
	for _, obj := range mob.Nears(constant.ObjectTypeCharacter, &game.SearchOption{
		IncludeHidden: true,
	}) {
		peer, ok := obj.(*game.Character)
		if !ok || peer == nil || peer.GetID() == controllerID {
			continue
		}
		_ = peer.Send(movePacket, types.SEND_POLICY_ENCRYPT)
	}
}

func (l *CharacterListenerImpl) OnAttack(ch *game.Character, attackPayload dto.AttackPayload, skillLevel uint8) {
	if ch.GetMap() == nil {
		return
	}
	ch.Broadcast(&response.Attack{
		AttackInfo:  attackPayload.ToAttackInfo(),
		CharacterId: ch.GetID(),
		SkillLevel:  skillLevel,
	}, nil)
}

func (l *CharacterListenerImpl) OnEnergyAttack(ch *game.Character, attackPayload dto.AttackPayload) {
	if ch.GetMap() == nil {
		return
	}
	ch.Broadcast(&response.EnergyAttack{
		Attack: response.Attack{
			AttackInfo:  attackPayload.ToAttackInfo(),
			CharacterId: ch.GetID(),
		},
	}, nil)
}

func (l *CharacterListenerImpl) OnRangedAttack(ch *game.Character, attackPayload dto.AttackPayload, skillLevel uint8, bullet uint32) {
	if ch.GetMap() == nil {
		return
	}
	ch.Broadcast(&response.RangedAttack{
		AttackInfo:  attackPayload.ToAttackInfo(),
		CharacterId: ch.GetID(),
		SkillLevel:  skillLevel,
		CashBullet:  bullet,
	}, nil)
}

func (l *CharacterListenerImpl) OnMagicAttack(ch *game.Character, attackPayload dto.AttackPayload, skillLevel uint8) {
	if ch.GetMap() == nil {
		return
	}
	ch.Broadcast(&response.MagicAttack{
		AttackInfo:  attackPayload.ToAttackInfo(),
		CharacterId: ch.GetID(),
		SkillLevel:  skillLevel,
	}, nil)
}

func (l *CharacterListenerImpl) OnBuffAdded(ch *game.Character, buffID int32, remainingDuration time.Duration, values map[constant.BuffFlag]int32) {
	if len(values) == 0 {
		return
	}
	dtoBuffs := make([]dto.BuffEntry, 0, len(values))
	for flag, value := range values {
		dtoBuff := dto.BuffEntry{Buff: flag, Value: value}
		dtoBuffs = append(dtoBuffs, dtoBuff)
	}

	var selfPacket types.Packet
	var remotePacket types.Packet
	if mountID, ok := values[constant.BuffFlagMonsterRiding]; ok {
		selfPacket = &response.UpdateSelfRidding{
			BuffID:  buffID,
			MountID: mountID,
			Buffs:   dtoBuffs,
		}
		remotePacket = &response.UpdateRidding{
			CharacterID: int32(ch.GetID()),
			MountID:     mountID,
			SkillID:     buffID,
			Buffs:       []dto.BuffEntry{{Buff: constant.BuffFlagMonsterRiding, Value: mountID}},
		}
	} else {
		selfPacket = &response.UpdateSelfBuff{
			BuffID:   buffID,
			Duration: remainingDuration,
			Buffs:    dtoBuffs,
		}
		remotePacket = &response.UpdateBuff{
			CharacterID: int32(ch.GetID()),
			BuffID:      buffID,
			Duration:    remainingDuration,
			Buffs:       dtoBuffs,
		}
	}

	ch.Send(selfPacket, types.SEND_POLICY_ENCRYPT)
	ch.Broadcast(remotePacket, nil)
}

func (l *CharacterListenerImpl) OnBuffRemoved(ch *game.Character, flags []constant.BuffFlag) {
	ch.Send(&response.CancelSelfBuff{Buffs: flags}, types.SEND_POLICY_ENCRYPT)

	ch.Broadcast(&response.CancelBuff{
		CharacterID: int32(ch.GetID()),
		Buffs:       flags,
	}, nil)
}

func (l *CharacterListenerImpl) OnDebuffAdded(ch *game.Character, debuff constant.DebuffFlag, x int16, skillID uint16, skillLevel uint16, durationMs int32) {
	ch.Send(&response.GiveSelfDebuff{
		Debuff:     debuff,
		X:          x,
		SkillID:    skillID,
		SkillLevel: skillLevel,
		DurationMs: durationMs,
	}, types.SEND_POLICY_ENCRYPT)

	ch.Broadcast(&response.GiveDebuff{
		CharacterID: int32(ch.GetID()),
		Debuff:      debuff,
		X:           x,
		SkillID:     skillID,
		SkillLevel:  skillLevel,
	}, nil)
}

func (l *CharacterListenerImpl) OnDebuffRemoved(ch *game.Character, flags []constant.DebuffFlag) {
	ch.Send(&response.RemoveSelfDebuff{Diseases: flags}, types.SEND_POLICY_ENCRYPT)

	ch.Broadcast(&response.RemoveDebuff{
		CharacterID: int32(ch.GetID()),
		Diseases:    flags,
	}, nil)
}

func (l *CharacterListenerImpl) OnUpdateSkill(ch *game.Character, skillID uint32, level int32, masterLevel int32) {
	ch.Send(&response.UpdateSkills{
		SkillID:     skillID,
		Level:       level,
		MasterLevel: masterLevel,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnSkillCooldown(ch *game.Character, skillID uint32, remainingSec uint16) {
	ch.Send(&response.SkillCooldown{
		SkillID:      skillID,
		RemainingSec: uint32(remainingSec),
	}, types.SEND_POLICY_ENCRYPT)
}
