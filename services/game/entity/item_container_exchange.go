package entity

import (
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
	emptySlots := uint16(0)
	for slot := int16(1); slot <= int16(inv.SlotLimit); slot++ {
		item := inv.Items[slot]
		if item == nil {
			emptySlots++
			continue
		}
		slots[slot] = itemSlotSnapshot{
			itemID: item.GetModel().GetID(),
			count:  item.GetCount(),
		}
	}

	slotsByID := make(map[uint32][]int16)
	for slot, snap := range slots {
		slotsByID[snap.itemID] = append(slotsByID[snap.itemID], slot)
	}

	effectiveFreeSlots := emptySlots

	for id, costCount := range cost {
		if costCount == 0 {
			continue
		}
		model := modelOf(id)
		if model == nil {
			return ExchangeLackCost
		}
		indices := slotsByID[id]
		if len(indices) == 0 {
			return ExchangeLackCost
		}

		if model.GetCapacity() > 1 {
			slot := indices[0]
			snap, ok := slots[slot]
			if !ok {
				return ExchangeLackCost
			}
			if snap.count < costCount {
				return ExchangeLackCost
			}
			remain := snap.count - costCount
			if remain == 0 {
				delete(slots, slot)
				effectiveFreeSlots++
			} else {
				snap.count = remain
				slots[slot] = snap
			}
		} else {
			if uint16(len(indices)) < costCount {
				return ExchangeLackCost
			}
			for i := uint16(0); i < costCount; i++ {
				slot := indices[i]
				if _, ok := slots[slot]; !ok {
					return ExchangeLackCost
				}
				delete(slots, slot)
				effectiveFreeSlots++
			}
		}
	}

	bundleRemaining := make(map[uint32]uint16)
	for _, snap := range slots {
		model := modelOf(snap.itemID)
		if model == nil {
			continue
		}
		if model.GetCapacity() > 1 {
			bundleRemaining[snap.itemID] += snap.count
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
		if model.GetCapacity() > 1 {
			requiredSlots++
		} else {
			requiredSlots += int(count)
		}
	}

	availableSlots := int(effectiveFreeSlots)
	for id, count := range reward {
		if count == 0 {
			continue
		}
		model := modelOf(id)
		if model == nil {
			return ExchangeLackCapacity
		}
		if model.GetCapacity() <= 1 {
			continue
		}
		existing := bundleRemaining[id]
		if uint16(model.GetCapacity()) < existing+count {
			return ExchangeLackCapacity
		}
		availableSlots++
	}

	if availableSlots < requiredSlots {
		return ExchangeLackCapacity
	}

	return ExchangeOK
}
