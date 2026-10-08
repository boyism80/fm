package server

import (
	"fmt"
	"log"

	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/protocol/request"
	"github.com/boyism80/fm/services/game/client"
	"github.com/boyism80/fm/services/game/constant"
)

type UseCashItem struct{}

func (UseCashItem) New(_ *GameServer) *UseCashItem {
	return &UseCashItem{}
}

func (*UseCashItem) Handle(ctx *core.ClientContext, req *request.UseCashItem) error {
	client, ok := ctx.Client.(*client.GameClient)
	if !ok {
		log.Printf("Client is not a GameClient")
		return fmt.Errorf("client is not a GameClient")
	}

	ch := client.GetCharacter()
	if ch == nil {
		return nil
	}

	if constant.IsCashTeleportStone(req.ItemID) {
		ch.TeleportStones.Use(ctx.ActorContext, constant.InventoryTypeCash, int16(req.Slot), req.ItemID, req.Target.MapID, req.Target.Name)
		return nil
	}
	if constant.IsShopScanner(req.ItemID) {
		ch.UseShopScanner(ctx.ActorContext, constant.InventoryTypeCash, int16(req.Slot), req.ItemID, req.SearchID, req.HighFirst)
		return nil
	}
	ch.UseCashItem(ctx.ActorContext, int16(req.Slot), req.ItemID, req.Text, req.Ear, req.PetSN)
	return nil
}
