package server

import (
	"fmt"
	"log"

	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/protocol/request"
	"github.com/boyism80/fm/services/game/client"
	"github.com/boyism80/fm/services/game/constant"
	"github.com/boyism80/fm/services/game/entity"
	"github.com/boyism80/fm/services/game/wz"
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

	useInventory := ch.Inventory[constant.InventoryTypeConsume]
	if useInventory == nil {
		ch.Listener.OnUpdateStats(ch, nil, true)
		return nil
	}

	item := useInventory.GetItem(uint8(req.Slot))
	if item == nil || item.GetCount() < 1 || item.GetModel().GetID() != req.ItemID {
		ch.Listener.OnUpdateStats(ch, nil, true)
		return nil
	}

	consumeEnt, ok := item.(*entity.Consume)
	if !ok || consumeEnt == nil {
		ch.Listener.OnUpdateStats(ch, nil, true)
		return nil
	}

	consumeModel, ok := consumeEnt.GetModel().(*wz.Consume)
	if !ok || consumeModel == nil || consumeModel.MoveTo == -1 {
		ch.Listener.OnUpdateStats(ch, nil, true)
		return nil
	}
	if ch.GameWorld == nil {
		ch.Listener.OnUpdateStats(ch, nil, true)
		return nil
	}

	var targetID uint32
	switch consumeModel.MoveTo {
	case wz.ConsumeMoveToReturnMap:
		m := ch.GetMap()
		if m == nil || m.Wz == nil || m.Wz.ReturnMapId <= 0 {
			ch.Listener.OnUpdateStats(ch, nil, true)
			return nil
		}
		targetID = uint32(m.Wz.ReturnMapId)
	default:
		if consumeModel.MoveTo <= 0 {
			ch.Listener.OnUpdateStats(ch, nil, true)
			return nil
		}
		targetID = uint32(consumeModel.MoveTo)
	}

	target := ch.GameWorld.GetMapSystem().Get(targetID)
	if target == nil {
		ch.Listener.OnUpdateStats(ch, nil, true)
		return nil
	}
	if err := ch.Warp(target, 0); err != nil {
		return err
	}

	item.Reduce(1)
	if item.GetCount() == 0 {
		useInventory.RemoveItem(uint8(req.Slot))
		ch.Listener.OnRemoveInventorySlot(ch, constant.InventoryTypeConsume, int16(req.Slot))
	} else {
		ch.Listener.OnInventorySlotUpdated(ch, constant.InventoryTypeConsume, int16(req.Slot), item)
	}
	ch.Listener.OnUpdateStats(ch, nil, true)
	return nil
}
