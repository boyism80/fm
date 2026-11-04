package server

import (
	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/protocol/request"
)

type PetAutoPotion struct {
	gs *GameServer
}

func (PetAutoPotion) New(gs *GameServer) *PetAutoPotion {
	return &PetAutoPotion{
		gs: gs,
	}
}

func (h *PetAutoPotion) Handle(ctx *core.ClientContext, req *request.PetAutoPotion) error {
	character := h.gs.character(ctx)
	if character == nil {
		return nil
	}
	if character.Pets().Active() == nil {
		return nil
	}
	if character.UsingCashItem() {
		character.Listener().OnUnlockAction(character)
		return nil
	}

	err := character.Pets().Active().UsePotion(req.Slot, req.ItemID)
	if err != nil {
		character.Listener().OnUnlockAction(character)
	}

	return nil
}
