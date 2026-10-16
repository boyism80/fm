package game

import (
	"time"

	"github.com/boyism80/fm/maple/wz"
)

type ReactorSpawn struct {
	id           uint32
	wz           *wz.ReactorSpawn
	template     *wz.Reactor
	state        byte
	spawned      bool
	activeOID    uint32
	respawnTimer *time.Timer
}

func (rs *ReactorSpawn) RespawnDelay() time.Duration {
	if rs == nil || rs.wz == nil {
		return 0
	}
	return rs.wz.RespawnDelay
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
	if delay <= 0 || mapInstance.reactorGenBlocked {
		return
	}

	rs.CancelRespawnTimer()

	pid := mapInstance.LogicActorPID()
	if pid == nil || mapInstance.gameWorld == nil {
		return
	}

	spawnID := rs.id
	gameWorld := mapInstance.gameWorld
	rs.respawnTimer = time.AfterFunc(delay, func() {
		if gameWorld != nil {
			gameWorld.Systems().Get[*SchedulerSystem]().RunReactorRespawn(pid, mapInstance.GetMapID(), spawnID)
		}
	})
}
