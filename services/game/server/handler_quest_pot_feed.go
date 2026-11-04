package server

import (
	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/protocol/request"
)

type QuestPotFeed struct {
	gs *GameServer
}

func (QuestPotFeed) New(gs *GameServer) *QuestPotFeed {
	return &QuestPotFeed{gs: gs}
}

func (h *QuestPotFeed) Handle(ctx *core.ClientContext, req *request.QuestPotFeed) error {
	ch, err := h.gs.character(ctx)
	if err != nil {
		return nil
	}
	defer ch.Listener().OnUnlockAction(ch)

	ch.Quests().FeedPot(req.Slot, req.ItemID, req.QuestID)
	return nil
}
