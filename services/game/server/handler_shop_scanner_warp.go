package server

import (
	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/protocol/request"
)

type ShopScannerWarp struct {
	gs *GameServer
}

func (ShopScannerWarp) New(gs *GameServer) *ShopScannerWarp {
	return &ShopScannerWarp{
		gs: gs,
	}
}

func (h *ShopScannerWarp) Handle(ctx *core.ClientContext, req *request.ShopScannerWarp) error {
	ch, err := h.gs.character(ctx)
	if err != nil {
		return nil
	}
	ch.Listener().OnMiniRoomFailed(ch, ch.RemoteShops().VisitFound(ctx.ActorContext, req.SN, req.MapID))
	return nil
}
