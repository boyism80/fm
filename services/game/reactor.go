package game

import (
	"fmt"
	"log"
	"time"

	"github.com/boyism80/fm/core/luax"
	"github.com/boyism80/fm/maple/constant"
	"github.com/boyism80/fm/maple/wz"
	"github.com/boyism80/fm/protocol/dto"
	"github.com/boyism80/fm/protocol/response"
	"github.com/boyism80/fm/types"
)

type Reactor struct {
	objectCore
	wz                 *wz.Reactor
	spawn              *ReactorSpawn
	state              byte
	triggerCharacterID uint32
	activatedItemID    uint32
}

func (r *Reactor) GetTrigger() *Character {
	if r == nil || r.onMap == nil {
		return nil
	}
	return r.onMap.GetPlayer(r.triggerCharacterID)
}

func (r *Reactor) GetObjectType() constant.ObjectType {
	return constant.ObjectTypeReactor
}

func (r *Reactor) Is(typ constant.ObjectType) bool {
	return r.GetObjectType().Has(typ)
}

func (r *Reactor) BuildS2C() *dto.Reactor {
	if r == nil {
		return nil
	}

	reactorID := uint32(0)
	facing := uint8(0)
	name := ""
	if r.wz != nil {
		reactorID = r.wz.ID
	}
	if r.spawn != nil && r.spawn.wz != nil {
		facing = uint8(r.spawn.wz.FacingDirection)
		name = r.spawn.wz.Name
	}

	return &dto.Reactor{
		OID:       r.oid,
		ReactorID: reactorID,
		State:     r.state,
		Position:  r.position,
		Facing:    facing,
		Name:      name,
	}
}

func (r *Reactor) SendSpawnSyncToViewer(viewer *Character) {
	if r == nil || viewer == nil {
		return
	}

	viewer.Send(&response.SpawnReactor{
		Reactor: r.BuildS2C(),
	}, types.SEND_POLICY_ENCRYPT)
}

func (r *Reactor) SendDestroySyncToViewer(viewer *Character) {
	if r == nil || viewer == nil {
		return
	}
	viewer.Send(&response.DestroyReactor{
		Reactor: r.BuildS2C(),
	}, types.SEND_POLICY_ENCRYPT)
}

func (r *Reactor) currentEvent() *wz.ReactorEvent {
	if r == nil {
		return nil
	}
	return r.eventAt(r.state)
}

func (r *Reactor) eventAt(state byte) *wz.ReactorEvent {
	if r == nil || r.wz == nil {
		return nil
	}
	return r.wz.States[state]
}

func (r *Reactor) EventType() constant.ReactorEventType {
	event := r.currentEvent()
	if event == nil {
		return constant.ReactorEventType(-1)
	}
	return event.Type
}

func (r *Reactor) ForceHitState(state byte) {
	if r == nil || r.onMap == nil {
		return
	}
	r.state = state
	r.onMap.listener.OnReactorTriggered(r.onMap, r, 0)
}

func (r *Reactor) Hit(trigger *Character, hitSide constant.ReactorHitSide, stance int32) {
	if r == nil || r.wz == nil || r.onMap == nil {
		return
	}
	r.triggerCharacterID = 0
	if trigger != nil {
		r.triggerCharacterID = trigger.GetID()
	}
	event := r.currentEvent()
	if event == nil {
		return
	}
	if event.Type == constant.ReactorEventTypeDirectionalHit && (hitSide == constant.ReactorHitAirLeft || hitSide == constant.ReactorHitGroundLeft) {
		return
	}

	oldState := r.state
	newState := event.NextState
	r.state = newState

	if r.eventAt(newState) == nil {
		if r.respawnDelay() > 0 {
			_ = r.onMap.RemoveReactor(r.oid, true)
			r.callScript("on_reactor")
			return
		}

		// The client plays the final collapse only when DESTROY arrives without a preceding TRIGGER.
		r.callScript("on_reactor")
		if r.onMap.GetReactor(r.oid) != r {
			return
		}
		r.onMap.listener.OnReactorTriggered(r.onMap, r, stance)
		return
	}

	done := false
	r.onMap.listener.OnReactorTriggered(r.onMap, r, stance)
	r.callScript("on_state")
	newEvent := r.eventAt(newState)
	if newEvent.NextState == newState {
		r.callScript("on_reactor")
		done = true
	}
	timeout := r.StateTimeOut(newState)
	if timeout > 0 {
		if !done {
			r.callScript("on_reactor")
		}
		r.ScheduleStateRevert(newState, oldState, timeout)
	}
	if done || timeout > 0 || trigger == nil || !trigger.gm.instantKill || r.cycles(newState) {
		return
	}
	if newEvent.Type == constant.ReactorEventTypeHit || newEvent.Type == constant.ReactorEventTypeDirectionalHit {
		r.Hit(trigger, hitSide, stance)
	}
}

// cycles reports whether hits from state come back around without breaking, as a lever does.
func (r *Reactor) cycles(state byte) bool {
	seen := map[byte]bool{}
	for event := r.eventAt(state); event != nil && event.NextState != state; event = r.eventAt(state) {
		if event.Type != constant.ReactorEventTypeHit && event.Type != constant.ReactorEventTypeDirectionalHit {
			return false
		}
		if seen[state] {
			return true
		}
		seen[state] = true
		state = event.NextState
	}
	return false
}

// A blocked map keeps a broken reactor in its final state, as if WZ gave it no respawn delay.
func (r *Reactor) respawnDelay() time.Duration {
	if r.spawn == nil || r.onMap.reactorGenBlocked {
		return 0
	}
	return r.spawn.RespawnDelay()
}

func (r *Reactor) callScript(hook string) {
	mapInstance := r.GetMap()
	if mapInstance == nil {
		return
	}
	root := mapInstance.GetLuaRoot()
	if root == nil {
		return
	}
	thread, err := luax.NewThread(root, fmt.Sprintf("script/reactor/%d.lua", r.wz.ID))
	if err != nil {
		log.Printf("reactor script %s: %v", hook, err)
		return
	}
	luax.SetConfiguration(thread, luax.Configuration{
		ActorPID: mapInstance.LogicActorPID(),
	})
	if root.G != nil && root.G.CurrentThread == root.G.MainThread {
		luax.CallAsync(nil, root, thread, hook, r).OnError(func(err error) {
			log.Printf("reactor script %s: %v", hook, err)
		})
		return
	}
	if _, err := luax.Call(thread, hook, r); err != nil {
		log.Printf("reactor script %s: %v", hook, err)
	}
}
