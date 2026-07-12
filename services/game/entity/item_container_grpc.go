package entity

import (
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
)

func (m *ItemContainer) ToProto(ownerID uint32) []*internal.InventoryPersisted {
	if m == nil {
		return nil
	}
	out := make([]*internal.InventoryPersisted, 0, len(m.Items))
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
