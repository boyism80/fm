package server

import (
	"fmt"

	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/protocol/request"
	"github.com/boyism80/fm/services/game/client"
)

type UseHiredMerchant struct {
	gs *GameServer
}

func (UseHiredMerchant) New(gs *GameServer) *UseHiredMerchant {
	return &UseHiredMerchant{
		gs: gs,
	}
}

func (h *UseHiredMerchant) Handle(ctx *core.ClientContext, req *request.UseHiredMerchant) error {
	client, ok := ctx.Client.(*client.GameClient)
	if ok == false {
		return fmt.Errorf("client is not a GameClient")
	}

	ch := client.GetCharacter()
	if ch == nil {
		return nil
	}
	ch.UseHiredMerchant(ctx.ActorContext)
	return nil
}
