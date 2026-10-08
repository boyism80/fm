package server

import (
	"fmt"

	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/protocol/request"
	"github.com/boyism80/fm/services/game/client"
)

type UseEntrustedShop struct {
	gs *GameServer
}

func (UseEntrustedShop) New(gs *GameServer) *UseEntrustedShop {
	return &UseEntrustedShop{
		gs: gs,
	}
}

func (h *UseEntrustedShop) Handle(ctx *core.ClientContext, req *request.UseEntrustedShop) error {
	client, ok := ctx.Client.(*client.GameClient)
	if ok == false {
		return fmt.Errorf("client is not a GameClient")
	}

	ch := client.GetCharacter()
	if ch == nil {
		return nil
	}
	ch.UseEntrustedShop(ctx.ActorContext)
	return nil
}
