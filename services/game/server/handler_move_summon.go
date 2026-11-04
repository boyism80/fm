package server

import (
	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/protocol/request"
)

type MoveSummon struct {
	gs *GameServer
}

func (MoveSummon) New(gs *GameServer) *MoveSummon {
	return &MoveSummon{
		gs: gs,
	}
}

func (h *MoveSummon) Handle(ctx *core.ClientContext, req *request.MoveSummon) error {
	character, err := h.gs.character(ctx)
	if err != nil {
		return nil
	}
	summon := character.Summons().FindByOID(req.OID)
	if summon == nil {
		return nil
	}
	summon.Move(req.Fragments)
	return nil
}
