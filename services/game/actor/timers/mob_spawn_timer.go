package timers

import (
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/services/game/constant"
	"github.com/boyism80/fm/services/game/entity"
	"github.com/boyism80/fm/types"
)

type MobSpawnTimer struct{}

func (*MobSpawnTimer) New() *MobSpawnTimer {
	return &MobSpawnTimer{}
}

func (t *MobSpawnTimer) GetName() string {
	return "MobSpawn"
}

func (t *MobSpawnTimer) GetInterval() time.Duration {
	return 8 * time.Second
}

func (t *MobSpawnTimer) GetInitialDelay() time.Duration {
	return 8 * time.Second
}

func (t *MobSpawnTimer) Handle(ctx actor.Context, mapData *entity.Map) error {
	if mapData.GetPlayerCount() == 0 {
		return nil
	}

	now := time.Now()

	for _, mobSpawn := range mapData.MobSpawns {
		if mobSpawn.Spawned {
			continue
		}

		if mobSpawn.Wz.MobTime > 0 {
			if now.Sub(mobSpawn.LastSpawnedAt) < mobSpawn.Wz.MobTime {
				continue
			}
		}

		position := types.Point[int16]{
			X: mobSpawn.Wz.Position.X,
			Y: mobSpawn.Wz.Position.Y,
		}

		_, err := mapData.SpawnMob(mobSpawn.Wz.ID, position, mobSpawn, constant.MOB_SPAWN_TYPE_ANIMATE, 0)
		if err != nil {
			continue
		}

		mobSpawn.Spawned = true
		mobSpawn.LastSpawnedAt = now
	}

	return nil
}
