package server

import (
	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/protocol/request"
	"github.com/boyism80/fm/services/game/client"
)

type DenyAllianceRequest struct {
	gs *GameServer
}

func (DenyAllianceRequest) New(gs *GameServer) *DenyAllianceRequest {
	return &DenyAllianceRequest{
		gs: gs,
	}
}

func (h *DenyAllianceRequest) Handle(ctx *core.ClientContext, req *request.DenyAllianceRequest) error {
	gameClient, ok := ctx.Client.(*client.GameClient)
	if !ok {
		return nil
	}
	ch := gameClient.GetCharacter()
	if ch == nil {
		return nil
	}
	op := &AllianceOperation{gs: h.gs}
	return op.handleDenyInvite(ctx, ch)
}
