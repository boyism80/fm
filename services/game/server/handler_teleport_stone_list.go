package server

import (
	"github.com/boyism80/fm/core"
	pconst "github.com/boyism80/fm/protocol/constant"
	"github.com/boyism80/fm/protocol/request"
)

type TeleportStoneList struct {
	gs *GameServer
}

func (TeleportStoneList) New(gs *GameServer) *TeleportStoneList {
	return &TeleportStoneList{gs: gs}
}

func (h *TeleportStoneList) Handle(ctx *core.ClientContext, req *request.TeleportStoneList) error {
	character := h.gs.character(ctx)
	if character == nil {
		return nil
	}

	switch req.Action {
	case pconst.TeleportStoneActionRegister:
		character.TeleportStones().Register(req.VIP)
	case pconst.TeleportStoneActionRemove:
		character.TeleportStones().Remove(req.VIP, req.MapID)
	}
	return nil
}
