package game

import (
	"encoding/json"
	"log"

	"github.com/boyism80/fm/maple/constant"
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
)

type escrowItem struct {
	invType constant.InventoryType
	slot    int16
	item    Item
}

type Escrow struct {
	owner *Character
	items []escrowItem
	meso  int32
}

func (e *Escrow) hold(invType constant.InventoryType, slot int16, count uint16) (Item, error) {
	item := e.owner.inventory.GetItem(invType, slot)
	if item == nil {
		return nil, ErrMiniRoomItemNotFound
	}
	if count == 0 || count > item.GetCount() {
		return nil, ErrMiniRoomInvalid
	}

	held := item.Clone(count)
	e.owner.inventory.RemoveItem(invType, slot, count)
	e.items = append(e.items, escrowItem{invType: invType, slot: slot, item: held})
	return held, nil
}

func (e *Escrow) holdMeso(amount int32) error {
	if amount <= 0 || !e.owner.inventory.RemoveMeso(amount) {
		return ErrMiniRoomInvalid
	}
	e.meso += amount
	return nil
}

func (e *Escrow) restore() {
	ch := e.owner
	var unplaced []Item
	for _, held := range e.items {
		inven := ch.inventory.tabs[held.invType]
		model := held.item.GetModel()
		exists := inven.items[held.slot]
		switch {
		case exists == nil:
			inven.items[held.slot] = held.item
			ch.listener.OnInventorySlotAdded(ch, held.invType, held.slot, held.item)
		case exists.GetModel() == model && exists.GetCount()+held.item.GetCount() <= model.GetCapacity():
			exists.Increase(held.item.GetCount())
			ch.listener.OnInventorySlotUpdated(ch, held.invType, held.slot, exists)
		default:
			slot, ok := inven.NextSlot()
			if !ok {
				unplaced = append(unplaced, held.item)
				continue
			}
			inven.items[int16(slot)] = held.item
			ch.listener.OnInventorySlotAdded(ch, held.invType, int16(slot), held.item)
		}
	}
	e.items = nil
	if len(unplaced) > 0 {
		e.keepUnplaced(unplaced)
	}

	if e.meso > 0 {
		ch.inventory.addMesoUnchecked(e.meso)
		e.meso = 0
	}
}

type tradeRestoreDetail struct {
	Items []json.RawMessage `json:"items"`
}

func tradeRestoreDetailFrom(items []Item) (tradeRestoreDetail, bool) {
	detail := tradeRestoreDetail{Items: make([]json.RawMessage, 0, len(items))}
	for _, item := range items {
		raw, err := json.Marshal(item.BuildRPC(0, 0))
		if err != nil {
			log.Printf("trade restore: %v", err)
			continue
		}
		detail.Items = append(detail.Items, raw)
	}
	return detail, len(detail.Items) > 0
}

func (e *Escrow) keepUnplaced(items []Item) {
	detail, ok := tradeRestoreDetailFrom(items)
	if !ok {
		return
	}
	id := e.owner.unsettled.Add(internal.UnsettledType_UNSETTLED_TYPE_TRADE_RESTORE, 0, detail)
	e.owner.unsettled.Lost(id)
}

func (e *Escrow) settleRestore(request *UnsettledRequest) {
	var detail tradeRestoreDetail
	if err := json.Unmarshal([]byte(request.detail), &detail); err != nil {
		request.inFlight = false
		log.Printf("Escrow.settleRestore character=%d id=%d: %v", e.owner.GetID(), request.id, err)
		return
	}

	var left []json.RawMessage
	for _, raw := range detail.Items {
		var pb internal.Inventory
		if err := json.Unmarshal(raw, &pb); err != nil {
			left = append(left, raw)
			continue
		}
		item, err := NewItemFromInternalProto(&pb, e.owner.gameWorld)
		if err != nil {
			left = append(left, raw)
			continue
		}
		if _, err := e.owner.inventory.AddItem(item, true); err != nil {
			left = append(left, raw)
		}
	}
	if len(left) > 0 {
		detail.Items = left
		if data, err := json.Marshal(detail); err == nil {
			request.detail = string(data)
		}
		request.inFlight = false
		return
	}
	e.owner.unsettled.Erase(request.id)
}

func (e *Escrow) release() ([]Item, int32) {
	items := make([]Item, 0, len(e.items))
	for _, held := range e.items {
		items = append(items, held.item)
	}
	meso := e.meso
	e.items = nil
	e.meso = 0
	return items, meso
}
