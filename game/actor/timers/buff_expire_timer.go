package timers

import (
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/game/entity"
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
	now := time.Now()
	players := mapData.GetAllPlayers()

	for _, obj := range players {
		ch, ok := obj.(*entity.Character)
		if !ok {
			continue
		}

		entities := ch.Buffs.Entities()
		for _, buff := range entities {
			if buff == nil || buff.Wz == nil {
				continue
			}
			levelData := buff.Wz.GetLevelData(int(buff.Level))
			if levelData == nil || levelData.Time <= 0 {
				continue
			}
			if now.After(buff.StartTime.Add(levelData.Time)) {
				ch.Buffs.RemoveBuff(buff.Flags)
			}
		}
	}

	return nil
}
