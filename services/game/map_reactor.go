package game

import (
	"fmt"

	"github.com/boyism80/fm/maple/constant"
	"github.com/boyism80/fm/maple/wz"
	"github.com/boyism80/fm/types"
)

func (m *Map) initReactors() {
	m.reactorSpawns = make(map[uint32]*ReactorSpawn)
	if m.wz == nil {
		return
	}

	for spawnID, wzSpawn := range m.wz.ReactorSpawns {
		rs := &ReactorSpawn{
			id: spawnID,
			wz: &wzSpawn,
		}
		if m.gameWorld != nil {
			rs.template = m.gameWorld.GetResources().GetReactor(wzSpawn.ReactorID)
		}
		m.reactorSpawns[spawnID] = rs
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
	if reactorSpawn == nil || reactorSpawn.wz == nil {
		return nil, fmt.Errorf("reactor spawn wz is nil")
	}
	if reactorSpawn.spawned {
		return nil, fmt.Errorf("reactor spawn slot %d already active", reactorSpawn.id)
	}

	oid := m.allocateOID()
	reactor := &Reactor{
		objectCore: ObjectCore{
			oid: oid,
			position: types.Point[int16]{
				X: reactorSpawn.wz.Position.X,
				Y: reactorSpawn.wz.Position.Y,
			},
			gameWorld: m.gameWorld,
			onMap:     m,
		},
		wz:    reactorSpawn.template,
		spawn: reactorSpawn,
		state: reactorSpawn.state,
	}
	reactor.objectCore.self = reactor
	reactor.initTimers()

	if m.objects[constant.ObjectTypeReactor] == nil {
		m.objects[constant.ObjectTypeReactor] = make(map[uint32]Object)
	}

	m.objects[constant.ObjectTypeReactor][oid] = reactor
	m.sections.add(reactor)

	reactorSpawn.CancelRespawnTimer()
	reactorSpawn.spawned = true
	reactorSpawn.activeOID = oid

	m.listener.OnReactorSpawned(m, reactor)

	return reactor, nil
}

func (m *Map) SpawnReactorTemplate(reactorID uint32, pos types.Point[int16], name string, state byte) (*Reactor, error) {
	if m == nil || m.gameWorld == nil {
		return nil, fmt.Errorf("map is nil")
	}
	template := m.gameWorld.GetResources().GetReactor(reactorID)
	if template == nil {
		return nil, fmt.Errorf("reactor %d not found", reactorID)
	}
	wzSpawn := &wz.ReactorSpawn{
		ReactorID: reactorID,
		Position:  types.Vector2[int16]{X: pos.X, Y: pos.Y},
		Name:      name,
	}
	rs := &ReactorSpawn{
		id:       reactorID,
		wz:       wzSpawn,
		template: template,
		state:    state,
	}
	return m.SpawnReactor(rs)
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

	reactorSpawn := reactor.spawn
	m.sections.remove(reactor)
	delete(m.objects[constant.ObjectTypeReactor], oid)
	m.releaseOID(oid)

	if reactorSpawn != nil {
		reactorSpawn.spawned = false
		reactorSpawn.activeOID = 0
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

	runtimeOIDs := make([]uint32, 0)
	for oid, obj := range m.GetReactors() {
		reactor, ok := obj.(*Reactor)
		if ok && reactor.spawn != nil && m.reactorSpawns[reactor.spawn.id] == reactor.spawn {
			continue
		}
		runtimeOIDs = append(runtimeOIDs, oid)
	}
	for _, oid := range runtimeOIDs {
		_ = m.RemoveReactor(oid, false)
	}

	count := 0
	for _, reactorSpawn := range m.reactorSpawns {
		if reactorSpawn == nil {
			continue
		}
		if reactorSpawn.spawned && reactorSpawn.activeOID != 0 {
			_ = m.RemoveReactor(reactorSpawn.activeOID, false)
		} else {
			reactorSpawn.CancelRespawnTimer()
		}
		reactorSpawn.spawned = false
		reactorSpawn.activeOID = 0
		if _, err := m.SpawnReactor(reactorSpawn); err == nil {
			count++
		}
	}
	return count
}

func (m *Map) FindReactor(reactorID uint32) *Reactor {
	if m == nil {
		return nil
	}
	for _, obj := range m.GetReactors() {
		reactor, ok := obj.(*Reactor)
		if !ok || reactor == nil || reactor.wz == nil {
			continue
		}
		if reactor.wz.ID == reactorID {
			return reactor
		}
	}
	return nil
}

func (m *Map) FindReactorName(name string) *Reactor {
	if m == nil || name == "" {
		return nil
	}
	for _, obj := range m.GetReactors() {
		reactor, ok := obj.(*Reactor)
		if !ok || reactor == nil || reactor.spawn == nil || reactor.spawn.wz == nil {
			continue
		}
		if reactor.spawn.wz.Name == name {
			return reactor
		}
	}
	return nil
}

func (m *Map) SetReactorGenEnabled(enabled bool) {
	if m == nil {
		return
	}
	m.reactorGenBlocked = !enabled
}

func (m *Map) GetReactorSpawn(spawnID uint32) *ReactorSpawn {
	if m == nil || m.reactorSpawns == nil {
		return nil
	}
	return m.reactorSpawns[spawnID]
}
