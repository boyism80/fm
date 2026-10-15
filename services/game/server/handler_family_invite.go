package server

import (
	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/protocol/request"
	"github.com/boyism80/fm/services/game"
	"github.com/boyism80/fm/services/game/client"
)

type FamilyInvite struct {
	gs *GameServer
}

func (FamilyInvite) New(gs *GameServer) *FamilyInvite {
	return &FamilyInvite{gs: gs}
}

func (h *FamilyInvite) Handle(ctx *core.ClientContext, req *request.FamilyInvite) error {
	gameClient, ok := ctx.Client.(*client.GameClient)
	if !ok {
		return nil
	}
	ch := gameClient.GetCharacter()
	if ch == nil || (req.Name != "" && !game.ClientText(req.Name).Valid(13)) {
		return nil
	}

	ch.Family.Invite(req.Name)
	return nil
}
