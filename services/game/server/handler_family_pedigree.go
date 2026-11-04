package server

import (
	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/protocol/request"
	"github.com/boyism80/fm/services/game"
)

type FamilyPedigree struct {
	gs *GameServer
}

func (FamilyPedigree) New(gs *GameServer) *FamilyPedigree {
	return &FamilyPedigree{gs: gs}
}

func (h *FamilyPedigree) Handle(ctx *core.ClientContext, req *request.FamilyPedigree) error {
	ch, err := h.gs.character(ctx)
	if err != nil {
		return nil
	}
	if req.Name != "" && !game.ClientText(req.Name).Valid(13) {
		return nil
	}

	ch.Family().ShowPedigree(ctx.ActorContext, req.Name)
	return nil
}
