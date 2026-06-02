package entity

import (
	"errors"
	"fmt"

	"github.com/boyism80/fm/protocol/response"
	"github.com/boyism80/fm/services/game/constant"
)

func (ch *Character) AddMeso(amount int32) {
	if amount <= 0 {
		return
	}

	if ch.Meso > 0 && amount > 0 && ch.Meso+amount < ch.Meso {
		ch.Meso = int32(^uint32(0) >> 1)
	} else {
		ch.Meso += amount
	}

	ch.Listener.OnMesoChanged(ch, ch.Meso)
}

func (ch *Character) RemoveMeso(amount int32) {
	if amount <= 0 {
		return
	}

	if ch.Meso < amount {
		ch.Meso = 0
	} else {
		ch.Meso -= amount
	}

	ch.Listener.OnMesoChanged(ch, ch.Meso)
}

func (ch *Character) GetItem(invType constant.InventoryType, slot int16) Item {
	inven := ch.Inventory[invType]
	if inven == nil {
		return nil
	}
	if slot < 1 || slot > int16(inven.SlotLimit) {
		return nil
	}
	return inven.GetItem(uint8(slot))
}

func (ch *Character) RemoveItem(invType constant.InventoryType, slot int16, count uint16) bool {
	inven := ch.Inventory[invType]
	if inven == nil {
		return false
	}
	if slot < 1 || slot > int16(inven.SlotLimit) {
		return false
	}
	item := inven.GetItem(uint8(slot))
	if item == nil {
		return false
	}
	if item.GetCount() < count {
		return false
	}
	item.Reduce(count)
	if item.GetCount() == 0 {
		if err := inven.RemoveItem(uint8(slot)); err != nil {
			return false
		}
		ch.Listener.OnRemoveInventorySlot(ch, invType, slot)
	} else {
		ch.Listener.OnInventorySlotUpdated(ch, invType, slot, item)
	}
	return true
}

func (ch *Character) FindSlots(itemID uint32) (invType constant.InventoryType, slots []int16) {
	invType = constant.GetInventoryTypeByItemID(itemID)
	inven := ch.Inventory[invType]
	if inven == nil {
		return invType, nil
	}
	for slot := int16(1); slot <= int16(inven.SlotLimit); slot++ {
		item := inven.GetItem(uint8(slot))
		if item != nil && item.GetModel().GetID() == itemID {
			slots = append(slots, slot)
		}
	}
	return invType, slots
}

type SlotPredicate func(invType constant.InventoryType, slot int16, item Item) bool

func (ch *Character) FindSlotsWhere(pred SlotPredicate) map[constant.InventoryType][]int16 {
	out := make(map[constant.InventoryType][]int16)
	if pred == nil {
		return out
	}
	for invType, inven := range ch.Inventory {
		if inven == nil {
			continue
		}
		for slot := int16(1); slot <= int16(inven.SlotLimit); slot++ {
			item := inven.GetItem(uint8(slot))
			if pred(invType, slot, item) {
				out[invType] = append(out[invType], slot)
			}
		}
	}
	return out
}

func (ch *Character) GetCountByItemID(itemID uint32) uint16 {
	invType, slots := ch.FindSlots(itemID)
	var total uint16
	for _, slot := range slots {
		if item := ch.GetItem(invType, slot); item != nil {
			total += item.GetCount()
		}
	}
	return total
}

func (ch *Character) HasItem(itemID uint32) bool {
	return ch.GetCountByItemID(itemID) >= 1
}

func (ch *Character) HasItemCount(itemID uint32, count uint16) bool {
	if count == 0 {
		return true
	}
	return ch.GetCountByItemID(itemID) >= count
}

func (ch *Character) RemoveByItemIDCount(itemID uint32, count uint16) bool {
	if count == 0 {
		return true
	}
	invType, slots := ch.FindSlots(itemID)
	var total uint16
	for _, slot := range slots {
		if item := ch.GetItem(invType, slot); item != nil {
			total += item.GetCount()
		}
	}
	if total < count {
		return false
	}
	remaining := count
	for _, slot := range slots {
		if remaining == 0 {
			break
		}
		item := ch.GetItem(invType, slot)
		if item == nil {
			continue
		}
		take := item.GetCount()
		if take > remaining {
			take = remaining
		}
		remaining -= take
		ch.RemoveItem(invType, slot, take)
	}
	return true
}

func (ch *Character) RemoveItemByID(itemID uint32) bool {
	return ch.RemoveByItemIDCount(itemID, 1)
}

func (ch *Character) AddItem(item Item, allOrNothing bool) (addedItems []Item, err error) {
	if item == nil {
		return nil, fmt.Errorf("item is nil")
	}

	invenType := item.GetInventoryType()
	inven := ch.Inventory[invenType]
	if inven == nil {
		return nil, fmt.Errorf("inventory type %d not found", invenType)
	}

	model := item.GetModel()
	requestedCount := item.GetCount()

	if allOrNothing {
		if !inven.IsFree(model, requestedCount) {
			return nil, fmt.Errorf("not enough inventory space for %d items", requestedCount)
		}
	}

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
		ch.Listener.OnShowItemGain(ch, model.GetID(), uint32(addedCount), constant.ShowItemGainTypeStatus)
	}

	return addedItems, nil
}

func (ch *Character) GainMeso(amount int32) {
	if amount <= 0 {
		return
	}

	if ch.Meso > 0 && amount > 0 && ch.Meso+amount < ch.Meso {
		ch.Meso = int32(^uint32(0) >> 1)
	} else {
		ch.Meso += amount
	}

	ch.Listener.OnMesoChanged(ch, ch.Meso)
	ch.Listener.OnShowMesoGain(ch, amount, constant.ShowMesoGainTypeStatus)
	ch.Listener.OnUpdateStats(ch, map[constant.Stat]int32{
		constant.StatMeso: ch.Meso,
	}, false)
}

func (ch *Character) FindSlot(invType constant.InventoryType, item Item) (int16, bool) {
	inven := ch.Inventory[invType]
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

func (ch *Character) UnequipToSlot(parts constant.EquipmentPartsType, destSlot int16) error {
	equipments := ch.Equipments
	inventory := ch.Inventory
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

func (ch *Character) Unequip(parts constant.EquipmentPartsType) error {
	inven := ch.Inventory[constant.InventoryTypeEquipment]
	if inven == nil {
		return errors.New("equipment inventory not found")
	}
	destSlot, ok := inven.NextSlot()
	if !ok {
		return ErrInventoryFull
	}
	return ch.UnequipToSlot(parts, int16(destSlot))
}

func (ch *Character) Equip(slot int16) error {
	equipments := ch.Equipments
	inventory := ch.Inventory
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
				if err := ch.UnequipToSlot(constant.EquipmentPartsPants, int16(storageSlot)); err != nil {
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
				if err := ch.UnequipToSlot(constant.EquipmentPartsTop, int16(storageSlot)); err != nil {
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
