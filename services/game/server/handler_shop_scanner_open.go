package server

import (
	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/protocol/request"
)

type ShopScannerOpen struct {
	gs *GameServer
}

func (ShopScannerOpen) New(gs *GameServer) *ShopScannerOpen {
	return &ShopScannerOpen{
		gs: gs,
	}
}

func (h *ShopScannerOpen) Handle(ctx *core.ClientContext, req *request.ShopScannerOpen) error {
	ch := h.gs.character(ctx)
	if ch == nil {
		return nil
	}
	ch.RemoteShops().FindPopularSearches(ctx.ActorContext)
	return nil
}
