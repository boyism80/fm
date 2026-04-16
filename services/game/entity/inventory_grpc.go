package entity

import (
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
)

// ToPersisted returns internal.InventoryPersisted rows for all non-nil slots.
func (m *Inventory) ToPersisted(ownerID uint32) []*internal.InventoryPersisted {
	if m == nil {
		return nil
	}
	out := make([]*internal.InventoryPersisted, 0, len(m.Items))
	for slot, item := range m.Items {
		if item == nil {
			continue
		}
		if pb := item.ToPersisted(ownerID, int32(slot)); pb != nil {
			out = append(out, pb)
		}
	}
	return out
}
