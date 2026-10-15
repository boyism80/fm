package game

import (
	"sort"

	"github.com/boyism80/fm/maple/constant"
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
)

func NewStorageFromInternalProto(pb *internal.Storage, owner *Character) *Storage {
	storage := &Storage{
		owner: owner,
		Slots: uint8(pb.GetSlots()),
		Meso:  pb.GetMeso(),
		Tabs:  map[constant.InventoryType][]Item{},
	}

	items := append([]*internal.Inventory(nil), pb.GetItems()...)
	sort.Slice(items, func(i, j int) bool {
		if items[i].GetInventoryType() != items[j].GetInventoryType() {
			return items[i].GetInventoryType() < items[j].GetInventoryType()
		}
		return items[i].GetSlot() < items[j].GetSlot()
	})
	for _, itemPB := range items {
		item, err := NewItemFromInternalProto(itemPB, owner.GameWorld)
		if err != nil {
			continue
		}
		storage.Tabs[item.GetInventoryType()] = append(storage.Tabs[item.GetInventoryType()], item)
	}
	return storage
}

func (s *Storage) BuildRPC(worldID uint32) *internal.Storage {
	pb := &internal.Storage{
		AccountId: s.owner.AccountID,
		WorldId:   worldID,
		Slots:     uint32(s.Slots),
		Meso:      s.Meso,
	}
	for _, items := range s.Tabs {
		for index, item := range items {
			pb.Items = append(pb.Items, item.BuildRPC(s.owner.AccountID, int32(index)))
		}
	}
	return pb
}
