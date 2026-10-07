package server

import (
	"github.com/boyism80/fm/core"
	pconst "github.com/boyism80/fm/protocol/constant"
	"github.com/boyism80/fm/protocol/request"
	"github.com/boyism80/fm/services/game/client"
)

type TeleportStoneList struct{}

func (TeleportStoneList) New(_ *GameServer) *TeleportStoneList {
	return &TeleportStoneList{}
}

func (h *TeleportStoneList) Handle(ctx *core.ClientContext, req *request.TeleportStoneList) error {
	gameClient, ok := ctx.Client.(*client.GameClient)
	if !ok {
		return nil
	}
	character := gameClient.GetCharacter()
	if character == nil {
		return nil
	}

	switch req.Action {
	case pconst.TeleportStoneActionRegister:
		character.RegisterTeleportStone(req.VIP)
	case pconst.TeleportStoneActionRemove:
		character.RemoveTeleportStone(req.VIP, req.MapID)
	}
	return nil
}
