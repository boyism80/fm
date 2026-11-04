package server

import (
	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/protocol/request"
	"github.com/boyism80/fm/services/game"
)

type SkillMacro struct {
	gs *GameServer
}

func (SkillMacro) New(gs *GameServer) *SkillMacro { return &SkillMacro{gs: gs} }
func (h *SkillMacro) Handle(ctx *core.ClientContext, req *request.SkillMacro) error {
	ch, err := h.gs.character(ctx)
	if err != nil {
		return nil
	}
	if req.Count != uint8(len(req.Entries)) {
		return nil
	}
	entries := make([]game.SkillMacro, 0, len(req.Entries))
	for _, entry := range req.Entries {
		if entry.Shout > 1 {
			return nil
		}
		entries = append(entries, game.SkillMacro{Name: entry.Name, Shout: entry.Shout == 1, Skills: entry.Skills})
	}
	ch.Skills().Macros().Replace(entries)
	return nil
}
