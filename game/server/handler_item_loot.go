package server

import (
	"fmt"
	"log"

	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/game/client"
	"github.com/boyism80/fm/game/constant"
	"github.com/boyism80/fm/game/entity"
	"github.com/boyism80/fm/protocol/request"
)

// ItemLoot handles item looting packet requests
type ItemLoot struct {
	gs     *GameServer
	opcode byte
}

func (ItemLoot) New(gs *GameServer) *ItemLoot {
	return &ItemLoot{
		gs:     gs,
		opcode: 0xA3,
	}
}

func (h *ItemLoot) GetOpcode() byte {
	return h.opcode
}

func (h *ItemLoot) Handle(ctx *core.ClientContext, req *request.ItemLoot) error {
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

	mapInstance := h.gs.GetMap(character.GetMap())
	if mapInstance == nil {
		log.Printf("Map %d not found for character %d", character.GetMap(), character.GetID())
		character.Listener.OnUpdateStats(nil, true)
		return fmt.Errorf("map not found")
	}

	lootedObject, reason := mapInstance.LootItem(req.OID, character, req.Position)
	if reason != constant.LOOT_SUCCESS {
		log.Printf("Failed to loot item %d for character %d, reason: %d", req.OID, character.GetID(), reason)

		if reason == constant.LOOT_FAILED_INVENTORY_FULL || reason == constant.LOOT_FAILED_MESO_FULL {
			character.Listener.OnItemGainFailed(constant.ITEM_GAIN_FAILED_TYPE_FULL)
		}

		character.Listener.OnUpdateStats(nil, true)
		return nil
	}

	switch obj := lootedObject.(type) {
	case entity.Item:
		item := obj
		invenType := item.GetInventoryType()
		inven := character.Inventory[invenType]
		model := item.GetModel()
		gain := uint16(0)

		for item.GetCount() > 0 {
			slot, ok := inven.FindSlot(model)
			if !ok {
				break
			}

			exists, ok := inven.Items[int16(slot)]
			cap := uint16(0)
			if ok {
				cap = min(model.GetCapacity()-exists.GetCount(), item.GetCount())
				exists.Increase(cap)
				character.Listener.OnInventorySlotUpdated(invenType, int16(slot), exists)
			} else {
				cap = min(model.GetCapacity(), item.GetCount())
				inven.Items[int16(slot)] = item.Clone(cap)
				character.Listener.OnInventorySlotAdded(invenType, int16(slot), inven.Items[int16(slot)])
			}
			if item.Reduce(cap) == 0 {
				break
			}
			gain += cap
		}

		character.Listener.OnShowItemGain(model.GetID(), uint32(gain), constant.ShowItemGainTypeStatus)

	case *entity.Meso:
		meso := obj
		mesoCount := meso.GetCount32()

		character.Meso += int32(mesoCount)

		character.Listener.OnShowMesoGain(int32(mesoCount), constant.ShowMesoGainTypeStatus)

		character.Listener.OnUpdateStats(map[constant.Stat]int32{
			constant.STAT_MESO: character.Meso,
		}, false)

	default:
		log.Printf("Unknown looted object type for OID %d", req.OID)
		character.Listener.OnUpdateStats(nil, true)
		return nil
	}

	character.Listener.OnUpdateStats(nil, true)

	return nil
}
