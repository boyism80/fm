package server

import (
	"fmt"
	"log"

	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/protocol/request"
	"github.com/boyism80/fm/services/game/client"
	"github.com/boyism80/fm/services/game/constant"
	"github.com/boyism80/fm/services/game/entity"
)

type UseItem struct{}

func (UseItem) New(_ *GameServer) *UseItem {
	return &UseItem{}
}

func (*UseItem) Handle(ctx *core.ClientContext, req *request.UseItem) error {
	client, ok := ctx.Client.(*client.GameClient)
	if !ok {
		log.Printf("Client is not a GameClient")
		return fmt.Errorf("client is not a GameClient")
	}

	ch := client.GetCharacter()
	if ch == nil {
		return nil
	}

	mapInstance := ch.GetMap()
	if mapInstance != nil && mapInstance.Wz.Limits(constant.FieldLimitPotion) {
		ch.Listener.OnUpdateStats(ch, nil, true)
		return nil
	}

	useInventory := ch.Inventory.Tabs[constant.InventoryTypeConsume]
	if useInventory == nil {
		ch.Listener.OnUpdateStats(ch, nil, true)
		return nil
	}

	item := useInventory.Get(uint8(req.Slot))
	if item == nil {
		ch.Listener.OnUpdateStats(ch, nil, true)
		return nil
	}

	if item.GetCount() < 1 {
		ch.Listener.OnUpdateStats(ch, nil, true)
		return nil
	}

	if item.GetModel().GetID() != req.ItemID {
		ch.Listener.OnUpdateStats(ch, nil, true)
		return nil
	}

	consumeEnt, ok := item.(*entity.Consume)
	if !ok || consumeEnt == nil {
		ch.Listener.OnUpdateStats(ch, nil, true)
		return nil
	}

	if !ch.UseConsume(consumeEnt) {
		ch.Listener.OnUpdateStats(ch, nil, true)
		return nil
	}

	ch.Inventory.RemoveItem(constant.InventoryTypeConsume, int16(req.Slot), 1)
	ch.Listener.OnUpdateStats(ch, nil, true)

	return nil
}
