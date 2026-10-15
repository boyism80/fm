package server

import (
	"fmt"
	"log"

	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/protocol/request"
	"github.com/boyism80/fm/services/game/client"
)

type ItemMaker struct{}

func (ItemMaker) New(_ *GameServer) *ItemMaker {
	return &ItemMaker{}
}

func (*ItemMaker) Handle(ctx *core.ClientContext, req *request.ItemMaker) error {
	gameClient, ok := ctx.Client.(*client.GameClient)
	if !ok {
		log.Printf("Client is not a GameClient")
		return fmt.Errorf("client is not a GameClient")
	}

	ch := gameClient.GetCharacter()
	if ch == nil {
		return nil
	}
	defer ch.Listener.OnUnlockAction(ch)
	if req.Malformed {
		return nil
	}

	ch.Inventory.Make(req.Type, req.ItemID, req.Catalyst, req.Enhancers, int16(req.Slot))
	return nil
}
