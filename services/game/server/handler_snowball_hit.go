package server

import (
	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/protocol/request"
)

type SnowballHit struct {
	gs *GameServer
}

func (SnowballHit) New(gs *GameServer) *SnowballHit {
	return &SnowballHit{gs: gs}
}

// The channel event judges the hit by where the character stands, not by the damage the client reports.
func (h *SnowballHit) Handle(ctx *core.ClientContext, req *request.SnowballHit) error {
	ch, err := h.gs.character(ctx)
	if err != nil {
		return err
	}

	if ch.StateMachine() == nil {
		return nil
	}
	ch.StateMachine().CallHook("on_snowball_hit", ch)
	return nil
}
