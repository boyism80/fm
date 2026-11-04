package server

import (
	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/protocol/request"
)

type MonsterBomb struct {
	gs *GameServer
}

func (MonsterBomb) New(gs *GameServer) *MonsterBomb {
	return &MonsterBomb{gs: gs}
}

func (h *MonsterBomb) Handle(ctx *core.ClientContext, req *request.MonsterBomb) error {
	character, err := h.gs.character(ctx)
	if err != nil {
		return err
	}
	m := character.GetMap()
	if m == nil {
		return nil
	}
	mob := m.GetMob(req.OID)
	if mob == nil {
		return nil
	}
	mob.SelfDestruct(character)
	return nil
}
