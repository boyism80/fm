package timers

import (
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/services/game/constant"
	"github.com/boyism80/fm/services/game/entity"
)

type ItemCleanupTimer struct{}

func (*ItemCleanupTimer) New() *ItemCleanupTimer {
	return &ItemCleanupTimer{}
}

func (t *ItemCleanupTimer) GetName() string {
	return "ItemCleanup"
}

func (t *ItemCleanupTimer) GetInterval() time.Duration {
	return 1 * time.Second
}

func (t *ItemCleanupTimer) GetInitialDelay() time.Duration {
	return 1 * time.Second
}

func (t *ItemCleanupTimer) Handle(ctx actor.Context, mapData *entity.Map) error {
	if mapData.GetPlayerCount() == 0 {
		return nil
	}

	now := time.Now()
	items := mapData.GetItems()

	for itemID, item := range items {
		dropable, ok := item.(entity.Dropable)
		if !ok {
			continue
		}

		drop := dropable.GetDrop()
		if drop == nil {
			continue
		}

		if drop.ShouldFFA(now) {
			drop.DropType = constant.DROP_TYPE_FFA
			drop.Owner = 0
		}

		if drop.ShouldExpire(now) {
			mapData.RemoveItem(itemID, constant.REMOVE_ITEM_TYPE_EXPIRED, 0)
		}
	}

	return nil
}
