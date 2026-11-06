package timers

import (
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/services/game"
)

// MobMpRecoveryTimer runs on maps without players too, as the original server did.
type MobMpRecoveryTimer struct{}

func (t *MobMpRecoveryTimer) GetName() string {
	return "MobMpRecovery"
}

func (t *MobMpRecoveryTimer) GetInterval() time.Duration {
	return time.Second
}

func (t *MobMpRecoveryTimer) GetInitialDelay() time.Duration {
	return time.Second
}

func (t *MobMpRecoveryTimer) Handle(_ actor.Context, mapData *game.Map) error {
	for _, obj := range mapData.GetMobs() {
		mob, ok := obj.(*game.Mob)
		if !ok || mob == nil || !mob.IsAlive() {
			continue
		}
		mob.RecoverMp()
	}
	return nil
}
