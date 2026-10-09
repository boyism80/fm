package timers

import (
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/core/clock"
	"github.com/boyism80/fm/services/game/entity"
)

type ItemExpireTimer struct{}

func (t *ItemExpireTimer) GetName() string {
	return "ItemExpire"
}

func (t *ItemExpireTimer) GetInterval() time.Duration {
	return 10 * time.Second
}

func (t *ItemExpireTimer) GetInitialDelay() time.Duration {
	return 1 * time.Second
}

func (t *ItemExpireTimer) Handle(ctx actor.Context, mapData *entity.Map) error {
	if mapData.GetPlayerCount() == 0 {
		return nil
	}
	now := clock.Now()
	for _, obj := range mapData.GetAllPlayers() {
		if ch, ok := obj.(*entity.Character); ok {
			ch.Inventory.ExpireItems(now)
		}
	}
	return nil
}
