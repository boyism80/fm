package server

import (
	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/protocol/request"
	"github.com/boyism80/fm/services/game/client"
)

type PetExceptions struct {
	gs *GameServer
}

func (PetExceptions) New(gs *GameServer) *PetExceptions {
	return &PetExceptions{
		gs: gs,
	}
}

func (h *PetExceptions) Handle(ctx *core.ClientContext, req *request.PetExceptions) error {
	gameClient, ok := ctx.Client.(*client.GameClient)
	if !ok {
		return nil
	}
	character := gameClient.GetCharacter()
	if character == nil || character.Pet == nil {
		return nil
	}

	character.Pet.SetExceptions(req.ItemIDs)

	return nil
}
