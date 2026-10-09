package entity

import (
	"fmt"
	"github.com/boyism80/fm/core/clock"
	"time"

	"github.com/boyism80/fm/services/game/constant"
	"github.com/boyism80/fm/services/game/wz"
	"github.com/boyism80/fm/types"
)

const summonedMobRespawn = 4 * time.Second

func (m *Map) Respawn(includeNegativeMobTime bool) int {
	if m == nil {
		return 0
	}

	now := clock.Now()
	spawned := 0
	for _, mobSpawn := range m.MobSpawns {
		if m.trySpawnMobRezen(mobSpawn, now, includeNegativeMobTime, true) {
			spawned++
		}
	}
	for _, mobSpawn := range m.SummonedMobSpawns {
		if m.trySpawnMobRezen(mobSpawn, now, includeNegativeMobTime, true) {
			spawned++
		}
	}
	return spawned
}

func (m *Map) TickMobSpawns() {
	if m == nil {
		return
	}

	now := clock.Now()
	for _, mobSpawn := range m.MobSpawns {
		m.trySpawnMobRezen(mobSpawn, now, false, false)
	}
	for _, mobSpawn := range m.SummonedMobSpawns {
		m.trySpawnMobRezen(mobSpawn, now, false, false)
	}
}

// Summoned spawn points stay until Reset, respawning their mob like a WZ life entry.
func (m *Map) SummonMob(mobID uint32, position types.Point[int16], team constant.CarnivalTeam) error {
	mobSpawn := &MobSpawn{
		Wz: &wz.MobSpawn{
			BaseSpawn: &wz.BaseSpawn{
				ID:       mobID,
				Position: types.Vector2[int16]{X: position.X, Y: position.Y},
				MobTime:  summonedMobRespawn,
				Team:     int(team),
			},
		},
	}
	if !m.trySpawnMobRezen(mobSpawn, clock.Now(), false, true) {
		return fmt.Errorf("failed to summon mob %d", mobID)
	}
	m.SummonedMobSpawns = append(m.SummonedMobSpawns, mobSpawn)
	return nil
}

func (m *Map) trySpawnMobRezen(mobSpawn *MobSpawn, now time.Time, includeNegativeMobTime bool, immediate bool) bool {
	if m == nil || mobSpawn == nil || mobSpawn.Spawned || mobSpawn.Wz == nil {
		return false
	}
	if m.IsMobGenBlocked(mobSpawn.Wz.ID) {
		return false
	}
	if mobSpawn.Wz.MobTime < 0 && !includeNegativeMobTime {
		return false
	}
	if !immediate && mobSpawn.Wz.MobTime > 0 {
		if now.Sub(mobSpawn.LastSpawnedAt) < mobSpawn.Wz.MobTime {
			return false
		}
	}

	position := types.Point[int16]{
		X: mobSpawn.Wz.Position.X,
		Y: mobSpawn.Wz.Position.Y,
	}

	_, err := m.SpawnMob(mobSpawn.Wz.ID, position, mobSpawn, constant.MobSpawnTypeAnimate, 0, constant.CarnivalTeam(mobSpawn.Wz.Team))
	if err != nil {
		return false
	}

	mobSpawn.Spawned = true
	mobSpawn.LastSpawnedAt = now
	return true
}

func (m *Map) SetAllMobGenEnabled(enabled bool) {
	if m == nil {
		return
	}
	m.mobGenBlockedAll = !enabled
}

func (m *Map) SetMobGenEnabled(mobWZID uint32, enabled bool) {
	if m == nil || mobWZID == 0 {
		return
	}
	if m.blockedMobGen == nil {
		m.blockedMobGen = make(map[uint32]struct{})
	}
	if enabled {
		delete(m.blockedMobGen, mobWZID)
	} else {
		m.blockedMobGen[mobWZID] = struct{}{}
	}
}

func (m *Map) IsMobGenBlocked(mobWZID uint32) bool {
	if m == nil {
		return false
	}
	if m.mobGenBlockedAll {
		return true
	}
	if mobWZID == 0 || m.blockedMobGen == nil {
		return false
	}
	_, ok := m.blockedMobGen[mobWZID]
	return ok
}

func (m *Map) ClearBlockedMobGen() {
	if m == nil {
		return
	}
	m.mobGenBlockedAll = false
	m.blockedMobGen = make(map[uint32]struct{})
}
