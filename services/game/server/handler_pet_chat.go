package server

import (
	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/protocol/request"
)

type PetChat struct {
	gs *GameServer
}

func (PetChat) New(gs *GameServer) *PetChat {
	return &PetChat{
		gs: gs,
	}
}

func (h *PetChat) Handle(ctx *core.ClientContext, req *request.PetChat) error {
	character, err := h.gs.character(ctx)
	if err != nil {
		return nil
	}
	if character.Pets().Active() == nil {
		return nil
	}

	character.Pets().Active().Chat(req.Type, req.Action, req.Text)

	return nil
}
