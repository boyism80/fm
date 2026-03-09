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

func (ch *Character) AddItem(item Item, allOrNothing bool) (uint16, error) {
	if item == nil {
		return 0, fmt.Errorf("item is nil")
	}

	invenType := item.GetInventoryType()
	inven := ch.Inventory[invenType]
	if inven == nil {
		return 0, fmt.Errorf("inventory type %d not found", invenType)
	}

	model := item.GetModel()
	requestedCount := item.GetCount()

	if allOrNothing {
		if !inven.IsFree(model, requestedCount) {
			return 0, fmt.Errorf("not enough inventory space for %d items", requestedCount)
		}
	}

	remainingCount := requestedCount
	addedCount := uint16(0)

	for remainingCount > 0 {
		slot, ok := inven.FindSlot(model)
		if !ok {
			if allOrNothing && addedCount == 0 {
				return 0, fmt.Errorf("no available slot found")
			}
			break
		}

		exists, ok := inven.Items[int16(slot)]
		cap := uint16(0)
		if ok {
			cap = min(model.GetCapacity()-exists.GetCount(), remainingCount)
			exists.Increase(cap)
			if ch.Listener != nil {
				ch.Listener.OnInventorySlotUpdated(invenType, int16(slot), exists)
			}
		} else {
			cap = min(model.GetCapacity(), remainingCount)
			inven.Items[int16(slot)] = item.Clone(cap)
			if ch.Listener != nil {
				ch.Listener.OnInventorySlotAdded(invenType, int16(slot), inven.Items[int16(slot)])
			}
		}
		remainingCount -= cap
		addedCount += cap
	}

	if !allOrNothing && addedCount < requestedCount {
		item.SetCount(remainingCount)
	}

	if ch.Listener != nil && addedCount > 0 {
		ch.Listener.OnShowItemGain(model.GetID(), uint32(addedCount), constant.ShowItemGainTypeStatus)
	}

	return addedCount, nil
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
