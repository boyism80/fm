package server

import (
	"fmt"
	"log"

	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/game/client"
	"github.com/boyism80/fm/game/constant"
	"github.com/boyism80/fm/game/wz"
	"github.com/boyism80/fm/protocol/request"
)

type UseItem struct {
	gs     *GameServer
	opcode byte
}

func (UseItem) New(gs *GameServer) *UseItem {
	return &UseItem{
		gs:     gs,
		opcode: 0x37,
	}
}

func (h *UseItem) GetOpcode() byte {
	return h.opcode
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

	if ch.Hp <= 0 {
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

	healMul := ch.PotionHealMultiplierPercent()
	flatHP := int(consumeItem.ActiveEffect.HP) * healMul / 100
	flatMP := int(consumeItem.ActiveEffect.MP) * healMul / 100
	hpRate := int(consumeItem.ActiveEffect.HPRate)
	mpRate := int(consumeItem.ActiveEffect.MPRate)

	hpChange := flatHP
	mpChange := flatMP

	if hpChange == 0 && mpChange == 0 && hpRate == 0 && mpRate == 0 {
		ch.Listener.OnUpdateStats(ch, nil, true)
		return nil
	}

	stats := make(map[constant.Stat]int32)

	maxHp := ch.GetMaxHp()
	maxMp := ch.GetMaxMp()
	if hpRate > 0 {
		hpRecovery := int(maxHp) * hpRate / 100
		hpChange += hpRecovery
	}

	if mpRate > 0 {
		mpRecovery := int(maxMp) * mpRate / 100
		mpChange += mpRecovery
	}

	actualHPChange := 0
	actualMPChange := 0

	if hpChange != 0 {
		newHP := int(ch.Hp) + hpChange
		if newHP < 1 {
			newHP = 1
		}
		if newHP > int(maxHp) {
			newHP = int(maxHp)
		}
		actualHPChange = newHP - int(ch.Hp)
		if actualHPChange != 0 {
			ch.Hp = uint32(newHP)
			stats[constant.STAT_HP] = int32(ch.Hp)
		}
	}

	if mpChange != 0 {
		newMP := int(ch.Mp) + mpChange
		if newMP < 0 {
			newMP = 0
		}
		if newMP > int(maxMp) {
			newMP = int(maxMp)
		}
		actualMPChange = newMP - int(ch.Mp)
		if actualMPChange != 0 {
			ch.Mp = uint32(newMP)
			stats[constant.STAT_MP] = int32(ch.Mp)
		}
	}

	if actualHPChange == 0 && actualMPChange == 0 {
		ch.Listener.OnUpdateStats(ch, nil, true)
		return nil
	}

	if len(stats) > 0 {
		ch.Listener.OnUpdateStats(ch, stats, false)
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
