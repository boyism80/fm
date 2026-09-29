package server

import (
	"fmt"
	"log"

	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/protocol/request"
	"github.com/boyism80/fm/services/game/client"
)

type SortInventory struct {
	gs *GameServer
}

func (SortInventory) New(gs *GameServer) *SortInventory {
	return &SortInventory{
		gs: gs,
	}
}

func (h *SortInventory) Handle(ctx *core.ClientContext, req *request.SortInventory) error {
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

	character.MergeItems(req.InventoryType)
	character.SortItems(req.InventoryType)

	character.Listener.OnEndSortInventory(character, req.InventoryType)
	character.Listener.OnUpdateStats(character, nil, true)

	return nil
}
