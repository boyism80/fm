package server

import (
	"fmt"
	"log"

	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/game/client"
	"github.com/boyism80/fm/game/constant"
	"github.com/boyism80/fm/game/entity"
	"github.com/boyism80/fm/protocol/request"
	"github.com/boyism80/fm/protocol/response"
)

type UseChair struct {
	gs     *GameServer
	opcode byte
}

func (UseChair) New(gs *GameServer) *UseChair {
	return &UseChair{
		gs:     gs,
		opcode: 0x1A,
	}
}

func (h *UseChair) GetOpcode() byte {
	return h.opcode
}

func (h *UseChair) Handle(ctx *core.ClientContext, req *request.UseChair) error {
	client, ok := ctx.Client.(*client.GameClient)
	if !ok {
		log.Printf("Client is not a GameClient")
		return fmt.Errorf("client is not a GameClient")
	}

	character := client.GetCharacter()
	if character == nil {
		return nil
	}

	mapInstance := character.GetMap()
	if mapInstance == nil {
		return nil
	}

	setupInventory := character.Inventory[constant.INVENTORY_TYPE_INSTALLATION]
	if setupInventory == nil {
		if character.Listener != nil {
			character.Listener.OnUpdateStats(nil, true)
		}
		return nil
	}

	item := setupInventory.FindById(req.ItemID)
	if item == nil {
		log.Printf("Chair item not found: %d", req.ItemID)
		if character.Listener != nil {
			character.Listener.OnUpdateStats(nil, true)
		}
		return nil
	}

	character.Chair = req.ItemID

	mapInstance.Broadcast(&response.ShowChair{
		CharacterID: character.GetID(),
		ItemID:      req.ItemID,
	}, &entity.BroadcastOption{
		ExceptPlayerIDs: []uint32{character.GetID()},
		Reference:       character,
		RecipientFilter: entity.BroadcastVisibleByReference,
	})

	if character.Listener != nil {
		character.Listener.OnUpdateStats(nil, true)
	}

	return nil
}
