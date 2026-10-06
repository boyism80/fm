package server

import (
	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/protocol/request"
	"github.com/boyism80/fm/services/game/client"
)

type CancelBuff struct {
	gs *GameServer
}

func (CancelBuff) New(gs *GameServer) *CancelBuff {
	return &CancelBuff{
		gs: gs,
	}
}

func (h *CancelBuff) Handle(ctx *core.ClientContext, req *request.CancelBuff) error {
	client, ok := ctx.Client.(*client.GameClient)
	if !ok {
		return nil
	}
	character := client.GetCharacter()
	if character == nil {
		return nil
	}
	character.Buffs.CancelBySource(req.SourceID)
	return nil
}
