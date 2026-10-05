package timers

import (
	"time"

	"github.com/asynkron/protoactor-go/actor"
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
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

	worldID := mapData.GameWorld.GetWorldID()
	entries := make([]*internal.CharacterSaveEntry, 0, len(allPlayers))
	for _, obj := range allPlayers {
		ch, ok := obj.(*entity.Character)
		if !ok || ch == nil || ch.LoggedOut() {
			continue
		}
		entry := ch.ToProto(worldID)
		if entry == nil {
			continue
		}
		entries = append(entries, entry)
	}

	if len(entries) > 0 {
		mapData.GameWorld.SaveAsync(ctx, entries)
	}
	return nil
}
