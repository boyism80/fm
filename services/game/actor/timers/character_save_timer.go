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

	// Build snapshots in the map actor's sync context: Map is non-nil here so
	// PersistMapID() is safe. LoggedOut and mapID==0 characters are filtered out
	// to avoid writing stale or invalid data from the async goroutine.
	snapshots := make([]*entity.CharacterSnapshot, 0, len(allPlayers))
	for _, obj := range allPlayers {
		ch, ok := obj.(*entity.Character)
		if !ok || ch == nil || ch.LoggedOut() {
			continue
		}
		mapID := ch.PersistMapID()
		if mapID == 0 {
			continue
		}
		snapshots = append(snapshots, &entity.CharacterSnapshot{Character: ch, MapID: mapID})
	}

	if len(snapshots) > 0 {
		mapData.GameWorld.SaveAsync(ctx, snapshots)
	}
	return nil
}
