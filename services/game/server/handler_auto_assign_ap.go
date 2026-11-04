package server

import (
	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/maple/constant"
	"github.com/boyism80/fm/protocol/request"
	"github.com/boyism80/fm/services/game"
)

type AutoAssignAP struct {
	gs *GameServer
}

func (AutoAssignAP) New(gs *GameServer) *AutoAssignAP {
	return &AutoAssignAP{
		gs: gs,
	}
}

func (h *AutoAssignAP) Handle(ctx *core.ClientContext, req *request.AutoAssignAP) error {
	character, err := h.gs.character(ctx)
	if err != nil {
		return err
	}

	totalAmount := uint32(0)
	for _, entry := range req.Entries {
		if int32(entry.Amount) < 0 {
			return nil
		}
		totalAmount += entry.Amount
	}
	if totalAmount == 0 {
		return nil
	}

	emptyStats := map[constant.Stat]int32{}
	character.Listener().OnUpdateStats(character, emptyStats, true)

	entries := make([]game.APEntry, len(req.Entries))
	for i, e := range req.Entries {
		entries[i] = game.APEntry{Stat: constant.StatType(e.Stat), Amount: uint16(e.Amount)}
	}
	if !character.Points().AssignAll(entries) {
		character.Listener().OnUpdateStats(character, emptyStats, true)
	}
	return nil
}
