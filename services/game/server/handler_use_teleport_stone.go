package server

import (
	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/maple/constant"
	"github.com/boyism80/fm/protocol/request"
)

type UseTeleportStone struct {
	gs *GameServer
}

func (UseTeleportStone) New(gs *GameServer) *UseTeleportStone {
	return &UseTeleportStone{gs: gs}
}

func (h *UseTeleportStone) Handle(ctx *core.ClientContext, req *request.UseTeleportStone) error {
	character := h.gs.character(ctx)
	if character == nil {
		return nil
	}

	character.TeleportStones().Use(ctx.ActorContext, constant.InventoryTypeConsume, int16(req.Slot), req.ItemID, req.Target.MapID, req.Target.Name)
	return nil
}
