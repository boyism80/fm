package server

import (
	"fmt"
	"log"

	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/core/luax"
	"github.com/boyism80/fm/protocol/request"
	"github.com/boyism80/fm/services/game/client"
	"github.com/boyism80/fm/services/game/constant"
	lua "github.com/yuin/gopher-lua"
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

	unlock := func() {
		ch.Listener.OnUpdateStats(ch, nil, true)
	}

	if ch.GetHp() <= 0 {
		unlock()
		return nil
	}

	cash := ch.Inventory.Tabs[constant.InventoryTypeCash]
	if cash == nil {
		unlock()
		return nil
	}
	item := cash.Get(uint8(req.Slot))
	if item == nil || item.GetCount() < 1 || item.GetModel().GetID() != req.ItemID {
		unlock()
		return nil
	}

	mapInstance := ch.GetMap()
	if mapInstance == nil || mapInstance.GetLuaRoot() == nil {
		unlock()
		return nil
	}
	scriptPath := fmt.Sprintf("script/item/%d.lua", req.ItemID)
	thread, err := luax.NewThread(mapInstance.GetLuaRoot(), scriptPath)
	if err != nil {
		unlock()
		return nil
	}
	luax.SetConfiguration(thread, luax.Configuration{ActorContext: ctx.ActorContext})
	luax.CallAsync(mapInstance.GetLuaRoot(), thread, "on_cash", ch, req.ItemID, req.Text, req.Ear).Then(func(value interface{}) (interface{}, error) {
		vals := luax.ResultValues(value)
		if len(vals) == 0 || vals[0] != lua.LTrue {
			unlock()
			return nil, nil
		}
		item.Reduce(1)
		if item.GetCount() == 0 {
			cash.Remove(uint8(req.Slot))
			ch.Listener.OnRemoveInventorySlot(ch, constant.InventoryTypeCash, int16(req.Slot))
		} else {
			ch.Listener.OnInventorySlotUpdated(ch, constant.InventoryTypeCash, int16(req.Slot), item)
		}
		unlock()
		return nil, nil
	}).OnError(func(err error) {
		log.Printf("cash item script %s: %v", scriptPath, err)
		unlock()
	})
	return nil
}
