package server

import (
	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/protocol/request"
	"github.com/boyism80/fm/services/game/client"
	"github.com/boyism80/fm/services/game/domain"
)

type FamilyAccept struct {
	gs *GameServer
}

func (FamilyAccept) New(gs *GameServer) *FamilyAccept {
	return &FamilyAccept{gs: gs}
}

func (h *FamilyAccept) Handle(ctx *core.ClientContext, req *request.FamilyAccept) error {
	gameClient, ok := ctx.Client.(*client.GameClient)
	if !ok {
		return nil
	}
	ch := gameClient.GetCharacter()
	if ch == nil || (req.Name != "" && !domain.ClientText(req.Name).Valid(13)) {
		return nil
	}

	ch.Family.AnswerInvite(ctx.ActorContext, req.SeniorID, req.Name, req.Accepted)
	return nil
}
