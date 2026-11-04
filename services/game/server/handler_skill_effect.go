package server

import (
	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/protocol/request"
)

type SkillEffect struct {
	gs *GameServer
}

func (SkillEffect) New(gs *GameServer) *SkillEffect { return &SkillEffect{gs: gs} }
func (h *SkillEffect) Handle(ctx *core.ClientContext, req *request.SkillEffect) error {
	ch := h.gs.character(ctx)
	if ch == nil {
		return nil
	}
	ch.ShowChargeSkillEffect(req.SkillID, req.Level, req.Flags, req.Speed)
	return nil
}
