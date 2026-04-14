package timers

import (
	"log"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/services/game/entity"
)

const CharacterSaveTimerName = "CharacterSave"
const characterSaveChunkSize = 100

type CharacterSaveTimer struct {
	saver func([]*entity.Character) error
}

func NewCharacterSaveTimer(saver func([]*entity.Character) error) *CharacterSaveTimer {
	return &CharacterSaveTimer{saver: saver}
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

func (t *CharacterSaveTimer) Handle(_ actor.Context, mapData *entity.Map) error {
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

	for i := 0; i < len(chars); i += characterSaveChunkSize {
		end := i + characterSaveChunkSize
		if end > len(chars) {
			end = len(chars)
		}
		if err := t.saver(chars[i:end]); err != nil {
			log.Printf("CharacterSave: chunk %d-%d: %v", i, end, err)
		}
	}
	return nil
}
