package server

import (
	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/protocol/request"
)

type SummonPet struct {
	gs *GameServer
}

func (SummonPet) New(gs *GameServer) *SummonPet {
	return &SummonPet{
		gs: gs,
	}
}

func (h *SummonPet) Handle(ctx *core.ClientContext, req *request.SummonPet) error {
	character, err := h.gs.character(ctx)
	if err != nil {
		return nil
	}

	err = character.Pets().Summon(req.Slot)
	if err != nil {
		character.Listener().OnUnlockAction(character)
	}
	return nil
}
