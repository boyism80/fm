package server

import (
	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/protocol/request"
)

type CoconutHit struct {
	gs *GameServer
}

func (CoconutHit) New(gs *GameServer) *CoconutHit {
	return &CoconutHit{gs: gs}
}

// The channel event decides what the struck coconut does.
func (h *CoconutHit) Handle(ctx *core.ClientContext, req *request.CoconutHit) error {
	ch := h.gs.character(ctx)
	if ch == nil {
		return nil
	}

	if ch.StateMachine() == nil {
		return nil
	}
	ch.StateMachine().CallHook("on_coconut_hit", ch, int(req.ID))
	return nil
}
