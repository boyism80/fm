package server

import (
	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/protocol/request"
	"github.com/boyism80/fm/services/game"
	"github.com/boyism80/fm/services/game/client"
)

type FamilyPedigree struct {
	gs *GameServer
}

func (FamilyPedigree) New(gs *GameServer) *FamilyPedigree {
	return &FamilyPedigree{gs: gs}
}

func (h *FamilyPedigree) Handle(ctx *core.ClientContext, req *request.FamilyPedigree) error {
	gameClient, ok := ctx.Client.(*client.GameClient)
	if !ok {
		return nil
	}
	ch := gameClient.GetCharacter()
	if ch == nil || (req.Name != "" && !game.ClientText(req.Name).Valid(13)) {
		return nil
	}

	ch.Family().ShowPedigree(ctx.ActorContext, req.Name)
	return nil
}
