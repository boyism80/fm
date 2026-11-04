package server

import (
	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/protocol/request"
)

type FamilyOpen struct {
	gs *GameServer
}

func (FamilyOpen) New(gs *GameServer) *FamilyOpen {
	return &FamilyOpen{gs: gs}
}

func (h *FamilyOpen) Handle(ctx *core.ClientContext, req *request.FamilyOpen) error {
	ch, err := h.gs.character(ctx)
	if err != nil {
		return nil
	}

	ch.Family().Open()
	return nil
}
