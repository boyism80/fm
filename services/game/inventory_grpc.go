package game

import (
	"github.com/boyism80/fm/maple/constant"
	"github.com/boyism80/fm/maple/wz"
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
)

func (inv *Inventory) Load(items []*internal.Inventory) {
	for _, pb := range items {
		if pb == nil {
			continue
		}
		item, err := NewItemFromInternalProto(pb, inv.owner.gameWorld)
		if err != nil {
			continue
		}
		if inv.owner.wedding.marriage == nil && item.GetModel().GetTraits()&(constant.ItemTraitEngagementRing|constant.ItemTraitWeddingRing) != 0 {
			continue
		}
		slot := int16(pb.GetSlot())
		if slot < 0 {
			eq, ok := item.(Equipment)
			if !ok {
				continue
			}
			model := eq.GetModel().(wz.Equipment)
			parts := constant.EquipmentPartsType(slot)
			if !model.CanEquipAt(parts) {
				if known := model.GetParts(); len(known) > 0 && inv.equipped[known[0]] == nil {
					parts = known[0]
				}
			}
			inv.equipped[parts] = eq
		} else {
			invType := constant.InventoryType(pb.GetInventoryType())
			if invType == 0 {
				invType = item.GetModel().GetInventoryType()
			}
			if tab, ok := inv.tabs[invType]; ok {
				tab.items[slot] = item
			}
		}
	}
}

func (inv *Inventory) BuildRPC() []*internal.Inventory {
	items := make([]*internal.Inventory, 0, len(inv.equipped)+64)
	ownerID := inv.owner.GetID()
	for parts, item := range inv.equipped {
		if item == nil {
			continue
		}
		if pb := item.BuildRPC(ownerID, int32(parts)); pb != nil {
			items = append(items, pb)
		}
	}
	for _, tab := range inv.tabs {
		items = append(items, tab.BuildRPC(ownerID)...)
	}
	return inv.owner.escrow.mergeProto(items)
}
