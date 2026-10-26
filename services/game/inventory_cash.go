package game

import (
	"github.com/boyism80/fm/core/clock"
	"github.com/boyism80/fm/maple/constant"
	"github.com/boyism80/fm/util"
)

// The magic hourglass may not push an equipment's expiration further than this from now.
const hourglassLimitDays = 100

// TagEquipment writes the owner's name on worn equipment that has none, as the item tag does.
func (inv *Inventory) TagEquipment(slot int16) bool {
	equipment := inv.equipped[constant.EquipmentPartsType(slot)]
	if equipment == nil || equipment.GetEquipmentCore().ownerName != "" {
		return false
	}
	equipment.GetEquipmentCore().ownerName = inv.owner.GetName()
	inv.owner.listener.OnInventorySlotUpdated(inv.owner, constant.InventoryTypeEquipment, slot, equipment)
	return true
}

// ExtendEquipment adds days to worn equipment that expires, as the magic hourglass does.
func (inv *Inventory) ExtendEquipment(slot int16, days int) bool {
	equipment := inv.equipped[constant.EquipmentPartsType(slot)]
	if equipment == nil || equipment.GetModel().GetTraits()&constant.ItemTraitCash != 0 {
		return false
	}
	expiration := equipment.GetExpiration()
	if !expiration.Before(util.TimeMax) {
		return false
	}
	extended := expiration.AddDate(0, 0, days)
	if extended.After(clock.Now().AddDate(0, 0, hourglassLimitDays)) {
		return false
	}
	equipment.SetExpiration(extended)
	inv.owner.listener.OnInventorySlotUpdated(inv.owner, constant.InventoryTypeEquipment, slot, equipment)
	return true
}

// LockItem seals an item that never expires so it cannot be dropped or traded, as the sealing lock does.
func (inv *Inventory) LockItem(invType constant.InventoryType, slot int16) bool {
	item := inv.GetItem(invType, slot)
	if item == nil || item.GetExpiration().Before(util.TimeMax) || item.GetFlags()&constant.ItemFlagLock != 0 {
		return false
	}
	item.SetFlags(item.GetFlags() | constant.ItemFlagLock)
	inv.owner.listener.OnInventorySlotUpdated(inv.owner, invType, slot, item)
	return true
}

// UseKarma lets an untradeable item be traded once, as the scissors of karma do; WZ names which scissors an item takes.
func (inv *Inventory) UseKarma(invType constant.InventoryType, slot int16, scissors constant.ItemTrait) bool {
	item := inv.GetItem(invType, slot)
	if item == nil || item.GetModel().GetTraits()&scissors == 0 || item.Tradable() {
		return false
	}
	flags := item.GetFlags()
	if flags&constant.ItemFlagUntradeable != 0 {
		flags &^= constant.ItemFlagUntradeable
	} else {
		flags |= item.karma()
	}
	item.SetFlags(flags)
	inv.owner.listener.OnInventorySlotUpdated(inv.owner, invType, slot, item)
	return true
}
