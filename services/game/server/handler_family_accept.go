package server

import (
	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/protocol/request"
	"github.com/boyism80/fm/services/game"
)

type FamilyAccept struct {
	gs *GameServer
}

func (FamilyAccept) New(gs *GameServer) *FamilyAccept {
	return &FamilyAccept{gs: gs}
}

func (h *FamilyAccept) Handle(ctx *core.ClientContext, req *request.FamilyAccept) error {
	ch := h.gs.character(ctx)
	if ch == nil {
		return nil
	}
	if req.Name != "" && !game.ClientText(req.Name).Valid(13) {
		return nil
	}

	ch.Family().AnswerInvite(ctx.ActorContext, req.SeniorID, req.Name, req.Accepted)
	return nil
}
