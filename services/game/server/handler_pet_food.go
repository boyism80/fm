package server

import (
	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/protocol/request"
	"github.com/boyism80/fm/services/game/client"
)

type PetFood struct {
	gs *GameServer
}

func (PetFood) New(gs *GameServer) *PetFood {
	return &PetFood{
		gs: gs,
	}
}

func (h *PetFood) Handle(ctx *core.ClientContext, req *request.PetFood) error {
	gameClient, ok := ctx.Client.(*client.GameClient)
	if !ok {
		return nil
	}
	character := gameClient.GetCharacter()
	if character == nil || character.Pets.Active == nil {
		return nil
	}

	err := character.Pets.Active.Feed(req.Slot, req.ItemID)
	if err != nil {
		character.Listener.OnUnlockAction(character)
	}

	return nil
}
