package server

import (
	"time"

	"github.com/boyism80/fm/game/constant"
	"github.com/boyism80/fm/game/entity"
	"github.com/boyism80/fm/game/wz"
	"github.com/boyism80/fm/protocol/dto"
	"github.com/boyism80/fm/protocol/response"
	"github.com/boyism80/fm/types"
)

type CharacterListenerImpl struct {
	gs *GameServer
	ch *entity.Character
}

func NewGameCharacterListener(gs *GameServer, ch *entity.Character) *CharacterListenerImpl {
	return &CharacterListenerImpl{
		gs: gs,
		ch: ch,
	}
}

func (l *CharacterListenerImpl) OnDialog(npc uint32, message string, prev bool, next bool) {
	dialogPacket := &response.Dialog{
		NPC:  npc,
		Type: constant.DIALOG_TYPE_DEFAULT,
		Text: message,
		Prev: prev,
		Next: next,
	}
	l.ch.Send(dialogPacket, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnDialogYesNo(npc uint32, message string, prev bool, next bool) {
	dialogPacket := &response.DialogYesNo{
		NPC:  npc,
		Text: message,
		Prev: prev,
		Next: next,
	}
	l.ch.Send(dialogPacket, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnDialogAccept(npc uint32, message string, enableEscape bool) {
	dialogPacket := &response.DialogAccept{
		NPC:          npc,
		Text:         message,
		EnableEscape: enableEscape,
	}
	l.ch.Send(dialogPacket, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnDialogList(npc uint32, message string, selections []string) {
	dialogPacket := &response.DialogList{
		NPC:        npc,
		Text:       message,
		Selections: selections,
	}
	l.ch.Send(dialogPacket, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnDialogInput(npc uint32, message string) {
	dialogPacket := &response.DialogInput{
		NPC:  npc,
		Text: message,
	}
	l.ch.Send(dialogPacket, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnChat(message string, highlight bool, dontRecordHistory bool) {
	chatPacket := &response.NormalChat{
		CharacterId:       l.ch.GetID(),
		Message:           message,
		Highlight:         highlight,
		DontRecordHistory: dontRecordHistory,
	}

	mapInstance := l.ch.GetMap()
	if mapInstance == nil {
		return
	}

	mapInstance.Broadcast(chatPacket, &entity.BroadcastOption{
		ReferenceCharacter: l.ch,
		RecipientFilter:    entity.BroadcastVisibleByReference,
	})
}

func (l *CharacterListenerImpl) OnMesoChanged(meso int32) {
	l.ch.Send(&response.UpdateStats{
		Stats: map[constant.Stat]int32{
			constant.STAT_MESO: meso,
		},
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnMessage(messageType constant.ServerMessageType, message string) {
	noticePacket := &response.Notice{
		Message: message,
		Type:    messageType,
	}
	l.ch.Send(noticePacket, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnExpGain(exp uint32) {
	expPacket := &response.GainExp{
		Gain:  exp,
		White: false,
	}
	l.ch.Send(expPacket, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnControlMoveMob(mob *entity.Mob, moveId uint8, enabledSkill bool, mp uint16, skillId uint32, skillLevel uint8) {
	l.ch.Send(&response.ControlMoveMob{
		OID:          mob.OID,
		MoveId:       uint16(moveId),
		EnabledSkill: enabledSkill,
		MP:           mp,
		SkillId:      uint8(skillId),
		SkillLevel:   skillLevel,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnShowMobHp(mob *entity.Mob, percentage uint8) {
	l.ch.Send(&response.ShowMobHp{
		OID:        mob.OID,
		Percentage: percentage,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnUnlockAction() {
	l.ch.Send(&response.UpdateStats{
		UnlockAction: true,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnItemGainFailed(mode constant.ItemGainFailedType) {
	l.ch.Send(&response.ItemGainFailed{
		Mode: mode,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnInventorySlotUpdated(inventoryType constant.InventoryType, slot int16, item entity.Item) {
	itemDTO := entity.ItemToDTO(item)
	l.ch.Send(&response.UpdateInventorySlot{
		InventoryType: inventoryType,
		Slot:          slot,
		Item:          itemDTO,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnInventorySlotAdded(inventoryType constant.InventoryType, slot int16, item entity.Item) {
	itemDTO := entity.ItemToDTO(item)

	// Determine FromDrop value based on item capacity
	// 0 = stackable (capacity >= 2), 1 = non-stackable (capacity < 2)
	fromDrop := true // default to non-stackable
	if item != nil {
		model := item.GetModel()
		if model != nil && model.GetCapacity() >= 2 {
			fromDrop = false // stackable
		}
	}

	l.ch.Send(&response.AddInventorySlot{
		InventoryType: inventoryType,
		Slot:          slot,
		Item:          itemDTO,
		FromDrop:      fromDrop,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnShowItemGain(itemId uint32, count uint32, mode constant.ShowItemGainType) {
	l.ch.Send(&response.ShowItemGain{
		ItemId: itemId,
		Count:  count,
		Mode:   mode,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnShowMesoGain(count int32, mode constant.ShowMesoGainType) {
	l.ch.Send(&response.ShowMesoGain{
		Count: count,
		Mode:  mode,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnUpdateStats(stats map[constant.Stat]int32, unlock bool) {
	l.ch.Send(&response.UpdateStats{
		Stats:        stats,
		UnlockAction: unlock,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnMobMoved(mob *entity.Mob, isAggroed bool, centerSplit int8, skill1 uint8, skill2 uint8, skill3 uint8, skill4 uint8, startPoint types.Vector2[int16], movements []dto.MoveFragment) {
	mapInstance := mob.GetObject().GetMap()
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

	var broadcastOption *entity.BroadcastOption
	if exists {
		broadcastOption = &entity.BroadcastOption{
			ExceptPlayerIDs: []uint32{controller.GetID()},
		}
	}

	mapInstance.Broadcast(movePacket, broadcastOption)
}

func (l *CharacterListenerImpl) OnPlayerMove(character *entity.Character, startPoint types.Vector2[int16], fragments []dto.MoveFragment) {
	mapInstance := character.GetMap()
	if mapInstance == nil {
		return
	}

	characterDTO := character.ToDTO()

	movePacket := &response.Move{
		Character:  characterDTO,
		Fragments:  fragments,
		StartPoint: startPoint,
	}

	mapInstance.Broadcast(movePacket, &entity.BroadcastOption{
		ExceptPlayerIDs:    []uint32{character.GetID()},
		ReferenceCharacter: character,
		RecipientFilter:    entity.BroadcastVisibleByReference,
	})
}

func (l *CharacterListenerImpl) OnAttack(character *entity.Character, attackInfo dto.AttackInfo, skillLevel uint8) {
	mapInstance := character.GetMap()
	if mapInstance == nil {
		return
	}

	attackPacket := &response.Attack{
		AttackInfo:  attackInfo,
		CharacterId: character.GetID(),
		SkillLevel:  skillLevel,
	}

	mapInstance.Broadcast(attackPacket, &entity.BroadcastOption{
		ExceptPlayerIDs:    []uint32{character.GetID()},
		ReferenceCharacter: character,
		RecipientFilter:    entity.BroadcastVisibleByReference,
	})
}

func (l *CharacterListenerImpl) OnEndSortInventory(inventoryType constant.InventoryType) {
	l.ch.Send(&response.EndSortInventory{
		InventoryType: inventoryType,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnSwapInventorySlot(inventoryType constant.InventoryType, source int16, dest int16, equipmentAction int8) {
	l.ch.Send(&response.SwapInventorySlot{
		InventoryType:   inventoryType,
		Source:          source,
		Dest:            dest,
		EquipmentAction: response.EquipmentActionType(equipmentAction),
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnRemoveInventorySlot(inventoryType constant.InventoryType, slot int16) {
	l.ch.Send(&response.RemoveInventorySlot{
		InventoryType: inventoryType,
		Slot:          slot,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnUpdateInventorySlot(inventoryType constant.InventoryType, slot int16, item entity.Item) {
	itemDTO := entity.ItemToDTO(item)
	l.ch.Send(&response.UpdateInventorySlot{
		InventoryType: inventoryType,
		Slot:          slot,
		Item:          itemDTO,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnFullMergeInventorySlot(inventoryType constant.InventoryType, source int16, dest int16, count uint16) {
	l.ch.Send(&response.FullMergeInventorySlot{
		InventoryType: inventoryType,
		Source:        source,
		Dest:          dest,
		Count:         count,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnPartialMergeInventorySlot(inventoryType constant.InventoryType, source int16, dest int16, sourceCount uint16, destCount uint16) {
	l.ch.Send(&response.PartialMergeInventorySlot{
		InventoryType: inventoryType,
		Source:        source,
		Dest:          dest,
		SourceCount:   sourceCount,
		DestCount:     destCount,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnUpdateCharacterLook(character *entity.Character) {
	mapInstance := character.GetMap()
	if mapInstance == nil {
		return
	}

	characterDTO := character.ToDTO()

	lookPacket := &response.UpdateCharacterLook{
		Character: characterDTO,
	}

	mapInstance.Broadcast(lookPacket, &entity.BroadcastOption{
		ExceptPlayerIDs:    []uint32{character.GetID()},
		ReferenceCharacter: character,
		RecipientFilter:    entity.BroadcastVisibleByReference,
	})
}

func (l *CharacterListenerImpl) OnNpcAction(bytes []byte) {
	l.ch.Send(&response.NpcAction{
		Bytes: bytes,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnClassChange(oldClass uint16, newClass uint16) {
	stats := map[constant.Stat]int32{
		constant.STAT_CLASS:        int32(newClass),
		constant.STAT_AVAILABLE_SP: int32(l.ch.SkillPoint),
	}

	l.ch.Send(&response.UpdateStats{
		Stats:        stats,
		UnlockAction: true,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnBuffAdded(character *entity.Character, wz *wz.Skill, level uint8, values map[constant.BuffFlag]int32) {
	if len(values) == 0 {
		return
	}
	dtoBuffs := make([]dto.BuffEntry, 0, len(values))
	for flag, value := range values {
		dtoBuff := dto.BuffEntry{Buff: flag, Value: value}
		dtoBuffs = append(dtoBuffs, dtoBuff)
	}

	duration := time.Duration(0)
	if wz != nil {
		duration = wz.GetLevelData(int(level)).Time
	}

	buffID := int32(0)
	if wz != nil {
		buffID = int32(wz.ID)
	}

	character.Send(&response.UpdateBuff{
		BuffID:   buffID,
		Duration: duration,
		Buffs:    dtoBuffs,
	}, types.SEND_POLICY_ENCRYPT)

	mapInstance := character.GetMap()
	if mapInstance != nil {
		mapInstance.Broadcast(&response.UpdateRemoteBuff{
			CharacterID: int32(character.GetID()),
			Buffs:       dtoBuffs,
		}, &entity.BroadcastOption{
			ExceptPlayerIDs:    []uint32{character.GetID()},
			ReferenceCharacter: character,
			RecipientFilter:    entity.BroadcastVisibleByReference,
		})
	}
}

func (l *CharacterListenerImpl) OnBuffRemoved(character *entity.Character, flags []constant.BuffFlag) {
	character.Send(&response.CancelBuff{Buffs: flags}, types.SEND_POLICY_ENCRYPT)

	mapInstance := character.GetMap()
	if mapInstance != nil {
		mapInstance.Broadcast(&response.CancelRemoteBuff{
			CharacterID: int32(character.GetID()),
			Buffs:       flags,
		}, &entity.BroadcastOption{
			ExceptPlayerIDs:    []uint32{character.GetID()},
			ReferenceCharacter: character,
			RecipientFilter:    entity.BroadcastVisibleByReference,
		})
	}
}

func (l *CharacterListenerImpl) OnSkillCooldown(skillID uint32, remainingSec uint16) {
	l.ch.Send(&response.SkillCooldown{
		SkillID:      skillID,
		RemainingSec: uint32(remainingSec),
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnHiddenChanged(hidden bool) {
	l.ch.Send(&response.SuperHide{Hidden: hidden}, types.SEND_POLICY_ENCRYPT)

	mapInstance := l.ch.GetMap()
	if mapInstance == nil {
		return
	}

	// Only players with lower role receive Leave/Spawn; same-or-higher role always see the character.
	if hidden {
		mapInstance.Broadcast(&response.LeavePlayer{ID: l.ch.GetID()}, &entity.BroadcastOption{
			ReferenceCharacter: l.ch,
			RecipientFilter:    entity.BroadcastRoleBelowReference,
		})
	} else {
		mapInstance.Broadcast(&response.SpawnPlayer{
			Character:       l.ch.ToDTO(),
			BuffStates:      [4]uint32{},
			Diseases:        [4]uint32{},
			CrushRings:      entity.RingsToDTO(l.ch.Rings.Left),
			FriendshipRings: entity.RingsToDTO(l.ch.Rings.Mid),
			MarriageRings:   entity.RingsToDTO(l.ch.Rings.Right),
		}, &entity.BroadcastOption{
			ReferenceCharacter: l.ch,
			RecipientFilter:    entity.BroadcastRoleBelowReference,
		})
	}
}
