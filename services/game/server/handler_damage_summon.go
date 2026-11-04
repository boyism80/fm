package server

import (
	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/protocol/request"
)

type DamageSummon struct {
	gs *GameServer
}

func (DamageSummon) New(gs *GameServer) *DamageSummon {
	return &DamageSummon{
		gs: gs,
	}
}

func (h *DamageSummon) Handle(ctx *core.ClientContext, req *request.DamageSummon) error {
	character := h.gs.character(ctx)
	if character == nil {
		return nil
	}
	summon := character.Summons().FindByOID(req.OID)
	if summon == nil {
		return nil
	}
	summon.TakeDamage(req.Unknown, req.Damage, req.MonsterIdFrom)
	return nil
}
