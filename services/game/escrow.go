package game

import (
	"log"

	"github.com/boyism80/fm/maple/constant"
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
				log.Printf("Escrow.restore character=%d: no slot for item %d x%d", ch.GetID(), model.GetID(), held.item.GetCount())
				continue
			}
			inven.items[int16(slot)] = held.item
			ch.listener.OnInventorySlotAdded(ch, held.invType, int16(slot), held.item)
		}
	}
	e.items = nil

	if e.meso > 0 {
		ch.inventory.addMesoUnchecked(e.meso)
		e.meso = 0
	}
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
