package dto

import (
	"sort"

	"github.com/boyism80/fm/game/constant"
	"github.com/boyism80/fm/stream"
)

// Inventory represents inventory data for protocol
type Inventory struct {
	Type      constant.InventoryType
	SlotLimit uint8
	Items     map[int16]Item // slot -> item (all item types, using interface)
}

// Serialize serializes inventory data
func (inv *Inventory) Serialize(writer *stream.StreamWriter) {
	if inv == nil || inv.Items == nil {
		writer.WriteU8(0)
		return
	}

	// Sort slots to ensure consistent serialization order
	slots := make([]int16, 0, len(inv.Items))
	for slot, item := range inv.Items {
		if item != nil {
			slots = append(slots, slot)
		}
	}
	sort.Slice(slots, func(i, j int) bool {
		return slots[i] < slots[j] // ascending order
	})

	for _, slot := range slots {
		if item := inv.Items[slot]; item != nil {
			item.Serialize(writer, true, slot)
		}
	}
	writer.WriteU8(0)
}
