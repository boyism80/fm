package server

import (
	"fmt"

	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/protocol/request"
	"github.com/boyism80/fm/services/game/client"
)

type ShopScannerWarp struct {
	gs *GameServer
}

func (ShopScannerWarp) New(gs *GameServer) *ShopScannerWarp {
	return &ShopScannerWarp{
		gs: gs,
	}
}

func (h *ShopScannerWarp) Handle(ctx *core.ClientContext, req *request.ShopScannerWarp) error {
	client, ok := ctx.Client.(*client.GameClient)
	if !ok {
		return fmt.Errorf("client is not a GameClient")
	}

	ch := client.GetCharacter()
	if ch == nil {
		return nil
	}
	ch.RejectMiniRoom(ch.VisitShopBySearch(ctx.ActorContext, req.SN, req.MapID))
	return nil
}
