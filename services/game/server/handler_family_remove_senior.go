package server

import (
	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/protocol/request"
	"github.com/boyism80/fm/services/game/client"
)

type FamilyRemoveSenior struct {
	gs *GameServer
}

func (FamilyRemoveSenior) New(gs *GameServer) *FamilyRemoveSenior {
	return &FamilyRemoveSenior{gs: gs}
}

func (h *FamilyRemoveSenior) Handle(ctx *core.ClientContext, req *request.FamilyRemoveSenior) error {
	gameClient, ok := ctx.Client.(*client.GameClient)
	if !ok {
		return nil
	}
	ch := gameClient.GetCharacter()
	if ch == nil {
		return nil
	}

	ch.Family.Separate(ctx.ActorContext, ch.GetID())
	return nil
}
