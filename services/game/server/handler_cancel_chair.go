package server

import (
	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/protocol/request"
)

type CancelChair struct {
	gs *GameServer
}

func (CancelChair) New(gs *GameServer) *CancelChair {
	return &CancelChair{
		gs: gs,
	}
}

func (h *CancelChair) Handle(ctx *core.ClientContext, req *request.CancelChair) error {
	character, err := h.gs.character(ctx)
	if err != nil {
		return nil
	}

	if req.ChairID == -1 {
		character.StandUp()
	} else {
		character.SitOnMapSeat(req.ChairID)
	}

	character.Listener().OnUpdateStats(character, nil, true)

	return nil
}
