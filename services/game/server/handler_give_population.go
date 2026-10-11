package server

import (
	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/protocol/request"
	"github.com/boyism80/fm/services/game/client"
)

type GivePopulation struct {
	gs *GameServer
}

func (GivePopulation) New(gs *GameServer) *GivePopulation {
	return &GivePopulation{
		gs: gs,
	}
}

func (h *GivePopulation) Handle(ctx *core.ClientContext, req *request.GivePopulation) error {
	gameClient, ok := ctx.Client.(*client.GameClient)
	if !ok {
		return nil
	}
	ch := gameClient.GetCharacter()
	if ch == nil {
		return nil
	}

	ch.GivePopulation(req.TargetID, req.Mode)
	return nil
}
