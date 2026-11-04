package server

import (
	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/protocol/request"
)

type QuestPotOpen struct {
	gs *GameServer
}

func (QuestPotOpen) New(gs *GameServer) *QuestPotOpen {
	return &QuestPotOpen{gs: gs}
}

func (h *QuestPotOpen) Handle(ctx *core.ClientContext, req *request.QuestPotOpen) error {
	ch, err := h.gs.character(ctx)
	if err != nil {
		return nil
	}

	ch.Quests().OpenPot(uint32(req.QuestID))
	return nil
}
