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

	if ch.GetHp() <= 0 {
		ch.Listener.OnUpdateStats(ch, nil, true)
		return nil
	}

	useInventory := ch.Inventory[constant.INVENTORY_TYPE_CONSUME]
	if useInventory == nil {
		ch.Listener.OnUpdateStats(ch, nil, true)
		return nil
	}

	item := useInventory.GetItem(uint8(req.Slot))
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

	if !ch.ApplyConsumeEffect(consumeEnt) {
		ch.Listener.OnUpdateStats(ch, nil, true)
		return nil
	}

	item.Reduce(1)
	if item.GetCount() == 0 {
		useInventory.RemoveItem(uint8(req.Slot))
		ch.Listener.OnRemoveInventorySlot(ch, constant.INVENTORY_TYPE_CONSUME, int16(req.Slot))
	} else {
		ch.Listener.OnInventorySlotUpdated(ch, constant.INVENTORY_TYPE_CONSUME, int16(req.Slot), item)
	}
	ch.Listener.OnUpdateStats(ch, nil, true)

	return nil
}
