package server

import (
	"fmt"
	"log"

	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/protocol/request"
	"github.com/boyism80/fm/services/game/client"
)

type UseRewardBox struct{}

func (UseRewardBox) New(_ *GameServer) *UseRewardBox {
	return &UseRewardBox{}
}

func (*UseRewardBox) Handle(ctx *core.ClientContext, req *request.UseRewardBox) error {
	gameClient, ok := ctx.Client.(*client.GameClient)
	if !ok {
		log.Printf("Client is not a GameClient")
		return fmt.Errorf("client is not a GameClient")
	}

	ch := gameClient.GetCharacter()
	if ch == nil {
		return nil
	}
	defer ch.Listener().OnUnlockAction(ch)

	ch.UseRewardBox(req.Slot, req.ItemID)
	return nil
}
