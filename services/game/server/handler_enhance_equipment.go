package server

import (
	"fmt"
	"log"

	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/core/luax"
	"github.com/boyism80/fm/protocol/request"
	"github.com/boyism80/fm/services/game/client"
	"github.com/boyism80/fm/services/game/constant"
	"github.com/boyism80/fm/services/game/entity"
	lua "github.com/yuin/gopher-lua"
)

type EnhanceEquipment struct{}

func (EnhanceEquipment) New(_ *GameServer) *EnhanceEquipment {
	return &EnhanceEquipment{}
}

func (h *EnhanceEquipment) Handle(ctx *core.ClientContext, req *request.EnhanceEquipment) error {
	gameClient, ok := ctx.Client.(*client.GameClient)
	if !ok {
		log.Printf("Client is not a GameClient")
		return fmt.Errorf("client is not a GameClient")
	}

	ch := gameClient.GetCharacter()
	if ch == nil {
		return nil
	}

	h.runEnhanceScript(ch, req.ScrollSlot, req.TargetSlot, func(ok bool) {
		if !ok {
			ch.Listener.OnUpdateStats(ch, nil, true)
		}
	})
	return nil
}

func (*EnhanceEquipment) runEnhanceScript(ch *entity.Character, scrollSlot int16, targetSlot int16, fn func(bool)) {
	if ch == nil || scrollSlot <= 0 {
		fn(false)
		return
	}
	useInventory := ch.Inventory.Tabs[constant.InventoryTypeConsume]
	if useInventory == nil {
		fn(false)
		return
	}
	scrollItem := useInventory.Get(uint8(scrollSlot))
	if scrollItem == nil || scrollItem.GetCount() < 1 {
		fn(false)
		return
	}
	scrollConsume, ok := scrollItem.(*entity.Consume)
	if !ok || scrollConsume == nil {
		fn(false)
		return
	}
	scrollWz := scrollConsume.GetModel()
	if scrollWz == nil || scrollWz.GetID() == 0 {
		fn(false)
		return
	}
	mapInstance := ch.GetMap()
	if mapInstance == nil {
		fn(false)
		return
	}
	root := mapInstance.GetLuaRoot()
	if root == nil {
		fn(false)
		return
	}
	scriptPath := fmt.Sprintf("script/item/%d.lua", scrollWz.GetID())
	thread, err := luax.NewThread(root, scriptPath)
	if err != nil {
		fn(false)
		return
	}
	luax.CallAsync(root, thread, "on_scroll", ch, int32(scrollSlot), int32(targetSlot)).Then(func(value interface{}) (interface{}, error) {
		vals := luax.ResultValues(value)
		if len(vals) == 0 || vals[0] == nil {
			fn(false)
			return nil, nil
		}
		if b, ok := vals[0].(lua.LBool); ok {
			fn(bool(b))
			return nil, nil
		}
		fn(false)
		return nil, nil
	}).OnError(func(err error) {
		log.Printf("enhance script failed %s: %v", scriptPath, err)
		fn(false)
	})
}
