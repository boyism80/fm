package server

import (
	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/protocol/request"
)

type Aggro struct {
	gs *GameServer
}

func (Aggro) New(gs *GameServer) *Aggro {
	return &Aggro{
		gs: gs,
	}
}

func (h *Aggro) Handle(ctx *core.ClientContext, req *request.Aggro) error {
	ch, err := h.gs.character(ctx)
	if err != nil {
		return err
	}

	m := ch.GetMap()
	if m == nil {
		return nil
	}
	mob := m.GetMob(req.OID)
	if mob == nil {
		return nil
	}
	m.GetControllerTable().AutoAggro(mob, ch)
	return nil
}
