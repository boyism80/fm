package server

import (
	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/maple/constant"
	"github.com/boyism80/fm/protocol/request"
)

type UseShopScanner struct {
	gs *GameServer
}

func (UseShopScanner) New(gs *GameServer) *UseShopScanner {
	return &UseShopScanner{
		gs: gs,
	}
}

func (h *UseShopScanner) Handle(ctx *core.ClientContext, req *request.UseShopScanner) error {
	ch := h.gs.character(ctx)
	if ch == nil {
		return nil
	}
	if ch.UsingCashItem() {
		ch.Listener().OnUpdateStats(ch, nil, true)
		return nil
	}
	ch.RemoteShops().Scan(ctx.ActorContext, constant.InventoryTypeConsume, req.Slot, req.ItemID, req.SearchID, req.HighFirst)
	return nil
}
