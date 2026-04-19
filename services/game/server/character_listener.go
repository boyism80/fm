package server

import (
	"time"

	pconst "github.com/boyism80/fm/protocol/constant"
	"github.com/boyism80/fm/protocol/dto"
	"github.com/boyism80/fm/protocol/response"
	"github.com/boyism80/fm/services/game/constant"
	"github.com/boyism80/fm/services/game/entity"
	"github.com/boyism80/fm/types"
)

type CharacterListenerImpl struct {
	gs *GameServer
}

func (l *CharacterListenerImpl) OnDialog(ch *entity.Character, npc uint32, message string, prev bool, next bool) {
	dialogPacket := &response.Dialog{
		NPC:  npc,
		Type: constant.DIALOG_TYPE_DEFAULT,
		Text: message,
		Prev: prev,
		Next: next,
	}
	ch.Send(dialogPacket, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnDialogYesNo(ch *entity.Character, npc uint32, message string, prev bool, next bool) {
	dialogPacket := &response.DialogYesNo{
		NPC:  npc,
		Text: message,
		Prev: prev,
		Next: next,
	}
	ch.Send(dialogPacket, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnDialogAccept(ch *entity.Character, npc uint32, message string, enableEscape bool) {
	dialogPacket := &response.DialogAccept{
		NPC:          npc,
		Text:         message,
		EnableEscape: enableEscape,
	}
	ch.Send(dialogPacket, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnDialogList(ch *entity.Character, npc uint32, message string, selections []string) {
	dialogPacket := &response.DialogList{
		NPC:        npc,
		Text:       message,
		Selections: selections,
	}
	ch.Send(dialogPacket, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnDialogInput(ch *entity.Character, npc uint32, message string) {
	dialogPacket := &response.DialogInput{
		NPC:  npc,
		Text: message,
	}
	ch.Send(dialogPacket, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnChat(ch *entity.Character, message string, highlight bool, dontRecordHistory bool) {
	chatPacket := &response.NormalChat{
		CharacterId:       ch.GetID(),
		Message:           message,
		Highlight:         highlight,
		DontRecordHistory: dontRecordHistory,
	}

	if ch.GetMap() == nil {
		return
	}

	ch.Broadcast(chatPacket, &entity.ObjectBroadcastOption{WithMe: true})
}

func (l *CharacterListenerImpl) OnMesoChanged(ch *entity.Character, meso int32) {
	ch.Send(&response.UpdateStats{
		Stats: map[constant.Stat]int32{
			constant.STAT_MESO: meso,
		},
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnMessage(ch *entity.Character, messageType constant.ServerMessageType, message string) {
	noticePacket := &response.Notice{
		Message: message,
		Type:    messageType,
	}
	ch.Send(noticePacket, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnPartyCreated(ch *entity.Character, partyID uint32) {
	if ch == nil {
		return
	}
	ch.Send(&response.PartyCreated{
		PartyID: partyID,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnPartyStatusMessage(ch *entity.Character, code pconst.PartyStatusCode) {
	if ch == nil {
		return
	}
	ch.Send(&response.PartyStatusMessage{Code: code}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnExpGain(ch *entity.Character, exp uint32) {
	expPacket := &response.GainExp{
		Gain:  exp,
		White: false,
	}
	ch.Send(expPacket, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnControlMoveMob(ch *entity.Character, mob *entity.Mob, moveId uint16, enabledSkill bool, mp uint16, skillId uint32, skillLevel uint8) {
	ch.Send(&response.ControlMoveMob{
		OID:          mob.OID,
		MoveId:       moveId,
		EnabledSkill: enabledSkill,
		MP:           mp,
		SkillId:      uint8(skillId),
		SkillLevel:   skillLevel,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnShowMobHp(ch *entity.Character, mob *entity.Mob, percentage uint8) {
	ch.Send(&response.ShowMobHp{
		OID:        mob.OID,
		Percentage: percentage,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnUnlockAction(ch *entity.Character) {
	ch.Send(&response.UpdateStats{
		UnlockAction: true,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnItemGainFailed(ch *entity.Character, mode constant.ItemGainFailedType) {
	ch.Send(&response.ItemGainFailed{
		Mode: mode,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnInventorySlotUpdated(ch *entity.Character, inventoryType constant.InventoryType, slot int16, item entity.Item) {
	itemDTO := entity.ItemToDTO(item)
	ch.Send(&response.UpdateInventorySlot{
		InventoryType: inventoryType,
		Slot:          slot,
		Item:          itemDTO,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnInventorySlotAdded(ch *entity.Character, inventoryType constant.InventoryType, slot int16, item entity.Item) {
	itemDTO := entity.ItemToDTO(item)

	fromDrop := true
	if item != nil {
		model := item.GetModel()
		if model != nil && model.GetCapacity() >= 2 {
			fromDrop = false
		}
	}

	ch.Send(&response.AddInventorySlot{
		InventoryType: inventoryType,
		Slot:          slot,
		Item:          itemDTO,
		FromDrop:      fromDrop,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnShowItemGain(ch *entity.Character, itemId uint32, count uint32, mode constant.ShowItemGainType) {
	ch.Send(&response.ShowItemGain{
		ItemId: itemId,
		Count:  count,
		Mode:   mode,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnShowMesoGain(ch *entity.Character, count int32, mode constant.ShowMesoGainType) {
	ch.Send(&response.ShowMesoGain{
		Count: count,
		Mode:  mode,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnUpdateStats(ch *entity.Character, stats map[constant.Stat]int32, unlock bool) {
	ch.Send(&response.UpdateStats{
		Stats:        stats,
		UnlockAction: unlock,
	}, types.SEND_POLICY_ENCRYPT)
	if stats != nil {
		_, hasHP := stats[constant.STAT_HP]
		_, hasMaxHP := stats[constant.STAT_MAX_HP]
		if hasHP || hasMaxHP {
			l.OnPartyMemberHPChanged(ch, nil)
		}
	}
}

func (l *CharacterListenerImpl) OnShowBuffEffect(ch *entity.Character, effectID uint8, skillID uint32, skillLevel uint8, additional *uint8) {
	if ch.GetMap() == nil {
		return
	}

	ch.Send(&response.ShowOwnBuffeffect{
		EffectID:   effectID,
		SkillID:    skillID,
		SkillLevel: skillLevel,
		Additional: additional,
	}, types.SEND_POLICY_ENCRYPT)

	ch.Broadcast(&response.ShowBuffeffect{
		CharacterID: ch.GetID(),
		EffectID:    effectID,
		SkillID:     skillID,
		SkillLevel:  skillLevel,
		Additional:  additional,
	}, nil)
}

func (l *CharacterListenerImpl) OnMobMoved(ch *entity.Character, mob *entity.Mob, isAggroed bool, centerSplit int8, skill1 uint8, skill2 uint8, skill3 uint8, skill4 uint8, startPoint types.Vector2[int16], movements []dto.MoveFragment) {
	mapInstance := mob.GetMap()
	if mapInstance == nil {
		return
	}

	controllerTable := mapInstance.GetControllerTable()
	controller, exists := controllerTable.GetController(mob)

	movePacket := &response.MoveMob{
		IsAggroed:   isAggroed,
		CenterSplit: centerSplit,
		Skill1:      skill1,
		Skill2:      skill2,
		Skill3:      skill3,
		Skill4:      skill4,
		OID:         mob.OID,
		StartPoint:  startPoint,
		Movements:   movements,
	}

	if exists {
		controller.Broadcast(movePacket, nil)
	} else {
		mob.Broadcast(movePacket, nil)
	}
}

func (l *CharacterListenerImpl) OnPlayerMove(ch *entity.Character, startPoint types.Vector2[int16], fragments []dto.MoveFragment) {
	if ch.GetMap() == nil {
		return
	}

	characterDTO := ch.ToDTO()

	movePacket := &response.Move{
		Character:  characterDTO,
		Fragments:  fragments,
		StartPoint: startPoint,
	}

	ch.Broadcast(movePacket, nil)
}

func (l *CharacterListenerImpl) broadcastAttack(ch *entity.Character, packet types.Packet) {
	if ch.GetMap() == nil {
		return
	}

	ch.Broadcast(packet, nil)
}

func (l *CharacterListenerImpl) OnAttack(ch *entity.Character, attackPayload dto.AttackPayload, skillLevel uint8) {
	l.broadcastAttack(ch, &response.Attack{
		AttackInfo:  attackPayload.ToAttackInfo(),
		CharacterId: ch.GetID(),
		SkillLevel:  skillLevel,
	})
}

func (l *CharacterListenerImpl) OnRangedAttack(ch *entity.Character, attackPayload dto.AttackPayload, skillLevel uint8) {
	l.broadcastAttack(ch, &response.RangedAttack{
		AttackInfo:  attackPayload.ToAttackInfo(),
		CharacterId: ch.GetID(),
		SkillLevel:  skillLevel,
		CashBullet:  0,
	})
}

func (l *CharacterListenerImpl) OnMagicAttack(ch *entity.Character, attackPayload dto.AttackPayload, skillLevel uint8) {
	l.broadcastAttack(ch, &response.MagicAttack{
		AttackInfo:  attackPayload.ToAttackInfo(),
		CharacterId: ch.GetID(),
		SkillLevel:  skillLevel,
	})
}

func (l *CharacterListenerImpl) OnEndSortInventory(ch *entity.Character, inventoryType constant.InventoryType) {
	ch.Send(&response.EndSortInventory{
		InventoryType: inventoryType,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnSwapInventorySlot(ch *entity.Character, inventoryType constant.InventoryType, source int16, dest int16, equipmentAction int8) {
	ch.Send(&response.SwapInventorySlot{
		InventoryType:   inventoryType,
		Source:          source,
		Dest:            dest,
		EquipmentAction: response.EquipmentActionType(equipmentAction),
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnRemoveInventorySlot(ch *entity.Character, inventoryType constant.InventoryType, slot int16) {
	ch.Send(&response.RemoveInventorySlot{
		InventoryType: inventoryType,
		Slot:          slot,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnUpdateInventorySlot(ch *entity.Character, inventoryType constant.InventoryType, slot int16, item entity.Item) {
	itemDTO := entity.ItemToDTO(item)
	ch.Send(&response.UpdateInventorySlot{
		InventoryType: inventoryType,
		Slot:          slot,
		Item:          itemDTO,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnFullMergeInventorySlot(ch *entity.Character, inventoryType constant.InventoryType, source int16, dest int16, count uint16) {
	ch.Send(&response.FullMergeInventorySlot{
		InventoryType: inventoryType,
		Source:        source,
		Dest:          dest,
		Count:         count,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnPartialMergeInventorySlot(ch *entity.Character, inventoryType constant.InventoryType, source int16, dest int16, sourceCount uint16, destCount uint16) {
	ch.Send(&response.PartialMergeInventorySlot{
		InventoryType: inventoryType,
		Source:        source,
		Dest:          dest,
		SourceCount:   sourceCount,
		DestCount:     destCount,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnUpdateCharacterLook(ch *entity.Character) {
	mapInstance := ch.GetMap()
	if mapInstance == nil {
		return
	}

	characterDTO := ch.ToDTO()

	lookPacket := &response.UpdateCharacterLook{
		Character: characterDTO,
	}

	ch.Broadcast(lookPacket, nil)
}

func (l *CharacterListenerImpl) OnNpcAction(ch *entity.Character, bytes []byte) {
	ch.Send(&response.NpcAction{
		Bytes: bytes,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnClassChange(ch *entity.Character, oldClass uint16, newClass uint16) {
	stats := map[constant.Stat]int32{
		constant.STAT_CLASS:        int32(newClass),
		constant.STAT_AVAILABLE_SP: int32(ch.SkillPoint),
	}

	ch.Send(&response.UpdateStats{
		Stats:        stats,
		UnlockAction: true,
	}, types.SEND_POLICY_ENCRYPT)
	l.gs.UpdatePartyMemberAsync(ch)
}

func (l *CharacterListenerImpl) OnPartyMemberFieldsChanged(ch *entity.Character) {
	l.gs.UpdatePartyMemberAsync(ch)
}

func (l *CharacterListenerImpl) OnPartyMemberHPChanged(ch *entity.Character, recipient *entity.Character) {
	if ch == nil {
		return
	}
	partyIDPtr := ch.GetPartyID()
	if partyIDPtr == nil {
		return
	}
	mapInstance := ch.GetMap()
	if mapInstance == nil {
		return
	}
	pkt := &response.UpdatePartyMemberHP{
		CharacterID: ch.GetID(),
		CurrentHP:   int32(ch.GetHp()),
		MaxHP:       int32(ch.GetMaxHp()),
	}
	if recipient != nil {
		if recipient.GetID() == ch.GetID() {
			return
		}
		if recipient.GetMap() != mapInstance {
			return
		}
		rPID := recipient.GetPartyID()
		if rPID == nil || *rPID != *partyIDPtr {
			return
		}
		_ = recipient.Send(pkt, types.SEND_POLICY_ENCRYPT)
		return
	}
	for _, obj := range mapInstance.GetObjects(constant.ObjectTypeCharacter) {
		peer, ok := obj.(*entity.Character)
		if !ok || peer == nil || peer.GetID() == ch.GetID() {
			continue
		}
		pid := peer.GetPartyID()
		if pid == nil || *pid != *partyIDPtr {
			continue
		}
		_ = peer.Send(pkt, types.SEND_POLICY_ENCRYPT)
	}
}

func (l *CharacterListenerImpl) OnBuffAdded(ch *entity.Character, buffID int32, remainingDuration time.Duration, values map[constant.BuffFlag]int32) {
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
		selfPacket = &response.UpdateRidding{
			BuffID:  buffID,
			MountID: mountID,
			Buffs:   dtoBuffs,
		}
		remotePacket = &response.UpdateRemoteRidding{
			CharacterID: int32(ch.GetID()),
			MountID:     mountID,
			Buffs:       dtoBuffs,
		}
	} else {
		selfPacket = &response.UpdateBuff{
			BuffID:   buffID,
			Duration: remainingDuration,
			Buffs:    dtoBuffs,
		}
		remotePacket = &response.UpdateRemoteBuff{
			CharacterID: int32(ch.GetID()),
			BuffID:      buffID,
			Duration:    remainingDuration,
			Buffs:       dtoBuffs,
		}
	}

	ch.Send(selfPacket, types.SEND_POLICY_ENCRYPT)
	ch.Broadcast(remotePacket, nil)
}

func (l *CharacterListenerImpl) OnBuffRemoved(ch *entity.Character, flags []constant.BuffFlag) {
	ch.Send(&response.CancelBuff{Buffs: flags}, types.SEND_POLICY_ENCRYPT)

	ch.Broadcast(&response.CancelRemoteBuff{
		CharacterID: int32(ch.GetID()),
		Buffs:       flags,
	}, nil)
}

func (l *CharacterListenerImpl) OnDebuffAdded(ch *entity.Character, disease constant.DebuffFlag, x int16, skillID uint16, skillLevel uint16, durationMs int32) {
	ch.Send(&response.GiveDebuff{
		Disease:    disease,
		X:          x,
		SkillID:    skillID,
		SkillLevel: skillLevel,
		DurationMs: durationMs,
	}, types.SEND_POLICY_ENCRYPT)

	ch.Broadcast(&response.GiveRemoteDebuff{
		CharacterID: int32(ch.GetID()),
		Disease:     disease,
		X:           x,
		SkillID:     skillID,
		SkillLevel:  skillLevel,
	}, nil)
}

func (l *CharacterListenerImpl) OnDebuffRemoved(ch *entity.Character, flags []constant.DebuffFlag) {
	ch.Send(&response.RemoveDebuff{Diseases: flags}, types.SEND_POLICY_ENCRYPT)

	ch.Broadcast(&response.RemoveRemoteDebuff{
		CharacterID: int32(ch.GetID()),
		Diseases:    flags,
	}, nil)
}

func (l *CharacterListenerImpl) OnSkillPassiveHook(ch *entity.Character, skillID uint32, hook string) {
	if ch == nil {
		return
	}
	CallPassiveSkillHook(nil, ch, skillID, hook)
}

func (l *CharacterListenerImpl) OnUpdateSkill(ch *entity.Character, skillID uint32, level int32, masterLevel int32) {
	if ch == nil {
		return
	}
	ch.Send(&response.UpdateSkills{
		SkillID:     skillID,
		Level:       level,
		MasterLevel: masterLevel,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnSkillCooldown(ch *entity.Character, skillID uint32, remainingSec uint16) {
	ch.Send(&response.SkillCooldown{
		SkillID:      skillID,
		RemainingSec: uint32(remainingSec),
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnHiddenChanged(ch *entity.Character, hidden bool) {
	ch.Send(&response.SuperHide{Hidden: hidden}, types.SEND_POLICY_ENCRYPT)

	mapInstance := ch.GetMap()
	if mapInstance == nil {
		return
	}

	if hidden {
		ch.Broadcast(&response.LeavePlayer{ID: ch.GetID()}, &entity.ObjectBroadcastOption{
			RecipientsRoleBelowPivot: true,
		})
	} else {
		for _, obj := range mapInstance.GetObjects(constant.ObjectTypeCharacter) {
			if viewer, ok := obj.(*entity.Character); ok {
				ch.SendSpawnSyncToViewer(viewer)
			}
		}
	}
}

func (l *CharacterListenerImpl) OnSummonSpawn(ch *entity.Character, summon *entity.Summon) {
	if summon.GetMap() == nil {
		return
	}
	packet := &response.SpawnSummon{
		OwnerID:      ch.GetID(),
		OID:          summon.OID,
		SkillID:      summon.SkillID,
		SkillLevel:   summon.SkillLevel,
		Position:     summon.Position,
		MovementType: summon.MovementType,
		SummonType:   summon.SummonType,
		Animated:     true,
	}
	ch.Broadcast(packet, &entity.ObjectBroadcastOption{WithMe: true})
}

func (l *CharacterListenerImpl) OnSummonRemove(ch *entity.Character, summon *entity.Summon, animated bool) {
	if summon.GetMap() == nil {
		return
	}
	packet := &response.RemoveSummon{
		OwnerID:  ch.GetID(),
		OID:      summon.OID,
		Animated: animated,
	}
	ch.Broadcast(packet, &entity.ObjectBroadcastOption{WithMe: true})
}

func (l *CharacterListenerImpl) OnSummonMove(ch *entity.Character, summon *entity.Summon, startPoint types.Vector2[int16], movements []dto.MoveFragment) {
	if summon.GetMap() == nil {
		return
	}
	packet := &response.MoveSummon{
		CharacterID: ch.GetID(),
		OID:         summon.OID,
		StartPoint:  startPoint,
		Fragments:   movements,
	}
	ch.Broadcast(packet, nil)
}

func (l *CharacterListenerImpl) OnSummonAttack(ch *entity.Character, summon *entity.Summon, animation uint8, targets []entity.SummonAttackTarget) {
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
		CharacterID:   ch.GetID(),
		SummonSkillID: uint32(summon.SkillID),
		Animation:     animation,
		Targets:       respTargets,
	}
	ch.Broadcast(packet, &entity.ObjectBroadcastOption{WithMe: true})
}

func (l *CharacterListenerImpl) OnSummonSkill(ch *entity.Character, summon *entity.Summon, newStance uint8) {
	if summon.GetMap() == nil {
		return
	}
	packet := &response.SummonSkill{
		CharacterID: ch.GetID(),
		SummonOID:   summon.OID,
		NewStance:   newStance,
	}
	ch.Broadcast(packet, &entity.ObjectBroadcastOption{WithMe: true})
}

func (l *CharacterListenerImpl) OnSummonDamaged(ch *entity.Character, summon *entity.Summon, unknown uint8, damage uint32, monsterIdFrom uint32) {
	if summon.GetMap() == nil {
		return
	}
	packet := &response.DamageSummon{
		CharacterID:   ch.GetID(),
		SummonSkillID: uint32(summon.SkillID),
		Unknown:       unknown,
		Damage:        damage,
		MonsterIDFrom: monsterIdFrom,
	}
	ch.Broadcast(packet, &entity.ObjectBroadcastOption{WithMe: true})
}
