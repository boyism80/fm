package timers

import (
	"github.com/boyism80/fm/core/clock"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/maple/constant"
	"github.com/boyism80/fm/services/game"
)

type MistExpireTimer struct{}

func (t *MistExpireTimer) GetName() string {
	return "MistExpire"
}

func (t *MistExpireTimer) GetInterval() time.Duration {
	return 300 * time.Millisecond
}

func (t *MistExpireTimer) GetInitialDelay() time.Duration {
	return 300 * time.Millisecond
}

func (t *MistExpireTimer) Handle(ctx actor.Context, mapData *game.Map) error {
	now := clock.Now()
	objects := mapData.GetObjects(constant.ObjectTypeMist)
	for _, obj := range objects {
		mist, ok := obj.(*game.Mist)
		if !ok || mist == nil || mist.GetOID() == 0 {
			continue
		}
		if mist.IsExpired(now) {
			mapData.RemoveMist(mist.GetOID())
		}
	}
	return nil
}
