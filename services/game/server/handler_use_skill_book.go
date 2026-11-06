package server

import (
	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/protocol/request"
)

type UseSkillBook struct {
	gs *GameServer
}

func (UseSkillBook) New(gs *GameServer) *UseSkillBook {
	return &UseSkillBook{gs: gs}
}

func (h *UseSkillBook) Handle(ctx *core.ClientContext, req *request.UseSkillBook) error {
	ch := h.gs.character(ctx)
	if ch == nil {
		return nil
	}
	defer ch.Listener().OnUnlockAction(ch)
	if ch.Occupied() {
		return nil
	}

	ch.UseSkillBook(req.Slot, req.ItemID)
	return nil
}
