package server

import (
	"fmt"

	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/protocol/request"
	"github.com/boyism80/fm/services/game/client"
)

type DenyGuildRequest struct {
	gs *GameServer
}

func (DenyGuildRequest) New(gs *GameServer) *DenyGuildRequest {
	return &DenyGuildRequest{
		gs: gs,
	}
}

func (h *DenyGuildRequest) Handle(ctx *core.ClientContext, req *request.DenyGuildRequest) error {
	gameClient, ok := ctx.Client.(*client.GameClient)
	if !ok {
		return fmt.Errorf("client is not a GameClient")
	}
	ch := gameClient.GetCharacter()
	if ch == nil {
		return fmt.Errorf("character not found")
	}
	ch.Guild.ClearInvites()
	return nil
}
