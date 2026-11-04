package server

import (
	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/protocol/request"
)

type HypnotizeDamage struct {
	gs *GameServer
}

func (HypnotizeDamage) New(gs *GameServer) *HypnotizeDamage {
	return &HypnotizeDamage{gs: gs}
}

func (h *HypnotizeDamage) Handle(ctx *core.ClientContext, req *request.HypnotizeDamage) error {
	character, err := h.gs.character(ctx)
	if err != nil {
		return err
	}
	if req.CharacterID != character.GetID() {
		return nil
	}
	character.HypnotizeDamage(req.SourceOID, req.TargetOID, req.Damage)
	return nil
}
