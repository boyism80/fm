package server

import (
	"log"

	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/protocol/request"
)

type UseChair struct {
	gs *GameServer
}

func (UseChair) New(gs *GameServer) *UseChair {
	return &UseChair{
		gs: gs,
	}
}

func (h *UseChair) Handle(ctx *core.ClientContext, req *request.UseChair) error {
	character := h.gs.character(ctx)
	if character == nil {
		return nil
	}

	if character.GetMap() == nil {
		return nil
	}

	if err := character.SitOnChair(req.ItemID); err != nil {
		log.Printf("Chair %d: %v", req.ItemID, err)
	}
	character.Listener().OnUpdateStats(character, nil, true)

	return nil
}
