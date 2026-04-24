package server

import (
	"fmt"
	"log"

	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/protocol/request"
	"github.com/boyism80/fm/services/game/client"
	"github.com/boyism80/fm/services/game/constant"
	"github.com/boyism80/fm/services/game/entity"
)

type ItemLoot struct {
	gs *GameServer
}

func (ItemLoot) New(gs *GameServer) *ItemLoot {
	return &ItemLoot{
		gs: gs,
	}
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

	mapInstance := character.GetMap()
	if mapInstance == nil {
		log.Printf("Map not found for character %d", character.GetID())
		character.Listener.OnUpdateStats(character, nil, true)
		return fmt.Errorf("map not found")
	}

	obj := mapInstance.GetItems()[req.OID]
	if obj == nil {
		log.Printf("Loot object not found for OID %d", req.OID)
		character.Listener.OnUpdateStats(character, nil, true)
		return fmt.Errorf("item not found")
	}

	if item, ok := obj.(entity.Item); ok && item.GetModel().IsConsumeOnPickup() {
		consume, ok := item.(*entity.Consume)
		if !ok {
			log.Printf("Item is not a Consume for OID %d", req.OID)
			character.Listener.OnUpdateStats(character, nil, true)
			return fmt.Errorf("item is not a Consume for OID %d", req.OID)
		}

		if !character.ApplyConsumeEffect(consume) {
			log.Printf("Failed to apply consume effect for OID %d", req.OID)
			character.Listener.OnUpdateStats(character, nil, true)
			return fmt.Errorf("failed to apply consume effect for OID %d", req.OID)
		}
	} else {
		reason := mapInstance.LootItem(obj, character, req.Position)
		if reason != constant.LOOT_SUCCESS {
			log.Printf("Failed to loot item %d for character %d, reason: %d", req.OID, character.GetID(), reason)
			if reason == constant.LOOT_FAILED_INVENTORY_FULL || reason == constant.LOOT_FAILED_MESO_FULL {
				character.Listener.OnItemGainFailed(character, constant.ITEM_GAIN_FAILED_TYPE_FULL)
			}
			character.Listener.OnUpdateStats(character, nil, true)
			return fmt.Errorf("failed to loot item for OID %d", req.OID)
		}
	}

	if err := mapInstance.RemoveItem(req.OID, constant.REMOVE_ITEM_TYPE_ANIMATED, character.GetID()); err != nil {
		log.Printf("Failed to remove item %d for character %d, reason: %v", req.OID, character.GetID(), err)
		character.Listener.OnUpdateStats(character, nil, true)
		return fmt.Errorf("failed to remove item")
	}

	character.Listener.OnUpdateStats(character, nil, true)
	return nil
}
