package server

import (
	"fmt"
	"log"

	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/protocol/request"
	"github.com/boyism80/fm/services/game/client"
)

type UseCatchItem struct{}

func (UseCatchItem) New(_ *GameServer) *UseCatchItem {
	return &UseCatchItem{}
}

func (*UseCatchItem) Handle(ctx *core.ClientContext, req *request.UseCatchItem) error {
	client, ok := ctx.Client.(*client.GameClient)
	if !ok {
		log.Printf("Client is not a GameClient")
		return fmt.Errorf("client is not a GameClient")
	}

	ch := client.GetCharacter()
	if ch == nil {
		return nil
	}

	ch.UseCatchItem(int16(req.Slot), req.ItemID, req.MobOID)
	ch.Listener.OnUpdateStats(ch, nil, true)
	return nil
}
