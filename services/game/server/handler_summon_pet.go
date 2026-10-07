package server

import (
	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/protocol/request"
	"github.com/boyism80/fm/services/game/client"
)

type SummonPet struct {
	gs *GameServer
}

func (SummonPet) New(gs *GameServer) *SummonPet {
	return &SummonPet{
		gs: gs,
	}
}

func (h *SummonPet) Handle(ctx *core.ClientContext, req *request.SummonPet) error {
	gameClient, ok := ctx.Client.(*client.GameClient)
	if !ok {
		return nil
	}
	character := gameClient.GetCharacter()
	if character == nil {
		return nil
	}

	err := character.SummonPet(req.Slot)
	if err != nil {
		character.Listener.OnUnlockAction(character)
	}
	return nil
}
