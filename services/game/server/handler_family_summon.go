package server

import (
	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/protocol/request"
	"github.com/boyism80/fm/services/game/client"
	"github.com/boyism80/fm/services/game/domain"
)

type FamilySummon struct {
	gs *GameServer
}

func (FamilySummon) New(gs *GameServer) *FamilySummon {
	return &FamilySummon{gs: gs}
}

func (h *FamilySummon) Handle(ctx *core.ClientContext, req *request.FamilySummon) error {
	gameClient, ok := ctx.Client.(*client.GameClient)
	if !ok {
		return nil
	}
	ch := gameClient.GetCharacter()
	if ch == nil || (req.Name != "" && !domain.ClientText(req.Name).Valid(13)) {
		return nil
	}

	ch.Family.AnswerSummon(req.Name, req.Accepted)
	return nil
}
