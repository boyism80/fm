package game

import (
	"log"
	"math/rand"

	"github.com/boyism80/fm/types"
)

type reactorDropSpawn struct {
	isMeso bool
	count  int32
	item   Item
	quest  uint32
}

func (r *Reactor) DropItems() {
	if r == nil || r.wz == nil || r.onMap == nil || r.gameWorld == nil {
		return
	}

	resources := r.gameWorld.GetResources()
	if resources == nil {
		return
	}

	reactorDrops, ok := resources.ReactorDrops[r.wz.ID]
	if !ok || len(reactorDrops) == 0 {
		return
	}

	trigger := r.GetTrigger()
	mesoRate := float32(r.gameWorld.GetMesoRate())
	gmDrop := trigger != nil && trigger.ActsAsGM()
	spawns := make([]reactorDropSpawn, 0, len(reactorDrops))

	for _, entry := range reactorDrops {
		quest := resources.DropQuest(entry)
		if quest > 0 && !questInProgress(trigger, quest) {
			continue
		}
		if !gmDrop && rand.Float32() > entry.Prob {
			continue
		}

		if entry.Item == 0 {
			count := entry.Meso()
			if count == 0 {
				continue
			}
			count = int32(float32(count) * mesoRate)
			if count == 0 {
				continue
			}
			spawns = append(spawns, reactorDropSpawn{isMeso: true, count: count})
		} else {
			item, err := newDropItem(entry, r.gameWorld)
			if err != nil {
				continue
			}
			spawns = append(spawns, reactorDropSpawn{isMeso: false, count: int32(item.GetCount()), item: item, quest: quest})
		}
	}

	if len(spawns) == 0 {
		return
	}

	spawnPoint := r.position
	dropX := spawnPoint.X - int16(12*len(spawns))
	ownerID, dropType := trigger.dropOwner()

	for _, spawn := range spawns {
		destPoint := types.Point[int16]{X: dropX, Y: spawnPoint.Y}
		dropX += 25

		if spawn.isMeso {
			if _, err := r.onMap.SpawnMeso(spawn.count, destPoint, spawnPoint, ownerID, dropType, false); err != nil {
				log.Printf("Failed to spawn reactor meso drop: %v", err)
			}
		} else {
			if err := r.onMap.SpawnItem(spawn.item, ItemSpawn{
				Position: destPoint,
				From:     spawnPoint,
				Owner:    ownerID,
				DropType: dropType,
				Quest:    spawn.quest,
			}); err != nil {
				log.Printf("Failed to spawn reactor item drop: %v", err)
			}
		}
	}
}

func questInProgress(ch *Character, questID uint32) bool {
	if questID == 0 {
		return true
	}
	if ch == nil {
		return false
	}
	return ch.quests != nil && ch.quests.Get(questID).IsStarted()
}
