package entity

import (
	"errors"
	"math/rand"
	"slices"

	"github.com/boyism80/fm/services/game/constant"
	"github.com/boyism80/fm/services/game/wz"
)

var ErrCannotScroll = errors.New("cannot scroll equipment")

func (ch *Character) Scroll(scrollSlot int16, targetSlot int16) error {
	if scrollSlot <= 0 {
		return ErrCannotScroll
	}
	useInventory := ch.Inventory.Tabs[constant.InventoryTypeConsume]
	scrollItem, ok := useInventory.Get(uint8(scrollSlot)).(*Consume)
	if ok == false {
		return ErrCannotScroll
	}
	scroll, ok := scrollItem.GetModel().(*wz.Consume)
	if ok == false {
		return ErrCannotScroll
	}

	equipInventory := ch.Inventory.Tabs[constant.InventoryTypeEquipment]
	var target Equipment
	if targetSlot < 0 {
		target = ch.Inventory.Equipped[constant.EquipmentPartsType(targetSlot)]
	} else {
		target, _ = equipInventory.Get(uint8(targetSlot)).(Equipment)
	}
	if target == nil {
		return ErrCannotScroll
	}
	model, ok := target.GetModel().(wz.Equipment)
	if ok == false {
		return ErrCannotScroll
	}
	core := target.GetEquipmentCore()

	recovery := scroll.ScrollRecover > 0
	chaos := scroll.ScrollRandStat > 0
	special := scroll.ScrollFlag != 0
	if recovery == false && special == false && core.EnhanceChance < 1 {
		return ErrCannotScroll
	}
	if (recovery || chaos || special) && model.IsCash() {
		return ErrCannotScroll
	}
	switch constant.GetEquipmentType(model.GetID()) {
	case constant.EquipmentTypeTaming, constant.EquipmentTypeDragon, constant.EquipmentTypeMechanic, constant.EquipmentTypeAndroid:
		return ErrCannotScroll
	}
	if len(scroll.ScrollReqs) > 0 && slices.Contains(scroll.ScrollReqs, model.GetID()) == false {
		return ErrCannotScroll
	}
	if recovery == false && chaos == false && special == false && scroll.GetID()/100%100 != model.GetID()/10000%100 {
		return ErrCannotScroll
	}
	if scroll.ScrollReqRUC > 0 && int32(model.GetEnhanceChance())-int32(core.EnhanceCount) < scroll.ScrollReqRUC {
		return ErrCannotScroll
	}

	if core.BonusStats == nil {
		core.BonusStats = &EquipmentBonusStats{}
	}
	success := special || rand.Intn(100) < int(scroll.ScrollSuccess)
	destroyed := false
	switch {
	case success && recovery:
		if core.EnhanceCount+core.EnhanceChance < model.GetEnhanceChance() {
			core.EnhanceChance += uint8(scroll.ScrollRecover)
		} else {
			success = false
		}
	case success && special:
		core.Flag |= uint16(scroll.ScrollFlag)
	case success:
		if chaos {
			core.rollChaosStats()
		} else {
			core.addScrollStats(scroll)
		}
		core.EnhanceChance--
		core.EnhanceCount++
	default:
		if recovery == false {
			core.EnhanceChance--
		}
		destroyed = rand.Intn(100) < int(scroll.ScrollCursed)
	}

	if scrollItem.Reduce(1) == 0 {
		useInventory.Remove(uint8(scrollSlot))
	}
	if destroyed {
		if targetSlot < 0 {
			delete(ch.Inventory.Equipped, constant.EquipmentPartsType(targetSlot))
			ch.Listener.OnUpdateCharacterLook(ch)
		} else {
			equipInventory.Remove(uint8(targetSlot))
		}
	}

	ch.Listener.OnScrolledItem(ch, scrollItem.GetInventoryType(), scrollSlot, scrollItem.GetCount(), targetSlot, destroyed, false, target)
	ch.Listener.OnShowScrollEffect(ch, success, destroyed)
	ch.Listener.OnUpdateStats(ch, nil, true)
	return nil
}

func (c *EquipmentCore) addScrollStats(scroll *wz.Consume) {
	c.BonusStats.Str += scroll.ScrollIncStr
	c.BonusStats.Dex += scroll.ScrollIncDex
	c.BonusStats.Int += scroll.ScrollIncInt
	c.BonusStats.Luk += scroll.ScrollIncLuk
	c.BonusStats.MaxHP += scroll.ScrollIncMaxHP
	c.BonusStats.MaxMP += scroll.ScrollIncMaxMP
	c.BonusStats.PAD += scroll.ScrollIncPAD
	c.BonusStats.MAD += scroll.ScrollIncMAD
	c.BonusStats.PDD += scroll.ScrollIncPDD
	c.BonusStats.MDD += scroll.ScrollIncMDD
	c.BonusStats.ACC += scroll.ScrollIncACC
	c.BonusStats.Avoid += scroll.ScrollIncAvoid
	c.BonusStats.Hands += scroll.ScrollIncHands
	c.BonusStats.Speed += scroll.ScrollIncSpeed
	c.BonusStats.Jump += scroll.ScrollIncJump
}

func (c *EquipmentCore) rollChaosStats() {
	ability := c.Wz.(wz.Equipment).GetAbility()
	stats := []struct {
		base  uint16
		bonus *int16
	}{
		{ability.Str, &c.BonusStats.Str},
		{ability.Dex, &c.BonusStats.Dex},
		{ability.Int, &c.BonusStats.Int},
		{ability.Luk, &c.BonusStats.Luk},
		{ability.PAD, &c.BonusStats.PAD},
		{ability.PDD, &c.BonusStats.PDD},
		{ability.MAD, &c.BonusStats.MAD},
		{ability.MDD, &c.BonusStats.MDD},
		{ability.ACC, &c.BonusStats.ACC},
		{ability.Avoid, &c.BonusStats.Avoid},
		{ability.Speed, &c.BonusStats.Speed},
		{ability.Jump, &c.BonusStats.Jump},
		{ability.MaxHP, &c.BonusStats.MaxHP},
		{ability.MaxMP, &c.BonusStats.MaxMP},
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
