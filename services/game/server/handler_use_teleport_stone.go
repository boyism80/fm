package server

import (
	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/protocol/request"
	"github.com/boyism80/fm/services/game/client"
	"github.com/boyism80/fm/services/game/constant"
)

type UseTeleportStone struct{}

func (UseTeleportStone) New(_ *GameServer) *UseTeleportStone {
	return &UseTeleportStone{}
}

func (h *UseTeleportStone) Handle(ctx *core.ClientContext, req *request.UseTeleportStone) error {
	gameClient, ok := ctx.Client.(*client.GameClient)
	if !ok {
		return nil
	}
	character := gameClient.GetCharacter()
	if character == nil {
		return nil
	}

	character.UseTeleportStone(ctx.ActorContext, constant.InventoryTypeConsume, int16(req.Slot), req.ItemID, req.Target.MapID, req.Target.Name)
	return nil
}
