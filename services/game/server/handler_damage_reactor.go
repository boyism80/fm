package server

import (
	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/protocol/request"
)

type DamageReactor struct {
	gs *GameServer
}

func (DamageReactor) New(gs *GameServer) *DamageReactor {
	return &DamageReactor{
		gs: gs,
	}
}

func (h *DamageReactor) Handle(ctx *core.ClientContext, req *request.DamageReactor) error {
	character, err := h.gs.character(ctx)
	if err != nil {
		return err
	}

	if character.Spectating() {
		return nil
	}

	mapInstance := character.GetMap()
	if mapInstance == nil {
		return nil
	}

	reactor := mapInstance.GetReactor(req.OID)
	if reactor == nil {
		return nil
	}

	reactor.Hit(character, req.HitSide, int32(req.Stance))
	return nil
}
