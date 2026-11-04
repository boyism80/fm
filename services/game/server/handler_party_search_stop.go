package server

import (
	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/protocol/request"
)

type PartySearchStop struct {
	gs *GameServer
}

func (PartySearchStop) New(gs *GameServer) *PartySearchStop {
	return &PartySearchStop{
		gs: gs,
	}
}

func (h *PartySearchStop) Handle(ctx *core.ClientContext, req *request.PartySearchStop) error {
	ch, err := h.gs.character(ctx)
	if err != nil {
		return nil
	}

	ch.Party().SetSearch(nil)
	return nil
}
