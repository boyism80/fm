package server

import (
	"fmt"
	"log"

	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/protocol/request"
	"github.com/boyism80/fm/protocol/response"
	"github.com/boyism80/fm/services/game/client"
	"github.com/boyism80/fm/services/game/constant"
)

type UseChair struct {
	gs *GameServer
}

func (UseChair) New(gs *GameServer) *UseChair {
	return &UseChair{
		gs: gs,
	}
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

	setupInventory := character.Inventory[constant.InventoryTypeInstallation]
	if setupInventory == nil {
		character.Listener.OnUpdateStats(character, nil, true)
		return nil
	}

	item := setupInventory.FindById(req.ItemID)
	if item == nil {
		log.Printf("Chair item not found: %d", req.ItemID)
		character.Listener.OnUpdateStats(character, nil, true)
		return nil
	}

	character.Chair = req.ItemID

	character.Broadcast(&response.ShowChair{
		CharacterID: character.GetID(),
		ItemID:      req.ItemID,
	}, nil)

	character.Listener.OnUpdateStats(character, nil, true)

	return nil
}
