package game

import (
	"encoding/json"
	"log"
	"math"

	"github.com/asynkron/protoactor-go/actor"
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
	var left []json.RawMessage
	var stayed []escrowItem
	for _, held := range e.items {
		if e.putBack(held) {
			continue
		}
		raw, err := json.Marshal(held.item.BuildRPC(ch.GetID(), 0))
		if err != nil {
			stayed = append(stayed, held)
			continue
		}
		left = append(left, raw)
	}
	e.items = stayed

	returnMeso := int32(0)
	if e.meso > 0 {
		if ch.inventory.meso > math.MaxInt32-e.meso {
			returnMeso = e.meso
		} else {
			ch.inventory.addMesoUnchecked(e.meso)
		}
		e.meso = 0
	}
	if len(left) == 0 && returnMeso == 0 {
		return
	}

	requestID := ch.unsettled.Add(internal.UnsettledType_UNSETTLED_TYPE_TRADE_RESTORE, 0, tradeRestoreDetail{Items: left, Meso: returnMeso})
	ch.Post(func(actx actor.Context) {
		ch.escrow.settleRestore(actx, ch.unsettled.Find(requestID))
	})
}

func (e *Escrow) putBack(held escrowItem) bool {
	ch := e.owner
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
			return false
		}
		inven.items[int16(slot)] = held.item
		ch.listener.OnInventorySlotAdded(ch, held.invType, int16(slot), held.item)
	}
	return true
}

func (e *Escrow) settleRestore(actx actor.Context, request *UnsettledRequest) {
	if request == nil {
		return
	}
	var detail tradeRestoreDetail
	if err := json.Unmarshal([]byte(request.detail), &detail); err != nil {
		request.inFlight = false
		log.Printf("Escrow.settleRestore character=%d id=%d: %v", e.owner.GetID(), request.id, err)
		return
	}

	parcels := make([]*internal.Parcel, 0, len(detail.Items)+1)
	meso := detail.Meso
	for _, raw := range detail.Items {
		var pb internal.Inventory
		if err := json.Unmarshal(raw, &pb); err != nil {
			request.inFlight = false
			log.Printf("Escrow.settleRestore character=%d id=%d: %v", e.owner.GetID(), request.id, err)
			return
		}
		parcels = append(parcels, systemParcel(&pb, meso))
		meso = 0
	}
	if meso > 0 {
		parcels = append(parcels, systemParcel(nil, meso))
	}
	e.owner.duey.Return(actx, request, parcels)
}

type tradeRestoreDetail struct {
	Items []json.RawMessage `json:"items"`
	Meso  int32             `json:"meso,omitempty"`
}

func (e *Escrow) release() ([]Item, int32) {
	items := make([]Item, 0, len(e.items))
	for _, held := range e.items {
		items = append(items, held.item)
	}
	e.items = nil
	meso := e.meso
	e.meso = 0
	return items, meso
}
