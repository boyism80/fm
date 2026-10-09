package server

import (
	"fmt"

	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/protocol/request"
	"github.com/boyism80/fm/services/game/client"
	"github.com/boyism80/fm/services/game/constant"
)

type UseShopScanner struct {
	gs *GameServer
}

func (UseShopScanner) New(gs *GameServer) *UseShopScanner {
	return &UseShopScanner{
		gs: gs,
	}
}

func (h *UseShopScanner) Handle(ctx *core.ClientContext, req *request.UseShopScanner) error {
	client, ok := ctx.Client.(*client.GameClient)
	if !ok {
		return fmt.Errorf("client is not a GameClient")
	}

	ch := client.GetCharacter()
	if ch == nil {
		return nil
	}
	ch.UseShopScanner(ctx.ActorContext, constant.InventoryTypeConsume, req.Slot, req.ItemID, req.SearchID, req.HighFirst)
	return nil
}
