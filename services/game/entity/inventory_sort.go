package entity

import (
	"github.com/boyism80/fm/services/game/constant"
	"github.com/boyism80/fm/services/game/wz"
)

func (ch *Character) MergeItems(inventoryType constant.InventoryType) {
	inven := ch.Inventory.Tabs[inventoryType]
	buckets := map[wz.Item]map[int16]Item{}

	for i := range inven.SlotLimit {
		item := inven.Items[int16(i+1)]
		if item == nil {
			continue
		}
		model := item.GetModel()
		if buckets[model] == nil {
			buckets[model] = map[int16]Item{}
		}
		buckets[model][int16(i+1)] = item
	}

	for model, bucket := range buckets {
		count := uint16(0)
		for _, v := range bucket {
			count += v.GetCount()
		}
		capacity := model.GetCapacity()
		for slot, item := range bucket {
			value := min(capacity, count)
			if item.GetCount() != value {
				item.SetCount(value)
				if value == 0 {
					ch.Listener.OnRemoveInventorySlot(ch, inventoryType, slot)
					delete(inven.Items, slot)
				} else {
					ch.Listener.OnUpdateInventorySlot(ch, inventoryType, slot, item)
				}
			}
			count -= value
		}
	}
}

func (ch *Character) SortItems(inventoryType constant.InventoryType) {
	inven := ch.Inventory.Tabs[inventoryType]
	n := inven.SlotLimit
	buffer := make([]Item, n)
	for i := range n {
		buffer[i] = inven.Items[int16(i+1)]
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
				ch.Listener.OnSwapInventorySlot(ch, inventoryType, int16(i1+1), int16(i2+1), 0)
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

	inven.Items = map[int16]Item{}
	for i := range buffer {
		if buffer[i] != nil {
			inven.Items[int16(i+1)] = buffer[i]
		}
	}
}
