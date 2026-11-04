package server

import (
	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/maple/constant"
	"github.com/boyism80/fm/protocol/request"
	"github.com/boyism80/fm/services/game"
)

type DistributeAP struct {
	gs *GameServer
}

func (DistributeAP) New(gs *GameServer) *DistributeAP {
	return &DistributeAP{
		gs: gs,
	}
}

func (h *DistributeAP) Handle(ctx *core.ClientContext, req *request.DistributeAP) error {
	character, err := h.gs.character(ctx)
	if err != nil {
		return err
	}

	stat := constant.StatType(req.StatType)
	assigned := false
	switch stat {
	case constant.StatTypeHP, constant.StatTypeMP:
		assigned = character.Points().AssignHPMP(stat)
	default:
		assigned = character.Points().AssignStats([]game.APEntry{{Stat: stat, Amount: 1}})
	}
	if !assigned {
		character.Listener().OnUpdateStats(character, map[constant.Stat]int32{}, true)
	}
	return nil
}
