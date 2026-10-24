package game

import (
	"time"

	"github.com/boyism80/fm/maple/constant"
	"github.com/boyism80/fm/maple/wz"
)

type planSlot struct {
	model      wz.Item
	count      uint16
	held       uint16
	holds      int
	expiration time.Time
}

type placement struct {
	slot  int16
	count uint16
	merge bool
}

// slotPlan is a copy of one tab; checks and real adds both choose slots through it, so they cannot disagree.
type slotPlan struct {
	limit   int16
	slots   map[int16]*planSlot
	blocked map[int16]struct{}
	// full is set when the escrow's holds no longer fit, which makes newSlotPlan refuse every add to the tab.
	full bool
}

// newSlotPlan keeps room for what the escrow holds, since a cancelled room returns it: its own empty slot,
// the stack it came from, or else another free slot. It reports false when there is no room for that.
func newSlotPlan(tab *InventoryTab, invType constant.InventoryType, escrow *Escrow) (*slotPlan, bool) {
	p := &slotPlan{
		limit:   int16(tab.slotLimit),
		slots:   make(map[int16]*planSlot, len(tab.items)),
		blocked: make(map[int16]struct{}),
	}
	for slot, item := range tab.items {
		if item == nil {
			continue
		}
		p.slots[slot] = &planSlot{model: item.GetModel(), count: item.GetCount(), expiration: item.GetExpiration()}
	}
	if escrow == nil {
		return p, true
	}
	for _, held := range escrow.items {
		if held.invType != invType {
			continue
		}
		if s := p.slots[held.slot]; s == nil {
			if _, taken := p.blocked[held.slot]; !taken {
				p.blocked[held.slot] = struct{}{}
				continue
			}
		} else if s.stacks(held.item.GetModel(), held.item.GetExpiration()) && int(s.count)+int(s.held)+int(held.item.GetCount()) <= int(s.model.GetCapacity()) {
			s.held += held.item.GetCount()
			s.holds++
			continue
		}
		slot, ok := p.free()
		if !ok {
			return p, false
		}
		p.blocked[slot] = struct{}{}
	}
	return p, true
}

func (s *planSlot) stacks(model wz.Item, expiration time.Time) bool {
	return s.model.GetID() == model.GetID() &&
		model.GetCapacity() > 1 &&
		model.GetTraits()&constant.ItemTraitRechargeable == 0 &&
		s.expiration.Equal(expiration)
}

func (p *slotPlan) free() (int16, bool) {
	for slot := int16(1); slot <= p.limit; slot++ {
		if p.slots[slot] != nil {
			continue
		}
		if _, taken := p.blocked[slot]; taken {
			continue
		}
		return slot, true
	}
	return 0, false
}

// take removes count of the item from the lowest slots first, like removeByItemIDCountUnchecked.
// A slot it empties is gone after the real removal, so what the escrow held there is reserved as newSlotPlan
// would: the slot itself for the first hold, another free slot for each further one.
func (p *slotPlan) take(itemID uint32, count uint16) bool {
	for slot := int16(1); slot <= p.limit && count > 0; slot++ {
		s := p.slots[slot]
		if s == nil || s.model.GetID() != itemID {
			continue
		}
		n := min(s.count, count)
		s.count -= n
		count -= n
		if s.count > 0 {
			continue
		}
		delete(p.slots, slot)
		if s.holds == 0 {
			continue
		}
		p.blocked[slot] = struct{}{}
		for range s.holds - 1 {
			free, ok := p.free()
			if !ok {
				p.full = true
				break
			}
			p.blocked[free] = struct{}{}
		}
	}
	return count == 0
}

// place fills stacks with room in slot order, then free slots; a rechargeable takes one free slot whatever its count.
// It returns where the count goes and whether all of it fits.
func (p *slotPlan) place(model wz.Item, count uint16, expiration time.Time) ([]placement, bool) {
	if p.full {
		return nil, false
	}
	if model.GetTraits()&constant.ItemTraitRechargeable != 0 {
		slot, ok := p.free()
		if !ok {
			return nil, false
		}
		p.slots[slot] = &planSlot{model: model, count: count, expiration: expiration}
		return []placement{{slot: slot, count: count}}, true
	}

	capacity := max(model.GetCapacity(), 1)
	var out []placement
	for slot := int16(1); slot <= p.limit && count > 0; slot++ {
		s := p.slots[slot]
		if s == nil || !s.stacks(model, expiration) {
			continue
		}
		room := int(capacity) - int(s.count) - int(s.held)
		if room <= 0 {
			continue
		}
		n := min(uint16(room), count)
		s.count += n
		count -= n
		out = append(out, placement{slot: slot, count: n, merge: true})
	}
	for count > 0 {
		slot, ok := p.free()
		if !ok {
			break
		}
		n := min(capacity, count)
		p.slots[slot] = &planSlot{model: model, count: n, expiration: expiration}
		count -= n
		out = append(out, placement{slot: slot, count: n})
	}
	return out, count == 0
}

// slotFor picks the slot the item would go to.
func (inv *Inventory) slotFor(invType constant.InventoryType, item Item) (int16, bool) {
	tab := inv.tabs[invType]
	if tab == nil || item == nil {
		return 0, false
	}
	p, ok := newSlotPlan(tab, invType, inv.owner.escrow)
	if !ok {
		return 0, false
	}
	placements, _ := p.place(item.GetModel(), max(item.GetCount(), 1), item.GetExpiration())
	if len(placements) == 0 {
		return 0, false
	}
	return placements[0].slot, true
}
