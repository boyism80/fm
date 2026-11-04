package server

import (
	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/core/clock"
	"github.com/boyism80/fm/protocol/request"
)

type HealOverTime struct {
	gs *GameServer
}

func (HealOverTime) New(gs *GameServer) *HealOverTime {
	return &HealOverTime{gs: gs}
}

func (h *HealOverTime) Handle(ctx *core.ClientContext, req *request.HealOverTime) error {
	character, err := h.gs.character(ctx)
	if err != nil {
		return nil
	}

	character.HealOverTime(req.HealHP, req.HealMP, req.PRate&1 == 1, clock.Now())
	return nil
}
