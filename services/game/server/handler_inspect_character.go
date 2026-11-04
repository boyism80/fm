package server

import (
	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/protocol/request"
)

type InspectCharacter struct {
	gs *GameServer
}

func (InspectCharacter) New(gs *GameServer) *InspectCharacter {
	return &InspectCharacter{
		gs: gs,
	}
}

func (h *InspectCharacter) Handle(ctx *core.ClientContext, req *request.InspectCharacter) error {
	character, err := h.gs.character(ctx)
	if err != nil {
		return nil
	}

	character.Inspect(req.CharacterID)
	return nil
}
