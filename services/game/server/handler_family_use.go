package server

import (
	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/protocol/request"
	"github.com/boyism80/fm/services/game/client"
	"github.com/boyism80/fm/services/game/domain"
)

type FamilyUse struct {
	gs *GameServer
}

func (FamilyUse) New(gs *GameServer) *FamilyUse {
	return &FamilyUse{gs: gs}
}

func (h *FamilyUse) Handle(ctx *core.ClientContext, req *request.FamilyUse) error {
	gameClient, ok := ctx.Client.(*client.GameClient)
	if !ok {
		return nil
	}
	ch := gameClient.GetCharacter()
	if ch == nil || (req.Name != "" && !domain.ClientText(req.Name).Valid(13)) {
		return nil
	}

	ch.Family.Use(ctx.ActorContext, req.Feature, req.Name)
	return nil
}
