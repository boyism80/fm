package server

import (
	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/protocol/request"
)

type UseSummonBag struct {
	gs *GameServer
}

func (UseSummonBag) New(gs *GameServer) *UseSummonBag {
	return &UseSummonBag{gs: gs}
}

func (h *UseSummonBag) Handle(ctx *core.ClientContext, req *request.UseSummonBag) error {
	ch := h.gs.character(ctx)
	if ch == nil {
		return nil
	}
	defer ch.Listener().OnUnlockAction(ch)

	ch.UseSummonBag(req.Slot, req.ItemID)
	return nil
}
