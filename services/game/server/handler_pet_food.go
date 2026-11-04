package server

import (
	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/protocol/request"
)

type PetFood struct {
	gs *GameServer
}

func (PetFood) New(gs *GameServer) *PetFood {
	return &PetFood{
		gs: gs,
	}
}

func (h *PetFood) Handle(ctx *core.ClientContext, req *request.PetFood) error {
	character, err := h.gs.character(ctx)
	if err != nil {
		return nil
	}
	if character.Pets().Active() == nil {
		return nil
	}
	if character.UsingCashItem() {
		character.Listener().OnUnlockAction(character)
		return nil
	}

	err = character.Pets().Active().Feed(req.Slot, req.ItemID)
	if err != nil {
		character.Listener().OnUnlockAction(character)
	}

	return nil
}
