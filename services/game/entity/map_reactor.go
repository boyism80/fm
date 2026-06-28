package entity

import (
	"fmt"
	"time"

	"github.com/boyism80/fm/services/game/constant"
	"github.com/boyism80/fm/types"
)

func (r *Reactor) ForceHitState(state byte) {
	if r == nil || r.Map == nil {
		return
	}
	r.State = state
	r.TimerActive = false
	r.triggerReactor(0)
}

func (r *Reactor) triggerReactor(stance int32) {
	if r == nil || r.Map == nil {
		return
	}
	r.Map.listener.OnReactorTriggered(r.Map, r, stance)
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

	if r.isFinalState(newState) {
		if r.shouldDestroyOnFinal(newState) {
			_ = r.Map.RemoveReactor(r.OID, true)
		} else {
			r.triggerReactor(stance)
		}
		r.runHitScript()
		return
	}

	done := false
	r.triggerReactor(stance)
	rid := r.Wz.ID
	newEvent := r.eventAt(newState)
	if newEvent != nil && (newEvent.NextState == newState || rid == 2618000 || rid == 2309000) {
		if rid > 200011 {
			r.runHitScript()
		}
		done = true
	}
	if timeout := r.StateTimeOut(newState); timeout > 0 {
		if !done && rid > 200011 {
			r.runHitScript()
		}
		r.ScheduleStateRevert(newState, oldState, timeout)
	}
}

func (r *Reactor) isFinalState(state byte) bool {
	return r.eventAt(state) == nil
}

func (r *Reactor) shouldDestroyOnFinal(state byte) bool {
	delay := r.respawnDelay()
	if delay <= 0 {
		return false
	}
	event := r.eventAt(state)
	if event == nil {
		return true
	}
	return int(event.Type) < int(constant.ReactorEventTypeItem)
}

func (r *Reactor) respawnDelay() time.Duration {
	if r == nil || r.Spawn == nil {
		return 0
	}
	return r.Spawn.RespawnDelay()
}

func (r *Reactor) runHitScript() {
	hook := fmt.Sprintf("on_reactor_%d", r.Wz.ID)
	_, _ = r.callReactorScript(hook, true, r)
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

func (m *Map) initReactors() {
	m.ReactorSpawns = make(map[uint32]*ReactorSpawn)
	if m.Wz == nil {
		return
	}

	for spawnID, wzSpawn := range m.Wz.ReactorSpawns {
		rs := &ReactorSpawn{
			ID: spawnID,
			Wz: &wzSpawn,
		}
		if m.GameWorld != nil {
			rs.Template = m.GameWorld.GetResources().GetReactor(wzSpawn.ReactorID)
		}
		m.ReactorSpawns[spawnID] = rs
		_, err := m.SpawnReactor(rs)
		if err != nil {
			continue
		}
	}
}

func (m *Map) SpawnReactor(reactorSpawn *ReactorSpawn) (*Reactor, error) {
	if m == nil {
		return nil, fmt.Errorf("map is nil")
	}
	if reactorSpawn == nil || reactorSpawn.Wz == nil {
		return nil, fmt.Errorf("reactor spawn wz is nil")
	}
	if reactorSpawn.Spawned {
		return nil, fmt.Errorf("reactor spawn slot %d already active", reactorSpawn.ID)
	}

	oid := m.allocateOID()
	reactor := &Reactor{
		ObjectCore: ObjectCore{
			OID: oid,
			Position: types.Point[int16]{
				X: reactorSpawn.Wz.Position.X,
				Y: reactorSpawn.Wz.Position.Y,
			},
			GameWorld: m.GameWorld,
			Map:       m,
		},
		Wz:    reactorSpawn.Template,
		Spawn: reactorSpawn,
		State: 0,
	}
	reactor.ObjectCore.self = reactor
	reactor.initTimers()

	if m.objects[constant.ObjectTypeReactor] == nil {
		m.objects[constant.ObjectTypeReactor] = make(map[uint32]Object)
	}

	m.objects[constant.ObjectTypeReactor][oid] = reactor

	reactorSpawn.CancelRespawnTimer()
	reactorSpawn.Spawned = true
	reactorSpawn.ActiveOID = oid

	m.listener.OnReactorSpawned(m, reactor)

	return reactor, nil
}

func (m *Map) RemoveReactor(oid uint32, scheduleRespawn bool) error {
	if m == nil {
		return fmt.Errorf("map is nil")
	}
	if m.objects[constant.ObjectTypeReactor] == nil {
		return fmt.Errorf("no reactors on map")
	}

	obj, ok := m.objects[constant.ObjectTypeReactor][oid]
	if !ok {
		return fmt.Errorf("reactor %d not found on map", oid)
	}

	reactor, ok := obj.(*Reactor)
	if !ok || reactor == nil {
		return fmt.Errorf("reactor %d is not a reactor", oid)
	}

	reactor.ClearReactorTimers()

	reactorSpawn := reactor.Spawn
	delete(m.objects[constant.ObjectTypeReactor], oid)
	m.releaseOID(oid)

	if reactorSpawn != nil {
		reactorSpawn.Spawned = false
		reactorSpawn.ActiveOID = 0
		if scheduleRespawn {
			reactorSpawn.ScheduleRespawn(m)
		} else {
			reactorSpawn.CancelRespawnTimer()
		}
	}

	m.listener.OnReactorRemoved(m, reactor)

	return nil
}

func (m *Map) ReloadReactors() int {
	if m == nil {
		return 0
	}

	count := 0
	for _, reactorSpawn := range m.ReactorSpawns {
		if reactorSpawn == nil {
			continue
		}
		if reactorSpawn.Spawned && reactorSpawn.ActiveOID != 0 {
			_ = m.RemoveReactor(reactorSpawn.ActiveOID, false)
		} else {
			reactorSpawn.CancelRespawnTimer()
		}
		reactorSpawn.Spawned = false
		reactorSpawn.ActiveOID = 0
		if _, err := m.SpawnReactor(reactorSpawn); err == nil {
			count++
		}
	}
	return count
}

func (m *Map) ReactorByTemplate(reactorID uint32) *Reactor {
	if m == nil {
		return nil
	}
	for _, obj := range m.GetReactors() {
		reactor, ok := obj.(*Reactor)
		if !ok || reactor == nil || reactor.Wz == nil {
			continue
		}
		if reactor.Wz.ID == reactorID {
			return reactor
		}
	}
	return nil
}

func (m *Map) ReactorByName(name string) *Reactor {
	if m == nil || name == "" {
		return nil
	}
	for _, obj := range m.GetReactors() {
		reactor, ok := obj.(*Reactor)
		if !ok || reactor == nil || reactor.Spawn == nil || reactor.Spawn.Wz == nil {
			continue
		}
		if reactor.Spawn.Wz.Name == name {
			return reactor
		}
	}
	return nil
}

func (m *Map) GetReactorSpawn(spawnID uint32) *ReactorSpawn {
	if m == nil || m.ReactorSpawns == nil {
		return nil
	}
	return m.ReactorSpawns[spawnID]
}
