package server

import (
	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/protocol/request"
)

type EnhanceEquipment struct {
	gs *GameServer
}

func (EnhanceEquipment) New(gs *GameServer) *EnhanceEquipment {
	return &EnhanceEquipment{gs: gs}
}

func (h *EnhanceEquipment) Handle(ctx *core.ClientContext, req *request.EnhanceEquipment) error {
	ch, err := h.gs.character(ctx)
	if err != nil {
		return nil
	}
	if ch.UsingCashItem() {
		ch.Listener().OnUpdateStats(ch, nil, true)
		return nil
	}

	if err := ch.Scroll(ctx.ActorContext, req.ScrollSlot, req.TargetSlot); err != nil {
		ch.Listener().OnUpdateStats(ch, nil, true)
	}
	return nil
}
