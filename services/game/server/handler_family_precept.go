package server

import (
	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/protocol/request"
	"github.com/boyism80/fm/services/game"
	"github.com/boyism80/fm/services/game/client"
)

type FamilyPrecept struct {
	gs *GameServer
}

func (FamilyPrecept) New(gs *GameServer) *FamilyPrecept {
	return &FamilyPrecept{gs: gs}
}

func (h *FamilyPrecept) Handle(ctx *core.ClientContext, req *request.FamilyPrecept) error {
	gameClient, ok := ctx.Client.(*client.GameClient)
	if !ok {
		return nil
	}
	ch := gameClient.GetCharacter()
	if ch == nil || (req.Precept != "" && !game.ClientText(req.Precept).Valid(200)) {
		return nil
	}

	ch.Family().SetPrecept(ctx.ActorContext, req.Precept)
	return nil
}
