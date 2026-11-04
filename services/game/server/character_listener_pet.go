package server

import (
	"github.com/boyism80/fm/maple/constant"
	"github.com/boyism80/fm/protocol/dto"
	"github.com/boyism80/fm/protocol/response"
	"github.com/boyism80/fm/services/game"
	"github.com/boyism80/fm/types"
)

func (l *CharacterListenerImpl) OnPetSpawn(ch *game.Character, sn uint64) {
	ch.Broadcast(&response.SpawnPet{
		CharacterID: ch.GetID(),
		Pet:         ch.Pets().Active().BuildS2C(),
	}, &game.ObjectBroadcastOption{WithMe: true})
	ch.Pets().ShowExceptions()
	ch.Send(&response.UpdateStats{Pet: &sn, UnlockAction: true}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnPetRemove(ch *game.Character, reason constant.PetRemoveReason) {
	ch.Broadcast(&response.SpawnPet{
		CharacterID: ch.GetID(),
		Reason:      reason,
	}, &game.ObjectBroadcastOption{WithMe: true})
	sn := uint64(0)
	ch.Send(&response.UpdateStats{Pet: &sn, UnlockAction: true}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnPetMove(ch *game.Character, startPoint types.Vector2[int16], movements []dto.MoveFragment) {
	ch.Broadcast(&response.MovePet{
		CharacterID: ch.GetID(),
		StartPoint:  startPoint,
		Fragments:   movements,
	}, nil)
}

func (l *CharacterListenerImpl) OnPetChat(ch *game.Character, typ uint8, action uint8, text string) {
	ch.Broadcast(&response.PetChat{
		CharacterID: ch.GetID(),
		Type:        typ,
		Action:      action,
		Text:        text,
	}, &game.ObjectBroadcastOption{WithMe: true})
}

func (l *CharacterListenerImpl) OnPetCommand(ch *game.Character, index uint8, success bool) {
	ch.Broadcast(&response.PetCommand{CharacterID: ch.GetID(), Index: index, Success: success}, &game.ObjectBroadcastOption{WithMe: true})
}

func (l *CharacterListenerImpl) OnPetFood(ch *game.Character, success bool) {
	ch.Broadcast(&response.PetCommand{CharacterID: ch.GetID(), Food: true, Success: success}, &game.ObjectBroadcastOption{WithMe: true})
}

func (l *CharacterListenerImpl) OnPetUpdated(ch *game.Character, pet *game.Pet) {
	slot, ok := ch.Inventory().FindSlot(constant.InventoryTypeCash, pet)
	if !ok {
		return
	}
	l.OnInventorySlotUpdated(ch, constant.InventoryTypeCash, slot, pet)
}

func (l *CharacterListenerImpl) OnMountUpdated(ch *game.Character, level uint32, exp uint32, fatigue uint32, levelUp bool) {
	ch.Broadcast(&response.UpdateMount{
		CharacterID: ch.GetID(),
		Level:       level,
		Exp:         exp,
		Fatigue:     fatigue,
		LevelUp:     levelUp,
	}, &game.ObjectBroadcastOption{WithMe: true})
}

func (l *CharacterListenerImpl) OnPetExceptions(ch *game.Character, sn uint64, itemIDs []uint32) {
	ch.Send(&response.PetExceptions{
		CharacterID: ch.GetID(),
		SN:          sn,
		ItemIDs:     itemIDs,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnPetNameChanged(ch *game.Character, name string) {
	ch.Broadcast(&response.PetNameChanged{
		CharacterID: ch.GetID(),
		Name:        name,
	}, &game.ObjectBroadcastOption{WithMe: true})
}

func (l *CharacterListenerImpl) OnPetSkillChanged(ch *game.Character, pet *game.Pet, sn uint64, skill constant.PetSkill, add bool) {
	l.OnPetUpdated(ch, pet)
	ch.Send(&response.PetSkillChanged{SN: sn, Add: add, Skill: skill}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnPetLevelUp(ch *game.Character) {
	ch.Send(&response.ShowSelfEffect{Type: response.EffectTypePet}, types.SEND_POLICY_ENCRYPT)
	ch.Broadcast(&response.ShowEffect{CharacterID: ch.GetID(), Type: response.EffectTypePet}, nil)
}

func (l *CharacterListenerImpl) OnSummonSpawn(ch *game.Character, summon *game.Summon) {
	if summon.GetMap() == nil {
		return
	}
	packet := &response.SpawnSummon{
		OwnerID:      ch.GetID(),
		OID:          summon.GetOID(),
		SkillID:      summon.SkillID(),
		SkillLevel:   summon.SkillLevel(),
		Position:     summon.GetPosition(),
		MovementType: summon.MovementType(),
		SummonType:   summon.SummonType(),
		Animated:     true,
	}
	ch.Broadcast(packet, &game.ObjectBroadcastOption{WithMe: true})
}

func (l *CharacterListenerImpl) OnSummonRemove(ch *game.Character, summon *game.Summon, animated bool) {
	if summon.GetMap() == nil {
		return
	}
	packet := &response.RemoveSummon{
		OwnerID:  ch.GetID(),
		OID:      summon.GetOID(),
		Animated: animated,
	}
	ch.Broadcast(packet, &game.ObjectBroadcastOption{WithMe: true})
}

func (l *CharacterListenerImpl) OnSummonMove(ch *game.Character, summon *game.Summon, startPoint types.Vector2[int16], movements []dto.MoveFragment) {
	if summon.GetMap() == nil {
		return
	}
	packet := &response.MoveSummon{
		CharacterID: ch.GetID(),
		OID:         summon.GetOID(),
		StartPoint:  startPoint,
		Fragments:   movements,
	}
	ch.Broadcast(packet, nil)
}

func (l *CharacterListenerImpl) OnSummonAttack(ch *game.Character, summon *game.Summon, animation uint8, targets []game.SummonAttackTarget) {
	if summon.GetMap() == nil {
		return
	}
	respTargets := make([]response.SummonAttackTarget, len(targets))
	for i, t := range targets {
		respTargets[i] = response.SummonAttackTarget{
			OID:    t.OID,
			Damage: t.Damage,
		}
	}
	packet := &response.SummonAttack{
		CharacterID: ch.GetID(),
		SummonOID:   summon.GetOID(),
		Animation:   animation,
		Targets:     respTargets,
	}
	ch.Broadcast(packet, nil)
}

func (l *CharacterListenerImpl) OnSummonSkill(ch *game.Character, summon *game.Summon, newStance uint8) {
	if summon.GetMap() == nil {
		return
	}
	packet := &response.SummonSkill{
		CharacterID: ch.GetID(),
		SummonOID:   summon.GetOID(),
		NewStance:   newStance,
	}
	ch.Broadcast(packet, &game.ObjectBroadcastOption{WithMe: true})
}

func (l *CharacterListenerImpl) OnSummonDamaged(ch *game.Character, summon *game.Summon, unknown uint8, damage uint32, monsterIdFrom uint32) {
	if summon.GetMap() == nil {
		return
	}
	packet := &response.DamageSummon{
		CharacterID:   ch.GetID(),
		SummonOID:     summon.GetOID(),
		Unknown:       unknown,
		Damage:        damage,
		MonsterIDFrom: monsterIdFrom,
	}
	ch.Broadcast(packet, nil)
}
