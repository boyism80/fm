package entity

import (
	"fmt"
	"log"

	"github.com/boyism80/fm/core/luax"
	"github.com/boyism80/fm/services/game/constant"
	"github.com/boyism80/fm/services/game/wz"
	"github.com/boyism80/fm/types"
	lua "github.com/yuin/gopher-lua"
)

func (r *Reactor) currentEvent() *wz.ReactorEvent {
	if r == nil {
		return nil
	}
	return r.eventAt(r.State)
}

func (r *Reactor) eventAt(state byte) *wz.ReactorEvent {
	if r == nil || r.Wz == nil {
		return nil
	}
	return r.Wz.States[state]
}

func (r *Reactor) EventType() constant.ReactorEventType {
	event := r.currentEvent()
	if event == nil {
		return constant.ReactorEventType(-1)
	}
	return event.Type
}

func (r *Reactor) ReactItemID() int {
	event := r.currentEvent()
	if event == nil {
		return 0
	}
	return event.ItemID
}

func (r *Reactor) ReactItemQuantity() int {
	event := r.currentEvent()
	if event == nil || event.ItemQuantity <= 0 {
		return 1
	}
	return event.ItemQuantity
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
	if !r.matchesItemDrop(item) {
		return false
	}
	if !r.ContainsPoint(fp.Position) {
		return false
	}
	if r.TimerActive {
		return false
	}

	itemOID := fp.OID
	return r.ScheduleItemActivation(constant.ItemReactorActivateDelay, func() {
		mapInstance := r.GetMap()
		if mapInstance == nil || mapInstance.GetItem(itemOID) == nil {
			return
		}
		_ = mapInstance.RemoveItem(itemOID, constant.RemoveItemTypeExpired, 0)
		r.Hit(owner, constant.ReactorHitAirLeft, 0)
		if r.Spawn != nil && r.Spawn.RespawnDelay() > 0 {
			r.ScheduleResetState(r.Spawn.RespawnDelay())
		}
	})
}

func (r *Reactor) matchesItemDrop(item Item) bool {
	reactorID := r.Wz.ID
	hook := fmt.Sprintf("on_reactor_%d", reactorID)
	result, err := r.callReactorScript(hook, false, r, item)
	if err != nil || result == nil {
		event := r.currentEvent()
		if event == nil || event.Type != constant.ReactorEventTypeItem {
			return false
		}
		model := item.GetModel()
		if model == nil {
			return false
		}
		if model.GetID() != uint32(event.ItemID) {
			return false
		}
		return int(item.GetCount()) == r.ReactItemQuantity()
	}
	return lua.LVAsBool(result)
}

func (r *Reactor) callReactorScript(hook string, yield bool, args ...interface{}) (lua.LValue, error) {
	mapInstance := r.GetMap()
	if mapInstance == nil {
		return nil, fmt.Errorf("reactor has no map")
	}
	root := mapInstance.GetLuaRoot()
	if root == nil {
		return nil, fmt.Errorf("map has no lua root")
	}
	thread, err := luax.NewThread(root, fmt.Sprintf("script/reactor/%d.lua", r.Wz.ID))
	if err != nil {
		return nil, err
	}
	luax.SetConfiguration(thread, luax.Configuration{
		ActorPID: mapInstance.LogicActorPID(),
	})
	filtered := make([]interface{}, 0, len(args))
	for _, arg := range args {
		if arg != nil {
			filtered = append(filtered, arg)
		}
	}
	callArgs := filtered
	if len(callArgs) == 0 {
		callArgs = nil
	}
	if yield && root.G != nil && root.G.CurrentThread == root.G.MainThread {
		if len(callArgs) == 0 {
			luax.CallAsync(root, thread, hook).OnError(func(err error) {
				log.Printf("reactor script %s: %v", hook, err)
			})
		} else {
			luax.CallAsync(root, thread, hook, callArgs...).OnError(func(err error) {
				log.Printf("reactor script %s: %v", hook, err)
			})
		}
		return nil, nil
	}
	var result lua.LValue
	var callErr error
	var promise = luax.CallAsync(root, thread, hook)
	if len(callArgs) > 0 {
		promise = luax.CallAsync(root, thread, hook, callArgs...)
	}
	promise.Then(func(value interface{}) (interface{}, error) {
		vals := luax.ResultValues(value)
		if len(vals) > 0 {
			result = vals[0]
		}
		return nil, nil
	}).OnError(func(err error) {
		callErr = err
	})
	return result, callErr
}
