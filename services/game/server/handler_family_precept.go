package server

import (
	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/protocol/request"
	"github.com/boyism80/fm/services/game"
)

type FamilyPrecept struct {
	gs *GameServer
}

func (FamilyPrecept) New(gs *GameServer) *FamilyPrecept {
	return &FamilyPrecept{gs: gs}
}

func (h *FamilyPrecept) Handle(ctx *core.ClientContext, req *request.FamilyPrecept) error {
	ch, err := h.gs.character(ctx)
	if err != nil {
		return nil
	}
	if req.Precept != "" && !game.ClientText(req.Precept).Valid(200) {
		return nil
	}

	ch.Family().SetPrecept(ctx.ActorContext, req.Precept)
	return nil
}
