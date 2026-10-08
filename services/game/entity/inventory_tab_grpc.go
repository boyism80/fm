package entity

import (
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
)

func (m *InventoryTab) ToProto(ownerID uint32) []*internal.Inventory {
	if m == nil {
		return nil
	}
	out := make([]*internal.Inventory, 0, len(m.Items))
	for slot, item := range m.Items {
		if item == nil {
			continue
		}
		if pb := item.ToProto(ownerID, int32(slot)); pb != nil {
			out = append(out, pb)
		}
	}
	return out
}
