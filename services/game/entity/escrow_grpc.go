package entity

import (
	"log"
	"math"

	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
	"github.com/boyism80/fm/services/game/constant"
)

type escrowSlotKey struct {
	invType constant.InventoryType
	slot    int16
}

func (e *Escrow) mergeProto(items []*internal.Inventory) []*internal.Inventory {
	if len(e.items) == 0 {
		return items
	}

	ownerID := e.owner.GetID()
	taken := make(map[escrowSlotKey]*internal.Inventory, len(items))
	for _, pb := range items {
		if pb.GetSlot() > 0 {
			taken[escrowSlotKey{invType: constant.InventoryType(pb.GetInventoryType()), slot: int16(pb.GetSlot())}] = pb
		}
	}

	for _, held := range e.items {
		model := held.item.GetModel()
		count := uint32(held.item.GetCount())
		key := escrowSlotKey{invType: held.invType, slot: held.slot}
		if exists, ok := taken[key]; ok {
			if exists.GetItemId() == model.GetID() && exists.GetCount()+count <= uint32(model.GetCapacity()) {
				exists.Count += count
				continue
			}

			free, found := e.freeProtoSlot(taken, held.invType)
			if found == false {
				log.Printf("Escrow.mergeProto character=%d: no slot for item %d x%d", ownerID, model.GetID(), count)
				continue
			}
			key = free
		}

		pb := held.item.ToProto(ownerID, int32(key.slot))
		taken[key] = pb
		items = append(items, pb)
	}
	return items
}

func (e *Escrow) freeProtoSlot(taken map[escrowSlotKey]*internal.Inventory, invType constant.InventoryType) (escrowSlotKey, bool) {
	for slot := int16(1); slot <= int16(e.owner.Inventory.Tabs[invType].SlotLimit); slot++ {
		key := escrowSlotKey{invType: invType, slot: slot}
		if _, ok := taken[key]; ok == false {
			return key, true
		}
	}
	return escrowSlotKey{}, false
}

func (e *Escrow) mergeMeso(meso int32) int32 {
	return int32(min(int64(meso)+int64(e.Meso), math.MaxInt32))
}
