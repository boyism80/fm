package server

import (
	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/protocol/request"
)

type UseCatchItem struct {
	gs *GameServer
}

func (UseCatchItem) New(gs *GameServer) *UseCatchItem {
	return &UseCatchItem{gs: gs}
}

func (h *UseCatchItem) Handle(ctx *core.ClientContext, req *request.UseCatchItem) error {
	ch := h.gs.character(ctx)
	if ch == nil {
		return nil
	}
	if ch.UsingCashItem() {
		ch.Listener().OnUpdateStats(ch, nil, true)
		return nil
	}

	ch.UseCatchItem(int16(req.Slot), req.ItemID, req.MobOID)
	ch.Listener().OnUpdateStats(ch, nil, true)
	return nil
}
