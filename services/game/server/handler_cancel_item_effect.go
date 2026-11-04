package server

import (
	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/protocol/request"
)

type CancelItemEffect struct {
	gs *GameServer
}

func (CancelItemEffect) New(gs *GameServer) *CancelItemEffect {
	return &CancelItemEffect{
		gs: gs,
	}
}

func (h *CancelItemEffect) Handle(ctx *core.ClientContext, req *request.CancelItemEffect) error {
	character, err := h.gs.character(ctx)
	if err != nil {
		return nil
	}
	character.Buffs().CancelBySource(req.SourceID)
	return nil
}
