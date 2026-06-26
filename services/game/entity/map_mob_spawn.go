package entity

import (
	"time"

	"github.com/boyism80/fm/services/game/constant"
	"github.com/boyism80/fm/types"
)

func (m *Map) Respawn(includeNegativeMobTime bool) int {
	if m == nil {
		return 0
	}

	now := time.Now()
	spawned := 0
	for _, mobSpawn := range m.MobSpawns {
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

	now := time.Now()
	for _, mobSpawn := range m.MobSpawns {
		m.trySpawnMobRezen(mobSpawn, now, false, false)
	}
}

func (m *Map) trySpawnMobRezen(mobSpawn *MobSpawn, now time.Time, includeNegativeMobTime bool, immediate bool) bool {
	if m == nil || mobSpawn == nil || mobSpawn.Spawned || mobSpawn.Wz == nil {
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

	_, err := m.SpawnMob(mobSpawn.Wz.ID, position, mobSpawn, constant.MobSpawnTypeAnimate, 0)
	if err != nil {
		return false
	}

	mobSpawn.Spawned = true
	mobSpawn.LastSpawnedAt = now
	return true
}
