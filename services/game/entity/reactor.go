package entity

import (
	"fmt"
	"log"
	"time"

	"github.com/boyism80/fm/core/luax"
	"github.com/boyism80/fm/protocol/dto"
	"github.com/boyism80/fm/protocol/response"
	"github.com/boyism80/fm/services/game/constant"
	"github.com/boyism80/fm/services/game/wz"
	"github.com/boyism80/fm/types"
)

type Reactor struct {
	ObjectCore
	Wz                 *wz.Reactor
	Spawn              *ReactorSpawn
	State              byte
	TriggerCharacterID uint32
	ActivatedItemID    uint32
}

func (r *Reactor) GetTrigger() *Character {
	if r == nil || r.Map == nil || r.TriggerCharacterID == 0 {
		return nil
	}
	return r.Map.GetPlayer(r.TriggerCharacterID)
}

func (r *Reactor) GetObjectType() constant.ObjectType {
	return constant.ObjectTypeReactor
}

func (r *Reactor) Is(typ constant.ObjectType) bool {
	return r.GetObjectType().Has(typ)
}

func (r *Reactor) ToDTO() *dto.Reactor {
	if r == nil {
		return nil
	}

	reactorID := uint32(0)
	facing := uint8(0)
	name := ""
	if r.Wz != nil {
		reactorID = r.Wz.ID
	}
	if r.Spawn != nil && r.Spawn.Wz != nil {
		facing = uint8(r.Spawn.Wz.FacingDirection)
		name = r.Spawn.Wz.Name
	}

	return &dto.Reactor{
		OID:       r.OID,
		ReactorID: reactorID,
		State:     r.State,
		Position:  r.Position,
		Facing:    facing,
		Name:      name,
	}
}

func (r *Reactor) SendSpawnSyncToViewer(viewer *Character) {
	if r == nil || viewer == nil {
		return
	}

	viewer.Send(&response.SpawnReactor{
		Reactor: r.ToDTO(),
	}, types.SEND_POLICY_ENCRYPT)
}

func (r *Reactor) SendDestroySyncToViewer(viewer *Character) {
	if r == nil || viewer == nil {
		return
	}
	viewer.Send(&response.DestroyReactor{
		Reactor: r.ToDTO(),
	}, types.SEND_POLICY_ENCRYPT)
}

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

func (r *Reactor) ForceHitState(state byte) {
	if r == nil || r.Map == nil {
		return
	}
	r.State = state
	r.Map.listener.OnReactorTriggered(r.Map, r, 0)
}

func (r *Reactor) Hit(trigger *Character, hitSide constant.ReactorHitSide, stance int32) {
	if r == nil || r.Wz == nil || r.Map == nil {
		return
	}
	r.TriggerCharacterID = 0
	if trigger != nil {
		r.TriggerCharacterID = trigger.GetID()
	}
	event := r.currentEvent()
	if event == nil {
		return
	}
	if event.Type == constant.ReactorEventTypeDirectionalHit && hitSide.IsLeft() {
		return
	}

	oldState := r.State
	newState := event.NextState
	r.State = newState

	if r.eventAt(newState) == nil {
		if r.respawnDelay() > 0 {
			_ = r.Map.RemoveReactor(r.OID, true)
			r.callScript("on_reactor")
			return
		}

		// The client plays the final collapse only when DESTROY arrives without a preceding TRIGGER.
		r.callScript("on_reactor")
		if r.Map.GetReactor(r.OID) != r {
			return
		}
		r.Map.listener.OnReactorTriggered(r.Map, r, stance)
		return
	}

	done := false
	r.Map.listener.OnReactorTriggered(r.Map, r, stance)
	r.callScript("on_state")
	newEvent := r.eventAt(newState)
	if newEvent.NextState == newState {
		r.callScript("on_reactor")
		done = true
	}
	timeout := r.StateTimeOut(newState)
	if timeout > 0 {
		if done == false {
			r.callScript("on_reactor")
		}
		r.ScheduleStateRevert(newState, oldState, timeout)
	}
	if done || timeout > 0 || trigger == nil || trigger.GetInstantKill() == false {
		return
	}
	if newEvent.Type == constant.ReactorEventTypeHit || newEvent.Type == constant.ReactorEventTypeDirectionalHit {
		r.Hit(trigger, hitSide, stance)
	}
}

// A blocked map keeps a broken reactor in its final state, as if WZ gave it no respawn delay.
func (r *Reactor) respawnDelay() time.Duration {
	if r.Spawn == nil || r.Map.reactorGenBlocked {
		return 0
	}
	return r.Spawn.RespawnDelay()
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
	thread, err := luax.NewThread(root, fmt.Sprintf("script/reactor/%d.lua", r.Wz.ID))
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
