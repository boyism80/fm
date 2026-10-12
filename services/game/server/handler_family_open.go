package server

import (
	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/protocol/request"
	"github.com/boyism80/fm/services/game/client"
)

type FamilyOpen struct {
	gs *GameServer
}

func (FamilyOpen) New(gs *GameServer) *FamilyOpen {
	return &FamilyOpen{gs: gs}
}

func (h *FamilyOpen) Handle(ctx *core.ClientContext, req *request.FamilyOpen) error {
	gameClient, ok := ctx.Client.(*client.GameClient)
	if !ok {
		return nil
	}
	ch := gameClient.GetCharacter()
	if ch == nil {
		return nil
	}

	ch.Family.Open()
	return nil
}
