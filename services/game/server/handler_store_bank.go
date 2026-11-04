package server

import (
	"github.com/boyism80/fm/core"
	pconst "github.com/boyism80/fm/protocol/constant"
	"github.com/boyism80/fm/protocol/request"
)

type StoreBank struct {
	gs *GameServer
}

func (StoreBank) New(gs *GameServer) *StoreBank {
	return &StoreBank{
		gs: gs,
	}
}

func (h *StoreBank) Handle(ctx *core.ClientContext, req *request.StoreBank) error {
	ch, err := h.gs.character(ctx)
	if err != nil {
		return nil
	}

	switch req.Mode {
	case pconst.StoreBankWithdraw:
		err = ch.StoreBank().Withdraw()
	case pconst.StoreBankConfirm:
		err = ch.StoreBank().Confirm(ctx.ActorContext)
	case pconst.StoreBankClose:
		ch.StoreBank().Close()
	}
	if err != nil {
		ch.Listener().OnUnlockAction(ch)
	}
	return nil
}
