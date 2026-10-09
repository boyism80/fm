package server

import (
	"errors"
	"fmt"
	"log"

	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/core/luax"
	"github.com/boyism80/fm/protocol/request"
	"github.com/boyism80/fm/services/game/client"
	"github.com/boyism80/fm/services/game/constant"
	"github.com/boyism80/fm/services/game/entity"
)

type MoveItem struct {
	gs *GameServer
}

func (MoveItem) New(gs *GameServer) *MoveItem {
	return &MoveItem{
		gs: gs,
	}
}

func (h *MoveItem) Handle(ctx *core.ClientContext, req *request.MoveItem) error {
	client, ok := ctx.Client.(*client.GameClient)
	if !ok {
		log.Printf("Client is not a GameClient")
		return fmt.Errorf("client is not a GameClient")
	}

	character := client.GetCharacter()
	if character == nil {
		log.Printf("Character is nil for client")
		return fmt.Errorf("character is nil")
	}
	if character.Trading() {
		character.Listener.OnUpdateStats(character, nil, true)
		return nil
	}

	switch {
	case req.Dest == 0:
		h.drop(ctx, character, req.InventoryType, req.Source, req.Count)
	case req.Source < 0:
		parts := constant.EquipmentPartsType(req.Source)
		before := character.Inventory.Equipped[parts]
		if err := character.Inventory.UnequipToSlot(parts, req.Dest); err != nil {
			character.Listener.OnUpdateStats(character, nil, true)
			return nil
		}
		callOnEquipmentChanged(ctx, character, parts, before, nil)
	case req.Dest < 0:
		parts := constant.EquipmentPartsType(req.Dest)
		before := character.Inventory.Equipped[parts]
		inven := character.Inventory.Tabs[constant.InventoryTypeEquipment]
		var after entity.Equipment
		if item := inven.Items[req.Source]; item != nil {
			after, _ = item.(entity.Equipment)
		}
		if err := character.Inventory.Equip(req.Source, parts); err != nil {
			if errors.Is(err, entity.ErrInventoryFull) {
				character.Listener.OnItemGainFailed(character, constant.ItemGainFailedTypeFull)
			}
			character.Listener.OnUpdateStats(character, nil, true)
			return nil
		}
		callOnEquipmentChanged(ctx, character, parts, before, after)
	default:
		h.move(character, req.InventoryType, req.Source, req.Dest)
	}

	return nil
}

func (h *MoveItem) drop(ctx *core.ClientContext, ch *entity.Character, invenType constant.InventoryType, slot int16, count uint16) {
	if ch.Spectating() {
		ch.Listener.OnUpdateStats(ch, nil, true)
		return
	}

	var spawned entity.Item
	if slot < 0 {
		parts := constant.EquipmentPartsType(slot)
		equipment := ch.Inventory.RemoveEquipped(parts)
		if equipment == nil {
			ch.Listener.OnUpdateStats(ch, nil, true)
			return
		}
		callOnEquipmentChanged(ctx, ch, parts, equipment, nil)
		spawned = equipment
	} else {
		inven := ch.Inventory.Tabs[invenType]
		if inven == nil {
			return
		}
		item, ok := inven.Items[slot]
		if !ok {
			return
		}
		if _, isPet := item.(*entity.Pet); isPet {
			ch.Listener.OnUpdateStats(ch, nil, true)
			return
		}

		actualCount := min(count, item.GetCount())
		if actualCount == 0 {
			return
		}

		spawned = item.Clone(actualCount)
		if ch.Inventory.RemoveItem(invenType, slot, actualCount) == false {
			return
		}
	}

	mapInstance := ch.GetMap()
	if mapInstance == nil {
		return
	}
	if err := mapInstance.SpawnItem(spawned, entity.ItemSpawn{
		Position:   ch.Position,
		From:       ch.Position,
		Owner:      ch.GetID(),
		DropType:   constant.DropTypeFFA,
		PlayerDrop: true,
	}); err != nil {
		log.Printf("Failed to spawn item on map: %v", err)
	}
}

func (h *MoveItem) move(ch *entity.Character, invenType constant.InventoryType, sourceSlot int16, destSlot int16) {
	inven := ch.Inventory.Tabs[invenType]
	if inven == nil {
		return
	}
	if destSlot < 1 || destSlot > int16(inven.SlotLimit) {
		ch.Listener.OnUpdateStats(ch, nil, true)
		return
	}
	src, ok := inven.Items[sourceSlot]
	if !ok {
		return
	}

	dst, ok := inven.Items[destSlot]
	if !ok {
		inven.Items[destSlot] = inven.Items[sourceSlot]
		delete(inven.Items, sourceSlot)
		ch.Listener.OnSwapInventorySlot(ch, invenType, sourceSlot, destSlot, 0)
		return
	}

	specSrc := src.GetModel()
	specDst := dst.GetModel()
	if specSrc != specDst {
		inven.Items[sourceSlot], inven.Items[destSlot] = inven.Items[destSlot], inven.Items[sourceSlot]
		ch.Listener.OnSwapInventorySlot(ch, invenType, sourceSlot, destSlot, 0)
		return
	}

	limit := min(src.GetCount(), specSrc.GetCapacity()-dst.GetCount())
	dst.Increase(limit)
	if src.Reduce(limit) == 0 {
		ch.Listener.OnFullMergeInventorySlot(ch, invenType, sourceSlot, destSlot, dst.GetCount())
		delete(inven.Items, sourceSlot)
	} else {
		ch.Listener.OnPartialMergeInventorySlot(ch, invenType, sourceSlot, destSlot, src.GetCount(), dst.GetCount())
	}
}

func callOnEquipmentChanged(ctx *core.ClientContext, character *entity.Character, parts constant.EquipmentPartsType, before, after entity.Equipment) {
	mapInstance := character.GetMap()
	if mapInstance == nil {
		return
	}
	root := mapInstance.GetLuaRoot()
	if root == nil {
		return
	}
	var beforeArg, afterArg interface{}
	if before != nil {
		beforeArg = luax.NewLuable(root, before.(luax.Luable))
	}
	if after != nil {
		afterArg = luax.NewLuable(root, after.(luax.Luable))
	}
	thread, err := luax.NewThread(root, constant.CharacterHookScriptPath)
	if err != nil {
		log.Printf("Failed to call script on_equipment_changed: %v", err)
		return
	}
	luax.CallAsync(ctx.ActorContext, root, thread, "on_equipment_changed", character, int32(parts), beforeArg, afterArg).OnError(func(err error) {
		log.Printf("Failed to call script on_equipment_changed: %v", err)
	})
}
