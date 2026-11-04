package server

import (
	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/protocol/request"
)

type CancelBuff struct {
	gs *GameServer
}

func (CancelBuff) New(gs *GameServer) *CancelBuff {
	return &CancelBuff{
		gs: gs,
	}
}

func (h *CancelBuff) Handle(ctx *core.ClientContext, req *request.CancelBuff) error {
	character := h.gs.character(ctx)
	if character == nil {
		return nil
	}
	character.Buffs().Cancel(req.SourceID)
	return nil
}
