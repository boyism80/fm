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

	mapInstance := character.GetMap()
	if mapInstance == nil {
		log.Printf("Map not found for character %d", character.GetID())
		character.Listener.OnUpdateStats(character, nil, true)
		return fmt.Errorf("map not found")
	}

	lootedObject, reason := mapInstance.LootItem(req.OID, character, req.Position)
	if reason != constant.LOOT_SUCCESS {
		log.Printf("Failed to loot item %d for character %d, reason: %d", req.OID, character.GetID(), reason)

		if reason == constant.LOOT_FAILED_INVENTORY_FULL || reason == constant.LOOT_FAILED_MESO_FULL {
			character.Listener.OnItemGainFailed(character, constant.ITEM_GAIN_FAILED_TYPE_FULL)
		}

		character.Listener.OnUpdateStats(character, nil, true)
		return nil
	}

	switch obj := lootedObject.(type) {
	case entity.Item:
		item := obj
		_, err := character.AddItem(item, false)
		if err != nil {
			log.Printf("Failed to add item: %v", err)
		}

	case *entity.Meso:
		mesoCount := obj.GetCount32()
		character.GainMeso(int32(mesoCount))

	default:
		log.Printf("Unknown looted object type for OID %d", req.OID)
		character.Listener.OnUpdateStats(character, nil, true)
		return nil
	}

	character.Listener.OnUpdateStats(character, nil, true)

	return nil
}
