package server

import (
	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/protocol/request"
)

type UseDoor struct {
	gs *GameServer
}

func (UseDoor) New(gs *GameServer) *UseDoor {
	return &UseDoor{
		gs: gs,
	}
}

func (h *UseDoor) Handle(ctx *core.ClientContext, req *request.UseDoor) error {
	character := h.gs.character(ctx)
	if character == nil {
		return nil
	}

	m := character.GetMap()
	if m == nil {
		character.Listener().OnUpdateStats(character, nil, true)
		return nil
	}
	door := m.FindDoorOwner(req.OwnerCharacterID)
	if door == nil || !door.Use(ctx.ActorContext, character) {
		character.Listener().OnUpdateStats(character, nil, true)
	}
	return nil
}
