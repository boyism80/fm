package timers

import (
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/game/entity"
)

type CooldownCheckTimer struct{}

func (*CooldownCheckTimer) New() *CooldownCheckTimer {
	return &CooldownCheckTimer{}
}

func (t *CooldownCheckTimer) GetName() string {
	return "CooldownCheck"
}

func (t *CooldownCheckTimer) GetInterval() time.Duration {
	return 1 * time.Second
}

func (t *CooldownCheckTimer) GetInitialDelay() time.Duration {
	return 1 * time.Second
}

func (t *CooldownCheckTimer) Handle(ctx actor.Context, mapData *entity.Map) error {
	_ = ctx
	now := time.Now()
	players := mapData.GetAllPlayers()
	for _, obj := range players {
		ch, ok := obj.(*entity.Character)
		if !ok {
			continue
		}
		skills := ch.Skills
		if skills == nil {
			continue
		}
		for _, entry := range skills {
			if entry == nil {
				continue
			}

			if entry.CooldownEnd == nil {
				continue
			}

			if now.Before(*entry.CooldownEnd) {
				continue
			}

			entry.ClearCooldown()
		}
	}
	return nil
}
