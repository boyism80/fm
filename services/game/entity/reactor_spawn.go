package entity

import (
	"time"

	"github.com/boyism80/fm/services/game/wz"
)

type ReactorSpawn struct {
	ID           uint32
	Wz           *wz.ReactorSpawn
	Template     *wz.Reactor
	Spawned      bool
	ActiveOID    uint32
	respawnTimer *time.Timer
}

func (rs *ReactorSpawn) RespawnDelay() time.Duration {
	if rs == nil || rs.Wz == nil {
		return 0
	}
	return rs.Wz.RespawnDelay
}

func (rs *ReactorSpawn) CancelRespawnTimer() {
	if rs == nil || rs.respawnTimer == nil {
		return
	}
	rs.respawnTimer.Stop()
	rs.respawnTimer = nil
}

func (rs *ReactorSpawn) ScheduleRespawn(mapInstance *Map) {
	if rs == nil || mapInstance == nil {
		return
	}

	delay := rs.RespawnDelay()
	if delay <= 0 {
		return
	}

	rs.CancelRespawnTimer()

	pid := mapInstance.GetActorPID()
	if pid == nil || mapInstance.GameWorld == nil {
		return
	}

	spawnID := rs.ID
	gameWorld := mapInstance.GameWorld
	rs.respawnTimer = time.AfterFunc(delay, func() {
		if gameWorld != nil {
			gameWorld.GetSchedulerSystem().RunReactorRespawn(pid, mapInstance.GetMapID(), spawnID)
		}
	})
}
