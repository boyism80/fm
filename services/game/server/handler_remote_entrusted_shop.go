package server

import (
	"fmt"

	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/protocol/request"
	"github.com/boyism80/fm/services/game/client"
)

type RemoteEntrustedShop struct {
	gs *GameServer
}

func (RemoteEntrustedShop) New(gs *GameServer) *RemoteEntrustedShop {
	return &RemoteEntrustedShop{
		gs: gs,
	}
}

func (h *RemoteEntrustedShop) Handle(ctx *core.ClientContext, req *request.RemoteEntrustedShop) error {
	client, ok := ctx.Client.(*client.GameClient)
	if ok == false {
		return fmt.Errorf("client is not a GameClient")
	}

	ch := client.GetCharacter()
	if ch == nil {
		return nil
	}
	ch.RejectMiniRoom(ch.UseRemoteEntrustedShop(ctx.ActorContext, req.Slot))
	return nil
}
