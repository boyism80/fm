package server

import (
	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/protocol/request"
	"github.com/boyism80/fm/services/game"
)

type FamilySummon struct {
	gs *GameServer
}

func (FamilySummon) New(gs *GameServer) *FamilySummon {
	return &FamilySummon{gs: gs}
}

func (h *FamilySummon) Handle(ctx *core.ClientContext, req *request.FamilySummon) error {
	ch := h.gs.character(ctx)
	if ch == nil {
		return nil
	}
	if req.Name != "" && !game.ClientText(req.Name).Valid(13) {
		return nil
	}

	ch.Family().AnswerSummon(req.Name, req.Accepted)
	return nil
}
