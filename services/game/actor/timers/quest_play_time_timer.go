package timers

import (
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/services/game"
)

// QuestPlayTimeTimer counts play time for the daily play time quests every 9 seconds, as the original server did.
// The cash shop is no map, so time there does not count.
type QuestPlayTimeTimer struct{}

const questPlayTimeTick = 9 * time.Second

func (t *QuestPlayTimeTimer) GetName() string {
	return "QuestPlayTime"
}

func (t *QuestPlayTimeTimer) GetInterval() time.Duration {
	return questPlayTimeTick
}

func (t *QuestPlayTimeTimer) GetInitialDelay() time.Duration {
	return questPlayTimeTick
}

func (t *QuestPlayTimeTimer) Handle(_ actor.Context, mapData *game.Map) error {
	for _, obj := range mapData.GetAllPlayers() {
		if ch, ok := obj.(*game.Character); ok {
			ch.Quests().CountPlayTime(questPlayTimeTick)
		}
	}
	return nil
}
