package game

import (
	"errors"
	"math/rand"
	"slices"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/maple/constant"
	"github.com/boyism80/fm/maple/wz"
)

var ErrCannotScroll = errors.New("cannot scroll equipment")

func (ch *Character) Scroll(actx actor.Context, scrollSlot int16, targetSlot int16) error {
	if scrollSlot <= 0 {
		return ErrCannotScroll
	}
	useInventory := ch.inventory.tabs[constant.InventoryTypeConsume]
	scrollItem, ok := ch.inventory.GetItem(constant.InventoryTypeConsume, scrollSlot).(*Consume)
	if !ok {
		return ErrCannotScroll
	}
	scroll, ok := scrollItem.GetModel().(*wz.Consume)
	if !ok || scroll.GetTraits()&constant.ItemTraitScroll == 0 {
		return ErrCannotScroll
	}

	equipInventory := ch.inventory.tabs[constant.InventoryTypeEquipment]
	var target Equipment
	if targetSlot < 0 {
		target = ch.inventory.equipped[constant.EquipmentPartsType(targetSlot)]
	} else {
		target, _ = ch.inventory.GetItem(constant.InventoryTypeEquipment, targetSlot).(Equipment)
	}
	if target == nil {
		return ErrCannotScroll
	}
	model, ok := target.GetModel().(wz.Equipment)
	if !ok {
		return ErrCannotScroll
	}
	core := target.GetEquipmentCore()

	recovery := scroll.ScrollRecover > 0
	chaos := scroll.ScrollRandStat > 0
	special := scroll.ScrollFlag != 0
	if !recovery && !special && core.enhanceChance < 1 {
		return ErrCannotScroll
	}
	if (recovery || chaos || special) && model.GetTraits()&constant.ItemTraitCash != 0 {
		return ErrCannotScroll
	}
	switch model.GetEquipmentType() {
	case constant.EquipmentTypeTaming, constant.EquipmentTypeDragon, constant.EquipmentTypeMechanic, constant.EquipmentTypeAndroid:
		return ErrCannotScroll
	}
	if len(scroll.ScrollReqs) > 0 && !slices.Contains(scroll.ScrollReqs, model.GetID()) {
		return ErrCannotScroll
	}
	if !recovery && !chaos && !special && scroll.GetID()/100%100 != model.GetID()/10000%100 {
		return ErrCannotScroll
	}
	if scroll.ScrollReqRUC > 0 && int32(model.GetEnhanceChance())-int32(core.enhanceCount) < scroll.ScrollReqRUC {
		return ErrCannotScroll
	}

	if core.bonusStats == nil {
		core.bonusStats = &EquipmentBonusStats{}
	}
	success := special || rand.Intn(100) < int(scroll.ScrollSuccess)
	destroyed := false
	switch {
	case success && recovery:
		if core.enhanceCount+core.enhanceChance < model.GetEnhanceChance() {
			core.enhanceChance += uint8(scroll.ScrollRecover)
		} else {
			success = false
		}
	case success && special:
		core.flags |= uint16(scroll.ScrollFlag)
	case success:
		if chaos {
			core.rollChaosStats()
		} else {
			core.addScrollStats(scroll)
		}
		core.enhanceChance--
		core.enhanceCount++
	default:
		if !recovery {
			core.enhanceChance--
		}
		destroyed = rand.Intn(100) < int(scroll.ScrollCursed)
	}

	if scrollItem.Reduce(1) == 0 {
		useInventory.Remove(uint8(scrollSlot))
	}
	if destroyed {
		if targetSlot < 0 {
			parts := constant.EquipmentPartsType(targetSlot)
			delete(ch.inventory.equipped, parts)
			ch.listener.OnUpdateCharacterLook(ch)
			ch.onEquipmentChanged(actx, parts, target, nil)
		} else {
			equipInventory.Remove(uint8(targetSlot))
		}
	}

	ch.listener.OnScrolledItem(ch, scrollItem.GetInventoryType(), scrollSlot, scrollItem.GetCount(), targetSlot, destroyed, false, target)
	ch.listener.OnShowScrollEffect(ch, success, destroyed)
	ch.listener.OnUpdateStats(ch, nil, true)
	return nil
}

func (c *EquipmentCore) addScrollStats(scroll *wz.Consume) {
	c.bonusStats.Str += scroll.ScrollIncStr
	c.bonusStats.Dex += scroll.ScrollIncDex
	c.bonusStats.Int += scroll.ScrollIncInt
	c.bonusStats.Luk += scroll.ScrollIncLuk
	c.bonusStats.MaxHP += scroll.ScrollIncMaxHP
	c.bonusStats.MaxMP += scroll.ScrollIncMaxMP
	c.bonusStats.PAD += scroll.ScrollIncPAD
	c.bonusStats.MAD += scroll.ScrollIncMAD
	c.bonusStats.PDD += scroll.ScrollIncPDD
	c.bonusStats.MDD += scroll.ScrollIncMDD
	c.bonusStats.ACC += scroll.ScrollIncACC
	c.bonusStats.Avoid += scroll.ScrollIncAvoid
	c.bonusStats.Hands += scroll.ScrollIncHands
	c.bonusStats.Speed += scroll.ScrollIncSpeed
	c.bonusStats.Jump += scroll.ScrollIncJump
}

func (c *EquipmentCore) rollChaosStats() {
	ability := c.wz.(wz.Equipment).GetAbility()
	stats := []struct {
		base  uint16
		bonus *int16
	}{
		{ability.Str, &c.bonusStats.Str},
		{ability.Dex, &c.bonusStats.Dex},
		{ability.Int, &c.bonusStats.Int},
		{ability.Luk, &c.bonusStats.Luk},
		{ability.PAD, &c.bonusStats.PAD},
		{ability.PDD, &c.bonusStats.PDD},
		{ability.MAD, &c.bonusStats.MAD},
		{ability.MDD, &c.bonusStats.MDD},
		{ability.ACC, &c.bonusStats.ACC},
		{ability.Avoid, &c.bonusStats.Avoid},
		{ability.Speed, &c.bonusStats.Speed},
		{ability.Jump, &c.bonusStats.Jump},
		{ability.MaxHP, &c.bonusStats.MaxHP},
		{ability.MaxMP, &c.bonusStats.MaxMP},
	}
	for _, stat := range stats {
		if int32(stat.base)+int32(*stat.bonus) <= 0 {
			continue
		}
		delta := int16(rand.Intn(constant.ChaosScrollRange))
		if rand.Intn(2) == 0 {
			delta = -delta
		}
		*stat.bonus += delta
	}
}
