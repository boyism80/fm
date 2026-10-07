package server

import (
	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/protocol/request"
	"github.com/boyism80/fm/services/game/client"
)

type PetReviveInquiry struct {
	gs *GameServer
}

func (PetReviveInquiry) New(gs *GameServer) *PetReviveInquiry {
	return &PetReviveInquiry{
		gs: gs,
	}
}

func (h *PetReviveInquiry) Handle(ctx *core.ClientContext, req *request.PetReviveInquiry) error {
	gameClient, ok := ctx.Client.(*client.GameClient)
	if !ok {
		return nil
	}
	character := gameClient.GetCharacter()
	if character == nil {
		return nil
	}

	character.Listener.OnUnlockAction(character)
	return nil
}
