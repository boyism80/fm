package server

import (
	"fmt"
	"log"

	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/protocol/request"
	"github.com/boyism80/fm/services/game/client"
)

type UseReturnScroll struct{}

func (UseReturnScroll) New(_ *GameServer) *UseReturnScroll {
	return &UseReturnScroll{}
}

func (*UseReturnScroll) Handle(ctx *core.ClientContext, req *request.UseReturnScroll) error {
	client, ok := ctx.Client.(*client.GameClient)
	if !ok {
		log.Printf("Client is not a GameClient")
		return fmt.Errorf("client is not a GameClient")
	}

	ch := client.GetCharacter()
	if ch == nil {
		return nil
	}

	return ch.UseReturnScroll(ctx.ActorContext, int16(req.Slot), req.ItemID)
}
