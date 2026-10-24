package game

import (
	"github.com/boyism80/fm/maple/constant"
	"github.com/boyism80/fm/maple/wz"
)

func (inv *Inventory) Merge(inventoryType constant.InventoryType) {
	inven := inv.tabs[inventoryType]
	if inven == nil {
		return
	}
	type stack struct {
		model      wz.Item
		expiration int64
	}
	buckets := map[stack]map[int16]Item{}

	for i := range inven.slotLimit {
		item := inven.items[int16(i+1)]
		if item == nil || !inven.Stackable(item, item) {
			continue
		}
		key := stack{model: item.GetModel(), expiration: item.GetExpiration().UnixNano()}
		if buckets[key] == nil {
			buckets[key] = map[int16]Item{}
		}
		buckets[key][int16(i+1)] = item
	}

	for key, bucket := range buckets {
		model := key.model
		count := 0
		for _, v := range bucket {
			count += int(v.GetCount())
		}
		capacity := int(model.GetCapacity())
		for slot, item := range bucket {
			value := uint16(min(capacity, count))
			if item.GetCount() != value {
				item.SetCount(value)
				if value == 0 {
					inv.owner.listener.OnRemoveInventorySlot(inv.owner, inventoryType, slot)
					delete(inven.items, slot)
				} else {
					inv.owner.listener.OnInventorySlotUpdated(inv.owner, inventoryType, slot, item)
				}
			}
			count -= int(value)
		}
	}
}

func (inv *Inventory) Sort(inventoryType constant.InventoryType) {
	inven := inv.tabs[inventoryType]
	if inven == nil {
		return
	}
	n := inven.slotLimit
	buffer := make([]Item, n)
	for i := range n {
		buffer[i] = inven.items[int16(i+1)]
	}

	less := func(item1, item2 Item) bool {
		if item1 == nil && item2 == nil {
			return false
		}
		if item1 == nil {
			return false
		}
		if item2 == nil {
			return true
		}
		id1, id2 := item1.GetModel().GetID(), item2.GetModel().GetID()
		if id1 != id2 {
			return id1 < id2
		}
		return item1.GetCount() > item2.GetCount()
	}

	partition := func(low, high int) int {
		pivot := buffer[(low+high)/2]
		i1, i2 := low, high
		for i1 <= i2 {
			for less(buffer[i1], pivot) {
				i1++
			}
			for less(pivot, buffer[i2]) {
				i2--
			}
			if i1 <= i2 {
				buffer[i1], buffer[i2] = buffer[i2], buffer[i1]
				inv.owner.listener.OnSwapInventorySlot(inv.owner, inventoryType, int16(i1+1), int16(i2+1), 0)
				i1++
				i2--
			}
		}
		return i1
	}

	var qsort func(low, high int)
	qsort = func(low, high int) {
		if low < high {
			p := partition(low, high)
			qsort(low, p-1)
			qsort(p, high)
		}
	}

	qsort(0, int(n-1))

	inven.items = map[int16]Item{}
	for i := range buffer {
		if buffer[i] != nil {
			inven.items[int16(i+1)] = buffer[i]
		}
	}
}
