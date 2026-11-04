package server

import (
	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/protocol/request"
)

type LieDetectorItem struct {
	gs *GameServer
}

func (LieDetectorItem) New(gs *GameServer) *LieDetectorItem {
	return &LieDetectorItem{
		gs: gs,
	}
}

func (h *LieDetectorItem) Handle(ctx *core.ClientContext, req *request.LieDetectorItem) error {
	character := h.gs.character(ctx)
	if character == nil {
		return nil
	}

	err := character.UseLieDetectorItem(req.Target, req.Slot, req.ItemID)
	if err != nil {
		character.Listener().OnUnlockAction(character)
	}

	return nil
}
