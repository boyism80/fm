package game

import (
	"sort"

	"github.com/boyism80/fm/maple/constant"
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
)

type storageSlotKey struct {
	invType constant.InventoryType
	slot    int16
}

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
			storage.kept = append(storage.kept, itemPB)
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
	taken := map[storageSlotKey]struct{}{}
	for _, items := range s.tabs {
		for index, item := range items {
			row := item.BuildRPC(s.owner.accountID, int32(index))
			pb.Items = append(pb.Items, row)
			taken[storageSlotKey{invType: constant.InventoryType(row.GetInventoryType()), slot: int16(row.GetSlot())}] = struct{}{}
		}
	}
	for _, item := range s.kept {
		key := storageSlotKey{invType: constant.InventoryType(item.GetInventoryType()), slot: int16(item.GetSlot())}
		if _, ok := taken[key]; ok {
			slot := int16(0)
			for {
				key.slot = slot
				if _, used := taken[key]; !used {
					item.Slot = int32(slot)
					break
				}
				slot++
			}
		}
		taken[key] = struct{}{}
		pb.Items = append(pb.Items, item)
	}
	return pb
}
