package timers

import (
	"github.com/boyism80/fm/core/clock"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/services/game/entity"
)

type BuffExpireTimer struct{}

func (*BuffExpireTimer) New() *BuffExpireTimer {
	return &BuffExpireTimer{}
}

func (t *BuffExpireTimer) GetName() string {
	return "BuffExpire"
}

func (t *BuffExpireTimer) GetInterval() time.Duration {
	return 1 * time.Second
}

func (t *BuffExpireTimer) GetInitialDelay() time.Duration {
	return 1 * time.Second
}

func (t *BuffExpireTimer) Handle(ctx actor.Context, mapData *entity.Map) error {
	_ = ctx
	if mapData.GetPlayerCount() == 0 {
		return nil
	}
	now := clock.Now()
	players := mapData.GetAllPlayers()

	for _, obj := range players {
		ch, ok := obj.(*entity.Character)
		if !ok {
			continue
		}

		entities := ch.Buffs.Entities()
		for _, buff := range entities {
			if buff == nil {
				continue
			}

			switch b := buff.(type) {
			case *entity.SkillBuff:
				if b.Duration <= 0 {
					continue
				}
				if now.After(b.StartTime.Add(b.Duration)) {
					ch.Buffs.RemoveBuff(b.GetFlags())
				}
			case *entity.ItemBuff:
				if b.Duration <= 0 {
					continue
				}
				if now.After(b.StartTime.Add(b.Duration)) {
					ch.Buffs.RemoveBuff(b.GetFlags())
				}
			}
		}
	}

	return nil
}
