package server

import (
	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/protocol/request"
	"github.com/boyism80/fm/services/game"
)

type FamilyUse struct {
	gs *GameServer
}

func (FamilyUse) New(gs *GameServer) *FamilyUse {
	return &FamilyUse{gs: gs}
}

func (h *FamilyUse) Handle(ctx *core.ClientContext, req *request.FamilyUse) error {
	ch := h.gs.character(ctx)
	if ch == nil {
		return nil
	}
	if req.Name != "" && !game.ClientText(req.Name).Valid(13) {
		return nil
	}

	ch.Family().Use(ctx.ActorContext, req.Feature, req.Name)
	return nil
}
