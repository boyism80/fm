package server

import (
	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/protocol/request"
)

type RemoteEntrustedShop struct {
	gs *GameServer
}

func (RemoteEntrustedShop) New(gs *GameServer) *RemoteEntrustedShop {
	return &RemoteEntrustedShop{
		gs: gs,
	}
}

func (h *RemoteEntrustedShop) Handle(ctx *core.ClientContext, req *request.RemoteEntrustedShop) error {
	ch := h.gs.character(ctx)
	if ch == nil {
		return nil
	}
	if ch.UsingCashItem() {
		ch.Listener().OnUpdateStats(ch, nil, true)
		return nil
	}
	ch.Listener().OnMiniRoomFailed(ch, ch.RemoteShops().Manage(ctx.ActorContext, req.Slot))
	return nil
}
