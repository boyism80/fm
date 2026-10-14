package timers

import (
	"github.com/boyism80/fm/core/clock"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/services/game/domain"
)

type CooldownCheckTimer struct{}

func (t *CooldownCheckTimer) GetName() string {
	return "CooldownCheck"
}

func (t *CooldownCheckTimer) GetInterval() time.Duration {
	return 1 * time.Second
}

func (t *CooldownCheckTimer) GetInitialDelay() time.Duration {
	return 1 * time.Second
}

func (t *CooldownCheckTimer) Handle(ctx actor.Context, mapData *domain.Map) error {
	_ = ctx
	if mapData.GetPlayerCount() == 0 {
		return nil
	}
	now := clock.Now()
	players := mapData.GetAllPlayers()
	for _, obj := range players {
		ch, ok := obj.(*domain.Character)
		if !ok {
			continue
		}
		ch.Skills.ForEach(func(_ uint32, entry *domain.SkillEntry) {
			if entry.CooldownEnd == nil {
				return
			}

			if now.Before(*entry.CooldownEnd) {
				return
			}

			entry.ClearCooldown()
		})
	}
	return nil
}
