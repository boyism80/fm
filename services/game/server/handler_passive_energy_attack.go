package server

import (
	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/protocol/request"
	"github.com/boyism80/fm/services/game"
)

type PassiveEnergyAttack struct {
	gs *GameServer
}

func (PassiveEnergyAttack) New(gs *GameServer) *PassiveEnergyAttack {
	return &PassiveEnergyAttack{gs: gs}
}

func (h *PassiveEnergyAttack) Handle(ctx *core.ClientContext, req *request.PassiveEnergyAttack) error {
	character := h.gs.character(ctx)
	if character == nil {
		return nil
	}

	return character.Attack(ctx.ActorContext, game.AttackEnergy, req)
}
