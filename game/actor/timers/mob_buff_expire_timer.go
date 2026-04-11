package timers

import (
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/game/entity"
)

type MobBuffExpireTimer struct{}

func (*MobBuffExpireTimer) New() *MobBuffExpireTimer {
	return &MobBuffExpireTimer{}
}

func (t *MobBuffExpireTimer) GetName() string {
	return "MobBuffExpire"
}

func (t *MobBuffExpireTimer) GetInterval() time.Duration {
	return 1 * time.Second
}

func (t *MobBuffExpireTimer) GetInitialDelay() time.Duration {
	return 1 * time.Second
}

func (t *MobBuffExpireTimer) Handle(ctx actor.Context, mapData *entity.Map) error {
	_ = ctx
	if mapData.GetPlayerCount() == 0 {
		return nil
	}
	now := time.Now()
	for _, obj := range mapData.GetMobs() {
		mob, ok := obj.(*entity.Mob)
		if !ok || mob == nil || !mob.IsAlive() {
			continue
		}
		mob.RemoveExpiredMobBuffs(now)
	}
	return nil
}
