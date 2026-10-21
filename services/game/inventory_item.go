package game

import (
	"errors"
	"fmt"
	"log"
	"math"
	"slices"
	"time"

	"github.com/boyism80/fm/core/luax"
	"github.com/boyism80/fm/maple/constant"
	"github.com/boyism80/fm/maple/wz"
	"github.com/boyism80/fm/protocol/response"
)

func (inv *Inventory) AddMeso(amount int32) {
	if amount <= 0 {
		return
	}
	if inv.validateMeso(0, amount) != ExchangeOK {
		return
	}
	inv.addMesoUnchecked(amount)
}

func (inv *Inventory) addMesoUnchecked(amount int32) {
	ch := inv.owner
	if amount <= 0 {
		return
	}

	if inv.meso > 0 && amount > 0 && inv.meso+amount < inv.meso {
		inv.meso = int32(^uint32(0) >> 1)
	} else {
		inv.meso += amount
	}

	ch.listener.OnMesoChanged(ch, inv.meso)
}

func (inv *Inventory) RemoveMeso(amount int32) bool {
	if amount <= 0 {
		return false
	}
	if inv.validateMeso(amount, 0) != ExchangeOK {
		return false
	}
	inv.removeMesoUnchecked(amount)
	return true
}

func (inv *Inventory) removeMesoUnchecked(amount int32) {
	ch := inv.owner
	if amount <= 0 {
		return
	}

	if inv.meso < amount {
		inv.meso = 0
	} else {
		inv.meso -= amount
	}

	ch.listener.OnMesoChanged(ch, inv.meso)
}

func (inv *Inventory) GetItem(invType constant.InventoryType, slot int16) Item {
	inven := inv.tabs[invType]
	if inven == nil {
		return nil
	}
	if slot < 1 || slot > int16(inven.slotLimit) {
		return nil
	}
	return inven.Get(uint8(slot))
}

func (inv *Inventory) RemoveItem(invType constant.InventoryType, slot int16, count uint16) bool {
	ch := inv.owner
	inven := inv.tabs[invType]
	if inven == nil {
		return false
	}
	if slot < 1 || slot > int16(inven.slotLimit) {
		return false
	}
	item := inven.Get(uint8(slot))
	if item == nil {
		return false
	}
	if item.GetCount() < count {
		return false
	}
	if item.GetCount() == count {
		itemID := item.GetModel().GetID()
		if err := inven.Remove(uint8(slot)); err != nil {
			return false
		}
		ch.listener.OnRemoveInventorySlot(ch, invType, slot)
		if invType == constant.InventoryTypeCash && ch.itemEffectID == itemID && inven.Find(itemID) == nil {
			ch.clearItemEffect()
		}
	} else {
		item.Reduce(count)
		ch.listener.OnInventorySlotUpdated(ch, invType, slot, item)
	}
	return true
}

// useAmmo keeps a used-up rechargeable in its slot so it can be recharged at a shop.
func (inv *Inventory) useAmmo(slot int16, count uint16) bool {
	item := inv.GetItem(constant.InventoryTypeConsume, slot)
	if item == nil || item.GetCount() < count {
		return false
	}
	if item.GetModel().GetTraits()&constant.ItemTraitRechargeable == 0 {
		return inv.RemoveItem(constant.InventoryTypeConsume, slot, count)
	}

	item.Reduce(count)
	inv.owner.listener.OnInventorySlotUpdated(inv.owner, constant.InventoryTypeConsume, slot, item)
	return true
}

func (inv *Inventory) FindSlots(itemID uint32) (invType constant.InventoryType, slots []int16) {
	model := inv.owner.itemModel(itemID)
	if model == nil {
		return 0, nil
	}
	invType = model.GetInventoryType()
	inven := inv.tabs[invType]
	if inven == nil {
		return invType, nil
	}
	for slot := int16(1); slot <= int16(inven.slotLimit); slot++ {
		item := inven.Get(uint8(slot))
		if item != nil && item.GetModel().GetID() == itemID {
			slots = append(slots, slot)
		}
	}
	return invType, slots
}

func (inv *Inventory) GetCountByItemID(itemID uint32) uint16 {
	invType, slots := inv.FindSlots(itemID)
	var total uint16
	for _, slot := range slots {
		if item := inv.GetItem(invType, slot); item != nil {
			total += item.GetCount()
		}
	}
	return total
}

func (inv *Inventory) HasItem(itemID uint32) bool {
	model := inv.owner.itemModel(itemID)
	if model == nil {
		return false
	}
	inven := inv.tabs[model.GetInventoryType()]
	if inven == nil {
		return false
	}
	for slot := int16(1); slot <= int16(inven.slotLimit); slot++ {
		item := inven.Get(uint8(slot))
		if item != nil && item.GetModel().GetID() == itemID {
			return true
		}
	}
	return false
}

// hasUniqueItem also counts what the character put in a trade or shop, which returns to the inventory if the room is cancelled.
func (inv *Inventory) hasUniqueItem(model wz.Item) bool {
	if inv.carriesUniqueItem(model) {
		return true
	}
	if model == nil || model.GetTraits()&constant.ItemTraitUnique == 0 {
		return false
	}
	for _, held := range inv.owner.escrow.items {
		if held.item.GetModel().GetID() == model.GetID() {
			return true
		}
	}
	return false
}

// carriesUniqueItem looks at the inventory and the worn equipment only.
func (inv *Inventory) carriesUniqueItem(model wz.Item) bool {
	if model == nil || model.GetTraits()&constant.ItemTraitUnique == 0 {
		return false
	}
	if inv.HasItem(model.GetID()) {
		return true
	}
	if model.GetInventoryType() != constant.InventoryTypeEquipment {
		return false
	}
	for _, eq := range inv.equipped {
		if eq != nil && eq.GetModel().GetID() == model.GetID() {
			return true
		}
	}
	return false
}

func (inv *Inventory) HasItemCount(itemID uint32, count uint16) bool {
	if count == 0 {
		return true
	}
	model := inv.owner.itemModel(itemID)
	if model == nil {
		return false
	}
	inven := inv.tabs[model.GetInventoryType()]
	if inven == nil {
		return false
	}
	var total uint16
	for slot := int16(1); slot <= int16(inven.slotLimit); slot++ {
		item := inven.Get(uint8(slot))
		if item == nil || item.GetModel().GetID() != itemID {
			continue
		}
		total += item.GetCount()
		if total >= count {
			return true
		}
	}
	return false
}

func (inv *Inventory) RemoveByItemIDCount(itemID uint32, count uint16) bool {
	ch := inv.owner
	if count == 0 {
		return true
	}
	spec := ExchangeSpec{
		Cost: ExchangeSide{
			Items: map[uint32]uint16{itemID: count},
		},
	}
	if spec.Valid(ch, ch.escrow) != ExchangeOK {
		return false
	}
	return inv.removeByItemIDCountUnchecked(itemID, count)
}

func (inv *Inventory) removeByItemIDCountUnchecked(itemID uint32, count uint16) bool {
	if count == 0 {
		return true
	}
	invType, slots := inv.FindSlots(itemID)
	remaining := count
	for _, slot := range slots {
		if remaining == 0 {
			break
		}
		item := inv.GetItem(invType, slot)
		if item == nil {
			continue
		}
		take := item.GetCount()
		if take > remaining {
			take = remaining
		}
		remaining -= take
		inv.RemoveItem(invType, slot, take)
	}
	return remaining == 0
}

func (inv *Inventory) ClearInventory() int {
	if inv == nil || inv.owner == nil {
		return 0
	}
	ch := inv.owner
	cleared := 0
	for invType, inven := range inv.tabs {
		if inven == nil || inven.items == nil {
			continue
		}
		slots := make([]int16, 0, len(inven.items))
		for slot, item := range inven.items {
			if item == nil {
				continue
			}
			slots = append(slots, slot)
		}
		for _, slot := range slots {
			delete(inven.items, slot)
			ch.listener.OnRemoveInventorySlot(ch, invType, slot)
			cleared++
		}
	}
	ch.clearItemEffect()
	return cleared
}

func (inv *Inventory) ExpireItems(now time.Time) {
	var expired []wz.Item
	for parts, equipment := range inv.equipped {
		if equipment == nil || !now.After(equipment.GetExpiration()) {
			continue
		}
		inv.RemoveEquipped(parts)
		expired = append(expired, equipment.GetModel())
	}
	for invType, inven := range inv.tabs {
		for slot, item := range inven.items {
			if item == nil || !now.After(item.GetExpiration()) {
				continue
			}
			if _, ok := item.(*Pet); ok {
				continue
			}
			inv.RemoveItem(invType, slot, item.GetCount())
			expired = append(expired, item.GetModel())
		}
	}
	if len(expired) == 0 {
		return
	}
	inv.owner.listener.OnItemExpired(inv.owner, expired)
}

func (inv *Inventory) AddItem(item Item, allOrNothing bool) (addedItems []Item, err error) {
	ch := inv.owner
	if item == nil {
		return nil, fmt.Errorf("item is nil")
	}

	if allOrNothing {
		spec := ExchangeSpec{
			Reward: ExchangeSide{
				Items: map[uint32]uint16{item.GetModel().GetID(): item.GetCount()},
			},
		}
		if spec.Valid(ch, ch.escrow) != ExchangeOK {
			return nil, fmt.Errorf("not enough inventory space for %d items", item.GetCount())
		}
	}

	return inv.addItemUnchecked(item, allOrNothing)
}

func (inv *Inventory) addItemUnchecked(item Item, allOrNothing bool) (addedItems []Item, err error) {
	ch := inv.owner
	if item == nil {
		return nil, fmt.Errorf("item is nil")
	}

	invenType := item.GetInventoryType()
	inven := inv.tabs[invenType]
	if inven == nil {
		return nil, fmt.Errorf("inventory type %d not found", invenType)
	}

	model := item.GetModel()
	requestedCount := item.GetCount()

	remainingCount := requestedCount
	var addedCount uint16

	for remainingCount > 0 {
		slot, ok := inv.owner.escrow.gainSlot(invenType, item)
		if !ok {
			if allOrNothing && addedCount == 0 {
				return nil, fmt.Errorf("no available slot found")
			}
			break
		}

		exists, ok := inven.items[int16(slot)]
		cap := uint16(0)
		var slotItem Item
		if ok {
			cap = min(model.GetCapacity()-exists.GetCount(), remainingCount)
			exists.Increase(cap)
			slotItem = exists
			ch.listener.OnInventorySlotUpdated(ch, invenType, int16(slot), exists)
		} else {
			cap = min(model.GetCapacity(), remainingCount)
			placed := item.Clone(cap)
			inven.items[int16(slot)] = placed
			slotItem = placed
			ch.listener.OnInventorySlotAdded(ch, invenType, int16(slot), placed)
		}
		addedItems = append(addedItems, slotItem)
		remainingCount -= cap
		addedCount += cap
	}

	if !allOrNothing && addedCount < requestedCount {
		item.SetCount(remainingCount)
	}

	if addedCount > 0 {
		inv.runGainHook(model.GetID())
		ch.listener.OnShowItemGain(ch, model.GetID(), uint32(addedCount), constant.ShowItemGainTypeStatus)
	}

	return addedItems, nil
}

func (inv *Inventory) runGainHook(itemID uint32) {
	if itemID == 0 {
		return
	}
	mapInstance := inv.owner.GetMap()
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
	luax.SetConfiguration(thread, luax.Configuration{
		ActorPID: mapInstance.LogicActorPID(),
	})
	hook := "on_item_gain"
	luax.CallAsync(nil, root, thread, hook, inv.owner, itemID).OnError(func(err error) {
		log.Printf("item gain script %s: %v", hook, err)
	})
}

func (inv *Inventory) GainMeso(amount int32) {
	if amount <= 0 {
		return
	}
	if inv.validateMeso(0, amount) != ExchangeOK {
		return
	}
	inv.gainMesoUnchecked(amount)
}

func (inv *Inventory) gainMesoUnchecked(amount int32) {
	ch := inv.owner
	if amount <= 0 {
		return
	}

	if inv.meso > 0 && amount > 0 && inv.meso+amount < inv.meso {
		inv.meso = int32(^uint32(0) >> 1)
	} else {
		inv.meso += amount
	}

	ch.listener.OnMesoChanged(ch, inv.meso)
	ch.listener.OnShowMesoGain(ch, amount, constant.ShowMesoGainTypeStatus)
	ch.listener.OnUpdateStats(ch, map[constant.Stat]int32{
		constant.StatMeso: inv.meso,
	}, false)
}

func (inv *Inventory) FindSlot(invType constant.InventoryType, item Item) (int16, bool) {
	inven := inv.tabs[invType]
	if inven == nil || item == nil {
		return 0, false
	}
	for slot := int16(1); slot <= int16(inven.slotLimit); slot++ {
		if inven.items[slot] == item {
			return slot, true
		}
	}
	return 0, false
}

func (inv *Inventory) UnequipToSlot(parts constant.EquipmentPartsType, destSlot int16) error {
	ch := inv.owner
	equipments := inv.equipped
	inventory := inv.tabs
	if equipments[parts] == nil {
		return nil
	}
	inven := inventory[constant.InventoryTypeEquipment]
	if inven == nil {
		return errors.New("equipment inventory not found")
	}
	if destSlot < 1 || destSlot > int16(inven.slotLimit) {
		return ErrSlotNotFound
	}
	if inven.items[destSlot] != nil {
		return ErrSlotAlreadyOccupied
	}
	inven.items[destSlot] = equipments[parts]
	delete(equipments, parts)
	ch.listener.OnSwapInventorySlot(ch, constant.InventoryTypeEquipment, int16(parts), destSlot, int8(response.EQUIPMENT_ACTION_TYPE_OFF))
	ch.listener.OnUpdateCharacterLook(ch)
	ch.mount.changeEquipment(parts)
	return nil
}

func (inv *Inventory) RemoveEquipped(parts constant.EquipmentPartsType) Equipment {
	ch := inv.owner
	equipment := inv.equipped[parts]
	if equipment == nil {
		return nil
	}
	delete(inv.equipped, parts)
	ch.listener.OnRemoveInventorySlot(ch, constant.InventoryTypeEquipment, int16(parts))
	ch.listener.OnUpdateCharacterLook(ch)
	ch.mount.changeEquipment(parts)
	return equipment
}

func (inv *Inventory) Unequip(parts constant.EquipmentPartsType) error {
	inven := inv.tabs[constant.InventoryTypeEquipment]
	if inven == nil {
		return errors.New("equipment inventory not found")
	}
	destSlot, ok := inv.owner.escrow.gainSlot(constant.InventoryTypeEquipment, inv.equipped[parts])
	if !ok {
		return ErrInventoryFull
	}
	return inv.UnequipToSlot(parts, int16(destSlot))
}

func (inv *Inventory) Equip(slot int16, parts constant.EquipmentPartsType) error {
	ch := inv.owner
	equipments := inv.equipped
	inventory := inv.tabs
	inven := inventory[constant.InventoryTypeEquipment]
	if inven == nil {
		return errors.New("equipment inventory not found")
	}
	item := inven.items[slot]
	if item == nil {
		return ErrSourceSlotEmpty
	}
	newEq, ok := item.(Equipment)
	if !ok {
		return ErrItemNotEquipment
	}
	model := newEq.GetModel().(wz.Equipment)
	if !model.CanEquipAt(parts) {
		return ErrInvalidEquipmentPart
	}
	if !ch.canEquip(model.GetRequired()) {
		return ErrEquipmentRequirements
	}
	if armor, ok := newEq.GetModel().(*wz.Armor); ok && armor.EquipmentType == constant.EquipmentTypePetEquip {
		if ch.pets.active == nil || !slices.Contains(armor.Pets, ch.pets.active.item.GetModel().GetID()) {
			return ErrInvalidEquipmentPart
		}
	}

	old, swap := equipments[parts]

	switch parts {
	case constant.EquipmentPartsTop:
		if newEq.GetModel().(wz.Equipment).GetEquipmentType() == constant.EquipmentTypeLongcoat {
			if equipments[constant.EquipmentPartsPants] != nil {
				storageSlot, isFree := inv.owner.escrow.gainSlot(constant.InventoryTypeEquipment, equipments[constant.EquipmentPartsPants])
				if !isFree {
					return ErrInventoryFull
				}
				if err := inv.UnequipToSlot(constant.EquipmentPartsPants, int16(storageSlot)); err != nil {
					return err
				}
			}
		}
	case constant.EquipmentPartsPants:
		topEq := equipments[constant.EquipmentPartsTop]
		if topEq != nil {
			if topEq.GetModel().(wz.Equipment).GetEquipmentType() == constant.EquipmentTypeLongcoat {
				storageSlot, isFree := inv.owner.escrow.gainSlot(constant.InventoryTypeEquipment, topEq)
				if swap && !isFree {
					return ErrInventoryFull
				}
				if err := inv.UnequipToSlot(constant.EquipmentPartsTop, int16(storageSlot)); err != nil {
					return err
				}
			}
		}
	}

	if newEq.GetEquipmentCore().bindToWearer() {
		ch.listener.OnRemoveInventorySlot(ch, constant.InventoryTypeEquipment, slot)
		ch.listener.OnInventorySlotAdded(ch, constant.InventoryTypeEquipment, slot, newEq)
	}

	equipments[parts] = newEq
	if old == nil {
		delete(inven.items, slot)
	} else {
		inven.items[slot] = old
	}

	ch.listener.OnSwapInventorySlot(ch, constant.InventoryTypeEquipment, slot, int16(parts), int8(response.EQUIPMENT_ACTION_TYPE_ON))
	ch.listener.OnUpdateCharacterLook(ch)
	if swap {
		ch.mount.changeEquipment(parts)
	}
	return nil
}

func (inv *Inventory) validateMeso(costMeso, rewardMeso int32) ExchangeResult {
	if costMeso < 0 || rewardMeso < 0 {
		return ExchangeLackCost
	}
	if costMeso == 0 && rewardMeso == 0 {
		return ExchangeOK
	}
	if inv.meso < costMeso {
		return ExchangeLackCost
	}
	after := inv.meso - costMeso
	if rewardMeso > math.MaxInt32-after {
		return ExchangeLackCapacity
	}
	return ExchangeOK
}

func (inv *Inventory) Move(invType constant.InventoryType, src, dst int16) bool {
	inven := inv.tabs[invType]
	if inven == nil {
		return false
	}
	if dst < 1 || dst > int16(inven.slotLimit) || src == dst {
		return false
	}
	srcItem, ok := inven.items[src]
	if !ok {
		return false
	}

	dstItem, ok := inven.items[dst]
	if !ok {
		inven.items[dst] = srcItem
		delete(inven.items, src)
		inv.owner.listener.OnSwapInventorySlot(inv.owner, invType, src, dst, 0)
		return true
	}

	if !inven.Stackable(srcItem, dstItem) {
		inven.items[src], inven.items[dst] = dstItem, srcItem
		inv.owner.listener.OnSwapInventorySlot(inv.owner, invType, src, dst, 0)
		return true
	}

	moved := min(srcItem.GetCount(), srcItem.GetModel().GetCapacity()-dstItem.GetCount())
	dstItem.Increase(moved)
	if srcItem.Reduce(moved) == 0 {
		inv.owner.listener.OnFullMergeInventorySlot(inv.owner, invType, src, dst, dstItem.GetCount())
		delete(inven.items, src)
	} else {
		inv.owner.listener.OnPartialMergeInventorySlot(inv.owner, invType, src, dst, srcItem.GetCount(), dstItem.GetCount())
	}
	return true
}

func (inv *Inventory) canReceive(itemID uint32, count uint16) bool {
	return (ExchangeSpec{Reward: ExchangeSide{Items: map[uint32]uint16{itemID: count}}}).Valid(inv.owner, inv.owner.escrow) == ExchangeOK
}
