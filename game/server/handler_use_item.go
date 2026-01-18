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

	character := client.GetCharacter()
	if character == nil {
		return nil
	}

	if character.Hp <= 0 {
		if character.Listener != nil {
			character.Listener.OnUpdateStats(nil, true)
		}
		return nil
	}

	useInventory := character.Inventory[constant.INVENTORY_TYPE_CONSUME]
	if useInventory == nil {
		if character.Listener != nil {
			character.Listener.OnUpdateStats(nil, true)
		}
		return nil
	}

	item := useInventory.GetItem(uint8(req.Slot))
	if item == nil {
		if character.Listener != nil {
			character.Listener.OnUpdateStats(nil, true)
		}
		return nil
	}

	if item.GetCount() < 1 {
		if character.Listener != nil {
			character.Listener.OnUpdateStats(nil, true)
		}
		return nil
	}

	if item.GetModel().GetID() != req.ItemID {
		if character.Listener != nil {
			character.Listener.OnUpdateStats(nil, true)
		}
		return nil
	}

	resources := h.gs.GetResources()
	if resources == nil {
		if character.Listener != nil {
			character.Listener.OnUpdateStats(nil, true)
		}
		return nil
	}

	wzItem := resources.Items[req.ItemID]
	if wzItem == nil {
		if character.Listener != nil {
			character.Listener.OnUpdateStats(nil, true)
		}
		return nil
	}

	consumeItem, ok := wzItem.(*wz.Consume)
	if !ok {
		if character.Listener != nil {
			character.Listener.OnUpdateStats(nil, true)
		}
		return nil
	}

	hpChange := int(consumeItem.ActiveEffect.HP)
	mpChange := int(consumeItem.ActiveEffect.MP)
	hpRate := int(consumeItem.ActiveEffect.HPRate)
	mpRate := int(consumeItem.ActiveEffect.MPRate)

	if hpChange == 0 && mpChange == 0 && hpRate == 0 && mpRate == 0 {
		if character.Listener != nil {
			character.Listener.OnUpdateStats(nil, true)
		}
		return nil
	}

	stats := make(map[constant.Stat]int32)

	if hpRate > 0 {
		hpRecovery := int(character.MaxHp) * hpRate / 100
		hpChange += hpRecovery
	}

	if mpRate > 0 {
		mpRecovery := int(character.MaxMp) * mpRate / 100
		mpChange += mpRecovery
	}

	actualHPChange := 0
	actualMPChange := 0

	if hpChange != 0 {
		newHP := int(character.Hp) + hpChange
		if newHP < 1 {
			newHP = 1
		}
		if newHP > int(character.MaxHp) {
			newHP = int(character.MaxHp)
		}
		actualHPChange = newHP - int(character.Hp)
		if actualHPChange != 0 {
			character.Hp = uint16(newHP)
			stats[constant.STAT_HP] = int32(character.Hp)
		}
	}

	if mpChange != 0 {
		newMP := int(character.Mp) + mpChange
		if newMP < 0 {
			newMP = 0
		}
		if newMP > int(character.MaxMp) {
			newMP = int(character.MaxMp)
		}
		actualMPChange = newMP - int(character.Mp)
		if actualMPChange != 0 {
			character.Mp = uint16(newMP)
			stats[constant.STAT_MP] = int32(character.Mp)
		}
	}

	if actualHPChange == 0 && actualMPChange == 0 {
		if character.Listener != nil {
			character.Listener.OnUpdateStats(nil, true)
		}
		return nil
	}

	if len(stats) > 0 && character.Listener != nil {
		character.Listener.OnUpdateStats(stats, false)
	}

	item.Reduce(1)
	if character.Listener != nil {
		if item.GetCount() == 0 {
			useInventory.RemoveItem(uint8(req.Slot))
			character.Listener.OnRemoveInventorySlot(constant.INVENTORY_TYPE_CONSUME, int16(req.Slot))
		} else {
			character.Listener.OnInventorySlotUpdated(constant.INVENTORY_TYPE_CONSUME, int16(req.Slot), item)
		}
		character.Listener.OnUpdateStats(nil, true)
	} else {
		if item.GetCount() == 0 {
			useInventory.RemoveItem(uint8(req.Slot))
		}
	}

	return nil
}
