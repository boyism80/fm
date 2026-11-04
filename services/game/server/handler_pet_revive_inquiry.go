package server

import (
	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/protocol/request"
)

type PetReviveInquiry struct {
	gs *GameServer
}

func (PetReviveInquiry) New(gs *GameServer) *PetReviveInquiry {
	return &PetReviveInquiry{
		gs: gs,
	}
}

func (h *PetReviveInquiry) Handle(ctx *core.ClientContext, req *request.PetReviveInquiry) error {
	character, err := h.gs.character(ctx)
	if err != nil {
		return nil
	}

	character.Listener().OnUnlockAction(character)
	return nil
}
