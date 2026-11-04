package server

import (
	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/protocol/request"
)

type FamilyRemoveJunior struct {
	gs *GameServer
}

func (FamilyRemoveJunior) New(gs *GameServer) *FamilyRemoveJunior {
	return &FamilyRemoveJunior{gs: gs}
}

func (h *FamilyRemoveJunior) Handle(ctx *core.ClientContext, req *request.FamilyRemoveJunior) error {
	ch, err := h.gs.character(ctx)
	if err != nil {
		return nil
	}

	ch.Family().Separate(ctx.ActorContext, req.JuniorID)
	return nil
}
