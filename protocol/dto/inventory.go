package dto

import (
	"sort"

	"github.com/boyism80/fm/services/game/constant"
	"github.com/boyism80/fm/stream"
)

type Inventory struct {
	Type      constant.InventoryType
	SlotLimit uint8
	Items     map[int16]Item
}

func (inv *Inventory) Serialize(writer *stream.StreamWriter) {
	if inv == nil || inv.Items == nil {
		writer.WriteU8(0)
		return
	}

	slots := make([]int16, 0, len(inv.Items))
	for slot, item := range inv.Items {
		if item != nil {
			slots = append(slots, slot)
		}
	}
	sort.Slice(slots, func(i, j int) bool {
		return slots[i] < slots[j]
	})

	for _, slot := range slots {
		if item := inv.Items[slot]; item != nil {
			item.Serialize(writer, true, slot)
		}
	}
	writer.WriteU8(0)
}
