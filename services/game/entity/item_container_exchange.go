package entity

import (
	"github.com/boyism80/fm/services/game/constant"
	"github.com/boyism80/fm/services/game/wz"
)

type itemSlotSnapshot struct {
	itemID uint32
	count  uint16
}

func (inv *ItemContainer) validateItemExchange(cost map[uint32]uint16, reward map[uint32]uint16, modelOf func(uint32) wz.Item) ExchangeResult {
	if inv == nil {
		if !exchangeItemsEmpty(cost) {
			return ExchangeLackCost
		}
		if !exchangeItemsEmpty(reward) {
			return ExchangeLackCapacity
		}
		return ExchangeOK
	}
	if exchangeItemsEmpty(cost) && exchangeItemsEmpty(reward) {
		return ExchangeOK
	}

	slots := make(map[int16]itemSlotSnapshot, inv.SlotLimit)
	slotsByID := make(map[uint32][]int16)
	freeSlots := 0
	for slot := int16(1); slot <= int16(inv.SlotLimit); slot++ {
		item := inv.Items[slot]
		if item == nil {
			freeSlots++
			continue
		}
		id := item.GetModel().GetID()
		slots[slot] = itemSlotSnapshot{itemID: id, count: item.GetCount()}
		slotsByID[id] = append(slotsByID[id], slot)
	}

	for id, costCount := range cost {
		if costCount == 0 {
			continue
		}
		if modelOf(id) == nil {
			return ExchangeLackCost
		}
		remaining := costCount
		for _, slot := range slotsByID[id] {
			if remaining == 0 {
				break
			}
			snap := slots[slot]
			take := min(snap.count, remaining)
			remaining -= take
			snap.count -= take
			if snap.count == 0 {
				delete(slots, slot)
				freeSlots++
			} else {
				slots[slot] = snap
			}
		}
		if remaining > 0 {
			return ExchangeLackCost
		}
	}

	requiredSlots := 0
	for id, count := range reward {
		if count == 0 {
			continue
		}
		model := modelOf(id)
		if model == nil {
			return ExchangeLackCapacity
		}
		capacity := max(int(model.GetCapacity()), 1)
		need := int(count)
		if capacity > 1 && constant.IsRechargeable(id) == false {
			for _, slot := range slotsByID[id] {
				snap, ok := slots[slot]
				if ok == false {
					continue
				}
				need -= min(need, max(capacity-int(snap.count), 0))
			}
		}
		requiredSlots += (need + capacity - 1) / capacity
	}

	if requiredSlots > freeSlots {
		return ExchangeLackCapacity
	}
	return ExchangeOK
}
