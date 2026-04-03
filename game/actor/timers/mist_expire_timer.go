package timers

import (
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/game/constant"
	"github.com/boyism80/fm/game/entity"
)

type MistExpireTimer struct{}

func (*MistExpireTimer) New() *MistExpireTimer {
	return &MistExpireTimer{}
}

func (t *MistExpireTimer) GetName() string {
	return "MistExpire"
}

func (t *MistExpireTimer) GetInterval() time.Duration {
	return 300 * time.Millisecond
}

func (t *MistExpireTimer) GetInitialDelay() time.Duration {
	return 300 * time.Millisecond
}

func (t *MistExpireTimer) Handle(ctx actor.Context, mapData *entity.Map) error {
	now := time.Now()
	objects := mapData.GetObjects(constant.ObjectTypeMist, nil)
	for _, obj := range objects {
		mist, ok := obj.(*entity.Mist)
		if !ok || mist == nil || mist.OID == 0 {
			continue
		}
		if mist.IsExpired(now) {
			mapData.RemoveMist(mist.OID)
		}
	}
	return nil
}
