package entity

import (
	"errors"
	"fmt"
	"log"

	"github.com/boyism80/fm/core/luax"
	"github.com/boyism80/fm/protocol/response"
	"github.com/boyism80/fm/services/game/constant"
)

func (inv *Inventory) AddMeso(amount int32) {
	ch := inv.owner
	if amount <= 0 {
		return
	}
	if ch.validateMesoExchange(0, amount) != ExchangeOK {
		return
	}
	inv.addMesoUnchecked(amount)
}

func (inv *Inventory) addMesoUnchecked(amount int32) {
	ch := inv.owner
	if amount <= 0 {
		return
	}

	if inv.Meso > 0 && amount > 0 && inv.Meso+amount < inv.Meso {
		inv.Meso = int32(^uint32(0) >> 1)
	} else {
		inv.Meso += amount
	}

	ch.Listener.OnMesoChanged(ch, inv.Meso)
}

func (inv *Inventory) RemoveMeso(amount int32) {
	ch := inv.owner
	if amount <= 0 {
		return
	}
	if ch.validateMesoExchange(amount, 0) != ExchangeOK {
		return
	}
	inv.removeMesoUnchecked(amount)
}

func (inv *Inventory) removeMesoUnchecked(amount int32) {
	ch := inv.owner
	if amount <= 0 {
		return
	}

	if inv.Meso < amount {
		inv.Meso = 0
	} else {
		inv.Meso -= amount
	}

	ch.Listener.OnMesoChanged(ch, inv.Meso)
}

func (inv *Inventory) GetItem(invType constant.InventoryType, slot int16) Item {
	inven := inv.Tabs[invType]
	if inven == nil {
		return nil
	}
	if slot < 1 || slot > int16(inven.SlotLimit) {
		return nil
	}
	return inven.Get(uint8(slot))
}

func (inv *Inventory) RemoveItem(invType constant.InventoryType, slot int16, count uint16) bool {
	ch := inv.owner
	inven := inv.Tabs[invType]
	if inven == nil {
		return false
	}
	if slot < 1 || slot > int16(inven.SlotLimit) {
		return false
	}
	item := inven.Get(uint8(slot))
	if item == nil {
		return false
	}
	if item.GetCount() < count {
		return false
	}
	item.Reduce(count)
	if item.GetCount() == 0 {
		if err := inven.Remove(uint8(slot)); err != nil {
			return false
		}
		ch.Listener.OnRemoveInventorySlot(ch, invType, slot)
	} else {
		ch.Listener.OnInventorySlotUpdated(ch, invType, slot, item)
	}
	return true
}

func (inv *Inventory) FindSlots(itemID uint32) (invType constant.InventoryType, slots []int16) {
	invType = constant.GetInventoryTypeByItemID(itemID)
	inven := inv.Tabs[invType]
	if inven == nil {
		return invType, nil
	}
	for slot := int16(1); slot <= int16(inven.SlotLimit); slot++ {
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
	invType := constant.GetInventoryTypeByItemID(itemID)
	inven := inv.Tabs[invType]
	if inven == nil {
		return false
	}
	for slot := int16(1); slot <= int16(inven.SlotLimit); slot++ {
		item := inven.Get(uint8(slot))
		if item != nil && item.GetModel().GetID() == itemID {
			return true
		}
	}
	return false
}

func (inv *Inventory) HasItemCount(itemID uint32, count uint16) bool {
	if count == 0 {
		return true
	}
	invType := constant.GetInventoryTypeByItemID(itemID)
	inven := inv.Tabs[invType]
	if inven == nil {
		return false
	}
	var total uint16
	for slot := int16(1); slot <= int16(inven.SlotLimit); slot++ {
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
	if spec.Valid(ch) != ExchangeOK {
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
	for invType, inven := range inv.Tabs {
		if inven == nil || inven.Items == nil {
			continue
		}
		slots := make([]int16, 0, len(inven.Items))
		for slot, item := range inven.Items {
			if item == nil {
				continue
			}
			slots = append(slots, slot)
		}
		for _, slot := range slots {
			delete(inven.Items, slot)
			if ch.Listener != nil {
				ch.Listener.OnRemoveInventorySlot(ch, invType, slot)
			}
			cleared++
		}
	}
	return cleared
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
		if spec.Valid(ch) != ExchangeOK {
			return nil, fmt.Errorf("not enough inventory space for %d items", item.GetCount())
		}
	}

	return inv.applyAddItem(item, allOrNothing)
}

func (inv *Inventory) applyAddItem(item Item, allOrNothing bool) (addedItems []Item, err error) {
	ch := inv.owner
	if item == nil {
		return nil, fmt.Errorf("item is nil")
	}

	invenType := item.GetInventoryType()
	inven := inv.Tabs[invenType]
	if inven == nil {
		return nil, fmt.Errorf("inventory type %d not found", invenType)
	}

	model := item.GetModel()
	requestedCount := item.GetCount()

	remainingCount := requestedCount
	var addedCount uint16

	for remainingCount > 0 {
		slot, ok := inven.FindSlot(model)
		if !ok {
			if allOrNothing && addedCount == 0 {
				return nil, fmt.Errorf("no available slot found")
			}
			break
		}

		exists, ok := inven.Items[int16(slot)]
		cap := uint16(0)
		var slotItem Item
		if ok {
			cap = min(model.GetCapacity()-exists.GetCount(), remainingCount)
			exists.Increase(cap)
			slotItem = exists
			ch.Listener.OnInventorySlotUpdated(ch, invenType, int16(slot), exists)
		} else {
			cap = min(model.GetCapacity(), remainingCount)
			placed := item.Clone(cap)
			inven.Items[int16(slot)] = placed
			slotItem = placed
			ch.Listener.OnInventorySlotAdded(ch, invenType, int16(slot), placed)
		}
		addedItems = append(addedItems, slotItem)
		remainingCount -= cap
		addedCount += cap
	}

	if !allOrNothing && addedCount < requestedCount {
		item.SetCount(remainingCount)
	}

	if addedCount > 0 {
		ch.NotifyItemGained(model.GetID())
		ch.Listener.OnShowItemGain(ch, model.GetID(), uint32(addedCount), constant.ShowItemGainTypeStatus)
	}

	return addedItems, nil
}

func (ch *Character) NotifyItemGained(itemID uint32) {
	if ch == nil || itemID == 0 {
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
	luax.SetConfiguration(thread, luax.Configuration{
		MapActorPID: mapInstance.GetActorPID(),
	})
	hook := fmt.Sprintf("on_item_gain_%d", itemID)
	luax.CallAsync(root, thread, hook, ch, itemID).OnError(func(err error) {
		log.Printf("item gain script %s: %v", hook, err)
	})
}

func (inv *Inventory) GainMeso(amount int32) {
	ch := inv.owner
	if amount <= 0 {
		return
	}
	if ch.validateMesoExchange(0, amount) != ExchangeOK {
		return
	}
	inv.gainMesoUnchecked(amount)
}

func (inv *Inventory) gainMesoUnchecked(amount int32) {
	ch := inv.owner
	if amount <= 0 {
		return
	}

	if inv.Meso > 0 && amount > 0 && inv.Meso+amount < inv.Meso {
		inv.Meso = int32(^uint32(0) >> 1)
	} else {
		inv.Meso += amount
	}

	ch.Listener.OnMesoChanged(ch, inv.Meso)
	ch.Listener.OnShowMesoGain(ch, amount, constant.ShowMesoGainTypeStatus)
	ch.Listener.OnUpdateStats(ch, map[constant.Stat]int32{
		constant.StatMeso: inv.Meso,
	}, false)
}

func (inv *Inventory) FindSlot(invType constant.InventoryType, item Item) (int16, bool) {
	inven := inv.Tabs[invType]
	if inven == nil || item == nil {
		return 0, false
	}
	for slot := int16(1); slot <= int16(inven.SlotLimit); slot++ {
		if inven.Items[slot] == item {
			return slot, true
		}
	}
	return 0, false
}

func (inv *Inventory) UnequipToSlot(parts constant.EquipmentPartsType, destSlot int16) error {
	ch := inv.owner
	equipments := inv.Equipped
	inventory := inv.Tabs
	if equipments[parts] == nil {
		return nil
	}
	inven := inventory[constant.InventoryTypeEquipment]
	if inven == nil || inven.Items[destSlot] != nil {
		return ErrSlotAlreadyOccupied
	}
	inven.Items[destSlot] = equipments[parts]
	delete(equipments, parts)
	ch.Listener.OnSwapInventorySlot(ch, constant.InventoryTypeEquipment, int16(parts), destSlot, int8(response.EQUIPMENT_ACTION_TYPE_OFF))
	ch.Listener.OnUpdateCharacterLook(ch)
	return nil
}

func (inv *Inventory) Unequip(parts constant.EquipmentPartsType) error {
	inven := inv.Tabs[constant.InventoryTypeEquipment]
	if inven == nil {
		return errors.New("equipment inventory not found")
	}
	destSlot, ok := inven.NextSlot()
	if !ok {
		return ErrInventoryFull
	}
	return inv.UnequipToSlot(parts, int16(destSlot))
}

func (inv *Inventory) Equip(slot int16) error {
	ch := inv.owner
	equipments := inv.Equipped
	inventory := inv.Tabs
	inven := inventory[constant.InventoryTypeEquipment]
	if inven == nil {
		return errors.New("equipment inventory not found")
	}
	item := inven.Items[slot]
	if item == nil {
		return ErrSourceSlotEmpty
	}
	newEq, ok := item.(Equipment)
	if !ok {
		return ErrItemNotEquipment
	}
	parts := constant.GetEquipmentPartsType(newEq.GetModel().GetID())
	if parts == 0 {
		return ErrInvalidEquipmentPart
	}

	old, swap := equipments[parts]

	switch parts {
	case constant.EquipmentPartsTop:
		if topNew, ok := newEq.(*Top); ok && topNew.IsOverall() {
			if equipments[constant.EquipmentPartsPants] != nil {
				storageSlot, isFree := inven.NextSlot()
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
			if top, ok := topEq.(*Top); ok && top.IsOverall() {
				storageSlot, isFree := inven.NextSlot()
				if swap && !isFree {
					return ErrInventoryFull
				}
				if err := inv.UnequipToSlot(constant.EquipmentPartsTop, int16(storageSlot)); err != nil {
					return err
				}
			}
		}
	}

	equipments[parts] = newEq
	if old == nil {
		delete(inven.Items, slot)
	} else {
		inven.Items[slot] = old
	}

	ch.Listener.OnSwapInventorySlot(ch, constant.InventoryTypeEquipment, slot, int16(parts), int8(response.EQUIPMENT_ACTION_TYPE_ON))
	ch.Listener.OnUpdateCharacterLook(ch)
	return nil
}
