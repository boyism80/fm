package timers

import (
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/services/game/domain"
)

type MobBuffExpireTimer struct{}

func (t *MobBuffExpireTimer) GetName() string {
	return "MobBuffExpire"
}

func (t *MobBuffExpireTimer) GetInterval() time.Duration {
	return 1 * time.Second
}

func (t *MobBuffExpireTimer) GetInitialDelay() time.Duration {
	return 1 * time.Second
}

func (t *MobBuffExpireTimer) Handle(ctx actor.Context, mapData *domain.Map) error {
	_ = ctx
	if mapData.GetPlayerCount() == 0 {
		return nil
	}
	for _, obj := range mapData.GetMobs() {
		mob, ok := obj.(*domain.Mob)
		if !ok || mob == nil || !mob.IsAlive() {
			continue
		}
		mob.Buffs.RemoveExpired()
	}
	return nil
}
