package server

import (
	"fmt"
	"log"
	"math/rand"

	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/protocol/request"
	"github.com/boyism80/fm/services/game/client"
	"github.com/boyism80/fm/services/game/constant"
	"github.com/boyism80/fm/services/game/entity"
	"github.com/boyism80/fm/services/game/wz"
	"github.com/boyism80/fm/util"
)

type EnhanceEquipment struct{}

func (EnhanceEquipment) New(_ *GameServer) *EnhanceEquipment {
	return &EnhanceEquipment{}
}

func (*EnhanceEquipment) canScroll(scrollID uint32, itemID uint32) bool {
	return (scrollID/100)%100 == (itemID/10000)%100
}

func (*EnhanceEquipment) addBonusWithClamp(base int16, add int16) int16 {
	return int16(util.ClampInt32(int32(base)+int32(add), -32768, 32767))
}

func (h *EnhanceEquipment) applyScrollBonusStats(targetCore *entity.EquipmentCore, scroll *wz.Consume) {
	if targetCore == nil || scroll == nil {
		return
	}
	if targetCore.BonusStats == nil {
		targetCore.BonusStats = &entity.EquipmentBonusStats{}
	}
	b := targetCore.BonusStats
	b.Str = h.addBonusWithClamp(b.Str, scroll.ScrollIncStr)
	b.Dex = h.addBonusWithClamp(b.Dex, scroll.ScrollIncDex)
	b.Int = h.addBonusWithClamp(b.Int, scroll.ScrollIncInt)
	b.Luk = h.addBonusWithClamp(b.Luk, scroll.ScrollIncLuk)
	b.MaxHP = h.addBonusWithClamp(b.MaxHP, scroll.ScrollIncMaxHP)
	b.MaxMP = h.addBonusWithClamp(b.MaxMP, scroll.ScrollIncMaxMP)
	b.PAD = h.addBonusWithClamp(b.PAD, scroll.ScrollIncPAD)
	b.MAD = h.addBonusWithClamp(b.MAD, scroll.ScrollIncMAD)
	b.PDD = h.addBonusWithClamp(b.PDD, scroll.ScrollIncPDD)
	b.MDD = h.addBonusWithClamp(b.MDD, scroll.ScrollIncMDD)
	b.ACC = h.addBonusWithClamp(b.ACC, scroll.ScrollIncACC)
	b.Avoid = h.addBonusWithClamp(b.Avoid, scroll.ScrollIncAvoid)
	b.Hands = h.addBonusWithClamp(b.Hands, scroll.ScrollIncHands)
	b.Speed = h.addBonusWithClamp(b.Speed, scroll.ScrollIncSpeed)
	b.Jump = h.addBonusWithClamp(b.Jump, scroll.ScrollIncJump)
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

	useInventory := ch.Inventory[constant.INVENTORY_TYPE_CONSUME]
	if req.ScrollSlot <= 0 {
		ch.Listener.OnUpdateStats(ch, nil, true)
		return nil
	}
	scrollItem := useInventory.GetItem(uint8(req.ScrollSlot))
	if scrollItem == nil || scrollItem.GetCount() < 1 {
		ch.Listener.OnUpdateStats(ch, nil, true)
		return nil
	}
	scrollConsume, ok := scrollItem.(*entity.Consume)
	if !ok || scrollConsume == nil {
		ch.Listener.OnUpdateStats(ch, nil, true)
		return nil
	}
	scrollWz, ok := scrollConsume.GetModel().(*wz.Consume)
	if !ok || scrollWz == nil || scrollWz.ScrollSuccess <= 0 {
		ch.Listener.OnUpdateStats(ch, nil, true)
		return nil
	}

	var targetEquip entity.Equipment
	equipInventory := ch.Inventory[constant.INVENTORY_TYPE_EQUIPMENT]
	if req.TargetSlot < 0 {
		targetEquip = ch.Equipments[constant.EquipmentPartsType(req.TargetSlot)]
	} else {
		if item := equipInventory.Items[req.TargetSlot]; item != nil {
			targetEquip, _ = item.(entity.Equipment)
		}
	}
	if targetEquip == nil {
		ch.Listener.OnUpdateStats(ch, nil, true)
		return nil
	}
	targetCore := targetEquip.GetEquipmentCore()
	if targetCore == nil || targetCore.EnhanceChance == 0 {
		ch.Listener.OnUpdateStats(ch, nil, true)
		return nil
	}
	if !h.canScroll(scrollWz.GetID(), targetEquip.GetModel().GetID()) {
		ch.Listener.OnUpdateStats(ch, nil, true)
		return nil
	}

	scrollConsume.Reduce(1)
	if scrollConsume.GetCount() == 0 {
		if err := useInventory.RemoveItem(uint8(req.ScrollSlot)); err != nil {
			ch.Listener.OnUpdateStats(ch, nil, true)
			return nil
		}
	}
	targetCore.EnhanceChance--

	success := rand.Intn(100) < int(min(max(scrollWz.ScrollSuccess, 0), 100))
	destroyed := false
	if !success && scrollWz.ScrollCursed > 0 {
		destroyed = rand.Intn(100) < int(min(max(scrollWz.ScrollCursed, 0), 100))
	}
	if success {
		h.applyScrollBonusStats(targetCore, scrollWz)
		targetCore.EnhanceCount++
	}
	if destroyed {
		if req.TargetSlot < 0 {
			delete(ch.Equipments, constant.EquipmentPartsType(req.TargetSlot))
			ch.Listener.OnUpdateCharacterLook(ch)
		} else {
			if err := equipInventory.RemoveItem(uint8(req.TargetSlot)); err != nil {
				ch.Listener.OnUpdateStats(ch, nil, true)
				return nil
			}
		}
	}
	ch.Listener.OnScrolledItem(
		ch,
		constant.GetInventoryTypeByItemID(scrollWz.GetID()),
		req.ScrollSlot,
		scrollConsume.GetCount(),
		req.TargetSlot,
		destroyed,
		false,
		targetEquip,
	)

	ch.Listener.OnShowScrollEffect(ch, success, destroyed)
	ch.Listener.OnUpdateStats(ch, nil, true)
	return nil
}
