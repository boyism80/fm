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

	if req.Source < 0 {
		parts := constant.EquipmentPartsType(req.Source)
		before := character.Inventory.Equipped[parts]
		if err := character.Inventory.UnequipToSlot(parts, req.Dest); err != nil {
			return nil
		}
		callOnEquipmentChanged(ctx, character, parts, before, nil)
	} else if req.Dest < 0 {
		parts := constant.EquipmentPartsType(req.Dest)
		before := character.Inventory.Equipped[parts]
		inven := character.Inventory.Containers[constant.InventoryTypeEquipment]
		var after entity.Equipment
		if item := inven.Items[req.Source]; item != nil {
			after, _ = item.(entity.Equipment)
		}
		if err := character.Inventory.Equip(req.Source); err != nil {
			if errors.Is(err, entity.ErrInventoryFull) {
				character.Listener.OnItemGainFailed(character, constant.ItemGainFailedTypeFull)
			}
			return nil
		}
		callOnEquipmentChanged(ctx, character, parts, before, after)
	} else if req.Dest == 0 {
		h.handleDrop(client, character, req.InventoryType, req.Source, req.Count)
	} else {
		h.handleMoveItemInternal(client, character, req.InventoryType, req.Source, req.Dest)
	}

	return nil
}

func (h *MoveItem) handleDrop(client *client.GameClient, character *entity.Character, invenType constant.InventoryType, slot int16, count uint16) {
	item, ok := character.Inventory.Containers[invenType].Items[slot]
	if !ok {
		return
	}

	removed := (item.Reduce(count) == 0)
	if removed {
		character.Listener.OnRemoveInventorySlot(character, invenType, slot)
		delete(character.Inventory.Containers[invenType].Items, slot)
	} else {
		character.Listener.OnUpdateInventorySlot(character, invenType, slot, item)
	}

	spawned := item.Clone(count)
	spawned.BindFieldPlacement(&entity.FieldPlacement{
		ObjectCore: &entity.ObjectCore{
			Position: character.Position,
		},
		Owner:        character.GetID(),
		SpawnedPoint: character.Position,
		DropType:     constant.DropTypeFFA,
		PlayerDrop:   true,
	})

	mapInstance := character.GetMap()
	if mapInstance != nil {
		if err := mapInstance.SpawnItem(spawned, character.GetID(), constant.DropTypeFFA); err != nil {
			log.Printf("Failed to spawn item on map: %v", err)
		}
	}
}

func (h *MoveItem) handleMoveItemInternal(client *client.GameClient, character *entity.Character, invenType constant.InventoryType, sourceSlot int16, destSlot int16) {
	inven := character.Inventory.Containers[invenType]
	src, ok := inven.Items[sourceSlot]
	if !ok {
		return
	}

	dst, ok := inven.Items[destSlot]

	if !ok {
		inven.Items[destSlot] = inven.Items[sourceSlot]
		delete(inven.Items, sourceSlot)
		character.Listener.OnSwapInventorySlot(character, invenType, sourceSlot, destSlot, 0)
		return
	}

	specSrc := src.GetModel()
	specDst := dst.GetModel()
	if specSrc != specDst {
		inven.Items[sourceSlot], inven.Items[destSlot] = inven.Items[destSlot], inven.Items[sourceSlot]
		character.Listener.OnSwapInventorySlot(character, invenType, sourceSlot, destSlot, 0)
		return
	}

	limit := min(src.GetCount(), specSrc.GetCapacity()-dst.GetCount())
	dst.Increase(limit)
	if src.Reduce(limit) == 0 {
		character.Listener.OnFullMergeInventorySlot(character, invenType, sourceSlot, destSlot, dst.GetCount())
		delete(inven.Items, sourceSlot)
	} else {
		character.Listener.OnPartialMergeInventorySlot(character, invenType, sourceSlot, destSlot, src.GetCount(), dst.GetCount())
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
	luax.CallAsync(root, thread, "on_equipment_changed", character, int32(parts), beforeArg, afterArg).OnError(func(err error) {
		log.Printf("Failed to call script on_equipment_changed: %v", err)
	})
}
