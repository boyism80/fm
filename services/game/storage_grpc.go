package game

import (
	"sort"

	"github.com/boyism80/fm/maple/constant"
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
)

func NewStorageFromInternalProto(pb *internal.Storage, owner *Character) *Storage {
	storage := &Storage{
		owner: owner,
		slots: uint8(pb.GetSlots()),
		meso:  pb.GetMeso(),
		tabs:  map[constant.InventoryType][]Item{},
	}

	items := append([]*internal.Inventory(nil), pb.GetItems()...)
	sort.Slice(items, func(i, j int) bool {
		if items[i].GetInventoryType() != items[j].GetInventoryType() {
			return items[i].GetInventoryType() < items[j].GetInventoryType()
		}
		return items[i].GetSlot() < items[j].GetSlot()
	})
	for _, itemPB := range items {
		item, err := NewItemFromInternalProto(itemPB, owner.gameWorld)
		if err != nil {
			continue
		}
		storage.tabs[item.GetInventoryType()] = append(storage.tabs[item.GetInventoryType()], item)
	}
	return storage
}

func (s *Storage) BuildRPC(worldID uint32) *internal.Storage {
	pb := &internal.Storage{
		AccountId: s.owner.accountID,
		WorldId:   worldID,
		Slots:     uint32(s.slots),
		Meso:      s.meso,
	}
	for _, items := range s.tabs {
		for index, item := range items {
			pb.Items = append(pb.Items, item.BuildRPC(s.owner.accountID, int32(index)))
		}
	}
	return pb
}
