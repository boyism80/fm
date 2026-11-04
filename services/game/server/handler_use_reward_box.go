package server

import (
	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/protocol/request"
)

type UseRewardBox struct {
	gs *GameServer
}

func (UseRewardBox) New(gs *GameServer) *UseRewardBox {
	return &UseRewardBox{gs: gs}
}

func (h *UseRewardBox) Handle(ctx *core.ClientContext, req *request.UseRewardBox) error {
	ch := h.gs.character(ctx)
	if ch == nil {
		return nil
	}
	defer ch.Listener().OnUnlockAction(ch)

	ch.UseRewardBox(req.Slot, req.ItemID)
	return nil
}
