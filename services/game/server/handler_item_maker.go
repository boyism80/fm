package server

import (
	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/protocol/request"
)

type ItemMaker struct {
	gs *GameServer
}

func (ItemMaker) New(gs *GameServer) *ItemMaker {
	return &ItemMaker{gs: gs}
}

func (h *ItemMaker) Handle(ctx *core.ClientContext, req *request.ItemMaker) error {
	ch := h.gs.character(ctx)
	if ch == nil {
		return nil
	}
	defer ch.Listener().OnUnlockAction(ch)
	if req.Malformed {
		return nil
	}

	ch.Inventory().Make(req.Type, req.ItemID, req.Catalyst, req.Enhancers, int16(req.Slot))
	return nil
}
