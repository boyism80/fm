package server

import (
	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/protocol/request"
)

type FamilyRemoveSenior struct {
	gs *GameServer
}

func (FamilyRemoveSenior) New(gs *GameServer) *FamilyRemoveSenior {
	return &FamilyRemoveSenior{gs: gs}
}

func (h *FamilyRemoveSenior) Handle(ctx *core.ClientContext, req *request.FamilyRemoveSenior) error {
	ch := h.gs.character(ctx)
	if ch == nil {
		return nil
	}

	ch.Family().Separate(ctx.ActorContext, ch.GetID())
	return nil
}
