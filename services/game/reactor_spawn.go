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
	rs.respawnTimer = mapInstance.After(delay, func() {
		if rs.spawned {
			return
		}
		_, _ = mapInstance.SpawnReactor(rs)
	})
}
