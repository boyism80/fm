package server

import (
	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/maple/constant"
	"github.com/boyism80/fm/protocol/request"
)

type DistributeSP struct {
	gs *GameServer
}

func (DistributeSP) New(gs *GameServer) *DistributeSP {
	return &DistributeSP{
		gs: gs,
	}
}

func (h *DistributeSP) Handle(ctx *core.ClientContext, req *request.DistributeSP) error {
	character := h.gs.character(ctx)
	if character == nil {
		return nil
	}

	if !character.Points().DistributeSkill(req.SkillID) {
		character.Listener().OnUpdateStats(character, map[constant.Stat]int32{}, true)
	}
	return nil
}
