package server

import (
	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/protocol/request"
)

type UseReturnScroll struct {
	gs *GameServer
}

func (UseReturnScroll) New(gs *GameServer) *UseReturnScroll {
	return &UseReturnScroll{gs: gs}
}

func (h *UseReturnScroll) Handle(ctx *core.ClientContext, req *request.UseReturnScroll) error {
	ch := h.gs.character(ctx)
	if ch == nil {
		return nil
	}
	if ch.UsingCashItem() {
		ch.Listener().OnUpdateStats(ch, nil, true)
		return nil
	}

	return ch.UseReturnScroll(ctx.ActorContext, int16(req.Slot), req.ItemID)
}
