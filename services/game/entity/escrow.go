package entity

import (
	"log"

	"github.com/boyism80/fm/services/game/constant"
)

type escrowItem struct {
	invType constant.InventoryType
	slot    int16
	item    Item
}

type Escrow struct {
	owner *Character
	items []escrowItem
	Meso  int32
}

func (e *Escrow) hold(invType constant.InventoryType, slot int16, count uint16) (Item, error) {
	item := e.owner.Inventory.GetItem(invType, slot)
	if item == nil {
		return nil, ErrMiniRoomItemNotFound
	}
	if count == 0 || count > item.GetCount() {
		return nil, ErrMiniRoomInvalid
	}

	held := item.Clone(count)
	e.owner.Inventory.RemoveItem(invType, slot, count)
	e.items = append(e.items, escrowItem{invType: invType, slot: slot, item: held})
	return held, nil
}

func (e *Escrow) holdMeso(amount int32) error {
	if amount <= 0 || e.owner.Inventory.RemoveMeso(amount) == false {
		return ErrMiniRoomInvalid
	}
	e.Meso += amount
	return nil
}

func (e *Escrow) restore() {
	ch := e.owner
	for _, held := range e.items {
		inven := ch.Inventory.Tabs[held.invType]
		model := held.item.GetModel()
		exists := inven.Items[held.slot]
		switch {
		case exists == nil:
			inven.Items[held.slot] = held.item
			ch.Listener.OnInventorySlotAdded(ch, held.invType, held.slot, held.item)
		case exists.GetModel() == model && exists.GetCount()+held.item.GetCount() <= model.GetCapacity():
			exists.Increase(held.item.GetCount())
			ch.Listener.OnInventorySlotUpdated(ch, held.invType, held.slot, exists)
		default:
			slot, ok := inven.NextSlot()
			if ok == false {
				log.Printf("Escrow.restore character=%d: no slot for item %d x%d", ch.GetID(), model.GetID(), held.item.GetCount())
				continue
			}
			inven.Items[int16(slot)] = held.item
			ch.Listener.OnInventorySlotAdded(ch, held.invType, int16(slot), held.item)
		}
	}
	e.items = nil

	if e.Meso > 0 {
		ch.Inventory.addMesoUnchecked(e.Meso)
		e.Meso = 0
	}
}

func (e *Escrow) release() ([]Item, int32) {
	items := make([]Item, 0, len(e.items))
	for _, held := range e.items {
		items = append(items, held.item)
	}
	meso := e.Meso
	e.items = nil
	e.Meso = 0
	return items, meso
}
