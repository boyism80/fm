package timers

import (
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/services/game/entity"
)

const CharacterSaveTimerName = "CharacterSave"

type CharacterSaveTimer struct{}

func (*CharacterSaveTimer) New() *CharacterSaveTimer {
	return &CharacterSaveTimer{}
}

func (t *CharacterSaveTimer) GetName() string {
	return CharacterSaveTimerName
}

func (t *CharacterSaveTimer) GetInterval() time.Duration {
	return 5 * time.Minute
}

func (t *CharacterSaveTimer) GetInitialDelay() time.Duration {
	return 5 * time.Minute
}

func (t *CharacterSaveTimer) Handle(ctx actor.Context, mapData *entity.Map) error {
	if mapData == nil || mapData.GameWorld == nil {
		return nil
	}
	allPlayers := mapData.GetAllPlayers()
	if len(allPlayers) == 0 {
		return nil
	}

	chars := make([]*entity.Character, 0, len(allPlayers))
	for _, obj := range allPlayers {
		if ch, ok := obj.(*entity.Character); ok {
			chars = append(chars, ch)
		}
	}

	mapData.GameWorld.SaveAsync(ctx, chars)
	return nil
}
