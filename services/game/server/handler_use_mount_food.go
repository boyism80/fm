package server

import (
	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/protocol/request"
)

type UseMountFood struct {
	gs *GameServer
}

func (UseMountFood) New(gs *GameServer) *UseMountFood {
	return &UseMountFood{
		gs: gs,
	}
}

func (h *UseMountFood) Handle(ctx *core.ClientContext, req *request.UseMountFood) error {
	character, err := h.gs.character(ctx)
	if err != nil {
		return nil
	}
	if character.UsingCashItem() {
		character.Listener().OnUpdateStats(character, nil, true)
		return nil
	}

	// The client locks its actions on use and waits for an unlock even on success.
	character.Mount().Feed(req.Slot, req.ItemID)
	character.Listener().OnUnlockAction(character)
	return nil
}
