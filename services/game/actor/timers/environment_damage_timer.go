package timers

import (
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/core/clock"
	"github.com/boyism80/fm/services/game"
)

type EnvironmentDamageTimer struct{}

func (t *EnvironmentDamageTimer) GetName() string {
	return "EnvironmentDamage"
}

func (t *EnvironmentDamageTimer) GetInterval() time.Duration {
	return time.Second
}

func (t *EnvironmentDamageTimer) GetInitialDelay() time.Duration {
	return time.Second
}

func (t *EnvironmentDamageTimer) Handle(_ actor.Context, mapData *game.Map) error {
	mapData.HurtByEnvironment(clock.Now())
	return nil
}
