package server

import (
	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/protocol/request"
	"github.com/boyism80/fm/services/game/client"
)

type InspectCharacter struct {
	gs *GameServer
}

func (InspectCharacter) New(gs *GameServer) *InspectCharacter {
	return &InspectCharacter{
		gs: gs,
	}
}

func (h *InspectCharacter) Handle(ctx *core.ClientContext, req *request.InspectCharacter) error {
	gameClient, ok := ctx.Client.(*client.GameClient)
	if !ok {
		return nil
	}
	character := gameClient.GetCharacter()
	if character == nil {
		return nil
	}

	character.Inspect(req.CharacterID)
	return nil
}
