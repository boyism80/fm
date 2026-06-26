package timers

import (
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/services/game/entity"
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

	mapData.TickMobSpawns()

	return nil
}
