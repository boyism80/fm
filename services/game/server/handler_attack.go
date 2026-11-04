package server

import (
	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/protocol/request"
	"github.com/boyism80/fm/services/game"
)

type Attack struct {
	gs *GameServer
}

func (Attack) New(gs *GameServer) *Attack {
	return &Attack{
		gs: gs,
	}
}

func (h *Attack) Handle(ctx *core.ClientContext, req *request.Attack) error {
	character := h.gs.character(ctx)
	if character == nil {
		return nil
	}

	return character.Attack(ctx.ActorContext, game.AttackClose, req)
}
