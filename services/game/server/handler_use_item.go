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
	"github.com/boyism80/fm/services/game/wz"
)

type UseItem struct {
	gs *GameServer
}

func (UseItem) New(gs *GameServer) *UseItem {
	return &UseItem{
		gs: gs,
	}
}

func (h *UseItem) Handle(ctx *core.ClientContext, req *request.UseItem) error {
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

	resources := h.gs.GetResources()
	if resources == nil {
		ch.Listener.OnUpdateStats(ch, nil, true)
		return nil
	}

	wzItem := resources.Items[req.ItemID]
	if wzItem == nil {
		ch.Listener.OnUpdateStats(ch, nil, true)
		return nil
	}

	consumeItem, ok := wzItem.(*wz.Consume)
	if !ok {
		ch.Listener.OnUpdateStats(ch, nil, true)
		return nil
	}

	appliedItemBuff := h.applyConsumeBuff(ch, consumeItem)
	appliedRecovery := h.applyConsumeRecovery(ch, consumeItem)
	if !appliedRecovery && !appliedItemBuff {
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
	h.callActiveItemScript(ch, item)

	return nil
}

func (h *UseItem) applyConsumeBuff(ch *entity.Character, consumeItem *wz.Consume) bool {
	if ch == nil || consumeItem == nil {
		return false
	}
	buffValues := consumeItem.BuffSpecValues()
	if len(buffValues) == 0 {
		return false
	}
	ch.Buffs.AddItemBuff(consumeItem, consumeItem.BuffDuration, buffValues)
	return true
}

func (h *UseItem) applyConsumeRecovery(ch *entity.Character, consumeItem *wz.Consume) bool {
	if ch == nil || consumeItem == nil {
		return false
	}
	healMul := ch.PotionHealMultiplierPercent()
	hpChange := int(consumeItem.ActiveEffect.HP) * healMul / 100
	mpChange := int(consumeItem.ActiveEffect.MP) * healMul / 100
	hpRate := int(consumeItem.ActiveEffect.HPRate)
	mpRate := int(consumeItem.ActiveEffect.MPRate)
	if hpRate > 0 {
		hpChange += int(ch.GetMaxHp()) * hpRate / 100
	}
	if mpRate > 0 {
		mpChange += int(ch.GetMaxMp()) * mpRate / 100
	}
	if hpChange == 0 && mpChange == 0 {
		return false
	}

	stats := make(map[constant.Stat]int32)
	if hpChange != 0 {
		newHP := int(ch.GetHp()) + hpChange
		if newHP < 1 {
			newHP = 1
		}
		maxHp := int(ch.GetMaxHp())
		if newHP > maxHp {
			newHP = maxHp
		}
		if newHP != int(ch.GetHp()) {
			ch.SetHp(uint32(newHP), false)
			stats[constant.STAT_HP] = int32(ch.GetHp())
		}
	}
	if mpChange != 0 {
		newMP := int(ch.GetMp()) + mpChange
		if newMP < 0 {
			newMP = 0
		}
		maxMp := int(ch.GetMaxMp())
		if newMP > maxMp {
			newMP = maxMp
		}
		if newMP != int(ch.GetMp()) {
			ch.SetMp(uint32(newMP), false)
			stats[constant.STAT_MP] = int32(ch.GetMp())
		}
	}
	if len(stats) == 0 {
		return false
	}
	ch.Listener.OnUpdateStats(ch, stats, false)
	return true
}

func (h *UseItem) callActiveItemScript(ch *entity.Character, item entity.Item) {
	if ch == nil || item == nil {
		return
	}

	consumeEntity, ok := item.(*entity.Consume)
	if !ok || consumeEntity == nil {
		return
	}

	consumeWz, ok := consumeEntity.GetModel().(*wz.Consume)
	if !ok || consumeWz == nil {
		return
	}
	itemID := consumeWz.ID
	if itemID == 0 {
		return
	}

	mapInstance := ch.GetMap()
	if mapInstance == nil {
		return
	}

	root := mapInstance.GetLuaRoot()
	if root == nil {
		return
	}

	scriptPath := fmt.Sprintf("script/item/%d.lua", itemID)
	thread, err := luax.NewThread(root, scriptPath)
	if err != nil {
		return
	}

	hookName := fmt.Sprintf("on_active_item_%d", itemID)
	_, err = luax.Call(thread, hookName, ch, consumeEntity)
	if err != nil {
		log.Printf("item script failed %s: %v", scriptPath, err)
		return
	}
}
