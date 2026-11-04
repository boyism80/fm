package server

import (
	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/protocol/request"
)

type PulleyPull struct {
	gs *GameServer
}

func (PulleyPull) New(gs *GameServer) *PulleyPull {
	return &PulleyPull{gs: gs}
}

func (h *PulleyPull) Handle(ctx *core.ClientContext, req *request.PulleyPull) error {
	ch, err := h.gs.character(ctx)
	if err != nil {
		return err
	}

	m := ch.GetMap()
	if m == nil {
		return nil
	}
	m.PullPulley(ch)
	return nil
}
