package server

import (
	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/protocol/request"
	"github.com/boyism80/fm/services/game/client"
)

type FamilyRemoveJunior struct {
	gs *GameServer
}

func (FamilyRemoveJunior) New(gs *GameServer) *FamilyRemoveJunior {
	return &FamilyRemoveJunior{gs: gs}
}

func (h *FamilyRemoveJunior) Handle(ctx *core.ClientContext, req *request.FamilyRemoveJunior) error {
	gameClient, ok := ctx.Client.(*client.GameClient)
	if !ok {
		return nil
	}
	ch := gameClient.GetCharacter()
	if ch == nil {
		return nil
	}

	ch.Family().Separate(ctx.ActorContext, req.JuniorID)
	return nil
}
