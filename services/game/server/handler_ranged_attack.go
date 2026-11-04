package server

import (
	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/protocol/request"
	"github.com/boyism80/fm/services/game"
)

type RangedAttack struct {
	gs *GameServer
}

func (RangedAttack) New(gs *GameServer) *RangedAttack {
	return &RangedAttack{
		gs: gs,
	}
}

func (h *RangedAttack) Handle(ctx *core.ClientContext, req *request.RangedAttack) error {
	character := h.gs.character(ctx)
	if character == nil {
		return nil
	}

	return character.Attack(ctx.ActorContext, game.AttackRanged, req)
}
