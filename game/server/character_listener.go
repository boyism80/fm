package server

import (
	"github.com/boyism80/fm/core/types"
	"github.com/boyism80/fm/game/action"
	"github.com/boyism80/fm/game/constant"
	"github.com/boyism80/fm/game/entity"
	"github.com/boyism80/fm/game/protocol/resp"
)

// CharacterListenerImpl implements CharacterListener for game server
type CharacterListenerImpl struct {
	gameServer *GameServer
	character  *entity.Character
}

// NewGameCharacterListener creates a new GameCharacterListener instance
func NewGameCharacterListener(gameServer *GameServer, character *entity.Character) *CharacterListenerImpl {
	return &CharacterListenerImpl{
		gameServer: gameServer,
		character:  character,
	}
}

// OnDialog handles dialog events
func (l *CharacterListenerImpl) OnDialog(npc uint32, message string, prev bool, next bool) {
	// Send dialog packet to client
	dialogPacket := &resp.Dialog{
		NPC:  npc,
		Type: constant.DIALOG_TYPE_DEFAULT,
		Text: message,
		Prev: prev,
		Next: next,
	}
	l.character.Send(dialogPacket, types.SEND_POLICY_ENCRYPT)
}

// OnDialogYesNo handles yes/no dialog events
func (l *CharacterListenerImpl) OnDialogYesNo(npc uint32, message string, prev bool, next bool) {
	// Send yes/no dialog packet to client
	dialogPacket := &resp.DialogYesNo{
		NPC:  npc,
		Text: message,
		Prev: prev,
		Next: next,
	}
	l.character.Send(dialogPacket, types.SEND_POLICY_ENCRYPT)
}

// OnDialogAccept handles dialog accept events
func (l *CharacterListenerImpl) OnDialogAccept(npc uint32, message string, enableEscape bool) {
	// Send dialog accept packet to client
	dialogPacket := &resp.DialogAccept{
		NPC:          npc,
		Text:         message,
		EnableEscape: enableEscape,
	}
	l.character.Send(dialogPacket, types.SEND_POLICY_ENCRYPT)
}

// OnDialogList handles dialog list events
func (l *CharacterListenerImpl) OnDialogList(npc uint32, message string, selections []string) {
	// Send dialog list packet to client
	dialogPacket := &resp.DialogList{
		NPC:        npc,
		Text:       message,
		Selections: selections,
	}
	l.character.Send(dialogPacket, types.SEND_POLICY_ENCRYPT)
}

// OnDialogInput handles dialog input events
func (l *CharacterListenerImpl) OnDialogInput(npc uint32, message string) {
	// Send dialog input packet to client
	dialogPacket := &resp.DialogInput{
		NPC:  npc,
		Text: message,
	}
	l.character.Send(dialogPacket, types.SEND_POLICY_ENCRYPT)
}

// OnChat handles chat events by broadcasting chat packet to all players on the map
func (l *CharacterListenerImpl) OnChat(message string, highlight bool, dontRecordHistory bool) {
	// Create chat packet
	chatPacket := &resp.NormalChat{
		CharacterId:       l.character.ID,
		Message:           message,
		Highlight:         highlight,
		DontRecordHistory: dontRecordHistory,
	}

	// Get map instance
	mapInstance := l.gameServer.GetMap(l.character.Map)
	if mapInstance == nil {
		return
	}

	mapInstance.BroadcastToAllPlayers(chatPacket, types.SEND_POLICY_ENCRYPT)
}

// OnMesoChanged handles meso change events
func (l *CharacterListenerImpl) OnMesoChanged(meso int32) {
	// TODO: Implement meso change handling
	// This could involve sending a packet to update the client's meso display
	l.character.Send(&resp.UpdateStats{
		Stats: map[constant.Stat]int32{
			constant.STAT_MESO: meso,
		},
	}, types.SEND_POLICY_ENCRYPT)
}

// OnMessage handles general message events
func (l *CharacterListenerImpl) OnMessage(messageType constant.ServerMessageType, message string) {
	// Send notice packet to client
	noticePacket := &resp.Notice{
		Message: message,
		Type:    messageType,
	}
	l.character.Send(noticePacket, types.SEND_POLICY_ENCRYPT)
}

// OnExpGain handles experience gain events
func (l *CharacterListenerImpl) OnExpGain(exp uint32) {
	// Send experience gain packet to client
	expPacket := &resp.GainExp{
		Gain:  exp,
		White: false,
	}
	l.character.Send(expPacket, types.SEND_POLICY_ENCRYPT)
}

// OnControlMoveMob handles control move mob events
func (l *CharacterListenerImpl) OnControlMoveMob(oid uint32, moveId uint8, enabledSkill bool, mp uint16, skillId uint32, skillLevel uint8) {
	l.character.Send(&resp.ControlMoveMob{
		OID:          oid,
		MoveId:       uint16(moveId),
		EnabledSkill: enabledSkill,
		MP:           mp,
		SkillId:      uint8(skillId),
		SkillLevel:   skillLevel,
	}, types.SEND_POLICY_ENCRYPT)
}

// OnShowMobHp handles mob HP display events
func (l *CharacterListenerImpl) OnShowMobHp(oid uint32, percentage uint8) {
	l.character.Send(&resp.ShowMobHp{
		OID:        oid,
		Percentage: percentage,
	}, types.SEND_POLICY_ENCRYPT)
}

// OnUnlockAction handles unlock action events
func (l *CharacterListenerImpl) OnUnlockAction() {
	l.character.Send(&resp.UpdateStats{
		UnlockAction: true,
	}, types.SEND_POLICY_ENCRYPT)
}

// OnItemGainFailed handles item gain failed events
func (l *CharacterListenerImpl) OnItemGainFailed(mode constant.ItemGainFailedType) {
	l.character.Send(&resp.ItemGainFailed{
		Mode: mode,
	}, types.SEND_POLICY_ENCRYPT)
}

// OnInventorySlotUpdated handles inventory slot update events
func (l *CharacterListenerImpl) OnInventorySlotUpdated(inventoryType constant.InventoryType, slot int16, item entity.Item) {
	l.character.Send(&resp.UpdateInventorySlot{
		InventoryType: inventoryType,
		Slot:          slot,
		Item:          item,
	}, types.SEND_POLICY_ENCRYPT)
}

// OnInventorySlotAdded handles inventory slot add events
func (l *CharacterListenerImpl) OnInventorySlotAdded(inventoryType constant.InventoryType, slot int16, item entity.Item) {
	l.character.Send(&resp.AddInventorySlot{
		InventoryType: inventoryType,
		Slot:          slot,
		Item:          item,
	}, types.SEND_POLICY_ENCRYPT)
}

// OnShowItemGain handles item gain display events
func (l *CharacterListenerImpl) OnShowItemGain(itemId uint32, count uint32, mode constant.ShowItemGainType) {
	l.character.Send(&resp.ShowItemGain{
		ItemId: itemId,
		Count:  count,
		Mode:   mode,
	}, types.SEND_POLICY_ENCRYPT)
}

// OnShowMesoGain handles meso gain display events
func (l *CharacterListenerImpl) OnShowMesoGain(count int32, mode constant.ShowMesoGainType) {
	l.character.Send(&resp.ShowMesoGain{
		Count: count,
		Mode:  mode,
	}, types.SEND_POLICY_ENCRYPT)
}

// OnUpdateStats handles general stat update events
func (l *CharacterListenerImpl) OnUpdateStats(stats map[constant.Stat]int32, unlock bool) {
	l.character.Send(&resp.UpdateStats{
		Stats:        stats,
		UnlockAction: unlock,
	}, types.SEND_POLICY_ENCRYPT)
}

// OnMobMoved broadcasts mob movement to all players on the map
// Following mob branch pattern: broadcast to all players except the controller (sender)
func (l *CharacterListenerImpl) OnMobMoved(mapID uint32, mobID uint32, isAggroed bool, centerSplit int8, skill1 uint8, skill2 uint8, skill3 uint8, skill4 uint8, startPoint types.Vector2[int16], movements []action.MoveFragment) {
	// Get the map instance
	mapInstance := l.gameServer.GetMap(mapID)
	if mapInstance == nil {
		return
	}

	// Get mob to find controller
	mob := mapInstance.GetMob(mobID)
	if mob == nil {
		return
	}

	// Get controller to exclude from broadcast (following mob branch pattern)
	controllerTable := mapInstance.GetControllerTable()
	controller, exists := controllerTable.GetController(mob)
	var exceptPlayerID uint32
	if exists {
		exceptPlayerID = controller.GetID()
	}

	// Create move mob packet
	movePacket := &resp.MoveMob{
		IsAggroed:   isAggroed,
		CenterSplit: centerSplit,
		Skill1:      skill1,
		Skill2:      skill2,
		Skill3:      skill3,
		Skill4:      skill4,
		OID:         mobID,
		StartPoint:  startPoint,
		Movements:   movements,
	}

	// Broadcast to all players on the map except the controller (following mob branch pattern)
	// The controller already received ControlMoveMob, so exclude from MoveMob broadcast
	mapInstance.BroadcastToPlayers(movePacket, types.SEND_POLICY_ENCRYPT, exceptPlayerID)
}

// OnPlayerMove sends move packet with fragments to all other players on the map
func (l *CharacterListenerImpl) OnPlayerMove(mapID uint32, playerID uint32, character *entity.Character, startPoint types.Vector2[int16], fragments []action.MoveFragment) {
	// Get the map instance
	mapInstance := l.gameServer.GetMap(mapID)
	if mapInstance == nil {
		return
	}

	// Create move packet
	movePacket := &resp.Move{
		Character:  character,
		Fragments:  fragments,
		StartPoint: startPoint,
	}

	// Broadcast to all players on the map except the moving player
	mapInstance.BroadcastToPlayers(movePacket, types.SEND_POLICY_ENCRYPT, playerID)
}

// OnAttack broadcasts attack to all players on the map
func (l *CharacterListenerImpl) OnAttack(mapID uint32, characterID uint32, attackInfo action.AttackInfo, skillLevel uint8) {
	// Get the map instance
	mapInstance := l.gameServer.GetMap(mapID)
	if mapInstance == nil {
		return
	}

	// Create attack packet
	attackPacket := &resp.Attack{
		AttackInfo: attackInfo,
		SkillLevel: skillLevel,
	}

	// Broadcast to all players on the map
	mapInstance.BroadcastToAllPlayers(attackPacket, types.SEND_POLICY_ENCRYPT)
}

// OnEndSortInventory handles end sort inventory events
func (l *CharacterListenerImpl) OnEndSortInventory(inventoryType constant.InventoryType) {
	l.character.Send(&resp.EndSortInventory{
		InventoryType: inventoryType,
	}, types.SEND_POLICY_ENCRYPT)
}

// OnSwapInventorySlot handles inventory slot swapping events
func (l *CharacterListenerImpl) OnSwapInventorySlot(inventoryType constant.InventoryType, source int16, dest int16, equipmentAction int8) {
	l.character.Send(&resp.SwapInventorySlot{
		InventoryType:   inventoryType,
		Source:          source,
		Dest:            dest,
		EquipmentAction: resp.EquipmentActionType(equipmentAction),
	}, types.SEND_POLICY_ENCRYPT)
}

// OnRemoveInventorySlot handles inventory slot removal events
func (l *CharacterListenerImpl) OnRemoveInventorySlot(inventoryType constant.InventoryType, slot int16) {
	l.character.Send(&resp.RemoveInventorySlot{
		InventoryType: inventoryType,
		Slot:          slot,
	}, types.SEND_POLICY_ENCRYPT)
}

// OnUpdateInventorySlot handles inventory slot update events
func (l *CharacterListenerImpl) OnUpdateInventorySlot(inventoryType constant.InventoryType, slot int16, item entity.Item) {
	l.character.Send(&resp.UpdateInventorySlot{
		InventoryType: inventoryType,
		Slot:          slot,
		Item:          item,
	}, types.SEND_POLICY_ENCRYPT)
}

// OnFullMergeInventorySlot handles full merge inventory slot events
func (l *CharacterListenerImpl) OnFullMergeInventorySlot(inventoryType constant.InventoryType, source int16, dest int16, count uint16) {
	l.character.Send(&resp.FullMergeInventorySlot{
		InventoryType: inventoryType,
		Source:        source,
		Dest:          dest,
		Count:         count,
	}, types.SEND_POLICY_ENCRYPT)
}

// OnPartialMergeInventorySlot handles partial merge inventory slot events
func (l *CharacterListenerImpl) OnPartialMergeInventorySlot(inventoryType constant.InventoryType, source int16, dest int16, sourceCount uint16, destCount uint16) {
	l.character.Send(&resp.PartialMergeInventorySlot{
		InventoryType: inventoryType,
		Source:        source,
		Dest:          dest,
		SourceCount:   sourceCount,
		DestCount:     destCount,
	}, types.SEND_POLICY_ENCRYPT)
}

// OnUpdateCharacterLook handles character look update broadcast events
func (l *CharacterListenerImpl) OnUpdateCharacterLook(character *entity.Character) {
	// Get map instance
	mapInstance := l.gameServer.GetMap(character.GetMap())
	if mapInstance == nil {
		return
	}

	// Create update character look packet
	lookPacket := &resp.UpdateCharacterLook{
		Character: character,
	}

	// Broadcast to all players on the map except the character
	mapInstance.BroadcastToPlayers(lookPacket, types.SEND_POLICY_ENCRYPT, character.GetID())
}

// OnNpcAction handles NPC action events
func (l *CharacterListenerImpl) OnNpcAction(bytes []byte) {
	l.character.Send(&resp.NpcAction{
		Bytes: bytes,
	}, types.SEND_POLICY_ENCRYPT)
}
