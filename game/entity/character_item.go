package entity

import (
	"fmt"

	"github.com/boyism80/fm/game/constant"
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

	if ch.Listener != nil {
		ch.Listener.OnMesoChanged(ch.Meso)
	}
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

	if ch.Listener != nil {
		ch.Listener.OnMesoChanged(ch.Meso)
	}
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
			if ch.Listener != nil {
				ch.Listener.OnInventorySlotUpdated(invenType, int16(slot), exists)
			}
		} else {
			cap = min(model.GetCapacity(), remainingCount)
			placed := item.Clone(cap)
			inven.Items[int16(slot)] = placed
			slotItem = placed
			if ch.Listener != nil {
				ch.Listener.OnInventorySlotAdded(invenType, int16(slot), placed)
			}
		}
		addedItems = append(addedItems, slotItem)
		remainingCount -= cap
		addedCount += cap
	}

	if !allOrNothing && addedCount < requestedCount {
		item.SetCount(remainingCount)
	}

	if ch.Listener != nil && addedCount > 0 {
		ch.Listener.OnShowItemGain(model.GetID(), uint32(addedCount), constant.ShowItemGainTypeStatus)
	}

	return addedItems, nil
}

// RemoveItemByID removes one slot of the first matching item by ID from any inventory.
// Returns true if an item was removed.
func (ch *Character) RemoveItemByID(itemId uint32) bool {
	for invType, inven := range ch.Inventory {
		if inven == nil {
			continue
		}
		for slot, item := range inven.Items {
			if item != nil && item.GetModel().GetID() == itemId {
				if err := inven.RemoveItem(uint8(slot)); err != nil {
					return false
				}
				if ch.Listener != nil {
					ch.Listener.OnRemoveInventorySlot(invType, slot)
				}
				return true
			}
		}
	}
	return false
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

	if ch.Listener != nil {
		ch.Listener.OnMesoChanged(ch.Meso)
		ch.Listener.OnShowMesoGain(amount, constant.ShowMesoGainTypeStatus)
		ch.Listener.OnUpdateStats(map[constant.Stat]int32{
			constant.STAT_MESO: ch.Meso,
		}, false)
	}
}
