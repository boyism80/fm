package server

import (
	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/protocol/request"
	"github.com/boyism80/fm/services/game"
)

type MagicAttack struct {
	gs *GameServer
}

func (MagicAttack) New(gs *GameServer) *MagicAttack {
	return &MagicAttack{
		gs: gs,
	}
}

func (h *MagicAttack) Handle(ctx *core.ClientContext, req *request.MagicAttack) error {
	character := h.gs.character(ctx)
	if character == nil {
		return nil
	}

	return character.Attack(ctx.ActorContext, game.AttackMagic, req)
}
