package entity

import (
	"time"

	"github.com/boyism80/fm/services/game/constant"
	"github.com/boyism80/fm/types"
)

func (r *Reactor) ReactItemID() int {
	event := r.currentEvent()
	if event == nil || len(event.Items) == 0 {
		return 0
	}
	return event.Items[0].ID
}

func (r *Reactor) ReactItemQuantity() int {
	event := r.currentEvent()
	if event == nil || len(event.Items) == 0 {
		return 1
	}
	return event.Items[0].Quantity
}

func (r *Reactor) ContainsPoint(point types.Point[int16]) bool {
	event := r.currentEvent()
	if event == nil || !event.HasLT || !event.HasRB {
		return false
	}

	left := r.Position.X + int16(event.LT.X)
	top := r.Position.Y + int16(event.LT.Y)
	right := r.Position.X + int16(event.RB.X)
	bottom := r.Position.Y + int16(event.RB.Y)
	if left > right {
		left, right = right, left
	}
	if top > bottom {
		top, bottom = bottom, top
	}

	return point.X >= left && point.X <= right && point.Y >= top && point.Y <= bottom
}

func (m *Map) activateItemReactors(item Item, owner *Character) {
	if m == nil || item == nil {
		return
	}
	fp := item.GetFieldPlacement()
	if fp == nil {
		return
	}
	if m.Wz != nil && m.Wz.Everlast && fp.PlayerDrop {
		return
	}
	for _, obj := range m.GetReactors() {
		reactor, ok := obj.(*Reactor)
		if !ok || reactor == nil {
			continue
		}
		if reactor.Activate(item, owner) {
			break
		}
	}
}

func (r *Reactor) Activate(item Item, owner *Character) bool {
	if r == nil || item == nil || r.Wz == nil {
		return false
	}
	if r.EventType() != constant.ReactorEventTypeItem {
		return false
	}
	fp := item.GetFieldPlacement()
	if fp == nil {
		return false
	}
	if r.ContainsPoint(fp.Position) == false {
		return false
	}
	if r.itemActivationPending() {
		return false
	}
	if r.matchesItemDrop(item) == false {
		return false
	}

	delay := constant.ItemReactorActivateDelay
	if owner != nil && owner.GM.InstantKill {
		delay = time.Millisecond
	}

	itemOID := fp.OID
	itemID := item.GetModel().GetID()
	return r.AddTimer(reactorTimerItemActivateKey, delay, false, func() {
		mapInstance := r.GetMap()
		if mapInstance == nil || mapInstance.GetItem(itemOID) != item {
			return
		}
		_ = mapInstance.RemoveItem(itemOID, constant.RemoveItemTypeExpired, 0)
		r.ActivatedItemID = itemID
		r.Hit(owner, constant.ReactorHitAirLeft, 0)
		r.ScheduleResetState(r.respawnDelay())
	})
}

func (r *Reactor) matchesItemDrop(item Item) bool {
	event := r.currentEvent()
	model := item.GetModel()
	if event == nil || model == nil {
		return false
	}
	for _, reactItem := range event.Items {
		if model.GetID() == uint32(reactItem.ID) && int(item.GetCount()) == reactItem.Quantity {
			return true
		}
	}
	return false
}

func (r *Reactor) itemActivationPending() bool {
	return r.GetTimerEntry(reactorTimerItemActivateKey) != nil
}
