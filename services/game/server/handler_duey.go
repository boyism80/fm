package server

import (
	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/maple/constant"
	pconst "github.com/boyism80/fm/protocol/constant"
	"github.com/boyism80/fm/protocol/request"
)

type Duey struct {
	gs *GameServer
}

func (Duey) New(gs *GameServer) *Duey {
	return &Duey{
		gs: gs,
	}
}

func (h *Duey) Handle(ctx *core.ClientContext, req *request.Duey) error {
	ch := h.gs.character(ctx)
	if ch == nil {
		return nil
	}
	if ch.UsingCashItem() {
		ch.Listener().OnUpdateStats(ch, nil, true)
		return nil
	}

	var err error
	switch req.Mode {
	case pconst.DueyOpenFromArrival:
		ch.Duey().Open(ctx.ActorContext, true)
	case pconst.DueyIdentity:
		err = ch.Duey().ConfirmIdentity()
	case pconst.DueySend:
		err = ch.Duey().Send(ctx.ActorContext, constant.InventoryType(req.InventoryType), req.Slot, req.Count, req.Meso, req.Recipient, req.Quick, req.Message)
	case pconst.DueyReceive:
		err = ch.Duey().Receive(ctx.ActorContext, req.ParcelID)
	case pconst.DueyDelete:
		err = ch.Duey().Delete(ctx.ActorContext, req.ParcelID)
	case pconst.DueyClose:
		ch.Duey().Close()
	}

	ch.Duey().Report(err)
	return nil
}
