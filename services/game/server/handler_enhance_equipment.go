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

	if ok := h.runEnhanceScript(ch, req.ScrollSlot, req.TargetSlot); !ok {
		ch.Listener.OnUpdateStats(ch, nil, true)
	}
	return nil
}

func (*EnhanceEquipment) runEnhanceScript(ch *entity.Character, scrollSlot int16, targetSlot int16) bool {
	if ch == nil || scrollSlot <= 0 {
		return false
	}
	useInventory := ch.Inventory[constant.InventoryTypeConsume]
	if useInventory == nil {
		return false
	}
	scrollItem := useInventory.GetItem(uint8(scrollSlot))
	if scrollItem == nil || scrollItem.GetCount() < 1 {
		return false
	}
	scrollConsume, ok := scrollItem.(*entity.Consume)
	if !ok || scrollConsume == nil {
		return false
	}
	scrollWz := scrollConsume.GetModel()
	if scrollWz == nil || scrollWz.GetID() == 0 {
		return false
	}
	mapInstance := ch.GetMap()
	if mapInstance == nil {
		return false
	}
	root := mapInstance.GetLuaRoot()
	if root == nil {
		return false
	}
	scriptPath := fmt.Sprintf("script/item/%d.lua", scrollWz.GetID())
	thread, err := luax.NewThread(root, scriptPath)
	if err != nil {
		return false
	}
	ret, err := luax.Call(thread, "on_scroll", ch, int32(scrollSlot), int32(targetSlot))
	if err != nil {
		log.Printf("enhance script failed %s: %v", scriptPath, err)
		return false
	}
	if b, ok := ret.(lua.LBool); ok {
		return bool(b)
	}
	return false
}
