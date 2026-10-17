package timers

import (
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/core/clock"
	"github.com/boyism80/fm/maple/constant"
	"github.com/boyism80/fm/services/game"
)

type MistTimer struct{}

func (t *MistTimer) GetName() string {
	return "Mist"
}

func (t *MistTimer) GetInterval() time.Duration {
	return 300 * time.Millisecond
}

func (t *MistTimer) GetInitialDelay() time.Duration {
	return 300 * time.Millisecond
}

func (t *MistTimer) Handle(ctx actor.Context, mapData *game.Map) error {
	now := clock.Now()
	for _, obj := range mapData.GetObjects(constant.ObjectTypeMist) {
		obj.(*game.Mist).Tick(ctx, now)
	}
	return nil
}
