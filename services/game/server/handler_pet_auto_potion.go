package server

import (
	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/protocol/request"
	"github.com/boyism80/fm/services/game/client"
)

type PetAutoPotion struct {
	gs *GameServer
}

func (PetAutoPotion) New(gs *GameServer) *PetAutoPotion {
	return &PetAutoPotion{
		gs: gs,
	}
}

func (h *PetAutoPotion) Handle(ctx *core.ClientContext, req *request.PetAutoPotion) error {
	gameClient, ok := ctx.Client.(*client.GameClient)
	if !ok {
		return nil
	}
	character := gameClient.GetCharacter()
	if character == nil || character.Pet == nil {
		return nil
	}

	err := character.Pet.UsePotion(req.Slot, req.ItemID)
	if err != nil {
		character.Listener.OnUnlockAction(character)
	}

	return nil
}
