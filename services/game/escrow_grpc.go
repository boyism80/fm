package game

import (
	"encoding/json"
	"log"
	"math"

	"github.com/boyism80/fm/maple/constant"
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
)

type escrowSlotKey struct {
	invType constant.InventoryType
	slot    int16
}

func (e *Escrow) mergeProto(items []*internal.Inventory) ([]*internal.Inventory, []Item) {
	if len(e.items) == 0 {
		return items, nil
	}

	var unplaced []Item

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
			if !found {
				unplaced = append(unplaced, held.item)
				continue
			}
			key = free
		}

		pb := held.item.BuildRPC(ownerID, int32(key.slot))
		taken[key] = pb
		items = append(items, pb)
	}
	return items, unplaced
}

func (e *Escrow) restoreRequests(items []Item) []*internal.UnsettledRequest {
	detail, ok := tradeRestoreDetailFrom(items)
	if !ok {
		return nil
	}
	data, err := json.Marshal(detail)
	if err != nil {
		log.Printf("Escrow.restoreRequests character=%d: %v", e.owner.GetID(), err)
		return nil
	}
	return []*internal.UnsettledRequest{{
		Id:     e.owner.gameWorld.NewUniqueID(),
		Type:   internal.UnsettledType_UNSETTLED_TYPE_TRADE_RESTORE,
		Detail: string(data),
	}}
}

func (e *Escrow) freeProtoSlot(taken map[escrowSlotKey]*internal.Inventory, invType constant.InventoryType) (escrowSlotKey, bool) {
	for slot := int16(1); slot <= int16(e.owner.inventory.tabs[invType].slotLimit); slot++ {
		key := escrowSlotKey{invType: invType, slot: slot}
		if _, ok := taken[key]; !ok {
			return key, true
		}
	}
	return escrowSlotKey{}, false
}

func (e *Escrow) mergeMeso(meso int32) int32 {
	return int32(min(int64(meso)+int64(e.meso), math.MaxInt32))
}
