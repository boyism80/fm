package server

import (
	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/protocol/request"
)

type CloseChalkboard struct {
	gs *GameServer
}

func (CloseChalkboard) New(gs *GameServer) *CloseChalkboard { return &CloseChalkboard{gs: gs} }
func (h *CloseChalkboard) Handle(ctx *core.ClientContext, _ *request.CloseChalkboard) error {
	ch := h.gs.character(ctx)
	if ch == nil {
		return nil
	}
	ch.CloseChalkboard()
	return nil
}
