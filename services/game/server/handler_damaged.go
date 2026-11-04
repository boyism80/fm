package server

import (
	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/protocol/request"
	"github.com/boyism80/fm/services/game"
)

type Damaged struct {
	gs *GameServer
}

func (Damaged) New(gs *GameServer) *Damaged {
	return &Damaged{
		gs: gs,
	}
}

func (h *Damaged) Handle(ctx *core.ClientContext, req *request.Damaged) error {
	character := h.gs.character(ctx)
	if character == nil {
		return nil
	}

	character.TakeHit(ctx.ActorContext, game.Hit{
		Type:       req.Type,
		Damage:     req.Damage,
		MobOID:     req.OID,
		SkillID:    uint32(req.SkillID),
		SkillLevel: req.Level,
		Reflect:    req.Reflect,
	})
	return nil
}
