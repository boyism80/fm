package server

import (
	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/protocol/request"
	"github.com/boyism80/fm/services/game/client"
)

type MovePet struct {
	gs *GameServer
}

func (MovePet) New(gs *GameServer) *MovePet {
	return &MovePet{
		gs: gs,
	}
}

func (h *MovePet) Handle(ctx *core.ClientContext, req *request.MovePet) error {
	gameClient, ok := ctx.Client.(*client.GameClient)
	if !ok {
		return nil
	}
	character := gameClient.GetCharacter()
	if character == nil || character.Pets.Active == nil {
		return nil
	}

	character.Pets.Active.Move(req.Position, req.Fragments)

	return nil
}
