package server

import (
	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/protocol/request"
)

type PetCommand struct {
	gs *GameServer
}

func (PetCommand) New(gs *GameServer) *PetCommand {
	return &PetCommand{
		gs: gs,
	}
}

func (h *PetCommand) Handle(ctx *core.ClientContext, req *request.PetCommand) error {
	character := h.gs.character(ctx)
	if character == nil {
		return nil
	}
	if character.Pets().Active() == nil {
		return nil
	}

	character.Pets().Active().Command(req.Index, req.CalledByName)

	return nil
}
