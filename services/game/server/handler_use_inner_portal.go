package server

import (
	"log"

	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/protocol/request"
	"github.com/boyism80/fm/types"
)

type UseInnerPortal struct {
	gs *GameServer
}

func (UseInnerPortal) New(gs *GameServer) *UseInnerPortal {
	return &UseInnerPortal{
		gs: gs,
	}
}

func (h *UseInnerPortal) Handle(ctx *core.ClientContext, req *request.UseInnerPortal) error {
	ch := h.gs.character(ctx)
	if ch == nil {
		return nil
	}
	if req.FieldKey != ch.FieldKey() {
		return nil
	}
	if err := ch.EnterInnerPortal(req.PortalName, types.Vector2[int16]{X: req.ToX, Y: req.ToY}); err != nil {
		log.Printf("UseInnerPortal: character=%d: %v", ch.GetID(), err)
	}
	return nil
}
