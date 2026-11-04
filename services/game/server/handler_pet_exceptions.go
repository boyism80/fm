package server

import (
	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/protocol/request"
)

type PetExceptions struct {
	gs *GameServer
}

func (PetExceptions) New(gs *GameServer) *PetExceptions {
	return &PetExceptions{
		gs: gs,
	}
}

func (h *PetExceptions) Handle(ctx *core.ClientContext, req *request.PetExceptions) error {
	character, err := h.gs.character(ctx)
	if err != nil {
		return nil
	}
	if character.Pets().Active() == nil {
		return nil
	}

	character.Pets().Active().SetExceptions(req.ItemIDs)

	return nil
}
