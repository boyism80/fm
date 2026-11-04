package server

import (
	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/protocol/request"
)

type MovePet struct {
	gs *GameServer
}

func (MovePet) New(gs *GameServer) *MovePet {
	return &MovePet{
		gs: gs,
	}
}

func (h *MovePet) Handle(ctx *core.ClientContext, req *request.MovePet) error {
	character := h.gs.character(ctx)
	if character == nil {
		return nil
	}
	if character.Pets().Active() == nil {
		return nil
	}

	character.Pets().Active().Move(req.Position, req.Fragments)

	return nil
}
