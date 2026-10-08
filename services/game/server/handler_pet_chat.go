package server

import (
	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/protocol/request"
	"github.com/boyism80/fm/services/game/client"
)

type PetChat struct {
	gs *GameServer
}

func (PetChat) New(gs *GameServer) *PetChat {
	return &PetChat{
		gs: gs,
	}
}

func (h *PetChat) Handle(ctx *core.ClientContext, req *request.PetChat) error {
	gameClient, ok := ctx.Client.(*client.GameClient)
	if !ok {
		return nil
	}
	character := gameClient.GetCharacter()
	if character == nil || character.Pets.Active == nil {
		return nil
	}

	character.Pets.Active.Chat(req.Type, req.Action, req.Text)

	return nil
}
