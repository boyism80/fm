package server

import (
	"errors"
	"fmt"
	"log"

	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/core/luax"
	"github.com/boyism80/fm/game/client"
	"github.com/boyism80/fm/game/constant"
	"github.com/boyism80/fm/game/entity"
	"github.com/boyism80/fm/protocol/request"
)

type MoveItem struct {
	gs     *GameServer
	opcode byte
}

func (MoveItem) New(gs *GameServer) *MoveItem {
	return &MoveItem{
		gs:     gs,
		opcode: 0x36,
	}
}

func (h *MoveItem) GetOpcode() byte {
	return h.opcode
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
		before := character.Equipments[parts]
		if err := character.UnequipToSlot(parts, req.Dest); err != nil {
			return nil
		}
		callOnEquipmentChanged(ctx, character, parts, before, nil)
	} else if req.Dest < 0 {
		parts := constant.EquipmentPartsType(req.Dest)
		before := character.Equipments[parts]
		inven := character.Inventory[constant.INVENTORY_TYPE_EQUIPMENT]
		var after entity.Equipment
		if item := inven.Items[req.Source]; item != nil {
			after, _ = item.(entity.Equipment)
		}
		if err := character.Equip(req.Source); err != nil {
			if errors.Is(err, entity.ErrInventoryFull) && character.Listener != nil {
				character.Listener.OnItemGainFailed(constant.ITEM_GAIN_FAILED_TYPE_FULL)
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
	item, ok := character.Inventory[invenType].Items[slot]
	if !ok {
		return
	}

	removed := (item.Reduce(count) == 0)
	if removed {
		character.Listener.OnRemoveInventorySlot(invenType, slot)
		delete(character.Inventory[invenType].Items, slot)
	} else {
		character.Listener.OnUpdateInventorySlot(invenType, slot, item)
	}

	spawned := item.Clone(count)
	spawned.BindDrop(&entity.Drop{
		Object: &entity.Object{
			Position: character.Position,
		},
		Owner:        character.GetID(),
		SpawnedPoint: character.Position,
		DropType:     constant.DROP_TYPE_FFA,
	})

	mapInstance := character.GetMap()
	if mapInstance != nil {
		if err := mapInstance.SpawnItem(spawned, character.GetID(), constant.DROP_TYPE_FFA); err != nil {
			log.Printf("Failed to spawn item on map: %v", err)
		}
	}
}

func (h *MoveItem) handleMoveItemInternal(client *client.GameClient, character *entity.Character, invenType constant.InventoryType, sourceSlot int16, destSlot int16) {
	inven := character.Inventory[invenType]
	src, ok := inven.Items[sourceSlot]
	if !ok {
		return
	}

	dst, ok := inven.Items[destSlot]

	if !ok {
		inven.Items[destSlot] = inven.Items[sourceSlot]
		delete(inven.Items, sourceSlot)
		character.Listener.OnSwapInventorySlot(invenType, sourceSlot, destSlot, 0)
		return
	}

	specSrc := src.GetModel()
	specDst := dst.GetModel()
	if specSrc != specDst {
		inven.Items[sourceSlot], inven.Items[destSlot] = inven.Items[destSlot], inven.Items[sourceSlot]
		character.Listener.OnSwapInventorySlot(invenType, sourceSlot, destSlot, 0)
		return
	}

	limit := min(src.GetCount(), specSrc.GetCapacity()-dst.GetCount())
	dst.Increase(limit)
	if src.Reduce(limit) == 0 {
		character.Listener.OnFullMergeInventorySlot(invenType, sourceSlot, destSlot, dst.GetCount())
		delete(inven.Items, sourceSlot)
	} else {
		character.Listener.OnPartialMergeInventorySlot(invenType, sourceSlot, destSlot, src.GetCount(), dst.GetCount())
	}
}

func callOnEquipmentChanged(ctx *core.ClientContext, character *entity.Character, parts constant.EquipmentPartsType, before, after entity.Equipment) {
	if ctx.LogicActorPID == nil {
		return
	}
	root := luax.GetRootLuaState(ctx.LogicActorPID.String())
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
	_, thread, err := luax.Call(root, "script/script.lua", "on_equipment_changed", character, int32(parts), beforeArg, afterArg)
	if thread != nil {
		thread.Close()
	}
	if err != nil {
		log.Printf("Failed to call script on_equipment_changed: %v", err)
	}
}
