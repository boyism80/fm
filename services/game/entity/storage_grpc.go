package entity

import (
	"sort"

	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
	"github.com/boyism80/fm/services/game/constant"
)

func NewStorageFromInternalProto(pb *internal.StoragePersisted, owner *Character) *Storage {
	storage := &Storage{
		owner: owner,
		Slots: uint8(pb.GetSlots()),
		Meso:  pb.GetMeso(),
		Tabs:  map[constant.InventoryType][]Item{},
	}

	items := append([]*internal.InventoryPersisted(nil), pb.GetItems()...)
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

func (s *Storage) ToProto(worldID uint32) *internal.StoragePersisted {
	if s == nil {
		return nil
	}
	pb := &internal.StoragePersisted{
		AccountId: s.owner.AccountID,
		WorldId:   worldID,
		Slots:     uint32(s.Slots),
		Meso:      s.Meso,
	}
	for _, items := range s.Tabs {
		for index, item := range items {
			pb.Items = append(pb.Items, item.ToProto(s.owner.AccountID, int32(index)))
		}
	}
	return pb
}
