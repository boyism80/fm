package server

import (
	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/protocol/request"
	"github.com/boyism80/fm/services/game/client"
)

type MoveSummon struct {
	gs *GameServer
}

func (MoveSummon) New(gs *GameServer) *MoveSummon {
	return &MoveSummon{
		gs: gs,
	}
}

func (h *MoveSummon) Handle(ctx *core.ClientContext, req *request.MoveSummon) error {
	gameClient, ok := ctx.Client.(*client.GameClient)
	if !ok {
		return nil
	}
	character := gameClient.GetCharacter()
	if character == nil {
		return nil
	}
	summon := character.Summons().FindByOID(req.OID)
	if summon == nil {
		return nil
	}
	summon.Move(req.Fragments)
	return nil
}
