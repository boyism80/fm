package server

import (
	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/protocol/dto"
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
	mapInstance := character.GetMap()
	if mapInstance == nil {
		return nil
	}

	summon := mapInstance.GetSummon(req.OID)
	if summon == nil {
		return nil
	}
	if summon.OwnerID != character.GetID() {
		return nil
	}

	start := summon.Position
	for _, m := range req.Fragments {
		if move, ok := m.(*dto.AbsoluteLifeMovement); ok {
			summon.Position = move.Position
		}
		summon.Stance = m.GetStance()
	}

	character.Listener.OnSummonMove(character, summon, start, req.Fragments)
	return nil
}
