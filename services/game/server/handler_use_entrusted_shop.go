package server

import (
	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/protocol/request"
)

type UseEntrustedShop struct {
	gs *GameServer
}

func (UseEntrustedShop) New(gs *GameServer) *UseEntrustedShop {
	return &UseEntrustedShop{
		gs: gs,
	}
}

func (h *UseEntrustedShop) Handle(ctx *core.ClientContext, req *request.UseEntrustedShop) error {
	ch := h.gs.character(ctx)
	if ch == nil {
		return nil
	}
	if ch.UsingCashItem() {
		ch.Listener().OnUpdateStats(ch, nil, true)
		return nil
	}
	ch.StoreBank().CheckShopOpen(ctx.ActorContext)
	return nil
}
