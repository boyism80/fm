package server

import (
	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/protocol/request"
	"github.com/boyism80/fm/services/game/client"
)

type PetCommand struct {
	gs *GameServer
}

func (PetCommand) New(gs *GameServer) *PetCommand {
	return &PetCommand{
		gs: gs,
	}
}

func (h *PetCommand) Handle(ctx *core.ClientContext, req *request.PetCommand) error {
	gameClient, ok := ctx.Client.(*client.GameClient)
	if !ok {
		return nil
	}
	character := gameClient.GetCharacter()
	if character == nil || character.Pets.Active == nil {
		return nil
	}

	character.Pets.Active.Command(req.Index, req.CalledByName)

	return nil
}
