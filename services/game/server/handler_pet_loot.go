package server

import (
	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/protocol/request"
)

type PetLoot struct {
	gs *GameServer
}

func (PetLoot) New(gs *GameServer) *PetLoot {
	return &PetLoot{
		gs: gs,
	}
}

func (h *PetLoot) Handle(ctx *core.ClientContext, req *request.PetLoot) error {
	character, err := h.gs.character(ctx)
	if err != nil {
		return nil
	}
	if character.Pets().Active() == nil {
		return nil
	}
	if character.Spectating() || character.Occupied() || req.FieldKey != character.FieldKey() {
		character.Listener().OnUnlockAction(character)
		return nil
	}

	err = character.Pets().Active().Loot(req.OID, req.Position)
	if err != nil {
		character.Listener().OnUnlockAction(character)
	}

	return nil
}
