package timers

import (
	"time"

	"github.com/boyism80/fm/core/clock"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/maple/constant"
	"github.com/boyism80/fm/services/game"
)

type ItemCleanupTimer struct{}

func (t *ItemCleanupTimer) GetName() string {
	return "ItemCleanup"
}

func (t *ItemCleanupTimer) GetInterval() time.Duration {
	return 1 * time.Second
}

func (t *ItemCleanupTimer) GetInitialDelay() time.Duration {
	return 1 * time.Second
}

func (t *ItemCleanupTimer) Handle(ctx actor.Context, mapData *game.Map) error {
	if mapData.GetPlayerCount() == 0 {
		return nil
	}

	now := clock.Now()
	items := mapData.GetItems()

	for itemID, item := range items {
		placeable, ok := item.(game.FieldPlaceable)
		if !ok {
			continue
		}

		fp := placeable.GetFieldPlacement()
		if fp == nil {
			continue
		}

		if fp.ShouldFFA(now) {
			fp.DropType = constant.DropTypeFFA
			fp.Owner = 0
		}

		if fp.ShouldExpire(now) {
			mapData.RemoveItem(itemID, constant.RemoveItemTypeExpired, 0)
		}
	}

	return nil
}
