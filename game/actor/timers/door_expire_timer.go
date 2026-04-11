package timers

import (
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/game/constant"
	"github.com/boyism80/fm/game/entity"
)

type DoorExpireTimer struct{}

func (*DoorExpireTimer) New() *DoorExpireTimer {
	return &DoorExpireTimer{}
}

func (t *DoorExpireTimer) GetName() string {
	return "DoorExpireTimer"
}

func (t *DoorExpireTimer) GetInterval() time.Duration {
	return 500 * time.Millisecond
}

func (t *DoorExpireTimer) GetInitialDelay() time.Duration {
	return 500 * time.Millisecond
}

func (t *DoorExpireTimer) Handle(ctx actor.Context, mapData *entity.Map) error {
	if mapData == nil {
		return nil
	}
	now := time.Now()
	objects := mapData.GetObjects(constant.ObjectTypeDoor)
	for _, object := range objects {
		door, ok := object.(*entity.Door)
		if !ok || door == nil {
			continue
		}
		if !door.IsExpired(now) {
			continue
		}
		mapData.RemoveDoor(door.OID, true)
	}
	return nil
}
