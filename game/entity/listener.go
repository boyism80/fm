package entity

import (
	"github.com/boyism80/fm/game/constant"
	"github.com/boyism80/fm/game/wz"
	"github.com/boyism80/fm/protocol/dto"
	"github.com/boyism80/fm/types"
)

type CharacterListener interface {
	OnDialog(npc uint32, message string, prev bool, next bool)
	OnDialogYesNo(npc uint32, message string, prev bool, next bool)
	OnDialogAccept(npc uint32, message string, enableEscape bool)
	OnDialogList(npc uint32, message string, selections []string)
	OnDialogInput(npc uint32, message string)
	OnChat(message string, highlight bool, dontRecordHistory bool)
	OnMesoChanged(meso int32)
	OnMessage(messageType constant.ServerMessageType, message string)
	OnExpGain(exp uint32)
	OnControlMoveMob(oid uint32, moveId uint8, enabledSkill bool, mp uint16, skillId uint32, skillLevel uint8)
	OnShowMobHp(oid uint32, percentage uint8)
	OnUnlockAction()
	OnItemGainFailed(mode constant.ItemGainFailedType)
	OnInventorySlotUpdated(inventoryType constant.InventoryType, slot int16, item Item)
	OnInventorySlotAdded(inventoryType constant.InventoryType, slot int16, item Item)
	OnShowItemGain(itemId uint32, count uint32, mode constant.ShowItemGainType)
	OnShowMesoGain(count int32, mode constant.ShowMesoGainType)
	OnUpdateStats(stats map[constant.Stat]int32, unlock bool)
	OnMobMoved(mapID uint32, mobID uint32, isAggroed bool, centerSplit int8, skill1 uint8, skill2 uint8, skill3 uint8, skill4 uint8, startPoint types.Vector2[int16], movements []dto.MoveFragment)
	OnPlayerMove(mapID uint32, playerID uint32, character *Character, startPoint types.Vector2[int16], fragments []dto.MoveFragment)
	OnAttack(mapID uint32, characterID uint32, attackInfo dto.AttackInfo, skillLevel uint8)
	// New methods for remaining client.Send calls
	OnEndSortInventory(inventoryType constant.InventoryType)
	OnSwapInventorySlot(inventoryType constant.InventoryType, source int16, dest int16, equipmentAction int8)
	OnRemoveInventorySlot(inventoryType constant.InventoryType, slot int16)
	OnUpdateInventorySlot(inventoryType constant.InventoryType, slot int16, item Item)
	OnFullMergeInventorySlot(inventoryType constant.InventoryType, source int16, dest int16, count uint16)
	OnPartialMergeInventorySlot(inventoryType constant.InventoryType, source int16, dest int16, sourceCount uint16, destCount uint16)
	OnUpdateCharacterLook(character *Character)
	OnNpcAction(bytes []byte)
	OnClassChange(oldClass uint16, newClass uint16)
	OnBuffAdded(character *Character, wz *wz.Skill, level uint8, values map[constant.BuffFlag]int32)
	OnBuffRemoved(character *Character, flags []constant.BuffFlag)
	OnSkillCooldown(skillID uint32, remainingSec uint16)
	OnHiddenChanged(hidden bool)
}
