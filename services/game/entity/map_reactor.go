package entity

import (
	"fmt"

	"github.com/boyism80/fm/protocol/response"
	"github.com/boyism80/fm/services/game/constant"
	"github.com/boyism80/fm/types"
)

func (r *Reactor) ForceHitState(state byte) {
	if r == nil || r.Map == nil {
		return
	}
	r.State = state
	r.TimerActive = false
	r.BroadcastTrigger(0)
}

func (r *Reactor) BroadcastSpawn() {
	if r == nil {
		return
	}
	r.Broadcast(&response.SpawnReactor{
		Reactor: r.ToDTO(),
	}, nil)
}

func (r *Reactor) BroadcastDestroy() {
	if r == nil {
		return
	}
	r.Broadcast(&response.DestroyReactor{
		Reactor: r.ToDTO(),
	}, nil)
}

func (r *Reactor) BroadcastTrigger(stance int32) {
	if r == nil {
		return
	}
	r.Broadcast(&response.TriggerReactor{
		Reactor: r.ToDTO(),
		Stance:  stance,
	}, nil)
}

func (r *Reactor) Hit(trigger *Character) {
	if r == nil || r.Wz == nil || r.Map == nil {
		return
	}
	event := r.currentEvent()
	if event == nil {
		return
	}
	eventType := event.Type
	if int(eventType) < 0 {
		return
	}

	oldState := r.State
	newState := event.NextState
	r.State = newState
	r.BroadcastTrigger(0)
	hook := fmt.Sprintf("on_reactor_hit_%d", r.ReactorID())
	_, _ = r.callReactorScript(hook, r)

	if timeout := r.StateTimeOut(newState); timeout > 0 {
		r.ScheduleStateRevert(newState, oldState, timeout)
	}
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
		return nil, fmt.Errorf("reactor spawn spec is nil")
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

func (m *Map) GetReactorSpawn(spawnID uint32) *ReactorSpawn {
	if m == nil || m.ReactorSpawns == nil {
		return nil
	}
	return m.ReactorSpawns[spawnID]
}
